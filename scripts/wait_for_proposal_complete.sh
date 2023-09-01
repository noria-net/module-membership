# Check if $PROPOSAL_ID was passed in as an argument, error out if not
if [ -z "$1" ]; then
    echo "ERROR: No proposal ID provided"
    exit 1
fi
PROPOSAL_ID=$1

# Loop until the STATUS is either "PROPOSAL_STATUS_REJECTED" or "PROPOSAL_STATUS_PASSED"
while true; do
    STATUS=$(membershipd query gov proposal $PROPOSAL_ID --output json | jq -r '.status')
    echo "Proposal status: $STATUS"
    if [ "$STATUS" = "PROPOSAL_STATUS_REJECTED" ] || [ "$STATUS" = "PROPOSAL_STATUS_PASSED" ]; then
        break
    fi
    sleep 2
done