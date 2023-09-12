package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noria-net/module-membership/x/membership/types"
)

func (k msgServer) ApproveMember(goCtx context.Context, msg *types.MsgApproveMember) (*types.MsgApproveMemberResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgApproveMemberResponse{}, nil
}
