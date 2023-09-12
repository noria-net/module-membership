package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgApproveMember = "approve_member"

var _ sdk.Msg = &MsgApproveMember{}

func NewMsgApproveMember(creator string) *MsgApproveMember {
	return &MsgApproveMember{
		Creator: creator,
	}
}

func (msg *MsgApproveMember) Route() string {
	return RouterKey
}

func (msg *MsgApproveMember) Type() string {
	return TypeMsgApproveMember
}

func (msg *MsgApproveMember) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
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
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
