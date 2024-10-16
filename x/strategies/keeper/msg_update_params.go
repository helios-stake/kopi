package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/kopi-money/kopi/cache"
	"github.com/kopi-money/kopi/x/strategies/types"
)

func (k msgServer) UpdateAutomationsCosts(ctx context.Context, req *types.MsgUpdateAutomationsCosts) (*types.Void, error) {
	err := cache.Transact(ctx, func(innerCtx context.Context) error {
		if k.GetAuthority() != req.Authority {
			return errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
		}

		params := k.GetParams(innerCtx)
		params.AutomationFeeCondition = req.ConditionFee
		params.AutomationFeeAction = req.ActionFee

		if err := k.SetParams(innerCtx, params); err != nil {
			return err
		}

		return nil
	})

	return &types.Void{}, err
}
