package keeper_test

import (
	"context"
	"fmt"
	denomtypes "github.com/kopi-money/kopi/x/denominations/types"
	"testing"

	"github.com/kopi-money/kopi/x/dex/constant_product"

	"github.com/kopi-money/kopi/cache"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/kopi-money/kopi/constants"
	keepertest "github.com/kopi-money/kopi/testutil/keeper"
	dexkeeper "github.com/kopi-money/kopi/x/dex/keeper"
	"github.com/kopi-money/kopi/x/dex/types"
	"github.com/stretchr/testify/require"
)

func TestCalculateSingleMaximumTradableAmount1(t *testing.T) {
	actualFrom := math.LegacyNewDec(1000)
	virtualFrom := math.LegacyNewDec(0)

	actualTo := math.LegacyNewDec(1000)
	virtualTo := math.LegacyNewDec(0)

	maximum := dexkeeper.CalculateSingleMaximumGivableAmount(actualFrom, actualTo, virtualFrom, virtualTo, nil)
	require.Nil(t, maximum)
}

func TestCalculateSingleMaximumTradableAmount2(t *testing.T) {
	actualFrom := math.LegacyNewDec(1000)
	virtualFrom := math.LegacyNewDec(0)

	actualTo := math.LegacyNewDec(500)
	virtualTo := math.LegacyNewDec(500)

	maximum := dexkeeper.CalculateSingleMaximumGivableAmount(actualFrom, actualTo, virtualFrom, virtualTo, nil)
	require.NotNil(t, maximum)

	receive, _ := constant_product.ConstantProductTradeSell(actualFrom.Add(virtualFrom), actualTo.Add(virtualTo), *maximum, math.LegacyZeroDec())
	require.Equal(t, receive, actualTo)
}

func TestCalculateSingleMaximumTradableAmount3(t *testing.T) {
	actualFrom := math.LegacyNewDec(1000)
	virtualFrom := math.LegacyNewDec(0)

	actualTo := math.LegacyNewDec(100)
	virtualTo := math.LegacyNewDec(900)

	maximum := dexkeeper.CalculateSingleMaximumGivableAmount(actualFrom, actualTo, virtualFrom, virtualTo, nil)
	require.NotNil(t, maximum)

	receive, _ := constant_product.ConstantProductTradeSell(actualFrom.Add(virtualFrom), actualTo.Add(virtualTo), *maximum, math.LegacyZeroDec())

	// Due to rounding we don't get exactly 100, but 99.999999999999999910
	diff := actualTo.Sub(receive).Abs()
	require.True(t, diff.LT(math.LegacyNewDecWithPrec(1, 10)))
}

func TestCalculateSingleMaximumTradableAmount4(t *testing.T) {
	liqFrom := math.LegacyNewDec(4966641376348)
	liqTo, _ := math.LegacyNewDecFromStr("1187591536070.805216643324621576")

	maxPrice, _ := math.LegacyNewDecFromStr("4.18418")
	//fee, _ := math.LegacyNewDecFromStr("0.0009")
	fee := math.LegacyZeroDec()

	fmt.Println(liqFrom.Quo(liqTo).String())

	maximumGiving := constant_product.CalculateMaximumGiving(liqFrom, liqTo, maxPrice, fee)
	receiving, _ := constant_product.ConstantProductTradeSell(liqFrom, liqTo, maximumGiving, fee)
	price := maximumGiving.Quo(receiving)

	require.True(t, maxPrice.Equal(price))
	require.Greater(t, maximumGiving.TruncateInt64(), int64(0))
}

