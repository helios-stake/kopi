package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/denominations module sentinel errors
var (
	ErrInvalidSigner          = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidDexAsset        = sdkerrors.Register(ModuleName, 1101, "given denom is no dex asset")
	ErrInvalidCAsset          = sdkerrors.Register(ModuleName, 1102, "given denom is no c asset")
	ErrInvalidKCoin           = sdkerrors.Register(ModuleName, 1103, "given denom is no kcoin")
	ErrInvalidCollateralDenom = sdkerrors.Register(ModuleName, 1104, "given collateral denom is no collateral denom")
	ErrInvalidArbitrageDenom  = sdkerrors.Register(ModuleName, 1105, "given denom is no arbitrage denom")
	ErrInvalidAmount          = sdkerrors.Register(ModuleName, 1106, "invalid amount")
	ErrInvalidFactorReference = sdkerrors.Register(ModuleName, 1107, "invalid factor reference")
	ErrInvalidFactor          = sdkerrors.Register(ModuleName, 1108, "invalid factor")
)
