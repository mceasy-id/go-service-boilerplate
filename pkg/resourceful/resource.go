package resourceful

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
)

type Definition struct {
	tableDefinition        *Table
	fieldsMap              map[string]*Field
	searchFields           []*Field
	defaultWhereStatements []string
	defaultSortFields      []Field
	defaultUsedTableMap    map[*Table]map[string]bool
}

// Create a new Definition
func NewDefinition(tableDefinition *Table) (*Definition, error) {
	fieldsMap, searchFields, defaultWhereStatements, defaultSortFields, defaultUsedTableMap, err := tableDefinition.init()
	if err != nil {
		return nil, err
	}

	return &Definition{
		tableDefinition:        tableDefinition,
		fieldsMap:              fieldsMap,
		searchFields:           searchFields,
		defaultWhereStatements: defaultWhereStatements,
		defaultSortFields:      defaultSortFields,
		defaultUsedTableMap:    defaultUsedTableMap,
	}, nil
}

type Resource[IDType, Model any] struct {
	// Refreshed State (on SetParam)
	Parameter        *Parameter
	isProcessed      bool
	selectStatements []string
	whereStatements  []string
	sortStatements   []string
	queryArgs        []any
	usedTablesMap    map[*Table]map[string]bool
	result           *Result[IDType, Model]

	// Persistance State
	tableDefinition        *Table
	fieldsMap              map[string]*Field
	searchFields           []*Field
	defaultWhereStatements []string
	defaultSortFields      []Field
	defaultUsedTableMap    map[*Table]map[string]bool
	isAPIResource          bool
}

type Parameter struct {
	Limit        int
	Page         int
	Search       string
	LocalFilters []string
	Filters      []string
	Sorts        []string
}

type Result[IDType, Model any] struct {
	Ids             []IDType
	PaginatedResult []Model
}

// Create a new Resource instance with definition
func NewResource[IDType, Model any](definition *Definition) *Resource[IDType, Model] {
	var (
		fieldsMap              map[string]*Field
		searchFields           []*Field
		defaultWhereStatements []string
		defaultSortFields      []Field
		defaultUsedTableMap    map[*Table]map[string]bool
	)
	fieldsMap = make(map[string]*Field)
	defaultUsedTableMap = make(map[*Table]map[string]bool)

	// Copy definition
	for fieldKey, field := range definition.fieldsMap {
		fieldsMap[fieldKey] = field
	}
	searchFields = append(searchFields, definition.searchFields...)
	defaultWhereStatements = append(defaultWhereStatements, definition.defaultWhereStatements...)
	defaultSortFields = append(defaultSortFields, definition.defaultSortFields...)
	for table, usageMap := range definition.defaultUsedTableMap {
		newUsageMap := make(map[string]bool)
		for usage := range usageMap {
			newUsageMap[usage] = true
		}

		defaultUsedTableMap[table] = newUsageMap
	}

	return &Resource[IDType, Model]{
		tableDefinition:        definition.tableDefinition,
		fieldsMap:              fieldsMap,
		searchFields:           searchFields,
		defaultWhereStatements: defaultWhereStatements,
		defaultSortFields:      defaultSortFields,
		defaultUsedTableMap:    defaultUsedTableMap,
	}
}
func NewAPIResource[IDType, Model comparable](param Parameter) *Resource[IDType, Model] {
	return &Resource[IDType, Model]{
		isAPIResource: true,
		Parameter:     &param,
		isProcessed:   true,
	}
}

