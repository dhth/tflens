package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dhth/tflens/internal/domain"
	"github.com/spf13/cobra"
)

const (
	configFileName = "tflens.yml"
)

var (
	ErrCouldntReadConfigFile = errors.New("couldn't read config file")
	ErrComparisonNotFound    = errors.New("comparison not found")
)

func Execute(version string) error {
	rootCmd, err := NewRootCommand(version)
	if err != nil {
		return err
	}

	return rootCmd.Execute()
}

func NewRootCommand(version string) (*cobra.Command, error) {
	var config domain.Config

	preRunE := func(_ *cobra.Command, _ []string) error {
		configBytes, err := os.ReadFile(configFileName)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrCouldntReadConfigFile, err)
		}

		config, err = getConfig(configBytes)
		if err != nil {
			return err
		}

		return nil
	}

	rootCmd := &cobra.Command{
		Use:           "tflens",
		Short:         "tflens offers tiny utilities for terraform/opentofu/terragrunt codebases",
		SilenceErrors: true,
		Version:       version,
	}

	compareModulesCmd := newCompareModulesCmd(
		preRunE,
		&config,
	)

	configCmd := newConfigCmd()

	rootCmd.AddCommand(compareModulesCmd)
	rootCmd.AddCommand(configCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd, nil
}
