package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/noria-net/module-membership/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgApproveMember_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgApproveMember
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgApproveMember{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgApproveMember{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
