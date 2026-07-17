// Package errdom is the domain error model for the error-norms worked example
// (the Go mirror of examples/errorspy).
//
// Two-level identity:
//   - Kind is a CLOSED set (validation / not_found / conflict) that the transport
//     boundary maps to an HTTP status. The mapping (StatusFor) is a pure switch;
//     the exhaustive analyzer proves it covers every Kind — add a Kind, forget a
//     case, and the check fails.
//   - Code is an OPEN, stable, machine-readable identifier for the SPECIFIC
//     problem (e.g. "duplicate_slug"). Two codes may share one Kind. Becomes the
//     RFC 9457 "type". Product semantics live here so callers never parse
//     messages or grow the closed Kind set.
//
// Domain code builds a *DomainError only through Invalid / NotFound / Conflict —
// never with a raw Kind. Infrastructure failures use *InfraError, which is NOT a
// domain kind: the boundary maps it to 503. The adapter raises it so nothing
// vendor-shaped crosses into the domain.
package errdom

import (
	"errors"
	"fmt"
)

// Kind is the closed set of domain error kinds. Closed on purpose: StatusFor is
// exhaustive over exactly these values.
type Kind int

const (
	KindValidation Kind = iota
	KindNotFound
	KindConflict
)

// FieldProblem is one field's failure inside an aggregated validation error
// (B6). It becomes an RFC 9457 invalid-params entry.
type FieldProblem struct {
	Code    string
	Field   string
	Message string
}

// DomainError carries an intrinsic Kind (-> status) and a stable Code (-> RFC
// 9457 type). Build it via Invalid / NotFound / Conflict, not a struct literal,
// so a raw Kind can never enter. Problems is non-empty only for an aggregated
// multi-field validation error (B6).
type DomainError struct {
	Kind     Kind
	Code     string
	Field    string // "" if none
	Message  string
	Problems []FieldProblem
	wrapped  error
}

func (e *DomainError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("[%s] (%s) %s", e.Code, e.Field, e.Message)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error { return e.wrapped }

// InfraError is an infrastructure failure (storage unavailable, timeout, driver
// error). NOT a domain kind — the boundary maps it to 503. The adapter raises
// it so nothing vendor-shaped crosses into the domain.
type InfraError struct {
	Message string
	wrapped error
}

func (e *InfraError) Error() string { return e.Message }
func (e *InfraError) Unwrap() error { return e.wrapped }

// Infra builds an *InfraError wrapping cause.
func Infra(cause error, format string, args ...any) *InfraError {
	return &InfraError{Message: fmt.Sprintf(format, args...), wrapped: cause}
}

// Invalid is a validation failure: input the domain refuses. Maps to 422.
func Invalid(code, field, message string) *DomainError {
	return &DomainError{Kind: KindValidation, Code: code, Field: field, Message: message}
}

// NotFound: the asked-for thing is absent. Maps to 404. Detected at the adapter
// but a domain-meaningful outcome.
func NotFound(code, message string) *DomainError {
	return &DomainError{Kind: KindNotFound, Code: code, Message: message}
}

// Conflict: valid input the current state disallows (illegal transition,
// duplicate, lost update). Maps to 409.
func Conflict(code, message string) *DomainError {
	return &DomainError{Kind: KindConflict, Code: code, Message: message}
}

// Wrap re-raises a child domain error with added positional/path context,
// PRESERVING its kind and code — so the boundary still maps it and the client
// still sees the specific problem. Use when a parent knows context a child
// cannot (e.g. which collection index failed). It never invents a new kind.
func Wrap(err *DomainError, field, message string) *DomainError {
	f := field
	if f == "" {
		f = err.Field
	}
	return &DomainError{Kind: err.Kind, Code: err.Code, Field: f, Message: message, wrapped: err}
}

// Check is one named field validation for Collect.
type Check struct {
	Name string
	Fn   func() error
}

// Collect runs each check and AGGREGATES their validation failures (B6): it
// returns ONE validation *DomainError whose Problems list every field that
// failed, instead of stopping at the first. A non-validation failure (e.g. a
// conflict) is not aggregated — it is returned immediately, since folding a 409
// into a 422 batch would be wrong. Returns nil when all pass.
func Collect(checks ...Check) error {
	var problems []FieldProblem
	for _, c := range checks {
		err := c.Fn()
		if err == nil {
			continue
		}
		var de *DomainError
		if errors.As(err, &de) && de.Kind == KindValidation {
			field := de.Field
			if field == "" {
				field = c.Name
			}
			problems = append(problems, FieldProblem{Code: de.Code, Field: field, Message: de.Message})
			continue
		}
		return err // non-validation: propagate as-is
	}
	if len(problems) > 0 {
		return &DomainError{
			Kind:     KindValidation,
			Code:     "validation_failed",
			Message:  "one or more fields are invalid",
			Problems: problems,
		}
	}
	return nil
}

// StatusFor is the pure Kind -> HTTP status mapper. This is the ONE place the
// closed set is enforced: the exhaustive analyzer flags a missing case, and the
// panic is the runtime backstop. Runtime error recovery (errors.As) happens at
// the boundary; the exhaustiveness check fires HERE, on the typed Kind.
func StatusFor(k Kind) int {
	switch k {
	case KindValidation:
		return 422
	case KindNotFound:
		return 404
	case KindConflict:
		return 409
	}
	panic(fmt.Sprintf("unhandled kind: %d", k))
}
