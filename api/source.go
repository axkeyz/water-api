// source.go contains the functions that extract data from the
// Watercare API and save the data in this app's database.
package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/axkeyz/water-down-again/database"
)

// GetAPIData returns the latest data as array of WaterOutage
// structs from the Watercare Outage API.
func GetAPIData() []WaterOutage {
	var outages []WaterOutage

	// Get data from the Watercare Outage API.
	response, err := http.Get(os.Getenv("SRC_API"))

	if err != nil {
		log.Println(err.Error())
		return outages
	}

	// Read response data
	outagesJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return outages
	}

	// Fold responseData into WaterOutage structs
	json.Unmarshal([]byte(string(outagesJSON)), &outages)

	return outages
}

// UnpackAPIData converts an array of WaterOutage structs into
// a SQL string for bulk insert.
func UnpackAPIData(outages []WaterOutage) string {
	// Initialise all variables
	arrOutages := make([]string, len(outages))

	// Loop through WaterOutages, separate and assign to
	// individual array
	for i := range outages {
		arrOutages[i] = UnpackSingleAPIData(outages[i])
	}

	return strings.Join(arrOutages[:], ", ")
}

// UnpackSingleAPIData converts a single WaterOutage struct to
// a specific formatted string for bulk insert.
// Format:
// `(OutageID, "Street", "Suburb", "(Longitude, Latitude)",
// "StartDate", "EndDate", "OutageType")`
func UnpackSingleAPIData(outage WaterOutage) string {
	street, suburb := AddressToStreetSuburb(
		outage.Location)

	return fmt.Sprintf(
		"(%d,'%s','%s','POINT(%f %f)','%s', '%s', '%s')",
		outage.OutageID, strings.TrimSpace(street),
		strings.TrimSpace(suburb), outage.Longitude,
		outage.Latitude, outage.StartDate,
		outage.EndDate, outage.OutageType,
	)
}

// MakeWriteOutageQuery returns an SQL string to bulk insert
// multiple outages.
func MakeWriteOutageQuery(outage []WaterOutage) string {
	// Prepare SQL Statement
	sqlStatement := `insert into outage (outage_id, street, 
		suburb, location, start_date, end_date, outage_type)  
		values %s on conflict (outage_id) do update SET 
		end_date = excluded.end_date;`
	outages := UnpackAPIData(outage)

	return fmt.Sprintf(sqlStatement, outages)
}

// WriteOutage upserts outages in the database. If the outage
// exists in the database (based on outage_id), WriteOutage
// attempts to update the endDate if applicable. If the outage
// does not exist in the database, WriteOutage creates a
// new record.
func WriteOutage(outage []WaterOutage) {
	// Open database
	db := database.SetupDB()
	defer db.Close()

	// Prepare SQL Statement
	query := MakeWriteOutageQuery(outage)

	// run sql statement
	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

// UpdateOutages gets the latest data from the Watercare API
// and upserts the data into the database.
func UpdateOutages() {
	outages := GetAPIData()
	WriteOutage(outages)
	log.Println("Outage list has been updated.")
}

// GetCurrentOutageIDs returns an int slice of all currently
// active outage ids.
func GetCurrentOutageIDs() (current_outage_ids []int) {
	current_outages := GetAPIData()

	for _, current_outage := range current_outages {
		current_outage_ids = append(
			current_outage_ids, current_outage.OutageID)
	}

	return
}
