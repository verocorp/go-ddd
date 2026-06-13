package ext_test

import (
	"testing"

	"ext"
)

// An external test package (package ext_test) is still scanned: .String()
// outside a Test*_String accessor test is flagged here too.
func TestFoo(t *testing.T) {
	f := ext.Foo{}
	_ = f.String() // want `TestFoo calls \.String\(\) outside`
}
