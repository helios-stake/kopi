package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgAddArbitrageDenom{}
	_ sdk.Msg = &MsgArbitrageUpdateBuyThreshold{}
	_ sdk.Msg = &MsgArbitrageUpdateSellThreshold{}
	_ sdk.Msg = &MsgArbitrageUpdateBuyAmount{}
	_ sdk.Msg = &MsgArbitrageUpdateSellAmount{}
	_ sdk.Msg = &MsgArbitrageUpdateRedemptionFee{}
	_ sdk.Msg = &MsgArbitrageUpdateRedemptionFeeReserveShare{}
)

func (msg *MsgAddArbitrageDenom) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := IsDec(msg.BuyThreshold, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("buy_threshold: %w", err)
	}

	if err := IsDec(msg.SellThreshold, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("sell_threshold: %w", err)
	}

	if err := IsInt(msg.BuyTradeAmount, math.ZeroInt()); err != nil {
		return fmt.Errorf("buy_trade_amount: %w", err)
	}

	if err := IsInt(msg.SellTradeAmount, math.ZeroInt()); err != nil {
		return fmt.Errorf("sell_trade_amount: %w", err)
	}

	if err := IsDec(msg.RedemptionFee, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("redemption_fee: %w", err)
	}

	if err := IsDec(msg.RedemptionFeeReserveShare, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("redemption_fee_reserve_share: %w", err)
	}

	if err := validateNewDexDenom(msg); err != nil {
		return err
	}

	if err := ValidateDenomName(msg.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := ValidateDenomName(msg.Kcoin); err != nil {
		return fmt.Errorf("invalid kcoin name: %w", err)
	}

	if err := ValidateDenomName(msg.CAsset); err != nil {
		return fmt.Errorf("invalid casset name: %w", err)
	}

	return nil
}

func (msg *MsgArbitrageUpdateBuyThreshold) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := ValidateDenomName(msg.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := IsDec(msg.BuyThreshold, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("buy_threshold: %w", err)
	}

	return nil
}

func (msg *MsgArbitrageUpdateSellThreshold) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := ValidateDenomName(msg.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := IsDec(msg.SellThreshold, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("sell_threshold: %w", err)
	}

	return nil
}

func (msg *MsgArbitrageUpdateBuyAmount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := ValidateDenomName(msg.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := IsInt(msg.BuyAmount, math.ZeroInt()); err != nil {
		return fmt.Errorf("buy_amount: %w", err)
	}

	return nil
}

func (msg *MsgArbitrageUpdateSellAmount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := ValidateDenomName(msg.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := IsInt(msg.SellAmount, math.ZeroInt()); err != nil {
		return fmt.Errorf("sell_amount: %w", err)
	}

	return nil
}

func (msg *MsgArbitrageUpdateRedemptionFee) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := ValidateDenomName(msg.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := IsDec(msg.RedemptionFee, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("redemption_fee: %w", err)
	}

	return nil
}

func (msg *MsgArbitrageUpdateRedemptionFeeReserveShare) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := ValidateDenomName(msg.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := IsDec(msg.RedemptionFeeReserveShare, math.LegacyZeroDec()); err != nil {
		return fmt.Errorf("redemption_fee_reserve_share: %w", err)
	}

	return nil
}
