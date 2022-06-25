package api

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

// TestGetAPIData calls api.GetAPIData and checks if is a valid []WaterOutage.
func TestGetAPIData(t *testing.T) {
	// Set environmental variables to .env
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	// Get API Data
	src := GetAPIData()
	src_type := fmt.Sprintf("%T", src)

	// Check for issues
	if src_type != "[]api.WaterOutage" {
		t.Fatalf(
			`TestGetAPIData did not return a []api.WaterOutage, got type %s`,
			src_type,
		)
	}
}

// TestUnpackAPIData calls api.UnpackAPIData and checks if a correctly
// formatted SQL string is generated.
func TestUnpackAPIData(t *testing.T) {
	// Sample data
	packed_api := []WaterOutage{
		{
			OutageID:   15988,
			Location:   "52 Uranus Street",
			Latitude:   -36.9089914167,
			Longitude:  174.8325907667,
			StartDate:  "2022-06-20T22:00:00+12:00",
			EndDate:    "2022-06-21T03:00:00+12:00",
			OutageType: "Planned",
		},
		{
			OutageID:   26344,
			Location:   "34 Mercury Road",
			Latitude:   -23.9029914167,
			Longitude:  175.8343907667,
			StartDate:  "2022-05-15T24:00:00+12:00",
			EndDate:    "2022-07-27T05:00:00+12:00",
			OutageType: "Unplanned",
		},
	}

	// Get outputs
	actual_output := UnpackAPIData(packed_api)
	expected_output := "(15988,'Uranus Street','Unknown','POINT(174.832591 -36.908991)'," +
		"'2022-06-20T22:00:00+12:00', '2022-06-21T03:00:00+12:00', 'Planned'), " +
		"(26344,'Mercury Road','Unknown','POINT(175.834391 -23.902991)'," +
		"'2022-05-15T24:00:00+12:00', '2022-07-27T05:00:00+12:00', 'Unplanned')"

	// check for issues
	if actual_output != expected_output {
		t.Fatalf(
			`TestUnpackAPIData did not unpack data correctly, expected %s
			got %s`,
			expected_output, actual_output,
		)
	}
}

// TestGetCurrentOutageIDs calls api.GetCurrentOutageIDs and checks
// if an int slice is correctly returned.
func TestGetCurrentOutageIDs(t *testing.T) {
	// Set environmental variables to .env
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	// Get outage IDs & data type
	src := GetCurrentOutageIDs()
	src_type := fmt.Sprintf("%T", src)

	// Check for issues
	if src_type != "[]int" {
		t.Fatalf(
			`TestGetCurrentOutageIDs did not return a []int, got type %s of value %d`,
			src_type, src,
		)
	}
}
