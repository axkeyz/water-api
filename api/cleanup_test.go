package api

import (
	"testing"
)

// TestUnabbreviateAddressName calls api.UnabbreviateAddressName and
// checks if street/suburb shorthands are correctly replaced.
func TestUnabbreviateAddressName(t *testing.T) {
	street_inputs := []string{"st", "street", "jUpItER", "arc", "SCENIC"}
	street_outputs := []string{"Street", "Street", "Jupiter", "Arcade", "Scenic"}
	suburb_inputs := []string{"st", "mt", "Mount", "AlberTOS", "EARTHY"}
	suburb_outputs := []string{"Saint", "Mount", "Mount", "Albertos", "Earthy"}

	// Check for issues
	for i := 0; i < len(street_inputs); i++ {
		if input := UnabbreviateAddressName(street_inputs[i], "street"); input != street_outputs[i] {
			t.Fatalf(
				`TestUnabbreviateAddressName(street) did not return %s,
				got %s`,
				street_outputs[i], input,
			)
		}

		if input := UnabbreviateAddressName(suburb_inputs[i], "suburb"); input != suburb_outputs[i] {
			t.Fatalf(
				`TestUnabbreviateAddressName(suburb) did not return %s,
				got %s`,
				suburb_outputs[i], input,
			)
		}
	}
}

// TestAddressToStreetSuburb calls api.AddressToStreetSuburb and checks
// the output of the street and suburb.
func TestAddressToStreetSuburb(t *testing.T) {
	address1 := "Flat 21B St MEADOWLAND St Titirangi Auckland 1023"
	address2 := "23A Queen Street, auckland cbd"
	street1, suburb1 := AddressToStreetSuburb(address1)
	street2, suburb2 := AddressToStreetSuburb(address2)

	// Check for issues
	if street1 != "Saint Meadowland Street" && suburb1 != "Titirangi" {
		t.Fatalf(
			`TestAddressToStreetSuburb did not return %s,
			got %s, %s`, "Saint Meadowland Street, Titirangi",
			street1, suburb1,
		)
	}

	// Check for issues
	if street2 != "Queen Street" && suburb2 != "Auckland Central" {
		t.Fatalf(
			`TestAddressToStreetSuburb did not return %s,
			got %s, %s`, "Queen Street, Auckland Central",
			street2, suburb2,
		)
	}
}

// TestCleanAddressName calls api.CleanAddressName and checks if the output
// is the address in the expected format.
func TestCleanAddressName(t *testing.T) {
	street := "Flat 21B St MEADOWLAND St"
	suburb := "Mt URANUS 1023"

	// Check for issues
	if CleanAddressName(street, "street") != "Saint Meadowland Street" {
		t.Fatalf(
			`TestCleanAddressName did not return Saint Meadowland Street, got %s`,
			CleanAddressName(street, "street"),
		)
	}

	// Check for issues
	if CleanAddressName(suburb, "suburb") != "Mount Uranus" {
		t.Fatalf(
			`TestCleanAddressName did not return Mount Uranus, got %s`,
			suburb,
		)
	}
}
