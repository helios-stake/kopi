package keeper

import (
	"context"
	"github.com/kopi-money/kopi/cache"

	"cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	"github.com/kopi-money/kopi/x/denominations/types"
)

func (k msgServer) CollateralAddDenom(ctx context.Context, req *types.MsgCollateralAddDenom) (*types.MsgUpdateParamsResponse, error) {
	err := cache.Transact(ctx, func(innerCtx context.Context) error {
		if k.GetAuthority() != req.Authority {
			return errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
		}

		params := k.GetParams(innerCtx)
		ltv, _ := math.LegacyNewDecFromStr(req.Ltv)
		maxDeposit, _ := math.NewIntFromString(req.MaxDeposit)

		params.CollateralDenoms = append(params.CollateralDenoms, &types.CollateralDenom{
			DexDenom:   req.Denom,
			Ltv:        ltv,
			MaxDeposit: maxDeposit,
		})

		if err := k.SetParams(innerCtx, params); err != nil {
			return err
		}
		return nil
	})

	return &types.MsgUpdateParamsResponse{}, err
}

func (k msgServer) CollateralUpdateLTV(ctx context.Context, req *types.MsgCollateralUpdateLTV) (*types.MsgUpdateParamsResponse, error) {
	err := cache.Transact(ctx, func(innerCtx context.Context) error {
		if k.GetAuthority() != req.Authority {
			return errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
		}

		params := k.GetParams(innerCtx)
		ltv, _ := math.LegacyNewDecFromStr(req.Ltv)
		collateralDenoms := []*types.CollateralDenom{}
		found := false

		for _, collateralDenom := range params.CollateralDenoms {
			if collateralDenom.DexDenom == req.Denom {
				collateralDenom.Ltv = ltv
				found = true
			}

			collateralDenoms = append(collateralDenoms, collateralDenom)
		}

		if !found {
			return types.ErrInvalidCollateralDenom
		}

		params.CollateralDenoms = collateralDenoms

		if err := k.SetParams(innerCtx, params); err != nil {
			return err
		}

		return nil
	})

	return &types.MsgUpdateParamsResponse{}, err
}

func (k msgServer) CollateralUpdateDepositLimit(ctx context.Context, req *types.MsgCollateralUpdateDepositLimit) (*types.MsgUpdateParamsResponse, error) {
	err := cache.Transact(ctx, func(innerCtx context.Context) error {
		if k.GetAuthority() != req.Authority {
			return errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
		}

		params := k.GetParams(innerCtx)

		maxDeposit, _ := math.NewIntFromString(req.MaxDeposit)
		collateralDenoms := []*types.CollateralDenom{}
		found := false

		for _, collateralDenom := range params.CollateralDenoms {
			if collateralDenom.DexDenom == req.Denom {
				collateralDenom.MaxDeposit = maxDeposit
				found = true
			}

			collateralDenoms = append(collateralDenoms, collateralDenom)
		}

		if !found {
			return types.ErrInvalidCollateralDenom
		}

		params.CollateralDenoms = collateralDenoms

		if err := k.SetParams(innerCtx, params); err != nil {
			return err
		}

		return nil
	})

	return &types.MsgUpdateParamsResponse{}, err
}
