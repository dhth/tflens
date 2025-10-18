package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/dhth/tflens/internal/cmd"
	"github.com/dhth/tflens/internal/services"
)

var version = "dev"

func main() {
	err := cmd.Execute(version)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrModulesNotInSync):
		case errors.Is(err, cmd.ErrConfigValidationFoundErrors):
		default:
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())

		}
		os.Exit(1)
	}
}
