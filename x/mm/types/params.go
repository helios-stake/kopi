package types

import (
	"fmt"

	"cosmossdk.io/math"
)

var (
	CollateralDiscount      = math.LegacyNewDecWithPrec(95, 2) // 0.95
	ProtocolShare           = math.LegacyNewDecWithPrec(5, 1)  // 0.5
	MinRedemptionFee        = math.LegacyNewDecWithPrec(1, 2)  // 0.01
	MaxRedemptionFee        = math.LegacyNewDecWithPrec(5, 2)  // 0.05
	MinimumInterestRate     = math.LegacyNewDecWithPrec(5, 2)  // 0.05
	A                       = math.LegacyNewDec(12)
	B                       = math.LegacyNewDec(131072)
	BlockSpeedMovingAverage = math.LegacyNewDecWithPrec(999, 3) // 0.999
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		CollateralDiscount:      CollateralDiscount,
		ProtocolShare:           ProtocolShare,
		MinRedemptionFee:        MinRedemptionFee,
		MaxRedemptionFee:        MaxRedemptionFee,
		MinInterestRate:         MinimumInterestRate,
		A:                       A,
		B:                       B,
		BlockSpeedMovingAverage: BlockSpeedMovingAverage,
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateZeroOne(p.CollateralDiscount); err != nil {
		return fmt.Errorf("invalid collateral discount: %w", err)
	}

	if err := validateZeroOne(p.MinRedemptionFee); err != nil {
		return fmt.Errorf("invalid minimum redemption fee: %w", err)
	}

	if err := validateZeroOne(p.MaxRedemptionFee); err != nil {
		return fmt.Errorf("invalid maximum redemption fee: %w", err)
	}

	if !p.MinRedemptionFee.LT(p.MaxRedemptionFee) {
		return fmt.Errorf("minimum redemption fee must not be larger than maximum redemption fee")
	}

	if err := validateZeroOne(p.MinInterestRate); err != nil {
		return fmt.Errorf("invalid minimum interest rate: %w", err)
	}

	if err := validateNumber(p.A); err != nil {
		return fmt.Errorf("invalid A: %w", err)
	}

	if err := validateNumber(p.B); err != nil {
		return fmt.Errorf("invalid B: %w", err)
	}

	return nil
}

func validateNumber(d any) error {
	_, ok := d.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", d)
	}

	return nil
}

func validateZeroOne(d any) error {
	v, ok := d.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", d)
	}

	if v.IsNil() {
		return fmt.Errorf("value is nil")
	}

	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("fee must not be larger than 1")
	}

	if v.IsNegative() {
		return fmt.Errorf("fee must be smaller than 0")
	}

	return nil
}
