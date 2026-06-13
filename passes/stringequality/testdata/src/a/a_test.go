package a

import "testing"

// TestFoo_String is the accessor test — .String() is allowed here.
func TestFoo_String(t *testing.T) {
	f := Foo{v: "x"}
	_ = f.String()
}

// TestFoo_NewFoo is not an accessor test, so its .String() is flagged.
func TestFoo_NewFoo(t *testing.T) {
	f := Foo{v: "x"}
	_ = f.String() // want `TestFoo_NewFoo calls \.String\(\) outside a Test\*_String accessor test`
}

// Two .String() calls in one function produce two diagnostics.
func TestFoo_Pair(t *testing.T) {
	a := Foo{v: "a"}
	b := Foo{v: "b"}
	_ = a.String() // want `TestFoo_Pair calls \.String\(\) outside`
	_ = b.String() // want `TestFoo_Pair calls \.String\(\) outside`
}

// ToString is not String, so it is ignored.
func TestFoo_Other(t *testing.T) {
	f := Foo{v: "x"}
	_ = f.ToString()
}
