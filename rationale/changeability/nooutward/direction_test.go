// The decision-3 PRIMARY discriminator: the D3a direction guard. Per ../SCORING.md
// the harm of "no outward representation" is dependency direction, not fan-out, so
// the assertion that makes the rule specific is a COMPILE guard, not a count:
//
//	go build ./nav            → succeeds (the pure domain compiles)
//	go build -tags leak ./nav → FAILS   (import cycle: domain importing outward)
//
// The outward layer (telemetry) imports the domain (nav) to map it — the correct
// inward direction — so any attempt to make the domain emit its own outward
// representation (nav importing telemetry, in leak_bug.go under -tags leak) is an
// import cycle the compiler rejects. This is the decision-3 analog of the anchor's
// subst_bug.go: a structural property the language enforces at N=1. It holds even
// where the D3b fan-out is tied by a cheaper structure, which is why it, not the
// fan-out, is where the direction rule earns its place.
package nooutward_test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestDirectionGuard_D3a_DomainCannotEmit(t *testing.T) {
	// Positive control: the pure domain compiles on its own.
	if out, err := exec.Command("go", "build", "./nav").CombinedOutput(); err != nil {
		t.Fatalf("pure domain ./nav must compile, but build failed:\n%s", out)
	}

	// The guard: making the domain emit its outward representation is an import
	// cycle — it must NOT compile.
	out, err := exec.Command("go", "build", "-tags", "leak", "./nav").CombinedOutput()
	if err == nil {
		t.Fatalf("`go build -tags leak ./nav` compiled, but the domain importing the\n" +
			"outward layer must be an import cycle — the D3a direction guard is broken")
	}
	if !strings.Contains(string(out), "import cycle") {
		t.Errorf("build failed, but not on an import cycle — the failure must be the\n"+
			"dependency-direction violation, not something incidental:\n%s", out)
	}
}
