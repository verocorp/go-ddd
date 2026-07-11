package ddd

import (
	"fmt"
	"strings"
)

// PassengerName is a traveler's name as it appears on the manifest.
type PassengerName struct {
	value string
}

// NewPassengerName validates and constructs a PassengerName.
func NewPassengerName(value string) (PassengerName, error) {
	if strings.TrimSpace(value) == "" {
		return PassengerName{}, fmt.Errorf("passenger name must not be empty")
	}
	return PassengerName{value: value}, nil
}

// MustNewPassengerName panics if value is not a valid PassengerName. Use
// only with known-valid literals (tests, package-level vars).
func MustNewPassengerName(value string) PassengerName {
	n, err := NewPassengerName(value)
	if err != nil {
		panic(err)
	}
	return n
}

func (n PassengerName) String() string {
	return n.value
}
