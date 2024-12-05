package types

import (
	"fmt"

	"cosmossdk.io/math"
)

var (
	kCoinBurnShare = math.LegacyOneDec()
)

func DefaultParams() Params {
	return Params{
		KcoinBurnShare: kCoinBurnShare,
	}
}

func NewParams() Params {
	return DefaultParams()
}

func (p Params) Validate() error {
	if kCoinBurnShare.IsNegative() {
		return fmt.Errorf("kcoin burn share must not be below 0")
	}

	if kCoinBurnShare.GT(math.LegacyOneDec()) {
		return fmt.Errorf("kcoin burn share must not be larger than 1")
	}

	return nil
}
