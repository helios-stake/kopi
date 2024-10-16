package keeper

import (
	"context"
	"cosmossdk.io/math"

	"github.com/kopi-money/kopi/x/swap/types"
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

func (k Keeper) mintThreshold(ctx context.Context) math.LegacyDec {
	mintThreshold := k.GetParams(ctx).MintThreshold
	if !mintThreshold.IsNil() {
		return mintThreshold
	}

	return types.MintThreshold
}

func (k Keeper) burnThreshold(ctx context.Context) math.LegacyDec {
	burnThreshold := k.GetParams(ctx).BurnThreshold
	if !burnThreshold.IsNil() {
		return burnThreshold
	}

	return types.BurnThreshold
}
