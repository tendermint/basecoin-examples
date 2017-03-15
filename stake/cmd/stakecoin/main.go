package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	// import _ to register the stake plugin to apptx
	_ "github.com/tendermint/basecoin-examples/stake/cmd/stakecoin/commands"
	"github.com/tendermint/basecoin/cmd/commands"
)

func main() {

	var RootCmd = &cobra.Command{
		Use: "stakecoin",
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

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
