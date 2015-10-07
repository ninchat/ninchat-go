package ninchatmodel

// Change
type Change int

const (
	unchanged Change = iota
	Added
	Updated
	Removed
)

func (c Change) String() string {
	switch c {
	case unchanged:
		return "unchanged"

	case Added:
		return "added"

	case Updated:
		return "updated"

	case Removed:
		return "removed"

	default:
		return "invalid change value"
	}
}
