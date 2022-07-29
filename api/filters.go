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

		for key, element := range params {
			if IsFilterableParam(key) {
				// Only append fiterable outages to key parameters list
				param := params[key]
				if param != nil {
					if isDate, column := IsDateParam(key); isDate {
						query.SetSignedWhere(column, element[0])
					} else if key == "location" {
						query.SetMapWhere(
							element[0], element[1], element[2],
						)
					} else if key == "outage_type" {
						column = GetEquationSignedColumn(key, 0)
						query.SetSignedWhere(column, element[0])
					} else if key == "search" {
						query.SetSearchWhere(element)
					} else {
						// Is a street / suburb
						for _, i := range element {
							query.SetLocationWhere(i)
						}
					}
				}
				query.IsValidWheres = true
			} else if key == "sort" {
				// Get parameters for sorting
				orderby := query.MakeOrderbyString(element)

				pagination := query.MakePaginationString(
					params["limit"][0], params["offset"][0],
				)

				// Combine
				order = fmt.Sprintf(" %s %s", orderby, pagination)
			}
		}

		if query.IsValidWheres {
			filter = query.MakeWhereString(GetSQLCondition(params.Get("excl")))
		}
	}

	return filter, order
}
