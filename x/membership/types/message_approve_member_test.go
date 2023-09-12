package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgApproveMember_ValidateBasic(t *testing.T) {
	valid_1 := "cosmos1l0znsvddllw9knha3yx2svnlxny676d8ns7uys"
	valid_2 := "cosmos1j8pp7zvcu9z8vd882m284j29fn2dszh05cqvf9"
	invalid := "invalid_address"

	tests := []struct {
		name string
		msg  MsgApproveMember
		err  error
	}{
		{
			name: "invalid approver address",
			msg: MsgApproveMember{
				Approver: invalid,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "invalid member address",
			msg: MsgApproveMember{
				Approver: valid_1,
				Member:   invalid,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid message",
			msg: MsgApproveMember{
				Approver: valid_1,
				Member:   valid_2,
			},
			err: nil,
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
