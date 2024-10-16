package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	"github.com/kopi-money/kopi/x/dex/types"
)

func (k Keeper) OrdersSum(ctx context.Context, _ *types.QueryOrdersSumRequest) (*types.QueryOrdersSumResponse, error) {
	denomSums := make(map[string]math.Int)

	iterator := k.OrderIterator(ctx)
	for iterator.Valid() {
		order := iterator.GetNext()
		sum, has := denomSums[order.DenomGiving]
		if !has {
			sum = math.ZeroInt()
		}

		denomSums[order.DenomGiving] = sum.Add(order.AmountLeft)
	}

	sum := math.LegacyZeroDec()
	for denom, denomSum := range denomSums {
		value, err := k.GetValueInUSD(ctx, denom, denomSum.ToLegacyDec())
		if err != nil {
			return nil, fmt.Errorf("could not get order value in usd: %w", err)
		}

		sum = sum.Add(value)
	}

	return &types.QueryOrdersSumResponse{
		Sum: sum.String(),
	}, nil
}

func (k Keeper) OrdersDenomSum(ctx context.Context, _ *types.QueryOrdersDenomSumRequest) (*types.QueryOrdersDenomSumResponse, error) {
	ordersMap := make(map[string]math.Int)

	iterator := k.OrderIterator(ctx)
	for iterator.Valid() {
		order := iterator.GetNext()
		ordersMap[order.DenomGiving] = ordersMap[order.DenomGiving].Add(order.AmountLeft)
	}

	orderSums := []*types.OrdersSum{}
	for _, denom := range k.DenomKeeper.Denoms(ctx) {
		sum := "0"
		orderSum, has := ordersMap[denom]
		if has {
			sum = orderSum.String()
		}

		orderSums = append(orderSums, &types.OrdersSum{
			DenomGiving: denom,
			Sum:         sum,
		})
	}

	return &types.QueryOrdersDenomSumResponse{
		Denoms: orderSums,
	}, nil
}
