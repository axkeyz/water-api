// filters.go contains functions that check and/or generate an
// SQL query depending on the given (url) parameters.
package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
		var keyParams []string
		isValidFilter := false

		for key, element := range params {
			if IsFilterableParam(key) || key == "search" {
				log.Println("Received filter for " + key)

				// Only append fiterable outages to key parameters list
				param := params[key]
				if param != nil {
					if key == "after_start_date" {
						key = "start_date"
						keyParams = append(keyParams, fmt.Sprintf("start_date >= '%s'", element[0]))
					} else if key == "before_start_date" {
						key = "start_date"
						keyParams = append(keyParams, fmt.Sprintf("start_date >= '%s'", element[0]))
					} else if key == "before_end_date" {
						key = "end_date"
						keyParams = append(keyParams, fmt.Sprintf("end_date <= '%s'", element[0]))
					} else if key == "after_end_date" {
						key = "end_date"
						keyParams = append(keyParams, fmt.Sprintf("end_date >= '%s'", element[0]))
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
						if _, err := strconv.Atoi(element[0]); err == nil {
							// Is an integer - check if is outage id
							keyParams = append(keyParams,
								fmt.Sprintf(
									`lower(cast(outage_id as text)) 
									LIKE lower('%%%s%%')`,
									element[0],
								),
							)
						} else {
							keyParams = append(keyParams,
								fmt.Sprintf(
									`lower(suburb) LIKE lower('%%%s%%')
									OR lower(street) LIKE lower('%%%s%%') 
									OR lower(suburb) LIKE lower('%%%s%%')
									OR lower(street) LIKE lower('%%%s%%')`,
									element[0], element[0],
									CleanAddressName(element[0], "suburb"),
									CleanAddressName(element[0], "street"),
								),
							)
						}

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

				sort = append(sort, element...)

				sorted := strings.Join(sort, ", ")

				// Get sorting order (ascending / descending)
				pagination := ""
				limit := params["limit"]
				offset := params["offset"]
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

// IsFilterableParam returns true if a (url) query parameter is filterable.
func IsFilterableParam(param string) bool {
	filterables := []string{
		"suburb", "street", "outage_type",
		"before_start_date", "before_end_date", "after_end_date",
		"location", "outage_id", "start_date", "end_date",
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
