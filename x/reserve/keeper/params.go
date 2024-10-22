package keeper

import (
	"context"

	"cosmossdk.io/math"

	"github.com/kopi-money/kopi/x/reserve/types"
)

func (k Keeper) GetParams(ctx context.Context) types.Params {
	params, has := k.params.Get(ctx)
	if !has {
		return types.DefaultParams()
	}

	return params
}

func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	k.params.Set(ctx, params)
	return nil
}

func (k Keeper) getKCoinBurnShare(ctx context.Context) math.LegacyDec {
	kCoinBurnShare := k.GetParams(ctx).KcoinBurnShare
	if kCoinBurnShare.IsNil() {
		kCoinBurnShare = math.LegacyOneDec()
	}

	return kCoinBurnShare
}
