package main

import (
	"github.com/spf13/cobra"

	// import _ to register escrow and options to apptx
	_ "github.com/tendermint/basecoin-examples/trader/cmd/trader/commands"
	"github.com/tendermint/basecoin/cmd/commands"
)

func main() {

	var RootCmd = &cobra.Command{
		Use: "trader",
	}

	RootCmd.AddCommand(
		commands.InitCmd,
		commands.StartCmd,
		commands.TxCmd,
		commands.QueryCmd,
		commands.KeyCmd,
		commands.VerifyCmd,
		commands.BlockCmd,
		commands.AccountCmd,
		commands.UnsafeResetAllCmd,
		commands.QuickVersionCmd("0.0.0"),
	)

	commands.ExecuteWithDebug(RootCmd)
}
