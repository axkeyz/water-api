// source.go contains the functions that extract data from the Watercare API and
// save the data in this app's database.
package api

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
	"encoding/json"
	"strings"
	"github.com/axkeyz/water-down-again/database"
)

// A WaterOutage struct maps a water outage from the Watercare API instance.
type WaterOutage struct {
	OutageID int64 `json:"outageId"`
	Location string `json:"location"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
	OutageType string `json:"outageType"`
}

// GetAPIData returns the latest data as array of WaterOutage structs from the Watercare Outage API.
func GetAPIData() []WaterOutage {
	// Get data from the Watercare Outagee API.
    // response, err := http.Get("https://api.watercare.co.nz/outages/all")

	// Currently use test API
	response, err := http.Get("https://618a623134b4f400177c4603.mockapi.io/wateroutage")
    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }
	
	// Read response data
    outagesJSON, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

	// Fold responseData into WaterOutage structs
	var outages []WaterOutage
	json.Unmarshal([]byte(string(outagesJSON)), &outages)

	return outages
}

// UnpackAPIData converts an array of WaterOutage structs into a string for SQL insertion.
// The data is returned as comma separated list where each element looks like:
// (OutageID, 'Location', '(Longitude, Latitude)', 'StartDate', 'EndDate', 'OutageType')
func UnpackAPIData(outages []WaterOutage) (string){
	// Initialise all variables
	numOutage := len(outages)
	arrOutages := make([]string, numOutage)

	// Loop through WaterOutages, separate and assign to individual array
	for i := range outages {
		location := fmt.Sprintf("%s", outages[i].Location)
		if (location[0] == ' '){
			location = location[1:]
		}
		arrOutages[i] = fmt.Sprintf("(%d,'%s','(%f, %f)','%s', '%s', '%s')", 
		outages[i].OutageID, location, outages[i].Longitude, outages[i].Latitude,
		outages[i].StartDate, outages[i].EndDate, outages[i].OutageType)
	}

	return strings.Join(arrOutages[:], ", ")
}

// WriteOutage upserts outages in the database.
// If the outage exists in the database (based on outage_id), WriteOutage attempts to update the endDate if
// applicable. If the outage does not exist in the database, WriteOutage creates a new record.
func WriteOutage(outage []WaterOutage) {
	// Open database
	db := database.SetupDB()

	// Prepare SQL Statement
	sqlStatement := `insert into outage (outage_id, address, location, start_date, end_date, outage_type) values  
	%s on conflict (outage_id) do update set end_date = excluded.end_date;`
	outages := UnpackAPIData(outage)
	combined := fmt.Sprintf(sqlStatement, outages)

	// run sql statement
	_, err := db.Exec(combined)
	if err != nil {
		log.Fatal(err)
	}
}

// This function gets the latest data from the Watercare API and upserts the data into the database.
func UpdateOutages() {
	outages := GetAPIData()
	WriteOutage(outages)
}