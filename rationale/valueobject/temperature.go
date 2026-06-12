package valueobject

import "fmt"

// Temperature stores Kelvin internally, so values created in different scales
// compare correctly (0°C equals 273.15K), and construction rejects anything
// below absolute zero. The scale a caller used is normalized away; the
// invariant is enforced in one place.
type Temperature struct{ kelvin float64 }

// FromKelvin validates against absolute zero.
func FromKelvin(k float64) (Temperature, error) {
	if k < 0 {
		return Temperature{}, fmt.Errorf("%.2fK is below absolute zero", k)
	}
	return Temperature{kelvin: k}, nil
}

// FromCelsius converts to Kelvin, then validates.
func FromCelsius(c float64) (Temperature, error) {
	return FromKelvin(c + 273.15)
}

// Equal compares by physical value, not by the scale or number used to build it.
func (t Temperature) Equal(other Temperature) bool { return t.kelvin == other.kelvin }

func (t Temperature) Kelvin() float64 { return t.kelvin }
