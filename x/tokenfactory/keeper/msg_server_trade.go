package keeper

import (
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/x/dex/constant_product"
	dexkeeper "github.com/kopi-money/kopi/x/dex/keeper"
	dextypes "github.com/kopi-money/kopi/x/dex/types"
	"github.com/kopi-money/kopi/x/tokenfactory/types"
)

type CalcAmount func(math.LegacyDec, math.LegacyDec, math.LegacyDec, math.LegacyDec, bool) math.LegacyDec

func (k msgServer) Sell(ctx context.Context, msg *types.MsgSell) (*types.MsgTradeResponse, error) {
	return k.Keeper.Sell(ctx, TradeData{
		factoryDenom:    msg.FullFactoryDenomName,
		creator:         msg.Creator,
		denomGiving:     msg.DenomGiving,
		denomReceiving:  msg.DenomReceiving,
		maxPrice:        msg.MaxPrice,
		tradeAmount:     msg.Amount,
		allowIncomplete: msg.AllowIncomplete,
	})
}

func (k Keeper) Sell(ctx context.Context, tradeData TradeData) (*types.MsgTradeResponse, error) {
	acc, _ := sdk.AccAddressFromBech32(tradeData.creator)

	factoryDenom, has := k.GetDenomByFullName(ctx, tradeData.factoryDenom)
	if !has {
		return nil, types.ErrDenomDoesntExists
	}

	pool, has := k.liquidityPools.Get(ctx, factoryDenom.FullName)
	if !has {
		return nil, types.ErrPoolDoesNotExist
	}

	if tradeData.denomGiving == tradeData.denomReceiving {
		return nil, types.ErrSameDenom
	}

	amountToGiveGross, err := dexkeeper.ParseAmount(tradeData.tradeAmount)
	if err != nil {
		return nil, fmt.Errorf("could not parse amount: %w", err)
	}

	if tradeData.maxPrice != "" {
		var maxPrice math.LegacyDec
		maxPrice, err = getMaxPrice(tradeData.maxPrice)
		if err != nil {
			return nil, err
		}

		priceTradeAmount := k.calculateMaxAmount(ctx, pool, tradeData.denomGiving, maxPrice, pool.PoolFee, constant_product.CalculateMaximumGiving)
		if priceTradeAmount.LT(amountToGiveGross) {
			if priceTradeAmount.LT(math.ZeroInt()) || !tradeData.allowIncomplete {
				return nil, types.ErrMarketPriceTooHigh
			}

			amountToGiveGross = priceTradeAmount
		}

		if amountToGiveGross.LTE(math.ZeroInt()) {
			return nil, types.ErrEmptyTrade
		}
	}

	if k.BankKeeper.SpendableCoin(ctx, acc, tradeData.denomGiving).Amount.LT(amountToGiveGross) {
		return nil, types.ErrInsufficientFunds
	}

	feesGiving, err := k.applyFees(ctx, &pool, acc, amountToGiveGross, tradeData.denomGiving)
	if err != nil {
		return nil, err
	}

	amountToGiveNet := amountToGiveGross.Sub(feesGiving.Fee())

	if amountToGiveNet.LT(math.NewInt(1000)) {
		return nil, types.ErrTradeAmountTooSmall
	}

	coins := sdk.NewCoins(sdk.NewCoin(tradeData.denomGiving, amountToGiveGross))
	if err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, acc, types.PoolFactoryLiquidity, coins); err != nil {
		return nil, fmt.Errorf("could not send coins from account to liquidity pool: %w", err)
	}

	amountToReceiveGross := k.constantProductSell(pool, tradeData.denomGiving, amountToGiveNet)
	feesReceiving, err := k.applyFees(ctx, &pool, acc, amountToReceiveGross, tradeData.denomReceiving)
	if err != nil {
		return nil, err
	}
	amountToReceiveNet := amountToReceiveGross.Sub(feesReceiving.Fee())

	pool.FactoryDenomAmount = getNewFactoryAmount(pool, tradeData.denomGiving, amountToGiveNet, amountToReceiveNet)
	pool.KCoinAmount = getNewKCoinAmount(pool, tradeData.denomGiving, amountToGiveNet, amountToReceiveNet)
	k.liquidityPools.Set(ctx, factoryDenom.FullName, pool)

	coins = sdk.NewCoins(sdk.NewCoin(tradeData.denomReceiving, amountToReceiveNet))
	if err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.PoolFactoryLiquidity, acc, coins); err != nil {
		return nil, fmt.Errorf("could not send coins from liquidity pool to account: %w", err)
	}

	var feeData FeeData
	if tradeData.denomGiving == pool.KCoin {
		feeData = feesGiving
	} else {
		feeData = feesReceiving
	}

	return &types.MsgTradeResponse{
		AmountGivenGross:    amountToGiveGross.Int64(),
		AmountGivenNet:      amountToGiveGross.Sub(feesGiving.Fee()).Int64(),
		AmountReceivedGross: amountToReceiveNet.Add(feesReceiving.Fee()).Int64(),
		AmountReceivedNet:   amountToReceiveNet.Int64(),
		Fee:                 feeData.Fee().Int64(),
		FeePool:             feeData.feePool.Int64(),
		FeeReserve:          feeData.feeReserve.Int64(),
	}, nil
}

