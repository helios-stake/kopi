package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCollateralAddDenom{}
	_ sdk.Msg = &MsgCollateralUpdateDepositLimit{}
	_ sdk.Msg = &MsgCollateralUpdateLTV{}
)

func (msg *MsgCollateralAddDenom) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := IsDec(msg.Ltv, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("ltv: %w", err)
	}

	if err := IsInt(msg.MaxDeposit, math.ZeroInt()); err != nil {
		return fmt.Errorf("max_deposit: %w", err)
	}

	if err := ValidateDenomName(msg.Denom); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	return nil
}

func (msg *MsgCollateralUpdateDepositLimit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := ValidateDenomName(msg.Denom); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := IsInt(msg.MaxDeposit, math.ZeroInt()); err != nil {
		return fmt.Errorf("max_deposit: %w", err)
	}

	return nil
}

func (msg *MsgCollateralUpdateLTV) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := ValidateDenomName(msg.Denom); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := IsDec(msg.Ltv, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("ltv: %w", err)
	}

	return nil
}
