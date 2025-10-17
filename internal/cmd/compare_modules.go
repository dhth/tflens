package cmd

import (
	"fmt"
	"os"

	"github.com/dhth/tflens/internal/domain"
	"github.com/dhth/tflens/internal/services"
	"github.com/spf13/cobra"
)

func newCompareModulesCmd(
	preRunE func(*cobra.Command, []string) error,
	config *domain.Config,
) *cobra.Command {
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
  valueRegex: "v?(\\d+\\.\\d+\\.\\d)"
  comparisons:
    - name: all-envs
      attributeKey: source
      sources:
        - path: environments/qa/apps/main.tf
          label: qa
        - path: environments/staging/apps/main.tf
          label: staging
        - path: environments/prod/apps/main.tf
          label: prod
---

$ tflens compare-modules all-envs

module      dev      staging    prod     in-sync
module_a    1.1.1    1.1.1      1.1.1    ✓
module_b    1.0.8    1.0.1      1.0.0    ✗
module_c    1.0.5    1.0.5      -        ✗
`,
		Args:              cobra.ExactArgs(1),
		SilenceUsage:      true,
		PersistentPreRunE: preRunE,
		RunE: func(_ *cobra.Command, args []string) error {
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

			return services.ShowModuleComparison(os.Stdout, *comparisonToUse, config.CompareModules.ValueRegex)
		},
	}

	return cmd
}
