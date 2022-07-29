// filtermodels.go contains the SQL query components.
package api

import (
	"fmt"
	"strings"
)

type Query struct {
	Selects       []string
	Table         string
	Wheres        []string
	IsValidWheres bool
	Orderbys      []string
	Limit         string
	Offset        string
	GroupBy       []string
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
// WHERE statement, joined by "AND".
func (query *Query) MakeWhereString() string {
	return " WHERE " + strings.Join(query.Wheres, " AND ")
}

// MakeOrderByString returns an orderby string after adding
// the orderby strings to query.Orderbys.
// The orderby strings are in the format "column_name asc/desc"
func (query *Query) MakeOrderbyString(
	orderbys []string,
) string {
	query.SetOrderbysField(orderbys)
	return query.MakeOrderbyStringFromFields()
}

// SetOrderbyFields adds strings in the format "column_name asc/desc"
// to the orderbys field.
func (query *Query) SetOrderbysField(orderbys []string) {
	query.Orderbys = append(query.Orderbys, orderbys...)
}

// MakeOrderString makes a single SQL ORDER BY statement by
// combining the Orderbys array.
func (query *Query) MakeOrderbyStringFromFields() string {
	return " ORDER BY " + strings.Join(query.Orderbys, ", ")
}

// MakePaginationString returns a limit-offset string after setting the
// query to corresponding limit and offset values.
func (query *Query) MakePaginationString(
	limit string, offset string,
) string {
	query.SetPaginationFields(limit, offset)
	return query.MakePaginationStringFromFields()
}

// Set pagination fields sets the limit and offset pagination fields.
func (query *Query) SetPaginationFields(
	limit string, offset string,
) {
	query.Limit = limit
	query.Offset = offset
}

// MakePaginationString makes an SQL limit-offset pagination string
// by using the query limit and offset.
func (query *Query) MakePaginationStringFromFields() (
	pagination string) {
	// Get sorting order (ascending / descending)
	if query.Limit != "" && query.Offset != "" {
		// Pagination string
		pagination = fmt.Sprintf(
			"LIMIT %s OFFSET %s", query.Limit, query.Offset,
		)
	}

	return
}
