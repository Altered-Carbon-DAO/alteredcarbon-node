package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	minttypes "github.com/Altered-Carbon-DAO/alteredcarbon-node/x/mint/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"

	claimtypes "github.com/Altered-Carbon-DAO/alteredcarbon-node/x/claim/types"
)

const (
	HumanCoinUnit       = "acarb"
	BaseCoinUnit        = "uacarb"
	acarbExponent       = 6
	Bech32PrefixAccAddr = "acarb"
)

type Snapshot struct {
	TotalacarbAirdropAmount sdk.Int                    `json:"total_acarb_amount"`
	Accounts                map[string]SnapshotAccount `json:"accounts"`
}

type SnapshotAccount struct {
	AtomAddress              string  `json:"atom_address"`
	OsmoAddress              string  `json:"osmo_address"`
	RegenAddress             string  `json:"regen_address"`
	AtomStaker               bool    `json:"atom_staker"`
	OsmoStaker               bool    `json:"osmo_staker"`
	RegenStaker              bool    `json:"regen_staker"`
	OsmosisLiquidityProvider bool    `json:"osmosis_lp"`
	AirdropAmount            sdk.Int `json:"airdrop_amount"`
}

type GenesisParams struct {
	AirdropSupply sdk.Int

	StrategicReserveAccounts []banktypes.Balance

	ConsensusParams *tmproto.ConsensusParams

	GenesisTime         time.Time
	NativeCoinMetadatas []banktypes.Metadata

	StakingParams      stakingtypes.Params
	DistributionParams distributiontypes.Params
	GovParams          govtypes.Params

	CrisisConstantFee sdk.Coin

	SlashingParams slashingtypes.Params

	ClaimParams claimtypes.Params
	MintParams  minttypes.Params
}

