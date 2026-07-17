package errdom

import (
	"errors"
	"testing"
)

func TestInvalidSetsValidationKindCodeAndField(t *testing.T) {
	e := Invalid("bad_slug", "slug", "slug must be lowercase")
	if e.Kind != KindValidation || e.Code != "bad_slug" || e.Field != "slug" {
		t.Fatalf("got kind=%d code=%q field=%q", e.Kind, e.Code, e.Field)
	}
}

func TestNotFoundAndConflictKinds(t *testing.T) {
	if NotFound("campaign_missing", "no such campaign").Kind != KindNotFound {
		t.Fatal("not_found kind wrong")
	}
	if Conflict("duplicate_slug", "slug taken").Kind != KindConflict {
		t.Fatal("conflict kind wrong")
	}
}

func TestStatusForMapsEachKind(t *testing.T) {
	cases := map[Kind]int{KindValidation: 422, KindNotFound: 404, KindConflict: 409}
	for k, want := range cases {
		if got := StatusFor(k); got != want {
			t.Fatalf("StatusFor(%d) = %d, want %d", k, got, want)
		}
	}
}

func TestTwoCodesShareOneKind(t *testing.T) {
	dup := Conflict("duplicate_slug", "slug taken")
	deactivated := Conflict("already_deactivated", "link already off")
	if dup.Kind != deactivated.Kind {
		t.Fatal("expected same kind")
	}
	if dup.Code == deactivated.Code {
		t.Fatal("expected different codes")
	}
}

func TestWrapPreservesKindAndCodeAndChains(t *testing.T) {
	child := Invalid("bad_slug", "slug", "nope")
	wrapped := Wrap(child, "links[1].slug", "link 1: bad slug")
	if wrapped.Kind != KindValidation || wrapped.Code != "bad_slug" {
		t.Fatalf("wrap lost identity: kind=%d code=%q", wrapped.Kind, wrapped.Code)
	}
	if wrapped.Field != "links[1].slug" {
		t.Fatalf("wrap lost field context: %q", wrapped.Field)
	}
	var got *DomainError
	if !errors.As(wrapped, &got) {
		t.Fatal("errors.As failed to recover the DomainError")
	}
	if !errors.Is(wrapped.Unwrap(), child) {
		t.Fatal("wrap did not chain the child as cause")
	}
}

func TestCollectAggregatesValidationFailures(t *testing.T) {
	err := Collect(
		Check{"slug", func() error { return Invalid("bad_slug", "slug", "x") }},
		Check{"url", func() error { return Invalid("bad_target_url", "target_url", "y") }},
		Check{"ok", func() error { return nil }},
	)
	var de *DomainError
	if !errors.As(err, &de) || de.Code != "validation_failed" {
		t.Fatalf("expected aggregated validation_failed, got %v", err)
	}
	if len(de.Problems) != 2 {
		t.Fatalf("expected 2 problems, got %d", len(de.Problems))
	}
}

func TestCollectPropagatesNonValidation(t *testing.T) {
	err := Collect(Check{"x", func() error { return Conflict("dup", "taken") }})
	var de *DomainError
	if !errors.As(err, &de) || de.Kind != KindConflict {
		t.Fatalf("non-validation should propagate as-is, got %v", err)
	}
}

func TestInfraIsNotADomainError(t *testing.T) {
	var de *DomainError
	if errors.As(Infra(nil, "db down"), &de) {
		t.Fatal("InfraError must not be recoverable as a DomainError")
	}
}
