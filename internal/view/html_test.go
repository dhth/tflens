package view

import (
	"testing"
	"time"

	"github.com/dhth/tflens/internal/domain"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderHTML(t *testing.T) {
	customTemplate := `{{.Title}}
Generated at {{.Timestamp}}

<table>
    <thead>
        <tr>
            {{- range .Columns }}
            <th>{{ . }}</th>
            {{- end }}
        </tr>
    </thead>
    <tbody>
        {{- range .Rows }}
        {{- if eq .Status "in_sync" }}
        <tr class="in-sync">
            {{- else if eq .Status "out_of_sync" }}
        <tr class="out-of-sync">
            {{- else }}
        <tr class="na">
            {{- end }}
            {{- range .Data }}
            <td>{{ . }}</td>
            {{- end }}
        </tr>
        {{- end }}
    </tbody>
</table>`
	referenceTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)

	//-------------//
	//  SUCCESSES  //
	//-------------//

	t.Run("works for built in template", func(t *testing.T) {
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
						"prod-eu": "2.1.0",
					},
					Status: domain.StatusInSync,
				},
				{
					Name: "module_c",
					Values: map[string]string{
						"dev":     "1.1.0",
						"prod-us": "1.1.0",
					},
					Status: domain.StatusOutOfSync,
				},
			},
		}

		config := HTMLConfig{
			Title: "Test Comparison",
		}

		// WHEN
		output, err := RenderHTML(result, config, referenceTime)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, output)
	})

	t.Run("works for all in sync modules", func(t *testing.T) {
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

		config := HTMLConfig{
			CustomTemplate: &customTemplate,
			Title:          "Test Comparison",
		}

		// WHEN
		output, err := RenderHTML(result, config, referenceTime)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, output)
	})

	t.Run("works when modules are out of sync", func(t *testing.T) {
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

		config := HTMLConfig{
			CustomTemplate: &customTemplate,
			Title:          "Test Comparison",
		}

		// WHEN
		output, err := RenderHTML(result, config, referenceTime)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, output)
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

		config := HTMLConfig{
			CustomTemplate: &customTemplate,
			Title:          "Test Comparison",
		}

		// WHEN
		output, err := RenderHTML(result, config, referenceTime)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, output)
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

		config := HTMLConfig{
			CustomTemplate: &customTemplate,
			Title:          "Test Comparison",
		}

		// WHEN
		output, err := RenderHTML(result, config, referenceTime)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, output)
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

		config := HTMLConfig{
			CustomTemplate: &customTemplate,
			Title:          "Test Comparison",
		}

		// WHEN
		output, err := RenderHTML(result, config, referenceTime)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, output)
	})

	//------------//
	//  FAILURES  //
	//------------//

	t.Run("fails when provided a malformed template", func(t *testing.T) {
		// GIVEN
		result := domain.ComparisonResult{
			SourceLabels: []string{"dev", "prod"},
			Modules: []domain.ModuleResult{
				{
					Name: "module_x",
					Values: map[string]string{
						"dev":  "1.0.0",
						"prod": "1.0.0",
					},
					Status: domain.StatusInSync,
				},
			},
		}

		malformedTemplate := `<html>{{if .Title}}<h1>{{.Title}}</h1>`

		config := HTMLConfig{
			CustomTemplate: &malformedTemplate,
			Title:          "Test",
		}

		// WHEN
		_, err := RenderHTML(result, config, referenceTime)

		// THEN
		assert.ErrorIs(t, err, errCouldntParseCustomTemplate)
	})
}
