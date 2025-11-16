package cmd

import (
	_ "embed"
	"errors"
	"fmt"
	"os"

	"github.com/dhth/tflens/internal/domain"
	"github.com/spf13/cobra"
)

//go:embed assets/sample-config.yml
var sampleConfig string

var ErrConfigValidationFoundErrors = errors.New("config validation found errors")

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage tflens' configuration",
	}

	cmd.AddCommand(newConfigSampleCmd())
	cmd.AddCommand(newConfigValidateCmd())

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

func newConfigValidateCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:          "validate",
		Short:        "Validate tflens' configuration file",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			configBytes, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrCouldntReadConfigFile, err)
			}

			_, err = domain.GetConfig(configBytes)
			if errors.Is(err, domain.ErrCouldntParseConfig) {
				return err
			} else if err != nil {
				fmt.Println(err.Error())
				return ErrConfigValidationFoundErrors
			}

			fmt.Println("Configuration is valid")
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

	return cmd
}