func TestSingleTrade1(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 100))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 100))

	require.Equal(t, int64(100), k.GetLiquiditySum(ctx, constants.BaseCurrency).Int64())
	require.Equal(t, int64(100), k.GetLiquiditySum(ctx, constants.KUSD).Int64())
	require.Equal(t, int64(400), k.GetFullLiquidityBase(ctx, constants.KUSD).TruncateInt().Int64())
	require.Equal(t, int64(100), k.GetFullLiquidityOther(ctx, constants.KUSD).TruncateInt().Int64())

	dexAcc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolLiquidity).GetAddress()
	coins := k.BankKeeper.SpendableCoins(ctx, dexAcc)
	require.Equal(t, int64(100), coins.AmountOf(constants.BaseCurrency).Int64())

	fee := math.LegacyZeroDec()

	ratio1, err := k.DenomKeeper.GetRatio(ctx, constants.KUSD)
	require.NoError(t, err)
	require.Equal(t, math.LegacyNewDecWithPrec(25, 2), ratio1.Ratio)

	pair1, _ := k.GetLiquidityPair(ctx, constants.KUSD)
	require.Equal(t, math.LegacyNewDec(300), pair1.VirtualBase)
	require.Equal(t, math.LegacyZeroDec(), pair1.VirtualOther)

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         math.NewInt(100),
		TradeDenomGiving:    constants.KUSD,
		TradeDenomReceiving: constants.BaseCurrency,
		CoinSource:          keepertest.Bob,
		CoinTarget:          keepertest.Bob,
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		FlatPrice:           &constant_product.FlatPrice{},
		Fee:                 fee,
	}

	var (
		amountUsed     math.Int
		amountReceived math.Int
	)

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		stepCtx := tradeCtx.TradeStep1(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), types.TradeTypeSell)
		amountUsed, amountReceived, _, err = k.ExecuteTradeStep(stepCtx)
		return err
	}))

	require.Equal(t, int64(100), amountUsed.Int64())
	require.Equal(t, int64(100), amountReceived.Int64())

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		stepCtx := tradeCtx.TradeStep2(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), amountUsed, types.TradeTypeSell)
		amountUsed, amountReceived, _, err = k.ExecuteTradeStep(stepCtx)
		return err
	}))

	require.Equal(t, int64(100), amountUsed.Int64())
	require.Equal(t, int64(100), amountReceived.Int64())
	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	liquidityPoolAcc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolLiquidity)
	liquidityPool := k.BankKeeper.SpendableCoins(ctx, liquidityPoolAcc.GetAddress())
	require.Equal(t, int64(0), liquidityPool.AmountOf(constants.BaseCurrency).Int64())

	//ratio2, err := k.GetRatio(ctx, constants.KUSD)
	//require.NoError(t, err)
	//require.Equal(t, "0.666666666666666667", ratio2.Ratio.String())

	require.Equal(t, int64(0), k.GetLiquiditySum(ctx, constants.BaseCurrency).Int64())
	require.Equal(t, int64(200), k.GetLiquiditySum(ctx, constants.KUSD).Int64())
	require.Equal(t, int64(300), k.GetFullLiquidityBase(ctx, constants.KUSD).RoundInt().Int64())
	require.Equal(t, int64(200), k.GetFullLiquidityOther(ctx, constants.KUSD).TruncateInt().Int64())

	pair2, _ := k.GetLiquidityPair(ctx, constants.KUSD)
	require.Equal(t, int64(300), pair2.VirtualBase.RoundInt().Int64())
	require.Equal(t, math.LegacyZeroDec(), pair2.VirtualOther)

	coins = k.BankKeeper.SpendableCoins(ctx, dexAcc)
	require.Equal(t, int64(0), coins.AmountOf(constants.BaseCurrency).Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestSingleTrade1a(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 100))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 100))

	fee := math.LegacyZeroDec()

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         math.NewInt(100),
		TradeDenomGiving:    constants.KUSD,
		TradeDenomReceiving: constants.BaseCurrency,
		CoinSource:          keepertest.Bob,
		CoinTarget:          keepertest.Bob,
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		FlatPrice:           &constant_product.FlatPrice{},
		Fee:                 fee,
	}

	var (
		amountUsed     math.Int
		amountReceived math.Int
		err            error
	)

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		stepCtx := tradeCtx.TradeStep1(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), types.TradeTypeBuy)
		amountUsed, amountReceived, _, err = k.ExecuteTradeStep(stepCtx)
		return err
	}))

	require.Equal(t, int64(100), amountUsed.Int64())
	require.Equal(t, int64(100), amountReceived.Int64())

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		stepCtx := tradeCtx.TradeStep2(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), amountUsed, types.TradeTypeBuy)
		amountUsed, amountReceived, _, err = k.ExecuteTradeStep(stepCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.Equal(t, int64(100), amountUsed.Int64())

	require.Equal(t, int64(0), k.GetLiquiditySum(ctx, constants.BaseCurrency).Int64())
	require.Equal(t, int64(200), k.GetLiquiditySum(ctx, constants.KUSD).Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestSingleTrade2(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 100))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 100))

	dexAcc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolLiquidity).GetAddress()
	coins := k.BankKeeper.SpendableCoins(ctx, dexAcc)
	require.Equal(t, int64(100), coins.AmountOf(constants.BaseCurrency).Int64())

	offer := math.NewInt(100)
	fee := math.LegacyZeroDec()

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         offer,
		CoinSource:          keepertest.Bob,
		CoinTarget:          keepertest.Bob,
		TradeDenomGiving:    constants.BaseCurrency,
		TradeDenomReceiving: constants.KUSD,
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		FlatPrice:           &constant_product.FlatPrice{},
		Fee:                 fee,
	}

	var (
		amount         math.Int
		usedAmount     math.Int
		receivedAmount math.Int
		err            error
	)
	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		stepCtx := tradeCtx.TradeStep1(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), types.TradeTypeSell)
		_, amount, _, err = k.ExecuteTradeStep(stepCtx)
		return err
	}))

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		stepCtx := tradeCtx.TradeStep2(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), amount, types.TradeTypeSell)
		usedAmount, receivedAmount, _, err = k.ExecuteTradeStep(stepCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	coins = k.BankKeeper.SpendableCoins(ctx, dexAcc)
	require.Equal(t, int64(100), usedAmount.Int64())
	require.Equal(t, int64(100), receivedAmount.Int64())
	require.Equal(t, int64(200), coins.AmountOf(constants.BaseCurrency).Int64())
	require.Equal(t, int64(0), coins.AmountOf(constants.KUSD).Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestSingleTrade3(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 100))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 100))

	offer := math.NewInt(100)
	fee := math.LegacyNewDecWithPrec(4, 2) // 4%

	liq := k.GetLiquiditySum(ctx, constants.BaseCurrency)
	require.Equal(t, int64(100), liq.Int64())

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         offer,
		TradeDenomGiving:    constants.KUSD,
		TradeDenomReceiving: constants.BaseCurrency,
		CoinSource:          keepertest.Bob,
		CoinTarget:          keepertest.Bob,
		OrdersCaches:        k.NewOrdersCaches(ctx),
		FlatPrice:           &constant_product.FlatPrice{},
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 fee,
	}

	var (
		amountUsed     math.Int
		amountReceived math.Int
		feePaid        math.Int
		err            error
	)

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		stepCtx := tradeCtx.TradeStep1(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), types.TradeTypeSell)
		amountUsed, amountReceived, feePaid, err = k.ExecuteTradeStep(stepCtx)
		if err != nil {
			return err
		}

		stepCtx = tradeCtx.TradeStep2(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), amountReceived, types.TradeTypeSell)
		_, amountReceived, _, err = k.ExecuteTradeStep(stepCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.Equal(t, offer, amountUsed)
	require.Equal(t, int64(98), amountReceived.Int64())
	require.Equal(t, int64(2), feePaid.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestSingleTrade3a(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 100))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 100))

	offer := math.NewInt(100)
	fee := math.LegacyNewDecWithPrec(4, 2) // 4%

	liq := k.GetLiquiditySum(ctx, constants.BaseCurrency)
	require.Equal(t, int64(100), liq.Int64())

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         offer,
		TradeDenomGiving:    constants.KUSD,
		TradeDenomReceiving: constants.BaseCurrency,
		CoinSource:          keepertest.Bob,
		CoinTarget:          keepertest.Bob,
		OrdersCaches:        k.NewOrdersCaches(ctx),
		FlatPrice:           &constant_product.FlatPrice{},
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 fee,
	}

	var (
		amountUsed     math.Int
		amountReceived math.Int
		feePaid        math.Int
		err            error
	)

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		stepCtx := tradeCtx.TradeStep1(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), types.TradeTypeBuy)
		amountUsed, amountReceived, _, err = k.ExecuteTradeStep(stepCtx)
		if err != nil {
			return err
		}

		stepCtx = tradeCtx.TradeStep2(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), amountReceived, types.TradeTypeBuy)
		amountUsed, amountReceived, feePaid, err = k.ExecuteTradeStep(stepCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.Equal(t, int64(102), amountUsed.Int64())
	require.Equal(t, int64(2), feePaid.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestSingleTrade4(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 100))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 50))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Bob, constants.KUSD, 50))

	offer := math.NewInt(100)
	fee := math.LegacyNewDecWithPrec(4, 2) // 4%

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         offer,
		MaxPrice:            nil,
		TradeDenomGiving:    constants.KUSD,
		TradeDenomReceiving: constants.BaseCurrency,
		MinimumTradeAmount:  &offer,
		CoinSource:          keepertest.Bob,
		CoinTarget:          keepertest.Bob,
		FlatPrice:           &constant_product.FlatPrice{},
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 fee,
	}

	var (
		amountUsed     math.Int
		amountReceived math.Int
		feePaid        math.Int
		err            error
	)

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		stepCtx := tradeCtx.TradeStep1(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), types.TradeTypeSell)
		amountUsed, amountReceived, feePaid, err = k.ExecuteTradeStep(stepCtx)
		if err != nil {
			return err
		}

		stepCtx = tradeCtx.TradeStep2(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), amountReceived, types.TradeTypeSell)
		_, amountReceived, _, err = k.ExecuteTradeStep(stepCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.Equal(t, offer, amountUsed)
	require.Equal(t, int64(98), amountReceived.Int64())
	require.Equal(t, int64(2), feePaid.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestSingleTrade5(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 10_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 10_000))

	offer := math.NewInt(10_000)
	fee := math.LegacyNewDecWithPrec(4, 2) // 4%

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         offer,
		MinimumTradeAmount:  &offer,
		CoinSource:          keepertest.Carol,
		CoinTarget:          keepertest.Carol,
		TradeDenomGiving:    constants.BaseCurrency,
		TradeDenomReceiving: constants.KUSD,
		FlatPrice:           &constant_product.FlatPrice{},
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 fee,
	}

	var (
		amountUsed     math.Int
		amountReceived math.Int
		feePaid1       math.Int
		feePaid2       math.Int
		err            error
	)

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		stepCtx := tradeCtx.TradeStep1(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), types.TradeTypeSell)
		amountUsed, amountReceived, feePaid1, err = k.ExecuteTradeStep(stepCtx)
		if err != nil {
			return err
		}

		stepCtx = tradeCtx.TradeStep2(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), amountReceived, types.TradeTypeSell)
		_, amountReceived, feePaid2, err = k.ExecuteTradeStep(stepCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.Equal(t, offer, amountUsed)
	require.Equal(t, int64(9_800), amountReceived.Int64())
	require.Equal(t, int64(0), feePaid1.Int64())
	require.Equal(t, int64(200), feePaid2.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestSingleTrade6(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 10_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 10_000))

	offer := math.NewInt(10_000)
	fee := math.LegacyNewDecWithPrec(4, 2) // 4%

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         offer,
		MinimumTradeAmount:  &offer,
		CoinSource:          keepertest.Carol,
		CoinTarget:          keepertest.Carol,
		TradeDenomGiving:    constants.BaseCurrency,
		TradeDenomReceiving: constants.KUSD,
		FlatPrice:           &constant_product.FlatPrice{},
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 fee,
	}

	var (
		amountUsed     math.Int
		amountReceived math.Int
		feePaid1       math.Int
		feePaid2       math.Int
		err            error
	)

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		stepCtx := tradeCtx.TradeStep1(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), types.TradeTypeSell)
		amountUsed, amountReceived, feePaid1, err = k.ExecuteTradeStep(stepCtx)
		if err != nil {
			return err
		}

		stepCtx = tradeCtx.TradeStep2(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), amountReceived, types.TradeTypeSell)
		_, amountReceived, feePaid2, err = k.ExecuteTradeStep(stepCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.Equal(t, offer, amountUsed)
	require.Equal(t, int64(9_800), amountReceived.Int64())
	require.Equal(t, int64(0), feePaid1.Int64())
	require.Equal(t, int64(200), feePaid2.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestSingleTrade7(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 100))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Bob, constants.KUSD, 50))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Bob, constants.KUSD, 50))

	offer := math.NewInt(100)
	fee := math.LegacyNewDecWithPrec(4, 2) // 4%

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         offer,
		MinimumTradeAmount:  &offer,
		CoinSource:          keepertest.Carol,
		CoinTarget:          keepertest.Carol,
		TradeDenomGiving:    constants.BaseCurrency,
		TradeDenomReceiving: constants.KUSD,
		FlatPrice:           &constant_product.FlatPrice{},
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 fee,
	}

	var (
		amountUsed     math.Int
		amountReceived math.Int
		feePaid1       math.Int
		feePaid2       math.Int
		err            error
	)

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		stepCtx := tradeCtx.TradeStep1(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), types.TradeTypeSell)
		amountUsed, amountReceived, feePaid1, err = k.ExecuteTradeStep(stepCtx)
		if err != nil {
			return err
		}

		stepCtx = tradeCtx.TradeStep2(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), amountReceived, types.TradeTypeSell)
		_, amountReceived, feePaid2, err = k.ExecuteTradeStep(stepCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.Equal(t, offer, amountUsed)
	require.Equal(t, int64(98), amountReceived.Int64())
	require.Equal(t, int64(0), feePaid1.Int64())
	require.Equal(t, int64(2), feePaid2.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestSingleTrade8(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 500_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 5_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 5_000))

	pair, err := k.GetLiquidityPair(ctx, constants.KUSD)
	require.NoError(t, err)
	require.Equal(t, math.LegacyNewDec(115_000), pair.VirtualOther)

	ordersCache := k.NewOrdersCaches(ctx)
	maximum := k.CalculateSingleGivableAmount(ordersCache, constants.BaseCurrency, constants.KUSD, nil)
	require.NotNil(t, maximum)

	receivedAmount, _, err := k.CalculateSingleSell(ctx, constants.BaseCurrency, constants.KUSD, *maximum, math.LegacyZeroDec())
	require.NoError(t, err)

	liqSum := k.GetLiquiditySum(ctx, constants.KUSD)
	require.Equal(t, receivedAmount.TruncateInt().Int64(), liqSum.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestSingleTrade9(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 500_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 5_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 5_000))

	pair, err := k.GetLiquidityPair(ctx, constants.KUSD)
	require.NoError(t, err)
	require.Equal(t, math.LegacyNewDec(115_000), pair.VirtualOther)

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))

	offer := math.NewInt(10_000)
	tradeCtx := types.TradeContext{
		Context:             ctx,
		CoinSource:          keepertest.Carol,
		CoinTarget:          keepertest.Carol,
		TradeAmount:         offer,
		MinimumTradeAmount:  &offer,
		TradeDenomGiving:    constants.BaseCurrency,
		TradeDenomReceiving: constants.KUSD,
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 math.LegacyNewDecWithPrec(1, 2),
	}

	var tradeResult types.TradeResult
	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		tradeResult, err = k.ExecuteSell(tradeCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.NoError(t, err)
	require.Equal(t, int64(2438), tradeResult.AmountReceived.Int64())
	require.Equal(t, int64(10_000), tradeResult.AmountGiven.Int64())
	require.Equal(t, int64(12), tradeResult.FeeBase.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestSingleTrade10(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 500_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 5_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 5_000))

	pair, err := k.GetLiquidityPair(ctx, constants.KUSD)
	require.NoError(t, err)
	require.Equal(t, math.LegacyNewDec(115_000), pair.VirtualOther)

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))

	offer := math.NewInt(10_000)

	tradeCtx := types.TradeContext{
		Context:             ctx,
		CoinSource:          keepertest.Carol,
		CoinTarget:          keepertest.Carol,
		TradeAmount:         offer,
		MinimumTradeAmount:  &offer,
		TradeDenomGiving:    constants.KUSD,
		TradeDenomReceiving: constants.BaseCurrency,
		FlatPrice:           &constant_product.FlatPrice{},
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 math.LegacyNewDecWithPrec(2, 3),
	}

	var tradeResult types.TradeResult
	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		tradeResult, err = k.ExecuteSell(tradeCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.Equal(t, int64(9_990), tradeResult.AmountReceived.Int64())
	require.Equal(t, int64(10_000), tradeResult.AmountGiven.Int64())
	require.Equal(t, int64(10), tradeResult.FeeOther.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestSingleTrade11(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 500_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 5_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 5_000))

	startAmount := int64(1000)
	fee := math.LegacyZeroDec()

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         math.NewInt(startAmount),
		TradeDenomGiving:    constants.BaseCurrency,
		TradeDenomReceiving: constants.KUSD,
	}

	tradeResult1, err := k.SimulateSellWithFee(tradeCtx, fee)
	require.NoError(t, err)

	tradeCtx = types.TradeContext{
		Context:             ctx,
		TradeAmount:         tradeResult1.AmountReceived,
		TradeDenomGiving:    constants.KUSD,
		TradeDenomReceiving: constants.BaseCurrency,
	}

	tradeResult2, err := k.SimulateSellWithFee(tradeCtx, fee)
	require.NoError(t, err)

	// Not exactly 1000 due to rounding
	require.Equal(t, int64(994), tradeResult2.AmountReceived.Int64())
}

