package commands

import (
	"errors"
	"fmt"

	"github.com/tendermint/basecoin-examples/paytovote"
	bcmd "github.com/tendermint/basecoin/cmd/commands"
	"github.com/tendermint/basecoin/types"
	cmn "github.com/tendermint/go-common"
	"github.com/tendermint/go-wire"

	"github.com/urfave/cli"
)

const PaytovoteName = "paytovote"

var (
	//common flag
	IssueFlag = cli.StringFlag{
		Name:  "issue",
		Value: "default issue",
		Usage: "name of the issue to generate or vote for",
	}

	//createIssue flags
	VoteFeeFlag = cli.StringFlag{
		Name:  "voteFee",
		Value: "1voteToken",
		Usage: "the fees required to  vote on this new issue, uses the format <amt><coin>,<amt2><coin2>,... (eg: 1gold,2silver,5btc)",
	}

	//vote flag
	VoteForFlag = cli.BoolFlag{
		Name:  "voteFor",
		Usage: "if present vote will be a vote-for, if absent a vote-against",
	}
)

var (
	P2VTxCmd = cli.Command{
		Name:  "paytovote",
		Usage: "Send transactions to the paytovote plugin",
		Subcommands: []cli.Command{
			P2VCreateIssueCmd,
			P2VVoteCmd,
		},
	}

	P2VQueryIssueCmd = cli.Command{
		Name:  "p2vIssue",
		Usage: "Query a paytovote issue",
		Action: func(c *cli.Context) error {
			return cmdQueryIssue(c)
		},
	}

	P2VCreateIssueCmd = cli.Command{
		Name:  "create-issue",
		Usage: "Create an issue which can be voted for",
		Action: func(c *cli.Context) error {
			return cmdCreateIssue(c)
		},
		Flags: append(bcmd.TxFlags,
			IssueFlag,
			VoteFeeFlag,
		),
	}

	P2VVoteCmd = cli.Command{
		Name:  "vote",
		Usage: "Vote for an existing issue",
		Action: func(c *cli.Context) error {
			return cmdVote(c)
		},
		Flags: append(bcmd.TxFlags,
			IssueFlag,
			VoteForFlag,
		),
	}
)

func init() {
	bcmd.RegisterTxSubcommand(P2VTxCmd)
	bcmd.RegisterQuerySubcommand(P2VQueryIssueCmd)
	bcmd.RegisterStartPlugin(PaytovoteName,
		func() types.Plugin { return paytovote.New() })
}

func cmdCreateIssue(c *cli.Context) error {
	issue := c.String(IssueFlag.Name)
	feeStr := c.String(VoteFeeFlag.Name)

	voteFee, err := bcmd.ParseCoins(feeStr)
	if err != nil {
		return err
	}

	createIssueFee := types.Coins{{"issueToken", 1}} //manually set the cost to create a new issue here

	txBytes := paytovote.NewCreateIssueTxBytes(issue, voteFee, createIssueFee)

	fmt.Println("Issue creation transaction sent")
	return bcmd.AppTx(c, PaytovoteName, txBytes)
}

func cmdVote(c *cli.Context) error {
	issue := c.String(IssueFlag.Name)
	voteFor := c.Bool(VoteForFlag.Name)

	var voteTB byte = paytovote.TypeByteVoteFor
	if !voteFor {
		voteTB = paytovote.TypeByteVoteAgainst
	}

	txBytes := paytovote.NewVoteTxBytes(issue, voteTB)

	fmt.Println("Vote transaction sent")
	return bcmd.AppTx(c, PaytovoteName, txBytes)
}

func cmdQueryIssue(c *cli.Context) error {

	//get the parent context
	parentContext := c.Parent()

	//get the issue, generate issue key
	if len(c.Args()) != 1 {
		return errors.New("query command requires an argument ([issue])")
	}
	issue := c.Args()[0]
	issueKey := paytovote.IssueKey(issue)

	//perform the query, get response
	resp, err := bcmd.Query(parentContext.String("node"), issueKey)
	if err != nil {
		return err
	}
	if !resp.Code.IsOK() {
		return errors.New(cmn.Fmt("Query for issueKey (%v) returned non-zero code (%v): %v",
			string(issueKey), resp.Code, resp.Log))
	}

	//get the paytovote issue object and print it
	p2vIssue, err := paytovote.GetIssueFromWire(resp.Value)
	if err != nil {
		return err
	}
	fmt.Println(string(wire.JSONBytes(p2vIssue)))

	return nil
}