// Set the limit, page, search, filters, localFilters, and sort param to the Resource
func (r *Resource[IDType, Model]) SetParam(param Parameter) error {
	r.cleanState()

	if r.isAPIResource {
		return nil
	}

	r.usedTablesMap = make(map[*Table]map[string]bool)
	for key, value := range r.defaultUsedTableMap {
		newUsageMap := make(map[string]bool)
		for usageKey, usageValue := range value {
			newUsageMap[usageKey] = usageValue
		}

		r.usedTablesMap[key] = newUsageMap
	}

	r.Parameter = &param

	var validationError ValidationErrors

	err := r.processSearch(param.Search)
	if err != nil {
		validationErr, _ := err.(ValidationErrors)
		validationError = append(validationError, validationErr...)
	}

	err = r.processFilters(param.Filters)
	if err != nil {
		validationErr, _ := err.(ValidationErrors)
		validationError = append(validationError, validationErr...)
	}

	err = r.processLocalFilters(param.LocalFilters)
	if err != nil {
		validationErr, _ := err.(ValidationErrors)
		validationError = append(validationError, validationErr...)
	}

	err = r.processSorts(param.Sorts)
	if err != nil {
		validationErr, _ := err.(ValidationErrors)
		validationError = append(validationError, validationErr...)
	}

	if len(validationError) > 0 {
		return validationError
	}

	r.isProcessed = true

	return nil
}

func (r *Resource[IDType, Model]) GetParamUri() string {
	stringBuilder := new(strings.Builder)

	stringBuilder.WriteString(fmt.Sprintf("limit=%d&", r.Parameter.Limit))
	stringBuilder.WriteString(fmt.Sprintf("page=%d&", r.Parameter.Page))
	if r.Parameter.Search != "" {
		stringBuilder.WriteString(fmt.Sprintf("search=%s&", url.QueryEscape(r.Parameter.Search)))
	}
	query := stringBuilder.String()
	return "?" + query[:len(query)-1]
}

// Select the table fields
func (r *Resource[IDType, Model]) Select(fields []*Field) {
	// Reset the select state
	if r.selectStatements != nil {
		r.selectStatements = nil
	}
	r.unuseTableFor(SELECT)

	for _, field := range fields {
		if field != nil {
			r.selectStatements = append(r.selectStatements, field.statement)
			r.useFieldFor(field, SELECT)
		}
	}
}

// Get the Resource query & args
func (r *Resource[IDType, Model]) QueryAndArgs() (string, []any, error) {
	if !r.isProcessed {
		return "", nil, createError("SetParam method must be called")
	}
	if len(r.selectStatements) == 0 {
		return "", nil, createError("selected fields can't be empty. Use the Select method instead")
	}

	var queryStatement string

	// Select Statement
	queryStatement += fmt.Sprintf("SELECT %s FROM %s", strings.Join(r.selectStatements, ", "), r.tableDefinition.statement)

	// Join Statements
	joinStatements := r.getJoinStatements([]string{SELECT, SEARCH, FILTER, SORT, MANDATORY})
	if len(joinStatements) > 0 {
		queryStatement += "\n" + strings.Join(joinStatements, "\n")
	}

	// Where Statement
	r.whereStatements = append(r.whereStatements, r.defaultWhereStatements...)
	if len(r.whereStatements) > 0 {
		queryStatement += "\nWHERE " + strings.Join(r.whereStatements, "\nAND ")
	}

	// Sort Statement
	if len(r.sortStatements) > 0 {
		queryStatement += "\nORDER BY " + strings.Join(r.sortStatements, ", ")
	}

	return queryStatement, r.queryArgs, nil
}

func (r *Resource[IDType, Model]) PopulateQueryArgs(idField *Field, ids []IDType) (string, []any, error) {
	if !r.isProcessed {
		return "", nil, createError("SetParam method must be called")
	}
	if len(r.selectStatements) == 0 {
		return "", nil, createError("selected fields can't be empty. Use the Select method instead")
	}

	var (
		queryStatement    string
		whereInStatements []string
		queryArgs         []any
	)

	// Select Statement
	queryStatement += fmt.Sprintf("SELECT %s FROM %s", strings.Join(r.selectStatements, ", "), r.tableDefinition.statement)

	// Join Statements
	joinStatements := r.getJoinStatements([]string{SELECT, SORT})
	if len(joinStatements) > 0 {
		queryStatement += "\n" + strings.Join(joinStatements, "\n")
	}

	// Where Statement
	for _, id := range ids {
		queryArgs = append(queryArgs, id)
		whereInStatements = append(whereInStatements, fmt.Sprintf("$%d", len(queryArgs)))
	}
	queryStatement += fmt.Sprintf("\nWHERE %s IN (%s)", idField.statement, strings.Join(whereInStatements, ", "))

	// Sort Statement
	if len(r.sortStatements) > 0 {
		queryStatement += "\nORDER BY " + strings.Join(r.sortStatements, ", ")
	}

	return queryStatement, queryArgs, nil
}

