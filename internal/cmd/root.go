package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

const (
	configFileName = "tflens.yml"
)

var ErrComparisonNotFound = errors.New("comparison not found")

func Execute(version string) error {
	rootCmd, err := NewRootCommand(version)
	if err != nil {
		return err
	}

	return rootCmd.Execute()
}

func NewRootCommand(version string) (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:           "tflens",
		Short:         "tflens offers tiny utilities for terraform/opentofu/terragrunt codebases",
		SilenceErrors: true,
		Version:       version,
	}

	compareModulesCmd := newCompareModulesCmd()
	configCmd := newConfigCmd()

	rootCmd.AddCommand(compareModulesCmd)
	rootCmd.AddCommand(configCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd, nil
}
