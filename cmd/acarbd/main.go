package main

import (
	"os"

	"github.com/Altered-Carbon-DAO/alteredcarbon-node/app"
	"github.com/Altered-Carbon-DAO/alteredcarbon-node/cmd/acarbd/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/tendermint/starport/starport/pkg/cosmoscmd"
	tmcmds "github.com/tendermint/tendermint/cmd/tendermint/commands"
)

func main() {
	rootCmd, _ := cosmoscmd.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		app.Name,
		app.ModuleBasics,
		app.NewAlteredCarbonApp,
		cosmoscmd.AddSubCmd(cmd.TestnetCmd(app.ModuleBasics)),
		cosmoscmd.AddSubCmd(cmd.PrepareGenesisCmd(app.DefaultNodeHome, app.ModuleBasics)),
		cosmoscmd.AddSubCmd(tmcmds.RollbackStateCmd),
		// this line is used by starport scaffolding # root/arguments
	)
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
