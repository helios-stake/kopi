package keeper

import (
	"context"

	"cosmossdk.io/math"
	denomtypes "github.com/kopi-money/kopi/x/denominations/types"
	"github.com/kopi-money/kopi/x/mm/types"
)

// GetVaultAmount return the amount of funds held in the base denom of an CAsset. For example, when akUSD is the CAsset,
// this functions return the amount of available kUSD
func (k Keeper) GetVaultAmount(ctx context.Context, cAsset *denomtypes.CAsset) math.Int {
	address := k.AccountKeeper.GetModuleAccount(ctx, types.PoolVault).GetAddress()
	amount := k.BankKeeper.SpendableCoins(ctx, address).AmountOf(cAsset.BaseDexDenom)
	return amount
}

func (k Keeper) getCAssetSupply(ctx context.Context, cAsset *denomtypes.CAsset) math.Int {
	return k.BankKeeper.GetSupply(ctx, cAsset.DexDenom).Amount
}

// CalculateNewCAssetAmount calculates how much new c-tokens have to be minted given how much value is being added to
// the vault.
func (k Keeper) CalculateNewCAssetAmount(ctx context.Context, cAsset *denomtypes.CAsset, addedAmount math.Int) math.Int {
	cAssetSupply := k.getCAssetSupply(ctx, cAsset)
	if cAssetSupply.IsZero() {
		return addedAmount
	}

	loanSum := k.GetLoanSumWithDefault(ctx, cAsset.BaseDexDenom).LoanSum
	vaultSize := k.GetVaultAmount(ctx, cAsset).ToLegacyDec()

	cAssetValue := loanSum.Add(vaultSize)

	newTotalValue := addedAmount.ToLegacyDec().Add(cAssetValue)
	valueShare := addedAmount.ToLegacyDec().Quo(newTotalValue)

	var newTokens math.Int
	if valueShare.Equal(math.LegacyOneDec()) {
		newTokens = addedAmount
	} else {
		newTokens = cAssetSupply.ToLegacyDec().Quo(math.LegacyOneDec().Sub(valueShare)).TruncateInt().Sub(cAssetSupply)
	}

	return newTokens
}

// CalculateCAssetValue calculates the total underlying of an CAsset. This includes funds lying in the vault as well as
// funds in outstanding loans.
func (k Keeper) CalculateCAssetValue(ctx context.Context, cAsset *denomtypes.CAsset) math.LegacyDec {
	loanSum := k.GetLoanSumWithDefault(ctx, cAsset.BaseDexDenom).LoanSum
	vaultSize := k.GetVaultAmount(ctx, cAsset).ToLegacyDec()

	return vaultSize.Add(loanSum)
}

func (k Keeper) CalculateCAssetRedemptionValue(ctx context.Context, cAsset *denomtypes.CAsset) math.LegacyDec {
	supply := k.BankKeeper.GetSupply(ctx, cAsset.DexDenom)
	if supply.Amount.IsZero() {
		return math.LegacyZeroDec()
	}

	value := k.CalculateCAssetValue(ctx, cAsset)
	redemptionValue := value.Quo(supply.Amount.ToLegacyDec())
	return redemptionValue
}

// calculateCAssetPrice calculates the price of a CAsset in relation to its base denomination.
func (k Keeper) calculateCAssetPrice(ctx context.Context, cAsset *denomtypes.CAsset) math.LegacyDec {
	CAssetValue := k.CalculateCAssetValue(ctx, cAsset)
	CAssetSupply := math.LegacyNewDecFromInt(k.BankKeeper.GetSupply(ctx, cAsset.DexDenom).Amount)

	CAssetPrice := math.LegacyOneDec()
	if CAssetSupply.IsPositive() {
		CAssetPrice = CAssetValue.Quo(CAssetSupply)
	}

	return CAssetPrice
}

func (k Keeper) ConvertToBaseAmount(ctx context.Context, cAsset *denomtypes.CAsset, amountCAsset math.LegacyDec) math.LegacyDec {
	if amountCAsset.IsZero() {
		return math.LegacyZeroDec()
	}

	cAssetValue := k.CalculateCAssetValue(ctx, cAsset)
	cAssetSupply := k.getCAssetSupply(ctx, cAsset)

	return convertToBaseAmount(cAssetSupply.ToLegacyDec(), cAssetValue, amountCAsset)
}

func convertToBaseAmount(supply, value, amountCAsset math.LegacyDec) math.LegacyDec {
	if amountCAsset.IsZero() {
		return math.LegacyZeroDec()
	}

	return amountCAsset.Quo(supply).Mul(value)
}
