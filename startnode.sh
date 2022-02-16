#!/bin/sh

# create users
rm -rf $HOME/.acarbd
CHAINID=acarbd
DENOM=uacarb
acarbd config chain-id $CHAINID
GENESIS="$HOME/.acarbd/config/genesis.json"

echo "Setting up keyring-backend..."
acarbd config keyring-backend test
acarbd config output json
echo "Adding up validator..."
yes | acarbd keys add validator
echo "Adding up treasury..."
yes | acarbd keys add treasury
echo "Adding up founder1..."
yes | acarbd keys add founder1
echo "Adding up founder2..."
yes | acarbd keys add founder2
echo "Adding up founder3..."
yes | acarbd keys add founder3
echo "Adding up founder4..."
yes | acarbd keys add founder4

VALIDATOR=$(acarbd keys show validator -a)
TREASURY=$(acarbd keys show treasury -a)
FOUNDER1=$(acarbd keys show founder1 -a)
FOUNDER2=$(acarbd keys show founder2 -a)
FOUNDER3=$(acarbd keys show founder3 -a)
FOUNDER4=$(acarbd keys show founder4 -a)
echo "Got VALIDATOR $VALIDATOR"
echo "Got TREASURY $TREASURY"
echo "Got FOUNDER1 $FOUNDER1"
echo "Got FOUNDER2 $FOUNDER2"
echo "Got FOUNDER3 $FOUNDER3"
echo "Got FOUNDER4 $FOUNDER4"

# setup chain
acarbd init $CHAINID --chain-id $CHAINID
jq '.app_state.staking.params.bond_denom = "uacarb"' $GENESIS > temp.json && mv temp.json $GENESIS
jq '.app_state.crisis.constant_fee.denom = "uacarb"' $GENESIS > temp.json && mv temp.json $GENESIS
jq '.app_state.gov.min_deposit.denom = "uacarb"' $GENESIS > temp.json && mv temp.json $GENESIS
jq '.app_state.mint.params.mint_denom = "uacarb"' $GENESIS > temp.json && mv temp.json $GENESIS

# modify config for development
# config="$HOME/.acarbd/config/config.toml"
# if [ "$(uname)" = "Linux" ]; then
#   sed -e "s/cors_allowed_origins = \[\]/cors_allowed_origins = [\"*\"]/g" $config
# else
#   sed -e '' "s/cors_allowed_origins = \[\]/cors_allowed_origins = [\"*\"]/g" $config
# fi

acarbd add-genesis-account $VALIDATOR 100000000000uacarb
acarbd add-genesis-account $TREASURY 2500000000000000uacarb
acarbd add-genesis-account $FOUNDER1 250000000000000uacarb
acarbd add-genesis-account $FOUNDER2 250000000000000uacarb
acarbd add-genesis-account $FOUNDER3 250000000000000uacarb
acarbd add-genesis-account $FOUNDER4 250000000000000uacarb

acarbd prepare-genesis mainnet $CHAINID
acarbd gentx validator 10000000000uacarb --chain-id $CHAINID --keyring-backend test
acarbd collect-gentxs
acarbd validate-genesis
acarbd start
