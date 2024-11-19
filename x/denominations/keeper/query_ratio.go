package keeper

import (
	"context"
	"fmt"
	"github.com/kopi-money/kopi/x/denominations/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Ratio(ctx context.Context, req *types.QueryGetRatioRequest) (*types.QueryGetRatioResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if !k.IsValidDenom(ctx, req.Denom) {
		return nil, fmt.Errorf("no ratio for denom: %v", req.Denom)
	}

	ratio, err := k.GetRatio(ctx, req.Denom)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetRatioResponse{Ratio: types.RatioResponse{
		Denom: req.Denom,
		Ratio: ratio.Ratio.String(),
	}}, nil
}

func (k Keeper) Ratios(ctx context.Context, _ *types.QueryGetRatiosRequest) (*types.QueryGetRatiosResponse, error) {
	var ratios []types.RatioResponse

	for _, ratio := range k.GetAllRatios(ctx) {
		ratios = append(ratios, types.RatioResponse{
			Denom: ratio.Denom,
			Ratio: ratio.Ratio.String(),
		})
	}

	return &types.QueryGetRatiosResponse{
		Ratios: ratios,
	}, nil
}
