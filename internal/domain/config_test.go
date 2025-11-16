package domain

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestRawDiffConfigParse(t *testing.T) {
	t.Run("parsing correct config works", func(t *testing.T) {
		// GIVEN
		rawCfg := rawDiffConfig{
			baseLabel: "base",
			headLabel: "head",
			cmd:       []string{"git", "diff", "$TFLENS_BASE_REF..$TFLENS_HEAD_REF", "--", "modules/applications"},
		}

		// WHEN
		result, errors := rawCfg.parse()

		// THEN
		require.Empty(t, errors)
		snaps.MatchYAML(t, result)
	})

	t.Run("parsing correct config with whitespace works", func(t *testing.T) {
		// GIVEN
		rawCfg := rawDiffConfig{
			baseLabel: " base",
			headLabel: "head ",
			cmd:       []string{"git", " diff ", "$TFLENS_BASE_REF..$TFLENS_HEAD_REF", "--", "modules/applications"},
		}

		// WHEN
		result, errors := rawCfg.parse()

		// THEN
		require.Empty(t, errors)
		snaps.MatchYAML(t, result)
	})

	t.Run("parsing invalid config fails", func(t *testing.T) {
		// GIVEN
		rawCfg := rawDiffConfig{
			baseLabel: "",
			headLabel: "",
			cmd:       []string{},
		}

		// WHEN
		_, errors := rawCfg.parse()

		// THEN
		require.NotEmpty(t, errors)
		snaps.MatchYAML(t, errors)
	})

	t.Run("parsing invalid config with whitespace only fails", func(t *testing.T) {
		// GIVEN
		rawCfg := rawDiffConfig{
			baseLabel: " ",
			headLabel: "\t",
			cmd:       []string{" ", "\t", "\n", "\t\n"},
		}

		// WHEN
		_, errors := rawCfg.parse()

		// THEN
		require.NotEmpty(t, errors)
		snaps.MatchYAML(t, errors)
	})
}
