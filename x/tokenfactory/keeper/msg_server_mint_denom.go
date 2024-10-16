package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/x/tokenfactory/types"
)

func (k msgServer) MintDenom(ctx context.Context, msg *types.MsgMintDenom) (*types.Void, error) {
	factoryDenom, has := k.GetDenomByFullName(ctx, msg.FullFactoryDenomName)
	if !has {
		return nil, types.ErrDenomDoesntExists
	}

	if factoryDenom.Admin != msg.Creator {
		return nil, types.ErrIncorrectAdmin
	}

	amount, ok := math.NewIntFromString(msg.Amount)
	if !ok {
		return nil, types.ErrInvalidAmountFormat
	}

	if !amount.GT(math.ZeroInt()) {
		return nil, types.ErrNonPositiveAmount
	}

	coins := sdk.NewCoins(sdk.NewCoin(factoryDenom.FullName, amount))
	if err := k.BankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return nil, err
	}

	targetAddr, err := sdk.AccAddressFromBech32(msg.TargetAddress)
	if err != nil {
		return nil, types.ErrInvalidAddress
	}

	if err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, targetAddr, coins); err != nil {
		return nil, err
	}

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"factory_denom_coins_minted",
			sdk.NewAttribute("full_name", factoryDenom.FullName),
			sdk.NewAttribute("amount", msg.Amount),
			sdk.NewAttribute("target_address", msg.TargetAddress),
		),
	})

	return &types.Void{}, nil
}
