package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	denomtypes "github.com/kopi-money/kopi/x/denominations/types"
)

var _ sdk.Msg = &MsgUpdateMovingAverageFactor{}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgUpdateMovingAverageFactor) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.IsDec(msg.MovingAverageFactor, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("moving_average_factor: %w", err)
	}

	return nil
}
