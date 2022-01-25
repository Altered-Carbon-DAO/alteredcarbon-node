set -ex
set -o pipefail
DENOM=uacarb
CHAINID=alteredcarbon
RLYKEY=acarb12g0xe2ld0k5ws3h7lmxc39d4rpl3fyxp5qys69
acarbd version --long
apk add -U --no-cache jq tree
STARGAZE_HOME=/alteredcarbon/acarbd

# Setup alteredcarbon
acarbd init --chain-id $CHAINID $CHAINID --home $STARGAZE_HOME
acarbd config keyring-backend test --home $STARGAZE_HOME
sed -i 's#tcp://127.0.0.1:26657#tcp://0.0.0.0:26657#g' $STARGAZE_HOME/config/config.toml
sed -i "s/\"stake\"/\"$DENOM\"/g" $STARGAZE_HOME/config/genesis.json
sed -i 's/pruning = "syncable"/pruning = "nothing"/g' $STARGAZE_HOME/config/app.toml
sed -i 's/enable = false/enable = true/g' $STARGAZE_HOME/config/app.toml
sed -i 's/172800s/60s/g'  $STARGAZE_HOME/config/genesis.json
acarbd keys --keyring-backend test add validator --home $STARGAZE_HOME
acarbd add-genesis-account $(acarbd keys --keyring-backend test show validator -a --home $STARGAZE_HOME) 10000000000000$DENOM --home $STARGAZE_HOME
acarbd add-genesis-account $RLYKEY 10000000000000$DENOM --home $STARGAZE_HOME
acarbd gentx validator 900000000$DENOM --keyring-backend test --chain-id alteredcarbon --home $STARGAZE_HOME
acarbd collect-gentxs --home $STARGAZE_HOME
/alteredcarbon/bin/upgrade-watcher acarbd start --pruning nothing --home $STARGAZE_HOME
