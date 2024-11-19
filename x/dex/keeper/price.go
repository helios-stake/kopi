package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	"github.com/kopi-money/kopi/constants"
	"github.com/kopi-money/kopi/x/dex/types"
)

// CalculatePrice returns the price of a given currency pair. The price is expressed how much "FROM" you need to give
// get one unit of "TO". I.e., the lower the returned value, the more valuable "FROM" is (or the less valuable "TO" is).
func (k Keeper) CalculatePrice(ctx context.Context, denomGiving, denomReceiving string) (math.LegacyDec, error) {
	price := math.LegacyOneDec()

	if denomGiving != constants.BaseCurrency {
		ratio, err := k.DenomKeeper.GetRatio(ctx, denomGiving)
		if err != nil {
			return price, err
		}

		price = price.Quo(ratio.Ratio)
	}

	if denomReceiving != constants.BaseCurrency {
		ratio, err := k.DenomKeeper.GetRatio(ctx, denomReceiving)
		if err != nil {
			return price, err
		}

		price = price.Mul(ratio.Ratio)
	}

	if price.IsZero() {
		return math.LegacyDec{}, types.ErrZeroPrice
	}

	price = math.LegacyOneDec().Quo(price)
	return price, nil
}

func (k Keeper) GetPriceInUSD(ctx context.Context, denom string) (math.LegacyDec, error) {
	referenceDenom, err := k.GetHighestUSDReference(ctx)
	if err != nil {
		return math.LegacyDec{}, fmt.Errorf("could not get highest usd reference: %w", err)
	}

	return k.CalculatePrice(ctx, denom, referenceDenom)
}

func (k Keeper) GetHighestUSDReference(ctx context.Context) (string, error) {
	var (
		ratio math.LegacyDec
		denom string
	)

	for _, usd := range k.DenomKeeper.ReferenceDenoms(ctx, constants.KUSD) {
		r, err := k.DenomKeeper.GetRatio(ctx, usd)
		if err != nil {
			return "", err
		}

		if ratio.IsNil() || ratio.GT(r.Ratio) {
			ratio = r.Ratio
			denom = usd
		}
	}

	return denom, nil
}

func (k Keeper) GetValueInBase(ctx context.Context, denom string, amount math.LegacyDec) (math.LegacyDec, error) {
	return k.GetValueIn(ctx, denom, constants.BaseCurrency, amount)
}

func (k Keeper) GetValueInUSD(ctx context.Context, denom string, amount math.LegacyDec) (math.LegacyDec, error) {
	if amount.IsZero() {
		return math.LegacyZeroDec(), nil
	}

	price, err := k.GetPriceInUSD(ctx, denom)
	if err != nil {
		return math.LegacyDec{}, err
	}

	value := amount.Quo(price)
	return k.DenomKeeper.ConvertToExponent(ctx, denom, value, 6)
}

func (k Keeper) GetValueIn(ctx context.Context, denomFrom, denomTo string, amount math.LegacyDec) (math.LegacyDec, error) {
	if amount.IsZero() {
		return math.LegacyZeroDec(), nil
	}

	price, err := k.CalculatePrice(ctx, denomFrom, denomTo)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return amount.Quo(price), nil
}
