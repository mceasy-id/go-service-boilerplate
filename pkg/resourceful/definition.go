package resourceful

import (
	"fmt"
	"sort"
)

type Table struct {
	Name      string
	Alias     string
	Fields    []*Field
	Relations []*Relation

	statement string
}

type Relation struct {
	IsMandatory       bool
	ForeignKeyField   *Field
	ReferenceKeyField *Field
	Table             *Table
}

type Field struct {
	Name            string
	Alias           string
	Type            string
	Searchable      bool
	Filterable      bool
	LocalFilterable bool
	Sortable        bool
	Sort            string
	SoftDeleteField bool

	statement string
	table     *Table
}

func (table *Table) init() (map[string]*Field, []*Field, []string, []Field, map[*Table]map[string]bool, error) {
	fieldsMap := make(map[string]*Field)
	defaultUsedTablesMap := make(map[*Table]map[string]bool)

	searchFields, defaultWhereStatements, defaultSortFields, err := table.dfsInit(fieldsMap, defaultUsedTablesMap)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return fieldsMap, searchFields, defaultWhereStatements, defaultSortFields, defaultUsedTablesMap, nil
}

func (table *Table) dfsInit(fieldsMap map[string]*Field, defaultUsedTablesMap map[*Table]map[string]bool) ([]*Field, []string, []Field, error) {
	selfFieldsMap := make(map[string]bool)
	var (
		searchFields           []*Field
		defaultWhereStatements []string
		defaultSortFields      []Field
	)

	// Set Table Selector & Statement
	tableSelector := table.Name
	table.statement = table.Name
	if table.Alias != "" {
		tableSelector = table.Alias
		table.statement += " " + table.Alias
	}

	// Set Field Statement
	for _, field := range table.Fields {
		// Check self duplicate Field
		if _, ok := selfFieldsMap[field.Name]; ok {
			return nil, nil, nil, createError(fmt.Sprintf(`duplicate field name on "%s" table`, table.Name))
		}

		// Set the Field properties
		field.statement = fmt.Sprintf(`%s."%s"`, tableSelector, field.Name)
		field.table = table

		// Get the Field Key
		fieldKey := field.Name
		if field.Alias != "" {
			fieldKey = field.Alias
		}
		if _, ok := fieldsMap[fieldKey]; ok {
			fieldKey = fmt.Sprintf("%s.%s", table.Name, field.Name)
		}

		// Append the fieldsMap
		if field.Filterable || field.LocalFilterable || field.Sortable {
			fieldsMap[fieldKey] = field
		}

		// Search
		if field.Searchable {
			searchFields = append(searchFields, field)
		}

		// defaultWhereStatements (soft delete only for now)
		if field.SoftDeleteField {
			defaultWhereStatements = append(defaultWhereStatements, fmt.Sprintf("%s is false", field.statement))
			useFieldAs(defaultUsedTablesMap, field, FILTER)
		}

		// Sort
		if field.Sort != "" {
			if _, ok := sortOperators[field.Sort]; !ok {
				return nil, nil, nil, createError(fmt.Sprintf(`invalid sort operator at "%s" field`, field.Name))
			}
			defaultSortFields = append(defaultSortFields, *field)
			useFieldAs(defaultUsedTablesMap, field, SORT)
		}

		selfFieldsMap[field.Name] = true
	}

	// DFS Relations
	for _, relation := range table.Relations {
		if relation.Table == table {
			return nil, nil, nil, createError(fmt.Sprintf(`circular relations on "%s" table`, table.Name))
		}

		if relation.IsMandatory {
			useTableAs(defaultUsedTablesMap, relation.Table, MANDATORY)
		}

		relationSearchFields, relationDefaultWhereStatement, relationDefaultSortFields, err := relation.Table.dfsInit(fieldsMap, defaultUsedTablesMap)
		if err != nil {
			return nil, nil, nil, err
		}

		// Append all the relation result
		searchFields = append(searchFields, relationSearchFields...)
		defaultWhereStatements = append(defaultWhereStatements, relationDefaultWhereStatement...)
		defaultSortFields = append(defaultSortFields, relationDefaultSortFields...)
	}

	return searchFields, defaultWhereStatements, defaultSortFields, nil
}

func (table *Table) getRequiredRelations(usedTables map[*Table]bool) []*Relation {
	var requiredRelations []*Relation
	var requiredRelationsMap = make(map[*Relation]int)

	table.dfsGetRequiredRelations(nil, usedTables, requiredRelationsMap)

	for table := range requiredRelationsMap {
		requiredRelations = append(requiredRelations, table)
	}

	sort.Slice(requiredRelations, func(i, j int) bool {
		return requiredRelationsMap[requiredRelations[i]] < requiredRelationsMap[requiredRelations[j]]
	})

	return requiredRelations
}

func (table *Table) dfsGetRequiredRelations(relationPaths []*Relation, usedTables map[*Table]bool, result map[*Relation]int) {
	for _, relation := range table.Relations {
		relationPaths = append(relationPaths, relation)

		// Search if exists, append the result Map
		for usedTable := range usedTables {
			if usedTable == relation.Table {
				for level, relationPath := range relationPaths {
					result[relationPath] = level
				}
			}
		}

		relation.Table.dfsGetRequiredRelations(relationPaths, usedTables, result)
	}
}

func (table *Table) Field(name string) *Field {
	for _, field := range table.Fields {
		if field.Name == name {
			return field
		}
	}

	return nil
}

func useFieldAs(usedTableMap map[*Table]map[string]bool, field *Field, usage string) {
	useTableAs(usedTableMap, field.table, usage)
}

func useTableAs(usedTableMap map[*Table]map[string]bool, table *Table, usage string) {
	if usedTableMap[table] == nil {
		usedTableMap[table] = make(map[string]bool)
	}

	usedTableMap[table][usage] = true
}
