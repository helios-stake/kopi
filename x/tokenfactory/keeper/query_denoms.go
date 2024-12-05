package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/kopi-money/kopi/x/tokenfactory/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Denoms(ctx context.Context, req *types.QueryDenomsRequest) (*types.QueryDenomsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	factoryDenoms, pageRes, err := query.CollectionPaginate(
		ctx,
		k.factoryDenoms,
		req.Pagination,
		func(key string, value types.FactoryDenom) (*types.FactoryDenom, error) {
			return &value, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("could not get factory denoms from pagination: %w", err)
	}

	return &types.QueryDenomsResponse{
		Denoms:     factoryDenoms,
		Pagination: pageRes,
	}, nil
}