func PrepareGenesisCmd(defaultNodeHome string, mbm module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepare-genesis [network] [chainID] [file]",
		Short: "Prepare a genesis file with initial setup",
		Long: `Prepare a genesis file with initial setup.
Examples include:
	- Setting module initial params
	- Setting denom metadata
Example:
	acarbd prepare-genesis mainnet alteredcarbon-1 snapshot.json
	- Check input genesis:
		file is at ~/.acarbd/config/genesis.json
`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			// read genesis file
			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			// get genesis params
			genesisParams := MainnetGenesisParams()
			switch args[0] {
			case "localnet":
				genesisParams = LocalnetGenesisParams()
			case "testnet":
				genesisParams = TestnetGenesisParams()
			case "devnet":
				genesisParams = DevnetGenesisParams()
			}
			// get genesis params
			chainID := args[1]

			// read snapshot.json and parse into struct
			snapshotFile, _ := ioutil.ReadFile(args[2])
			snapshot := Snapshot{}
			err = json.Unmarshal(snapshotFile, &snapshot)
			if err != nil {
				panic(err)
			}

			// run Prepare Genesis
			appState, genDoc, err = PrepareGenesis(clientCtx, appState, genDoc, genesisParams, chainID, snapshot)
			if err != nil {
				return fmt.Errorf("failed to prepare genesis: %w", err)
			}

			// validate genesis state
			if err = mbm.ValidateGenesis(cdc, clientCtx.TxConfig, appState); err != nil {
				return fmt.Errorf("error validating genesis file: %s", err.Error())
			}

			// save genesis
			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			err = genutil.ExportGenesisFile(genDoc, genFile)
			return err
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// fill with data
func PrepareGenesis(
	clientCtx client.Context,
	appState map[string]json.RawMessage,
	genDoc *tmtypes.GenesisDoc,
	genesisParams GenesisParams,
	chainID string,
	snapshot Snapshot,
) (map[string]json.RawMessage, *tmtypes.GenesisDoc, error) {
	cdc := clientCtx.Codec

	// chain params genesis
	genDoc.GenesisTime = genesisParams.GenesisTime
	genDoc.ChainID = chainID
	genDoc.ConsensusParams = genesisParams.ConsensusParams

	// IBC transfer module genesis
	ibcGenState := ibctransfertypes.DefaultGenesisState()
	ibcGenState.Params.SendEnabled = true
	ibcGenState.Params.ReceiveEnabled = true
	ibcGenStateBz, err := cdc.MarshalJSON(ibcGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal IBC transfer genesis state: %w", err)
	}
	appState[ibctransfertypes.ModuleName] = ibcGenStateBz

	// mint module genesis
	mintGenState := minttypes.DefaultGenesisState()
	mintGenState.Params = genesisParams.MintParams

	mintGenStateBz, err := cdc.MarshalJSON(mintGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal mint genesis state: %w", err)
	}
	appState[minttypes.ModuleName] = mintGenStateBz

	// staking module genesis
	stakingGenState := stakingtypes.GetGenesisStateFromAppState(cdc, appState)
	stakingGenState.Params = genesisParams.StakingParams
	stakingGenStateBz, err := cdc.MarshalJSON(stakingGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal staking genesis state: %w", err)
	}
	appState[stakingtypes.ModuleName] = stakingGenStateBz

	// distribution module genesis
	distributionGenState := distributiontypes.DefaultGenesisState()
	distributionGenState.Params = genesisParams.DistributionParams
	distributionGenStateBz, err := cdc.MarshalJSON(distributionGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal distribution genesis state: %w", err)
	}
	appState[distributiontypes.ModuleName] = distributionGenStateBz

	// gov module genesis
	govGenState := govtypes.DefaultGenesisState()
	govGenState.DepositParams = genesisParams.GovParams.DepositParams
	govGenState.TallyParams = genesisParams.GovParams.TallyParams
	govGenState.VotingParams = genesisParams.GovParams.VotingParams
	govGenStateBz, err := cdc.MarshalJSON(govGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal gov genesis state: %w", err)
	}
	appState[govtypes.ModuleName] = govGenStateBz

	// crisis module genesis
	crisisGenState := crisistypes.DefaultGenesisState()
	crisisGenState.ConstantFee = genesisParams.CrisisConstantFee
	crisisGenStateBz, err := cdc.MarshalJSON(crisisGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal crisis genesis state: %w", err)
	}
	appState[crisistypes.ModuleName] = crisisGenStateBz

	// slashing module genesis
	slashingGenState := slashingtypes.DefaultGenesisState()
	slashingGenState.Params = genesisParams.SlashingParams
	slashingGenStateBz, err := cdc.MarshalJSON(slashingGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal slashing genesis state: %w", err)
	}
	appState[slashingtypes.ModuleName] = slashingGenStateBz

	// auth accounts
	authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)
	accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get accounts from any: %w", err)
	}

	// ---
	// bank module genesis
	bankGenState := banktypes.GetGenesisStateFromAppState(cdc, appState)
	bankGenState.Params.DefaultSendEnabled = true
	bankGenState.DenomMetadata = genesisParams.NativeCoinMetadatas
	balances := bankGenState.Balances

	// claim module genesis
	claimGenState := claimtypes.GetGenesisStateFromAppState(cdc, appState)
	claimGenState.Params = genesisParams.ClaimParams
	claimRecords := make([]claimtypes.ClaimRecord, 0, len(snapshot.Accounts))
	claimsTotal := sdk.ZeroInt()
	// check from preexisint accounts in genesis
	preExistingAccounts := make(map[string]bool)
	for _, b := range balances {
		preExistingAccounts[b.Address] = true
	}
	for addr, acc := range snapshot.Accounts {
		claimRecord := claimtypes.ClaimRecord{
			Address:                addr,
			InitialClaimableAmount: sdk.NewCoins(sdk.NewCoin(BaseCoinUnit, acc.AirdropAmount)),
			ActionCompleted:        []bool{false, false, false, false, false},
		}
		claimsTotal = claimsTotal.Add(acc.AirdropAmount)
		claimRecords = append(claimRecords, claimRecord)
		// skip account addition if existent
		exists := preExistingAccounts[addr]
		if exists {
			continue
		}
		balances = append(balances, banktypes.Balance{
			Address: addr,
			Coins:   sdk.NewCoins(sdk.NewInt64Coin(BaseCoinUnit, 1)),
		})

		address, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return nil, nil, err
		}
		// add base account
		// Add the new account to the set of genesis accounts
		baseAccount := authtypes.NewBaseAccount(address, nil, 0, 0)
		if err := baseAccount.Validate(); err != nil {
			return nil, nil, fmt.Errorf("failed to validate new genesis account: %w", err)
		}
		accs = append(accs, baseAccount)
	}
	claimGenState.ClaimRecords = claimRecords
	claimGenState.ModuleAccountBalance = sdk.NewCoin(BaseCoinUnit, claimsTotal)
	claimGenStateBz, err := cdc.MarshalJSON(claimGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal claim genesis state: %w", err)
	}
	appState[claimtypes.ModuleName] = claimGenStateBz

	// save accounts

	// auth module genesis
	accs = authtypes.SanitizeGenesisAccounts(accs)
	genAccs, err := authtypes.PackAccounts(accs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert accounts into any's: %w", err)
	}
	authGenState.Accounts = genAccs
	authGenStateBz, err := cdc.MarshalJSON(&authGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal auth genesis state: %w", err)
	}
	appState[authtypes.ModuleName] = authGenStateBz

	// save balances
	bankGenState.Balances = banktypes.SanitizeGenesisBalances(balances)
	bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}
	appState[banktypes.ModuleName] = bankGenStateBz

	return appState, genDoc, nil
}

