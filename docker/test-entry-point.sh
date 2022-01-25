#!/bin/sh

CHAINID=$1
GENACCT=$2

if [ -z "$1" ]; then
  echo "Need to input chain id..."
  exit 1
fi

if [ -z "$2" ]; then
  echo "Need to input genesis account address..."
  exit 1
fi

# Build genesis file incl account for passed address
coins="100000000000uacarb,100000000000ucredits"
acarbd init --chain-id $CHAINID $CHAINID
acarbd keys add validator --keyring-backend="test"
acarbd add-genesis-account validator $coins --keyring-backend="test"
acarbd add-genesis-account $GENACCT $coins --keyring-backend="test"
acarbd gentx --name validator --amount 1000000uacarb --keyring-backend="test"
acarbd collect-gentxs


# Set proper defaults and change ports
sed -i 's/"leveldb"/"goleveldb"/g' ~/.acarbd/config/config.toml
sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.acarbd/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.acarbd/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.acarbd/config/config.toml
sed -i 's/index_all_keys = false/index_all_keys = true/g' ~/.acarbd/config/config.toml

# Start the stake
acarbd start --pruning=nothing
