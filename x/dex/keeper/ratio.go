package keeper

import (
	"context"
	"cosmossdk.io/math"
	"fmt"
	"github.com/kopi-money/kopi/constants"

	"github.com/kopi-money/kopi/x/dex/types"
)

func (k Keeper) SetRatio(ctx context.Context, ratio types.Ratio) {
	k.ratios.Set(ctx, ratio.Denom, ratio)
}

func (k Keeper) RemoveRatio(ctx context.Context, ratio types.Ratio) {
	k.ratios.Remove(ctx, ratio.Denom)
}

func (k Keeper) GetRatio(ctx context.Context, denom string) (types.Ratio, error) {
	ratio, has := k.ratios.Get(ctx, denom)
	if !has {
		referenceFactor, err := k.DenomKeeper.InitialVirtualLiquidityFactor(ctx, denom)
		if err != nil {
			return types.Ratio{}, err
		}

		var factor math.LegacyDec
		if referenceFactor.Denom == constants.BaseCurrency {
			factor = referenceFactor.Factor
		} else {
			var otherRatio types.Ratio
			otherRatio, err = k.GetRatio(ctx, referenceFactor.Denom)
			if err != nil {
				return otherRatio, fmt.Errorf("unable to find ratio for %s: %w", referenceFactor.Denom, err)
			}

			factor = otherRatio.Ratio.Quo(referenceFactor.Factor)
		}

		ratio = types.Ratio{
			Denom: denom,
			Ratio: factor,
		}
	}

	return ratio, nil
}

func (k Keeper) GetAllRatio(ctx context.Context) (list []types.Ratio) {
	return k.ratios.Iterator(ctx, nil).GetAll()
}
