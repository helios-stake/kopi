package keeper

import (
	"context"
	"testing"

	"github.com/kopi-money/kopi/constants"

	tokenfactorytypes "github.com/kopi-money/kopi/x/tokenfactory/types"

	"cosmossdk.io/math"
	"github.com/kopi-money/kopi/cache"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	denomkeeper "github.com/kopi-money/kopi/x/denominations/keeper"
	denomtypes "github.com/kopi-money/kopi/x/denominations/types"
	dextypes "github.com/kopi-money/kopi/x/dex/types"
	mmtypes "github.com/kopi-money/kopi/x/mm/types"
	reservetypes "github.com/kopi-money/kopi/x/reserve/types"
	strategiestypes "github.com/kopi-money/kopi/x/strategies/types"
	swaptypes "github.com/kopi-money/kopi/x/swap/types"
	"github.com/stretchr/testify/require"
)

type Keys struct {
	cdc      *codec.ProtoCodec
	registry codectypes.InterfaceRegistry

	acc *storetypes.KVStoreKey
	dex *storetypes.KVStoreKey
	bnk *storetypes.KVStoreKey
	dnm *storetypes.KVStoreKey
	mm  *storetypes.KVStoreKey
	res *storetypes.KVStoreKey
	swp *storetypes.KVStoreKey
	str *storetypes.KVStoreKey
	tof *storetypes.KVStoreKey
}

func DenomKeeper(t *testing.T) (denomkeeper.Keeper, context.Context, *Keys) {
	initSDKConfig()
	cache.NewTranscationHandler()

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	keys := Keys{
		acc: storetypes.NewKVStoreKey(authtypes.StoreKey),
		bnk: storetypes.NewKVStoreKey(banktypes.StoreKey),
		dex: storetypes.NewKVStoreKey(dextypes.StoreKey),
		dnm: storetypes.NewKVStoreKey(denomtypes.StoreKey),
		mm:  storetypes.NewKVStoreKey(mmtypes.StoreKey),
		res: storetypes.NewKVStoreKey(reservetypes.StoreKey),
		swp: storetypes.NewKVStoreKey(swaptypes.StoreKey),
		str: storetypes.NewKVStoreKey(strategiestypes.StoreKey),
		tof: storetypes.NewKVStoreKey(tokenfactorytypes.StoreKey),

		cdc:      cdc,
		registry: registry,
	}

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(keys.acc, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keys.bnk, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keys.dex, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keys.dnm, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keys.mm, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keys.res, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keys.swp, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keys.str, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keys.tof, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())
	require.NoError(t, stateStore.LoadLatestVersion())

	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	denomKeeper := denomkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keys.dnm),
		log.NewNopLogger(),
		authority.String(),
	)
	cache.AddCache(denomKeeper)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())
	params := createDenomTestParams()
	params.DexDenoms = append(params.DexDenoms,
		&denomtypes.DexDenom{
			Name:         "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
			MinLiquidity: math.NewInt(100_000),
			MinOrderSize: math.NewInt(1_000_000),
			Exponent:     6,
		},
		&denomtypes.DexDenom{
			Name:         "uawusdc",
			MinLiquidity: math.NewInt(1000),
			MinOrderSize: math.NewInt(1000),
			Exponent:     6,
		},
		&denomtypes.DexDenom{
			Name:         "inj",
			MinLiquidity: math.NewInt(1000),
			MinOrderSize: math.NewInt(1000),
			Exponent:     18,
		},
	)

	params.StrategyDenoms = &denomtypes.StrategyDenoms{
		ArbitrageDenoms: []*denomtypes.ArbitrageDenom{
			{
				DexDenom:                  "uawusdc",
				KCoin:                     "ukusd",
				CAsset:                    "ucwusdc",
				BuyThreshold:              math.LegacyOneDec(),
				SellThreshold:             math.LegacyOneDec(),
				BuyTradeAmount:            math.NewInt(2000),
				SellTradeAmount:           math.NewInt(2000),
				RedemptionFee:             math.LegacyNewDecWithPrec(1, 2),
				RedemptionFeeReserveShare: math.LegacyNewDecWithPrec(5, 1),
			},
		},
	}

	require.NoError(t, cache.Transact(ctx, func(innerContext context.Context) error {
		denomKeeper.SetRatio(innerContext, denomtypes.Ratio{
			Denom: "uwusdc",
			Ratio: math.LegacyNewDecWithPrec(25, 2),
		})

		denomKeeper.SetRatio(innerContext, denomtypes.Ratio{
			Denom: "uwusdt",
			Ratio: math.LegacyNewDecWithPrec(25, 2),
		})

		denomKeeper.SetRatio(innerContext, denomtypes.Ratio{
			Denom: "ukusd",
			Ratio: math.LegacyNewDecWithPrec(25, 2),
		})

		denomKeeper.SetRatio(innerContext, denomtypes.Ratio{
			Denom: "uckusd",
			Ratio: math.LegacyNewDecWithPrec(25, 2),
		})

		denomKeeper.SetRatio(innerContext, denomtypes.Ratio{
			Denom: "ucwusdc",
			Ratio: math.LegacyNewDecWithPrec(25, 2),
		})

		denomKeeper.SetRatio(innerContext, denomtypes.Ratio{
			Denom: "swbtc",
			Ratio: math.LegacyNewDecWithPrec(1, 3),
		})

		denomKeeper.SetRatio(innerContext, denomtypes.Ratio{
			Denom: "skbtc",
			Ratio: math.LegacyNewDecWithPrec(1, 3),
		})

		denomKeeper.SetRatio(innerContext, denomtypes.Ratio{
			Denom: "sckbtc",
			Ratio: math.LegacyNewDecWithPrec(1, 3),
		})

		denomKeeper.SetRatio(innerContext, denomtypes.Ratio{
			Denom: "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
			Ratio: math.LegacyNewDec(10),
		})

		denomKeeper.SetRatio(innerContext, denomtypes.Ratio{
			Denom: "uawusdc",
			Ratio: math.LegacyNewDecWithPrec(25, 2),
		})

		injRatio, err := math.LegacyNewDecFromStr("11310893732.791635615102371449")
		require.NoError(t, err)
		denomKeeper.SetRatio(innerContext, denomtypes.Ratio{
			Denom: "inj",
			Ratio: injRatio,
		})

		return denomKeeper.SetParams(innerContext, params)
	}))

	return denomKeeper, ctx, &keys
}

