package keeper

import (
	"context"
	"cosmossdk.io/math"
	"github.com/kopi-money/kopi/x/blockspeed/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx context.Context) types.Params {
	params, _ := k.params.Get(ctx)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	k.params.Set(ctx, params)
	return nil
}

func (k Keeper) movingAverageFactor(ctx context.Context) math.LegacyDec {
	fac := k.GetParams(ctx).MovingAverageFactor
	if !fac.IsNil() {
		return fac
	}

	return math.LegacyNewDecWithPrec(999, 3)
}
