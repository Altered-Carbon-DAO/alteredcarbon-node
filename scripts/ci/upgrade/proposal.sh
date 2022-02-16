set -ex
set -o pipefail
DENOM=uacarb
CHAINID=alteredcarbon
RLYKEY=acarb12g0xe2ld0k5ws3h7lmxc39d4rpl3fyxp5qys69
acarbd version --long
apk add -U --no-cache jq tree
STARGAZE_HOME=/alteredcarbon/acarbd
acarbd config keyring-backend test --home $STARGAZE_HOME

HEIGHT=$(acarbd status --node http://alteredcarbon:26657 --home $ACARB_HOME | jq .SyncInfo.latest_block_height -r)
tree -L 2 /alteredcarbon/acarbd/
echo "current height $HEIGHT"
HEIGHT=$(expr $HEIGHT + 20) 
echo "submit with height $HEIGHT"
acarbd tx gov submit-proposal software-upgrade v2 --upgrade-height $HEIGHT  \
--deposit 10000000uacarb \
--description "Upgrade contains fix for claiming airdrop with Keplr and Ledger" \
--title "V2 Upgrade" \
--gas-prices 0.025uacarb --gas auto --gas-adjustment 1.5 --from validator  \
--chain-id alteredcarbon -b block --yes --node http://alteredcarbon:26657 --home $STARGAZE_HOME --keyring-backend test

acarbd q gov proposals --node http://alteredcarbon:26657 --home $STARGAZE_HOME


acarbd tx gov vote 1 "yes" --gas-prices 0.025uacarb --gas auto --gas-adjustment 1.5 --from validator  \
--chain-id alteredcarbon -b block --yes --node http://alteredcarbon:26657 --home $STARGAZE_HOME --keyring-backend test
sleep 60
acarbd q gov proposals --node http://alteredcarbon:26657 --home $STARGAZE_HOME
sleep 60