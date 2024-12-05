package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	denomtypes "github.com/kopi-money/kopi/x/denominations/types"
)

var (
	_ sdk.Msg = &MsgUpdateCollateralDiscount{}
	_ sdk.Msg = &MsgUpdateInterestRateParameters{}
	_ sdk.Msg = &MsgUpdateProtocolShare{}
	_ sdk.Msg = &MsgUpdateRedemptionFees{}
)

func (msg *MsgUpdateCollateralDiscount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.IsDec(msg.CollateralDiscount, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("collateral_discount: %w", err)
	}

	return nil
}

func (msg *MsgUpdateInterestRateParameters) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.IsDec(msg.MinInterestRate, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("min_interest_rate: %w", err)
	}

	if err := denomtypes.IsDec(msg.A, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("a: %w", err)
	}

	if err := denomtypes.IsDec(msg.B, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("b: %w", err)
	}

	return nil
}

func (msg *MsgUpdateProtocolShare) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.IsDec(msg.ProtocolShare, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("protocol_share: %w", err)
	}

	return nil
}

func (msg *MsgUpdateRedemptionFees) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.IsDec(msg.MinRedemptionFee, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("min_redemption_fee: %w", err)
	}

	if err := denomtypes.IsDec(msg.MaxRedemptionFee, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("max_redemption_fee: %w", err)
	}

	return nil
}
