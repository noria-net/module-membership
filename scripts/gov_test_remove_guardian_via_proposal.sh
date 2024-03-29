#!/bin/bash

echo "GOV TEST: Remove a guardian via proposal"

# VAL1 is a guardian
ADDRESS_VAL1=mm1e7gp56hf85nk0qtg0542gmmmwq753ww2tg7dws

# Submit the proposal
./tx-gas.sh membershipd tx gov submit-legacy-proposal remove-guardians \
  $ADDRESS_VAL1 \
  --title "remove a guardian" \
  --description "remove a guardian" \
  --from val1 \
  --deposit 1000000unoria

# Get the ID of the  latest proposal
PROPOSAL_ID=$(membershipd query gov proposals --output json --reverse --limit 1 | jq -r '.proposals[].id')

# deposit to a proposal
echo "Depositing to proposal $PROPOSAL_ID"
./tx-gas.sh membershipd tx gov deposit $PROPOSAL_ID 1000000unoria --from val1

echo "val1 voting Yes on proposal $PROPOSAL_ID"
./tx-gas.sh membershipd tx gov vote $PROPOSAL_ID Yes --from val1

echo "Listing all guardians"
membershipd query membership guardians | grep address

# Wait for it to pass or fail
./wait_for_proposal_complete.sh $PROPOSAL_ID

echo "Listing all guardians"
membershipd query membership guardians | grep address

