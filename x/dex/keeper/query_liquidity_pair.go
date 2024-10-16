package keeper

import (
	"context"

	"github.com/kopi-money/kopi/x/dex/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) LiquidityPairAll(ctx context.Context, _ *types.QueryAllLiquidityPairRequest) (*types.QueryAllLiquidityPairResponse, error) {
	return &types.QueryAllLiquidityPairResponse{
		LiquidityPair: k.GetAllLiquidityPair(ctx),
	}, nil
}

func (k Keeper) LiquidityPair(ctx context.Context, req *types.QueryGetLiquidityPairRequest) (*types.QueryGetLiquidityPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	liquidityPair, err := k.GetLiquidityPair(ctx, req.Denom)
	if err != nil {
		return nil, err
	}

	fullOther := k.GetFullLiquidityOther(ctx, req.Denom)
	fullBase := k.GetFullLiquidityBase(ctx, req.Denom)

	return &types.QueryGetLiquidityPairResponse{
		Denom:        liquidityPair.Denom,
		VirtualBase:  liquidityPair.VirtualBase.String(),
		VirtualOther: liquidityPair.VirtualOther.String(),
		FullBase:     fullBase.String(),
		FullOther:    fullOther.String(),
	}, nil
}
