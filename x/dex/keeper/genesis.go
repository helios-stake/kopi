package keeper

import (
	"context"

	"github.com/kopi-money/kopi/x/dex/types"
)

func (k Keeper) ExportGenesis(ctx context.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.LiquidityList = k.GetAllLiquidity(ctx)
	genesis.LiquidityNextIndex, _ = k.liquidityEntriesNextIndex.Get(ctx)
	genesis.OrderNextIndex = k.GetOrderNextIndex(ctx)

	orderIterator := k.orders.Iterator(ctx, nil)
	for orderIterator.Valid() {
		genesis.OrderList = append(genesis.OrderList, orderIterator.GetNext())
	}

	return genesis
}

func (k Keeper) ExportGenesisBytes(ctx context.Context) []byte {
	return k.cdc.MustMarshal(k.ExportGenesis(ctx))
}
