package stringequality_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/chrisconley/go-ddd/passes/stringequality"
)

func TestStringEquality(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), stringequality.Analyzer, "a", "ext")
}
