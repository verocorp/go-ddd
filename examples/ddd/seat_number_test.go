package ddd

import "testing"

func TestNewSeatNumber_Valid(t *testing.T) {
	valid := []string{"1A", "12A", "99K", "5Z"}
	for _, v := range valid {
		if _, err := NewSeatNumber(v); err != nil {
			t.Errorf("NewSeatNumber(%q) returned unexpected error: %v", v, err)
		}
	}
}

func TestNewSeatNumber_InvalidRejected(t *testing.T) {
	invalid := []string{
		"",     // empty
		"12",   // missing letter
		"A12",  // letter first
		"12a",  // lowercase letter
		"0A",   // row must be >= 1
		"100A", // row too large (3 digits)
		"12AA", // two letters
		" 12A", // leading whitespace
	}
	for _, v := range invalid {
		if _, err := NewSeatNumber(v); err == nil {
			t.Errorf("NewSeatNumber(%q) = nil error, want error", v)
		}
	}
}

func TestSeatNumber_Equality(t *testing.T) {
	a := MustNewSeatNumber("12A")
	b := MustNewSeatNumber("12A")
	c := MustNewSeatNumber("12B")

	if a != b {
		t.Error("seat numbers built from the same value must be equal")
	}
	if a == c {
		t.Error("seat numbers built from different values must not be equal")
	}
}

func TestMustNewSeatNumber_PanicsOnInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("MustNewSeatNumber did not panic on invalid input")
		}
	}()
	MustNewSeatNumber("not-a-seat")
}

func TestSeatNumber_String(t *testing.T) {
	s := MustNewSeatNumber("12A")
	if got, want := s.String(), "12A"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
