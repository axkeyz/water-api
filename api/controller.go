// controller.go creates controller functions for this app's routing engine.
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/axkeyz/water-down-again/database"
	_ "github.com/lib/pq"
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
	filter, order := MakeFilterQuery(r, false)

	// Setup the database & model
	db := database.SetupDB()
	defer db.Close()

	var outages []DBWaterOutage

	// Assemble query and get data from database
	rows, err := db.Query(main + filter + order)

	log.Println(main + filter + order)

	if err != nil {
		// Filter or order string is invalid.
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AppError{
			ErrorCode: 3440,
			Message:   "invalid parameters",
			Details:   "Parameters given for this API were invalid.",
		})
	}

	// Get current outage IDs
	current_outage_ids := GetCurrentOutageIDs()

	// Map each row of the database to a DBWaterOutage struct
	defer rows.Close()
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
				Message:   "unknown error",
				Details:   "Please contact me at xahkun@gmail.com to figure out this issue.",
			})
		}

		// Save data to struct
		outages = append(outages, DBWaterOutage{
			OutageID:   outageID,
			Street:     street,
			Suburb:     suburb,
			Location:   location,
			StartDate:  startDate[:19] + "+13:00",
			EndDate:    endDate[:19] + "+13:00",
			OutageType: outageType,
			Status:     IsCurrentOutageID(outageID, current_outage_ids),
		})

		// log.Println(outages)
	}

	// Setup output headers & JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(outages)
}

// CountOutages JSON-encodes outages from the database of this app in a count-based format.
func CountOutages(w http.ResponseWriter, r *http.Request) {
	log.Println("Received CountOutages request.")

	// Setup CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	// Setup database & output model
	db := database.SetupDB()
	defer db.Close()
	var outages []DBWaterOutage

	// Get parameters
	params := r.URL.Query()
	if params == nil {
		// Setup output headers & JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AppError{
			ErrorCode: 3443,
			Message:   "no parameters set",
			Details:   "No parameters were found. This API needs them to work.",
		})
	}

	// Create "filter" phrase, or the "WHERE" part in an SQL query
	fields := params["get"]

	filter, order := MakeFilterQuery(r, true)
	var grouped, selected []string

	for _, element := range fields {
		if IsFilterableParam(element) {
			if strings.Contains(element, "start_date") && element != "start_date" ||
				strings.Contains(element, "end_date") && element != "end_date" {
				element = strings.Join(strings.Split(element, "_")[1:2], "_")
			}
			grouped = append(grouped, element)
			selected = append(selected, element)
		} else if element == "total_hours" {
			selected = append(selected,
				`SUM(CASE WHEN outage_type = 'Planned' AND 
				EXTRACT(day from end_date - start_date) > 0
				THEN (EXTRACT(day from end_date - start_date) * 2.85)::float
				ELSE (EXTRACT(EPOCH FROM end_date-start_date)/3600)::float
				END) total_hours`,
			)
		}
	}

	// Create the GROUP BY part of an SQL query
	group := "GROUP BY " + strings.Join(grouped, ", ")

	if len(grouped) == 0 {
		group = ""
	}

	// Create select string
	var selects string
	if len(selected) > 0 {
		selects = strings.Join(selected, ", ") + ","
	} else {
		selects = ""
	}

	// Generate main query string
	main := fmt.Sprintf(
		`SELECT %s count(outage_id) as total_outages FROM outage %s %s 
		%s`, selects, filter, group, order,
	)

	// Assemble query and get data from database
	rows, err := db.Query(main)
	log.Println(main)

	if err != nil {
		// Filter string is invalid.
		log.Println(err)
		log.Println(main)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AppError{
			ErrorCode: 3445,
			Message:   "unknown error",
			Details:   "Please contact me at xahkun@gmail.com to figure out this issue.",
		})
	} else {
		// get the column names
		columns, err := rows.Columns()
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(AppError{
				ErrorCode: 3446,
				Message:   "unknown error",
				Details:   "Please contact me at xahkun@gmail.com to figure out this issue.",
			})
		}

		numColumns := len(columns)

		// Get current outage IDs
		current_outage_ids := GetCurrentOutageIDs()

		defer rows.Close()
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
					Message:   "unknown error",
					Details:   "Please contact me at xahkun@gmail.com to figure out this issue.",
				})
			}

			outage.Status = IsCurrentOutageID(outage.OutageID, current_outage_ids)

			// Append outage to all outages
			outages = append(outages, outage)
			// log.Println(outage)
		}

		// Setup output headers & JSON
		w.Header().Set("Content-Type", "application/json")
		//Allow CORS here By * or specific origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		json.NewEncoder(w).Encode(outages)
	}
	// }
}
