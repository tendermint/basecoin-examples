package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/tendermint/basecoin-examples/stake"
	bcmd "github.com/tendermint/basecoin/cmd/commands"
	"github.com/tendermint/basecoin/types"
	wire "github.com/tendermint/go-wire"
)

var (
	//flags
	validatorFlag string
	amountFlag    int

	CmdBond = &cobra.Command{
		Use:   "bond",
		Short: "Bond some coins to give voting power to a validator",
		RunE:  cmdBond,
	}
)

func init() {

	flags := []bcmd.Flag2Register{
		{&validatorFlag, "validator", "", "Validator's public key"},
		{&amountFlag, "amount", 0, "Amount of coins"},
	}

	bcmd.RegisterFlags(CmdBond, flags)

	bcmd.RegisterTxSubcommand(CmdBond)
	bcmd.RegisterStartPlugin("stake",
		func() types.Plugin {
			return stake.New(stake.Params{
				UnbondingPeriod: 100,
				TokenDenom:      "atom",
			})
		},
	)
}

func cmdBond(cmd *cobra.Command, args []string) error {

	// convert validator pubkey to bytes
	validator, err := hex.DecodeString(bcmd.StripHex(validatorFlag))
	if err != nil {
		return errors.Errorf("Validator is invalid hex: %v\n", err)
	}

	bondTx := stake.BondTx{ValidatorPubKey: validator}
	fmt.Println("BondTx:", string(wire.JSONBytes(bondTx)))
	bytes := wire.BinaryBytes(bondTx)
	return bcmd.AppTx("stake", bytes)
}
