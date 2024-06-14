package keeper

import (
	"context"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/kopi-money/kopi/utils"
	dexkeeper "github.com/kopi-money/kopi/x/dex/keeper"
	dextypes "github.com/kopi-money/kopi/x/dex/types"
	"github.com/kopi-money/kopi/x/swap/types"
	"github.com/pkg/errors"
)

// Burn is called at the end of each block to check whether the prices of the kCoins are lower than their
// "real" counterparts. If yes, funds for the base currency are minted, the kCoin is bought and received
// funds are burned such as to lower the supply of the kCoin and slightly increase its price. The amount
// that is minted is limited depending on the currency to not mint too much per block.
func (k Keeper) Burn(ctx context.Context) error {
	for _, kCoin := range k.DenomKeeper.KCoins(ctx) {
		maxBurnAmount := k.DenomKeeper.MaxBurnAmount(ctx, kCoin)
		if err := k.CheckBurn(ctx, kCoin, maxBurnAmount); err != nil {
			return errors.Wrap(err, fmt.Sprintf("could not burn denom %v", kCoin))
		}
	}

	return nil
}

func (k Keeper) CheckBurn(ctx context.Context, kCoin string, maxBurnAmount math.Int) error {
	parity, referenceDenom, err := k.DexKeeper.CalculateParity(ctx, kCoin)
	if err != nil {
		return errors.Wrap(err, "could not calculate parity")
	}

	// parity can be nil at initialization of the chain when not all currencies have liquidity. It is an edge case.
	if parity == nil || parity.GT(math.LegacyOneDec()) {
		return nil
	}

	referenceRatio, _ := k.DexKeeper.GetRatio(ctx, referenceDenom)
	mintAmountBase := k.calcBaseMintAmount(ctx, referenceRatio.Ratio, kCoin)
	mintAmountBase = math.MinInt(mintAmountBase, maxBurnAmount)
	if mintAmountBase.LTE(math.ZeroInt()) {
		return nil
	}

	// Liquidity of the kCoin is removed if present
	liq := k.DexKeeper.GetLiquidityByAddress(ctx, kCoin, dextypes.PoolReserve)
	if liq.GT(math.ZeroInt()) {
		if err = k.DexKeeper.RemoveAllLiquidityForModule(ctx, kCoin, dextypes.PoolReserve); err != nil {
			return errors.Wrap(err, "could not remove all liquidity for module")
		}
	}

	// New coins of the base currency are minted, used to buy the kCoin and burn
	if err = k.mintTradeBurn(ctx, kCoin, mintAmountBase); err != nil {
		return errors.Wrap(err, "could not mintTradeBurn")
	}

	return nil
}

func (k Keeper) calcBaseMintAmount(ctx context.Context, referenceRatio math.LegacyDec, kCoin string) math.Int {
	liqBase := k.DexKeeper.GetFullLiquidityBase(ctx, kCoin)
	liqKCoin := k.DexKeeper.GetFullLiquidityOther(ctx, kCoin)
	constantProductRoot, _ := liqBase.Mul(liqKCoin).Quo(referenceRatio).ApproxSqrt()
	mintAmount := constantProductRoot.Sub(liqBase)
	return mintAmount.TruncateInt()
}

// This function mints new XKP, buys the kCoin and then burns the tokens it has bought.
func (k Keeper) mintTradeBurn(ctx context.Context, kCoin string, mintAmountBase math.Int) error {
	mintCoins := sdk.NewCoins(sdk.NewCoin(utils.BaseCurrency, mintAmountBase))
	if err := k.BankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return errors.Wrap(err, "could not mint new XKP")
	}

	address := k.AccountKeeper.GetModuleAccount(ctx, types.ModuleName).GetAddress()

	tradeCtx := dextypes.TradeContext{
		Context:             ctx,
		GivenAmount:         mintAmountBase,
		CoinSource:          address.String(),
		CoinTarget:          address.String(),
		MaxPrice:            nil,
		TradeDenomStart:     utils.BaseCurrency,
		TradeDenomEnd:       kCoin,
		AllowIncomplete:     true,
		ExcludeFromDiscount: true,
		ProtocolTrade:       true,
		TradeBalances:       dexkeeper.NewTradeBalances(),
	}

	if _, _, _, _, _, err := k.DexKeeper.ExecuteTrade(tradeCtx); err != nil {
		if errors.Is(err, dextypes.ErrTradeAmountTooSmall) {
			return nil
		}
		if errors.Is(err, dextypes.ErrNotEnoughLiquidity) {
			return nil
		}

		return errors.Wrap(err, "could not execute trade")
	}

	if err := tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper); err != nil {
		return errors.Wrap(err, "could not settle trade balances")
	}

	if _, err := k.burnFunds(ctx, kCoin); err != nil {
		return errors.Wrap(err, "could not burn funds")
	}

	return nil
}

func (k Keeper) burnFunds(ctx context.Context, denom string) (math.Int, error) {
	burnableAmount := k.getUsableAmount(ctx, denom, types.ModuleName)

	if denom == utils.BaseCurrency {
		rewards := k.GetParams(ctx).StakingShare.Mul(math.LegacyNewDecFromInt(burnableAmount))
		if rewards.GT(math.LegacyZeroDec()) {
			rewardCoins := sdk.NewCoins(sdk.NewCoin(utils.BaseCurrency, rewards.RoundInt()))
			if err := k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, distributiontypes.ModuleName, rewardCoins); err != nil {
				return math.Int{}, errors.Wrap(err, "could not send coins to distribution")
			}

			burnableAmount = burnableAmount.Sub(rewards.RoundInt())
		}
	}

	burnCoins := sdk.NewCoins(sdk.NewCoin(denom, burnableAmount))
	if err := k.BankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins); err != nil {
		return burnableAmount, err
	}

	return burnableAmount, nil
}
