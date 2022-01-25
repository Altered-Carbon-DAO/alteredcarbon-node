package app_test

import (
	"testing"

	"github.com/Altered-Carbon-DAO/alteredcarbon-node/v2/testutil/simapp"
)

func TestAnteHandler(t *testing.T) {
	simapp.New(t.TempDir())
	// suite.app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "alteredcarbon-1", Time: time.Now().UTC()})

}
