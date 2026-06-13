package a

// Encapsulated is a well-formed value object: all fields unexported.
type Encapsulated struct{ v string }

func NewEncapsulated(v string) (Encapsulated, error) { return Encapsulated{v: v}, nil }

// Leaky is a value object that exposes its representation through an exported
// field — the encapsulation leak this analyzer catches.
type Leaky struct {
	Name string // want `value object Leaky exposes exported field Name`
	tag  string
}

func NewLeaky(name, tag string) (Leaky, error) { return Leaky{Name: name, tag: tag}, nil }
