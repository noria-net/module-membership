package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noria-net/module-membership/x/membership/types"
)

func (k msgServer) AddGuardians(goCtx context.Context, msg *types.MsgAddGuardians) (*types.MsgAddGuardiansResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgAddGuardiansResponse{}, nil
}
