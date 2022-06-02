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
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/axkeyz/water-down-again/database"
)

var street_abbreviations = map[string]string{
	"pl": "Place", "rd": "Road", "st": "Street", "ave": "Avenue", "dr": "Drive",
	"cr": "Crescent", "cres": "Crescent", "cresc": "Crescent", "blvd": "Boulevard",
	"bvd": "Boulevard", "cl": "Close", "hl": "Hill", "tce": "Terrace", "gln": "Glen",
	"pde": "Parade", "ct": "Court", "plza": "Plaza", "prom": "Promenade", "rt": "Retreat",
	"rtt": "Retreat", "rdge": "Ridge", "sq": "Square", "arc": "Arcade", "bwlk": "Boardwalk",
	"ch": "Chase", "cct": "Circuit", "bp": "Bypass", "bypa": "Bypass", "brk": "Break",
	"crst": "Crest", "ent": "Entrance", "esp": "Esplanade", "fwy": "Freeway",
	"glde": "Glade", "gd": "Glade", "gra": "Grange",
}

var suburb_abbreviations = map[string]string{
	"mt": "Mount", "pt": "Point", "st": "Saint",
}

// CleanAddressName removes all numbers (unit/street numbers, postcodes) from
// an address string and returns the address in title-case.
func CleanAddressName(address string, address_type string) string {
	location := strings.Split(address, " ")
	var cleaned_location []string
	var add_location bool

	// Iterate through each word
	for _, j := range location {
		add_location = true

		// Iterate through each letter
		for _, k := range j {
			if unicode.IsNumber(k) {
				// If a letter contains a number (unit, street, postcode), do not
				// add word to final address
				add_location = false
				break
			}
		}

		j = UnabbreviateAddressName(j, address_type)

		if add_location {
			// Only add words that do not contain numbers
			cleaned_location = append(cleaned_location, j)
		}
	}

	// Return the cleaned_location slice as a title-cased string
	caser := cases.Title(language.English)
	return caser.String(strings.Join(cleaned_location, " "))
}

// UnabbreviateAddressName replaces short-hands in addresses with the long form.
func UnabbreviateAddressName(address string, address_type string) string {
	// Convert to lowercase
	address = strings.ToLower(address)

	if address_type == "street" {
		if unabbreviated, ok := street_abbreviations[address]; ok {
			return unabbreviated
		}
	} else {
		if unabbreviated, ok := suburb_abbreviations[address]; ok {
			return unabbreviated
		}
	}
	return address
}

// GetAPIData returns the latest data as array of WaterOutage structs from the Watercare Outage API.
func GetAPIData() []WaterOutage {
	// Get data from the Watercare Outage API.
	response, err := http.Get("https://api.watercare.co.nz/outages/all")

	// Test API Route
	// response, err := http.Get("https://618a623134b4f400177c4603.mockapi.io/wateroutage")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// Read response data
	outagesJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	// Fold responseData into WaterOutage structs
	var outages []WaterOutage
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
		location := strings.Split(outages[i].Location, ",")

		if len(location) == 1 {
			location = strings.Split(location[0], " ")
			street := strings.Join(location[0:len(location)-1], " ")
			suburb := location[len(location)-1]
			location[0], location[1] = street, suburb
		}

		arrOutages[i] = fmt.Sprintf("(%d,'%s','%s','POINT(%f %f)','%s', '%s', '%s')",
			outages[i].OutageID, strings.TrimSpace(location[0]), strings.TrimSpace(location[1]), outages[i].Longitude,
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
