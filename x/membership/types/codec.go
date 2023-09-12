package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	gov_v1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgEnroll{}, "membership/Enroll", nil)
	cdc.RegisterConcrete(&MsgUpdateStatus{}, "membership/UpdateStatus", nil)
	cdc.RegisterConcrete(&AddGuardiansProposal{}, "membership/AddGuardiansProposal", nil)
	cdc.RegisterConcrete(&RemoveGuardiansProposal{}, "membership/RemoveGuardiansProposal", nil)
	cdc.RegisterConcrete(&UpdateTotalVotingWeightProposal{}, "membership/UpdateTotalVotingWeightProposal", nil)
	cdc.RegisterConcrete(&MsgApproveMember{}, "membership/ApproveMember", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgEnroll{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateStatus{},
	)
	registry.RegisterImplementations((*gov_v1beta1.Content)(nil),
		&AddGuardiansProposal{},
	)
	registry.RegisterImplementations((*gov_v1beta1.Content)(nil),
		&RemoveGuardiansProposal{},
	)
	registry.RegisterImplementations((*gov_v1beta1.Content)(nil),
		&UpdateTotalVotingWeightProposal{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgApproveMember{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
