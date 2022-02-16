package keeper_test

import (
	"encoding/json"

	"github.com/tendermint/starport/starport/pkg/cosmoscmd"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	alteredcarbonapp "github.com/Altered-Carbon-DAO/alteredcarbon-node/app"
	"github.com/Altered-Carbon-DAO/alteredcarbon-node/x/mint/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool) (*alteredcarbonapp.App, sdk.Context) {
	app := setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, types.DefaultParams())
	app.MintKeeper.SetMinter(ctx, types.DefaultInitialMinter())

	return app, ctx
}

func setup(isCheckTx bool) *alteredcarbonapp.App {
	app, genesisState := genApp(!isCheckTx, 5)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func genApp(withGenesis bool, invCheckPeriod uint) (*alteredcarbonapp.App, alteredcarbonapp.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := cosmoscmd.MakeEncodingConfig(alteredcarbonapp.ModuleBasics)
	app := alteredcarbonapp.NewAlteredCarbonApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		simapp.DefaultNodeHome,
		invCheckPeriod,
		encCdc,
		simapp.EmptyAppOptions{})

	originalApp := app.(*alteredcarbonapp.App)

	if withGenesis {
		return originalApp, alteredcarbonapp.NewDefaultGenesisState(encCdc.Marshaler)
	}

	return originalApp, alteredcarbonapp.GenesisState{}
}
