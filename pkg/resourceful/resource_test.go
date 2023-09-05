package resourceful

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	productResource *Resource[string, string]
)

func init() {
	productDefinition, _ := NewDefinition(&productTable)
	productResource = NewResource[string, string](productDefinition)
}

func TestResource_SetParam(t *testing.T) {
	t.Run("valid param", func(t *testing.T) {
		productResource.cleanState()

		var param Parameter
		param.Search = "search value"
		param.Filters = []string{"name eq product_name", "count gt 2", "product_type.name eq product_type_name"}
		param.LocalFilters = []string{"company_id eq 1"}
		param.Sorts = []string{"name asc"}

		err := productResource.SetParam(param)
		require.NoError(t, err)

		productResource.Select([]*Field{productTable.Field("id")})

		query, _, err := productResource.QueryAndArgs()
		require.NoError(t, err)
		require.NotEmpty(t, query)
	})
	t.Run("valid param IN with valid format", func(t *testing.T) {
		productResource.cleanState()

		var param Parameter
		param.Search = "search value"
		param.Filters = []string{"name eq \"product_name\"", "count gt 2", "name in (\"perabotan dan alat rumah tangga\" hobi 'elektronik' 'makanan dan minuman')"}
		param.LocalFilters = []string{"company_id eq 1"}
		param.Sorts = []string{"name asc"}

		err := productResource.SetParam(param)
		require.NoError(t, err)

		productResource.Select([]*Field{productTable.Field("id")})

		query, _, err := productResource.QueryAndArgs()
		require.NoError(t, err)
		require.NotEmpty(t, query)
	})

}
func TestResource_processSearch(t *testing.T) {
	productResource.cleanState()
	productResource.usedTablesMap = make(map[*Table]map[string]bool)
	productResource.processSearch("find something")

	require.Equal(t, "%find something%", productResource.queryArgs[0])
	require.Equal(t, `(lower(product."name") like $1 OR lower(pt."name") like $2)`, productResource.whereStatements[0])
	require.Equal(t, map[*Table]map[string]bool{&productTable: {SEARCH: true}, &productTypeTable: {SEARCH: true}}, productResource.usedTablesMap)
}

func TestResource_processFilters(t *testing.T) {
	productResource.cleanState()
	productResource.usedTablesMap = make(map[*Table]map[string]bool)

	t.Run("error if invalid format", func(t *testing.T) {
		err := productResource.processFilters([]string{"invalid_format"})
		require.Error(t, err)
	})

	t.Run("error if field doesnt exists", func(t *testing.T) {
		err := productResource.processFilters([]string{"random_field eq 21"})
		require.Error(t, err)
	})

	t.Run("error if field cant filtered", func(t *testing.T) {
		err := productResource.processFilters([]string{"is_deleted eq true"})
		require.Error(t, err)
	})

	t.Run("error if invalid operator", func(t *testing.T) {
		err := productResource.processFilters([]string{"name asdf 21"})
		require.Error(t, err)
	})

	t.Run("error if invalid value type numeric", func(t *testing.T) {
		err := productResource.processFilters([]string{"count eq a"})
		require.Error(t, err)
	})

	t.Run("no error if valid", func(t *testing.T) {
		err := productResource.processFilters([]string{"name eq search_name"})
		require.NoError(t, err)
	})
	t.Run("error if no value", func(t *testing.T) {
		err := productResource.processFilters([]string{"name eq "})
		require.Error(t, err)
	})

	t.Run("no error if valid relation field", func(t *testing.T) {
		err := productResource.processFilters([]string{"product_type.name eq search_name"})
		require.NoError(t, err)
	})

	t.Run("no error if valid value type numeric", func(t *testing.T) {
		err := productResource.processFilters([]string{"count eq 21"})
		require.NoError(t, err)
	})
	t.Run("ok when IN format is valid", func(t *testing.T) {
		err := productResource.processFilters([]string{"name in (aasd asd)"})
		require.NoError(t, err)
	})

	t.Run("ok when IN format is valid", func(t *testing.T) {
		err := productResource.processFilters([]string{"name in (aasd)"})
		require.NoError(t, err)
	})
	t.Run("error when in field is invalid", func(t *testing.T) {
		err := productResource.processFilters([]string{"name in ())"})
		require.Error(t, err)
	})

	t.Run("error when in field is empty with spaces", func(t *testing.T) {
		err := productResource.processFilters([]string{"name in (    )"})
		require.Error(t, err)
	})

	require.Equal(t,
		map[*Table]map[string]bool{
			&productTable:     {FILTER: true},
			&productTypeTable: {FILTER: true},
		},
		productResource.usedTablesMap,
	)

}

