// filters.go contains functions that check and/or generate an
// SQL query depending on the given (url) parameters.
package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// MakeFilterQuery generates an SQL WHERE string and a string containing
// ORDER BY, LIMIT and OFFSET statements based on given parameters.
func MakeFilterQuery(r *http.Request, isCount bool) (where string, sort string) {
	// Get url params
	params := r.URL.Query()

	// Set up query object
	query := new(Query)
	query.IsCount = isCount

	// Make strings from params and query objects
	sort = query.MakeOrderbyPaginationString(params)

	// if parameters exist
	where = query.MakeWhereString(params)

	return
}

// MakeOrderbyPaginationString makes a string with an SQL
// order by, limit and offset string based on sort, limit
// and offset parameters if any.
// Default value:
//		" ORDER BY outage_id LIMIT 50 OFFSET 0"
// where outage_id can be replaced by the sort, while 50
// and 0 can be replaced by limit and offset parameters.
func (query *Query) MakeOrderbyPaginationString(
	params url.Values) string {
	if query.IsCount {
		return ""
	}

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

// MakeWhereString combines the Wheres into a single SQL
// WHERE statement, joined by an " AND " or " OR " SQL
// condition.
func (query *Query) MakeWhereString(
	params url.Values) (filter string) {
	// Make SQL Wheres from url params
	query.SetWheres(params)

	// Get SQL condition from excl param
	condition := GetSQLCondition(params.Get("excl"))

	// Join strings
	if len(query.Wheres) > 0 {
		return " WHERE " + strings.Join(query.Wheres, condition)
	}
	return
}

// SetWheres adds all SQL Wheres of the equivalent
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
			query.SetAddressOfTypeWhere(i, addressType)
		}
	}
}

// SetOutageIDWhere adds a SQL WHERE condition for all address
// values and types (street and suburb).
func (query *Query) SetAllAddressWheres(params url.Values) {
	query.SetAllAddressTypeWheres(params, "street")
	query.SetAllAddressTypeWheres(params, "suburb")
}
