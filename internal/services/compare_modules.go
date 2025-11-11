package services

import (
	"regexp"
	"slices"
	"sort"

	"github.com/dhth/tflens/internal/domain"
	"github.com/dhth/tflens/internal/hcl"
)

func GetComparisonResult(
	comparison domain.Comparison,
	globalValueRegex *regexp.Regexp,
	ignoreMissingModules bool,
) (domain.ComparisonResult, error) {
	var zero domain.ComparisonResult
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
			return zero, err
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

	return buildComparisonResult(store, sourceLabels, ignoreMissingModules), nil
}

func buildComparisonResult(store map[string]map[string]string, sourceLabels []string, ignoreMissingModules bool) domain.ComparisonResult {
	modules := make([]string, 0, len(store))
	for k := range store {
		modules = append(modules, k)
	}
	sort.Strings(modules)

	moduleResults := make([]domain.ModuleResult, 0, len(modules))
	for _, moduleName := range modules {
		labelToAttr := store[moduleName]

		values := make(map[string]string)

		isMissing := false
		for _, label := range sourceLabels {
			value, exists := labelToAttr[label]
			if exists {
				values[label] = value
			} else {
				isMissing = true
			}
		}

		status := determineModuleStatus(values, isMissing, ignoreMissingModules)

		moduleResults = append(moduleResults, domain.ModuleResult{
			Name:   moduleName,
			Values: values,
			Status: status,
		})
	}

	return domain.ComparisonResult{
		SourceLabels: sourceLabels,
		Modules:      moduleResults,
	}
}

func determineModuleStatus(values map[string]string, isMissing, ignoreMissingModules bool) domain.ModuleStatus {
	if isMissing && !ignoreMissingModules {
		return domain.StatusOutOfSync
	}

	var nonEmptyVersions []string
	for _, v := range values {
		if v != "" {
			nonEmptyVersions = append(nonEmptyVersions, v)
		}
	}

	if len(nonEmptyVersions) <= 1 {
		return domain.StatusUnknown
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
		return domain.StatusInSync
	}

	return domain.StatusOutOfSync
}
