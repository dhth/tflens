package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dhth/tflens/internal/domain"
	"github.com/dhth/tflens/internal/services"
	"github.com/dhth/tflens/internal/view"
	"github.com/spf13/cobra"
)

var (
	errInvalidOutputFormat = errors.New("invalid output format provided")
	ErrModulesNotInSync    = errors.New("modules not in sync")
)

func newCompareModulesCmd() *cobra.Command {
	var config domain.Config
	var configPath string
	var outputFmtStr string
	var ignoreMissingModules bool

	cmd := &cobra.Command{
		Use:   "compare-modules <COMPARISON>",
		Short: "Compare modules by an attribute across multiple Terraform sources",
		Long: `Compare modules by an attribute across multiple Terraform sources.

This reads module blocks from the specified sources and compares a given attribute
(like 'source' or 'version') across them. Can be useful for ensuring consistency
across environments.

Example tflens.yml:
---
compareModules:
  # list of configured comparisons
  comparisons:
    # will be used when specifying the comparison to be run
    - name: apps
      # the attribute to use for comparison
      attributeKey: source
      # where to look for terraform files
      sources:
        - path: environments/dev/virginia/apps/main.tf
          # this label will appear in the comparison output
          label: dev
        - path: environments/prod/virginia/apps/main.tf
          label: prod-us
        - path: environments/prod/frankfurt/apps/main.tf
          # regex to extract the desired string from the attribute value
          # only applies to this source, overrides the global valueRegex
          # optional
          valueRegex: "v?(\\d+\\.\\d+\\.\\d+)"
          label: prod-eu

  # regex to extract the desired string from the attribute value
  # applies to all comparisons
  # optional
  valueRegex: "v?(\\d+\\.\\d+\\.\\d+)"
---

$ tflens compare-modules apps

module      dev       prod-us    prod-eu    in-sync
module_a    1.0.24    1.0.24     1.0.24     ✓
module_b    0.2.0     0.2.0      -          ✗
module_c    1.1.1     1.1.1      1.1.0      ✗
`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,

		PreRunE: func(_ *cobra.Command, _ []string) error {
			configBytes, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrCouldntReadConfigFile, err)
			}
			config, err = getConfig(configBytes)
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			outputFmt, outputFmtOk := domain.ParseOutputFormat(outputFmtStr)
			if !outputFmtOk {
				return fmt.Errorf("%w: %q; allowed values: %v", errInvalidOutputFormat, outputFmtStr, domain.GetOutputFormatValues())
			}

			comparisonName := args[0]
			var comparisonToUse *domain.Comparison
			for i := range config.CompareModules.Comparisons {
				if config.CompareModules.Comparisons[i].Name == comparisonName {
					comparisonToUse = &config.CompareModules.Comparisons[i]
					break
				}
			}

			if comparisonToUse == nil {
				return fmt.Errorf("%w: %q", ErrComparisonNotFound, comparisonName)
			}

			result, err := services.GetComparisonResult(*comparisonToUse, config.CompareModules.ValueRegex, ignoreMissingModules)
			if err != nil {
				return err
			}

			switch outputFmt {
			case domain.StdoutOutput:
				err := view.RenderStdout(os.Stdout, result)
				if err != nil {
					return fmt.Errorf("failed to render stdout: %w", err)
				}

				for _, moduleRes := range result.Modules {
					if moduleRes.Status == domain.StatusOutOfSync {
						return ErrModulesNotInSync
					}
				}

			case domain.HtmlOutput:
				fmt.Println("todo")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(
		&configPath,
		"config-path",
		"c",
		configFileName,
		"path to tflens' configuration file",
	)

	cmd.Flags().BoolVarP(
		&ignoreMissingModules,
		"ignore-missing-modules",
		"i",
		false,
		"to not have the absence of a module lead to an out-of-sync status",
	)

	cmd.Flags().StringVarP(
		&outputFmtStr,
		"output-format",
		"o",
		"stdout",
		fmt.Sprintf("output format for results; allowed values: %v", domain.GetOutputFormatValues()),
	)

	return cmd
}
