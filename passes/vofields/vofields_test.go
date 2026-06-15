package vofields_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/verocorp/go-ddd/passes/vofields"
)

func TestVOFields(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), vofields.Analyzer, "a")
}
