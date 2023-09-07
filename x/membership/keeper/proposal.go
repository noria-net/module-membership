package keeper

import (
	errors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/noria-net/module-membership/x/membership/types"
)

// HandleAddGuardiansProposal adds new guardians when the proposal passes
func HandleAddGuardiansProposal(ctx sdk.Context, k Keeper, p *types.AddGuardiansProposal) error {
	// Must be a valid creator address
	creator, err := sdk.AccAddressFromBech32(p.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	// Only guardians can create this proposal
	if !k.IsGuardian(ctx, creator) {
		return errors.Wrapf(sdkerrors.ErrUnauthorized, "creator is not a guardian")
	}

	// Create an array of guardian addresses to add
	var guardiansToAdd []sdk.AccAddress

	// Iterate through the GuardiansToAdd and ensure there are
	// no empty or invalid addresses
	// The whole proposal fails if any of the addresses are invalid
	for _, addr := range p.GuardiansToAdd {
		// Get the member
		member, err := validateAndFetchMember(ctx, k, addr)
		if err != nil {
			return err
		}

		// Ensure this member's electorate status is active
		if member.Status != types.MembershipStatus_MemberElectorate {
			return errors.Wrapf(sdkerrors.ErrUnauthorized, "member is not active: %s", addr)
		}
		// Ensure this member is not already a guardian
		if member.IsGuardian {
			return errors.Wrapf(sdkerrors.ErrUnauthorized, "member is already a guardian: %s", addr)
		}
		// Add to the list
		guardiansToAdd = append(guardiansToAdd, sdk.MustAccAddressFromBech32(addr))
	}

	// Ensure we have more than one guardian to add
	if len(guardiansToAdd) < 1 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "no guardians to add")
	}

	// Loop through guardiansToAdd and add them to the guardian set
	for _, addr := range guardiansToAdd {
		// Add to the guardian set
		err := k.SetMemberGuardianStatus(ctx, addr, true)
		if err != nil {
			return err
		}
	}

	// Add these guardians to the DirectDemocracySettings
	dd := k.GetDirectDemocracySettings(ctx)
	dd.Guardians = append(dd.Guardians, p.GuardiansToAdd...)
	k.SetDirectDemocracySettings(ctx, dd)

	return nil
}

// HandleRemoveGuardiansProposal removes guardians when the proposal passes
func HandleRemoveGuardiansProposal(ctx sdk.Context, k Keeper, p *types.RemoveGuardiansProposal) error {
	// Must be a valid creator address
	creator, err := sdk.AccAddressFromBech32(p.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	// Only guardians can create this proposal
	if !k.IsGuardian(ctx, creator) {
		return errors.Wrapf(sdkerrors.ErrUnauthorized, "creator is not a guardian")
	}

	// Create an array of guardian addresses to remove
	var guardiansToRemove []sdk.AccAddress

	// Iterate through the GuardiansToRemove and ensure there are
	// no empty or invalid addresses
	// The whole proposal fails if any of the addresses are invalid
	for _, addr := range p.GuardiansToRemove {
		// Get the member
		member, err := validateAndFetchMember(ctx, k, addr)
		if err != nil {
			return err
		}

		// Ensure this member is a guardian
		if !member.IsGuardian {
			return errors.Wrapf(sdkerrors.ErrUnauthorized, "member is not a guardian: %s", addr)
		}
		// Add to the list
		guardiansToRemove = append(guardiansToRemove, sdk.MustAccAddressFromBech32(addr))
	}

	// Ensure we have more than one guardian to remove
	if len(guardiansToRemove) < 1 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "no guardians to remove")
	}

	// Loop through guardiansToRemove and remove their guardian status
	for _, addr := range guardiansToRemove {
		// Remove from the guardian set
		err := k.SetMemberGuardianStatus(ctx, addr, false)
		if err != nil {
			return err
		}
	}

	// Remove these guardians from the DirectDemocracySettings
	dd := k.GetDirectDemocracySettings(ctx)
	dd.Guardians = removeFromSlice(dd.Guardians, p.GuardiansToRemove)
	k.SetDirectDemocracySettings(ctx, dd)

	return nil
}

// HandleUpdateTotalVotingWeightProposal updates the total voting weight when the proposal passes
func HandleUpdateTotalVotingWeightProposal(ctx sdk.Context, k Keeper, p *types.UpdateTotalVotingWeightProposal) error {
	// Creator must be a guardian
	creator, err := sdk.AccAddressFromBech32(p.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	// Only guardians can create this proposal
	if !k.IsGuardian(ctx, creator) {
		return errors.Wrapf(sdkerrors.ErrUnauthorized, "creator is not a guardian")
	}

	// Ensure the total voting weight is > 0 and <= 1
	if p.NewTotalVotingWeight.LT(sdk.ZeroDec()) || p.NewTotalVotingWeight.GT(sdk.OneDec()) {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "total voting weight must be > 0 and <= 1")
	}

	// Update the total voting weight
	dd := k.GetDirectDemocracySettings(ctx)
	// Save the old value and update to the new value
	oldTotalVotingWeight := dd.TotalVotingWeight
	dd.TotalVotingWeight = p.NewTotalVotingWeight
	k.SetDirectDemocracySettings(ctx, dd)

	// Emit an event saying the total voting weight has changed
	err = ctx.EventManager().EmitTypedEvent(
		&types.EventTotalVotingWeightChanged{
			OldTotalVotingWeight: oldTotalVotingWeight,
			NewTotalVotingWeight: p.NewTotalVotingWeight,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// validateAndFetchMember ensures the address is valid and returns the member
func validateAndFetchMember(ctx sdk.Context, k Keeper, addr string) (*types.Member, error) {
	// Address cannot be empty
	if addr == "" {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidAddress, "empty guardian address")
	}
	// Address must be valid
	bechAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid guardian address: %s", err)
	}
	// Ensure this address is for an existing member
	if !k.IsMember(ctx, bechAddr) {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "member not found at this address: %s", addr)
	}
	// Get the member
	member, _ := k.GetMemberAccount(ctx, bechAddr)
	return &member, nil
}

// removeFromSlice excludes itemsToRemove from slice
func removeFromSlice(slice []string, itemsToRemove []string) []string {
	var result []string
	for _, s := range slice {

		for _, item := range itemsToRemove {
			if s == item {
				continue
			}
		}
		result = append(result, s)
	}
	return result
}
