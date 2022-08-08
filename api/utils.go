package api

import (
	"strconv"
	"strings"
)

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

// GetNWordsRemovedFromStart returns a string after removing
// n words from the start of a string.
func GetNWordsRemovedFromStart(
	s string, sep string, n int) string {
	return strings.Join(strings.Split(s, sep)[n:], sep)
}

// StripNonLetters is a utility function that strips all non-letters
// from the given string (str).
func StripNonLetters(str string) string {
	s := []byte(str)
	n := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') {
			s[n] = b
			n++
		}
	}
	return string(s[:n])
}
