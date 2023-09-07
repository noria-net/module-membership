#!/bin/bash

echo "GOV TEST: Update total voting weight via proposal"

# VAL1 is a guardian
ADDRESS_VAL1=mm1e7gp56hf85nk0qtg0542gmmmwq753ww2tg7dws

# Show current total voting weight
membershipd query membership guardians | grep total_voting_weight

# Submit the proposal
./tx-gas.sh membershipd tx gov submit-legacy-proposal update-total-voting-weight \
  0.6 \
  --title "Update the total voting weight" \
  --description "Update the total voting weight for guardians" \
  --from val1 \
  --deposit 1000000unoria

# Get the ID of the  latest proposal
PROPOSAL_ID=$(membershipd query gov proposals --output json --reverse --limit 1 | jq -r '.proposals[].id')

# deposit to a proposal
echo "Depositing to proposal $PROPOSAL_ID"
./tx-gas.sh membershipd tx gov deposit $PROPOSAL_ID 1000000unoria --from val1

echo "val1 voting Yes on proposal $PROPOSAL_ID"
./tx-gas.sh membershipd tx gov vote $PROPOSAL_ID Yes --from val1

# Wait for it to pass or fail
./wait_for_proposal_complete.sh $PROPOSAL_ID

echo "Show updated voting weight"
membershipd query membership guardians | grep total_voting_weight
