package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/kopi-money/kopi/cache"

	errorsmod "cosmossdk.io/errors"
	"github.com/kopi-money/kopi/x/reserve/types"
)

func (k msgServer) UpdateKCoinBurnShare(goCtx context.Context, req *types.MsgUpdateKCoinBurnShare) (*types.Void, error) {
	err := cache.Transact(goCtx, func(ctx context.Context) error {
		if k.GetAuthority() != req.Authority {
			return errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
		}

		kCoinBurnShare, err := math.LegacyNewDecFromStr(req.KcoinBurnShare)
		if err != nil {
			return err
		}

		params := k.GetParams(ctx)
		params.KcoinBurnShare = kCoinBurnShare

		if err = k.SetParams(ctx, params); err != nil {
			return err
		}

		return nil
	})

	return &types.Void{}, err
}
