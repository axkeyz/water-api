// controller.go creates controller functions for this app's routing engine.
package api

import (
	"log"
	"encoding/json"
	"net/http"
	"github.com/axkeyz/water-down-again/database"
	_ "github.com/lib/pq"
)

// A DBWaterOutage struct maps a water outage from the database of this app.
type DBWaterOutage struct {
	OutageID int `json:"outage_id"`
	Address string `json:"address"`
	Location string `json:"location"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	OutageType string `json:"outage_type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetOutages json encodes all outages from the database of this app
func GetOutages(w http.ResponseWriter, r *http.Request) {
    db := database.SetupDB()
	var outages []DBWaterOutage

	// Get all outages from outage table
	rows, err := db.Query( `SELECT outage_id, address, location, start_date, end_date, outage_type, 
	created_at, updated_at FROM outage`)
	
	if err != nil {
		log.Fatal(err)
	}

	// Map each row of the database to a DBWaterOutage struct
	for rows.Next() {
		var outageID int
		var address, location, startDate, endDate, outageType, createdAt, updatedAt string

		// Get data in the row
		err = rows.Scan(&outageID, &address, &location, &startDate, &endDate, 
			&outageType, &createdAt, &updatedAt)
		if err != nil {
			log.Fatal(err)
		}

		// Save data to struct
		outages = append(outages, DBWaterOutage{
			OutageID: outageID, 
			Address: address,
			Location: location,
			StartDate: startDate[:19] + "+13:00",
			EndDate: endDate[:19] + "+13:00",
			OutageType: outageType,
			CreatedAt: createdAt[:19] + "+13:00",
			UpdatedAt: updatedAt[:19] + "+13:00",
		})
	}

	// Setup output headers & JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(outages)

	// fmt.Printf("Someone retrieved all api data.\n")
}