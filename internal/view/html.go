package view

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/dhth/tflens/internal/domain"
)

//go:embed assets/template.html
var builtInTemplate string

var (
	errCouldntParseCustomTemplate  = errors.New("couldn't parse custom template")
	ErrCouldntParseBuiltInTemplate = errors.New("couldn't parse built-in template")
	errCouldntPopulateTemplate     = errors.New("couldn't populate template")
)

func RenderHTML(result domain.ComparisonResult, config HTMLConfig, referenceTime time.Time) (string, error) {
	htmlData := NewHTMLData(config.Title, referenceTime)
	htmlData.Columns = append([]string{"module"}, result.SourceLabels...)
	htmlData.Columns = append(htmlData.Columns, "in-sync")

	for _, moduleResult := range result.Modules {
		row := HTMLRow{
			Data:   []string{moduleResult.Name},
			Status: moduleResult.Status.String(),
		}

		for _, label := range result.SourceLabels {
			row.Data = append(row.Data, moduleResult.Values[label])
		}

		row.Data = append(row.Data, moduleResult.Status.Symbol())

		htmlData.Rows = append(htmlData.Rows, row)

		if moduleResult.DiffResult == nil {
			continue
		}

		htmlData.Diffs = append(htmlData.Diffs, HTMLDiff{
			ModuleName: moduleResult.Name,
			Output:     string(moduleResult.DiffResult.Output),
			BaseLabel:  moduleResult.DiffResult.BaseLabel,
			HeadLabel:  moduleResult.DiffResult.HeadLabel,
			BaseRef:    moduleResult.DiffResult.BaseRef,
			HeadRef:    moduleResult.DiffResult.HeadRef,
		})
	}

	var tmpl *template.Template
	var templErr error

	var zero string

	if config.CustomTemplate != nil {
		tmpl, templErr = template.New("custom").Parse(*config.CustomTemplate)
		if templErr != nil {
			return zero, fmt.Errorf("%w: %w", errCouldntParseCustomTemplate, templErr)
		}
	} else {
		tmpl, templErr = template.New("built-in").Parse(builtInTemplate)
		if templErr != nil {
			return zero, fmt.Errorf("%w: %w", ErrCouldntParseBuiltInTemplate, templErr)
		}
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, htmlData)
	if err != nil {
		return zero, fmt.Errorf("%w: %w", errCouldntPopulateTemplate, err)
	}

	return buf.String(), nil
}
