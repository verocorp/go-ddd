package ddd

import "testing"

func TestNewPassengerName_Valid(t *testing.T) {
	valid := []string{"Ada Lovelace", "X", "O'Brien-Smith"}
	for _, v := range valid {
		if _, err := NewPassengerName(v); err != nil {
			t.Errorf("NewPassengerName(%q) returned unexpected error: %v", v, err)
		}
	}
}

func TestNewPassengerName_InvalidRejected(t *testing.T) {
	invalid := []string{"", "   ", "\t\n"}
	for _, v := range invalid {
		if _, err := NewPassengerName(v); err == nil {
			t.Errorf("NewPassengerName(%q) = nil error, want error", v)
		}
	}
}

func TestPassengerName_Equality(t *testing.T) {
	a := MustNewPassengerName("Ada Lovelace")
	b := MustNewPassengerName("Ada Lovelace")
	c := MustNewPassengerName("Grace Hopper")

	if a != b {
		t.Error("passenger names built from the same value must be equal")
	}
	if a == c {
		t.Error("passenger names built from different values must not be equal")
	}
}

func TestMustNewPassengerName_PanicsOnInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("MustNewPassengerName did not panic on invalid input")
		}
	}()
	MustNewPassengerName("")
}

func TestPassengerName_String(t *testing.T) {
	n := MustNewPassengerName("Ada Lovelace")
	if got, want := n.String(), "Ada Lovelace"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
