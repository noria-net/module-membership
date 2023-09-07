package cli

import (
	"fmt"
	"strconv"

	"strings"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	gov_v1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/noria-net/module-membership/x/membership/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func NewSubmitUpdateTotalVotingWeightProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-total-voting-weight [weight]",
		Short: "Submit a proposal to update the total voting weight of the guardians",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a proposal to update the total voting weight of the guardians.

NOTE: Only existing members with status 'electorate' may be added as guardians.

NOTE: Only a guardian may submit this proposal.

NOTE: The total voting weight must be > 0 and <= 1

Example: Updating the total voting weight
$ %s tx membership update-total-voting-weight <total_voting_weight> --deposit=1000000unoria --from=<key_or_address>

`, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argWeight, err := math.LegacyNewDecFromStr(args[0])
			if err != nil {
				return fmt.Errorf("invalid weight: %w", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			content := types.NewUpdateTotalVotingWeightProposal(
				title,
				description,
				from.String(),
				argWeight,
			)
			// Validate the proposal
			err = content.ValidateBasic()
			if err != nil {
				return err
			}

			msg, err := gov_v1beta1.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "The proposal title")
	cmd.Flags().String(cli.FlagDescription, "", "The proposal description")
	cmd.Flags().String(cli.FlagDeposit, "", "The proposal deposit")
	if err := cmd.MarkFlagRequired(cli.FlagTitle); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired(cli.FlagDescription); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired(cli.FlagDeposit); err != nil {
		panic(err)
	}

	return cmd
}
