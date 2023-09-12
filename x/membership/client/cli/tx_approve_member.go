package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/noria-net/module-membership/x/membership/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdApproveMember() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-member [address]",
		Short: "Approve a member's pending enrollment",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Approve a member's pending enrollment.

NOTE: Only Guardians may execute this command.

Example:
$ %s tx membership approve-member <address> --from=<key_or_address>
`, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			argMemberAddress := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgApproveMember(
				clientCtx.GetFromAddress().String(),
				argMemberAddress,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
