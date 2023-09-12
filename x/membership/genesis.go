package membership

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noria-net/module-membership/x/membership/keeper"
	"github.com/noria-net/module-membership/x/membership/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetDirectDemocracySettings(ctx, &genState.DirectDemocracy)

	// Enroll and add guardians
	for _, address := range genState.DirectDemocracy.Guardians {
		guardian := sdk.MustAccAddressFromBech32(address)
		if !k.IsMember(ctx, guardian) {
			// Add the member
			k.AppendMember(ctx, guardian)
			// Set their status to Electorate
			k.UpdateMemberStatus(ctx, guardian, types.MembershipStatus_MemberElectorate)
		}
		if !k.IsGuardian(ctx, guardian) {
			k.SetMemberGuardianStatus(ctx, guardian, true)
		}
	}

	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.DirectDemocracy = *k.GetDirectDemocracySettings(ctx)
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