func TestTrade1(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, keepertest.Pow(2)))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, keepertest.Pow(2)))

	liq := k.GetLiquiditySum(ctx, constants.BaseCurrency)
	require.Equal(t, math.NewInt(2_000_000), liq)
	liq = k.GetLiquiditySum(ctx, constants.KUSD)
	require.Equal(t, math.NewInt(2_000_000), liq)

	pair, err := k.GetLiquidityPair(ctx, constants.KUSD)
	require.NoError(t, err)
	require.Equal(t, keepertest.PowDec(6), pair.VirtualBase)
	require.Equal(t, math.LegacyZeroDec(), pair.VirtualOther)

	res, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Bob,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         keepertest.PowInt64String(1),
	})

	require.Nil(t, err)
	require.Equal(t, int64(222111), res.AmountReceived)

	addr, _ := sdk.AccAddressFromBech32(keepertest.Bob)

	coins := k.BankKeeper.SpendableCoins(ctx, addr)
	coinBase := getCoin(coins, constants.BaseCurrency)
	require.Equal(t, "99999000000", coinBase.Amount.String())

	coinKUSD := getCoin(coins, constants.KUSD)
	expected := 100000000000 + res.AmountReceived
	require.Equal(t, expected, coinKUSD.Amount.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade2(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, keepertest.Pow(1)))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, keepertest.Pow(1)))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))

	_, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         keepertest.PowInt64String(1),
	})

	require.NoError(t, err)
	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))

	liq := k.GetLiquiditySum(ctx, constants.BaseCurrency)
	require.Equal(t, math.NewInt(2000000), liq)
}