// params only
func MainnetGenesisParams() GenesisParams {
	genParams := GenesisParams{}

	genParams.AirdropSupply = sdk.NewInt(250_000_000_000_000)              // 250M acarb
	genParams.GenesisTime = time.Date(2022, 01, 02, 17, 0, 0, 0, time.UTC) //

	genParams.NativeCoinMetadatas = []banktypes.Metadata{
		{
			Description: "The native token of Altered Carbon",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    BaseCoinUnit,
					Exponent: 0,
					Aliases:  nil,
				},
				{
					Denom:    HumanCoinUnit,
					Exponent: acarbExponent,
					Aliases:  nil,
				},
			},
			Name:    "acarb",
			Base:    BaseCoinUnit,
			Display: HumanCoinUnit,
			Symbol:  "acarb",
		},
	}

	// mint
	genParams.MintParams = minttypes.DefaultParams()
	genParams.MintParams.MintDenom = BaseCoinUnit
	genParams.MintParams.StartTime = genParams.GenesisTime.AddDate(1, 0, 0)
	genParams.MintParams.InitialAnnualProvisions = sdk.NewDec(1_000_000_000_000_000)
	genParams.MintParams.ReductionFactor = sdk.NewDec(2).QuoInt64(3)
	genParams.MintParams.BlocksPerYear = uint64(5737588)

	genParams.StakingParams = stakingtypes.DefaultParams()
	genParams.StakingParams.UnbondingTime = time.Hour * 24 * 7 * 2 // 2 weeks
	genParams.StakingParams.MaxValidators = 100
	genParams.StakingParams.BondDenom = genParams.NativeCoinMetadatas[0].Base
	// MinCommissionRate is enforced in ante-handler

	genParams.DistributionParams = distributiontypes.DefaultParams()

	genParams.DistributionParams.BaseProposerReward = sdk.MustNewDecFromStr("0.01")
	genParams.DistributionParams.BonusProposerReward = sdk.MustNewDecFromStr("0.04")
	genParams.DistributionParams.CommunityTax = sdk.MustNewDecFromStr("0.05")
	genParams.DistributionParams.WithdrawAddrEnabled = true

	genParams.GovParams = govtypes.DefaultParams()
	genParams.GovParams.DepositParams.MaxDepositPeriod = time.Hour * 24 * 14 // 2 weeks
	genParams.GovParams.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(1_000_000_000),
	))
	genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.2") // 40%
	genParams.GovParams.VotingParams.VotingPeriod = time.Hour * 24 * 3    // 3 days

	genParams.CrisisConstantFee = sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(100_000_000_000),
	)

	genParams.SlashingParams = slashingtypes.DefaultParams()
	genParams.SlashingParams.SignedBlocksWindow = int64(25000)                       // ~41 hr at 6 second blocks
	genParams.SlashingParams.MinSignedPerWindow = sdk.MustNewDecFromStr("0.05")      // 5% minimum liveness
	genParams.SlashingParams.DowntimeJailDuration = time.Minute                      // 1 minute jail period
	genParams.SlashingParams.SlashFractionDoubleSign = sdk.MustNewDecFromStr("0.05") // 5% double sign slashing
	genParams.SlashingParams.SlashFractionDowntime = sdk.MustNewDecFromStr("0.0001") // 0.01% liveness slashing

	genParams.ClaimParams = claimtypes.Params{
		AirdropEnabled:     false,
		AirdropStartTime:   genParams.GenesisTime.Add(time.Hour * 24 * 365), // 1 year (will be changed by gov)
		DurationUntilDecay: time.Hour * 24 * 120,                            // 120 days = ~4 months
		DurationOfDecay:    time.Hour * 24 * 120,                            // 120 days = ~4 months
		ClaimDenom:         genParams.NativeCoinMetadatas[0].Base,
	}

	genParams.ConsensusParams = tmtypes.DefaultConsensusParams()
	genParams.ConsensusParams.Block.MaxBytes = 25 * 1024 * 1024 // 26,214,400 for cosmwasm
	genParams.ConsensusParams.Block.MaxGas = 10_000_000
	genParams.ConsensusParams.Evidence.MaxAgeDuration = genParams.StakingParams.UnbondingTime
	genParams.ConsensusParams.Evidence.MaxAgeNumBlocks = int64(genParams.StakingParams.UnbondingTime.Seconds()) / 3
	genParams.ConsensusParams.Version.AppVersion = 1

	return genParams
}

