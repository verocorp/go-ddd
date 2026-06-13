package a

// Foo is a value object with both a String() (display) and a ToString().
type Foo struct{ v string }

func (f Foo) String() string   { return f.v }
func (f Foo) ToString() string { return f.v }

// Render calls .String() in production code, which must NOT be flagged — only
// test files are scanned.
func Render(f Foo) string { return f.String() }
