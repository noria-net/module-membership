package keeper

import (
	errors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/noria-net/module-membership/x/membership/types"
)

// HandleAddGuardiansProposal adds new guardians when the proposal passes
func HandleAddGuardiansProposal(ctx sdk.Context, k Keeper, p *types.AddGuardiansProposal) error {
	// Creator must be a guardian
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
		// Address cannot be empty
		if addr == "" {
			return errors.Wrapf(sdkerrors.ErrInvalidAddress, "empty guardian address")
		}
		// Address must be valid
		bechAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid guardian address: %s", err)
		}
		// Ensure this address is for an existing member
		if !k.IsMember(ctx, bechAddr) {
			return errors.Wrapf(sdkerrors.ErrUnauthorized, "member not found at this address: %s", addr)
		}
		// Get the member
		member, _ := k.GetMemberAccount(ctx, bechAddr)
		// Ensure this member's electorate status is active
		if member.Status != types.MembershipStatus_MemberElectorate {
			return errors.Wrapf(sdkerrors.ErrUnauthorized, "member is not active: %s", addr)
		}
		// Ensure this member is not already a guardian
		if member.IsGuardian {
			return errors.Wrapf(sdkerrors.ErrUnauthorized, "member is already a guardian: %s", addr)
		}
		// Add to the list
		guardiansToAdd = append(guardiansToAdd, bechAddr)
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
