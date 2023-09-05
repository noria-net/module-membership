package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddGuardians = "add_guardians"

var _ sdk.Msg = &MsgAddGuardians{}

func NewMsgAddGuardians(creator string, guardians []string) *MsgAddGuardians {
	return &MsgAddGuardians{
		Creator:   creator,
		Guardians: guardians,
	}
}

func (msg *MsgAddGuardians) Route() string {
	return RouterKey
}

func (msg *MsgAddGuardians) Type() string {
	return TypeMsgAddGuardians
}

func (msg *MsgAddGuardians) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddGuardians) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddGuardians) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate guardians
	if len(msg.Guardians) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "guardians cannot be empty")
	}

	// Ensure every address is valid
	for _, addr := range msg.Guardians {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid guardian address (%s)", err)
		}
	}

	return nil
}
