package api

import "strconv"

// isStringInArray returns true if a string is in a
// string array.
func isStringInArray(item string, array []string) bool {
	for _, i := range array {
		if item == i {
			return true
		}
	}
	return false
}

// isStringInArray returns true if an integer is in an
// integer array.
func isIntInArray(item int, array []int) bool {
	for _, i := range array {
		if item == i {
			return true
		}
	}
	return false
}

// isInt returns true if a string is an integer.
func isInt(integer string) bool {
	if _, err := strconv.Atoi(integer); err == nil {
		return true
	}
	return false
}
