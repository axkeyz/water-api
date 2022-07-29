// filters.go contains functions that check and/or generate an
// SQL query depending on the given (url) parameters.
package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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
				log.Println("Received filter for " + key)

				// Only append fiterable outages to key parameters list
				param := params[key]
				if param != nil {
					if isDate, column := IsDateParam(key); isDate {
						query.Wheres = append(
							query.Wheres,
							fmt.Sprintf("%s '%s'", column, element[0]),
						)
					} else if key == "location" {
						radius := element[2]
						longitude := element[0]
						latitude := element[1]

						query.Wheres = append(query.Wheres, fmt.Sprintf(
							"ST_DWithin(location, ST_SetSRID(ST_Point(%s, %s), 4326), %s)",
							longitude, latitude, radius,
						))
					} else if key == "outage_type" {
						query.Wheres = append(query.Wheres, fmt.Sprintf("%s = '%s'", key, element[0]))
					} else if key == "search" {
						query.SetSearchWhere(element)
					} else {
						// Allow chaining (query equivalent = OR)
						var elems []string
						for _, i := range element {
							elems = append(elems, fmt.Sprintf("lower(cast(%s as text)) LIKE lower('%%%s%%')", key, i))
						}

						query.Wheres = append(query.Wheres, strings.Join(elems, " OR "))
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
			// Join key parameters into final parameter string
			filter = query.MakeWhereString()
		}
	}

	return filter, order
}
