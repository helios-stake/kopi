package keeper

import (
	"context"
	"fmt"
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

		dexDenom, ratio, err := k.createDexDenom(ctx, req.Name, req.Factor, req.MinLiquidity, req.MinOrderSize, req.Exponent)
		if err != nil {
			return err
		}

		params.DexDenoms = append(params.DexDenoms, &dexDenom)

		if err = k.SetParams(innerCtx, params); err != nil {
			return err
		}

		k.ratios.Set(innerCtx, req.Name, ratio)

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

func (k Keeper) createDexDenom(ctx context.Context, name, factorStr, minLiquidityStr, minOrderSizeStr string, exponent uint64) (types.DexDenom, types.Ratio, error) {
	referenceFactor, referenceDenom, err := types.ExtractNumberAndString(factorStr)
	if err != nil {
		return types.DexDenom{}, types.Ratio{}, err
	}

	if referenceDenom != constants.BaseCurrency && referenceDenom != "" {
		var otherRatio types.Ratio
		otherRatio, err = k.GetRatio(ctx, referenceDenom)
		if err != nil {
			return types.DexDenom{}, types.Ratio{}, fmt.Errorf("unable to find ratio for %s: %w", referenceDenom, err)
		}

		referenceFactor = otherRatio.Ratio.Quo(referenceFactor)
	} else {
		referenceDenom = constants.BaseCurrency
	}

	otherDenom, has := k.GetDexDenom(ctx, referenceDenom)
	if !has {
		return types.DexDenom{}, types.Ratio{}, fmt.Errorf("unable to find other denom: %v", referenceDenom)
	}

	referenceFactor = adjustForExponent(referenceFactor, otherDenom.Exponent, exponent)

	if referenceFactor.LTE(math.LegacyZeroDec()) {
		return types.DexDenom{}, types.Ratio{}, types.ErrInvalidFactor
	}

	minLiquidity, _ := math.NewIntFromString(minLiquidityStr)
	minOrderSize, _ := math.NewIntFromString(minOrderSizeStr)

	dexDenom := types.DexDenom{
		Name:         name,
		MinLiquidity: minLiquidity,
		MinOrderSize: minOrderSize,
		Exponent:     exponent,
	}

	ratio := types.Ratio{
		Denom: name,
		Ratio: referenceFactor,
	}

	return dexDenom, ratio, nil
}

func adjustForExponent(value math.LegacyDec, exp1, exp2 uint64) math.LegacyDec {
	if exp1 == exp2 {
		return value
	}

	value = value.Mul(math.LegacyNewDec(10).Power(exp2))
	value = value.Quo(math.LegacyNewDec(10).Power(exp1))
	return value
}