func (k msgServer) Buy(ctx context.Context, msg *types.MsgBuy) (*types.MsgTradeResponse, error) {
	return k.Keeper.Buy(ctx, TradeData{
		factoryDenom:    msg.FullFactoryDenomName,
		creator:         msg.Creator,
		denomGiving:     msg.DenomGiving,
		denomReceiving:  msg.DenomReceiving,
		maxPrice:        msg.MaxPrice,
		tradeAmount:     msg.Amount,
		allowIncomplete: msg.AllowIncomplete,
	})
}

func (k Keeper) Buy(ctx context.Context, tradeData TradeData) (*types.MsgTradeResponse, error) {
	acc, _ := sdk.AccAddressFromBech32(tradeData.creator)

	factoryDenom, has := k.GetDenomByFullName(ctx, tradeData.factoryDenom)
	if !has {
		return nil, types.ErrDenomDoesntExists
	}

	pool, has := k.liquidityPools.Get(ctx, factoryDenom.FullName)
	if !has {
		return nil, types.ErrPoolDoesNotExist
	}

	if tradeData.denomGiving == tradeData.denomReceiving {
		return nil, types.ErrSameDenom
	}

	amountToReceiveNet, err := dexkeeper.ParseAmount(tradeData.tradeAmount)
	if err != nil {
		return nil, fmt.Errorf("could not parse amount: %w", err)
	}

	feesGiving, err := k.applyFees(ctx, &pool, acc, amountToReceiveNet, tradeData.denomGiving)
	if err != nil {
		return nil, err
	}

	amountToReceiveGross := amountToReceiveNet.Add(feesGiving.Fee())

	if tradeData.maxPrice != "" {
		var maxPrice math.LegacyDec
		maxPrice, err = getMaxPrice(tradeData.maxPrice)
		if err != nil {
			return nil, err
		}

		maxPrice = math.LegacyOneDec().Quo(maxPrice)
		priceTradeAmount := k.calculateMaxAmount(ctx, pool, tradeData.denomGiving, maxPrice, pool.PoolFee, constant_product.CalculateMaximumReceiving)
		if priceTradeAmount.LT(amountToReceiveGross) {
			if priceTradeAmount.LT(math.ZeroInt()) || !tradeData.allowIncomplete {
				return nil, types.ErrMarketPriceTooHigh
			}

			amountToReceiveGross = priceTradeAmount
		}

		if amountToReceiveGross.LTE(math.ZeroInt()) {
			return nil, types.ErrEmptyTrade
		}
	}

	if amountToReceiveGross.LT(math.NewInt(1000)) {
		return nil, types.ErrTradeAmountTooSmall
	}

	amountToGiveNet := k.constantProductBuy(pool, tradeData.denomGiving, amountToReceiveGross)
	feesReceiving, err := k.applyFees(ctx, &pool, acc, amountToReceiveGross, tradeData.denomReceiving)
	if err != nil {
		return nil, err
	}

	amountToGiveGross := amountToGiveNet.Add(feesReceiving.Fee())

	if k.BankKeeper.SpendableCoin(ctx, acc, tradeData.denomGiving).Amount.LT(amountToGiveGross) {
		return nil, types.ErrInsufficientFunds
	}

	coins := sdk.NewCoins(sdk.NewCoin(tradeData.denomGiving, amountToGiveGross))
	if err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, acc, types.PoolFactoryLiquidity, coins); err != nil {
		return nil, fmt.Errorf("could not send coins from account to liquidity pool: %w", err)
	}

	pool.FactoryDenomAmount = getNewFactoryAmount(pool, tradeData.denomGiving, amountToGiveNet, amountToReceiveNet)
	pool.KCoinAmount = getNewKCoinAmount(pool, tradeData.denomGiving, amountToGiveNet, amountToReceiveNet)
	k.liquidityPools.Set(ctx, factoryDenom.FullName, pool)

	coins = sdk.NewCoins(sdk.NewCoin(tradeData.denomReceiving, amountToReceiveNet))
	if err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.PoolFactoryLiquidity, acc, coins); err != nil {
		return nil, fmt.Errorf("could not send coins from liquidity pool to account: %w", err)
	}

	var feeData FeeData
	if tradeData.denomGiving == pool.KCoin {
		feeData = feesGiving
	} else {
		feeData = feesReceiving
	}

	return &types.MsgTradeResponse{
		AmountGivenGross:    amountToGiveGross.Int64(),
		AmountGivenNet:      amountToGiveGross.Sub(feesGiving.Fee()).Int64(),
		AmountReceivedGross: amountToReceiveNet.Add(feesReceiving.Fee()).Int64(),
		AmountReceivedNet:   amountToReceiveNet.Int64(),
		Fee:                 feeData.Fee().Int64(),
		FeePool:             feeData.feePool.Int64(),
		FeeReserve:          feeData.feeReserve.Int64(),
	}, nil
}

