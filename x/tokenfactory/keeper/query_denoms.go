package keeper

import (
	"context"

	"github.com/kopi-money/kopi/x/tokenfactory/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Denoms(ctx context.Context, req *types.QueryDenomsRequest) (*types.QueryDenomsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	factoryDenoms := k.GetAllDenoms(ctx)

	response := types.QueryDenomsResponse{
		Denoms: make([]*types.FactoryDenom, len(factoryDenoms)),
	}

	for i, factoryDenom := range factoryDenoms {
		response.Denoms[i] = &factoryDenom
	}

	return &response, nil
}
