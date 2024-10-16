package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCAssetAddDenom{}
	_ sdk.Msg = &MsgCAssetUpdateDexFeeShare{}
	_ sdk.Msg = &MsgCAssetUpdateBorrowLimit{}
	_ sdk.Msg = &MsgCAssetUpdateMinimumLoanSize{}
)

func (msg *MsgCAssetAddDenom) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := IsDec(msg.DexFeeShare, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("dex_fee_share: %w", err)
	}

	if err := IsDec(msg.BorrowLimit, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("borrow_limit: %w", err)
	}

	if err := IsInt(msg.MinLoanSize, math.ZeroInt()); err != nil {
		return fmt.Errorf("min_loan_size: %w", err)
	}

	if err := validateNewDexDenom(msg); err != nil {
		return fmt.Errorf("invalid dex denom: %w", err)
	}

	if err := ValidateDenomName(msg.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := ValidateDenomName(msg.BaseDenom); err != nil {
		return fmt.Errorf("invalid base name: %w", err)
	}

	return nil
}

func (msg *MsgCAssetUpdateDexFeeShare) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := IsDec(msg.DexFeeShare, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("dex_fee_share: %w", err)
	}

	if err := ValidateDenomName(msg.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	return nil
}

func (msg *MsgCAssetUpdateBorrowLimit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := IsDec(msg.BorrowLimit, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("borrow_limit: %w", err)
	}

	if err := ValidateDenomName(msg.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	return nil
}

func (msg *MsgCAssetUpdateMinimumLoanSize) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := IsInt(msg.MinimumLoanSize, math.ZeroInt()); err != nil {
		return fmt.Errorf("minimum_loan_size: %w", err)
	}

	if err := ValidateDenomName(msg.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	return nil
}
