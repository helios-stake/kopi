package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	denomtypes "github.com/kopi-money/kopi/x/denominations/types"
	"github.com/kopi-money/kopi/x/strategies/types"
)

type CalculateValue []func() (math.LegacyDec, error)

func (cv CalculateValue) get() (math.LegacyDec, error) {
	value := math.LegacyZeroDec()
	for _, calcValue := range cv {
		v, err := calcValue()
		if err != nil {
			return math.LegacyDec{}, err
		}

		value = value.Add(v)
	}

	return value, nil
}

func (k Keeper) calculateNewStrategyAssetAmount(ctx context.Context, denom string, addedAmount math.Int, calculateValue CalculateValue) (math.Int, error) {
	assetSupply := k.BankKeeper.GetSupply(ctx, denom).Amount
	if assetSupply.IsZero() {
		return addedAmount, nil
	}

	aAssetValue, err := calculateValue.get()
	if err != nil {
		return math.Int{}, fmt.Errorf("could not calculate aasset value: %w", err)
	}

	valueShare := addedAmount.ToLegacyDec().Quo(aAssetValue)

	var newTokens math.Int
	if valueShare.Equal(math.LegacyOneDec()) {
		newTokens = addedAmount
	} else {
		newTokens = assetSupply.ToLegacyDec().Quo(math.LegacyOneDec().Sub(valueShare)).RoundInt().Sub(assetSupply)
	}

	return newTokens, nil
}

func (k Keeper) calculateRedemptionAmount(ctx context.Context, arbitrageDenom *denomtypes.ArbitrageDenom, requestedAAssetAmount, available math.Int, calculateValue CalculateValue, allowIncomplete bool) (math.Int, math.Int, error) {
	redemptionValue, err := k.calculateRedemptionValue(ctx, arbitrageDenom, requestedAAssetAmount, calculateValue)
	if err != nil {
		return math.Int{}, math.Int{}, fmt.Errorf("could not calculate redemption value: %w", err)
	}

	if redemptionValue.GT(available) && !allowIncomplete {
		return math.Int{}, math.Int{}, types.ErrNotEnoughVault
	}

	redeemAmount := math.MinInt(redemptionValue, available)
	requestedShare := redeemAmount.Quo(redemptionValue)

	// how much of the given cAssets have been used
	usedTokens := requestedAAssetAmount.Mul(requestedShare)
	return redeemAmount, usedTokens, nil
}

func (k Keeper) calculateRedemptionValue(ctx context.Context, arbitrageDenom *denomtypes.ArbitrageDenom, requestedAAssetAmount math.Int, calculateValue CalculateValue) (math.Int, error) {
	if requestedAAssetAmount.IsZero() {
		return math.ZeroInt(), nil
	}

	// First it is calculated how much of the total share the withdrawal request's given tokens represent.
	assetSupply := math.LegacyNewDecFromInt(k.BankKeeper.GetSupply(ctx, arbitrageDenom.DexDenom).Amount)
	assetValue, err := calculateValue.get()
	if err != nil {
		return math.Int{}, fmt.Errorf("could not calculate aAsset value: %w", err)
	}

	// how much value of all cAssetValue does the redemption request represent
	redemptionShare := requestedAAssetAmount.ToLegacyDec().Quo(assetSupply)
	redemptionValue := assetValue.Mul(redemptionShare).TruncateInt()

	msg := fmt.Sprintf("Share: %v, Total value: %v, Redemption value: %v", redemptionShare.String(), assetValue.String(), redemptionValue.String())
	k.Logger().Info(msg)

	return redemptionValue, nil
}
