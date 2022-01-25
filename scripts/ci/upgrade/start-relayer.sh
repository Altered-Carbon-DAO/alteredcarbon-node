set -ex
mkdir -p ~/.hermes/
cp ./scripts/ci/upgrade/config.toml ~/.hermes/
hermes keys add alteredcarbon -f $PWD/scripts/ci/hermes/alteredcarbon.json
hermes keys add gaia -f $PWD/scripts/ci/hermes/gaia.json
hermes keys add osmosis -f $PWD/scripts/ci/hermes/osmosis.json
hermes start
