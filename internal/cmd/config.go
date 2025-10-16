package cmd

import (
	"errors"
	"fmt"

	"github.com/dhth/tflens/internal/domain"
	yaml "github.com/goccy/go-yaml"
)

var ErrCouldntParseConfig = errors.New("couldn't parse config")

func getConfig(configBytes []byte) (domain.Config, error) {
	var raw domain.RawConfig

	err := yaml.Unmarshal(configBytes, &raw)
	if err != nil {
		return domain.Config{}, fmt.Errorf("%w: %w", ErrCouldntParseConfig, err)
	}

	return domain.NewConfigFromRaw(raw)
}
