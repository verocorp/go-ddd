package primitive

// Tags-as-a-raw-map: telemetry labels passed around as map[string]string with
// no wrapper. Reads are a direct map index (zero allocation), but nothing stops
// a caller from mutating shared state. This is the fast, unsafe baseline the
// collection-VO benchmark compares against.

// TagValue reads a key directly from the caller's map.
func TagValue(tags map[string]string, key string) string {
	return tags[key]
}
