package domain

type ModuleStatus int

const (
	StatusInSync ModuleStatus = iota
	StatusOutOfSync
	StatusUnknown
)

func (s ModuleStatus) String() string {
	switch s {
	case StatusInSync:
		return "in-sync"
	case StatusOutOfSync:
		return "out-of-sync"
	case StatusUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

func (s ModuleStatus) Symbol() string {
	switch s {
	case StatusInSync:
		return "✓"
	case StatusOutOfSync:
		return "✗"
	case StatusUnknown:
		return "-"
	default:
		return "-"
	}
}

type ModuleResult struct {
	Name   string
	Values map[string]string
	Status ModuleStatus
}

type ComparisonResult struct {
	SourceLabels []string
	Modules      []ModuleResult
}
