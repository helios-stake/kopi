package keeper

import (
	"context"
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/x/tokenfactory/types"
)

func (k msgServer) Buyback(ctx context.Context, msg *types.MsgBuyback) (*types.Void, error) {
	factoryDenom, has := k.GetDenomByFullName(ctx, msg.FullFactoryDenomName)
	if !has {
		return nil, types.ErrDenomDoesntExists
	}

	pool, has := k.liquidityPools.Get(ctx, factoryDenom.FullName)
	if !has {
		return nil, types.ErrPoolDoesNotExist
	}

	acc, _ := sdk.AccAddressFromBech32(msg.Creator)
	res, err := k.Keeper.Sell(ctx, TradeData{
		factoryDenom:    msg.FullFactoryDenomName,
		creator:         msg.Creator,
		denomGiving:     pool.KCoin,
		denomReceiving:  factoryDenom.FullName,
		tradeAmount:     msg.BuybackAmount,
		allowIncomplete: true,
	})

	if err != nil {
		return nil, err
	}

	coins := sdk.NewCoins(sdk.NewCoin(factoryDenom.FullName, math.NewInt(res.AmountReceivedNet)))
	if err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, acc, types.ModuleName, coins); err != nil {
		return nil, err
	}

	if err = k.BankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
		return nil, err
	}

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"factory_denom_buyback",
			sdk.NewAttribute("factor_denom_full_name", factoryDenom.FullName),
			sdk.NewAttribute("buyback_amount", strconv.Itoa(int(res.AmountGivenGross))),
			sdk.NewAttribute("amount_burned", strconv.Itoa(int(res.AmountReceivedNet))),
		),
	})

	return &types.Void{}, nil
}
