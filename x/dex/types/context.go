package types

import (
	"context"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/utils"
)

type TradeContext struct {
	context.Context

	GivenAmount math.Int
	MaxPrice    *math.LegacyDec

	TradeDenomStart string
	TradeDenomEnd   string

	AllowIncomplete     bool
	ExcludeFromDiscount bool
	ProtocolTrade       bool
	IsOrder             bool

	CoinSource      string
	CoinTarget      string
	DiscountAddress string

	TradeBalances    TradeBalances
	TradeCalculation TradeCalculation
	OrdersCaches     *OrdersCaches
}

type Sender interface {
	SendCoins(ctx context.Context, address sdk.AccAddress, accAddress sdk.AccAddress, coins sdk.Coins) error
}

type TradeBalances interface {
	AddTransfer(string, string, string, math.Int)
	NetBalance(string, string) math.Int
	Settle(context.Context, Sender) error
}

type TradeStepContext struct {
	TradeContext

	StepDenomFrom string
	StepDenomTo   string

	Amount          math.Int
	TradeFee        math.LegacyDec
	ReserveFeeShare math.LegacyDec
}

func (tc TradeContext) TradeToBase(tradeFee, reserveFeeShare math.LegacyDec) TradeStepContext {
	return TradeStepContext{
		TradeContext:    tc,
		StepDenomFrom:   tc.TradeDenomStart,
		StepDenomTo:     utils.BaseCurrency,
		Amount:          tc.GivenAmount,
		TradeFee:        tradeFee,
		ReserveFeeShare: reserveFeeShare,
	}
}

func (tc TradeContext) TradeToTarget(tradeFee, reserveFeeShare math.LegacyDec, amount math.Int) TradeStepContext {
	return TradeStepContext{
		TradeContext:    tc,
		StepDenomFrom:   utils.BaseCurrency,
		StepDenomTo:     tc.TradeDenomEnd,
		Amount:          amount,
		TradeFee:        tradeFee,
		ReserveFeeShare: reserveFeeShare,
	}
}

type TradeCalculation interface {
	Forward(poolFrom, poolTo, offer math.LegacyDec) math.Int
	Backward(poolFrom, poolTo, result math.LegacyDec) math.Int
}
