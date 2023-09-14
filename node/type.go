package node

import "fmt"

type NodeType int

const (
	Embedded NodeType = iota
	Light
	Bridge
	Full
)

// String returns the string representation of the NodeType.
func (n NodeType) String() string {
	switch n {
	case Embedded:
		return "Embedded"
	case Light:
		return "Light"
	case Bridge:
		return "Bridge"
	case Full:
		return "Full"
	default:
		return fmt.Sprintf("Unknown NodeType: %d", int(n))
	}
}

// FromString takes a string and converts it to its corresponding NodeType.
// Returns an error if the string doesn't match any known NodeType.
func FromString(s string) (NodeType, error) {
	switch s {
	case "Embedded":
		return Embedded, nil
	case "Light":
		return Light, nil
	case "Bridge":
		return Bridge, nil
	case "Full":
		return Full, nil
	default:
		return -1, fmt.Errorf("unknown NodeType: %s", s)
	}
}
