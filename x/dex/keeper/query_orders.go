package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/x/dex/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Order(ctx context.Context, req *types.QueryOrderRequest) (*types.QueryOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	order, ok := k.GetOrder(ctx, req.Index)
	if !ok {
		return nil, types.ErrOrderNotFound
	}

	referenceDenom, err := k.GetHighestUSDReference(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get reference denom: %w", err)
	}

	orderResponse, err := k.toOrderResponse(ctx, order, referenceDenom)
	if err != nil {
		return nil, err
	}

	return &types.QueryOrderResponse{
		Order: orderResponse,
	}, nil
}

func (k Keeper) Orders(ctx context.Context, req *types.QueryOrdersRequest) (*types.QueryOrdersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	referenceDenom, err := k.GetHighestUSDReference(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get reference denom: %w", err)
	}

	orders, pageRes, err := query.CollectionPaginate(
		ctx,
		k.orders,
		req.Pagination,
		func(_ uint64, order types.Order) (*types.OrderResponse, error) {
			return k.toOrderResponse(ctx, order, referenceDenom)
		},
	)

	if err != nil {
		return nil, fmt.Errorf("could not get orders from pagination: %w", err)
	}

	return &types.QueryOrdersResponse{
		Orders:     orders,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) OrdersNum(_ context.Context, _ *types.QueryOrdersNumRequest) (*types.QueryOrdersNumResponse, error) {
	return &types.QueryOrdersNumResponse{Num: int64(k.GetAllOrdersNum())}, nil
}

func (k Keeper) OrdersAddress(goCtx context.Context, req *types.QueryOrdersAddressRequest) (*types.QueryOrdersAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	referenceDenom, err := k.GetHighestUSDReference(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get reference denom: %w", err)
	}

	orders, pageRes, err := query.CollectionFilteredPaginate(
		ctx,
		k.orders,
		req.Pagination,
		func(_ uint64, order types.Order) (include bool, err error) {
			if order.Creator != req.Address {
				return false, nil
			}

			if req.DenomGiving != "" && order.DenomGiving != req.DenomGiving {
				return false, nil
			}

			if req.DenomReceiving != "" && order.DenomReceiving != req.DenomReceiving {
				return false, nil
			}

			return true, nil
		},
		func(_ uint64, order types.Order) (*types.OrderResponse, error) {
			return k.toOrderResponse(ctx, order, referenceDenom)
		},
	)

	return &types.QueryOrdersAddressResponse{
		Orders:     orders,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) toOrderResponse(ctx context.Context, order types.Order, referenceDenom string) (*types.OrderResponse, error) {
	amountLeftUSD, err := k.GetValueIn(ctx, order.DenomGiving, referenceDenom, order.AmountLeft.ToLegacyDec())
	if err != nil {
		return nil, fmt.Errorf("could not get amount left in usd: %w", err)
	}

	amountReceivedUSD, err := k.GetValueIn(ctx, order.DenomReceiving, referenceDenom, order.AmountReceived.ToLegacyDec())
	if err != nil {
		return nil, fmt.Errorf("could not get amount received in usd: %w", err)
	}

	maxPriceUSD, err := k.GetValueIn(ctx, order.DenomGiving, referenceDenom, order.MaxPrice)
	if err != nil {
		return nil, fmt.Errorf("could not get amount received in usd: %w", err)
	}

	return &types.OrderResponse{
		Index:             order.Index,
		Address:           order.Creator,
		DenomGiving:       order.DenomGiving,
		DenomReceiving:    order.DenomReceiving,
		TradeAmount:       order.TradeAmount.String(),
		AmountLeft:        order.AmountLeft.String(),
		AmountLeftUsd:     amountLeftUSD.String(),
		AmountGiven:       order.AmountGiven.String(),
		AmountReceived:    order.AmountReceived.String(),
		AmountReceivedUsd: amountReceivedUSD.String(),
		MaxPrice:          order.MaxPrice.String(),
		MaxPriceUsd:       maxPriceUSD.String(),
		NumBlocks:         order.NumBlocks,
		BlockEnd:          uint64(order.AddedAt) + order.NumBlocks,
		AllowIncomplete:   order.AllowIncomplete,
	}, nil
}

func (k Keeper) OrdersByPair(ctx context.Context, req *types.OrdersByPairRequest) (*types.QueryOrdersByPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	referenceDenom, err := k.GetHighestUSDReference(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get reference denom: %w", err)
	}

	var asks, bids []*types.OrderResponse

	iterator := k.OrderIterator(ctx)
	for iterator.Valid() {
		order := iterator.GetNext()

		if order.DenomGiving == req.DenomGiving && order.DenomReceiving == req.DenomReceiving {
			var orderResponse *types.OrderResponse
			orderResponse, err = k.toOrderResponse(ctx, order, referenceDenom)
			if err != nil {
				return nil, err
			}

			bids = append(bids, orderResponse)
		}

		if order.DenomGiving == req.DenomReceiving && order.DenomReceiving == req.DenomGiving {
			var orderResponse *types.OrderResponse
			orderResponse, err = k.toOrderResponse(ctx, order, referenceDenom)
			if err != nil {
				return nil, err
			}

			asks = append(asks, orderResponse)
		}
	}

	return &types.QueryOrdersByPairResponse{
		Bids: bids,
		Asks: asks,
	}, nil
}
