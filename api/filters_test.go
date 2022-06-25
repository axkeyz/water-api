package api

import (
	"testing"
)

// TestIsFilterableOutage calls api.IsFilterableOutage and checks
// whether the correct bool for the given param is returned.
func TestIsFilterableOutage(t *testing.T) {
	// Get column name for "outage_id"
	test1 := IsFilterableOutage("outage_id")
	test2 := IsFilterableOutage("outageid")

	// Check for issues
	if !test1 && test2 {
		t.Fatalf(
			`TestIsFilterableOutage did not return true, false got %v, %v`,
			test1, test2,
		)
	}
}

// TestIsFilterableCountOutage calls api.IsFilterableCountOutage and checks
// whether the correct bool for the given param is returned.
func TestIsFilterableCountOutage(t *testing.T) {
	// Get column name for "outage_id"
	test1 := IsFilterableCountOutage("total_hours")
	test2 := IsFilterableCountOutage("outage_id")

	// Check for issues
	if !test1 && test2 {
		t.Fatalf(
			`TestIsFilterableCountOutage did not return true, false got %v, %v`,
			test1, test2,
		)
	}
}

// TestIsCurrentOutage calls api.IsCurrentOutage and checks whether the
// an outage is correctly identified as a currently active one.
func TestIsCurrentOutage(t *testing.T) {
	// Create random []int of outage iids
	outage_ids := []int{5344, 23114, 3592, 54675, 54677, 12954}

	// Get column name for "outage_id"
	test1 := IsCurrentOutage(54677, outage_ids)
	test2 := IsCurrentOutage(54676, outage_ids)

	// Check for issues
	if !test1 && test2 {
		t.Fatalf(
			`TestIsCurrentOutage did not return true, false got %v, %v`,
			test1, test2,
		)
	}
}
