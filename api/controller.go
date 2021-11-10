// controller.go creates controller functions for this app's routing engine.
package api

import (
	// "fmt"
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

// This function gets all outages from the database of this app
func GetOutages(w http.ResponseWriter, r *http.Request) {
    db := database.SetupDB()
	var outages []DBWaterOutage

    // Get all outages from outage table
    rows, err := db.Query("SELECT * FROM outage")
    if err != nil {
        log.Fatal(err)
    }

    // Map each row of the database to a DBWaterOutage struct
    for rows.Next() {
        var id, outageID int
        var address, location, startDate, endDate, outageType, createdAt, updatedAt string

		// Get data in the row
        err = rows.Scan(&id, &outageID, &address, &location, &startDate, &endDate, 
			&outageType, &createdAt, &updatedAt)
        if err != nil {
			log.Fatal(err)
		}

		// Save data to struct
        outages = append(outages, DBWaterOutage{
			OutageID: outageID, 
			Address: address,
			Location: location,
			StartDate: startDate,
			EndDate: endDate,
			OutageType: outageType,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
    }

	// Setup output headers & JSON
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(outages)

	// fmt.Printf("Someone retrieved all api data.\n")
}