package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	denomtypes "github.com/kopi-money/kopi/x/denominations/types"
)

var (
	_ sdk.Msg = &MsgCreateRedemptionRequest{}
	_ sdk.Msg = &MsgCancelRedemptionRequest{}
	_ sdk.Msg = &MsgUpdateRedemptionRequest{}
)

func (msg *MsgCreateRedemptionRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.ValidateDenomName(msg.Denom); err != nil {
		return err
	}

	if err := denomtypes.ValidateDenomName(msg.CAssetAmount); err != nil {
		return err
	}

	if err := denomtypes.IsDec(msg.Fee, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("fee: %w", err)
	}

	return nil
}

func (msg *MsgCancelRedemptionRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.ValidateDenomName(msg.Denom); err != nil {
		return err
	}

	return nil
}

func (msg *MsgUpdateRedemptionRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.ValidateDenomName(msg.Denom); err != nil {
		return err
	}

	if err := denomtypes.ValidateDenomName(msg.CAssetAmount); err != nil {
		return err
	}

	if err := denomtypes.IsDec(msg.Fee, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("fee: %w", err)
	}

	return nil
}