// params only
func TestnetGenesisParams() GenesisParams {
	genParams := MainnetGenesisParams()

	genParams.AirdropSupply = sdk.NewInt(250_000_000_000_000) // 250M acarb
	genParams.GenesisTime = time.Now()

	// mint
	genParams.MintParams.StartTime = genParams.GenesisTime.Add(time.Minute * 5)

	genParams.GovParams.DepositParams.MaxDepositPeriod = time.Hour * 24 * 14 // 2 weeks
	genParams.GovParams.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(1),
	))
	genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.2") // 20%
	genParams.GovParams.VotingParams.VotingPeriod = time.Minute * 15      // 15 min

	return genParams
}

func DevnetGenesisParams() GenesisParams {
	genParams := MainnetGenesisParams()

	genParams.AirdropSupply = sdk.NewInt(250_000_000_000_000) // 250M acarb
	genParams.GenesisTime = time.Now()
	genParams.ClaimParams.AirdropEnabled = true
	genParams.ClaimParams.AirdropStartTime = genParams.GenesisTime
	// mint
	genParams.MintParams.StartTime = genParams.GenesisTime.Add(time.Hour * 10)

	genParams.GovParams.DepositParams.MaxDepositPeriod = time.Hour * 1 // 1 hour
	genParams.GovParams.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(1),
	))
	genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.1")    // 10%
	genParams.GovParams.TallyParams.Threshold = sdk.MustNewDecFromStr("0.5") // 50%
	genParams.GovParams.VotingParams.VotingPeriod = time.Minute * 5          // 5 min

	return genParams
}

func LocalnetGenesisParams() GenesisParams {
	genParams := MainnetGenesisParams()

	genParams.AirdropSupply = sdk.NewInt(250_000_000_000_000) // 250M acarb
	genParams.GenesisTime = time.Now()
	genParams.ClaimParams.AirdropEnabled = true
	genParams.ClaimParams.AirdropStartTime = genParams.GenesisTime
	// mint
	genParams.MintParams.StartTime = genParams.GenesisTime.Add(time.Hour * 10)

	genParams.GovParams.DepositParams.MaxDepositPeriod = time.Hour * 1 // 1 hour
	genParams.GovParams.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(1),
	))
	genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.1")    // 10%
	genParams.GovParams.TallyParams.Threshold = sdk.MustNewDecFromStr("0.5") // 50%
	genParams.GovParams.VotingParams.VotingPeriod = time.Minute * 1          // 5 min

	return genParams
}
