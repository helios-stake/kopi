package keeper

import (
	"context"
	"cosmossdk.io/math"
	"github.com/kopi-money/kopi/cache"

	errorsmod "cosmossdk.io/errors"
	"github.com/kopi-money/kopi/x/blockspeed/types"
)

func (k msgServer) UpdateMovingAverageFactor(ctx context.Context, req *types.MsgUpdateMovingAverageFactor) (*types.Void, error) {
	err := cache.Transact(ctx, func(innerCtx context.Context) error {
		if k.GetAuthority() != req.Authority {
			return errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
		}

		movingAverageFactor, err := math.LegacyNewDecFromStr(req.MovingAverageFactor)
		if err != nil {
			return err
		}

		params := k.GetParams(innerCtx)
		params.MovingAverageFactor = movingAverageFactor

		if err = k.SetParams(innerCtx, params); err != nil {
			return err
		}

		return nil
	})

	return &types.Void{}, err
}
