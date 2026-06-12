package changeability_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestCoverageMatrix_NoSilentGaps keeps coverage.md honest: every checker in
// ../cmd must appear in the matrix, and every Test* the matrix names must exist
// in this package. It tolerates the ❌/⚠️ rows by design (the rationale is
// broader than the enforcement); it forbids a SILENT gap — a checker missing
// from the matrix, or a dangling test reference that rotted.
func TestCoverageMatrix_NoSilentGaps(t *testing.T) {
	matrix, err := os.ReadFile("coverage.md")
	if err != nil {
		t.Fatalf("read coverage.md: %v", err)
	}
	content := string(matrix)

	// 1. Every cmd/check* checker is named in the matrix.
	entries, err := os.ReadDir("../cmd")
	if err != nil {
		t.Fatalf("read ../cmd: %v", err)
	}
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "check") {
			if !strings.Contains(content, e.Name()) {
				t.Errorf("checker %q has no row in coverage.md (silent gap)", e.Name())
			}
		}
	}

	// 2. Every concrete Test* the matrix references actually exists here.
	// (Patterns like "Test*_Equality" carry an asterisk and are not matched.)
	var src strings.Builder
	testFiles, _ := filepath.Glob("*_test.go")
	for _, f := range testFiles {
		b, _ := os.ReadFile(f)
		src.Write(b)
	}
	srcStr := src.String()

	for _, name := range regexp.MustCompile(`Test[A-Za-z0-9_]+`).FindAllString(content, -1) {
		if !strings.Contains(srcStr, "func "+name+"(") {
			t.Errorf("coverage.md references %q but no such test exists (dangling reference)", name)
		}
	}
}
