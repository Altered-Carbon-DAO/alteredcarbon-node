#!/bin/sh

# create users
rm -rf $HOME/.alteredcarbon
CHAINID=alteredcarbon
DENOM=uacarb
alteredcarbond config chain-id $CHAINID
GENESIS="$HOME/.alteredcarbon/config/genesis.json"

echo "Setting up keyring-backend..."
alteredcarbond config keyring-backend test
alteredcarbond config output json
echo "Adding up validator..."
yes | alteredcarbond keys add validator
echo "Adding up treasury..."
yes | alteredcarbond keys add treasury
echo "Adding up founder1..."
yes | alteredcarbond keys add founder1
echo "Adding up founder2..."
yes | alteredcarbond keys add founder2
echo "Adding up founder3..."
yes | alteredcarbond keys add founder3
echo "Adding up founder4..."
yes | alteredcarbond keys add founder4

VALIDATOR=$(alteredcarbond keys show validator -a)
TREASURY=$(alteredcarbond keys show treasury -a)
FOUNDER1=$(alteredcarbond keys show founder1 -a)
FOUNDER2=$(alteredcarbond keys show founder2 -a)
FOUNDER3=$(alteredcarbond keys show founder3 -a)
FOUNDER4=$(alteredcarbond keys show founder4 -a)
echo "Got VALIDATOR $VALIDATOR"
echo "Got TREASURY $TREASURY"
echo "Got FOUNDER1 $FOUNDER1"
echo "Got FOUNDER2 $FOUNDER2"
echo "Got FOUNDER3 $FOUNDER3"
echo "Got FOUNDER4 $FOUNDER4"

# setup chain
alteredcarbond init $CHAINID --chain-id $CHAINID
jq '.app_state.staking.params.bond_denom = "blah"' $GENESIS > temp.json && mv temp.json $GENESIS
jq '.app_state.crisis.constant_fee.denom = "blah"' $GENESIS > temp.json && mv temp.json $GENESIS
jq '.app_state.gov.min_deposit.denom = "blah"' $GENESIS > temp.json && mv temp.json $GENESIS
jq '.app_state.mint.params.mint_denom = "blah"' $GENESIS > temp.json && mv temp.json $GENESIS

# modify config for development
config="$HOME/.alteredcarbond/config/config.toml"
if [ "$(uname)" = "Linux" ]; then
  sed -e "s/cors_allowed_origins = \[\]/cors_allowed_origins = [\"*\"]/g" $config
else
  sed -e '' "s/cors_allowed_origins = \[\]/cors_allowed_origins = [\"*\"]/g" $config
fi

alteredcarbond add-genesis-account $VALIDATOR 1000000000000uacarb
alteredcarbond add-genesis-account $TREASURY 2500000000000000uacarb

# alteredcarbond add-genesis-account $FOUNDER1 2500000000000000uacarb
# alteredcarbond add-genesis-account $FOUNDER2 250000000000000uacarb
# alteredcarbond add-genesis-account $FOUNDER3 250000000000000uacarb
# TODO allocate coins to FOUNDER4

alteredcarbond prepare-genesis mainnet $CHAINID
alteredcarbond gentx validator 10000000000uacarb --chain-id $CHAINID --keyring-backend test
alteredcarbond collect-gentxs
alteredcarbond validate-genesis
alteredcarbond start
