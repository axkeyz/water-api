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

	// Get parameters and assemble filter query
	// params := r.URL.Query()
	main := `SELECT outage_id, street, suburb, st_astext(location), start_date, end_date, outage_type, 
	created_at, updated_at FROM outage`
	filter := MakeFilterQuery(r)
	// groupBy := ""

	// Assemble main query
	// count, _ := params["count"]
	// if count == nil {
	// 	main = `SELECT outage_id, street, suburb, st_astext(location), start_date, end_date, outage_type, 
	// 	created_at, updated_at FROM outage`
	// }else{
	// 	countStreet := ""
	// 	if count[0] == "street" {
	// 		// Add street to aggregate clause
	// 		countStreet = " street, "
	// 	}
	// 	// Create count api of data
	// 	main = fmt.Sprintf(
	// 		`SELECT %s suburb, outage_type, 
	// 		CASE 
	// 			WHEN outage_type = 'Planned' AND extract(day from end_date - start_date) > 0
	// 			THEN extract(day from end_date - start_date) * 2.85
	// 			ELSE EXTRACT(EPOCH FROM end_date-start_date)/3600
	// 		END total_hours, 
	// 		count(suburb) as total_outages FROM outage`, countStreet,
	// 	)
	// 	// Create group by clause
	// 	groupBy = fmt.Sprintf(
	// 		"GROUP BY %s suburb, outage_type, total_hours ORDER BY total_outages desc", countStreet,
	// 	)
	// }


	// Setup the database & model
    db := database.SetupDB()
	var outages []DBWaterOutage

	// Assemble query and get data from database
	rows, err := db.Query( main + filter)
	
	if err != nil {
		// Filter string is invalid.
		rows, _ = db.Query( main )
	}

	// Map each row of the database to a DBWaterOutage struct
	for rows.Next() {
		var outageID int
		var street, suburb, location, startDate, endDate, outageType, createdAt, updatedAt string

		// Get data in the row
		err = rows.Scan(&outageID, &street, &suburb, &location, &startDate, &endDate, 
			&outageType, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
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
			CreatedAt: createdAt[:19] + "+13:00",
			UpdatedAt: updatedAt[:19] + "+13:00",
		})

		log.Println(outages)
	}

	// Setup output headers & JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(outages)
}

// MakeFilterQuery generates an SQL WHERE string based on given parameters.
func MakeFilterQuery(r *http.Request) string {
	// Get params
	params := r.URL.Query()
	filter := ""

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
			}
		}

		if isValidFilter {
			// Join key parameters into final parameter string
			filter = " WHERE " + strings.Join(keyParams, " AND ")
		}
	}

	return filter
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