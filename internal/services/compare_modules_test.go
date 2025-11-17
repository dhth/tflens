package services

import (
	"regexp"
	"testing"

	"github.com/dhth/tflens/internal/domain"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestGetComparisonResult(t *testing.T) {
	t.Run("works for various cases", func(t *testing.T) {
		// GIVEN
		valueRegex := regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)

		comparison := domain.Comparison{
			Name:         "test-comparison",
			AttributeKey: "source",
			Sources: []domain.Source{
				{
					Path:  "testdata/environments/qa/main.tf",
					Label: "qa",
				},
				{
					Path:  "testdata/environments/staging/main.tf",
					Label: "staging",
				},
				{
					Path:  "testdata/environments/prod/main.tf",
					Label: "prod",
				},
			},
		}

		// WHEN
		result, err := GetComparisonResult(comparison, valueRegex, false)

		// THEN
		require.NoError(t, err)
		snaps.MatchYAML(t, result)
	})

	t.Run("works when missing modules are to be ignored", func(t *testing.T) {
		// GIVEN
		valueRegex := regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)

		comparison := domain.Comparison{
			Name:         "test-comparison",
			AttributeKey: "source",
			Sources: []domain.Source{
				{
					Path:  "testdata/environments/qa/main.tf",
					Label: "qa",
				},
				{
					Path:  "testdata/environments/staging/main.tf",
					Label: "staging",
				},
				{
					Path:  "testdata/environments/prod/main.tf",
					Label: "prod",
				},
			},
		}

		// WHEN
		result, err := GetComparisonResult(comparison, valueRegex, true)

		// THEN
		require.NoError(t, err)
		snaps.MatchYAML(t, result)
	})

	t.Run("ignoring modules works", func(t *testing.T) {
		// GIVEN
		valueRegex := regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)

		comparison := domain.Comparison{
			Name:         "test-comparison-sync",
			AttributeKey: "source",
			Sources: []domain.Source{
				{
					Path:  "testdata/environments/staging/main.tf",
					Label: "staging",
				},
				{
					Path:  "testdata/environments/prod/main.tf",
					Label: "prod",
				},
			},
			IgnoreModules: []string{"module_a", "module_b"},
		}

		// WHEN
		result, err := GetComparisonResult(comparison, valueRegex, false)

		// THEN
		require.NoError(t, err)
		snaps.MatchYAML(t, result)
	})
}

func TestBuildComparisonResult(t *testing.T) {
	store := map[string]map[string]string{
		"module_a": {
			"qa":      "1.2.3",
			"staging": "1.2.3",
			"prod":    "1.2.3",
		},
		"module_b": {
			"qa":      "1.0.0",
			"staging": "1.1.0",
			"prod":    "1.2.0",
		},
		"module_c": {
			"qa":      "2.0.0",
			"staging": "2.0.0",
			"prod":    "",
		},
		"module_d": {
			"qa":      "3.0.0",
			"staging": "",
			"prod":    "",
		},
	}

	t.Run("works when missing modules are not ignored", func(t *testing.T) {
		// GIVEN
		sourceLabels := []string{"qa", "staging", "prod"}

		// WHEN
		result, err := buildComparisonResult(store, sourceLabels, false, nil)

		// THEN
		require.NoError(t, err)
		snaps.MatchYAML(t, result)
	})

	t.Run("works when missing modules are ignored", func(t *testing.T) {
		// GIVEN
		sourceLabels := []string{"qa", "staging", "prod"}

		// WHEN
		result, err := buildComparisonResult(store, sourceLabels, true, nil)

		// THEN
		require.NoError(t, err)
		snaps.MatchYAML(t, result)
	})
}
