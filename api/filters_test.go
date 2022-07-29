package api

import (
	"testing"
)

// TestIsFilterableParam calls api.IsFilterableParam and checks
// whether the correct bool for the given param is returned.
func TestIsFilterableParam(t *testing.T) {
	// Get column name for "outage_id"
	test1 := IsFilterableParam("outage_id")
	test2 := IsFilterableParam("outageid")

	// Check for issues
	if !test1 && test2 {
		t.Fatalf(
			`TestIsFilterableParam did not return true, false got %v, %v`,
			test1, test2,
		)
	}
}

// TestIsFilterableCountParam calls api.IsFilterableCountParam and checks
// whether the correct bool for the given param is returned.
func TestIsFilterableCountParam(t *testing.T) {
	// Get column name for "outage_id"
	test1 := IsFilterableCountParam("total_hours")
	test2 := IsFilterableCountParam("outage_id")

	// Check for issues
	if !test1 && test2 {
		t.Fatalf(
			`TestIsFilterableCountParam did not return true, false got %v, %v`,
			test1, test2,
		)
	}
}

// TestIsCurrentOutageID calls api.IsCurrentOutageID and checks whether the
// an outage is correctly identified as a currently active one.
func TestIsCurrentOutageID(t *testing.T) {
	// Create random []int of outage iids
	outage_ids := []int{5344, 23114, 3592, 54675, 54677, 12954}

	// Get column name for "outage_id"
	test1 := IsCurrentOutageID(54677, outage_ids)
	test2 := IsCurrentOutageID(54676, outage_ids)

	// Check for issues
	if !test1 && test2 {
		t.Fatalf(
			`TestIsCurrentOutageID did not return true, false got %v, %v`,
			test1, test2,
		)
	}
}
