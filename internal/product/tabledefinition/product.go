package tabledefinition

import "mceasy/service-demo/pkg/resourceful"

var Product = &resourceful.Table{
	Name: "product",
	Fields: []*resourceful.Field{
		{Name: "company_id", Type: resourceful.NUMERIC, LocalFilterable: true},
		{Name: "uuid"},
		{Name: "name", Searchable: true, Sortable: true},
		{Name: "description", Searchable: true, Sortable: true},
		{Name: "price", Searchable: true, Sortable: true},
	},
}