func TestTrade3(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, keepertest.Pow(2)))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Bob, constants.BaseCurrency, keepertest.Pow(2)))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, keepertest.Pow(1)))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Bob, constants.KUSD, keepertest.Pow(1)))

	_, _ = keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         keepertest.PowInt64String(2),
	})

	liq := k.GetLiquiditySum(ctx, constants.BaseCurrency)
	require.Equal(t, liq, math.NewInt(6_000_000))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
}

func TestTrade4(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, keepertest.Pow(1)))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Bob, constants.KUSD, keepertest.Pow(1)))

	_, _ = keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         keepertest.PowInt64String(1),
	})

	liq := k.GetLiquiditySum(ctx, constants.BaseCurrency)
	require.Equal(t, math.NewInt(2_000_000), liq)

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade5(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, keepertest.Pow(2)))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, keepertest.Pow(2)))

	_, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Bob,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         keepertest.PowInt64String(1),
	})
	require.NoError(t, err)

	pair, _ := k.GetLiquidityPair(ctx, constants.KUSD)
	require.Equal(t, int64(6000000), pair.VirtualBase.RoundInt().Int64())
	require.Equal(t, math.LegacyNewDec(0), pair.VirtualOther)

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade6(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, keepertest.Pow(2)))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, keepertest.Pow(2)))

	_, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Bob,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         keepertest.PowInt64String(1),
	})

	require.NoError(t, err)

	acc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolReserve)
	coins := k.BankKeeper.SpendableCoins(ctx, acc.GetAddress())
	require.Equal(t, 1, len(coins))
	require.Equal(t, int64(56), coins[0].Amount.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade7(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 50000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 50000))

	_, _ = keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Bob,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         "2000",
	})

	pair, err := k.GetLiquidityPair(ctx, constants.KUSD)
	require.NoError(t, err)
	require.True(t, pair.VirtualBase.GTE(math.LegacyZeroDec()))
	require.True(t, pair.VirtualOther.GTE(math.LegacyZeroDec()))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade8(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 50000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 50000))

	price1, err := k.CalculatePrice(ctx, constants.BaseCurrency, constants.KUSD)
	require.NoError(t, err)

	_, err = keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Bob,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         "2000",
	})
	require.NoError(t, err)

	price2, err := k.CalculatePrice(ctx, constants.BaseCurrency, constants.KUSD)
	require.NoError(t, err)
	require.True(t, price2.GT(price1))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade9(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 50000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 50000))

	price1, err := k.CalculatePrice(ctx, constants.BaseCurrency, constants.KUSD)
	require.NoError(t, err)

	_, err = keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Bob,
		DenomGiving:    constants.KUSD,
		DenomReceiving: constants.BaseCurrency,
		Amount:         "2000",
	})
	require.NoError(t, err)

	price2, err := k.CalculatePrice(ctx, constants.BaseCurrency, constants.KUSD)
	require.NoError(t, err)
	require.True(t, price2.LT(price1))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade11(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 50_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 25_000))

	_, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         "100000",
	})

	require.NoError(t, err)
	require.Equal(t, 2, len(k.LiquidityIterator(ctx, constants.KUSD).GetAll()))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade12(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 50_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 25_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Bob, constants.KUSD, 25_000))

	_, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         "100000",
	})

	require.NoError(t, err)

	require.Equal(t, 3, len(k.LiquidityIterator(ctx, constants.KUSD).GetAll()))
	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade13(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 50_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 25_000))

	pair, err := k.GetLiquidityPair(ctx, constants.KUSD)
	require.NoError(t, err)
	require.Equal(t, math.LegacyZeroDec(), pair.VirtualOther)

	_, err = keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         "100000",
	})
	require.NoError(t, err)

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade14(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 50_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 25_000))

	pair, err := k.GetLiquidityPair(ctx, constants.KUSD)
	require.NoError(t, err)
	require.Equal(t, math.LegacyZeroDec(), pair.VirtualOther)

	response, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         "100000",
		MaxPrice:       "0.01",
	})

	require.Nil(t, response)
	require.Error(t, err)

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade15(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 50_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 25_000))

	pair, err := k.GetLiquidityPair(ctx, constants.KUSD)
	require.NoError(t, err)
	require.Equal(t, math.LegacyZeroDec(), pair.VirtualOther)

	response, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         "100000",
	})

	require.NoError(t, err)
	require.NotNil(t, response)

	price := float64(response.AmountGiven) / float64(response.AmountReceived)
	require.Equal(t, 8.004482510205715, price)

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade17(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 50_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 1_000))

	pair, err := k.GetLiquidityPair(ctx, constants.KUSD)
	require.NoError(t, err)
	require.Equal(t, math.LegacyNewDec(11500), pair.VirtualOther)

	response, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         "4000",
	})

	require.NoError(t, err)
	// AmountReceived is not 1000 because of fee
	require.Equal(t, int64(925), response.AmountReceived)
	require.Equal(t, int64(4000), response.AmountGiven)

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade18(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 500_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 5_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 5_000))

	pair, err := k.GetLiquidityPair(ctx, constants.KUSD)
	require.NoError(t, err)
	require.Equal(t, math.LegacyNewDec(115000), pair.VirtualOther)

	response, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         "10000",
	})

	require.NoError(t, err)
	// AmountReceived is not 1000 because of fee
	require.Equal(t, int64(2449), response.AmountReceived)
	require.Equal(t, int64(10000), response.AmountGiven)

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade22(t *testing.T) {
	amountReceived1 := testSmallDenomTrade(t, 20_000)
	amountReceived2 := testSmallDenomTrade(t, 1000)
	require.Greater(t, amountReceived1, amountReceived2)
}

