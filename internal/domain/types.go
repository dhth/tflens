package domain

type ModuleStatus int

const (
	StatusInSync ModuleStatus = iota
	StatusOutOfSync
	StatusNotApplicable
)

func (s ModuleStatus) String() string {
	switch s {
	case StatusInSync:
		return "in_sync"
	case StatusOutOfSync:
		return "out_of_sync"
	case StatusNotApplicable:
		return "not_applicable"
	default:
		return "not_applicable"
	}
}

func (s ModuleStatus) Symbol() string {
	switch s {
	case StatusInSync:
		return "✓"
	case StatusOutOfSync:
		return "✗"
	case StatusNotApplicable:
		return "-"
	default:
		return "-"
	}
}

type DiffResult struct {
	Output    []byte
	BaseLabel string
	HeadLabel string
	BaseRef   string
	HeadRef   string
}

type ModuleResult struct {
	Name       string
	Values     map[string]string
	Status     ModuleStatus
	DiffResult *DiffResult `yaml:"diffResult,omitempty"`
}

type ComparisonResult struct {
	SourceLabels []string
	Modules      []ModuleResult
}
