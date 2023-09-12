package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/membership module sentinel errors
var (
	ErrInvalidMembershipStatus          = errors.Register(ModuleName, 2, "invalid membership status")
	ErrMemberNotFound                   = errors.Register(ModuleName, 3, "member not found")
	ErrMembershipStatusChangeNotAllowed = errors.Register(ModuleName, 4, "membership status change not allowed")
	ErrInvalidNickname                  = errors.Register(ModuleName, 5, "invalid nickname")
	ErrInvalidGuardianList              = errors.Register(ModuleName, 6, "invalid guardian list")
	ErrInvalidTotalVotingWeight         = errors.Register(ModuleName, 7, "invalid total voting weight")
	ErrVoterNotAMember                  = errors.Register(ModuleName, 8, "voter is not a member")
	ErrMemberNotEligibleToVote          = errors.Register(ModuleName, 9, "member is not eligible to vote")
	ErrInvalidVoteWeighting             = errors.Register(ModuleName, 10, "invalid vote weighting")
	ErrMemberNotPendingApproval         = errors.Register(ModuleName, 11, "member's status is not pending")
)
