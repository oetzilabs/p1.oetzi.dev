package tabs

type TabGroup string

const (
	TabGroupsMain   TabGroup = "tabgroup:main"
	TabGroupsBottom TabGroup = "tabgroup:bottom"
)

// IsValid checks if a TabGroup is a valid value.
func (t TabGroup) IsValid() bool {
	switch t {
	case TabGroupsMain, TabGroupsBottom:
		return true
	default:
		return false
	}
}
