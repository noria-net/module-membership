package cli

import (
	"fmt"
	"strconv"

	"strings"

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

func NewSubmitAddGuardiansProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-guardians [addresses]",
		Short: "Submit a proposal to add one or more guardians",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a proposal to add one or more guardians.
Separate multiple addresses with commas.

NOTE: Only existing members with status 'electorate' may be added as guardians.

NOTE: Only a guardian may submit this proposal.

Example: Adding a single guardian
$ %s tx membership add-guardians <address> --deposit=1000000unoria --from=<key_or_address>

Example: Adding multiple guardians
$ %s tx membership add-guardians <address1,address2,address3> --deposit=1000000unoria --from=<key_or_address>

`, version.AppName, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argGuardians := strings.Split(args[0], listSeparator)

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
			content := types.NewAddGuardiansProposal(
				title,
				description,
				from.String(),
				argGuardians,
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
