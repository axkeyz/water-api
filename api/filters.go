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
	order := ""

	// if parameters exist
	if len(params) > 0 {
		var keyParams []string
		isValidFilter := false

		for key, element := range params {
			if IsFilterableOutage(key) || key == "search" {
				log.Println("Received filter for " + key)

				// Only append fiterable outages to key parameters list
				param, _ := params[key]
				if param != nil {
					if key == "start_date" {
						keyParams = append(keyParams, fmt.Sprintf("%s >= '%s'", key, element[0]))
					} else if key == "end_date" {
						keyParams = append(keyParams, fmt.Sprintf("end_date <= '%s'", element[0]))
					} else if key == "location" {
						radius := element[2]
						longitude := element[0]
						latitude := element[1]

						keyParams = append(keyParams, fmt.Sprintf(
							"ST_DWithin(location, ST_SetSRID(ST_Point(%s, %s), 4326), %s)",
							longitude, latitude, radius,
						))
					} else if key == "outage_type" {
						keyParams = append(keyParams, fmt.Sprintf("%s = '%s'", key, element[0]))
					} else if key == "search" {
						keyParams = append(keyParams,
							fmt.Sprintf(
								`lower(suburb) LIKE lower('%%%s%%') OR
								lower(street) LIKE lower('%%%s%%') OR
								lower(cast(outage_id as text)) LIKE lower('%%%s%%')`,
								element[0], element[0], element[0],
							),
						)
					} else {
						// Allow chaining (query equivalent = OR)
						var elems []string
						for _, i := range element {
							elems = append(elems, fmt.Sprintf("lower(cast(%s as text)) LIKE lower('%%%s%%')", key, i))
						}

						keyParams = append(keyParams, strings.Join(elems, " OR "))
					}
				}
				isValidFilter = true
			} else if key == "sort" {
				// Get parameters for sorting
				var sort []string

				for _, i := range element {
					sort = append(sort, i)
				}

				sorted := strings.Join(sort, ", ")

				// Get sorting order (ascending / descending)
				pagination := ""
				limit, _ := params["limit"]
				offset, _ := params["offset"]
				if limit != nil && offset != nil {
					// Pagination string
					pagination = fmt.Sprintf("LIMIT %s OFFSET %s", limit[0], offset[0])
				}

				// Combine
				order = fmt.Sprintf(" ORDER BY %s %s", sorted, pagination)
			}
		}

		if isValidFilter {
			// Join key parameters into final parameter string
			filter = " WHERE " + strings.Join(keyParams, " AND ")
		}
	}

	return filter, order
}

// IsFilterableOutage returns true if a (url) parameter is filterable.
func IsFilterableOutage(param string) bool {
	if param == "suburb" || param == "street" || param == "outage_type" ||
		param == "start_date" || param == "end_date" ||
		param == "location" || param == "outage_id" {
		return true
	}

	// default false
	return false
}

// IsFilterableCountOutage extends IsFilterableOutage with extra filters. It
// is attended to be a companion to IsFilterableOutage for the Count API.
func IsFilterableCountOutage(param string) bool {
	if param == "total_hours" || param == "total_outages" {
		return true
	}
	// default false
	return false
}
