package keeper

import (
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/constants"
	"github.com/kopi-money/kopi/x/dex/types"
	"golang.org/x/net/context"
)

func (k Keeper) ResetTradeFeeTracker(ctx context.Context) {
	k.tradeFeeTracker.Set(ctx, 0)
}

func (k Keeper) addTradeFee(ctx context.Context, denom string, amount math.LegacyDec) {
	if denom != constants.BaseCurrency {
		amount, _ = k.GetValueInBase(ctx, denom, amount)
	}

	tracked, _ := k.tradeFeeTracker.Get(ctx)
	tracked += amount.TruncateInt().Int64()
	k.tradeFeeTracker.Set(ctx, tracked)
}

func (k Keeper) EmitTradeFeeEvent(ctx context.Context) {
	feeAmount, _ := k.tradeFeeTracker.Get(ctx)
	liquidityAmount := k.getLiquidityAmountInBase(ctx)

	if feeAmount > 0 {
		sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(
			sdk.NewEvent("block_trade_fee",
				sdk.Attribute{Key: "fee_amount", Value: strconv.Itoa(int(feeAmount))},
				sdk.Attribute{Key: "liquidity_amount", Value: strconv.Itoa(int(liquidityAmount.Int64()))},
			),
		)
	}
}

func (k Keeper) getLiquidityAmountInBase(ctx context.Context) math.Int {
	sum := math.LegacyZeroDec()

	poolAcc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolLiquidity)
	balance := k.BankKeeper.SpendableCoins(ctx, poolAcc.GetAddress())
	for _, coin := range balance {
		baseAmount, _ := k.GetValueInBase(ctx, coin.Denom, coin.Amount.ToLegacyDec())
		sum = sum.Add(baseAmount)
	}

	return sum.TruncateInt()
}
