// cleanup.go contains functions that format addresses.
package api

import (
	"log"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/axkeyz/water-down-again/database"
	_ "github.com/lib/pq"
)

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
	return strings.TrimSpace(strings.Join(cleaned_location, " "))
}

// UnabbreviateAddressName replaces short-hands in addresses with the long form.
func UnabbreviateAddressName(address string, address_type string) string {
	caser := cases.Title(language.English)

	// If address is an abbreviation, return the uncodensed form after checking the address type
	if address_type == "street" {
		if unabbreviated, ok := street_abbreviations[address]; ok {
			return caser.String(unabbreviated)
		}
	} else {
		if unabbreviated, ok := suburb_abbreviations[address]; ok {
			return caser.String(unabbreviated)
		}
	}

	return caser.String(address)
}

// AddressToStreetSuburb attempts to return the street and suburb from the given address.
func AddressToStreetSuburb(address string) (string, string) {
	address = strings.ToLower(address)
	address_slice := strings.Split(address, ",")
	length := len(address_slice)

	// Default returns if location is not detected
	street := address
	suburb := "UNKNOWN"

	if length >= 2 {
		// Address is comma-separated, usually in the format: [unit/flat, ]street, suburb[, auckland]
		index := length
		for i := 0; i < length; i++ {
			if strings.Contains(address_slice[i], "auckland") &&
				!strings.Contains(address_slice[i], "auckland central") {
				// Reassign index if there are extra commas after the suburb
				index = i
				break
			}
		}
		// Save the street and suburb
		street, suburb = address_slice[index-2], address_slice[index-1]
	} else {
		// Address is not comma-separated, usually in the format: street suburb[ auckland]
		var position int
		for _, i := range suburbs {
			// Find last occurence of suburb (i) to avoid road names that have the suburb name included
			position = strings.LastIndex(address, i)
			if position > 0 {
				// Position found, save the street and suburb
				suburb = address[position : position+len(i)]
				street = address[0:position]
				break
			}
		}
	}

	// Return formatted version of street and suburb
	return CleanAddressName(street, "street"), CleanAddressName(suburb, "suburb")
}

// CleanupOutages re-formats the all existing outages in the database
func CleanupOutages() {
	// Open database
	db := database.SetupDB()
	defer db.Close()

	// Prepare SQL Statement
	query := `SELECT id, lower(street), lower(suburb) FROM outage`
	rows, err := db.Query(query)

	var suburb, street string
	var id int

	if err != nil {
		log.Println("Cleanup outages failed:", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			// Get data in the row
			err = rows.Scan(&id, &street, &suburb)
			if err != nil {
				log.Println("Cleanup outages failed:", err)
			}

			_, err := db.Exec(
				"UPDATE outage SET street = $1, suburb = $2 where id = $3",
				CleanAddressName(street, "street"), CleanAddressName(suburb, "suburb"), id,
			)

			if err != nil {
				log.Println("Cleanup outages failed:", err)
			}
		}
	}

	log.Println("Outages have been cleaned up.")
}
