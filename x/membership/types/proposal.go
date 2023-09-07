package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
	gov_v1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeAddGuardians            = "AddGuardians"
	ProposalTypeRemoveGuardians         = "RemoveGuardians"
	ProposalTypeUpdateTotalVotingWeight = "UpdateTotalVotingWeight"
)

// Ensure all proposals implement govtypes.Content at compile time
var (
	_ gov_v1beta1.Content = &AddGuardiansProposal{}
	_ gov_v1beta1.Content = &RemoveGuardiansProposal{}
	_ gov_v1beta1.Content = &UpdateTotalVotingWeightProposal{}
)

func init() {
	gov_v1beta1.RegisterProposalType(ProposalTypeAddGuardians)
	gov_v1beta1.RegisterProposalType(ProposalTypeRemoveGuardians)
	gov_v1beta1.RegisterProposalType(ProposalTypeUpdateTotalVotingWeight)
}

////////
// Add Guardian Proposal
////////

// NewAddGuardiansProposal creates an empty proposal instance
func NewAddGuardiansProposal(title string, description string, creator string, guardiansToAdd []string) gov_v1beta1.Content {
	return &AddGuardiansProposal{
		Title:          title,
		Description:    description,
		Creator:        creator,
		GuardiansToAdd: guardiansToAdd,
	}
}

// GetTitle returns the title of a add guardians proposal.
func (p *AddGuardiansProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of a add guardians proposal.
func (p *AddGuardiansProposal) GetDescription() string { return p.Description }

// ProposalRoute ensures this proposal will be handled by the Membership Module
func (p *AddGuardiansProposal) ProposalRoute() string { return ModuleName }

// ProposalType defines the type for a AddGuardiansProposal
func (p *AddGuardiansProposal) ProposalType() string {
	return ProposalTypeAddGuardians
}

// ValidateBasic performs basic validation on the proposal
func (p *AddGuardiansProposal) ValidateBasic() error {
	if len(p.GuardiansToAdd) == 0 {
		return fmt.Errorf("no guardians to add")
	}
	if len(p.Creator) == 0 {
		return fmt.Errorf("creator address cannot be empty")
	}
	return nil
}

// String describes the proposal
func (p *AddGuardiansProposal) String() string {
	var b strings.Builder

	// Combine GuardiansToAdd into a CSV string
	var guardiansToAddCSV string
	if len(p.GuardiansToAdd) > 0 {
		guardiansToAddCSV = strings.Join(p.GuardiansToAdd, ", ")
	}

	b.WriteString(fmt.Sprintf(`Add Guardians Proposal:
  Title:              %s
  Description:        %s
  Guardians to Add:   %s
`, p.Title, p.Description, guardiansToAddCSV))
	return b.String()
}

////////
// Remove Guardian Proposal
////////

// NewRemoveGuardiansProposal creates an empty proposal instance
func NewRemoveGuardiansProposal(title string, description string, creator string, guardiansToRemove []string) gov_v1beta1.Content {
	return &RemoveGuardiansProposal{
		Title:             title,
		Description:       description,
		Creator:           creator,
		GuardiansToRemove: guardiansToRemove,
	}
}

// GetTitle returns the title of a remove guardians proposal.
func (p *RemoveGuardiansProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of a remove guardians proposal.
func (p *RemoveGuardiansProposal) GetDescription() string { return p.Description }

// ProposalRoute ensures this proposal will be handled by the Membership Module
func (p *RemoveGuardiansProposal) ProposalRoute() string { return ModuleName }

// ProposalType defines the type for a RemoveGuardiansProposal
func (p *RemoveGuardiansProposal) ProposalType() string {
	return ProposalTypeRemoveGuardians
}

// ValidateBasic performs basic validation on the proposal
func (p *RemoveGuardiansProposal) ValidateBasic() error {
	if len(p.GuardiansToRemove) == 0 {
		return fmt.Errorf("no guardians to remove")
	}
	if len(p.Creator) == 0 {
		return fmt.Errorf("creator address cannot be empty")
	}
	return nil
}

// String describes the proposal
func (p *RemoveGuardiansProposal) String() string {
	var b strings.Builder

	// Combine GuardiansToRemove into a CSV string
	var guardiansToRemoveCSV string
	if len(p.GuardiansToRemove) > 0 {
		guardiansToRemoveCSV = strings.Join(p.GuardiansToRemove, ", ")
	}

	b.WriteString(fmt.Sprintf(`Remove Guardians Proposal:
  Title:              %s
  Description:        %s
  Guardians to Remove:%s
`, p.Title, p.Description, guardiansToRemoveCSV))
	return b.String()
}

////////
// Update Total Voting Weight Proposal
////////

// NewUpdateTotalVotingWeightProposal creates an empty proposal instance
func NewUpdateTotalVotingWeightProposal(title string, description string, creator string, totalVotingWeight math.LegacyDec) gov_v1beta1.Content {
	return &UpdateTotalVotingWeightProposal{
		Title:                title,
		Description:          description,
		Creator:              creator,
		NewTotalVotingWeight: totalVotingWeight,
	}
}

// GetTitle returns the title of a update total voting weight proposal.
func (p *UpdateTotalVotingWeightProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of a update total voting weight proposal.
func (p *UpdateTotalVotingWeightProposal) GetDescription() string { return p.Description }

// ProposalRoute ensures this proposal will be handled by the Membership Module
func (p *UpdateTotalVotingWeightProposal) ProposalRoute() string { return ModuleName }

// ProposalType defines the type for a UpdateTotalVotingWeightProposal
func (p *UpdateTotalVotingWeightProposal) ProposalType() string {
	return ProposalTypeUpdateTotalVotingWeight
}

// ValidateBasic performs basic validation on the proposal
func (p *UpdateTotalVotingWeightProposal) ValidateBasic() error {
	if len(p.Creator) == 0 {
		return fmt.Errorf("creator address cannot be empty")
	}
	if p.NewTotalVotingWeight.LTE(math.LegacyZeroDec()) {
		return fmt.Errorf("total voting weight cannot be <= 0")
	}
	if p.NewTotalVotingWeight.GT(math.LegacyOneDec()) {
		return fmt.Errorf("total voting weight cannot be > 1")
	}
	return nil
}

// String describes the proposal
func (p *UpdateTotalVotingWeightProposal) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`Update Total Voting Weight Proposal:
  Title:                   %s
  Description:             %s
  New Total Voting Weight: %s
`, p.Title, p.Description, p.NewTotalVotingWeight))
	return b.String()
}
