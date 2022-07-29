// filtermodels.go contains the SQL query components.
package api

import (
	"fmt"
	"strings"
)

type Query struct {
	Selects []string
	Table   string
	Wheres  []string
	Order   string
	Sorts   []string
	Limit   int
	Offset  int
	GroupBy []string
}

// SetSearchWhere adds a search by location name (street/suburb)
// or outage id depending on the value of user input.
func (query *Query) SetSearchWhere(value []string) {
	for _, i := range value {
		if isInt(i) {
			query.SetSearchIDWhere(i)
		} else {
			query.SetSearchLocationWhere(i)
		}
	}
}

// SetSearchLocationWhere adds a search by location name
// (street/suburb) to the array of wheres.
func (query *Query) SetSearchLocationWhere(location string) {
	query.Wheres = append(
		query.Wheres, fmt.Sprintf(
			`(lower(suburb) LIKE lower('%%%s%%')
			OR lower(street) LIKE lower('%%%s%%') 
			OR lower(suburb) LIKE lower('%%%s%%')
			OR lower(street) LIKE lower('%%%s%%'))`,
			location, location,
			CleanAddressName(location, "suburb"),
			CleanAddressName(location, "street"),
		),
	)
}

// SetSearchIDWhere adds a search by outage id to the
// array of wheres.
func (query *Query) SetSearchIDWhere(id string) {
	query.Wheres = append(query.Wheres,
		fmt.Sprintf(
			`lower(cast(outage_id as text)) 
			LIKE lower('%%%s%%')`, id,
		),
	)
}

// MakeWhereString combines the Wheres into a single SQL
// WHERE statement.
//
// For example:
//		Query.Wheres = [
//			"street = 'URANIUM ROAD'",
//			"end_date >= '2006-02-01'",
//		]
// Output:
//		"WHERE street = 'URANIUM ROAD' AND end_date >=
//		'2006-02-01'"
func (query *Query) MakeWhereString() string {
	return " WHERE " + strings.Join(query.Wheres, " AND ")
}
