package ddd

import (
	"fmt"
	"regexp"
)

// flightNumberPattern matches an airline code (2-3 uppercase letters)
// followed by 1-4 digits, e.g. "AA100", "DL2703", "UAL204".
var flightNumberPattern = regexp.MustCompile(`^[A-Z]{2,3}[0-9]{1,4}$`)

// FlightNumber identifies a scheduled flight, e.g. "DL2703". Two
// FlightNumbers built from the same string are the same flight and are
// interchangeable.
type FlightNumber struct {
	value string
}

// NewFlightNumber validates and constructs a FlightNumber.
func NewFlightNumber(value string) (FlightNumber, error) {
	if !flightNumberPattern.MatchString(value) {
		return FlightNumber{}, fmt.Errorf("invalid flight number: %q", value)
	}
	return FlightNumber{value: value}, nil
}

// MustNewFlightNumber panics if value is not a valid flight number. Use
// only with known-valid literals (tests, package-level vars).
func MustNewFlightNumber(value string) FlightNumber {
	f, err := NewFlightNumber(value)
	if err != nil {
		panic(err)
	}
	return f
}

func (f FlightNumber) String() string {
	return f.value
}
