package cli

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestConfigValidateCmd(t *testing.T) {
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
			"config",
			"validate",
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
			"config",
			"validate",
			"--config-path", "testdata/config/good.yml",
		}

		// WHEN
		result, err := fx.runCmd(args)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})

	t.Run("finds issues in empty config", func(t *testing.T) {
		// GIVEN
		args := []string{
			"config",
			"validate",
			"--config-path", "testdata/config/empty.yml",
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

	t.Run("fails for invalid yaml config", func(t *testing.T) {
		// GIVEN
		args := []string{
			"config",
			"validate",
			"--config-path", "testdata/config/invalid-yaml.yml",
		}

		// WHEN
		result, err := fx.runCmd(args)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})
}
