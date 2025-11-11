package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/dhth/tflens/internal/cmd"
)

var version = "dev"

func main() {
	err := cmd.Execute(version)
	if err != nil {
		switch {
		case errors.Is(err, cmd.ErrModulesNotInSync):
		case errors.Is(err, cmd.ErrConfigValidationFoundErrors):
		case errors.Is(err, cmd.ErrCouldntParseConfig):
			fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		default:
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())

		}
		os.Exit(1)
	}
}