func TestTrade23(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 10_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 10_000))

	price1, err := k.CalculatePrice(ctx, constants.BaseCurrency, constants.KUSD)
	require.NoError(t, err)

	ratio1, err := k.DenomKeeper.GetRatio(ctx, constants.KUSD)
	require.NoError(t, err)

	_, err = keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         "1000",
		MaxPrice:       "",
	})
	require.NoError(t, err)

	ratio2, err := k.DenomKeeper.GetRatio(ctx, constants.KUSD)
	require.NoError(t, err)
	require.True(t, ratio1.Ratio.GT(ratio2.Ratio))

	price2, err := k.CalculatePrice(ctx, constants.BaseCurrency, constants.KUSD)
	require.NoError(t, err)
	require.True(t, price1.LT(price2))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func testSmallDenomTrade(t *testing.T, amount int64) int64 {
	_, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 25_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, amount))

	response, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    constants.BaseCurrency,
		DenomReceiving: constants.KUSD,
		Amount:         "1000",
	})

	require.NoError(t, err)
	return response.GetAmountReceived()
}

func getCoin(coins []sdk.Coin, denom string) sdk.Coin {
	for _, coin := range coins {
		if coin.Denom == denom {
			return coin
		}
	}

	return sdk.Coin{}
}

func TestAddress(t *testing.T) {
	bz, err := sdk.GetFromBech32("axelar1txu08a5y7mylplyyvn9pwnfcderrz28eag23zj", "axelar")
	require.NoError(t, err)

	addr := sdk.AccAddress(bz)
	addrStr, err := bech32.ConvertAndEncode("migaloo", addr.Bytes())
	_ = addrStr
	require.NoError(t, err)
}