// Helper function to return the new slice from the results based on Resource limit & page parameter
// !!! FIXED SOON
func (r *Resource[IDType, Model]) GetPaginatedResults(results []IDType) ([]IDType, error) {
	endOffset := (r.Parameter.Limit * r.Parameter.Page)
	if endOffset > len(results) {
		endOffset = len(results)
	}
	startOffset := (r.Parameter.Limit * r.Parameter.Page) - r.Parameter.Limit
	if startOffset > (endOffset - 1) {
		return nil, ErrPagination
	}

	return results[startOffset:endOffset], nil
}

// Set the result to the Resource so the ResourfulResponse function can be used later
func (r *Resource[IDType, Model]) SetResult(result Result[IDType, Model]) {
	r.result = &result
}

// Will return the Resource metadata from the result
func (r *Resource[IDType, Model]) Metadata() *Metadata {
	var metadata Metadata

	metadata.Page = r.Parameter.Page
	if r.result != nil {
		metadata.TotalCount = len(r.result.Ids)
		metadata.Count = len(r.result.PaginatedResult)
		metadata.TotalPage = int(math.Ceil(float64(metadata.TotalCount) / float64(r.Parameter.Limit)))
	}

	return &metadata
}

// Build the resrouceful response
func (r *Resource[IDType, Model]) Response() *Response[IDType, Model] {
	var response Response[IDType, Model]

	response.Metadata = r.Metadata()
	response.Data.Ids = make([]IDType, 0)
	response.Data.PaginatedResult = make([]Model, 0)

	if r.result != nil {
		response.Data.Ids = r.result.Ids
		response.Data.PaginatedResult = r.result.PaginatedResult
	}

	return &response
}

func (r *Resource[IDType, Model]) getJoinStatements(usages []string) []string {
	var joinStatements []string

	// Filter the usedTables by usages
	usedTables := make(map[*Table]bool)
	for usedTable, tableUsage := range r.usedTablesMap {
		for _, usage := range usages {
			if _, ok := tableUsage[usage]; ok {
				usedTables[usedTable] = true
			}
		}
	}

	relations := r.tableDefinition.getRequiredRelations(usedTables)

	for _, relation := range relations {
		joinStatement := fmt.Sprintf(
			"JOIN %s ON %s = %s",
			relation.Table.statement,
			relation.ForeignKeyField.statement,
			relation.ReferenceKeyField.statement,
		)

		if !relation.IsMandatory {
			// If the table usage is not only for FILTER
			if !(len(r.usedTablesMap[relation.Table]) == 1 && r.usedTablesMap[relation.Table][FILTER]) {
				joinStatement = "LEFT " + joinStatement
			}
		}
		joinStatements = append(joinStatements, joinStatement)
	}

	return joinStatements
}

func (r *Resource[IDType, Model]) processSearch(searchParam string) error {
	var validationErr ValidationErrors

	if len(r.searchFields) == 0 {
		validationErr.appendFieldError("search", "no available search fields for this entity")
		return validationErr
	}

	if searchParam != "" {

		var searchStatements []string
		for _, searchField := range r.searchFields {
			if searchField.Type == NUMERIC {
				_, err := strconv.ParseFloat(searchParam, 64)

				if err == nil {
					r.useFieldFor(searchField, SEARCH)
					r.queryArgs = append(r.queryArgs, searchParam+"%")
					searchStatements = append(searchStatements, fmt.Sprintf("%s::text like $%d", searchField.statement, len(r.queryArgs)))
				}
			} else {
				r.useFieldFor(searchField, SEARCH)
				r.queryArgs = append(r.queryArgs, "%"+strings.ToLower(searchParam)+"%")
				searchStatements = append(searchStatements, fmt.Sprintf("lower(%s) like $%d", searchField.statement, len(r.queryArgs)))
			}
		}

		searchStatement := fmt.Sprintf("(%s)", strings.Join(searchStatements, " OR "))
		r.whereStatements = append(r.whereStatements, searchStatement)
	}

	return nil
}

