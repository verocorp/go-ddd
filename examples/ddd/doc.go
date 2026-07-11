// Package ddd models a small flight-boarding domain: a SeatNumber and a
// FlightNumber (value objects), a Passenger (an entity — a specific
// traveler, tracked by identity, assigned to a seat), and a
// BoardingManifest (an aggregate — the flight's roster, which is the only
// way to seat a Passenger and which enforces that no two Passengers may
// occupy the same seat).
package ddd
