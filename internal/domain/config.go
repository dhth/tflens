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
	AttributeKey  string         `yaml:"attributeKey"`
	Sources       []rawSource    `yaml:"sources"`
	IgnoreModules []string       `yaml:"ignoreModules,omitempty"`
	ValueRegex    string         `yaml:"valueRegex,omitempty"`
	DiffCfg       *rawDiffConfig `yaml:"diffConfig"`
}

type rawSource struct {
	Path  string
	Label string
}

type rawDiffConfig struct {
	BaseLabel string   `yaml:"baseLabel"`
	HeadLabel string   `yaml:"headLabel"`
	Cmd       []string `yaml:"cmd"`
}

func (c rawDiffConfig) parse(labels map[string]struct{}) (DiffConfig, []string) {
	var errors []string

	baseLabel := strings.TrimSpace(c.BaseLabel)
	if len(baseLabel) == 0 {
		errors = append(errors, "base label is empty")
	} else {
		_, baseLabelOk := labels[baseLabel]
		if !baseLabelOk {
			errors = append(errors, fmt.Sprintf("base label %q is not in the list of defined labels", baseLabel))
		}
	}

	headLabel := strings.TrimSpace(c.HeadLabel)
	if len(headLabel) == 0 {
		errors = append(errors, "head label is empty")
	} else {
		_, headLabelOk := labels[headLabel]
		if !headLabelOk {
			errors = append(errors, fmt.Sprintf("head label %q is not in the list of defined labels", headLabel))
		}
	}

	if len(c.Cmd) == 0 {
		errors = append(errors, "cmd is empty")
	}

	trimmedCmd := make([]string, 0, len(c.Cmd))
	for i, cmdElement := range c.Cmd {
		trimmedElement := strings.TrimSpace(cmdElement)
		if len(trimmedElement) == 0 {
			errors = append(errors, fmt.Sprintf("cmd[%d] is empty", i+1))
			continue
		}

		trimmedCmd = append(trimmedCmd, trimmedElement)
	}

	if len(errors) > 0 {
		var zero DiffConfig
		return zero, errors
	}

	return DiffConfig{
		BaseLabel: baseLabel,
		HeadLabel: headLabel,
		Cmd:       trimmedCmd,
	}, nil
}
