package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/kopi-money/kopi/x/mm/types"
)

func (k Keeper) GetVaultValues(ctx context.Context, _ *types.GetVaultValuesQuery) (*types.GetVaultValuesResponse, error) {
	var vaults []*types.Vault

	for _, cAsset := range k.DenomKeeper.GetCAssets(ctx) {
		vault := types.Vault{
			Denom:   cAsset.BaseDexDenom,
			Balance: k.getBalance(ctx, cAsset.BaseDexDenom).String(),
			LoanSum: k.GetLoanSumWithDefault(ctx, cAsset.BaseDexDenom).LoanSum.String(),
			Supply:  k.BankKeeper.GetSupply(ctx, cAsset.DexDenom).Amount.String(),
		}

		vaults = append(vaults, &vault)
	}

	return &types.GetVaultValuesResponse{
		Vaults: vaults,
	}, nil
}

func (k Keeper) getBalance(ctx context.Context, denom string) math.Int {
	address := k.AccountKeeper.GetModuleAccount(ctx, types.PoolVault).GetAddress()
	coins := k.BankKeeper.SpendableCoins(ctx, address)

	for _, coin := range coins {
		if coin.Denom == denom {
			return coin.Amount
		}
	}

	return math.ZeroInt()
}
