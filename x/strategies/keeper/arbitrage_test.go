package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	mmtypes "github.com/kopi-money/kopi/x/mm/types"

	"cosmossdk.io/math"
	"github.com/kopi-money/kopi/cache"
	"github.com/kopi-money/kopi/constants"
	keepertest "github.com/kopi-money/kopi/testutil/keeper"
	dextypes "github.com/kopi-money/kopi/x/dex/types"
	"github.com/kopi-money/kopi/x/strategies/types"
	"github.com/stretchr/testify/require"
)

func TestHandle1(t *testing.T) {
	k, _, dexMsgServer, _, ctx := keepertest.SetupStrategiesMsgServer(t)

	dexKeeper := k.DexKeeper.(keepertest.LiquidityI)
	keepertest.TestAddLiquidity(ctx, dexKeeper, t, keepertest.Alice, constants.BaseCurrency, 100000)
	keepertest.TestAddLiquidity(ctx, dexKeeper, t, keepertest.Alice, constants.KUSD, 100000)
	keepertest.TestAddLiquidity(ctx, dexKeeper, t, keepertest.Alice, "uwusdc", 100000)

	_, err := keepertest.Sell(ctx, dexMsgServer, &dextypes.MsgSell{
		Creator:        keepertest.Bob,
		DenomGiving:    constants.KUSD,
		DenomReceiving: "uwusdc",
		Amount:         "10000",
	})
	require.NoError(t, err)

	parity1, _, err := k.DexKeeper.CalculateParity(ctx, constants.KUSD)
	require.NoError(t, err)
	require.NotNil(t, parity1)
	require.True(t, parity1.LT(math.LegacyOneDec()))

	require.NoError(t, k.HandleArbitrageDenoms(ctx))

	parity2, _, err := k.DexKeeper.CalculateParity(ctx, constants.KUSD)
	require.NoError(t, err)
	require.NotNil(t, parity2)
	require.True(t, parity1.Equal(*parity2))
}

func TestHandle2(t *testing.T) {
	k, msg, dexMsgServer, _, ctx := keepertest.SetupStrategiesMsgServer(t)
	moduleAcc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolArbitrage)
	liqPool := k.AccountKeeper.GetModuleAccount(ctx, dextypes.PoolLiquidity)

	dexKeeper := k.DexKeeper.(keepertest.LiquidityI)
	keepertest.TestAddLiquidity(ctx, dexKeeper, t, keepertest.Alice, constants.BaseCurrency, 1000000)
	keepertest.TestAddLiquidity(ctx, dexKeeper, t, keepertest.Alice, constants.KUSD, 1000000)
	keepertest.TestAddLiquidity(ctx, dexKeeper, t, keepertest.Alice, "uwusdc", 1000000)

	balance := k.BankKeeper.SpendableCoins(ctx, liqPool.GetAddress()).AmountOf(constants.KUSD)
	require.Equal(t, int64(1000000), balance.Int64())
	balance = k.BankKeeper.SpendableCoins(ctx, liqPool.GetAddress()).AmountOf("uwusdc")
	require.Equal(t, int64(1000000), balance.Int64())

	_, err := keepertest.Sell(ctx, dexMsgServer, &dextypes.MsgSell{
		Creator:        keepertest.Alice,
		DenomGiving:    constants.KUSD,
		DenomReceiving: "uwusdc",
		Amount:         "100000",
	})
	require.NoError(t, err)

	balance = k.BankKeeper.SpendableCoins(ctx, liqPool.GetAddress()).AmountOf(constants.KUSD)
	require.Equal(t, int64(1100000), balance.Int64())

	require.NoError(t, keepertest.AddArbitrageDeposit(ctx, msg, &types.MsgArbitrageDeposit{
		Creator: keepertest.Alice,
		Denom:   "uwusdc",
		Amount:  "10000",
	}))

	parity1, _, err := k.DexKeeper.CalculateParity(ctx, constants.KUSD)
	require.NoError(t, err)
	require.NotNil(t, parity1)
	require.True(t, parity1.LT(math.LegacyOneDec()))

	balance = k.BankKeeper.SpendableCoins(ctx, moduleAcc.GetAddress()).AmountOf("ucwusdc")
	require.Equal(t, int64(10000), balance.Int64())

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		return k.HandleArbitrageDenoms(innerCtx)
	}))

	balance = k.BankKeeper.SpendableCoins(ctx, moduleAcc.GetAddress()).AmountOf(constants.KUSD)
	require.True(t, balance.GT(math.ZeroInt()))

	parity2, _, err := k.DexKeeper.CalculateParity(ctx, constants.KUSD)
	require.NoError(t, err)
	require.NotNil(t, parity2)

	require.True(t, parity1.LT(*parity2))

	_, err = keepertest.Sell(ctx, dexMsgServer, &dextypes.MsgSell{
		Creator:        keepertest.Bob,
		DenomGiving:    "uwusdc",
		DenomReceiving: constants.KUSD,
		Amount:         "200000",
	})
	require.NoError(t, err)

	parity3, _, err := k.DexKeeper.CalculateParity(ctx, constants.KUSD)
	require.NoError(t, err)
	require.NotNil(t, parity3)
	require.True(t, parity3.GT(math.LegacyOneDec()))

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		return k.HandleArbitrageDenoms(innerCtx)
	}))

	parity4, _, err := k.DexKeeper.CalculateParity(ctx, constants.KUSD)
	require.NoError(t, err)
	require.NotNil(t, parity4)
	require.True(t, parity4.LT(*parity3))
}

