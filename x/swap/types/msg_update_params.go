package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	denomtypes "github.com/kopi-money/kopi/x/denominations/types"
)

var (
	_ sdk.Msg = &MsgUpdateBurnThreshold{}
	_ sdk.Msg = &MsgUpdateMintThreshold{}
	_ sdk.Msg = &MsgUpdateStakingShare{}
)

func (m *MsgUpdateBurnThreshold) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.IsDec(m.BurnThreshold, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("burn_threshold: %w", err)
	}

	return nil
}

func (m *MsgUpdateMintThreshold) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.IsDec(m.MintThreshold, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("burn_threshold: %w", err)
	}

	return nil
}

func (m *MsgUpdateStakingShare) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := denomtypes.IsDec(m.StakingShare, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("burn_threshold: %w", err)
	}

	return nil
}
