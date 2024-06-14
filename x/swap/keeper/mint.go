package keeper

import (
	"context"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/utils"
	dexkeeper "github.com/kopi-money/kopi/x/dex/keeper"
	dextypes "github.com/kopi-money/kopi/x/dex/types"
	"github.com/kopi-money/kopi/x/swap/types"
	"github.com/pkg/errors"
)

// Mint is called at the end of each block to check whether the prices of the kCoins are higher than their
// "real" counterparts. If yes, funds for the kCoin are minted, the kCoin is sold for the base
// currency and received funds are burned such as to increase the supply of the kCoin and slightly decrease
// its price. The amount that is minted is limited depending on the currency to not mint too much per block.
func (k Keeper) Mint(ctx context.Context) error {
	for _, kCoin := range k.DenomKeeper.KCoins(ctx) {
		maxMintAmount := k.DenomKeeper.MaxMintAmount(ctx, kCoin)
		if err := k.CheckMint(ctx, kCoin, maxMintAmount); err != nil {
			return errors.Wrap(err, "could not mint denom")
		}
	}

	return nil
}

// CheckMint checks the parity of a given kCoin. If it is above 1, new coins are minted and sold in favor of
// the base currency.
func (k Keeper) CheckMint(ctx context.Context, kCoin string, maxMintAmount math.Int) error {
	parity, referenceDenom, err := k.DexKeeper.CalculateParity(ctx, kCoin)
	if err != nil {
		return errors.Wrap(err, "could not calculate parity")
	}

	if parity == nil || parity.LT(math.LegacyOneDec()) {
		return nil
	}

	referenceRatio, _ := k.DexKeeper.GetRatio(ctx, referenceDenom)
	mintAmount := k.calcKCoinMintAmount(ctx, referenceRatio.Ratio, kCoin)
	mintAmount = math.MinInt(mintAmount, maxMintAmount)
	mintAmount = k.adjustForSupplyCap(ctx, kCoin, mintAmount)
	if mintAmount.LTE(math.OneInt()) {
		return nil
	}

	mintCoins := sdk.NewCoins(sdk.NewCoin(kCoin, mintAmount))
	if err = k.BankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return errors.Wrap(err, "could not mint coins")
	}

	address := k.AccountKeeper.GetModuleAccount(ctx, types.ModuleName).GetAddress()

	tradeCtx := dextypes.TradeContext{
		Context:             ctx,
		CoinSource:          address.String(),
		CoinTarget:          address.String(),
		GivenAmount:         mintAmount,
		MaxPrice:            nil,
		TradeDenomStart:     kCoin,
		TradeDenomEnd:       utils.BaseCurrency,
		AllowIncomplete:     true,
		ExcludeFromDiscount: true,
		ProtocolTrade:       true,
		TradeBalances:       dexkeeper.NewTradeBalances(),
	}

	if _, _, _, _, _, err = k.DexKeeper.ExecuteTrade(tradeCtx); err != nil {
		if errors.Is(err, dextypes.ErrTradeAmountTooSmall) {
			return nil
		}
		if errors.Is(err, dextypes.ErrNotEnoughLiquidity) {
			return nil
		}

		return errors.Wrap(err, "could not execute incomplete trade")
	}

	if err = tradeCtx.TradeBalances.Settle(ctx, k.BankKeeper); err != nil {
		return errors.Wrap(err, "could not settle trade balances")
	}

	if _, err = k.burnFunds(ctx, utils.BaseCurrency); err != nil {
		return errors.Wrap(err, "could not burn funds")
	}

	return nil
}

func (k Keeper) adjustForSupplyCap(ctx context.Context, kCoin string, amountToAdd math.Int) math.Int {
	supply := k.BankKeeper.GetSupply(ctx, kCoin).Amount
	maximumSupply := k.DenomKeeper.MaxSupply(ctx, kCoin)

	maximumAddableAmount := maximumSupply.Sub(supply.Add(amountToAdd))
	amountToAdd = math.MinInt(maximumAddableAmount, amountToAdd)

	return amountToAdd
}

func (k Keeper) calcKCoinMintAmount(ctx context.Context, referenceRatio math.LegacyDec, kCoin string) math.Int {
	referenceRatio = math.LegacyOneDec().Quo(referenceRatio)
	liqBase := k.DexKeeper.GetFullLiquidityBase(ctx, kCoin)
	liqKCoin := k.DexKeeper.GetFullLiquidityOther(ctx, kCoin)
	constantProductRoot, _ := liqBase.Mul(liqKCoin).Quo(referenceRatio).ApproxSqrt()
	mintAmount := constantProductRoot.Sub(liqKCoin)
	return mintAmount.TruncateInt()
}

func (k Keeper) getUsableAmount(ctx context.Context, denom, module string) math.Int {
	address := k.AccountKeeper.GetModuleAccount(ctx, module).GetAddress()
	coins := k.BankKeeper.SpendableCoins(ctx, address)

	for _, coin := range coins {
		if coin.Denom == denom {
			return coin.Amount
		}
	}

	return math.ZeroInt()
}
