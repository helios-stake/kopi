package keeper

import (
	"context"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/x/dex/types"
)

func (k msgServer) AddOrder(goCtx context.Context, msg *types.MsgAddOrder) (*types.Order, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.DenomFrom == msg.DenomTo {
		return nil, types.ErrSameDenom
	}

	amount, err := parseAmount(msg.Amount)
	if err != nil {
		return nil, err
	}

	if amount.LT(k.DenomKeeper.MinOrderSize(ctx, msg.DenomFrom)) {
		return nil, types.ErrOrderSizeTooSmall
	}

	if msg.TradeAmount == "" {
		msg.TradeAmount = "0"
	}

	tradeAmount, err := parseAmount(msg.TradeAmount)
	if err != nil {
		return nil, err
	}

	address, err := k.validateMsg(ctx, msg.Creator, msg.DenomFrom, amount)
	if err != nil {
		return nil, err
	}

	maxPrice, err := getMaxPrice(msg.MaxPrice)
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

	coins := sdk.NewCoins(sdk.NewCoin(msg.DenomFrom, amount))
	if err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, address, types.PoolOrders, coins); err != nil {
		return nil, errors.Wrap(err, "could not send coins to module")
	}

	order := types.Order{
		Creator:           msg.Creator,
		DenomFrom:         msg.DenomFrom,
		DenomTo:           msg.DenomTo,
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
