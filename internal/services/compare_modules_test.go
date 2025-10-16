package services

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/dhth/tflens/internal/domain"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShowModuleComparison(t *testing.T) {
	t.Run("modules with drift detected", func(t *testing.T) {
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

		var buf bytes.Buffer

		// WHEN
		err := ShowModuleComparison(&buf, comparison, valueRegex)

		// THEN
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrModulesNotInSync)

		output := buf.String()
		snaps.MatchStandaloneSnapshot(t, output)
	})

	t.Run("all modules in sync", func(t *testing.T) {
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
			IgnoreModules: []string{"module_b"},
		}

		var buf bytes.Buffer

		// WHEN
		err := ShowModuleComparison(&buf, comparison, valueRegex)

		// THEN
		require.NoError(t, err)

		output := buf.String()
		snaps.MatchStandaloneSnapshot(t, output)
	})
}

func TestGetStatus(t *testing.T) {
	tests := []struct {
		name       string
		values     []string
		hasMissing bool
		expected   syncStatus
	}{
		{
			name:     "all values match",
			values:   []string{"1.2.3", "1.2.3", "1.2.3"},
			expected: statusInSync,
		},
		{
			name:     "values don't match",
			values:   []string{"1.2.3", "1.2.4", "1.2.3"},
			expected: statusOutOfSync,
		},
		{
			name:       "module missing in some environments",
			values:     []string{"1.2.3", "1.2.3"},
			hasMissing: true,
			expected:   statusOutOfSync,
		},
		{
			name:     "all values match with empty string ignored",
			values:   []string{"1.2.3", "1.2.3", ""},
			expected: statusOutOfSync,
		},
		{
			name:     "values don't match with empty string",
			values:   []string{"1.2.3", "1.2.4", ""},
			expected: statusOutOfSync,
		},
		{
			name:     "single non-empty value with empty strings",
			values:   []string{"1.2.3", ""},
			expected: statusOutOfSync,
		},
		{
			name:     "all empty values",
			values:   []string{"", "", ""},
			expected: statusUnknown,
		},
		{
			name:     "empty values slice",
			values:   []string{},
			expected: statusUnknown,
		},
		{
			name:     "single matching value",
			values:   []string{"1.2.3"},
			expected: statusInSync,
		},
		{
			name:     "two different values",
			values:   []string{"1.2.3", "1.2.4"},
			expected: statusOutOfSync,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStatus(tt.values, tt.hasMissing)
			assert.Equal(t, tt.expected, got)
		})
	}
}
