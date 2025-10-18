package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed assets/sample-config.yml
var sampleConfig string

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage tflens' configuration",
	}

	cmd.AddCommand(newConfigSampleCmd())

	return cmd
}

func newConfigSampleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sample",
		Short: "Print a sample configuration to stdout",
		Long: `Print a sample configuration file for tflens.

This can be used as a starting point for creating your own configuration:

$ tflens config sample > tflens.yml
`,
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Print(sampleConfig)
			return nil
		},
	}

	return cmd
}