func createDenomTestParams() denomtypes.Params {
	return denomtypes.Params{
		CAssets:          createDefaultCAssets(),
		CollateralDenoms: createDefaultCollateralDenoms(),
		DexDenoms:        createDefaultDexDenoms(),
		KCoins:           createDefaultKCoins(),
	}
}

func createDefaultCollateralDenoms() []*denomtypes.CollateralDenom {
	return []*denomtypes.CollateralDenom{
		{
			DexDenom:   constants.BaseCurrency,
			Ltv:        math.LegacyNewDecWithPrec(5, 1),
			MaxDeposit: math.NewInt(1_000_000_000),
		},
		{
			DexDenom:   "uwusdc",
			Ltv:        math.LegacyNewDecWithPrec(9, 1),
			MaxDeposit: math.NewInt(1_000_000_000),
		},
		{
			DexDenom:   "ucwusdc",
			Ltv:        math.LegacyNewDecWithPrec(95, 2),
			MaxDeposit: math.NewInt(1_000_000_000),
		},
		{
			DexDenom:   constants.KUSD,
			Ltv:        math.LegacyNewDecWithPrec(9, 1),
			MaxDeposit: math.NewInt(1_000_000_000),
		},
		{
			DexDenom:   "uckusd",
			Ltv:        math.LegacyNewDecWithPrec(95, 2),
			MaxDeposit: math.NewInt(1_000_000_000),
		},
		{
			DexDenom:   "swbtc",
			Ltv:        math.LegacyNewDecWithPrec(8, 1),
			MaxDeposit: math.NewInt(1_000_000_000),
		},
		{
			DexDenom:   "skbtc",
			Ltv:        math.LegacyNewDecWithPrec(8, 1),
			MaxDeposit: math.NewInt(1_000_000_000),
		},
	}
}

