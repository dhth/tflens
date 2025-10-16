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
		if !errors.Is(err, services.ErrModulesNotInSync) {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		}
		os.Exit(1)
	}
}
