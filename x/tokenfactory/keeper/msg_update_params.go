package keeper

import (
	"context"
	"fmt"

	"github.com/kopi-money/kopi/cache"

	"cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	"github.com/kopi-money/kopi/x/tokenfactory/types"
)

func (k msgServer) UpdateFeeAmount(ctx context.Context, req *types.MsgUpdateFeeAmount) (*types.Void, error) {
	err := cache.Transact(ctx, func(innerCtx context.Context) error {
		if k.GetAuthority() != req.Authority {
			return errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
		}

		feeAmount, ok := math.NewIntFromString(req.FeeAmount)
		if !ok {
			return fmt.Errorf("invalid amount")
		}

		params := k.GetParams(ctx)
		params.CreationFee = feeAmount

		if err := params.Validate(); err != nil {
			return err
		}

		if err := k.SetParams(ctx, params); err != nil {
			return err
		}

		return nil
	})

	return &types.Void{}, err
}
