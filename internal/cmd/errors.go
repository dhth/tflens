package cmd

import (
	"errors"

	"github.com/dhth/tflens/internal/view"
)

func IsErrorUnexpected(err error) bool {
	return errors.Is(err, view.ErrCouldntParseBuiltInTemplate)
}
