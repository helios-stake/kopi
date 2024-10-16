package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	"github.com/kopi-money/kopi/constants"
	"github.com/kopi-money/kopi/x/denominations/types"
)

// IsValidDenom is used to check whether a given denom is included in the parameters
func (k Keeper) IsValidDenom(ctx context.Context, denom string) bool {
	for _, dexDenom := range k.GetParams(ctx).DexDenoms {
		if dexDenom.Name == denom {
			return true
		}
	}

	return false
}

func (k Keeper) GetDexDenom(ctx context.Context, denom string) (*types.DexDenom, bool) {
	for _, dexDenom := range k.GetParams(ctx).DexDenoms {
		if dexDenom.Name == denom {
			return dexDenom, true
		}
	}

	return &types.DexDenom{}, false
}

// Denoms returns a list of all denoms
func (k Keeper) Denoms(ctx context.Context) (denoms []string) {
	for _, dexDenom := range k.GetParams(ctx).DexDenoms {
		denoms = append(denoms, dexDenom.Name)
	}

	return
}

func (k Keeper) IsKCoin(ctx context.Context, denom string) bool {
	for _, kCoin := range k.GetParams(ctx).KCoins {
		if kCoin.DexDenom == denom {
			return true
		}
	}

	return false
}

func (k Keeper) IsBorrowableDenom(ctx context.Context, denom string) bool {
	for _, cAsset := range k.GetParams(ctx).CAssets {
		if cAsset.BaseDexDenom == denom {
			return true
		}
	}

	return false
}

func (k Keeper) IsCollateralDenom(ctx context.Context, denom string) bool {
	for _, collateralDenom := range k.GetParams(ctx).CollateralDenoms {
		if collateralDenom.DexDenom == denom {
			return true
		}
	}

	return false
}

func (k Keeper) IsCAsset(ctx context.Context, denom string) bool {
	for _, cAsset := range k.GetParams(ctx).CAssets {
		if cAsset.DexDenom == denom {
			return true
		}
	}

	return false
}

func (k Keeper) IsArbitrageDenom(ctx context.Context, denom string) bool {
	for _, arbitrageDenom := range k.GetArbitrageDenoms(ctx) {
		if arbitrageDenom.DexDenom == denom {
			return true
		}
	}

	return false
}

func (k Keeper) IsNativeDenom(ctx context.Context, denom string) bool {
	if denom == constants.BaseCurrency {
		return true
	}

	if k.IsKCoin(ctx, denom) {
		return true
	}

	if k.IsCAsset(ctx, denom) {
		return true
	}

	if k.IsArbitrageDenom(ctx, denom) {
		return true
	}

	return false
}

// KCoins returns a slice containing the kCoins of all denom groups.
func (k Keeper) KCoins(ctx context.Context) (kCoins []string) {
	for _, kCoin := range k.GetParams(ctx).KCoins {
		kCoins = append(kCoins, kCoin.DexDenom)
	}

	return kCoins
}

// NonKCoins returns a slice containing the non-kCoins of all denom groups.
func (k Keeper) NonKCoins(ctx context.Context) (nonKCoins []string) {
	for _, dexDenom := range k.GetParams(ctx).DexDenoms {
		if !k.IsValidDenom(ctx, dexDenom.Name) {
			nonKCoins = append(nonKCoins, dexDenom.Name)
		}
	}

	return
}

// ReferenceDenoms returns a list of denoms that are used as price reference for a kCoin. If the kCoin
// does not exist, an empty slice is created.
func (k Keeper) ReferenceDenoms(ctx context.Context, kCoinName string) []string {
	for _, kCoin := range k.GetParams(ctx).KCoins {
		if kCoin.DexDenom == kCoinName {
			return kCoin.References
		}
	}

	return nil
}

func (k Keeper) Exponent(ctx context.Context, denom string) (uint64, error) {
	for _, dexDenom := range k.GetParams(ctx).DexDenoms {
		if dexDenom.Name == denom {
			return dexDenom.Exponent, nil
		}
	}

	return 0, fmt.Errorf("could not find gien denom: %v", denom)
}

// InitialVirtualLiquidityFactor returns the factor used for initial virtual liquidity for a denom.
func (k Keeper) InitialVirtualLiquidityFactor(ctx context.Context, denom string) (math.LegacyDec, error) {
	for _, dexDenom := range k.GetParams(ctx).DexDenoms {
		if dexDenom.Name == denom {
			return *dexDenom.Factor, nil
		}
	}

	return math.LegacyDec{}, fmt.Errorf("no initial virtual liquidity factor found for %v", denom)
}

func (k Keeper) MaxSupply(ctx context.Context, kCoinName string) math.Int {
	for _, kCoin := range k.GetParams(ctx).KCoins {
		if kCoin.DexDenom == kCoinName {
			return kCoin.MaxSupply
		}
	}

	panic(fmt.Sprintf("no max burn amount found for %v", kCoinName))
}

func (k Keeper) MaxBurnAmount(ctx context.Context, kCoinName string) math.Int {
	for _, kCoin := range k.GetParams(ctx).KCoins {
		if kCoin.DexDenom == kCoinName {
			return kCoin.MaxBurnAmount
		}
	}

	panic(fmt.Sprintf("no max burn amount found for %v", kCoinName))
}

