package keeper

import (
	"context"
	"fmt"

	"github.com/kopi-money/kopi/x/dex/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type simulate func(ctx types.TradeContext) (types.TradeSimulationResult, error)

func (k Keeper) QuerySimulateSell(ctx context.Context, req *types.QuerySimulateTradeRequest) (*types.QuerySimulateTradeResponse, error) {
	return k.querySimulateTrade(ctx, req, k.SimulateSell)
}

func (k Keeper) QuerySimulateBuy(ctx context.Context, req *types.QuerySimulateTradeRequest) (*types.QuerySimulateTradeResponse, error) {
	return k.querySimulateTrade(ctx, req, k.SimulateBuy)
}

func (k Keeper) querySimulateTrade(ctx context.Context, req *types.QuerySimulateTradeRequest, simulateFunc simulate) (*types.QuerySimulateTradeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	amount, err := ParseAmount(req.Amount)
	if err != nil {
		return nil, err
	}

	if amount.IsZero() {
		return nil, types.ErrZeroAmount
	}

	tradeCtx := types.TradeContext{
		Context:             ctx,
		TradeAmount:         amount,
		TradeDenomGiving:    req.DenomGiving,
		TradeDenomReceiving: req.DenomReceiving,
		DiscountAddress:     req.Address,
		OrdersCaches:        k.NewOrdersCaches(ctx),
	}

	tradeResult, err := simulateFunc(tradeCtx)
	if err != nil {
		return nil, fmt.Errorf("could not simulate trade: %w", err)
	}

	priceGivingUSD, err := k.GetPriceInUSD(ctx, req.DenomGiving)
	if err != nil {
		return nil, fmt.Errorf("could not get price in USD: %w", err)
	}

	priceReceivingUSD, err := k.GetPriceInUSD(ctx, req.DenomReceiving)
	if err != nil {
		return nil, fmt.Errorf("could not get price in USD: %w", err)
	}

	price := tradeResult.AmountGiven.ToLegacyDec().Quo(tradeResult.AmountReceived.ToLegacyDec())

	return &types.QuerySimulateTradeResponse{
		AmountGiven:         tradeResult.AmountGiven.String(),
		AmountGivenInUsd:    tradeResult.AmountGiven.ToLegacyDec().Quo(priceGivingUSD).RoundInt().String(),
		AmountReceived:      tradeResult.AmountReceived.String(),
		AmountReceivedInUsd: tradeResult.AmountReceived.ToLegacyDec().Quo(priceReceivingUSD).RoundInt().String(),
		Fee:                 tradeResult.FeeGiven.String(),
		Price:               price.String(),
		PriceGivenInUsd:     priceGivingUSD.String(),
		PriceReceivedInUsd:  priceReceivingUSD.String(),
	}, nil
}
