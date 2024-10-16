package keeper

import (
	"context"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/constants"
	"github.com/kopi-money/kopi/x/blockspeed/types"
)

func (k Keeper) AdjustBlockspeed(ctx context.Context) {
	height := sdk.UnwrapSDKContext(ctx).BlockHeight()
	now := sdk.UnwrapSDKContext(ctx).BlockTime()

	timestamp := now.UnixMilli()
	blockspeed := k.GetBlockspeed(ctx)

	// There is some delay between the genesis creation and the first blocks. To ignore that delay, we use hardcoded
	// values for the first few blocks.
	if blockspeed.PreviousTimestamp == 0 || height <= 10 {
		blockspeed.AverageTime = math.LegacyNewDec(1000) // 1000ms = 1s
	} else {
		diff := timestamp - blockspeed.PreviousTimestamp
		blockspeed.AverageTime = k.calcAverageTime(ctx, blockspeed.AverageTime, diff)
	}

	blockspeed.PreviousTimestamp = timestamp
	k.blockspeed.Set(ctx, blockspeed)
}

func (k Keeper) GetBlockspeed(ctx context.Context) types.Blockspeed {
	blockspeed, _ := k.blockspeed.Get(ctx)
	return blockspeed
}

func (k Keeper) SetBlockspeed(ctx context.Context, blockspeed types.Blockspeed) {
	k.blockspeed.Set(ctx, blockspeed)
}

func (k Keeper) calcAverageTime(ctx context.Context, averageTime math.LegacyDec, timeDiff int64) math.LegacyDec {
	movingAverageFactor := k.movingAverageFactor(ctx)
	averageTime = averageTime.Mul(movingAverageFactor)
	toAdd := math.LegacyOneDec().Sub(movingAverageFactor).Mul(math.LegacyNewDec(timeDiff))
	return averageTime.Add(toAdd)
}

func (k Keeper) GetSecondsPerBlock(ctx context.Context) math.LegacyDec {
	blockspeed := k.GetBlockspeed(ctx)
	return blockspeed.AverageTime.Quo(math.LegacyNewDec(1000))
}

func (k Keeper) GetBlocksPerSecond(ctx context.Context) math.LegacyDec {
	return math.LegacyOneDec().Quo(k.GetSecondsPerBlock(ctx))
}

func (k Keeper) BlocksPerYear(ctx context.Context) (math.LegacyDec, error) {
	secondsPerYear := math.LegacyNewDec(constants.SecondsPerYear)
	blockPerSecond := k.GetBlocksPerSecond(ctx)

	if blockPerSecond.IsZero() {
		return math.LegacyDec{}, types.ErrDivisionByZero
	}

	return secondsPerYear.Quo(blockPerSecond), nil
}
