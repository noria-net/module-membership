#!/bin/bash

echo "GOV TEST: Add new guardian via proposal"

# Me is a regular member
ADDRESS_ME=mm1dmh80jwx0mv5khvqdz9sj28dmuhvems97wq628
# VAL1 is a guardian
ADDRESS_VAL1=mm1e7gp56hf85nk0qtg0542gmmmwq753ww2tg7dws
PROPOSAL_TEXT="proposal_add_guardian.json"

# Submit the proposal
./tx-gas.sh membershipd tx membership update-direct-democracy $PROPOSAL_TEXT \
  --from me \
  --deposit 1000000unoria

# Get the ID of the  latest proposal
PROPOSAL_ID=$(membershipd query gov proposals --output json --reverse --limit 1 | jq -r '.proposals[].id')

# deposit to a proposal
echo "Depositing to proposal $PROPOSAL_ID"
./tx-gas.sh membershipd tx gov deposit $PROPOSAL_ID 1000000unoria --from val1

echo "val2 voting Yes on proposal $PROPOSAL_ID"
./tx-gas.sh membershipd tx gov vote $PROPOSAL_ID Yes --from val1

# Wait for it to pass or fail
./wait_for_proposal_complete.sh $PROPOSAL_ID