func TestTrade24(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 10_000))

	liqOtherSum1 := k.GetLiquiditySum(ctx, constants.KUSD)

	_, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    constants.KUSD,
		DenomReceiving: constants.BaseCurrency,
		Amount:         "1000",
		MaxPrice:       "",
	})
	require.NoError(t, err)

	liqOtherSum2 := k.GetLiquiditySum(ctx, constants.KUSD)
	require.Equal(t, math.NewInt(1000), liqOtherSum2.Sub(liqOtherSum1))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade25(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	amount := int64(10_000)
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, amount))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, amount))

	pair, err := k.GetLiquidityPair(ctx, constants.KUSD)
	require.NoError(t, err)

	A := k.GetFullLiquidityOther(ctx, constants.KUSD)
	B := k.GetFullLiquidityBase(ctx, constants.KUSD)
	b := B.Sub(pair.VirtualBase)
	c := pair.VirtualBase

	var maximum1 *math.LegacyDec
	if c.IsPositive() {
		m := A.Mul(b.Add(c)).Quo(c).Sub(A)
		maximum1 = &m
	}

	ordersCache := k.NewOrdersCaches(ctx)
	maximum2 := k.CalculateSingleGivableAmount(ordersCache, constants.KUSD, constants.BaseCurrency, nil)
	require.NotNil(t, maximum2)
	require.Equal(t, (*maximum1).TruncateInt(), (*maximum2).TruncateInt())

	A = k.GetFullLiquidityBase(ctx, constants.KUSD)
	B = k.GetFullLiquidityOther(ctx, constants.KUSD)
	b = B.Sub(pair.VirtualOther)
	c = pair.VirtualOther

	maximum1 = nil
	if c.IsPositive() {
		m := A.Mul(b.Add(c)).Quo(c).Sub(A)
		maximum1 = &m
	}

	//maximum2 = k.CalculateSingleMaximumTradableAmount(ctx, constants.BaseCurrency, constants.KUSD, nil)
	//require.Equal(t, maximum1.RoundInt(), maximum2)

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade26(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	amount := int64(10_000)
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, amount))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, amount))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, "uwusdc", amount))

	ordersCache := k.NewOrdersCaches(ctx)
	var maximum1 *math.LegacyDec
	maximum1 = k.CalculateSingleGivableAmount(ordersCache, constants.BaseCurrency, constants.KUSD, nil)
	maximum1 = k.CalculateSingleGivableAmount(ordersCache, "uwusdc", constants.BaseCurrency, maximum1)

	var maximum2 *math.Int
	maximum2 = k.CalculateMaximumTradableAmount(ordersCache, "uwusdc", constants.KUSD)

	require.NotNil(t, maximum2)
	require.Equal(t, maximum1.TruncateInt().Int64(), maximum2.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade27(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 10_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 10))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, "uwusdc", 10_000))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade28(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 100_000_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 100))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade30(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 100_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 10_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, "uwusdc", 100_000))

	maxPrice := math.LegacyNewDecWithPrec(105, 1)
	res, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Carol,
		DenomGiving:    "uwusdc",
		DenomReceiving: constants.KUSD,
		Amount:         "1000",
		MaxPrice:       maxPrice.String(),
	})

	require.NoError(t, err)
	require.True(t, res.AmountReceived > 0)

	maxPriceF, _ := maxPrice.Float64()

	var paidPrice float64
	if res.AmountReceived > 0 {
		paidPrice = float64(res.AmountGiven) / float64(res.AmountReceived)
	}

	require.LessOrEqual(t, paidPrice, maxPriceF)

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade31(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 10))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 10))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, "uwusdc", 10))

	ordersCache := k.NewOrdersCaches(ctx)

	maximumTradableAmount := k.CalculateMaximumTradableAmount(ordersCache, constants.KUSD, constants.BaseCurrency)

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         *maximumTradableAmount,
		MinimumTradeAmount:  maximumTradableAmount,
		CoinSource:          keepertest.Bob,
		CoinTarget:          keepertest.Bob,
		TradeDenomGiving:    "uwusdc",
		TradeDenomReceiving: constants.KUSD,
		FlatPrice:           &constant_product.FlatPrice{},
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 math.LegacyZeroDec(),
	}

	var (
		amountReceivedNet math.Int
		err               error
	)

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		stepCtx := tradeCtx.TradeStep1(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), types.TradeTypeSell)
		_, amountReceivedNet, _, err = k.ExecuteTradeStep(stepCtx)
		if err != nil {
			return err
		}

		stepCtx = tradeCtx.TradeStep2(tradeCtx.OrdersCaches.ReserveFeeShare.Get(), amountReceivedNet, types.TradeTypeSell)
		_, _, _, err = k.ExecuteTradeStep(stepCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade32(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 10_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 10_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, "uwusdc", 10_000))

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         math.NewInt(1000),
		CoinSource:          keepertest.Bob,
		CoinTarget:          keepertest.Bob,
		TradeDenomGiving:    "uwusdc",
		TradeDenomReceiving: constants.KUSD,
		FlatPrice:           &constant_product.FlatPrice{},
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 math.LegacyNewDecWithPrec(4, 2),
	}

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		_, err := k.ExecuteSell(tradeCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade33(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 10_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, "uwusdc", 10_000))

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         math.NewInt(1000),
		CoinSource:          keepertest.Bob,
		CoinTarget:          keepertest.Bob,
		TradeDenomGiving:    "uwusdc",
		TradeDenomReceiving: constants.BaseCurrency,
		FlatPrice:           &constant_product.FlatPrice{},
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 math.LegacyNewDecWithPrec(4, 2),
	}

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		_, err := k.ExecuteSell(tradeCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade34(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 10))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, "uwusdc", 10))

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         math.NewInt(10000),
		CoinSource:          keepertest.Bob,
		CoinTarget:          keepertest.Bob,
		TradeDenomGiving:    constants.BaseCurrency,
		TradeDenomReceiving: "uwusdc",
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 math.LegacyNewDecWithPrec(4, 2),
	}

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		_, err := k.ExecuteSell(tradeCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade35(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 10_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 10_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, "uwusdc", 10_000))

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         math.NewInt(1000),
		CoinSource:          keepertest.Bob,
		CoinTarget:          keepertest.Bob,
		TradeDenomGiving:    "uwusdc",
		TradeDenomReceiving: constants.KUSD,
		FlatPrice:           &constant_product.FlatPrice{},
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 math.LegacyNewDecWithPrec(4, 2),
	}

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		_, err := k.ExecuteSell(tradeCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, tradePoolEmpty(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade36(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 10_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 10_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, "uwusdc", 10_000))

	addr, _ := sdk.AccAddressFromBech32(keepertest.Bob)
	coins1 := k.BankKeeper.SpendableCoins(ctx, addr)
	funds1_uwusdc := coins1.AmountOf("uwusdc").Int64()
	funds1_ukusd := coins1.AmountOf(constants.KUSD).Int64()

	tradeAmount := int64(1000)

	liq1_uwusdc := k.GetLiquiditySum(ctx, "uwusdc")
	liq1_ukusd := k.GetLiquiditySum(ctx, constants.KUSD)

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         math.NewInt(tradeAmount),
		CoinSource:          keepertest.Bob,
		CoinTarget:          keepertest.Bob,
		TradeDenomGiving:    "uwusdc",
		TradeDenomReceiving: constants.KUSD,
		FlatPrice:           &constant_product.FlatPrice{},
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 math.LegacyZeroDec(),
	}

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		_, err := k.ExecuteSell(tradeCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	coins2 := k.BankKeeper.SpendableCoins(ctx, addr)
	funds2_uwusdc := coins2.AmountOf("uwusdc").Int64()
	funds2_ukusd := coins2.AmountOf(constants.KUSD).Int64()

	liq2_uwusdc := k.GetLiquiditySum(ctx, "uwusdc")
	liq2_ukusd := k.GetLiquiditySum(ctx, constants.KUSD)

	require.Equal(t, funds1_uwusdc-tradeAmount, funds2_uwusdc)
	require.Equal(t, funds1_ukusd+tradeAmount, funds2_ukusd)

	require.Equal(t, liq1_uwusdc.Int64()+tradeAmount, liq2_uwusdc.Int64())
	require.Equal(t, liq1_ukusd.Int64()-tradeAmount, liq2_ukusd.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade37(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 10_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 10_000))

	liqBase := k.LiquidityIterator(ctx, constants.BaseCurrency).GetAll()
	require.Equal(t, 1, len(liqBase))
	require.Equal(t, int64(10_000), liqBase[0].Amount.Int64())

	liqOther := k.LiquidityIterator(ctx, constants.KUSD).GetAll()
	require.Equal(t, 1, len(liqOther))

	tradeAmount := int64(1000)

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         math.NewInt(tradeAmount),
		CoinSource:          keepertest.Bob,
		CoinTarget:          keepertest.Bob,
		TradeDenomGiving:    constants.BaseCurrency,
		TradeDenomReceiving: constants.KUSD,
		FlatPrice:           &constant_product.FlatPrice{},
		OrdersCaches:        k.NewOrdersCaches(ctx),
		TradeBalances:       dexkeeper.NewTradeBalances(),
		Fee:                 math.LegacyNewDecWithPrec(1, 2),
	}

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		tradeCtx.Context = innerCtx
		_, err := k.ExecuteSell(tradeCtx)
		return err
	}))

	require.NoError(t, tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper))

	liqBase = k.LiquidityIterator(ctx, constants.BaseCurrency).GetAll()
	require.Equal(t, 2, len(liqBase))
	require.Equal(t, int64(10_000), liqBase[0].Amount.Int64())
	require.Equal(t, int64(1_000), liqBase[1].Amount.Int64())

	liqOther = k.LiquidityIterator(ctx, constants.KUSD).GetAll()
	require.Equal(t, 2, len(liqOther))
	require.Equal(t, int64(9_000), liqOther[0].Amount.Int64())
	require.Equal(t, int64(3), liqOther[1].Amount.Int64())

	require.True(t, liquidityBalanced(ctx, k))
	require.NoError(t, checkCache(ctx, k))
}

