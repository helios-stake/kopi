package keeper

import (
	"context"

	"github.com/kopi-money/kopi/x/swap/types"
)

// Clean burns coins that have been created by the module but have not been used. It can happen that
func (k Keeper) Clean(ctx context.Context) error {
	moduleAcc := k.AccountKeeper.GetModuleAccount(ctx, types.ModuleName)
	coins := k.BankKeeper.SpendableCoins(ctx, moduleAcc.GetAddress())
	return k.BankKeeper.BurnCoins(ctx, types.ModuleName, coins)
}
