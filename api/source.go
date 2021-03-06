// source.go contains the functions that extract data from the Watercare API and
// save the data in this app's database.
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

// GetAPIData returns the latest data as array of WaterOutage structs from the Watercare Outage API.
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

// UnpackAPIData converts an array of WaterOutage structs into a string for SQL insertion.
// The data is returned as comma separated list where each element looks like:
// (OutageID, 'Street', 'Suburb', '(Longitude, Latitude)', 'StartDate', 'EndDate', 'OutageType')
func UnpackAPIData(outages []WaterOutage) string {
	// Initialise all variables
	numOutage := len(outages)
	arrOutages := make([]string, numOutage)

	// Loop through WaterOutages, separate and assign to individual array
	for i := range outages {
		street, suburb := AddressToStreetSuburb(outages[i].Location)

		arrOutages[i] = fmt.Sprintf("(%d,'%s','%s','POINT(%f %f)','%s', '%s', '%s')",
			outages[i].OutageID, strings.TrimSpace(street), strings.TrimSpace(suburb), outages[i].Longitude,
			outages[i].Latitude, outages[i].StartDate, outages[i].EndDate, outages[i].OutageType)
	}

	return strings.Join(arrOutages[:], ", ")
}

// WriteOutage upserts outages in the database.
// If the outage exists in the database (based on outage_id), WriteOutage attempts to update the endDate if
// applicable. If the outage does not exist in the database, WriteOutage creates a new record.
func WriteOutage(outage []WaterOutage) {
	// Open database
	db := database.SetupDB()
	defer db.Close()

	// Prepare SQL Statement
	sqlStatement := `insert into outage (outage_id, street, suburb, location, start_date, end_date, outage_type)  
	values %s on conflict (outage_id) do update set end_date = excluded.end_date;`
	outages := UnpackAPIData(outage)
	combined := fmt.Sprintf(sqlStatement, outages)

	// run sql statement
	_, err := db.Exec(combined)
	if err != nil {
		log.Println(err)
	}
}

// UpdateOutages gets the latest data from the Watercare API and upserts the data into the database.
func UpdateOutages() {
	outages := GetAPIData()
	WriteOutage(outages)
	log.Println("Outage list has been updated.")
}

// GetCurrentOutageIDs returns an int slice of all currently active outage ids.
func GetCurrentOutageIDs() []int {
	current_outages := GetAPIData()
	var current_outage_ids []int

	for _, current_outage := range current_outages {
		current_outage_ids = append(current_outage_ids, current_outage.OutageID)
	}

	return current_outage_ids
}
