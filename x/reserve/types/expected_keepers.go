package types

import (
	"context"

	"cosmossdk.io/math"
	denomtypes "github.com/kopi-money/kopi/x/denominations/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here

	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	SpendableCoin(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin

	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error

	BurnCoins(ctx context.Context, name string, amt sdk.Coins) error
	MintCoins(ctx context.Context, name string, amt sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

type DenomKeeper interface {
	GetCAssetByBaseName(ctx context.Context, baseDenom string) (*denomtypes.CAsset, error)
	IsKCoin(ctx context.Context, denom string) bool
	IsValidDenom(ctx context.Context, denom string) bool
}

type DexKeeper interface {
	AddLiquidity(ctx context.Context, address sdk.AccAddress, denom string, amount math.Int) (math.Int, error)
}

type MMKeeper interface {
	CalculateCAssetValue(ctx context.Context, cAsset *denomtypes.CAsset) math.LegacyDec
}
