package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/noria-net/module-membership/x/membership/types"
)

func (k msgServer) ApproveMember(goCtx context.Context, msg *types.MsgApproveMember) (*types.MsgApproveMemberResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	approverAddr := sdk.MustAccAddressFromBech32(msg.Approver)
	memberAddr := sdk.MustAccAddressFromBech32(msg.Member)

	// Only guardians can approve members
	if !k.Keeper.IsGuardian(ctx, approverAddr) {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "only guardians can approve members")
	}

	member, found := k.GetMemberAccount(ctx, memberAddr)
	// Member must exist
	if !found {
		return nil, errors.Wrap(types.ErrMemberNotFound, "member does not exist")
	}

	// Member status must be Pending
	if member.Status != types.MembershipStatus_MemberStatusPendingApproval {
		return nil, errors.Wrap(types.ErrMemberNotPendingApproval, "member is not pending approval")
	}

	// Update the member's status
	member.Status = types.MembershipStatus_MemberElectorate

	// Save the member back to the store
	k.UpdateMemberStatus(ctx, memberAddr, types.MembershipStatus_MemberElectorate)

	// Publish events
	err := ctx.EventManager().EmitTypedEvents(
		// A member was approved by a guardian
		&types.EventMemberApproved{
			MemberAddress:   msg.Member,
			ApproverAddress: msg.Approver,
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgApproveMemberResponse{}, nil
}
