package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	denomtypes "github.com/kopi-money/kopi/x/denominations/types"
)

var (
	_ sdk.Msg = &MsgUpdateFeeAmount{}
)

func (msg *MsgUpdateFeeAmount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.IsInt(msg.FeeAmount, math.ZeroInt()); err != nil {
		return fmt.Errorf("fee_amount: %w", err)
	}

	return nil
}
