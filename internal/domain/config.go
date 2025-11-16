package domain

import "regexp"

type Config struct {
	Version        int
	CompareModules CompareModules
}

type CompareModules struct {
	Comparisons []Comparison
	ValueRegex  *regexp.Regexp
}

type Comparison struct {
	Name          string
	AttributeKey  string
	Sources       []Source
	IgnoreModules []string
	ValueRegex    *regexp.Regexp
}

type Source struct {
	Path  string
	Label string
}

type OutputFormat uint8

const (
	StdoutOutput OutputFormat = iota
	HtmlOutput
)

func ParseOutputFormat(value string) (OutputFormat, bool) {
	switch value {
	case "stdout":
		return StdoutOutput, true
	case "html":
		return HtmlOutput, true
	default:
		return StdoutOutput, false
	}
}

func GetOutputFormatValues() []string {
	return []string{"stdout", "html"}
}

type rawConfig struct {
	Version        int
	CompareModules rawCompareModules `yaml:"compareModules"`
}

type rawCompareModules struct {
	Comparisons []rawComparison `yaml:"comparisons"`
	ValueRegex  string          `yaml:"valueRegex,omitempty"`
}

type rawComparison struct {
	Name          string
	AttributeKey  string      `yaml:"attributeKey"`
	Sources       []rawSource `yaml:"sources"`
	IgnoreModules []string    `yaml:"ignoreModules,omitempty"`
	ValueRegex    string      `yaml:"valueRegex,omitempty"`
}

type rawSource struct {
	Path  string
	Label string
}
