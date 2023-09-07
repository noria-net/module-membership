package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/noria-net/module-membership/x/membership/client/cli"
)

var (
	AddGuardiansProposalHandler    = govclient.NewProposalHandler(cli.NewSubmitAddGuardiansProposal)
	RemoveGuardiansProposalHandler = govclient.NewProposalHandler(cli.NewSubmitRemoveGuardiansProposal)
)
