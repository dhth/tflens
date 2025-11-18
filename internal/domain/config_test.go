package domain

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestRawDiffConfigParse(t *testing.T) {
	sourceLabels := make(map[string]struct{})
	sourceLabels["base"] = struct{}{}
	sourceLabels["head"] = struct{}{}

	//-------------//
	//  SUCCESSES  //
	//-------------//

	t.Run("parsing correct config works", func(t *testing.T) {
		// GIVEN
		rawCfg := rawDiffConfig{
			BaseLabel: "base",
			HeadLabel: "head",
			Cmd:       []string{"git", "diff", "$TFLENS_DIFF_BASE_REF..$TFLENS_DIFF_HEAD_REF", "--", "modules/applications"},
		}

		// WHEN
		result, errors := rawCfg.parse(sourceLabels)

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
		result, errors := rawCfg.parse(sourceLabels)

		// THEN
		require.Empty(t, errors)
		snaps.MatchYAML(t, result)
	})

	//------------//
	//  FAILURES  //
	//------------//

	t.Run("parsing invalid config fails", func(t *testing.T) {
		// GIVEN
		rawCfg := rawDiffConfig{
			BaseLabel: "",
			HeadLabel: "",
			Cmd:       []string{},
		}

		// WHEN
		_, errors := rawCfg.parse(sourceLabels)

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
		_, errors := rawCfg.parse(sourceLabels)

		// THEN
		require.NotEmpty(t, errors)
		snaps.MatchYAML(t, errors)
	})

	t.Run("using incorrect labels fails", func(t *testing.T) {
		// GIVEN
		rawCfg := rawDiffConfig{
			BaseLabel: "absent\t",
			HeadLabel: " this-too\n",
			Cmd:       []string{"./scripts/generate-diff.sh", "apps"},
		}

		// WHEN
		_, errors := rawCfg.parse(sourceLabels)

		// THEN
		require.NotEmpty(t, errors)
		snaps.MatchYAML(t, errors)
	})
}
