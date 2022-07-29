// filters_utils_test.go contains tests that test filters_utils_test.go
package api

import (
	"fmt"
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

// TestIsDateParam calls api.IsDateParam and checks if
func TestIsDateParam(t *testing.T) {
	tests := map[string][]string{
		"after_end_date":            {"end_date >=", "true"},
		"before_end_date":           {"end_date <=", "true"},
		"outage_type":               {"", "false"},
		"fake_date":                 {"", "false"},
		"super_fake_after_end_date": {"", "false"},
	}

	for test, expected := range tests {
		actualBool, actualColumn := IsDateParam(test)
		actualIsDate := fmt.Sprintf("%v", actualBool)

		if actualColumn != expected[0] &&
			actualIsDate == expected[1] {
			t.Fatalf(
				`TestIsDateParam did not return %s, %s got %s, %v`,
				expected[0], expected[1], actualColumn, actualIsDate,
			)
		}
	}
}

// TestGetEquationSignedColumn calls api.GetEquationSignedColumn and
// checks whether the correct column and operator sign is returned.
func TestGetEquationSignedColumn(t *testing.T) {
	// Get column name for "outage_id"
	tests := map[string]string{
		GetEquationSignedColumn("after_end_date", 1):    "end_date >=",
		GetEquationSignedColumn("before_start_date", 1): "start_date <=",
		GetEquationSignedColumn("outage_id", 0):         "outage_id =",
	}

	for actual, expected := range tests {
		if actual != expected {
			t.Fatalf(
				`TestGetEquationSignedColumn did not return %s got %s`,
				expected, actual,
			)
		}
	}
}

// TestGetSQLEquationSigns calls api.GetSQLEquationSigns and checks
// if the correct SQL sign is returned when the keyword is used.
func TestGetSQLEquationSigns(t *testing.T) {
	tests := map[string]string{
		"after": ">=",
		"like":  "LIKE",
		"other": "=",
	}

	for test, expected := range tests {
		if GetSQLEquationSigns(test) != expected {
			t.Fatalf(
				`TestGetSQLEquationSigns did not return %s got %s`,
				expected, GetSQLEquationSigns(test),
			)
		}
	}
}

// TestGetSQLCondition calls api.TestGetSQLCondition and checks
// if the correct AND/OR SQL condition is returned.
func TestGetSQLCondition(t *testing.T) {
	tests := map[string]string{
		"":        " AND ",
		"true":    " AND ",
		"false":   " OR ",
		"uranium": " OR ",
	}

	for test, expected := range tests {
		if GetSQLCondition(test) != expected {
			t.Fatalf(
				`TestGetSQLCondition did not return %s got %s`,
				expected, GetSQLCondition(test),
			)
		}
	}
}
