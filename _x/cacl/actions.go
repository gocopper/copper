package cacl

// Common actions used for managing permissions
const (
	ActionRead  = "READ"
	ActionWrite = "WRITE"
)

// Common combination of actions used for managing permission
var (
	ActionReadWrite = []string{ActionRead, ActionWrite}
	ActionReadOnly  = []string{ActionRead}
)
