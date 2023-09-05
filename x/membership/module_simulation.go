package membership

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/noria-net/module-membership/testutil/sample"
	membershipsimulation "github.com/noria-net/module-membership/x/membership/simulation"
	"github.com/noria-net/module-membership/x/membership/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = membershipsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgEnroll = "op_weight_msg_enroll"
	// TODO: Determine the simulation weight value
	defaultWeightMsgEnroll int = 100

	opWeightMsgUpdateStatus = "op_weight_msg_update_status"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateStatus int = 100

	opWeightMsgAddGuardians = "op_weight_msg_add_guardians"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddGuardians int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	membershipGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&membershipGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgEnroll int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgEnroll, &weightMsgEnroll, nil,
		func(_ *rand.Rand) {
			weightMsgEnroll = defaultWeightMsgEnroll
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgEnroll,
		membershipsimulation.SimulateMsgEnroll(am.accountKeeper, am.keeper),
	))

	var weightMsgUpdateStatus int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateStatus, &weightMsgUpdateStatus, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateStatus = defaultWeightMsgUpdateStatus
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateStatus,
		membershipsimulation.SimulateMsgUpdateStatus(am.accountKeeper, am.keeper),
	))

	var weightMsgAddGuardians int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAddGuardians, &weightMsgAddGuardians, nil,
		func(_ *rand.Rand) {
			weightMsgAddGuardians = defaultWeightMsgAddGuardians
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddGuardians,
		membershipsimulation.SimulateMsgAddGuardians(am.accountKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgEnroll,
			defaultWeightMsgEnroll,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				membershipsimulation.SimulateMsgEnroll(am.accountKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateStatus,
			defaultWeightMsgUpdateStatus,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				membershipsimulation.SimulateMsgUpdateStatus(am.accountKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddGuardians,
			defaultWeightMsgAddGuardians,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				membershipsimulation.SimulateMsgAddGuardians(am.accountKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
