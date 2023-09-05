package resourceful

const (
	SELECT    = "select"
	SEARCH    = "search"
	FILTER    = "filter"
	SORT      = "sort"
	MANDATORY = "mandatory"
)

var sortOperators = map[string]string{
	"desc": "DESC",
	"asc":  "ASC",
}

var filterOperators = map[string]string{
	"eq":  "=",
	"ne":  "!=",
	"lt":  "<",
	"lte": "<=",
	"gt":  ">",
	"gte": ">=",
	//in string with quoted string separated by spaces
}

func validFilterOperator(filterOperator string) bool {
	_, ok := filterOperators[filterOperator]
	return ok || filterOperator == "in"
}
