package keeper

import (
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/noria-net/module-membership/x/membership/types"
)

func (k Keeper) GetMemberAccount(ctx sdk.Context, address sdk.AccAddress) (types.Member, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	key := types.MemberKey(address)
	var member types.Member

	var b []byte
	if b = store.Get(key); b == nil {
		return member, false
	}

	if err := k.cdc.Unmarshal(b, &member); err != nil {
		panic(err)
	}

	// Validate guardianship status
	member.IsGuardian = member.IsGuardian &&
		member.Status == types.MembershipStatus_MemberElectorate

	return member, true
}

func (k Keeper) AppendMember(ctx sdk.Context, address sdk.AccAddress) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	key := types.MemberKey(address)

	// Must not already have a member account
	if k.IsMember(ctx, address) {
		return errors.Wrap(sdkerrors.ErrUnauthorized, "account has already been enrolled")
	}

	// Get or create a base account
	var baseAccount = k.accountKeeper.GetAccount(ctx, address)
	if baseAccount == nil {
		// Create a base baseAccount
		baseAccount = k.accountKeeper.NewAccountWithAddress(ctx, address)
		// Ensure account type is correct
		if _, ok := baseAccount.(*authtypes.BaseAccount); !ok {
			return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid account type; expected: BaseAccount, got: %T", baseAccount)
		}
		// Save the base account
		k.accountKeeper.SetAccount(ctx, baseAccount)
	}

	// Create a member account
	newMember := types.NewMemberAccountWithDefaultMemberStatus(
		baseAccount.(*authtypes.BaseAccount),
	)

	// Fetch member counts
	memberCount := k.GetMemberCount(ctx)
	memberStatusCount := k.GetMemberStatusCount(ctx, newMember.Status)

	// Marshal and Set
	memberData := k.cdc.MustMarshal(newMember)
	store.Set(key, memberData)

	// Bump member count
	k.SetMemberCount(ctx, memberCount+1)

	// Bump member status count
	k.SetMemberStatusCount(ctx, newMember.Status, memberStatusCount+1)

	return nil
}

func (k Keeper) UpdateMember(ctx sdk.Context, member types.Member) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})

	// Parse the address
	address, _ := sdk.AccAddressFromBech32(member.Address)
	key := types.MemberKey(address)

	// Fetch old member
	oldMember, found := k.GetMemberAccount(ctx, address)
	if !found {
		panic(errors.Wrapf(types.ErrMemberNotFound, "member not found: %s", address.String()))
	}

	// Marshal and Set
	memberData := k.cdc.MustMarshal(&member)
	store.Set(key, memberData)

	// Update member status count if the status has changed
	if oldMember.Status != member.Status {
		// Fetch member counts
		memberStatusCount := k.GetMemberStatusCount(ctx, oldMember.Status)
		k.SetMemberStatusCount(ctx, oldMember.Status, memberStatusCount-1)
		k.SetMemberStatusCount(ctx, member.Status, memberStatusCount+1)
	}
}

// SetMemberNickname sets the nickname of a member
// NOTE: Assumes the member exists
func (k Keeper) SetMemberNickname(ctx sdk.Context, address sdk.AccAddress, nickname string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	key := types.MemberMetadataKey(address, types.MemberMetadata_Nickname)
	store.Set(key, []byte(nickname))
}

// GetMemberNickname gets the nickname of a member
func (k Keeper) GetMemberNickname(ctx sdk.Context, address sdk.AccAddress) string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	key := types.MemberMetadataKey(address, types.MemberMetadata_Nickname)
	bz := store.Get(key)
	if bz == nil {
		return ""
	}

	return string(bz)
}

func (k Keeper) IsMember(ctx sdk.Context, address sdk.AccAddress) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	key := types.MemberKey(address)
	return store.Has(key)
}

func (k Keeper) GetMemberCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	bz := store.Get(types.MemberCountKey)
	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetMemberCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	bz := sdk.Uint64ToBigEndian(count)
	store.Set(types.MemberCountKey, bz)
}

func (k Keeper) GetMemberStatusCount(ctx sdk.Context, s types.MembershipStatus) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	bz := store.Get(types.MemberStatusCountKey(s))
	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetMemberStatusCount(ctx sdk.Context, s types.MembershipStatus, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	bz := sdk.Uint64ToBigEndian(count)
	store.Set(types.MemberStatusCountKey(s), bz)
}

func (k Keeper) UpdateMemberStatus(ctx sdk.Context, target sdk.AccAddress, s types.MembershipStatus) error {
	// Fetch the member
	member, found := k.GetMemberAccount(ctx, target)

	// Member must exist
	if !found {
		return errors.Wrapf(types.ErrMemberNotFound, "member not found: %s", target.String())
	}

	// Must be a valid status transition
	if !member.Status.CanTransitionTo(s) {
		return errors.Wrapf(types.ErrMembershipStatusChangeNotAllowed, "transition %s is not allowed", s.DescribeTransition(s))
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	key := types.MemberKey(target)

	// Fetch status counts
	oldStatus := member.Status
	newStatus := s
	oldStatusCount := k.GetMemberStatusCount(ctx, oldStatus)
	newStatusCount := k.GetMemberStatusCount(ctx, newStatus)

	// Update the member's status
	member.Status = s

	// Marshal and Set
	memberData := k.cdc.MustMarshal(&member)
	store.Set(key, memberData)

	// Update the status counts
	k.SetMemberStatusCount(ctx, oldStatus, oldStatusCount-1)
	k.SetMemberStatusCount(ctx, newStatus, newStatusCount+1)

	// Publish an update event
	ctx.EventManager().EmitTypedEvent(
		// A member's citizenship status has changed
		&types.EventMemberStatusChanged{
			MemberAddress: target.String(),
			// TODO: Change this
			Operator:       "",
			Status:         types.MembershipStatus_MemberElectorate,
			PreviousStatus: types.MembershipStatus_MemberStatusEmpty,
		},
	)

	return nil
}
