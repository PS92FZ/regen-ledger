package cli_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/regen-network/regen-ledger/v5/app"
	"github.com/regen-network/regen-ledger/v5/app/client/cli"
)

func TestInitCmd(t *testing.T) {
	nodeHome := os.TempDir()
	rootCmd, _ := cli.NewRootCmd()
	rootCmd.SetArgs([]string{
		"init",          // Test the init cmd
		"regenapp-test", // Moniker
		fmt.Sprintf("--%s=%s", flags.FlagHome, nodeHome), // Set home flag
	})

	err := cmd.Execute(rootCmd, app.EnvPrefix, app.DefaultNodeHome)
	require.NoError(t, err)
}
