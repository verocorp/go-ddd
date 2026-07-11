package ddd

import (
	"fmt"
	"strings"
)

// PassengerID identifies one specific traveler — e.g. a booking reference
// — independent of their name or seat assignment. It is the identity a
// Passenger entity carries for its whole lifecycle.
type PassengerID struct {
	value string
}

// NewPassengerID validates and constructs a PassengerID.
func NewPassengerID(value string) (PassengerID, error) {
	if strings.TrimSpace(value) == "" {
		return PassengerID{}, fmt.Errorf("passenger ID must not be empty")
	}
	return PassengerID{value: value}, nil
}

// MustNewPassengerID panics if value is not a valid PassengerID. Use only
// with known-valid literals (tests, package-level vars).
func MustNewPassengerID(value string) PassengerID {
	id, err := NewPassengerID(value)
	if err != nil {
		panic(err)
	}
	return id
}

func (p PassengerID) String() string {
	return p.value
}
