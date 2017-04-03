package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tendermint/basecoin-examples/trader/plugins/escrow"
	"github.com/tendermint/basecoin-examples/trader/types"
	bcmd "github.com/tendermint/basecoin/cmd/commands"
	bc "github.com/tendermint/basecoin/types"
	cmn "github.com/tendermint/go-common"
	wire "github.com/tendermint/go-wire"
)

const EscrowName = "escrow"

var (
	//flags
	EscrowNodeFlag    string
	EscrowRecvFlag    string
	EscrowArbiterFlag string
	EscrowAddrFlag    string
	EscrowExpireFlag  uint64
	EscrowPayoutFlag  bool

	//commands
	CmdEscrowTx = &cobra.Command{
		Use:   "escrow",
		Short: "Create and resolve escrows",
	}

	CmdEscrowCreateTx = &cobra.Command{
		Use:   "create",
		Short: "Create a new escrow by sending money",
		Run:   cmdEscrowCreateTx,
	}

	CmdEscrowResolveTx = &cobra.Command{
		Use:   "pay",
		Short: "Resolve the escrow by paying out of returning the money",
		Run:   cmdEscrowResolveTx,
	}

	CmdEscrowExpireTx = &cobra.Command{
		Use:   "expire",
		Short: "Call to expire the escrow if no action in a given time",
		Run:   cmdEscrowExpireTx,
	}

	CmdEscrowQuery = &cobra.Command{
		Use:   "query [address]",
		Short: "Return the contents of the given escrow",
		Run:   cmdEscrowQuery,
	}
)

func init() {

	//register flags
	queryFlags := []bcmd.Flag2Register{
		{&EscrowNodeFlag, "node", "tcp://localhost:46657", "Tendermint RPC address"},
	}
	addrFlag := bcmd.Flag2Register{
		&EscrowAddrFlag, "escrow", "", "The address of this escrow"}
	expireFlags := []bcmd.Flag2Register{
		addrFlag,
	}
	resolveFlags := []bcmd.Flag2Register{
		addrFlag,
		{&EscrowPayoutFlag, "abort-payout", false, "Set this flag if to return the money to the sender"},
	}
	createTxFlags := []bcmd.Flag2Register{
		{&EscrowRecvFlag, "recv", "", "Who is the intended recipient of the escrow"},
		{&EscrowArbiterFlag, "arbiter", "", "Who is the arbiter of the escrow"},
		{&EscrowExpireFlag, "expire", uint64(0), "The block height when the escrow expires"},
	}
	bcmd.RegisterFlags(CmdEscrowQuery, queryFlags)
	bcmd.RegisterFlags(CmdEscrowExpireTx, expireFlags)
	bcmd.RegisterFlags(CmdEscrowResolveTx, resolveFlags)
	bcmd.RegisterFlags(CmdEscrowCreateTx, createTxFlags)

	//register subcommands of EscrowTxCmd
	CmdEscrowTx.AddCommand(
		CmdEscrowCreateTx,
		CmdEscrowResolveTx,
		CmdEscrowExpireTx,
		CmdEscrowQuery,
	)

	//register with main tx command
	bcmd.RegisterTxSubcommand(CmdEscrowTx)
	bcmd.RegisterStartPlugin(EscrowName,
		func() bc.Plugin { return escrow.New(EscrowName) })
}

func cmdEscrowCreateTx(cmd *cobra.Command, args []string) {
	// convert destination address to bytes
	recv, err := hex.DecodeString(bcmd.StripHex(EscrowRecvFlag))
	if err != nil {
		cmn.Exit(fmt.Sprintf("Recv address is invalid hex: %+v\n", err))
	}

	// convert destination address to bytes
	arb, err := hex.DecodeString(bcmd.StripHex(EscrowArbiterFlag))
	if err != nil {
		cmn.Exit(fmt.Sprintf("Arbiter address is invalid hex: %+v\n", err))
	}

	tx := types.CreateEscrowTx{
		Recipient:  recv,
		Arbiter:    arb,
		Expiration: EscrowExpireFlag,
	}
	data := types.EscrowTxBytes(tx)
	bcmd.AppTx(EscrowName, data)
}

func cmdEscrowResolveTx(cmd *cobra.Command, args []string) {

	// convert destination address to bytes
	addr, err := hex.DecodeString(bcmd.StripHex(EscrowAddrFlag))
	if err != nil {
		cmn.Exit(fmt.Sprintf("Recv address is invalid hex: %+v\n", err))
	}

	tx := types.ResolveEscrowTx{
		Escrow: addr,
		Payout: !EscrowPayoutFlag,
	}
	data := types.EscrowTxBytes(tx)
	bcmd.AppTx(EscrowName, data)
}

func cmdEscrowExpireTx(cmd *cobra.Command, args []string) {

	// convert destination address to bytes
	addr, err := hex.DecodeString(bcmd.StripHex(EscrowAddrFlag))
	if err != nil {
		cmn.Exit(fmt.Sprintf("Recv address is invalid hex: %+v\n", err))
	}

	tx := types.ExpireEscrowTx{
		Escrow: addr,
	}
	data := types.EscrowTxBytes(tx)
	bcmd.AppTx(EscrowName, data)
}

func cmdEscrowQuery(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmn.Exit("account command requires an argument ([address])")
	}
	addrHex := bcmd.StripHex(args[0])

	// convert destination address to bytes
	addr, err := hex.DecodeString(addrHex)
	if err != nil {
		cmn.Exit(fmt.Sprintf("Recv address is invalid hex: %+v\n", err))
	}

	esc, err := getEscrow(EscrowNodeFlag, addr)
	if err != nil {
		cmn.Exit(fmt.Sprintf("%+v\n", err))
	}

	fmt.Println(string(wire.JSONBytes(esc)))
}

func getEscrow(tmAddr string, address []byte) (*types.EscrowData, error) {
	prefix := []byte(fmt.Sprintf("%s/", EscrowName))
	key := append(prefix, address...)
	response, err := bcmd.Query(tmAddr, key)
	if err != nil {
		return nil, err
	}

	escrowBytes := response.Value

	if len(escrowBytes) == 0 {
		return nil, fmt.Errorf("Escrow bytes are empty for address: %X ", address)
	}
	esc, err := types.ParseEscrow(escrowBytes)
	if err != nil {
		return nil, fmt.Errorf("Error reading account %X error: %v",
			escrowBytes, err.Error())
	}
	return &esc, nil

}
