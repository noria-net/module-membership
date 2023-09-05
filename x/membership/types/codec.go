package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	gov_v1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgEnroll{}, "membership/Enroll", nil)
	cdc.RegisterConcrete(&MsgUpdateStatus{}, "membership/UpdateStatus", nil)
	cdc.RegisterConcrete(&MsgAddGuardians{}, "membership/AddGuardians", nil)
	cdc.RegisterConcrete(&AddGuardiansProposal{}, "membership/AddGuardiansProposal", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgEnroll{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateStatus{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddGuardians{},
	)
	registry.RegisterImplementations((*gov_v1beta1.Content)(nil),
		&AddGuardiansProposal{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
