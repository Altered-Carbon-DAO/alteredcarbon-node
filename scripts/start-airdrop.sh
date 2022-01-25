#!/bin/bash

DENOM=uacarb
CHAIN_ID=localnet-1
ONE_HOUR=3600
ONE_DAY=$(($ONE_HOUR * 24))
ONE_YEAR=$(($ONE_DAY * 365))
VALIDATOR_COINS=100000000000000$DENOM

rm -rf $HOME/.acarbd

if [ "$1" == "mainnet" ]
then
    LOCKUP=ONE_YEAR
else
    LOCKUP=ONE_DAY
fi
echo "Lockup period is $LOCKUP"

echo "Processing airdrop snapshot..."
if ! [ -f genesis.json ]; then
    curl -O https://archive.interchain.io/4.0.2/genesis.json
fi
acarbd export-airdrop-snapshot uatom genesis.json snapshot.json
acarbd init testmoniker --chain-id $CHAIN_ID
acarbd prepare-genesis testnet $CHAIN_ID
acarbd import-genesis-accounts-from-snapshot snapshot.json

acarbd config chain-id localnet-1
acarbd config keyring-backend test
acarbd config output json
yes | acarbd keys add validator

acarbd add-genesis-account $(acarbd keys show validator -a) $VALIDATOR_COINS

echo "Adding vesting accounts..."
GENESIS_TIME=$(jq '.genesis_time' ~/.acarbd/config/genesis.json | tr -d '"')
echo "Genesis time is $GENESIS_TIME"
if [[ "$OSTYPE" == "darwin"* ]]; then
    GENESIS_UNIX_TIME=$(TZ=UTC gdate "+%s" -d $GENESIS_TIME)
else
    GENESIS_UNIX_TIME=$(TZ=UTC date "+%s" -d $GENESIS_TIME)
fi
vesting_start_time=$(($GENESIS_UNIX_TIME + $LOCKUP))
vesting_end_time=$(($vesting_start_time + $LOCKUP))

acarbd add-genesis-account acarb1s4ckh9405q0a3jhkwx9wkf9hsjh66nmuu53dwe 350000000000000$DENOM
acarbd add-genesis-account acarb13nh557xzyfdm6csyp0xslu939l753sdlgdc2q0 250000000000000$DENOM
acarbd add-genesis-account acarb12yxedm78tpptyhhasxrytyfyj7rg7dcqfgrdk4 16666666666667$DENOM \
    --vesting-amount 16666666666667$DENOM \
    --vesting-start-time $vesting_start_time \
    --vesting-end-time $vesting_end_time
acarbd add-genesis-account acarb1nek5njjd7uqn5zwf5zyl3xhejvd36er3qzp6x3 16666666666667$DENOM \
    --vesting-amount 16666666666667$DENOM \
    --vesting-start-time $vesting_start_time \
    --vesting-end-time $vesting_end_time
acarbd add-genesis-account acarb1avlcqcn4hsxrds2dgxmgrj244hu630kfl89vrt 16666666666667$DENOM \
    --vesting-amount 16666666666667$DENOM \
    --vesting-start-time $vesting_start_time \
    --vesting-end-time $vesting_end_time
acarbd add-genesis-account acarb1wppujuuqrv52atyg8uw3x779r8w72ehrr5a4yx 50000000000000$DENOM \
    --vesting-amount 50000000000000$DENOM \
    --vesting-start-time $vesting_start_time \
    --vesting-end-time $vesting_end_time

acarbd gentx validator 1000000000000uacarb --chain-id localnet-1 --keyring-backend test
acarbd collect-gentxs
acarbd validate-genesis
acarbd start
