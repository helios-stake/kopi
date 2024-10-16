package keeper

import (
	"context"
	"github.com/kopi-money/kopi/cache"
	"github.com/kopi-money/kopi/constants"

	"cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	"github.com/kopi-money/kopi/x/denominations/types"
)

func (k msgServer) DexAddDenom(ctx context.Context, req *types.MsgDexAddDenom) (*types.MsgUpdateParamsResponse, error) {
	err := cache.Transact(ctx, func(innerCtx context.Context) error {
		if k.GetAuthority() != req.Authority {
			return errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
		}

		params := k.GetParams(innerCtx)

		dexDenom, err := createDexDenom(params.DexDenoms, req.Name, req.Factor, req.MinLiquidity, req.MinOrderSize, req.Exponent)
		if err != nil {
			return err
		}

		params.DexDenoms = append(params.DexDenoms, dexDenom)

		if err = k.SetParams(innerCtx, params); err != nil {
			return err
		}

		return nil
	})

	return &types.MsgUpdateParamsResponse{}, err
}

func (k msgServer) DexUpdateMinimumLiquidity(ctx context.Context, req *types.MsgDexUpdateMinimumLiquidity) (*types.MsgUpdateParamsResponse, error) {
	err := cache.Transact(ctx, func(innerCtx context.Context) error {
		if k.GetAuthority() != req.Authority {
			return errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
		}

		params := k.GetParams(innerCtx)
		minLiquidity, _ := math.NewIntFromString(req.MinLiquidity)
		dexDenoms := []*types.DexDenom{}
		found := false

		for _, dexDenom := range params.DexDenoms {
			if dexDenom.Name == req.Name {
				dexDenom.MinLiquidity = minLiquidity
				found = true
			}

			dexDenoms = append(dexDenoms, dexDenom)
		}

		if !found {
			return types.ErrInvalidDexAsset
		}

		params.DexDenoms = dexDenoms

		if err := k.SetParams(innerCtx, params); err != nil {
			return err
		}

		return nil
	})

	return &types.MsgUpdateParamsResponse{}, err
}

func (k msgServer) DexUpdateMinimumOrderSize(ctx context.Context, req *types.MsgDexUpdateMinimumOrderSize) (*types.MsgUpdateParamsResponse, error) {
	err := cache.Transact(ctx, func(innerCtx context.Context) error {
		if k.GetAuthority() != req.Authority {
			return errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
		}

		params := k.GetParams(innerCtx)
		minOrderSize, _ := math.NewIntFromString(req.MinOrderSize)
		dexDenoms := []*types.DexDenom{}
		found := false

		for _, dexDenom := range params.DexDenoms {
			if dexDenom.Name == req.Name {
				dexDenom.MinOrderSize = minOrderSize
				found = true
			}

			dexDenoms = append(dexDenoms, dexDenom)
		}

		if !found {
			return types.ErrInvalidDexAsset
		}

		params.DexDenoms = dexDenoms

		if err := k.SetParams(innerCtx, params); err != nil {
			return err
		}

		return nil
	})

	return &types.MsgUpdateParamsResponse{}, err
}

func createDexDenom(dexDenoms types.DexDenoms, name, factorStr, minLiquidityStr, minOrderSizeStr string, exponent uint64) (*types.DexDenom, error) {
	factor, denom, err := types.ExtractNumberAndString(factorStr)
	if err != nil {
		return nil, err
	}

	if denom != "" {
		if denom == constants.BaseCurrency {
			return nil, types.ErrInvalidFactorReference
		}

		referenceDenom := dexDenoms.Get(denom)
		if referenceDenom == nil {
			return nil, types.ErrInvalidFactorReference
		}

		factor = referenceDenom.Factor.Mul(factor)
	}

	if factor.LTE(math.LegacyZeroDec()) {
		return nil, types.ErrInvalidFactor
	}

	minLiquidity, _ := math.NewIntFromString(minLiquidityStr)
	minOrderSize, _ := math.NewIntFromString(minOrderSizeStr)

	return &types.DexDenom{
		Name:         name,
		Factor:       &factor,
		MinLiquidity: minLiquidity,
		MinOrderSize: minOrderSize,
		Exponent:     exponent,
	}, nil
}
