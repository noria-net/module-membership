package keeper

import (
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noria-net/module-membership/x/membership/types"
)

func (k Keeper) SetMemberGuardianStatus(ctx sdk.Context, addr sdk.AccAddress, isGuardian bool) error {
	member, found := k.GetMemberAccount(ctx, addr)
	if !found {
		return errors.Wrapf(types.ErrMemberNotFound, "member not found: %s", addr.String())
	}

	// Make sure the status has actually changed
	if member.IsGuardian == isGuardian {
		return nil
	}

	// Set the guardianship status
	member.IsGuardian = isGuardian
	k.UpdateMember(ctx, member)

	// Publish an event for this change
	if isGuardian {
		ctx.EventManager().EmitTypedEvent(
			&types.EventMemberGrantedGuardianship{
				MemberAddress: addr.String(),
			},
		)
	} else {
		// Publish an event
		ctx.EventManager().EmitTypedEvent(
			&types.EventMemberRevokedGuardianship{
				MemberAddress: addr.String(),
			},
		)
	}

	return nil
}

// GetGuardians returns all guardians of the electorate
// NOTE: Only valid members with membership status of MemberElectorate are returned
func (k Keeper) GetGuardians(ctx sdk.Context) (guardians []*types.Member) {

	dd := k.GetDirectDemocracySettings(ctx)

	for _, guardianAddress := range dd.Guardians {

		// Unmarshal the address and panic if there's an error
		// NOTE: The key is prefixed with the guardian key prefix
		// so we need to remove it before unmarshalling
		acc := sdk.MustAccAddressFromBech32(guardianAddress)

		// Get the member
		member, found := k.GetMemberAccount(ctx, acc)
		// Guardian must be a member with a status of MemberElectorate
		if found && member.IsGuardian && member.Status == types.MembershipStatus_MemberElectorate {
			guardians = append(guardians, &member)
		}
		// TODO: emit an event when encountering an invalid guardian
	}

	return guardians
}

// IsGuardian returns true if:
// - The address is a guardian
// - The address is a member
// - The member has a status of MemberElectorate
func (k Keeper) IsGuardian(ctx sdk.Context, addr sdk.AccAddress) bool {
	dd := k.GetDirectDemocracySettings(ctx)

	for _, guardianAddress := range dd.Guardians {

		// Unmarshal the address and panic if there's an error
		// NOTE: The key is prefixed with the guardian key prefix
		// so we need to remove it before unmarshalling
		acc := sdk.MustAccAddressFromBech32(guardianAddress)

		if acc.Equals(addr) {
			// Get the member
			member, found := k.GetMemberAccount(ctx, acc)
			// Guardian must be a member with a status of MemberElectorate
			return found &&
				member.IsGuardian &&
				member.Status == types.MembershipStatus_MemberElectorate
		}
	}
	return false
}

func (k Keeper) GetDirectDemocracySettings(ctx sdk.Context) *types.DirectDemocracy {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	b := store.Get(types.DirectDemocracyKey)
	if b == nil {
		return nil
	}
	democracy := new(types.DirectDemocracy)
	err := democracy.Unmarshal(b)
	if err != nil {
		panic(err)
	}
	return democracy
}

func (k Keeper) SetDirectDemocracySettings(ctx sdk.Context, democracy *types.DirectDemocracy) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	b, err := democracy.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(types.DirectDemocracyKey, b)
}
