package keeper

import (
	"context"

	"github.com/kopi-money/kopi/x/strategies/types"
)

func (k Keeper) GetParams(ctx context.Context) types.Params {
	params, has := k.params.Get(ctx)
	if !has {
		return types.DefaultParams()
	}

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

func (k Keeper) GetAutomationFeeCondition(ctx context.Context) uint64 {
	return k.GetParams(ctx).AutomationFeeCondition
}

func (k Keeper) GetAutomationFeeAction(ctx context.Context) uint64 {
	return k.GetParams(ctx).AutomationFeeAction
}
