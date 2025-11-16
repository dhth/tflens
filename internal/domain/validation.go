package domain

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	yaml "github.com/goccy/go-yaml"
)

var (
	ErrConfigHasErrors    = errors.New("config has errors")
	ErrCouldntParseConfig = errors.New("couldn't parse config")
)

type comparisonValidationErrors struct {
	index  int
	errors []string
}

func GetConfig(configBytes []byte) (Config, error) {
	var raw rawConfig

	err := yaml.Unmarshal(configBytes, &raw)
	if err != nil {
		return Config{}, fmt.Errorf("%w: %w", ErrCouldntParseConfig, err)
	}

	return parseRawConfig(raw)
}

func parseRawConfig(raw rawConfig) (Config, error) {
	var errors []comparisonValidationErrors
	var globalErrors []string

	// This check will need to be changed when more commands are added
	if len(raw.CompareModules.Comparisons) == 0 {
		globalErrors = append(globalErrors, "config has no comparisons configured")
	}

	var globalPattern *regexp.Regexp
	var err error

	if raw.CompareModules.ValueRegex != "" {
		globalPattern, err = regexp.Compile(raw.CompareModules.ValueRegex)
		if err != nil {
			globalErrors = append(globalErrors, fmt.Sprintf("invalid global valueRegex: %s", err.Error()))
		}
	}

	var validatedConfig Config
	validatedConfig.CompareModules.ValueRegex = globalPattern

	for c, comparison := range raw.CompareModules.Comparisons {
		var comparisonErrors []string

		comparisonName := strings.TrimSpace(comparison.Name)
		if len(comparisonName) == 0 {
			comparisonErrors = append(comparisonErrors, "comparison has an empty name")
		}

		attributeKey := strings.TrimSpace(comparison.AttributeKey)
		if len(attributeKey) == 0 {
			comparisonErrors = append(comparisonErrors, "comparison has an empty attribute key")
		}

		if len(comparison.Sources) <= 1 {
			comparisonErrors = append(comparisonErrors, "comparison needs to have at least 2 sources")
		}

		var comparisonPattern *regexp.Regexp
		if comparison.ValueRegex != "" {
			comparisonPattern, err = regexp.Compile(comparison.ValueRegex)
			if err != nil {
				comparisonErrors = append(comparisonErrors, fmt.Sprintf("invalid valueRegex: %s", err.Error()))
			}
		}

		for s, source := range comparison.Sources {
			if len(strings.TrimSpace(source.Label)) == 0 {
				comparisonErrors = append(comparisonErrors, fmt.Sprintf("source #%d has an empty label", s+1))
			}

			if len(strings.TrimSpace(source.Path)) == 0 {
				comparisonErrors = append(comparisonErrors, fmt.Sprintf("source #%d is empty", s+1))
				continue
			}

			if !strings.HasSuffix(source.Path, ".tf") {
				comparisonErrors = append(comparisonErrors, fmt.Sprintf("source #%d should have the extension .tf", s+1))
				continue
			}

			_, err := os.Stat(source.Path)
			if os.IsNotExist(err) {
				comparisonErrors = append(comparisonErrors, fmt.Sprintf("source #%d does not exist: %s", s+1, source.Path))
			} else if err != nil {
				comparisonErrors = append(comparisonErrors, fmt.Sprintf("couldn't check if source #%d exists: %s", s+1, err.Error()))
			}
		}

		if len(comparisonErrors) > 0 {
			errors = append(errors, comparisonValidationErrors{index: c, errors: comparisonErrors})
		} else {
			validatedComparison := Comparison{
				Name:          comparisonName,
				AttributeKey:  attributeKey,
				IgnoreModules: comparison.IgnoreModules,
				ValueRegex:    comparisonPattern,
			}

			for _, source := range comparison.Sources {
				validatedSource := Source{
					Path:  strings.TrimSpace(source.Path),
					Label: strings.TrimSpace(source.Label),
				}
				validatedComparison.Sources = append(validatedComparison.Sources, validatedSource)
			}

			validatedConfig.CompareModules.Comparisons = append(validatedConfig.CompareModules.Comparisons, validatedComparison)
		}
	}

	if len(globalErrors) > 0 || len(errors) > 0 {
		var errorLines []string

		if len(globalErrors) > 0 {
			for _, err := range globalErrors {
				errorLines = append(errorLines, fmt.Sprintf("- %s", err))
			}
		}

		for _, cErr := range errors {
			errorLines = append(errorLines, fmt.Sprintf("- comparison #%d has errors:", cErr.index+1))
			for _, err := range cErr.errors {
				errorLines = append(errorLines, fmt.Sprintf("  - %s", err))
			}
		}
		return validatedConfig, fmt.Errorf("%w:\n%s", ErrConfigHasErrors, strings.Join(errorLines, "\n"))
	}

	return validatedConfig, nil
}
