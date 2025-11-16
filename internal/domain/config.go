package domain

import (
	"fmt"
	"regexp"
	"strings"
)

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
	DiffCfg       *DiffConfig
}

type Source struct {
	Path  string
	Label string
	Cmd   []string
}

type DiffConfig struct {
	BaseLabel string
	HeadLabel string
	Cmd       []string
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
	DiffCfg       *rawDiffConfig
}

type rawSource struct {
	Path  string
	Label string
}

type rawDiffConfig struct {
	baseLabel string   `yaml:"baseLabel"`
	headLabel string   `yaml:"headLabel"`
	cmd       []string `yaml:"cmd"`
}

func (c rawDiffConfig) parse() (DiffConfig, []string) {
	var errors []string

	baseLabel := strings.TrimSpace(c.baseLabel)
	if baseLabel == "" {
		errors = append(errors, "base label is empty")
	}

	headLabel := strings.TrimSpace(c.headLabel)
	if headLabel == "" {
		errors = append(errors, "head label is empty")
	}

	if len(c.cmd) == 0 {
		errors = append(errors, "cmd is empty")
	}

	for i, cmdElement := range c.cmd {
		if len(cmdElement) == 0 {
			errors = append(errors, fmt.Sprintf("cmd[%d] is empty", i))
		}
	}

	if len(errors) > 0 {
		var zero DiffConfig
		return zero, errors
	}

	return DiffConfig{
		BaseLabel: baseLabel,
		HeadLabel: headLabel,
		Cmd:       c.cmd,
	}, nil
}
