package api

import (
	"testing"
)

// TestDBWaterOutageCol calls api.DBWaterOutageCol and returns
// a reference for a column.
func TestDBWaterOutageCol(t *testing.T) {
	outage := DBWaterOutage{}
	json_col := "outage_id"

	// Get column name for "outage_id"
	actual_return := DBWaterOutageCol(json_col, &outage)
	expected_return := &outage.OutageID

	// Check for issues
	if actual_return != expected_return {
		t.Fatalf(
			`TestDBWaterOutageCol did not return %v, got %v`,
			expected_return, actual_return,
		)
	}
}
