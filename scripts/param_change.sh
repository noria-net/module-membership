#!/bin/bash

PROPOSAL_FILE=$1
KEY_NAME=$2
BINARY_DIR=".module-membership"
CHAIN_ID="mmchain-1"
DENOM="unoria"
GAS_PRICE_DENOM="ucrd"
GAS_PRICE="0.0025"
NODE="http://127.0.0.1:26657/"
export DAEMON_NAME="membershipd"
export DAEMON_HOME="$HOME/$BINARY_DIR"

# Block height and sender key name must be set
if [ -z "$1" ] && [ -z "$2" ]; then
  echo "Parameter file and/or sender key name is missing"
  exit 1
fi

exe() { echo "EXE\$ $@" ; ./scripts/tx.sh "$@" ; }

# submit parameter change proposal
exe $DAEMON_NAME tx gov submit-legacy-proposal param-change $PROPOSAL_FILE \
  --from $KEY_NAME \
  --home $DAEMON_HOME \
  --node $NODE \
  --yes \
  --gas-prices $GAS_PRICE$GAS_PRICE_DENOM \
  --gas auto \
  --gas-adjustment 1.5

# Get the proposal ID by fetching the most recent proposal
PROPOSAL_ID=$($DAEMON_NAME q gov proposals limit 1 --reverse --output json --home $DAEMON_HOME --node $NODE | jq '.proposals[0].id | tonumber')

# vote on the proposal
exe $DAEMON_NAME tx gov vote \
  $PROPOSAL_ID \
  yes \
  --from $KEY_NAME \
  --home $DAEMON_HOME \
  --node $NODE \
  --yes \
  --gas-prices $GAS_PRICE$GAS_PRICE_DENOM \
  --gas auto \
  --gas-adjustment 2

$DAEMON_NAME q gov votes $PROPOSAL_ID --output json | jq
