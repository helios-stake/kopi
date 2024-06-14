package keeper

import (
	"context"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/x/dex/types"
	"github.com/pkg/errors"
	"strconv"
)

func (k Keeper) ExecuteOrders(ctx context.Context, eventManager sdk.EventManagerI, blockHeight int64) error {
	ordersCaches := k.NewOrdersCaches(ctx)
	fee := k.GetTradeFee(ctx)
	maxOrderLife := int64(k.GetParams(ctx).MaxOrderLife)
	iterator := k.OrderIterator(ctx)
	tradeBalances := NewTradeBalances()

	numTrades := 0
	tradeVolumeBaseSum := math.ZeroInt()

	// At this point we know that there are no changes in the ongoing transaction. To avoid the costly iteration over
	// two lists (of which the second is empty but needs to be checked at every step), we just get all the items from
	// cache directly
	orders := iterator.GetAllFromCache()

	for _, keyValue := range orders {
		order := keyValue.Value().Value()
		if order == nil {
			k.logger.Info("order is nil, should not happen")
			continue
		}

		// First we check whether the order is expired. If yes, it is removed.
		blockEnd := k.calculateBlockEnd(maxOrderLife, blockHeight, int64(order.NumBlocks))
		if blockHeight > blockEnd {
			if !order.AmountLeft.IsNil() && order.AmountLeft.GT(math.ZeroInt()) {
				tradeBalances.AddTransfer(ordersCaches.AccPoolOrders.Get().String(), order.Creator, order.DenomFrom, order.AmountLeft)
			}

			eventManager.EmitEvent(
				sdk.NewEvent("order_expired",
					sdk.Attribute{Key: "index", Value: strconv.Itoa(int(order.Index))},
					sdk.Attribute{Key: "address", Value: order.Creator},
					sdk.Attribute{Key: "denom_from", Value: order.DenomFrom},
					sdk.Attribute{Key: "denom_to", Value: order.DenomTo},
					sdk.Attribute{Key: "amount_given", Value: order.AmountGiven.String()},
					sdk.Attribute{Key: "amount_used", Value: order.AmountGiven.Sub(order.AmountLeft).String()},
					sdk.Attribute{Key: "amount_received", Value: order.AmountReceived.String()},
					sdk.Attribute{Key: "max_price", Value: order.MaxPrice.String()},
				),
			)

			k.RemoveOrder(ctx, *order)
			continue
		}

		// Next we check whether the order is to be executed at this height
		if (order.AddedAt+blockHeight)%int64(order.ExecutionInterval) != 0 {
			continue
		}

		// Next we do the actual execution
		tradeVolumeBase, remove, err := k.executeOrder(ctx, ordersCaches, tradeBalances, fee, order)
		if err != nil {
			return errors.Wrap(err, "error executing order")
		}

		if tradeVolumeBase.GT(math.ZeroInt()) {
			tradeVolumeBaseSum = tradeVolumeBaseSum.Add(tradeVolumeBase)
			numTrades++
		}

		if remove {
			eventManager.EmitEvent(
				sdk.NewEvent("order_completed",
					sdk.Attribute{Key: "index", Value: strconv.Itoa(int(order.Index))},
					sdk.Attribute{Key: "address", Value: order.Creator},
					sdk.Attribute{Key: "denom_from", Value: order.DenomFrom},
					sdk.Attribute{Key: "denom_to", Value: order.DenomTo},
					sdk.Attribute{Key: "amount_given", Value: order.AmountGiven.String()},
					sdk.Attribute{Key: "amount_used", Value: order.AmountGiven.Sub(order.AmountLeft).String()},
					sdk.Attribute{Key: "amount_received", Value: order.AmountReceived.String()},
					sdk.Attribute{Key: "max_price", Value: order.MaxPrice.String()},
				),
			)

			k.RemoveOrder(ctx, *order)
		}
	}

	if numTrades > 0 {
		eventManager.EmitEvent(
			sdk.NewEvent("orders_executed",
				sdk.Attribute{Key: "num_trades", Value: strconv.Itoa(numTrades)},
				sdk.Attribute{Key: "amount_intermediate_base_currency", Value: tradeVolumeBaseSum.String()},
			),
		)
	}

	if err := tradeBalances.Settle(ctx, k.BankKeeper); err != nil {
		return errors.Wrap(err, "could not settle trade balances")
	}

	return nil
}

