package membership

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gov_v1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/noria-net/module-membership/x/membership/keeper"
	"github.com/noria-net/module-membership/x/membership/types"
)

func NewMembershipProposalHandler(k keeper.Keeper) gov_v1beta1.Handler {
	return func(ctx sdk.Context, content gov_v1beta1.Content) error {
		switch c := content.(type) {
		case *types.AddGuardiansProposal:
			return keeper.HandleAddGuardiansProposal(ctx, k, c)
		case *types.RemoveGuardiansProposal:
			return keeper.HandleRemoveGuardiansProposal(ctx, k, c)
		case *types.UpdateTotalVotingWeightProposal:
			return keeper.HandleUpdateTotalVotingWeightProposal(ctx, k, c)
		default:
			return errors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized membership proposal content type: %T", c)
		}
	}
}
