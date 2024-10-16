package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/kopi-money/kopi/constants"
	denomtypes "github.com/kopi-money/kopi/x/denominations/types"
)

func (k Keeper) CalculateBorrowableAmount(ctx context.Context, address, borrowDenom string) (math.LegacyDec, error) {
	collateralBaseValue, err := k.calculateCollateralBaseValue(ctx, address)
	if err != nil {
		return math.LegacyDec{}, err
	}

	loanBaseValue, err := k.calculateLoanBaseValue(ctx, address)
	if err != nil {
		return math.LegacyDec{}, err
	}

	borrowableBaseValue := collateralBaseValue.Sub(loanBaseValue)
	borrowableBaseValue = math.LegacyMaxDec(math.LegacyZeroDec(), borrowableBaseValue)

	borrowableValue, err := k.DexKeeper.GetValueIn(ctx, constants.BaseCurrency, borrowDenom, borrowableBaseValue)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return borrowableValue, nil
}

func (k Keeper) calculateCollateralBaseValue(ctx context.Context, address string) (math.LegacyDec, error) {
	borrowableAmount := math.LegacyZeroDec()
	for _, collateral := range k.DenomKeeper.GetCollateralDenoms(ctx) {
		amount, err := k.calculateCollateralValueForDenom(ctx, collateral, address)
		if err != nil {
			return math.LegacyDec{}, err
		}

		borrowableAmount = borrowableAmount.Add(amount)
	}

	return borrowableAmount, nil
}

func (k Keeper) calculateCollateralValueForDenom(ctx context.Context, collateralDenom *denomtypes.CollateralDenom, address string) (math.LegacyDec, error) {
	collateral, found := k.collateral.Get(ctx, collateralDenom.DexDenom, address)
	if !found {
		return math.LegacyZeroDec(), nil
	}

	price, err := k.DexKeeper.CalculatePrice(ctx, collateralDenom.DexDenom, constants.BaseCurrency)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return collateral.Amount.ToLegacyDec().Quo(price).Mul(collateralDenom.Ltv), nil
}

func (k Keeper) calculateLoanBaseValue(ctx context.Context, address string) (math.LegacyDec, error) {
	loanSum := math.LegacyZeroDec()

	for _, cAsset := range k.DenomKeeper.GetCAssets(ctx) {
		loanValue := k.GetLoanValue(ctx, cAsset.BaseDexDenom, address)

		loanValueBase, err := k.DexKeeper.GetValueInBase(ctx, cAsset.BaseDexDenom, loanValue)
		if err != nil {
			return math.LegacyDec{}, err
		}

		loanSum = loanSum.Add(loanValueBase)
	}

	return loanSum, nil
}
