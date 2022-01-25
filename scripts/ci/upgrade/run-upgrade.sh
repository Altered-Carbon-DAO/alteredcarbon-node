set -ex
set -o pipefail
DENOM=uacarb
CHAINID=alteredcarbon
RLYKEY=acarb12g0xe2ld0k5ws3h7lmxc39d4rpl3fyxp5qys69
acarbd version --long
apk add -U --no-cache jq tree curl wget
STARGAZE_HOME=/alteredcarbon/acarbd
curl -s -v http://alteredcarbon:8090/kill || echo "done"
sleep 10
acarbd start --pruning nothing --home $STARGAZE_HOME
