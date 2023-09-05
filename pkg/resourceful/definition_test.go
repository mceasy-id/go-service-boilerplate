package resourceful

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	productTypeTable Table
	productTable     Table
)

func init() {
	productTypeTable.Name = "product_type"
	productTypeTable.Alias = "pt"
	productTypeTable.Fields = []*Field{
		{Name: "id"},
		{Name: "name", Searchable: true, Filterable: true},
	}

	productTable.Name = "product"
	productTable.Fields = []*Field{
		{Name: "id"},
		{Name: "name", Type: STRING, Searchable: true, Filterable: true, Sortable: true},
		{Name: "count", Type: NUMERIC, Filterable: true},
		{Name: "product_type_id"},
		{Name: "company_id", LocalFilterable: true},
		{Name: "created_on", Sort: "desc", Sortable: true},
		{Name: "is_deleted", SoftDeleteField: true},
	}
	productTable.Relations = []*Relation{
		{
			IsMandatory:       true,
			Table:             &productTypeTable,
			ForeignKeyField:   productTable.Field("product_type_id"),
			ReferenceKeyField: productTypeTable.Field("id"),
		},
	}

}

func TestDefinition_init(t *testing.T) {
	fieldsMap, searchFields, defaultWhereStatements, defaultSortFields, defaultUsedTablesMap, err := productTable.init()
	require.NoError(t, err)

	t.Run("will_set_the_table_statement", func(t *testing.T) {
		// Parent Table
		require.Equal(t, "product", productTable.statement)

		// Related Table
		require.Equal(t, "product_type pt", productTypeTable.statement)
	})

	t.Run("will_set_the_field_statement_and_table", func(t *testing.T) {
		// Parent Table
		for _, field := range productTable.Fields {
			require.Equal(t, fmt.Sprintf(`product."%s"`, field.Name), field.statement)
			require.Equal(t, &productTable, field.table)
		}

		// Related Table
		for _, field := range productTypeTable.Fields {
			require.Equal(t, fmt.Sprintf(`pt."%s"`, field.Name), field.statement)
			require.Equal(t, &productTypeTable, field.table)
		}
	})

	t.Run("will_return_the_fields_map", func(t *testing.T) {
		require.Equal(t, productTable.Fields[1], fieldsMap["name"])
		require.Equal(t, productTable.Fields[2], fieldsMap["count"])
		require.Equal(t, productTable.Fields[5], fieldsMap["created_on"])
		require.Equal(t, productTypeTable.Fields[1], fieldsMap["product_type.name"])
	})

	t.Run("will_return_the_search_fields", func(t *testing.T) {
		require.Equal(t, searchFields[0], productTable.Fields[1])
		require.Equal(t, searchFields[1], productTypeTable.Fields[1])
	})

	t.Run("will_return_the_default_where_statements", func(t *testing.T) {
		require.Equal(t, `product."is_deleted" is false`, defaultWhereStatements[0])
	})

	t.Run("will_return_the_default_sort_fields", func(t *testing.T) {
		require.Equal(t, *productTable.Fields[5], defaultSortFields[0])
	})

	t.Run("will_return_the_default_used_tables_map", func(t *testing.T) {
		assertDefaultUsedTablesMap := map[*Table]map[string]bool{
			&productTable:     {FILTER: true, SORT: true},
			&productTypeTable: {MANDATORY: true},
		}
		require.Equal(t, assertDefaultUsedTablesMap, defaultUsedTablesMap)
	})

	t.Run("error_if_has_duplicate_fields", func(t *testing.T) {
		duplicateFieldTable := Table{
			Name: "duplicate_field_table",
			Fields: []*Field{
				{Name: "duplicate_field"},
				{Name: "duplicate_field"},
			},
		}

		_, _, _, _, _, err := duplicateFieldTable.init()
		require.Error(t, err)
	})
}

func TestDefinition_getRelationTables(t *testing.T) {
	t.Run("simple_relation", func(t *testing.T) {
		results := productTable.getRequiredRelations(map[*Table]bool{&productTypeTable: true})
		require.NotEmpty(t, results)
	})

	t.Run("advanced_relation", func(t *testing.T) {
		commentTable := Table{Name: "comment"}
		categoryTable := Table{Name: "category"}
		postTable := Table{Name: "post"}
		postTable.Relations = append(postTable.Relations, &Relation{Table: &commentTable})
		postTable.Relations = append(postTable.Relations, &Relation{Table: &categoryTable})

		photoTable := Table{Name: "photo"}
		profileTable := Table{Name: "profile"}
		profileTable.Relations = append(profileTable.Relations, &Relation{Table: &photoTable})

		userTable := Table{Name: "user"}
		userTable.Relations = append(userTable.Relations, &Relation{Table: &postTable})
		userTable.Relations = append(userTable.Relations, &Relation{Table: &profileTable})

		results := userTable.getRequiredRelations(map[*Table]bool{&commentTable: true, &photoTable: true})
		require.NotEmpty(t, results)
	})
}
