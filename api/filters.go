// filters.go contains functions that check and/or generate an
// SQL query depending on the given (url) parameters.
package api

import (
	"net/http"
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
