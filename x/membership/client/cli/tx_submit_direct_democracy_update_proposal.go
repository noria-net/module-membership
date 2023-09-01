package cli

import (
	"fmt"
	"strconv"

	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/noria-net/module-membership/x/membership/client/utils"
	"github.com/spf13/cobra"

	gov_v1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

var _ = strconv.Itoa(0)

const (
	FlagAddGuardians      = "add-guardians"
	FlagRemoveGuardians   = "remove-guardians"
	FlagTotalVotingWeight = "total-voting-weight"
)

func NewSubmitDirectDemocracyUpdateProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-direct-democracy",
		Short: "Submit a proposal to update the Direct Democracy settings",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a proposal to add or remove guardians from the electorate, 
or to update the total voting weight of the electorate. The proposal will be put to a governance vote
and only applied if it succeeds.

The format of the proposal json JSON is as follows:
{
	  "title": "<title>",
	  "description": "<description>",
	  "guardians_to_add": [
		"<guardian_address>"
	  ],
	  "guardians_to_remove": [
		"<guardian_address>"
	  ],
	  "total_voting_weight": "<total_voting_weight>"	
}

NOTE: There must be at least one guardian to add or remove, or there must be a total voting weight to update.
Otherwise the command will be rejected because there is no work to do.

NOTE: The total voting weight must be a decimal value between 0 and 1, inlcusive.

Example - Adding a guardian:
$ %s tx membership update-direct-democracy <path/to/proposal.json> --deposit %s --from=<key_or_address>

where proposal.json contains:
{
	  "title": "Add guardian",
	  "description": "Add guardian to the electorate",
	  guardians_to_add: [
		"cosmos1..."
	  ]
}

Example - Updating the total voting weight:
$ %s tx membership update-direct-democracy <path/to/proposal.json> --deposit %s --from=<key_or_address>

where proposal.json contains:
{
	  "title": "Update total voting weight",
	  "description": "Update the total voting weight of the electorate",
	  "total_voting_weight": "0.51
}

`,
				version.AppName,
				"1000000unoria",
				version.AppName,
				"1000000unoria",
			)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse the proposal
			proposal, err := utils.ParseDirectDemocracyUpdateProposal(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			// Validate the proposal
			err = proposal.ValidateBasic()
			if err != nil {
				return err
			}

			// Get initial deposit
			deposit, err := parseInitialDeposit(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			msg, err := gov_v1beta1.NewMsgSubmitProposal(&proposal, deposit, from)
			if err != nil {
				return err
			}

			// TODO: Not sure if this also validates the proposal? Can we remove the ValidateBasic call above?
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// Add flags
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