func (r *Resource[IDType, Model]) processFilters(filterParams []string) error {
	var validationErrors ValidationErrors

	for key, filterParam := range filterParams {
		filterParam = strings.TrimSpace(filterParam)
		filterParamSplit := strings.Split(filterParam, " ")

		// Validate Format
		if len(filterParamSplit) < 3 {
			validationErrors.appendFieldError(fmt.Sprintf("filters.%d", (key+1)), "invalid format")
			continue
		}

		// Validate Field
		filterKey := filterParamSplit[0]
		var (
			filterField *Field
			ok          bool
		)
		if filterField, ok = r.fieldsMap[filterKey]; !ok || !filterField.Filterable {
			validationErrors.appendFieldError(fmt.Sprintf("filters.%d", (key+1)), "invalid field")
			continue
		}

		// Validate Operator
		filterOperator := filterParamSplit[1]
		if !validFilterOperator(filterOperator) {
			validationErrors.appendFieldError(fmt.Sprintf("filters.%d", (key+1)), "invalid operator")
			continue
		}

		// Validate Value
		rawFilterValue := strings.Join(filterParamSplit[2:], " ")
		filterValues, ok := sanitizeFilterValueByType(rawFilterValue, filterField.Type)
		if !ok {
			validationErrors.appendFieldError(fmt.Sprintf("filters.%d", (key+1)), fmt.Sprintf("value must be a valid %s format", filterField.Type))
			continue
		}

		// Process Filter
		if len(validationErrors) != 0 {
			continue
		}

		if filterOperator == "in" {
			inStatement := make([]string, 0, len(filterValues))
			for _, filterValue := range filterValues {
				r.queryArgs = append(r.queryArgs, filterValue)
				inStatement = append(inStatement, fmt.Sprintf("$%d", len(r.queryArgs)))
			}

			whereStatement := fmt.Sprintf("%s in (%s)", filterField.statement, strings.Join(inStatement, ","))
			r.whereStatements = append(r.whereStatements, whereStatement)
			r.useFieldFor(filterField, FILTER)
			continue
		}

		r.queryArgs = append(r.queryArgs, filterValues[0])
		whereStatement := fmt.Sprintf("%s %s $%d", filterField.statement, filterOperators[filterOperator], len(r.queryArgs))
		r.whereStatements = append(r.whereStatements, whereStatement)
		r.useFieldFor(filterField, FILTER)

	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func (r *Resource[IDType, Model]) processLocalFilters(localFilterParams []string) error {
	var validationErrors ValidationErrors

	for key, localFilterParam := range localFilterParams {
		localFilterParam = strings.TrimSpace(localFilterParam)
		localFilterSplit := strings.Split(localFilterParam, " ")

		// Validate Format
		if len(localFilterSplit) < 3 {
			validationErrors.appendFieldError(fmt.Sprintf("localFilters.%d", (key+1)), "invalid format")
			continue
		}

		// Validate Field
		filterKey := localFilterSplit[0]
		var (
			filterField *Field
			ok          bool
		)
		if filterField, ok = r.fieldsMap[filterKey]; !ok || (!filterField.LocalFilterable && !filterField.Filterable) {
			validationErrors.appendFieldError(fmt.Sprintf("localFilters.%d", (key+1)), "invalid field")
			continue
		}

		// Validate Operator
		filterOperator := localFilterSplit[1]
		if _, ok := filterOperators[filterOperator]; !ok {
			validationErrors.appendFieldError(fmt.Sprintf("localFilters.%d", (key+1)), "invalid operator")
			continue
		}

		// Validate Value
		filterValue := strings.Join(localFilterSplit[2:], " ")
		//TODO:  implement IN in Local Filter
		if _, ok := sanitizeFilterValueByType(filterValue, filterField.Type); !ok {
			validationErrors.appendFieldError(fmt.Sprintf("localFilters.%d", (key+1)), fmt.Sprintf("value must be a valid %s", filterField.Type))
			continue
		}

		// Process Filter
		if len(validationErrors) != 0 {
			continue
		}

		r.queryArgs = append(r.queryArgs, filterValue)
		whereStatement := fmt.Sprintf("%s %s $%d", filterField.statement, filterOperators[filterOperator], len(r.queryArgs))
		r.whereStatements = append(r.whereStatements, whereStatement)

		r.useFieldFor(filterField, FILTER)
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func (r *Resource[IDType, Model]) processSorts(sortParams []string) error {
	var (
		validationErrors  ValidationErrors
		usedSortKey       map[string]bool
		defaultSortFields []Field
		sortStatements    []string
	)
	usedSortKey = make(map[string]bool)
	defaultSortFields = make([]Field, len(r.defaultSortFields))

	copy(defaultSortFields, r.defaultSortFields)

	for key, sortParam := range sortParams {
		sortParamSplit := strings.Split(sortParam, " ")

		// Validate Format
		if len(sortParamSplit) != 2 {
			validationErrors.appendFieldError(fmt.Sprintf("sorts.%d", (key+1)), "invalid sorts format")
			continue
		}

		// Validate Field
		sortKey := sortParamSplit[0]
		var (
			sortField *Field
			ok        bool
		)
		if sortField, ok = r.fieldsMap[sortKey]; !ok || !r.fieldsMap[sortKey].Sortable {
			validationErrors.appendFieldError(fmt.Sprintf("sorts.%d", (key+1)), "invalid sorts field")
			continue
		}
		if _, ok := usedSortKey[sortKey]; ok {
			validationErrors.appendFieldError(fmt.Sprintf("sorts.%d", (key+1)), fmt.Sprintf(`duplicate with "%s" field`, sortKey))
			continue
		}

		// Validate Operator
		sortOpeartor := sortParamSplit[1]
		if _, ok := sortOperators[sortOpeartor]; !ok {
			validationErrors.appendFieldError(fmt.Sprintf("sorts.%d", (key+1)), "invalid sorts operator")
			continue
		}

		if len(validationErrors) == 0 {
			// Check if sort key override the defaultSortFields
			for key, defaultSortField := range defaultSortFields {
				if defaultSortField == *sortField {
					defaultSortFields = append(defaultSortFields[:key], defaultSortFields[key+1:]...)
				}
			}
			sortStatement := fmt.Sprintf("%s %s", sortField.statement, sortOperators[sortOpeartor])
			sortStatements = append(sortStatements, sortStatement)

			usedSortKey[sortKey] = true
			r.useFieldFor(sortField, SORT)
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	for _, defaultSortField := range defaultSortFields {
		sortStatement := fmt.Sprintf("%s %s", defaultSortField.statement, sortOperators[defaultSortField.Sort])
		sortStatements = append(sortStatements, sortStatement)

		r.useFieldFor(&defaultSortField, SORT)
	}

	r.sortStatements = append(r.sortStatements, sortStatements...)

	return nil
}

func (r *Resource[IDType, Model]) useFieldFor(field *Field, usage string) {
	if r.usedTablesMap[field.table] == nil {
		r.usedTablesMap[field.table] = make(map[string]bool)
	}

	r.usedTablesMap[field.table][usage] = true
}

func (r *Resource[IDType, Model]) unuseTableFor(usage string) {
	for usedTable, usageMap := range r.usedTablesMap {
		if _, ok := usageMap[usage]; ok {
			delete(r.usedTablesMap[usedTable], usage)
		}

		if len(r.usedTablesMap[usedTable]) == 0 {
			delete(r.usedTablesMap, usedTable)
		}
	}
}

func (r *Resource[IDType, Model]) cleanState() {
	r.Parameter = nil
	r.isProcessed = false
	r.selectStatements = nil
	r.whereStatements = nil
	r.sortStatements = nil
	r.usedTablesMap = nil
	r.queryArgs = nil
	r.result = nil
	r.isAPIResource = false
}
