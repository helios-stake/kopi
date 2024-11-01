package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/x/tokenfactory/types"
)

func (k msgServer) CreateDenom(ctx context.Context, msg *types.MsgCreateDenom) (*types.MsgCreateDenomResponse, error) {
	factoryDenom, err := k.Keeper.CreateDenom(ctx, msg.Creator, msg.Name, msg.Symbol, msg.IconHash, msg.Exponent)
	if err != nil {
		return nil, err
	}

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"factory_denom_created",
			sdk.NewAttribute("full_name", factoryDenom.FullName),
			sdk.NewAttribute("creator", msg.Creator),
		),
	})

	return &types.MsgCreateDenomResponse{
		DisplayName: factoryDenom.DisplayName,
		FullName:    factoryDenom.FullName,
	}, nil
}
