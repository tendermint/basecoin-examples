package main

import (
	"github.com/spf13/cobra"

	// import _ to register the mint plugin to apptx
	_ "github.com/tendermint/basecoin-examples/mintcoin/cmd/mintcoin/commands"
	"github.com/tendermint/basecoin/cmd/commands"
)

func main() {

	var RootCmd = &cobra.Command{
		Use: "mintcoin",
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
		commands.QuickVersionCmd("0.1.0"),
	)

	commands.ExecuteWithDebug(RootCmd)
}
