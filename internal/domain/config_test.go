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
			BaseLabel: "base",
			HeadLabel: "head",
			Cmd:       []string{"git", "diff", "$TFLENS_DIFF_BASE_REF..$TFLENS_DIFF_HEAD_REF", "--", "modules/applications"},
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
			BaseLabel: " base",
			HeadLabel: "head ",
			Cmd:       []string{"git", " diff ", "$TFLENS_DIFF_BASE_REF..$TFLENS_DIFF_HEAD_REF", "--", "modules/applications"},
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
			BaseLabel: "",
			HeadLabel: "",
			Cmd:       []string{},
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
			BaseLabel: " ",
			HeadLabel: "\t",
			Cmd:       []string{" ", "\t", "\n", "\t\n"},
		}

		// WHEN
		_, errors := rawCfg.parse()

		// THEN
		require.NotEmpty(t, errors)
		snaps.MatchYAML(t, errors)
	})
}
