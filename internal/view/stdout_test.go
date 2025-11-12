package view

import (
	"bytes"
	"testing"

	"github.com/dhth/tflens/internal/domain"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestRenderStdout(t *testing.T) {
	t.Run("works for all in-sync modules", func(t *testing.T) {
		// GIVEN
		result := domain.ComparisonResult{
			SourceLabels: []string{"dev", "prod-us", "prod-eu"},
			Modules: []domain.ModuleResult{
				{
					Name: "module_a",
					Values: map[string]string{
						"dev":     "1.0.0",
						"prod-us": "1.0.0",
						"prod-eu": "1.0.0",
					},
					Status: domain.StatusInSync,
				},
				{
					Name: "module_b",
					Values: map[string]string{
						"dev":     "2.0.0",
						"prod-us": "2.0.0",
						"prod-eu": "2.0.0",
					},
					Status: domain.StatusInSync,
				},
			},
		}

		var buf bytes.Buffer

		// WHEN
		err := RenderStdout(&buf, result, true)

		// THEN
		require.NoError(t, err)

		output := buf.String()
		snaps.MatchSnapshot(t, output)
	})

	t.Run("works when modules are out-of-sync", func(t *testing.T) {
		// GIVEN
		result := domain.ComparisonResult{
			SourceLabels: []string{"dev", "prod-us", "prod-eu"},
			Modules: []domain.ModuleResult{
				{
					Name: "module_a",
					Values: map[string]string{
						"dev":     "1.0.0",
						"prod-us": "1.0.1",
						"prod-eu": "2.0.0",
					},
					Status: domain.StatusOutOfSync,
				},
				{
					Name: "module_b",
					Values: map[string]string{
						"dev":     "2.0.0",
						"prod-us": "2.1.0",
						"prod-eu": "2.0.0",
					},
					Status: domain.StatusOutOfSync,
				},
			},
		}

		var buf bytes.Buffer

		// WHEN
		err := RenderStdout(&buf, result, true)

		// THEN
		require.NoError(t, err)

		output := buf.String()
		snaps.MatchSnapshot(t, output)
	})

	t.Run("works when missing modules are ignored", func(t *testing.T) {
		// GIVEN
		result := domain.ComparisonResult{
			SourceLabels: []string{"dev", "prod-us", "prod-eu"},
			Modules: []domain.ModuleResult{
				{
					Name: "module_a",
					Values: map[string]string{
						"dev":     "1.0.0",
						"prod-eu": "1.0.0",
					},
					Status: domain.StatusInSync,
				},
				{
					Name: "module_b",
					Values: map[string]string{
						"dev":     "2.0.0",
						"prod-us": "2.0.0",
					},
					Status: domain.StatusInSync,
				},
			},
		}

		var buf bytes.Buffer

		// WHEN
		err := RenderStdout(&buf, result, true)

		// THEN
		require.NoError(t, err)

		output := buf.String()
		snaps.MatchSnapshot(t, output)
	})

	t.Run("works when missing modules are not ignored", func(t *testing.T) {
		// GIVEN
		result := domain.ComparisonResult{
			SourceLabels: []string{"dev", "prod-us", "prod-eu"},
			Modules: []domain.ModuleResult{
				{
					Name: "module_a",
					Values: map[string]string{
						"dev":     "1.0.0",
						"prod-eu": "1.0.0",
					},
					Status: domain.StatusOutOfSync,
				},
				{
					Name: "module_b",
					Values: map[string]string{
						"dev":     "2.0.0",
						"prod-us": "2.0.0",
					},
					Status: domain.StatusOutOfSync,
				},
			},
		}

		var buf bytes.Buffer

		// WHEN
		err := RenderStdout(&buf, result, true)

		// THEN
		require.NoError(t, err)

		output := buf.String()
		snaps.MatchSnapshot(t, output)
	})

	t.Run("works when only one label is present", func(t *testing.T) {
		// GIVEN
		result := domain.ComparisonResult{
			SourceLabels: []string{"dev", "prod-us", "prod-eu"},
			Modules: []domain.ModuleResult{
				{
					Name: "module_a",
					Values: map[string]string{
						"dev": "1.0.0",
					},
					Status: domain.StatusNotApplicable,
				},
				{
					Name: "module_b",
					Values: map[string]string{
						"prod-us": "2.0.0",
					},
					Status: domain.StatusNotApplicable,
				},
				{
					Name: "module_c",
					Values: map[string]string{
						"dev":     "2.0.0",
						"prod-us": "2.0.0",
						"prod-eu": "2.0.0",
					},
					Status: domain.StatusInSync,
				},
			},
		}

		var buf bytes.Buffer

		// WHEN
		err := RenderStdout(&buf, result, true)

		// THEN
		require.NoError(t, err)

		output := buf.String()
		snaps.MatchSnapshot(t, output)
	})
}
