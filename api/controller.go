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

// A DBWaterOutage struct maps a water outage from the database of this app.
type DBWaterOutage struct {
	OutageID int `json:"outage_id"`
	Street string `json:"street"`
	Suburb string `json:"suburb"`
	Location string `json:"location"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	OutageType string `json:"outage_type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetOutages JSON-encodes all outages from the database of this app
func GetOutages(w http.ResponseWriter, r *http.Request) {
	log.Println("Received GetOutage request.")

	// Query lines
	main := `SELECT outage_id, street, suburb, st_astext(location), start_date, end_date, outage_type, 
	created_at, updated_at FROM outage`

	// Get parameters and assembler filter query
	filter := ""
	params := r.URL.Query()

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
					}else if key == "end_date" {
						keyParams = append(keyParams, fmt.Sprintf("%s <= '%s'", key, element[0]))
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

	// Setup the database & model
    db := database.SetupDB()
	var outages []DBWaterOutage

	// Assemble query and get data from database
	rows, err := db.Query( main + filter )
	
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
			log.Fatal(err)
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

	// fmt.Printf("Someone retrieved all api data.\n")
}

// IsFilterableOutage returns true if a (url) parameter is filterable.
func IsFilterableOutage(param string) bool {
	if param == "suburb" || param == "street" || param == "outage_type" || 
	param == "start_date" || param == "end_date" || param == "location" {
		return true
	}

	// default false
	return false
}