func TestResource_processLocalFilters(t *testing.T) {
	productResource.cleanState()
	productResource.usedTablesMap = make(map[*Table]map[string]bool)

	t.Run("error if invalid format", func(t *testing.T) {
		err := productResource.processLocalFilters([]string{"min_temperature eq"})
		require.Error(t, err)
	})

	t.Run("error if field doesnt exists", func(t *testing.T) {
		err := productResource.processLocalFilters([]string{"random_field eq 21"})
		require.Error(t, err)
	})

	t.Run("error if field cant filtered", func(t *testing.T) {
		err := productResource.processLocalFilters([]string{"secret_field eq 12"})
		require.Error(t, err)
	})

	t.Run("no error if local filterable", func(t *testing.T) {
		err := productResource.processLocalFilters([]string{"company_id eq 21"})
		require.NoError(t, err)
	})

	t.Run("no error if filterable", func(t *testing.T) {
		err := productResource.processLocalFilters([]string{"product_type.name eq product_type_name"})
		require.NoError(t, err)
	})

	require.Equal(t,
		map[*Table]map[string]bool{
			&productTable:     {FILTER: true},
			&productTypeTable: {FILTER: true},
		},
		productResource.usedTablesMap,
	)

}

func TestResource_sortStatement(t *testing.T) {
	productResource.defaultSortFields = []Field{*productTable.Fields[5]}

	t.Run("error if invalid format", func(t *testing.T) {
		err := productResource.processSorts([]string{"asdfj sdfk asdf"})
		require.Error(t, err)
	})

	t.Run("error if forbidden field", func(t *testing.T) {
		err := productResource.processSorts([]string{"company_id desc"})
		require.Error(t, err)
	})

	t.Run("error if has duplicate field", func(t *testing.T) {
		err := productResource.processSorts([]string{"name asc", "name desc"})
		require.Error(t, err)
	})

	t.Run("no error if valid", func(t *testing.T) {
		productResource.cleanState()
		productResource.usedTablesMap = make(map[*Table]map[string]bool)
		err := productResource.processSorts([]string{"name asc"})
		require.NoError(t, err)

		require.Equal(t, []string{`product."name" ASC`, `product."created_on" DESC`}, productResource.sortStatements)
	})

	t.Run("correct order if override the default sort", func(t *testing.T) {
		productResource.cleanState()
		productResource.usedTablesMap = make(map[*Table]map[string]bool)
		err := productResource.processSorts([]string{"created_on desc", "name asc"})
		require.NoError(t, err)

		require.Equal(t, []string{`product."created_on" DESC`, `product."name" ASC`}, productResource.sortStatements)
	})

	require.Equal(t,
		map[*Table]map[string]bool{
			&productTable: {SORT: true},
		},
		productResource.usedTablesMap,
	)
}
func TestResource_GetQueryUri(t *testing.T) {
	t.Run("positive test with search", func(t *testing.T) {
		productResource.SetParam(Parameter{Limit: 10, Page: 10, Search: "uncle bob"})
		queryParam := productResource.GetParamUri()
		require.Equal(t, `?limit=10&page=10&search=uncle+bob`, queryParam)
	})
	t.Run("positive test only required", func(t *testing.T) {
		productResource.SetParam(Parameter{Limit: 10, Page: 10})
		queryParam := productResource.GetParamUri()
		require.Equal(t, `?limit=10&page=10`, queryParam)
	})
}
