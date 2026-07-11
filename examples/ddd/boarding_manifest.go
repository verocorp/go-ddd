package ddd

import "fmt"

// BoardingManifest is the flight's boarding manifest: the aggregate root
// that owns the flight's Passengers and is the only way to seat one. Its
// invariant — no two Passengers may occupy the same seat — is enforced
// here, both at construction and on every AddPassenger call, so a manifest
// with a duplicate seat assignment is unrepresentable. Boarding is a
// lifecycle: the manifest starts with whatever passengers are already
// checked in and grows as more board, so it is mutable and its
// transitions re-establish the invariant.
type BoardingManifest struct {
	flight     FlightNumber
	passengers []Passenger
	_          [0]func() // non-comparable — a manifest is never compared by value
}

// BoardingManifestSpec carries construction data across the layer
// boundary.
type BoardingManifestSpec struct {
	Flight     string
	Passengers []PassengerSpec
}

// NewBoardingManifest validates spec and constructs a BoardingManifest.
// Passengers are built and seated left to right; the first seat collision
// is rejected with a contextful error.
func NewBoardingManifest(spec BoardingManifestSpec) (BoardingManifest, error) {
	flight, err := NewFlightNumber(spec.Flight)
	if err != nil {
		return BoardingManifest{}, fmt.Errorf("invalid flight number: %w", err)
	}

	manifest := BoardingManifest{flight: flight}
	for i, pSpec := range spec.Passengers {
		passenger, err := NewPassenger(pSpec)
		if err != nil {
			return BoardingManifest{}, fmt.Errorf("invalid passenger at index %d: %w", i, err)
		}
		if err := manifest.addPassenger(passenger); err != nil {
			return BoardingManifest{}, err
		}
	}
	return manifest, nil
}

// Flight returns the flight this manifest belongs to.
func (m BoardingManifest) Flight() FlightNumber { return m.flight }

// Passengers returns a defensive copy of the boarded passengers; mutating
// the result never affects the manifest.
func (m BoardingManifest) Passengers() []Passenger {
	out := make([]Passenger, len(m.passengers))
	copy(out, m.passengers)
	return out
}

// AddPassenger boards one more passenger, re-establishing the
// no-shared-seat invariant. It errors — leaving the manifest unchanged —
// if the passenger's seat is already occupied.
func (m *BoardingManifest) AddPassenger(p Passenger) error {
	return m.addPassenger(p)
}

// addPassenger is the single site that appends to m.passengers, so the
// invariant check and the mutation can never drift apart.
func (m *BoardingManifest) addPassenger(p Passenger) error {
	for _, existing := range m.passengers {
		if existing.Seat() == p.Seat() {
			return fmt.Errorf("seat %s is already occupied", p.Seat())
		}
	}
	m.passengers = append(m.passengers, p)
	return nil
}
