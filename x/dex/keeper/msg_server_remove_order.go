package keeper

import (
	"context"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/x/dex/types"
)

func (k msgServer) RemoveOrder(goCtx context.Context, msg *types.MsgRemoveOrder) (*types.Void, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	order, found := k.GetOrder(ctx, msg.Index)
	if !found {
		return nil, types.ErrItemNotFound
	}

	if order.Creator != msg.Creator {
		return nil, types.ErrInvalidCreator
	}

	if !order.AmountLeft.IsNil() && order.AmountLeft.GT(math.ZeroInt()) {
		coins := sdk.NewCoins(sdk.NewCoin(order.DenomFrom, order.AmountLeft))
		address, _ := sdk.AccAddressFromBech32(order.Creator)
		if err := k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.PoolOrders, address, coins); err != nil {
			return nil, err
		}
	}

	k.Keeper.RemoveOrder(ctx, order)

	return &types.Void{}, nil
}

func (k msgServer) RemoveOrders(goCtx context.Context, msg *types.MsgRemoveOrders) (*types.Void, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	for _, order := range k.GetAllOrdersByAddress(ctx, msg.Creator) {
		k.Keeper.RemoveOrder(ctx, order)
	}

	return &types.Void{}, nil
}
