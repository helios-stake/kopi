package keeper

import (
	"cosmossdk.io/math"
	"github.com/kopi-money/kopi/x/dex/types"
)

type LiquidityProvider struct {
	index          uint64
	address        string
	amountProvided math.LegacyDec
	shareProvided  math.LegacyDec
}

type LiquidityProviders []*LiquidityProvider

func (lps *LiquidityProviders) sumProvided() math.LegacyDec {
	sum := math.LegacyZeroDec()
	for _, lp := range *lps {
		sum = sum.Add(lp.amountProvided)
	}

	return sum
}

func (lps *LiquidityProviders) provided() *LiquidityProviders {
	sumProvided := lps.sumProvided()
	for _, lp := range *lps {
		if sumProvided.GT(math.LegacyZeroDec()) {
			lp.shareProvided = lp.amountProvided.Quo(sumProvided)
		} else {
			lp.shareProvided = math.LegacyZeroDec()
		}
	}

	return lps
}

func (k Keeper) determineLiquidityProviders(ctx types.TradeStepContext, amountToReceiveLeft math.Int, denomTo string) (*LiquidityProviders, math.Int, error) {
	var (
		liquidityProviders LiquidityProviders
		liquidityUsed      math.Int
		deleteIndexes      []int
		sumUsed            = math.ZeroInt()
	)

	// Iterate over the existing liquidity entries for this currency
	liquidityList := ctx.OrdersCaches.LiquidityMap.Get(denomTo)
	for index, liq := range liquidityList {
		if amountToReceiveLeft.LTE(math.ZeroInt()) {
			break
		}

		if amountToReceiveLeft.LT(liq.Amount) {
			// the current liquidity entry will not be fully used
			liquidityUsed = amountToReceiveLeft
			amountToReceiveLeft = math.ZeroInt()
		} else {
			// the current liquidity entry will be fully used
			liquidityUsed = liq.Amount
			amountToReceiveLeft = amountToReceiveLeft.Sub(liq.Amount)
		}

		lp := LiquidityProvider{index: liq.Index, address: liq.Address, amountProvided: liquidityUsed.ToLegacyDec()}
		liquidityProviders = append(liquidityProviders, &lp)
		sumUsed = sumUsed.Add(liquidityUsed)
		liq.Amount = liq.Amount.Sub(liquidityUsed)

		if liq.Amount.IsZero() {
			k.RemoveLiquidity(ctx.TradeContext.Context, denomTo, liq.Index)
			deleteIndexes = append(deleteIndexes, index)
		} else {
			k.SetLiquidity(ctx.TradeContext.Context, liq)
		}
	}

	ctx.OrdersCaches.LiquidityPool.Get().Sub(denomTo, sumUsed)
	liquidityList = removeIndexes(liquidityList, deleteIndexes)
	ctx.OrdersCaches.LiquidityMap.Set(denomTo, liquidityList)
	ctx.TradeBalances.AddTransfer(
		ctx.OrdersCaches.AccPoolLiquidity.Get().String(),
		ctx.OrdersCaches.AccPoolTrade.Get().String(),
		denomTo, sumUsed,
	)

	return &liquidityProviders, amountToReceiveLeft, nil
}

func removeIndexes(liquidityList []types.Liquidity, indexes []int) []types.Liquidity {
	for len(indexes) > 0 {
		index := indexes[len(indexes)-1]
		indexes = indexes[:len(indexes)-1]
		liquidityList = append(liquidityList[:index], liquidityList[index+1:]...)
	}

	return liquidityList
}

func (k Keeper) distributeFeeForLiquidityProviders(ctx types.TradeStepContext, liquidityProviders *LiquidityProviders, feeForLiquidityProvidersLeft math.Int, denom string) error {
	providerFee := k.getProviderFee(ctx.TradeContext.Context)
	liquidityEntries := ctx.TradeContext.OrdersCaches.LiquidityMap.Get(denom)

	liquidityProviderIndex := 0
	sendSum := math.ZeroInt()

	for feeForLiquidityProvidersLeft.GT(math.ZeroInt()) {
		liquidityProvider := (*liquidityProviders)[liquidityProviderIndex]
		liquidityProviderIndex += 1

		amount := math.MinInt(feeForLiquidityProvidersLeft, liquidityProvider.amountProvided.RoundInt())
		sendSum = sendSum.Add(amount)
		feeForLiquidityProvidersLeft = feeForLiquidityProvidersLeft.Sub(amount)
		liquidityProvider.amountProvided = liquidityProvider.amountProvided.Mul(providerFee)
		liquidityEntries, _ = k.addLiquidity(ctx.TradeContext.Context, denom, liquidityProvider.address, amount, liquidityEntries)
	}

	ctx.OrdersCaches.LiquidityPool.Get().Add(denom, sendSum)
	ctx.OrdersCaches.LiquidityMap.Set(denom, liquidityEntries)

	ctx.TradeContext.TradeBalances.AddTransfer(
		ctx.OrdersCaches.AccPoolTrade.Get().String(),
		ctx.OrdersCaches.AccPoolLiquidity.Get().String(),
		denom, sendSum,
	)

	return nil
}

func (k Keeper) distributeGivenFunds(ctx types.TradeStepContext, ordersCaches *types.OrdersCaches, liquidityProviders *LiquidityProviders, fundsToDistribute math.Int, denom string) error {
	liquidityEntries := ordersCaches.LiquidityMap.Get(denom)
	provided := liquidityProviders.provided()

	for _, liquidityProvider := range *provided {
		eligable := liquidityProvider.shareProvided.Mul(fundsToDistribute.ToLegacyDec()).RoundInt()
		liquidityEntries, _ = k.addLiquidity(ctx.TradeContext.Context, denom, liquidityProvider.address, eligable, liquidityEntries)
	}

	ordersCaches.LiquidityPool.Get().Add(denom, fundsToDistribute)
	ctx.TradeBalances.AddTransfer(
		ordersCaches.AccPoolTrade.Get().String(),
		ordersCaches.AccPoolLiquidity.Get().String(),
		denom, fundsToDistribute,
	)

	return nil
}
