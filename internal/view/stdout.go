package view

import (
	"errors"
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/dhth/tflens/internal/domain"
)

var errCouldntPrintTable = errors.New("couldn't print table")

func RenderStdout(writer io.Writer, result domain.ComparisonResult, plain bool) error {
	rows := make([][]string, 0, len(result.Modules))

	rowStatuses := make(map[int]domain.ModuleStatus)

	for i, module := range result.Modules {
		row := make([]string, 0, len(result.SourceLabels)+2)
		row = append(row, module.Name)
		rowStatuses[i] = module.Status

		for _, label := range result.SourceLabels {
			value, ok := module.Values[label]
			if ok {
				row = append(row, value)
			} else {
				row = append(row, "-")
			}
		}

		row = append(row, module.Status.Symbol())
		rows = append(rows, row)
	}

	plainStyle := lipgloss.NewStyle().PaddingRight(4)
	outOfSyncStyle := plainStyle.Foreground(lipgloss.Color("1"))
	notApplicableStyle := plainStyle.Foreground(lipgloss.Color("8"))

	headers := make([]string, 0, len(result.SourceLabels)+2)
	headers = append(headers, "module")
	headers = append(headers, result.SourceLabels...)
	headers = append(headers, "in-sync")

	tbl := table.New().
		Border(lipgloss.HiddenBorder()).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if plain {
				return plainStyle
			}

			status, ok := rowStatuses[row]
			if !ok {
				return plainStyle
			}

			switch status {
			case domain.StatusOutOfSync:
				return outOfSyncStyle
			case domain.StatusNotApplicable:
				return notApplicableStyle
			default:
				return plainStyle
			}
		}).
		Headers(headers...).
		Rows(rows...)

	_, err := fmt.Fprintln(writer, tbl)
	if err != nil {
		return fmt.Errorf("%w: %w", errCouldntPrintTable, err)
	}

	return nil
}
