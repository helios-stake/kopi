package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
)

var RegexPercentage = regexp.MustCompile(`^(100|[1-9][0-9]?|0[1-9])%$`)

var (
	_ sdk.Msg = &MsgAutomationsAdd{}
	_ sdk.Msg = &MsgAutomationsUpdate{}
	_ sdk.Msg = &MsgAutomationsRemove{}
	_ sdk.Msg = &MsgAutomationsActive{}
)

func (msg *MsgAutomationsAdd) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	return validateAutomation(msg)
}

func (msg *MsgAutomationsUpdate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	return validateAutomation(msg)
}

func (msg *MsgAutomationsRemove) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	return nil
}

func (msg *MsgAutomationsActive) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	return nil
}
