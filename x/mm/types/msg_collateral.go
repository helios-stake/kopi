package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	denomtypes "github.com/kopi-money/kopi/x/denominations/types"
)

var (
	_ sdk.Msg = &MsgAddCollateral{}
	_ sdk.Msg = &MsgRemoveCollateral{}
)

func (msg *MsgAddCollateral) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.ValidateDenomName(msg.Denom); err != nil {
		return err
	}

	if err := denomtypes.IsInt(msg.Amount, math.ZeroInt()); err != nil {
		return fmt.Errorf("amount: %w", err)
	}

	return nil
}

func (msg *MsgRemoveCollateral) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.ValidateDenomName(msg.Denom); err != nil {
		return err
	}

	if err := denomtypes.IsInt(msg.Amount, math.ZeroInt()); err != nil {
		return fmt.Errorf("amount: %w", err)
	}

	return nil
}
