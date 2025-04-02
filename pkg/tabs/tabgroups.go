package tabs

type TabGroup string

const (
	AlignTop   TabGroup = "tabgroup:main"
	AlignBottom TabGroup = "tabgroup:bottom"
)

// IsValid checks if a TabGroup is a valid value.
func (t TabGroup) IsValid() bool {
	switch t {
	case AlignTop, AlignBottom:
		return true
	default:
		return false
	}
}
