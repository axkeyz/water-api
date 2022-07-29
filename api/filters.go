// filters.go contains functions that check and/or generate an
// SQL query depending on the given (url) parameters.
package api

import (
	"fmt"
	"net/http"
)

// MakeFilterQuery generates an SQL WHERE string based on given parameters.
func MakeFilterQuery(r *http.Request) (string, string) {
	// Get params
	params := r.URL.Query()
	filter := ""
	order := " LIMIT 50 OFFSET 0"

	// if parameters exist
	if len(params) > 0 {
		query := new(Query)

		query.SetWheres(params)

		if values := params["sort"]; len(values) > 0 {
			// Get parameters for sorting
			orderby := query.MakeOrderbyString(values)

			pagination := query.MakePaginationString(
				params["limit"][0], params["offset"][0],
			)

			// Combine
			order = fmt.Sprintf(" %s %s", orderby, pagination)
		}

		if len(query.Wheres) > 0 {
			filter = query.MakeWhereString(GetSQLCondition(params.Get("excl")))
		}
	}

	return filter, order
}