func (k Keeper) executeOrder(ctx context.Context, ordersCaches *types.OrdersCaches, tradeBalances *TradeBalances, fee math.LegacyDec, order *types.Order) (math.Int, bool, error) {
	denomPair := types.Pair{DenomFrom: order.DenomFrom, DenomTo: order.DenomTo}
	previousMaxPrice, has := ordersCaches.PriceAmounts[denomPair]
	if has && order.MaxPrice.LT(previousMaxPrice) {
		return math.ZeroInt(), false, nil
	}

	priceAmount := k.calculateAmountGivenPrice(ordersCaches, order.DenomFrom, order.DenomTo, order.MaxPrice, fee).TruncateInt()
	if priceAmount.LTE(math.ZeroInt()) {
		if !has || previousMaxPrice.LT(order.MaxPrice) {
			ordersCaches.PriceAmounts[denomPair] = order.MaxPrice
		}

		return math.ZeroInt(), false, nil
	}

	amount := math.MinInt(order.AmountLeft, priceAmount)
	if amount.LTE(math.ZeroInt()) {
		return math.ZeroInt(), false, nil
	}

	if order.TradeAmount.GT(math.ZeroInt()) {
		amount = math.MinInt(amount, order.TradeAmount)
	}

	address := sdk.MustAccAddressFromBech32(order.Creator)
	tradeCtx := types.TradeContext{
		Context:         ctx,
		CoinSource:      ordersCaches.AccPoolOrders.Get().String(),
		CoinTarget:      address.String(),
		GivenAmount:     amount,
		TradeDenomStart: order.DenomFrom,
		TradeDenomEnd:   order.DenomTo,
		MaxPrice:        &order.MaxPrice,
		AllowIncomplete: order.AllowIncomplete,
		TradeBalances:   tradeBalances,
		OrdersCaches:    ordersCaches,
		IsOrder:         true,
	}

	usedAmount, amountReceived1, amountReceived2, _, _, err := k.ExecuteTrade(tradeCtx)
	if err != nil {
		if errors.Is(err, types.ErrTradeAmountTooSmall) {
			return math.ZeroInt(), false, nil
		}
		if errors.Is(err, types.ErrNotEnoughLiquidity) {
			return math.ZeroInt(), false, nil
		}

		msg := fmt.Sprintf("could not execute trade (%v%v > %v)", order.AmountLeft.String(), order.DenomFrom, order.DenomTo)
		return math.Int{}, false, errors.Wrap(err, msg)
	}

	if usedAmount.IsZero() {
		return math.ZeroInt(), false, nil
	}

	order.AmountLeft = order.AmountLeft.Sub(usedAmount)
	order.AmountReceived = order.AmountReceived.Add(amountReceived2)

	// AmountLeft should never be negative zero. The comparison is still considering lower
	// than zero to cover potential rounding issues
	fullyExecuted := order.AmountLeft.LTE(math.ZeroInt())

	if order.AmountLeft.LT(math.ZeroInt()) {
		return math.Int{}, false, fmt.Errorf("order has negative amount left (%v, %v)", usedAmount.String(), order.AmountLeft.String())
	}

	if !fullyExecuted {
		k.SetOrder(ctx, *order)
	}

	return amountReceived1, fullyExecuted, nil
}

// calculateBlockEnd calculates the maximum block height that an order can be alive. If the requested block height is
// bigger than the time allowed by the parameter, the height is capped to the allowed limit.
func (k Keeper) calculateBlockEnd(maxOrderLife, addedAt, numBlocks int64) int64 {
	var life int64
	if numBlocks > maxOrderLife {
		life = maxOrderLife
	} else {
		life = numBlocks
	}

	return addedAt + life
}
