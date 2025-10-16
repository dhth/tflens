package services

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"slices"
	"sort"
	"text/tabwriter"

	"github.com/dhth/tflens/internal/domain"
	"github.com/dhth/tflens/internal/hcl"
)

var (
	ErrCouldntWriteTable = errors.New("couldn't write table")
	ErrModulesNotInSync  = errors.New("modules not in sync")
)

type syncStatus uint8

const (
	statusInSync syncStatus = iota
	statusOutOfSync
	statusUnknown
)

func ShowModuleComparison(writer io.Writer, comparison domain.Comparison, globalValueRegex *regexp.Regexp) error {
	sourceLabels := make([]string, len(comparison.Sources))
	for i, source := range comparison.Sources {
		sourceLabels[i] = source.Label
	}

	valueRegex := globalValueRegex
	if comparison.ValueRegex != nil {
		valueRegex = comparison.ValueRegex
	}

	//                module     label  attribute
	store := make(map[string]map[string]string)

	for _, source := range comparison.Sources {
		result, err := hcl.ParseModules(source.Path, comparison.AttributeKey, valueRegex)
		if err != nil {
			return err
		}

		for _, mod := range result {
			if slices.Contains(comparison.IgnoreModules, mod.Name) {
				continue
			}

			labelAttributeMap, ok := store[mod.Name]
			if !ok {
				labelAttributeMap = make(map[string]string)
			}

			labelAttributeMap[source.Label] = mod.Attribute
			store[mod.Name] = labelAttributeMap
		}
	}

	inSync, err := writeTable(writer, store, sourceLabels)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCouldntWriteTable, err)
	}

	if !inSync {
		return ErrModulesNotInSync
	}

	return nil
}

func writeTable(w io.Writer, store map[string]map[string]string, sourceLabels []string) (bool, error) {
	tw := tabwriter.NewWriter(w, 0, 4, 4, ' ', 0)

	fmt.Fprint(tw, "module")
	for _, label := range sourceLabels {
		fmt.Fprintf(tw, "\t%s", label)
	}
	fmt.Fprint(tw, "\tin-sync")
	fmt.Fprintln(tw)

	modules := make([]string, 0, len(store))
	for module := range store {
		modules = append(modules, module)
	}
	sort.Strings(modules)

	inSync := true
	for _, module := range modules {
		fmt.Fprint(tw, module)
		labelValues := store[module]

		var values []string
		hasMissing := false

		for _, label := range sourceLabels {
			value, exists := labelValues[label]
			if exists {
				fmt.Fprintf(tw, "\t%s", value)
				values = append(values, value)
			} else {
				fmt.Fprint(tw, "\t-")
				hasMissing = true
			}
		}

		status := getStatus(values, hasMissing)
		var statusSymbol string
		switch status {
		case statusInSync:
			statusSymbol = "✓"
		case statusOutOfSync:
			statusSymbol = "✗"
			inSync = false
		case statusUnknown:
			statusSymbol = "-"
		}
		fmt.Fprintf(tw, "\t%s", statusSymbol)
		fmt.Fprintln(tw)
	}

	return inSync, tw.Flush()
}

func getStatus(values []string, hasMissing bool) syncStatus {
	if hasMissing {
		return statusOutOfSync
	}

	if len(values) == 0 {
		return statusUnknown
	}

	var nonEmptyVersions []string
	for _, v := range values {
		if v != "" {
			nonEmptyVersions = append(nonEmptyVersions, v)
		}
	}

	if len(nonEmptyVersions) == 0 {
		return statusUnknown
	}

	if len(nonEmptyVersions) != len(values) {
		return statusOutOfSync
	}

	allMatch := true
	firstVersion := nonEmptyVersions[0]
	for _, v := range nonEmptyVersions[1:] {
		if v != firstVersion {
			allMatch = false
			break
		}
	}

	if allMatch {
		return statusInSync
	}

	return statusOutOfSync
}
