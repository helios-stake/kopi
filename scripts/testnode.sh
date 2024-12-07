#!/bin/bash
set -ux

CHAIN_ID="testing"
USER="myaccount"
MONIKER=${MONIKER:-node001}
HIDE_LOGS="/dev/null"
# PASSWORD=${PASSWORD:-$1}
NODE_HOME="$PWD/.kopid"
GENESIS=$NODE_HOME/config/genesis.json
TMP_GENESIS=$NODE_HOME/config/tmp_genesis.json
ARGS="--keyring-backend test --home $NODE_HOME"
START_ARGS="--home $NODE_HOME --pruning=nothing  --minimum-gas-prices=0stake"

rm -rf $NODE_HOME

kopid init --chain-id "$CHAIN_ID" "$MONIKER" --home $NODE_HOME >$HIDE_LOGS

kopid keys add $USER $ARGS 2>&1 | tee account.txt
kopid keys add $USER-eth $ARGS --eth 2>&1 | tee account-eth.txt

# hardcode the validator account for this instance
kopid genesis add-genesis-account $USER "100000000000000stake" $ARGS

jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="stake"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["gov"]["voting_params"]["voting_period"]="45s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
#jq '.app_state["gov"]["params"]["voting_period"]="45s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

# submit a genesis validator tx
# Workraround for https://github.com/cosmos/cosmos-sdk/issues/8251
kopid genesis gentx $USER "10000000000000stake" --chain-id="$CHAIN_ID" -y $ARGS >$HIDE_LOGS

kopid genesis collect-gentxs --home $NODE_HOME >$HIDE_LOGS

#kopid start $START_ARGS