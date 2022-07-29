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
					if key == "after_start_date" {
						key = "start_date"
						query.Wheres = append(query.Wheres, fmt.Sprintf("start_date >= '%s'", element[0]))
					} else if key == "before_start_date" {
						key = "start_date"
						query.Wheres = append(query.Wheres, fmt.Sprintf("start_date >= '%s'", element[0]))
					} else if key == "before_end_date" {
						key = "end_date"
						query.Wheres = append(query.Wheres, fmt.Sprintf("end_date <= '%s'", element[0]))
					} else if key == "after_end_date" {
						key = "end_date"
						query.Wheres = append(query.Wheres, fmt.Sprintf("end_date >= '%s'", element[0]))
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

// IsFilterableParam returns true if a (url) query parameter is filterable.
func IsFilterableParam(param string) bool {
	filterables := []string{
		"suburb", "street", "outage_type", "search",
		"before_start_date", "before_end_date", "after_end_date",
		"after_start_date", "start_date", "end_date",
		"location", "outage_id",
	}

	return isStringInArray(param, filterables)
}

// IsFilterableCountParam returns true if a query parameter is a count API
// parameter. It is intended to be an extension of IsFilterableParam.
func IsFilterableCountParam(param string) bool {
	filterables := []string{"total_hours", "total_outages"}

	return isStringInArray(param, filterables)
}

// Checks if an outage id is in a list of current outage ids.
func IsCurrentOutageID(outage_id int, current_outage_ids []int) bool {
	return isIntInArray(outage_id, current_outage_ids)
}
