package api

import "testing"

// TestSetSearchWhere tests SetSearchWhere and checks if the input
// is served as a correct WHERE string.
func TestSetSearchWhere(t *testing.T) {
	tests := [][]string{
		{"15899", "outage_id = '15899'"},
		{"21 Uranus Street", `(lower(suburb) LIKE lower('%%Uranus Street%%')
		OR lower(street) LIKE lower('%%Uranus Street%%') 
		OR lower(suburb) LIKE lower('%%Uranus Street%%')
		OR lower(street) LIKE lower('%%Uranus Street%%'))`},
	}

	for _, expected := range tests {
		query := Query{}
		query.SetSearchWhere([]string{expected[0]})
		actual := query.Wheres[0]

		if StripNonLetters(actual) != StripNonLetters(expected[1]) {
			t.Fatalf(
				`TestSearchWhere did not return 
				%s 
				got 
				%s`,
				expected[1], actual,
			)
		}
	}
}
