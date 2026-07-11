package ddd

import "fmt"

// Passenger is a specific traveler assigned to a specific seat. The system
// tracks a Passenger by identity: two Passengers are the same traveler iff
// they share a PassengerID, regardless of name or seat (a corrected name,
// a reassigned seat, is still the same passenger); two distinct travelers
// who happen to share a name are NOT the same Passenger. A Passenger
// records a completed seat assignment — it is a fact, not a lifecycle — so
// it is immutable; a reassignment is a new Passenger value with the same
// ID.
type Passenger struct {
	id   PassengerID
	name PassengerName
	seat SeatNumber
}

// PassengerSpec carries construction data across the layer boundary.
type PassengerSpec struct {
	ID   string
	Name string
	Seat string
}

// NewPassenger validates spec and constructs a Passenger. Each child value
// object performs its own validation; NewPassenger only wraps the errors
// with field context — it never re-checks a child's rules itself.
func NewPassenger(spec PassengerSpec) (Passenger, error) {
	id, err := NewPassengerID(spec.ID)
	if err != nil {
		return Passenger{}, fmt.Errorf("invalid passenger ID: %w", err)
	}
	name, err := NewPassengerName(spec.Name)
	if err != nil {
		return Passenger{}, fmt.Errorf("invalid passenger name: %w", err)
	}
	seat, err := NewSeatNumber(spec.Seat)
	if err != nil {
		return Passenger{}, fmt.Errorf("invalid seat number: %w", err)
	}
	return Passenger{id: id, name: name, seat: seat}, nil
}

// ID returns the Passenger's identity.
func (p Passenger) ID() PassengerID { return p.id }

// Name returns the Passenger's name.
func (p Passenger) Name() PassengerName { return p.name }

// Seat returns the Passenger's assigned seat.
func (p Passenger) Seat() SeatNumber { return p.seat }

// Equal reports whether two Passengers are the same traveler. Identity is
// by PassengerID alone — never by name or seat — so Equal, not native
// `==`, is the correct comparison for Passenger identity.
func (p Passenger) Equal(other Passenger) bool {
	return p.id == other.id
}
