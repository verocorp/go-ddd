package a

// Good is a value object whose constructor has the required MustNew pairing.
type Good struct{ v string }

func NewGood(v string) (Good, error) { return Good{v: v}, nil }
func MustNewGood(v string) Good      { g, _ := NewGood(v); return g }

// Bad is a value object whose MustNew helper is missing.
type Bad struct{ v string }

func NewBad(v string) (Bad, error) { return Bad{v: v}, nil } // want `value object Bad: constructor NewBad has no paired MustNewBad`

// makeThing is a factory, not a value-object constructor (suffix != return
// type), so it must NOT be flagged.
func makeThing(v string) (Good, error) { return NewGood(v) }
