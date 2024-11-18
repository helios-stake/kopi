package types

import (
	"cosmossdk.io/math"
	"fmt"
)

var (
	MovingAverageFactors = math.LegacyNewDecWithPrec(9999, 4) // 0.9999
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		MovingAverageFactor: MovingAverageFactors,
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.MovingAverageFactor.IsNil() {
		return fmt.Errorf("moving_average_factor must not be null")
	}

	if p.MovingAverageFactor.GT(math.LegacyOneDec()) {
		return fmt.Errorf("moving_average_factor must not be larger than 1")
	}

	if p.MovingAverageFactor.IsNegative() {
		return fmt.Errorf("moving_average_factor must not be smaller than 0")
	}

	return nil
}