type TradeData struct {
	factoryDenom    string
	creator         string
	denomGiving     string
	denomReceiving  string
	maxPrice        string
	tradeAmount     string
	allowIncomplete bool

	calcAmountToGive    CalcAmount
	calcAmountToReceive CalcAmount
}

type FeeData struct {
	feePool    math.Int
	feeReserve math.Int
}

func (fd FeeData) Fee() math.Int {
	return fd.feePool.Add(fd.feeReserve)
}

func (k Keeper) applyFees(ctx context.Context, pool *types.LiquidityPool, acc sdk.AccAddress, amount math.Int, kCoin string) (FeeData, error) {
	if !k.DenomKeeper.IsKCoin(ctx, kCoin) {
		return FeeData{
			feePool:    math.ZeroInt(),
			feeReserve: math.ZeroInt(),
		}, nil
	}

	reserveFeeAmount, err := k.applyReserveFee(ctx, acc, amount, kCoin)
	if err != nil {
		return FeeData{}, err
	}

	poolFeeAmount := applyPoolFee(pool, amount)
	return FeeData{
		feePool:    poolFeeAmount,
		feeReserve: reserveFeeAmount,
	}, nil
}

func applyPoolFee(pool *types.LiquidityPool, amountToReceive math.Int) math.Int {
	poolFeeAmount := pool.PoolFee.Mul(amountToReceive.ToLegacyDec()).TruncateInt()
	pool.KCoinAmount = pool.KCoinAmount.Add(poolFeeAmount)
	return poolFeeAmount
}

