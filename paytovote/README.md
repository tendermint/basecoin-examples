# PayToVote Plugin

### Description
Paytovote is a basic application which demonstrates how to create an instance of the basecoin system which
utilizes a custom paytovote plugin. The premise of this plugin is to allow users to pay a fee to create or
vote for user-specified issues. Unique fees are applied when voting or creating a new issue. Fees may use
coin types (for example "voteTokens" or "newIssueTokens"). Currently, the fee to cast a vote is decided by
the user when the issue is being generated, and the fee to create a new issue is defined globally within the
plugin CLI commands (cmd/commands)

### Install
Run `make all` in this directory. This will update all dependencies, run the
test suite, and install the `paytovote` binary to your `$GOPATH/bin`.  

### General Usage
 - create issues with `paytovote tx paytovote create-issue`
   - mandatory flags
     - --from value         Path to a private key to sign the transaction (default: "key.json")
     - --amount value       Amount of coins to send in the transaction (default: 0)
     - --issue value        Name of the issue to generate or vote for (default: "default issue")
     - --voteFee value      The fees required to  vote on this new issue, uses the format <amt><coin>,<amt2><coin2>,...
   - optional flags
     - --node value         Tendermint RPC address (default: "tcp://localhost:46657")
     - --chain_id value     ID of the chain for replay protection (default: "test_chain_id")
     - --coin value         Specify a coin denomination (default: "blank")
     - --gas value          The amount of gas for the transaction (default: 0)
     - --fee value          The transaction fee (default: 0)
     - --sequence value     Sequence number for the account (default: 0)
 - vote for issues with `paytovote tx paytovote vote` and flags listed below
   - mandatory flags
     - --from value      Path to a private key to sign the transaction (default: "key.json")
     - --amount value    Amount of coins to send in the transaction (default: 0)
     - --issue value     name of the issue to generate or vote for (default: "default issue")
     - --voteFor         set to true when vote be cast is a vote-for the issue, false if vote-against
   - optional flags
     - --node value      Tendermint RPC address (default: "tcp://localhost:46657")
     - --chain_id value  ID of the chain for replay protection (default: "test_chain_id")
     - --coin value      Specify a coin denomination (default: "blank")
     - --gas value       The amount of gas for the transaction (default: 0)
     - --fee value       The transaction fee (default: 0)
     - --sequence value  Sequence number for the account (default: 0)
 - query the state of an issue using the command `paytovote query p2vIssue <yourissuename>`

### Example CLI Usage
First perform the usual tendermint reset routine:

```
tendermint init
tendermint unsafe_reset_all
```

For the default genesis file provided in data/genesis.json we have specified a starting account at the hex 
address D397BC62B435F3CF50570FBAB4340FE52C60858F to have 1000 coins of "issueToken", and "voteToken" coins. 
The address, public key, and private key for this account are also stored under data/key.json
Now we can start basecoin with the paytovote plugin and the default genesis:

```
cd $GOPATH/src/github.com/tendermint/basecoin-examples/paytovote/data
paytovote start --in-proc 
```

Note that by navigating to the location of the genesis.json file we will automatically load it with `paytovote start`.

In another terminal window (or tab: ctrl-shift-t), we can run the client tool:

```
cd $GOPATH/src/github.com/tendermint/basecoin-examples/paytovote/data
paytovote account 0xD397BC62B435F3CF50570FBAB4340FE52C60858F
```

The above transaction will check for an account with the given hex address and list any coins within that account. We should see 
the initialized amount of 1000 issueToken, and voteToken. The default cost of generating a new issue 1 issueToken and is currently
hard coded into paytovote, let's create an issue which can be voted on. Notice the flags that are used in this proceedure:
 - `--from key.json` the transaction is coming from the account described within the key.json file under our current directory
 - `--voteFee 1voteToken` set the future cost of voting for this issue to 1 voteToken
 - `--amount 1issueToken` the amount of coins we are sending in with this transaction, in this case 1 issueToken 
 - `--issue freeFoobar` name of the issue we will be generating with this transaction

```
paytovote tx paytovote create-issue --from key.json --voteFee 1voteToken --amount 1issueToken --issue freeFoobar
```

Now we can query for our issue as see that it has been created and that no votes have yet been cast:

```
paytovote query p2vIssue freeFoobar
```

Next let's make a few votes, first we will vote for the issue once, and then against the issue twice

```
paytovote tx paytovote vote --from key.json --amount 1voteToken --issue freeFoobar --voteFor
paytovote tx paytovote vote --from key.json --amount 1voteToken --issue freeFoobar
paytovote tx paytovote vote --from key.json --amount 1voteToken --issue freeFoobar
```

To view the votes that have been cast query the transaction once more

```
paytovote query p2vIssue freeFoobar
```

Lastly we can verify that we have in fact spent 1 issueToken, and 3 voteToken

```
paytovote account 0xD397BC62B435F3CF50570FBAB4340FE52C60858F
```

### Thoughts for future development
 - Creating multiple vote options at issue generation (such as options or candidate name)
 - Alternative voting methods (for example ranked voting system)
   - Determine the type of voting mechanism when creating the issue
 - Allow votes to 'write in' their own candidate or spoil their ballot
 - Methods for the distributions of vote tokens

