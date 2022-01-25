set -ex
DENOM=uacarb
CHAINID=alteredcarbon
RLYKEY=acarb12g0xe2ld0k5ws3h7lmxc39d4rpl3fyxp5qys69
LEDGER_ENABLED=false make install
acarbd version --long



# Setup alteredcarbon
acarbd init --chain-id $CHAINID $CHAINID
sed -i 's#tcp://127.0.0.1:26657#tcp://0.0.0.0:26657#g' ~/.acarbd/config/config.toml
sed -i "s/\"stake\"/\"$DENOM\"/g" ~/.acarbd/config/genesis.json
sed -i 's/pruning = "syncable"/pruning = "nothing"/g' ~/.acarbd/config/app.toml
sed -i 's/enable = false/enable = true/g' ~/.acarbd/config/app.toml
acarbd keys --keyring-backend test add validator

acarbd add-genesis-account $(acarbd keys --keyring-backend test show validator -a) 100000000000$DENOM
acarbd add-genesis-account $RLYKEY 100000000000$DENOM
acarbd gentx validator 900000000$DENOM --keyring-backend test --chain-id alteredcarbon
acarbd collect-gentxs

acarbd start --pruning nothing
