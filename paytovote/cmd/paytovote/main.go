package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	_ "github.com/tendermint/basecoin-examples/paytovote/cmd/paytovote/commands"
	"github.com/tendermint/basecoin/cmd/commands"
)

func main() {

	var RootCmd = &cobra.Command{
		Use: "paytovote",
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
		commands.QuickVersionCmd("0.2.0"),
	)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
