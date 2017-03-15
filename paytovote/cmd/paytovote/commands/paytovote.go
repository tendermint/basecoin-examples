package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tendermint/basecoin-examples/paytovote"
	bcmd "github.com/tendermint/basecoin/cmd/commands"
	"github.com/tendermint/basecoin/types"
	cmn "github.com/tendermint/go-common"
	"github.com/tendermint/go-wire"
)

const PaytovoteName = "paytovote"

var (
	//flags
	issueFlag   string
	voteFeeFlag string
	voteForFlag bool

	//commands
	P2VTxCmd = &cobra.Command{
		Use:   "paytovote",
		Short: "Send transactions to the paytovote plugin",
	}

	P2VQueryIssueCmd = &cobra.Command{
		Use:   "p2vIssue",
		Short: "Query a paytovote issue",
		Run:   queryIssueCmd,
	}

	P2VCreateIssueCmd = &cobra.Command{
		Use:   "create-issue",
		Short: "Create an issue which can be voted for",
		Run:   createIssueCmd,
	}

	P2VVoteCmd = &cobra.Command{
		Use:   "vote",
		Short: "Vote for an existing issue",
		Run:   voteCmd,
	}
)

func init() {

	//register flags

	issueFlag2Reg := bcmd.Flag2Register{&issueFlag, "issue", "default issue", "name of the issue to generate or vote for"}

	createIssueFlags := []bcmd.Flag2Register{
		issueFlag2Reg,
		{&voteFeeFlag, "voteFee", "1voteToken",
			"the fees required to  vote on this new issue, uses the format <amt><coin>,<amt2><coin2>,... (eg: 1gold,2silver,5btc)"},
	}

	voteFlags := []bcmd.Flag2Register{
		issueFlag2Reg,
		{&voteForFlag, "voteFor", false, "if present vote will be a vote-for, if absent a vote-against"},
	}

	bcmd.RegisterFlags(P2VCreateIssueCmd, createIssueFlags)
	bcmd.RegisterFlags(P2VVoteCmd, voteFlags)

	//register commands
	P2VTxCmd.AddCommand(P2VCreateIssueCmd, P2VVoteCmd)

	bcmd.RegisterTxSubcommand(P2VTxCmd)
	bcmd.RegisterQuerySubcommand(P2VQueryIssueCmd)
	bcmd.RegisterStartPlugin(PaytovoteName, func() types.Plugin { return paytovote.New() })
}

func createIssueCmd(cmd *cobra.Command, args []string) {

	voteFee, err := bcmd.ParseCoins(voteFeeFlag)
	if err != nil {
		cmn.Exit(fmt.Sprintf("%+v\n", err))
	}

	createIssueFee := types.Coins{{"issueToken", 1}} //manually set the cost to create a new issue here

	txBytes := paytovote.NewCreateIssueTxBytes(issueFlag, voteFee, createIssueFee)

	fmt.Println("Issue creation transaction sent")
	bcmd.AppTx(PaytovoteName, txBytes)
}

func voteCmd(cmd *cobra.Command, args []string) {

	var voteTB byte = paytovote.TypeByteVoteFor
	if !voteForFlag {
		voteTB = paytovote.TypeByteVoteAgainst
	}

	txBytes := paytovote.NewVoteTxBytes(issueFlag, voteTB)

	fmt.Println("Vote transaction sent")
	bcmd.AppTx(PaytovoteName, txBytes)
}

func queryIssueCmd(cmd *cobra.Command, args []string) {

	//get the parent context
	parentContext := cmd.Parent()

	//get the issue, generate issue key
	if len(args) != 1 {
		cmn.Exit("query command requires an argument ([issue])")
	}
	issue := args[0]
	issueKey := paytovote.IssueKey(issue)

	//perform the query, get response
	resp, err := bcmd.Query(parentContext.Flag("node").Value.String(), issueKey)
	if err != nil {
		cmn.Exit(fmt.Sprintf("%+v\n", err))
	}
	if !resp.Code.IsOK() {
		cmn.Exit(fmt.Sprintf("Query for issueKey (%v) returned non-zero code (%v): %v",
			string(issueKey), resp.Code, resp.Log))
	}

	//get the paytovote issue object and print it
	p2vIssue, err := paytovote.GetIssueFromWire(resp.Value)
	if err != nil {
		cmn.Exit(fmt.Sprintf("%+v\n", err))
	}
	fmt.Println(string(wire.JSONBytes(p2vIssue)))
}
