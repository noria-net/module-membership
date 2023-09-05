package types

import (
	"fmt"
	"strings"

	govcdc "github.com/cosmos/cosmos-sdk/x/gov/codec"
	gov_v1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeAddGuardians = "AddGuardians"
)

// Ensure all proposals implement govtypes.Content at compile time
var (
	_ gov_v1beta1.Content = &AddGuardiansProposal{}
)

func init() {
	gov_v1beta1.RegisterProposalType(ProposalTypeAddGuardians)
	govcdc.ModuleCdc.Amino.RegisterConcrete(&AddGuardiansProposal{}, "membership/AddGuardiansProposal", nil)
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