func (k Keeper) applyReserveFee(ctx context.Context, acc sdk.AccAddress, tradeAmount math.Int, kCoin string) (math.Int, error) {
	if !k.DenomKeeper.IsKCoin(ctx, kCoin) {
		return tradeAmount, nil
	}

	reserveFee := k.GetParams(ctx).ReserveFee.Mul(tradeAmount.ToLegacyDec()).TruncateInt()
	coins := sdk.NewCoins(sdk.NewCoin(kCoin, reserveFee))
	if err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, acc, dextypes.PoolReserve, coins); err != nil {
		return math.Int{}, fmt.Errorf("could not send reserve fee to module: %w", err)
	}

	return reserveFee, nil
}

func getNewFactoryAmount(pool types.LiquidityPool, denomFrom string, tradeAmountGross, amountToReceive math.Int) math.Int {
	if pool.KCoin == denomFrom {
		// kCoin -> FactoryDenom => Factory denom decreases
		return pool.FactoryDenomAmount.Sub(amountToReceive)
	} else {
		// FactoryDenom -> kCoin => Factory denom increases
		return pool.FactoryDenomAmount.Add(tradeAmountGross)
	}
}

func getNewKCoinAmount(pool types.LiquidityPool, denomFrom string, tradeAmountGross, amountToReceive math.Int) math.Int {
	if pool.KCoin == denomFrom {
		// kCoin -> FactoryDenom => kCoin amount increases
		return pool.KCoinAmount.Add(tradeAmountGross)
	} else {
		// FactoryDenom -> kCoin => kCoin denom decreases
		return pool.KCoinAmount.Sub(amountToReceive)
	}
}

func (k Keeper) constantProductSell(pool types.LiquidityPool, denomGiving string, amount math.Int) math.Int {
	liqFrom, liqTo := getLiquidity(pool, denomGiving)
	amountDec, _ := constant_product.ConstantProductTradeSell(liqFrom, liqTo, amount.ToLegacyDec(), math.LegacyZeroDec())
	return amountDec.TruncateInt()
}

func (k Keeper) constantProductBuy(pool types.LiquidityPool, denomGiving string, amount math.Int) math.Int {
	liqFrom, liqTo := getLiquidity(pool, denomGiving)
	amountDec, _ := constant_product.ConstantProductTradeBuy(liqFrom, liqTo, amount.ToLegacyDec(), math.LegacyZeroDec())
	return amountDec.TruncateInt()
}

func (k Keeper) calculateMaxAmount(ctx context.Context, pool types.LiquidityPool, denomFrom string, maxPrice, poolFee math.LegacyDec, calculate constant_product.CalculateMaximumAmount) math.Int {
	liqFrom, liqTo := getLiquidity(pool, denomFrom)
	tradeFee := k.getTradeFee(ctx, poolFee)
	return calculate(liqFrom, liqTo, maxPrice, tradeFee).TruncateInt()
}

func getLiquidity(pool types.LiquidityPool, denomFrom string) (math.LegacyDec, math.LegacyDec) {
	var liqFrom, liqTo math.Int

	if denomFrom == pool.KCoin {
		liqFrom, liqTo = pool.KCoinAmount, pool.FactoryDenomAmount
	} else {
		liqFrom, liqTo = pool.FactoryDenomAmount, pool.KCoinAmount
	}

	return liqFrom.ToLegacyDec(), liqTo.ToLegacyDec()
}

func getMaxPrice(maxPriceString string) (math.LegacyDec, error) {
	maxPriceString = strings.ReplaceAll(maxPriceString, ",", "")
	maxPrice, err := math.LegacyNewDecFromStr(maxPriceString)
	if err != nil {
		return math.LegacyDec{}, types.ErrInvalidPriceFormat
	}

	return maxPrice, nil
}
