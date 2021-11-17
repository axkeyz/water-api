// controller.go creates controller functions for this app's routing engine.
package api

import (
	"log"
	"encoding/json"
	"net/http"
	"github.com/axkeyz/water-down-again/database"
	_ "github.com/lib/pq"
	"fmt"
	"strings"
)

// GetOutages JSON-encodes all outages from the database of this app.
func GetOutages(w http.ResponseWriter, r *http.Request) {
	log.Println("Received GetOutage request.")

	// Setup CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
    if r.Method == http.MethodOptions {
        return
    }

	// Get parameters and assemble filter query
	main := `SELECT outage_id, street, suburb, st_astext(location), start_date, end_date, 
	outage_type FROM outage`
	filter, order := MakeFilterQuery(r)

	// Setup the database & model
    db := database.SetupDB()
	var outages []DBWaterOutage

	// Assemble query and get data from database
	rows, err := db.Query(main + filter + order)
	log.Println(main + filter + order)
	
	if err != nil {
		// Filter or order string is invalid.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AppError{
			ErrorCode: 3440,
			Message: "invalid parameters",
			Details: "Parameters given for this API were invalid.",
		})
	}

	// Map each row of the database to a DBWaterOutage struct
	for rows.Next() {
		var outageID int
		var street, suburb, location, startDate, endDate, outageType string

		// Get data in the row
		err = rows.Scan(&outageID, &street, &suburb, &location, &startDate, &endDate, 
			&outageType)
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(AppError{
				ErrorCode: 3441,
				Message: "unknown error",
				Details: "Please contact me at xahkun@gmail.com to figure out this issue.",
			})
		}

		// Save data to struct
		outages = append(outages, DBWaterOutage{
			OutageID: outageID, 
			Street: street,
			Suburb: suburb,
			Location: location,
			StartDate: startDate[:19] + "+13:00",
			EndDate: endDate[:19] + "+13:00",
			OutageType: outageType,
		})

		// log.Println(outages)
	}

	// Setup output headers & JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(outages)
}

// CountOutages JSON-encodes outages from the database of this app in a count-based format.
func CountOutages(w http.ResponseWriter, r *http.Request){
	log.Println("Received CountOutages request.")

	// Setup CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
    if r.Method == http.MethodOptions {
        return
    }

	// Setup database & output model
	db := database.SetupDB()
	var outages []DBWaterOutage

	// Get parameters
	params := r.URL.Query()
	if params == nil {
		// Setup output headers & JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AppError{
			ErrorCode: 3443,
			Message: "no parameters set",
			Details: "No parameters were found. This API needs them to work.",
		})
	}

	fields,_ := params["get"]

	if fields == nil {
		// Setup output headers & JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AppError{
			ErrorCode: 3444,
			Message: "required parameters not found",
			Details: "This API requires a get parameter.",
		})
	}else{
		// Generate filter (where) & group by string
		filter, order := MakeFilterQuery(r)
		var grouped, selected []string

		for _, element := range fields {
			if IsFilterableOutage(element) && element != "aend_date" {
				grouped = append(grouped, element)
				selected = append(selected, element)
			}else if element == "total_hours" {
				selected = append(selected, 
					`SUM(CASE WHEN outage_type = 'Planned' AND 
					EXTRACT(day from end_date - start_date) > 0
					THEN (EXTRACT(day from end_date - start_date) * 2.85)::float
					ELSE (EXTRACT(EPOCH FROM end_date-start_date)/3600)::float
					END) total_hours`,
				)
			}
		}

		group := strings.Join(grouped, ", ")
		selects := strings.Join(selected, ", ")

		// Generate main query string
		main := fmt.Sprintf(
			`SELECT %s, count(outage_id) as total_outages FROM outage %s GROUP BY %s 
			%s`, selects, filter, group, order,
		)

		// Assemble query and get data from database
		rows, err := db.Query(main)
		log.Println(main)
		
		if err != nil {
			// Filter string is invalid.
			log.Println(main)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(AppError{
				ErrorCode: 3445,
				Message: "unknown error",
				Details: "Please contact me at xahkun@gmail.com to figure out this issue.",
			})
		}else{
			// get the column names
			columns, err := rows.Columns()
			if err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(AppError{
					ErrorCode: 3446,
					Message: "unknown error",
					Details: "Please contact me at xahkun@gmail.com to figure out this issue.",
				})
			}

			numColumns := len(columns)

			for rows.Next() {
				// Create new outage
				outage := DBWaterOutage{}

				// make references for the columns by calling DBWaterOutageCol
				column := make([]interface{}, numColumns)
				for i := 0; i < numColumns; i++ {
					column[i] = DBWaterOutageCol(columns[i], &outage)
				}

				err = rows.Scan(column...)
				if err != nil {
					log.Println(err)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(AppError{
						ErrorCode: 3447,
						Message: "unknown error",
						Details: "Please contact me at xahkun@gmail.com to figure out this issue.",
					})
				}

				// Append outage to all outages
				outages = append(outages, outage)
				// log.Println(outage)
			}

			// Setup output headers & JSON
			w.Header().Set("Content-Type", "application/json")
			//Allow CORS here By * or specific origin
			w.Header().Set("Access-Control-Allow-Origin", "*.aileenhuang.dev")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			json.NewEncoder(w).Encode(outages)
		}
	}
}

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
			if IsFilterableOutage(key) {
				log.Println("Received filter for "+key)

				// Only append fiterable outages to key parameters list
				param, _ := params[key]
				if param != nil {
					if key == "start_date" {
						keyParams = append(keyParams, fmt.Sprintf("%s >= '%s'", key, element[0]))
					}else if key == "aend_date" {
						keyParams = append(keyParams, fmt.Sprintf("end_date <= '%s'", element[0]))
					}else if key == "location" {
						radius, _ := params["radius"]
						longitude, _ := params["longitude"]
						latitude, _ := params["latitude"]

						keyParams = append(keyParams, fmt.Sprintf(
							"ST_DWithin(location, ST_SetSRID(ST_Point(%s, %s), 4326), %s)",
							longitude[0], latitude[0], radius[0],
						))
					}else{
						keyParams = append(keyParams, fmt.Sprintf("%s = '%s'", key, element[0]))
					}
				}
				isValidFilter = true
			}else if key == "sort" {
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
	param == "start_date" || param == "end_date" || param == "aend_date" ||
	param == "location" {
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