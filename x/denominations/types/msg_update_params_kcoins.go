package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgKCoinAddDenom{}
	_ sdk.Msg = &MsgKCoinUpdateSupplyLimit{}
	_ sdk.Msg = &MsgKCoinUpdateMintAmount{}
	_ sdk.Msg = &MsgKCoinUpdateBurnAmount{}
	_ sdk.Msg = &MsgKCoinRemoveReferences{}
	_ sdk.Msg = &MsgKCoinAddReferences{}
)

func (msg *MsgKCoinAddDenom) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := IsInt(msg.MaxSupply, math.ZeroInt()); err != nil {
		return fmt.Errorf("max_supply: %w", err)
	}

	if err := IsInt(msg.MaxBurnAmount, math.ZeroInt()); err != nil {
		return fmt.Errorf("max_burn_amount: %w", err)
	}

	if err := IsInt(msg.MaxMintAmount, math.ZeroInt()); err != nil {
		return fmt.Errorf("max_mint_amount: %w", err)
	}

	if err := validateNewDexDenom(msg); err != nil {
		return err
	}

	if err := ValidateDenomName(msg.Name); err != nil {
		return fmt.Errorf("invalid denom name: %w", err)
	}

	for _, referenceDenom := range msg.References {
		if err := ValidateDenomName(referenceDenom); err != nil {
			return fmt.Errorf("invalid denom name: %v: %w", referenceDenom, err)
		}
	}

	return nil
}

func (msg *MsgKCoinUpdateSupplyLimit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := IsInt(msg.MaxSupply, math.ZeroInt()); err != nil {
		return fmt.Errorf("max_supply: %w", err)
	}

	if err := ValidateDenomName(msg.Denom); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	return nil
}

func (msg *MsgKCoinUpdateMintAmount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := IsInt(msg.MaxMintAmount, math.ZeroInt()); err != nil {
		return fmt.Errorf("max_mint_amount: %w", err)
	}

	if err := ValidateDenomName(msg.Denom); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	return nil
}

func (msg *MsgKCoinUpdateBurnAmount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := IsInt(msg.MaxBurnAmount, math.ZeroInt()); err != nil {
		return fmt.Errorf("max_burn_amount: %w", err)
	}

	if err := ValidateDenomName(msg.Denom); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	return nil
}

func (msg *MsgKCoinRemoveReferences) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := ValidateDenomName(msg.Denom); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	return nil
}

func (msg *MsgKCoinAddReferences) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := ValidateDenomName(msg.Denom); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	return nil
}
