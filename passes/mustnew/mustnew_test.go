package mustnew_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/verocorp/go-ddd/passes/mustnew"
)

func TestMustNew(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), mustnew.Analyzer, "a")
}