func TestTrade38(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 1_000_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 1_000_000))

	offer := math.LegacyNewDec(10_000)
	amountReceived, _, err := k.CalculateSingleSell(ctx, constants.BaseCurrency, constants.KUSD, offer, math.LegacyZeroDec())
	require.NoError(t, err)

	amountToGive, _, err := k.CalculateSingleBuy(ctx, constants.BaseCurrency, constants.KUSD, amountReceived, math.LegacyZeroDec())
	require.NoError(t, err)

	require.Equal(t, offer.RoundInt64(), amountToGive.RoundInt64())
}

func TestTrade39(t *testing.T) {
	_, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 1_000_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 1_000_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, "uwusdc", 1_000_000))

	res, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Bob,
		DenomGiving:    constants.KUSD,
		DenomReceiving: "uwusdc",
		Amount:         "10000",
	})
	require.NoError(t, err)

	price := float64(res.AmountGiven) / float64(res.AmountReceived)

	_, err = keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:        keepertest.Bob,
		DenomGiving:    constants.KUSD,
		DenomReceiving: "uwusdc",
		Amount:         "10000",
		MaxPrice:       fmt.Sprintf("%.8f", price),
	})
	require.Error(t, err)
}

func TestTrade40(t *testing.T) {
	_, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 1_000_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 1_000_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, "uwusdc", 1_000_000))

	res, err := keepertest.Buy(ctx, msg, &types.MsgBuy{
		Creator:        keepertest.Bob,
		DenomGiving:    constants.KUSD,
		DenomReceiving: "uwusdc",
		Amount:         "10000",
	})
	require.NoError(t, err)

	price := float64(res.AmountGiven) / float64(res.AmountReceived)

	_, err = keepertest.Buy(ctx, msg, &types.MsgBuy{
		Creator:        keepertest.Bob,
		DenomGiving:    constants.KUSD,
		DenomReceiving: "uwusdc",
		Amount:         "10000",
		MaxPrice:       fmt.Sprintf("%.8f", price),
	})
	require.Error(t, err)
}

