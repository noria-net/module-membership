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

const (
	FlagNickname = "nickname"
)

var _ = strconv.Itoa(0)

func CmdEnroll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enroll",
		Short: "Enroll the caller as an electorate member",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Enroll the caller as an electorate member.

NOTE: A Guardian will need to approve the member's enrollment before it becomes active. This is done with the approve-member command.

Example:
$ %s tx membership enroll --from=<key_or_address> --nickname=<nickname>
`, version.AppName)),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			nickname, err := cmd.Flags().GetString(FlagNickname)
			if err != nil {
				return err
			}

			msg := types.NewMsgEnroll(
				clientCtx.GetFromAddress().String(),
				nickname,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// Add nickname flag
	cmd.Flags().String(FlagNickname, "", "nickname of the member")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
