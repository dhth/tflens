package domain

import "regexp"

// raw

type RawConfig struct {
	Version        int
	CompareModules RawCompareModules `yaml:"compareModules"`
}

type RawCompareModules struct {
	Comparisons []RawComparison `yaml:"comparisons"`
	ValueRegex  string          `yaml:"valueRegex,omitempty"`
}

type RawComparison struct {
	Name          string
	AttributeKey  string      `yaml:"attributeKey"`
	Sources       []RawSource `yaml:"sources"`
	IgnoreModules []string    `yaml:"ignoreModules,omitempty"`
	ValueRegex    string      `yaml:"valueRegex,omitempty"`
}

type RawSource struct {
	Path  string
	Label string
}

// validated

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
