package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/x/dex/types"
)

func (k msgServer) AddOrder(goCtx context.Context, msg *types.MsgAddOrder) (*types.Order, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.DenomGiving == msg.DenomReceiving {
		return nil, types.ErrSameDenom
	}

	amount, err := ParseAmount(msg.Amount)
	if err != nil {
		return nil, err
	}

	if amount.LT(k.DenomKeeper.MinOrderSize(ctx, msg.DenomGiving)) {
		return nil, types.ErrOrderSizeTooSmall
	}

	if msg.TradeAmount == "" {
		msg.TradeAmount = "0"
	}

	tradeAmount, err := ParseAmount(msg.TradeAmount)
	if err != nil {
		return nil, err
	}

	if err = k.precheckTrade(ctx, msg.Creator, msg.DenomGiving, &amount, false); err != nil {
		return nil, err
	}

	maxPrice, err := stringToDec(msg.MaxPrice)
	if err != nil {
		return nil, err
	}

	if msg.Interval < 1 {
		msg.Interval = 1
	}

	if maxPrice == nil || maxPrice.IsNil() {
		return nil, types.ErrMaxPriceNotSet
	}

	if maxPrice.LTE(math.LegacyZeroDec()) {
		return nil, types.ErrNegativePrice
	}

	acc, _ := sdk.AccAddressFromBech32(msg.Creator)
	coins := sdk.NewCoins(sdk.NewCoin(msg.DenomGiving, amount))
	if err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, acc, types.PoolOrders, coins); err != nil {
		return nil, fmt.Errorf("could not send coins to module: %w", err)
	}

	order := types.Order{
		Creator:           msg.Creator,
		DenomGiving:       msg.DenomGiving,
		DenomReceiving:    msg.DenomReceiving,
		AmountLeft:        amount,
		AmountGiven:       amount,
		AmountReceived:    math.ZeroInt(),
		TradeAmount:       tradeAmount,
		MaxPrice:          *maxPrice,
		AddedAt:           ctx.BlockHeight(),
		NumBlocks:         msg.Blocks,
		ExecutionInterval: msg.Interval,
		AllowIncomplete:   msg.AllowIncomplete,
	}

	order.Index = k.SetOrder(ctx, order)

	return &order, nil
}

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
		coins := sdk.NewCoins(sdk.NewCoin(order.DenomGiving, order.AmountLeft))
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

func (k msgServer) UpdateOrder(ctx context.Context, msg *types.MsgUpdateOrder) (*types.Order, error) {
	order, found := k.GetOrder(ctx, msg.Index)
	if !found {
		return nil, types.ErrOrderNotFound
	}

	address, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, types.ErrInvalidAddress
	}

	if order.Creator != msg.Creator {
		return nil, types.ErrInvalidCreator
	}

	amount, err := ParseAmount(msg.Amount)
	if err != nil {
		return nil, fmt.Errorf("could not parse amount: %w", err)
	}

	if amount.LT(k.DenomKeeper.MinOrderSize(ctx, order.DenomGiving)) {
		return nil, types.ErrOrderSizeTooSmall
	}

	if msg.TradeAmount == "" {
		msg.TradeAmount = "0"
	}

	tradeAmount, err := ParseAmount(msg.TradeAmount)
	if err != nil {
		return nil, fmt.Errorf("could not parse trade amount: %w", err)
	}

	coins := sdk.NewCoins(sdk.NewCoin(order.DenomGiving, order.AmountLeft))
	if err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.PoolOrders, address, coins); err != nil {
		return nil, fmt.Errorf("could not send coins to address: %w", err)
	}

	if k.BankKeeper.SpendableCoin(ctx, address, order.DenomGiving).Amount.LT(amount) {
		return nil, types.ErrNotEnoughFunds
	}

	coins = sdk.NewCoins(sdk.NewCoin(order.DenomGiving, amount))
	if err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, address, types.PoolOrders, coins); err != nil {
		return nil, fmt.Errorf("could not send coins to module: %w", err)
	}

	maxPrice, err := stringToDec(msg.MaxPrice)
	if err != nil {
		return nil, fmt.Errorf("could not get max price: %w", err)
	}

	if maxPrice == nil {
		maxPrice = &order.MaxPrice
	}

	amountChange := amount.Sub(order.AmountLeft)

	order.AmountGiven = order.AmountGiven.Add(amountChange)
	order.AmountLeft = amount
	order.TradeAmount = tradeAmount
	order.MaxPrice = *maxPrice

	k.SetOrder(ctx, order)

	return &order, nil
}
