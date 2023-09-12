package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgApproveMember = "approve_member"

var _ sdk.Msg = &MsgApproveMember{}

func NewMsgApproveMember(approver string, member string) *MsgApproveMember {
	return &MsgApproveMember{
		Approver: approver,
		Member:   member,
	}
}

func (msg *MsgApproveMember) Route() string {
	return RouterKey
}

func (msg *MsgApproveMember) Type() string {
	return TypeMsgApproveMember
}

func (msg *MsgApproveMember) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Approver)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgApproveMember) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApproveMember) ValidateBasic() error {
	// Approver and member addresses must be valid
	if _, err := sdk.AccAddressFromBech32(msg.Approver); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid approver address")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Member); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid member address")
	}

	return nil
}
