package cli

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestCompareModulesCmd(t *testing.T) {
	fx, err := newFixture()
	require.NoErrorf(t, err, "error setting up fixture: %s", err)

	defer func() {
		err := fx.cleanup()
		require.NoErrorf(t, err, "error cleaning up fixture: %s", err)
	}()

	//-------------//
	//  SUCCESSES  //
	//-------------//

	t.Run("help flag works", func(t *testing.T) {
		// GIVEN
		args := []string{
			"compare-modules",
			"--help",
		}

		// WHEN
		result, err := fx.runCmd(args)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})

	t.Run("works for correct config", func(t *testing.T) {
		// GIVEN
		args := []string{
			"compare-modules",
			"--config-path", "testdata/config/good.yml",
			"--stdout-plain",
			"apps",
		}

		// WHEN
		result, err := fx.runCmd(args)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})

	t.Run("ignoring missing modules works", func(t *testing.T) {
		// GIVEN
		args := []string{
			"compare-modules",
			"--ignore-missing-modules",
			"--config-path", "testdata/config/good.yml",
			"--stdout-plain",
			"apps",
		}

		// WHEN
		result, err := fx.runCmd(args)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})

	//------------//
	//  FAILURES  //
	//------------//

	t.Run("fails for incorrect config", func(t *testing.T) {
		// GIVEN
		args := []string{
			"compare-modules",
			"--config-path", "testdata/config/bad.yml",
			"apps",
		}

		// WHEN
		result, err := fx.runCmd(args)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})

	t.Run("fails for invalid output format flag", func(t *testing.T) {
		// GIVEN
		args := []string{
			"compare-modules",
			"--config-path", "testdata/config/good.yml",
			"--output-format", "invalid",
			"apps",
		}

		// WHEN
		result, err := fx.runCmd(args)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})
}
