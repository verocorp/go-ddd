package primitive

// Temperature as a bare float64 carries no scale. Two of the failures a value
// object prevents show up here:
//
//   - Equality: 0.0 (°C) and 273.15 (K) are the SAME temperature but compare
//     unequal, and nothing in the type records which scale a number is in.
//   - Validation: a physically impossible temperature below absolute zero
//     (negative Kelvin) flows through unchecked.

// TempEqual compares two temperatures with ==. Wrong across scales: it reports
// 0°C and 273.15K as different.
func TempEqual(a, b float64) bool {
	return a == b
}

// AsKelvin does no validation. A negative Kelvin (below absolute zero) passes
// through silently.
func AsKelvin(k float64) float64 {
	return k
}
