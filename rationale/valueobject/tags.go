package valueobject

// Tags is a collection value object over telemetry labels. It enforces the
// defensive-copy rule: copy on construction and copy on the bulk accessor, so a
// caller can never mutate the internal map. This is SAFE — and it is the
// clearest place a VO costs real runtime: every defensive copy is an allocation
// the raw map never pays. The benchmark measures exactly that trade.
type Tags struct {
	values map[string]string
}

// NewTags copies the input so later caller mutation can't reach inside.
func NewTags(m map[string]string) Tags {
	return Tags{values: copyMap(m)}
}

// Get is a single-key read — no copy, no allocation.
func (t Tags) Get(key string) string {
	return t.values[key]
}

// All returns a defensive copy. Safe, but allocates on every call.
func (t Tags) All() map[string]string {
	return copyMap(t.values)
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