func TestTrade41(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 1_000_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 1_000_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, "uwusdc", 1_000_000))

	acc1, _ := sdk.AccAddressFromBech32(keepertest.Alice)
	acc2, _ := sdk.AccAddressFromBech32(keepertest.Dave)

	amount := int64(10_000)
	coins := sdk.NewCoins(sdk.NewCoin(constants.KUSD, math.NewInt(amount)))
	_ = k.BankKeeper.SendCoinsFromAccountToModule(ctx, acc1, types.PoolReserve, coins)
	_ = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.PoolReserve, acc2, coins)

	res, err := keepertest.Buy(ctx, msg, &types.MsgBuy{
		Creator:        keepertest.Dave,
		DenomGiving:    constants.KUSD,
		DenomReceiving: "uwusdc",
		Amount:         "100000",
	})
	require.NoError(t, err)
	require.Equal(t, amount, res.AmountGiven)
}

func TestTrade42(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	pool := map[string]int64{
		"skbtc":   8496,
		"swbtc":   737567446,
		"ucwusdc": 4118700571,
		"ucwusdt": 492046007,
		"ukopi":   5223913832979,
		"ukusd":   479921845,
		"uwusdc":  38585199076,
		"uwusdt":  698434,
	}

	ratios := map[string]string{
		"sckbtc":     "0.000405912156352312",
		"skbtc":      "0.000896829507787797",
		"swbtc":      "0.000896829508394457",
		"uarbstusdc": "0.101556613612585061",
		"uarbstusdt": "0.101756063584109115",
		"uckusd":     "0.101478039088050841",
		"ucwusdc":    "0.249106333509990319",
		"ucwusdt":    "0.203250771531685959",
		"ukusd":      "0.196939223849083103",
		"uwusdc":     "0.211427997614305922",
		"uwusdt":     "0.186917656024168530",
	}

	balance := map[string]int64{
		"ucwusdc": 9180625,
		"ucwusdt": 3216676466,
		"ukusd":   3899152227,
	}

	for denom, amount := range balance {
		keepertest.AddFunds(ctx, t, k.BankKeeper, denom, keepertest.Dave, amount)
	}

	keepertest.SetLiquidity(ctx, k.BankKeeper, k, t, pool)

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		for denom, ratio := range ratios {
			r, _ := math.LegacyNewDecFromStr(ratio)
			k.DenomKeeper.SetRatio(innerCtx, denomtypes.Ratio{Denom: denom, Ratio: r})
		}

		return nil
	}))

	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.BaseCurrency, 1_000_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, constants.KUSD, 1_000_000))
	require.NoError(t, keepertest.AddLiquidity(ctx, msg, keepertest.Alice, "uwusdc", 1_000_000))

	acc1, _ := sdk.AccAddressFromBech32(keepertest.Alice)
	acc2, _ := sdk.AccAddressFromBech32(keepertest.Dave)

	amount := int64(10_000)
	coins := sdk.NewCoins(sdk.NewCoin(constants.KUSD, math.NewInt(amount)))
	_ = k.BankKeeper.SendCoinsFromAccountToModule(ctx, acc1, types.PoolReserve, coins)
	_ = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.PoolReserve, acc2, coins)

	_, err := keepertest.Buy(ctx, msg, &types.MsgBuy{
		Creator:        keepertest.Dave,
		DenomGiving:    "ucwusdc",
		DenomReceiving: constants.KUSD,
		Amount:         "10000000",
	})
	require.NoError(t, err)
}

func TestTrade43(t *testing.T) {
	k, msg, ctx := keepertest.SetupDexMsgServer(t)

	pool := map[string]int64{
		"inj":   829132866175164798,
		"ukusd": 1740434604,
		"ukopi": 49994130281493,
	}

	ratios := map[string]string{
		"inj":   "11310893732.791635615102371449",
		"skbtc": "0.235974469935270325",
	}

	balance := map[string]int64{
		"inj": 37582987632368903,
	}

	for denom, amount := range balance {
		keepertest.AddFunds(ctx, t, k.BankKeeper, denom, keepertest.Dave, amount)
	}

	keepertest.SetLiquidity(ctx, k.BankKeeper, k, t, pool)

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		for denom, ratio := range ratios {
			r, _ := math.LegacyNewDecFromStr(ratio)
			k.DenomKeeper.SetRatio(innerCtx, denomtypes.Ratio{Denom: denom, Ratio: r})
		}

		return nil
	}))

	_, err := keepertest.Sell(ctx, msg, &types.MsgSell{
		Creator:            keepertest.Dave,
		DenomGiving:        "inj",
		DenomReceiving:     "ukusd",
		Amount:             "1000000000000000",
		MaxPrice:           "48057996063.181141677468913534",
		MinimumTradeAmount: "1000000000000000",
	})
	require.NoError(t, err)
}

func liquidityBalanced(ctx context.Context, k dexkeeper.Keeper) bool {
	acc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolLiquidity)
	coins := k.BankKeeper.SpendableCoins(ctx, acc.GetAddress())

	for _, denom := range k.DenomKeeper.Denoms(ctx) {
		liqSum := k.GetLiquiditySum(ctx, denom).Int64()
		summedLiq := k.SumLiquidity(ctx, denom).Int64()
		funds := coins.AmountOf(denom).Int64()

		if liqSum != funds || summedLiq != funds {
			fmt.Println(denom)
			fmt.Println(fmt.Sprintf("liqSum: %v", liqSum))
			fmt.Println(fmt.Sprintf("summedLiq: %v", summedLiq))
			fmt.Println(fmt.Sprintf("funds: %v", funds))

			return false
		}
	}

	return true
}

func tradePoolEmpty(ctx context.Context, k dexkeeper.Keeper) error {
	acc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolTrade)
	coins := k.BankKeeper.SpendableCoins(ctx, acc.GetAddress())

	for _, coin := range coins {
		if coin.Amount.GT(math.ZeroInt()) {
			return fmt.Errorf("trade pool: %v %v", coin.Amount.String(), coin.Denom)
		}
	}

	return nil
}

func checkCache(ctx context.Context, k dexkeeper.Keeper) error {
	return k.CheckCache(ctx)
}
