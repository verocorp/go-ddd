package ddd

import "testing"

func TestNewPassengerID_Valid(t *testing.T) {
	valid := []string{"PNR-ABC123", "1", "traveler-42"}
	for _, v := range valid {
		if _, err := NewPassengerID(v); err != nil {
			t.Errorf("NewPassengerID(%q) returned unexpected error: %v", v, err)
		}
	}
}

func TestNewPassengerID_InvalidRejected(t *testing.T) {
	invalid := []string{"", "   ", "\t"}
	for _, v := range invalid {
		if _, err := NewPassengerID(v); err == nil {
			t.Errorf("NewPassengerID(%q) = nil error, want error", v)
		}
	}
}

func TestPassengerID_Equality(t *testing.T) {
	a := MustNewPassengerID("PNR-ABC123")
	b := MustNewPassengerID("PNR-ABC123")
	c := MustNewPassengerID("PNR-XYZ999")

	if a != b {
		t.Error("passenger IDs built from the same value must be equal")
	}
	if a == c {
		t.Error("passenger IDs built from different values must not be equal")
	}
}

func TestMustNewPassengerID_PanicsOnInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("MustNewPassengerID did not panic on invalid input")
		}
	}()
	MustNewPassengerID("")
}

func TestPassengerID_String(t *testing.T) {
	id := MustNewPassengerID("PNR-ABC123")
	if got, want := id.String(), "PNR-ABC123"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
