#!/bin/bash

echo "GOV TEST: Everyone votes yes, Proposal passes"

# Me is a regular member
ADDRESS_ME=mm1dmh80jwx0mv5khvqdz9sj28dmuhvems97wq628
# VAL1 is a guardian
ADDRESS_VAL1=mm1e7gp56hf85nk0qtg0542gmmmwq753ww2tg7dws
PROPOSAL_TEXT="proposal_text.json"

# Test sequence for proposing and voting with direct democratic tallying

# Enroll "me" and "val1"
echo "Enrolling me"
./enroll.sh me --ignore
echo "Enrolling val1"
./enroll.sh val1 --ignore

# Get me account membership
#membershipd query membership get-member $ADDRESS_ME \
#    --output json | jq --color-output

# submit a text proposal
echo "Submitting proposal"
./tx-gas.sh membershipd tx gov submit-proposal $PROPOSAL_TEXT --from val1

#membershipd query gov proposals \
#    --output json | jq --color-output

# Get the ID of the  latest proposal
PROPOSAL_ID=$(membershipd query gov proposals --output json --reverse --limit 1 | jq -r '.proposals[].id')

# deposit to a proposal
echo "Depositing to proposal $PROPOSAL_ID"
./tx-gas.sh membershipd tx gov deposit $PROPOSAL_ID 1000000unoria --from val1

# Get the proposal status
membershipd query gov proposal $PROPOSAL_ID --output json | jq --color-output

##########################
# Scenario: Me and Val1 vote "yes", proposal passes
##########################

# vote on a proposal
echo "val2 voting Yes on proposal $PROPOSAL_ID"
./tx-gas.sh membershipd tx gov vote $PROPOSAL_ID Yes --from val1
echo "me voting Yes on proposal $PROPOSAL_ID"
./tx-gas.sh membershipd tx gov vote $PROPOSAL_ID Yes --from me

# Wait for it to pass or fail
./wait_for_proposal_complete.sh $PROPOSAL_ID