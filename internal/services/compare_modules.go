package services

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"sort"

	"github.com/dhth/tflens/internal/domain"
	"github.com/dhth/tflens/internal/hcl"
)

var ErrCouldntComputeDiff = errors.New("couldn't compute diff")

func GetComparisonResult(
	comparison domain.Comparison,
	globalValueRegex *regexp.Regexp,
	ignoreMissingModules, includeDiffs bool,
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

	var diffCfg *domain.DiffConfig
	if includeDiffs {
		diffCfg = comparison.DiffCfg
	}
	result, err := buildComparisonResult(store, sourceLabels, ignoreMissingModules, diffCfg)
	if err != nil {
		return zero, err
	}

	return result, nil
}

func buildComparisonResult(
	store map[string]map[string]string,
	sourceLabels []string,
	ignoreMissingModules bool,
	diffCfg *domain.DiffConfig,
) (domain.ComparisonResult, error) {
	modules := make([]string, 0, len(store))
	for k := range store {
		modules = append(modules, k)
	}
	sort.Strings(modules)

	moduleResults := make([]domain.ModuleResult, 0, len(modules))
	for _, moduleName := range modules {
		labelToAttr := store[moduleName]

		//                 label  attribute
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

		var diffResult *domain.DiffResult
		if status == domain.StatusOutOfSync && diffCfg != nil {
			baseRef, baseExists := values[diffCfg.BaseLabel]
			headRef, headExists := values[diffCfg.HeadLabel]

			if baseExists && headExists && (baseRef != headRef) {
				diffOutput, diffErr := generateDiff(
					moduleName,
					baseRef,
					headRef,
					diffCfg.Cmd,
				)
				if diffErr != nil {
					return domain.ComparisonResult{}, fmt.Errorf("%w for module %q (command: %v): %w", ErrCouldntComputeDiff, moduleName, diffCfg.Cmd, diffErr)
				}

				diffResult = &domain.DiffResult{
					Output:    diffOutput,
					BaseLabel: diffCfg.BaseLabel,
					HeadLabel: diffCfg.HeadLabel,
					BaseRef:   baseRef,
					HeadRef:   headRef,
				}
			}
		}

		moduleResults = append(moduleResults, domain.ModuleResult{
			Name:       moduleName,
			Values:     values,
			Status:     status,
			DiffResult: diffResult,
		})
	}

	return domain.ComparisonResult{
		SourceLabels: sourceLabels,
		Modules:      moduleResults,
	}, nil
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
		return domain.StatusNotApplicable
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

func generateDiff(moduleName, baseLabel, headLabel string, command []string) ([]byte, error) {
	var zero []byte
	if len(command) == 0 {
		return zero, fmt.Errorf("empty command")
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("TFLENS_DIFF_BASE_REF=%s", baseLabel),
		fmt.Sprintf("TFLENS_DIFF_HEAD_REF=%s", headLabel),
		fmt.Sprintf("TFLENS_DIFF_MODULE_NAME=%s", moduleName),
	)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode := exitError.ExitCode()
			return zero, fmt.Errorf(`command exited with non success exit code

exit_code: %d
----- stdout -----
%s
----- stderr -----
%s`, exitCode, stdoutBuf.String(), stderrBuf.String())
		}

		return zero, fmt.Errorf("couldn't run command: %w", err)
	}

	return stdoutBuf.Bytes(), nil
}
