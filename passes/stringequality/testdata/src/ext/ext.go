package ext

// Foo is a value object whose .String() is display-only.
type Foo struct{ v string }

func (f Foo) String() string { return f.v }
