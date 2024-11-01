package types

import (
	errorsmod "cosmossdk.io/errors"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
)

var RegexPercentage = regexp.MustCompile(`^(100|[1-9][0-9]?|0[1-9])%$`)

var (
	_ sdk.Msg = &MsgAutomationsAdd{}
	_ sdk.Msg = &MsgAutomationsImport{}
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

func (msg *MsgAutomationsImport) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	_, err := msg.Convert()
	if err != nil {
		return fmt.Errorf("could not convert: %w", err)
	}

	return nil
}

func validateImportAutomation(automation AutomationImport) error {
	if len(automation.GetTitle()) == 0 {
		return ErrAutomationTitleEmpty
	}

	if len(automation.GetTitle()) > 30 {
		return ErrAutomationTitleTooLong
	}

	for _, condition := range automation.GetConditions() {
		if err := checkAutomationString(condition.GetString1()); err != nil {
			return fmt.Errorf("invalid string1: %w", err)
		}

		if err := checkAutomationString(condition.GetString2()); err != nil {
			return fmt.Errorf("invalid string2: %w", err)
		}

		if !IsValidComparison(condition.GetConditionType(), condition.GetComparison()) {
			return fmt.Errorf("invalid comparison: %v", condition.GetComparison())
		}
	}

	if len(automation.Actions) > 16 {
		return fmt.Errorf("must not contain more than 16 actions")
	}

	if err := checkActions(automation.Actions); err != nil {
		return err
	}

	return nil
}

func (msg *MsgAutomationsImport) Convert() ([]MsgAutomationsAdd, error) {
	var automations []AutomationImport
	if err := json.Unmarshal([]byte(msg.Automations), &automations); err != nil {
		return nil, fmt.Errorf("could not unmarshal: %w", err)
	}

	var newAutomations []MsgAutomationsAdd
	for _, automation := range automations {
		conditions, err := json.Marshal(automation.Conditions)
		if err != nil {
			return nil, fmt.Errorf("could not marshal conditions: %w", err)
		}

		actions, err := json.Marshal(automation.Actions)
		if err != nil {
			return nil, fmt.Errorf("could not marshal conditions: %w", err)
		}

		newAutomations = append(newAutomations, MsgAutomationsAdd{
			Creator:        msg.Creator,
			Title:          automation.Title,
			IntervalType:   automation.IntervalType,
			IntervalLength: automation.IntervalLength,
			ValidityType:   automation.ValidityType,
			ValidityValue:  automation.ValidityValue,
			Conditions:     string(conditions),
			Actions:        string(actions),
		})
	}

	return newAutomations, nil
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
