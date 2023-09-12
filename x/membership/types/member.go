package types

import (
	"strings"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NicknameMaxLength is the maximum number of characters allowed for
// a member's nickname
const NicknameMaxLength = 30

const MembershipStatusPrefix = "MEMBERSHIP_STATUS_"
const MemberMetadata_Nickname = "nickname"

// Members is a collection of Member objects
type Members []*Member

// AllowedMembershipStatusTransitions holds a truth table of all permissable membership status changes
var AllowedMembershipStatusTransitions = map[MembershipStatus][]MembershipStatus{
	MembershipStatus_MemberElectorate: {MembershipStatus_MemberInactive, MembershipStatus_MemberExpulsed},
	MembershipStatus_MemberInactive:   {MembershipStatus_MemberElectorate},
	MembershipStatus_MemberRecalled:   {MembershipStatus_MemberElectorate},
	MembershipStatus_MemberExpulsed:   {MembershipStatus_MemberElectorate},
}

// NewMemberAccountWithDefaultMemberStatus creats a new member account with a
// default member status of Electorate.
func NewMemberAccountWithDefaultMemberStatus(baseAccount *authtypes.BaseAccount) *Member {
	acc := &Member{
		BaseAccount: baseAccount,
		Status:      MembershipStatus_MemberStatusPendingApproval,
	}
	return acc
}

// NewMemberAccountWithStatus parses the raw status and returns a valid status value
func ParseMembershipStatus(s string) MembershipStatus {
	value, ok := MembershipStatus_value[s]
	if !ok {
		return MembershipStatus_MemberStatusEmpty
	}
	return MembershipStatus(value)
}

func ParseShortFormMembershipStatus(s string) MembershipStatus {
	status := MembershipStatusPrefix + strings.ToUpper(s)
	return ParseMembershipStatus(status)
}

func GetAllShortFormMembershipStatuses() []string {
	var statuses []string
	for _, status := range MembershipStatus_value {
		ms := MembershipStatus(status)
		if ms == MembershipStatus_MemberStatusEmpty {
			continue
		}
		statuses = append(statuses, ms.ToLowerCaseShortForm())
	}
	return statuses
}

func GetAllShortFormMembershipStatusesAsString() string {
	return strings.Join(GetAllShortFormMembershipStatuses(), ", ")
}

// CanTransitionTo returns true if the receiver MembershipStatus can transition to the desired MembershipStatus.
func (m MembershipStatus) CanTransitionTo(desired MembershipStatus) bool {
	// Get the valid targets for this transition
	validTargets, ok := AllowedMembershipStatusTransitions[m]
	if !ok {
		return false
	}
	for _, t := range validTargets {
		if t == desired {
			return true
		}
	}
	return false
}

// IsValid returns true if the status is within range and is not zero / empty
func (m MembershipStatus) IsValid() bool {
	_, ok := MembershipStatus_name[int32(m)]
	return ok && m != MembershipStatus_MemberStatusEmpty
}

func (m MembershipStatus) ToLowerCaseShortForm() string {
	name := MembershipStatus_name[int32(m)]
	return strings.ToLower(strings.TrimPrefix(name, MembershipStatusPrefix))
}

func (m MembershipStatus) DescribeTransition(to MembershipStatus) string {
	return "from " + m.ToLowerCaseShortForm() + " to " + to.ToLowerCaseShortForm()
}