func TestArbitrage1(t *testing.T) {
	_, msg, _, _, ctx := keepertest.SetupStrategiesMsgServer(t)

	require.Error(t, keepertest.AddArbitrageDeposit(ctx, msg, &types.MsgArbitrageDeposit{
		Creator: keepertest.Alice,
		Denom:   "uawusdc",
		Amount:  "1000",
	}))
}

func TestArbitrage2(t *testing.T) {
	k, msg, _, _, ctx := keepertest.SetupStrategiesMsgServer(t)
	moduleAcc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolArbitrage)

	require.NoError(t, keepertest.AddArbitrageDeposit(ctx, msg, &types.MsgArbitrageDeposit{
		Creator: keepertest.Alice,
		Denom:   "uwusdc",
		Amount:  "1000",
	}))

	supply := k.BankKeeper.GetSupply(ctx, "uawusdc")
	require.Equal(t, int64(1000), supply.Amount.Int64())

	balance := k.BankKeeper.SpendableCoins(ctx, moduleAcc.GetAddress()).AmountOf("ucwusdc")
	require.Equal(t, int64(1000), balance.Int64())

	require.NoError(t, keepertest.AddArbitrageDeposit(ctx, msg, &types.MsgArbitrageDeposit{
		Creator: keepertest.Alice,
		Denom:   "uwusdc",
		Amount:  "1000",
	}))

	supply = k.BankKeeper.GetSupply(ctx, "uawusdc")
	require.Equal(t, int64(2000), supply.Amount.Int64())

	balance = k.BankKeeper.SpendableCoins(ctx, moduleAcc.GetAddress()).AmountOf("ucwusdc")
	require.Equal(t, int64(2000), balance.Int64())
}

func TestArbitrage3(t *testing.T) {
	_, msg, _, _, ctx := keepertest.SetupStrategiesMsgServer(t)

	require.NoError(t, keepertest.AddArbitrageDeposit(ctx, msg, &types.MsgArbitrageDeposit{
		Creator: keepertest.Alice,
		Denom:   "uwusdc",
		Amount:  "1000",
	}))

	require.Error(t, keepertest.Redeem(ctx, msg, &types.MsgArbitrageRedeem{
		Creator: keepertest.Alice,
		Denom:   "uwusdc",
		Amount:  "1000",
	}))
}

func TestArbitrage4(t *testing.T) {
	k, msg, _, _, ctx := keepertest.SetupStrategiesMsgServer(t)
	moduleAcc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolArbitrage)

	require.NoError(t, keepertest.AddArbitrageDeposit(ctx, msg, &types.MsgArbitrageDeposit{
		Creator: keepertest.Alice,
		Denom:   "uwusdc",
		Amount:  "1000",
	}))

	balance := k.BankKeeper.SpendableCoins(ctx, moduleAcc.GetAddress()).AmountOf("ucwusdc")
	require.Equal(t, int64(1000), balance.Int64())

	require.NoError(t, keepertest.Redeem(ctx, msg, &types.MsgArbitrageRedeem{
		Creator: keepertest.Alice,
		Denom:   "uawusdc",
		Amount:  "1000",
	}))

	supply := k.BankKeeper.GetSupply(ctx, "uawusdc")
	require.Equal(t, int64(0), supply.Amount.Int64())

	require.NoError(t, keepertest.AddArbitrageDeposit(ctx, msg, &types.MsgArbitrageDeposit{
		Creator: keepertest.Alice,
		Denom:   "uwusdc",
		Amount:  "1000",
	}))

	balance = k.BankKeeper.SpendableCoins(ctx, moduleAcc.GetAddress()).AmountOf("ucwusdc")
	require.Equal(t, int64(1005), balance.Int64())
}

func TestArbitrage5(t *testing.T) {
	k, msg, _, mmMsg, ctx := keepertest.SetupStrategiesMsgServer(t)

	acc, _ := sdk.AccAddressFromBech32(keepertest.Alice)

	require.NoError(t, keepertest.AddDeposit(ctx, mmMsg, &mmtypes.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   "uwusdc",
		Amount:  "1000",
	}))

	balance1 := k.BankKeeper.SpendableCoins(ctx, acc).AmountOf("ucwusdc").Int64()

	require.NoError(t, keepertest.AddArbitrageDeposit(ctx, msg, &types.MsgArbitrageDeposit{
		Creator: keepertest.Alice,
		Denom:   "ucwusdc",
		Amount:  "1000",
	}))

	balance2 := k.BankKeeper.SpendableCoins(ctx, acc).AmountOf("ucwusdc").Int64()
	require.Less(t, balance2, balance1)
}
