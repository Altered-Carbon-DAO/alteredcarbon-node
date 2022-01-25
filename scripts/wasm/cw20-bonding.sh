#!/bin/sh

TXFLAG="--gas-prices 0.01uacarb --gas auto --gas-adjustment 1.3 -y -b block"

CREATOR=$(acarbd keys show creator -a)
INVESTOR=$(acarbd keys show investor -a)

# see contracts code that have been uploaded
acarbd q wasm list-code

# download cw20-bonding contract code
curl -LO https://github.com/CosmWasm/cosmwasm-plus/releases/download/v0.9.0/cw20_bonding.wasm

# upload contract code
acarbd tx wasm store cw20_bonding.wasm --from validator $TXFLAG

# instantiate contract
INIT='{
  "name": "sirbobo",
  "symbol": "BOBO",
  "decimals": 2,
  "reserve_denom": "uacarb",
  "reserve_decimals": 8,
  "curve_type": { "linear": { "slope": "1", "scale": 1 } }
}'
acarbd tx wasm instantiate 1 "$INIT" --from creator --label "social token" $TXFLAG

# get contract address
acarbd q wasm list-contract-by-code 1 --output json
CONTRACT=$(acarbd q wasm list-contract-by-code 1 --output json | jq -r '.contracts[-1]')

# query contract
acarbd q wasm contract-state smart $CONTRACT '{"token_info":{}}'
acarbd q wasm contract-state smart $CONTRACT '{"curve_info":{}}'
acarbd q wasm contract-state smart $CONTRACT "{\"balance\":{\"address\":\"$INVESTOR\"}}"

# execute a buy order
BUY='{"buy":{}}'
acarbd tx wasm execute $CONTRACT $BUY --from investor --amount=500000000uacarb $TXFLAG

# check balances
acarbd q bank balances $INVESTOR
acarbd q wasm contract-state smart $CONTRACT "{\"balance\":{\"address\":\"$INVESTOR\"}}"

# execute a burn / sell order
SELL='{"burn":{"amount":"500"}}'
acarbd tx wasm execute $CONTRACT $SELL --from investor $TXFLAG
acarbd q wasm contract-state smart $CONTRACT "{\"balance\":{\"address\":\"$INVESTOR\"}}"
acarbd q bank balances $INVESTOR

