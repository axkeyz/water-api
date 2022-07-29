// filters.go contains functions that check and/or generate an
// SQL query depending on the given (url) parameters.
package api

import (
	"fmt"
	"net/http"
	"net/url"
)

// MakeFilterQuery generates an SQL WHERE string based on given parameters.
func MakeFilterQuery(r *http.Request) (string, string) {
	// Get params
	params := r.URL.Query()
	var filter, ordering string

	// if parameters exist
	if len(params) > 0 {
		query := new(Query)

		ordering = query.MakeOrderbyPaginationString(params)

		query.SetWheres(params)
		if len(query.Wheres) > 0 {
			filter = query.MakeWhereString(GetSQLCondition(params.Get("excl")))
		}
	}

	return filter, ordering
}

// MakeOrderbyPaginationString creates a string with an SQL
// order by, limit and offset string based on sort, limit
// and offset parameters if any.
// Default value:
//		" ORDER BY outage_id LIMIT 50 OFFSET 0"
// where outage_id can be replaced by the sort, while 50
// and 0 can be replaced by limit and offset parameters.
func (query *Query) MakeOrderbyPaginationString(
	params url.Values) string {
	if values := params["sort"]; len(values) > 0 {
		// Get parameters for sorting
		orderby := query.MakeOrderbyString(values)

		pagination := query.MakePaginationString(
			params["limit"][0], params["offset"][0],
		)

		// Combine
		return fmt.Sprintf(" %s %s", orderby, pagination)
	}
	return " ORDER BY outage_id LIMIT 50 OFFSET 0"
}

// SetWheres sets the Wheres file with the equivalent
// string for valid API parameters.
func (query *Query) SetWheres(params url.Values) {
	query.SetSearchWhere(params["search"])
	query.SetOutageTypeWhere(params.Get("outage_type"))
	query.SetOutageIDWhere(params.Get("outage_id"))

	query.SetDateWheres(params)
	query.SetAllAddressWheres(params)

	query.SetLocationRadiusWhere(
		params.Get("longitude"), params.Get("latitude"),
		params.Get("radius"),
	)
}

// SetDateWheres adds a SQL WHERE condition for all date
// parameters.
func (query *Query) SetDateWheres(params url.Values) {
	for _, param := range DateColumns {
		if value := params.Get(param); len(value) > 0 {
			_, column := IsDateParam(param)
			query.SetSignedWhere(column, value)
		}
	}
}

// SetAllAddressTypeWheres adds a SQL WHERE condition for
// all address values of a single address type.
func (query *Query) SetAllAddressTypeWheres(
	params url.Values, addressType string) {
	if values := params[addressType]; len(values) > 0 {
		for _, i := range values {
			query.SetOneLocationWhere(i, addressType)
		}
	}
}

// SetOutageIDWhere adds a SQL WHERE condition for all address
// values and types (street and suburb).
func (query *Query) SetAllAddressWheres(params url.Values) {
	query.SetAllAddressTypeWheres(params, "street")
	query.SetAllAddressTypeWheres(params, "suburb")
}
