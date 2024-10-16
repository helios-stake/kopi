package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/blockspeed module sentinel errors
var (
	ErrInvalidSigner  = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrDivisionByZero = sdkerrors.Register(ModuleName, 1101, "division by  zero")
)