func createDefaultCAssets() []*denomtypes.CAsset {
	return []*denomtypes.CAsset{
		{
			DexDenom:        "uckusd",
			BaseDexDenom:    constants.KUSD,
			DexFeeShare:     math.LegacyNewDecWithPrec(5, 1),
			BorrowLimit:     math.LegacyNewDecWithPrec(99, 2),
			MinimumLoanSize: math.NewInt(1000),
		},
		{
			DexDenom:        "ucwusdc",
			BaseDexDenom:    "uwusdc",
			DexFeeShare:     math.LegacyNewDecWithPrec(5, 1),
			BorrowLimit:     math.LegacyNewDecWithPrec(99, 2),
			MinimumLoanSize: math.NewInt(1000),
		},
		{
			DexDenom:        "sckbtc",
			BaseDexDenom:    "skbtc",
			DexFeeShare:     math.LegacyNewDecWithPrec(5, 1),
			BorrowLimit:     math.LegacyNewDecWithPrec(99, 2),
			MinimumLoanSize: math.NewInt(1000),
		},
	}
}

func createDefaultDexDenoms() []*denomtypes.DexDenom {
	return []*denomtypes.DexDenom{
		{
			Name:         constants.BaseCurrency,
			MinLiquidity: math.NewInt(10_000),
			MinOrderSize: math.NewInt(1),
			Exponent:     6,
		},
		{
			Name:         "uwusdc",
			MinLiquidity: math.NewInt(10_000_000),
			MinOrderSize: math.NewInt(1),
			Exponent:     6,
		},
		{
			Name:         "uwusdt",
			MinLiquidity: math.NewInt(10_000_000),
			MinOrderSize: math.NewInt(1_000_000),
			Exponent:     6,
		},
		{
			Name:         constants.KUSD,
			MinLiquidity: math.NewInt(10_000_000),
			MinOrderSize: math.NewInt(1),
			Exponent:     6,
		},
		{
			Name:         "uckusd",
			MinLiquidity: math.NewInt(10_000_000),
			MinOrderSize: math.NewInt(1_000_000),
			Exponent:     6,
		},
		{
			Name:         "ucwusdc",
			MinLiquidity: math.NewInt(10_000_000),
			MinOrderSize: math.NewInt(1_000_000),
			Exponent:     6,
		},
		{
			Name:         "swbtc",
			MinLiquidity: math.NewInt(1_000),
			MinOrderSize: math.NewInt(1_000_000),
			Exponent:     8,
		},
		{
			Name:         "skbtc",
			MinLiquidity: math.NewInt(1_000),
			MinOrderSize: math.NewInt(1_000_000),
			Exponent:     8,
		},
		{
			Name:         "sckbtc",
			MinLiquidity: math.NewInt(1_000),
			MinOrderSize: math.NewInt(1_000_000),
			Exponent:     8,
		},
	}
}

func AddDexDenom(ctx context.Context, k denomtypes.MsgServer, msg *denomtypes.MsgDexAddDenom) error {
	_, err := k.DexAddDenom(ctx, msg)
	return err
}

func createDefaultKCoins() []*denomtypes.KCoin {
	return []*denomtypes.KCoin{
		{
			DexDenom:      constants.KUSD,
			References:    []string{"uwusdc", "uwusdt"},
			MaxSupply:     math.NewInt(1_000_000_000_000),
			MaxMintAmount: math.NewInt(1_000_000),
			MaxBurnAmount: math.NewInt(1_000_000),
		},
		{
			DexDenom:      "skbtc",
			References:    []string{"swbtc"},
			MaxSupply:     math.NewInt(100_000_000),
			MaxMintAmount: math.NewInt(10_000),
			MaxBurnAmount: math.NewInt(10_000),
		},
	}
}
