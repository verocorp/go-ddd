package ddd

import (
	"fmt"
	"regexp"
)

// seatNumberPattern matches a row of 1-2 digits (1-99) followed by a
// single uppercase letter, e.g. "12A", "1F", "99K".
var seatNumberPattern = regexp.MustCompile(`^[1-9][0-9]?[A-Z]$`)

// SeatNumber identifies a single seat on an aircraft, e.g. "12A". Two
// SeatNumbers built from the same string are the same seat and are
// interchangeable — a SeatNumber has no identity of its own.
type SeatNumber struct {
	value string
}

// NewSeatNumber validates and constructs a SeatNumber.
func NewSeatNumber(value string) (SeatNumber, error) {
	if !seatNumberPattern.MatchString(value) {
		return SeatNumber{}, fmt.Errorf("invalid seat number: %q", value)
	}
	return SeatNumber{value: value}, nil
}

// MustNewSeatNumber panics if value is not a valid seat number. Use only
// with known-valid literals (tests, package-level vars).
func MustNewSeatNumber(value string) SeatNumber {
	s, err := NewSeatNumber(value)
	if err != nil {
		panic(err)
	}
	return s
}

func (s SeatNumber) String() string {
	return s.value
}