func (k Keeper) MaxMintAmount(ctx context.Context, kCoinName string) math.Int {
	for _, kCoin := range k.GetParams(ctx).KCoins {
		if kCoin.DexDenom == kCoinName {
			return kCoin.MaxMintAmount
		}
	}

	panic(fmt.Sprintf("no max mint amount found for %v", kCoinName))
}

func (k Keeper) MinLiquidity(ctx context.Context, denom string) math.Int {
	for _, dexDenom := range k.GetParams(ctx).DexDenoms {
		if dexDenom.Name == denom {
			return dexDenom.MinLiquidity
		}
	}

	panic(fmt.Sprintf("no minimum liquidity found for %v", denom))
}

func (k Keeper) MinOrderSize(ctx context.Context, denom string) math.Int {
	for _, dexDenom := range k.GetParams(ctx).DexDenoms {
		if dexDenom.Name == denom {
			return dexDenom.MinOrderSize
		}
	}

	panic(fmt.Sprintf("no minimum order size found for %v", denom))
}

func (k Keeper) GetArbitrageDenomByCAsset(ctx context.Context, cAsset string) (*types.ArbitrageDenom, error) {
	for _, arbitrageDenom := range k.GetParams(ctx).StrategyDenoms.ArbitrageDenoms {
		if arbitrageDenom.CAsset == cAsset {
			return arbitrageDenom, nil
		}
	}

	return nil, types.ErrInvalidArbitrageDenom
}

func (k Keeper) GetArbitrageDenomByName(ctx context.Context, name string) (*types.ArbitrageDenom, error) {
	for _, arbitrageDenom := range k.GetParams(ctx).StrategyDenoms.ArbitrageDenoms {
		if arbitrageDenom.DexDenom == name {
			return arbitrageDenom, nil
		}
	}

	return nil, types.ErrInvalidArbitrageDenom
}

func (k Keeper) GetCAssets(ctx context.Context) []*types.CAsset {
	return k.GetParams(ctx).CAssets
}

func (k Keeper) GetCAsset(ctx context.Context, name string) (*types.CAsset, error) {
	cAsset, _ := k.GetCAssetByName(ctx, name)
	if cAsset != nil {
		return cAsset, nil
	}

	cAsset, _ = k.GetCAssetByBaseName(ctx, name)
	if cAsset != nil {
		return cAsset, nil
	}

	return nil, types.ErrInvalidCAsset
}

func (k Keeper) GetCAssetByBaseName(ctx context.Context, baseDenom string) (*types.CAsset, error) {
	for _, cAsset := range k.GetParams(ctx).CAssets {
		if cAsset.BaseDexDenom == baseDenom {
			return cAsset, nil
		}
	}

	return nil, types.ErrInvalidCAsset
}

func (k Keeper) GetCAssetByName(ctx context.Context, name string) (*types.CAsset, error) {
	for _, aasset := range k.GetParams(ctx).CAssets {
		if aasset.DexDenom == name {
			return aasset, nil
		}
	}

	return nil, types.ErrInvalidCAsset
}

func (k Keeper) GetCollateralDenoms(ctx context.Context) []*types.CollateralDenom {
	return k.GetParams(ctx).CollateralDenoms
}

func (k Keeper) GetCollateralDenom(ctx context.Context, denom string) *types.CollateralDenom {
	for _, collateralDenom := range k.GetParams(ctx).CollateralDenoms {
		if collateralDenom.DexDenom == denom {
			return collateralDenom
		}
	}

	return nil
}

func (k Keeper) GetDepositCap(ctx context.Context, denom string) (math.Int, error) {
	for _, collateralDenom := range k.GetParams(ctx).CollateralDenoms {
		if collateralDenom.DexDenom == denom {
			return collateralDenom.MaxDeposit, nil
		}
	}

	return math.Int{}, types.ErrInvalidCollateralDenom
}

func (k Keeper) GetLTV(ctx context.Context, denom string) (math.LegacyDec, error) {
	for _, collateralDenom := range k.GetParams(ctx).CollateralDenoms {
		if collateralDenom.DexDenom == denom {
			return collateralDenom.Ltv, nil
		}
	}

	return math.LegacyDec{}, types.ErrInvalidCollateralDenom
}

func (k Keeper) IsValidCollateralDenom(ctx context.Context, denom string) bool {
	for _, depositDenom := range k.GetParams(ctx).CollateralDenoms {
		if depositDenom.DexDenom == denom {
			return true
		}
	}

	return false
}

func (k Keeper) GetArbitrageDenoms(ctx context.Context) []*types.ArbitrageDenom {
	strategyDenoms := k.GetParams(ctx).StrategyDenoms
	if strategyDenoms == nil {
		return nil
	}

	return strategyDenoms.ArbitrageDenoms
}

func (k Keeper) ConvertToExponent(ctx context.Context, denom string, amount math.LegacyDec, targetExponent uint64) (math.LegacyDec, error) {
	sourceExponent, err := k.Exponent(ctx, denom)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return ConvertToExponent(amount, sourceExponent, targetExponent), nil
}

func ConvertToExponent(amount math.LegacyDec, sourceExponent, targetExponent uint64) math.LegacyDec {
	switch {
	case sourceExponent > targetExponent:
		factor := math.LegacyNewDec(10).Power(sourceExponent - targetExponent)
		return amount.Quo(factor)
	case sourceExponent < targetExponent:
		factor := math.LegacyNewDec(10).Power(targetExponent - sourceExponent)
		return amount.Mul(factor)
	default:
		return amount
	}
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}

	return false
}
