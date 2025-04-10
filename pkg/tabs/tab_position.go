package tabs

type TabPosition string

const (
	AlignTop    TabPosition = "tabgroup:top"
	AlignBottom TabPosition = "tabgroup:bottom"
)

// IsValid checks if a TabGroup is a valid value.
func (t TabPosition) IsValid() bool {
	switch t {
	case AlignTop, AlignBottom:
		return true
	default:
		return false
	}
}
