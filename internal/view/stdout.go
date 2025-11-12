package view

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/dhth/tflens/internal/domain"
)

func RenderStdout(w io.Writer, result domain.ComparisonResult) error {
	tw := tabwriter.NewWriter(w, 0, 4, 4, ' ', 0)

	fmt.Fprint(tw, "module")
	for _, label := range result.SourceLabels {
		fmt.Fprintf(tw, "\t%s", label)
	}
	fmt.Fprint(tw, "\tin-sync")
	fmt.Fprintln(tw)

	for _, module := range result.Modules {
		fmt.Fprint(tw, module.Name)

		for _, label := range result.SourceLabels {
			value, ok := module.Values[label]
			if ok {
				fmt.Fprintf(tw, "\t%s", value)
			} else {
				fmt.Fprint(tw, "\t-")
			}
		}

		fmt.Fprintf(tw, "\t%s", module.Status.Symbol())
		fmt.Fprintln(tw)
	}

	if err := tw.Flush(); err != nil {
		return fmt.Errorf("failed to write table: %w", err)
	}

	return nil
}
