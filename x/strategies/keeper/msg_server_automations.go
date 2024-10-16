package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/x/strategies/types"
)

func (k msgServer) AutomationsAdd(ctx context.Context, msg *types.MsgAutomationsAdd) (*types.Void, error) {
	height := sdk.UnwrapSDKContext(ctx).BlockHeight()

	conditions, actions, err := k.checkAutomationMessage(ctx, msg.Creator, msg)
	if err != nil {
		return nil, err
	}

	intervalType, _ := strconv.Atoi(msg.IntervalType)
	intervalLength, _ := strconv.Atoi(msg.IntervalLength)
	validityType, _ := strconv.Atoi(msg.ValidityType)
	validitValue, _ := strconv.Atoi(msg.ValidityValue)

	k.SetAutomation(ctx, types.Automation{
		Address:        msg.Creator,
		Title:          msg.Title,
		Active:         true,
		AddedAt:        height,
		PeriodStart:    height,
		IntervalType:   int64(intervalType),
		IntervalLength: int64(intervalLength),
		ValidityType:   int64(validityType),
		ValidityValue:  int64(validitValue),
		Conditions:     conditions,
		Actions:        actions,
	})

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(
		sdk.NewEvent("automation_created",
			sdk.Attribute{Key: "address", Value: msg.Creator},
		),
	)

	return &types.Void{}, nil
}

func (k msgServer) AutomationsUpdate(ctx context.Context, msg *types.MsgAutomationsUpdate) (*types.Void, error) {
	address, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, types.ErrInvalidAddress
	}

	automation, has := k.automations.Get(ctx, msg.Index)
	if !has {
		return nil, types.ErrAutomationNotFound
	}

	if automation.Address != address.String() {
		return nil, types.ErrAutomationInvalidCreator
	}

	automation.Conditions, automation.Actions, err = k.checkAutomationMessage(ctx, msg.Creator, msg)
	if err != nil {
		return nil, err
	}

	intervalType, _ := strconv.Atoi(msg.IntervalType)
	intervalLength, _ := strconv.Atoi(msg.IntervalLength)
	validityType, _ := strconv.Atoi(msg.ValidityType)
	validitValue, _ := strconv.Atoi(msg.ValidityValue)

	if automation.IntervalType != int64(intervalType) {
		automation.PeriodTimesChecked = 0
	}

	automation.Title = msg.Title
	automation.IntervalType = int64(intervalType)
	automation.IntervalLength = int64(intervalLength)
	automation.ValidityType = int64(validityType)
	automation.ValidityValue = int64(validitValue)

	k.SetAutomation(ctx, automation)

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(
		sdk.NewEvent("automation_updated",
			sdk.Attribute{Key: "address", Value: msg.Creator},
			sdk.Attribute{Key: "index", Value: strconv.Itoa(int(msg.Index))},
		),
	)

	return &types.Void{}, nil
}

func (k Keeper) checkAutomationMessage(ctx context.Context, address string, am types.AutomationMessage) ([]*types.Condition, []*types.Action, error) {
	if k.GetAutomationFunds(ctx, address).LTE(math.ZeroInt()) {
		return nil, nil, types.ErrEmptyAutomationFunds
	}

	intervalType, err := strconv.Atoi(am.GetIntervalType())
	if err != nil {
		return nil, nil, fmt.Errorf("invalid interval type: %v", am.GetIntervalType())
	}

	intervalLength, err := strconv.Atoi(am.GetIntervalLength())
	if err != nil {
		return nil, nil, fmt.Errorf("invalid interval length: %v", am.GetIntervalLength())
	}

	validityType, err := strconv.Atoi(am.GetValidityType())
	if err != nil {
		return nil, nil, fmt.Errorf("invalid validity type: %v", am.GetValidityType())
	}

	validityValue, err := strconv.Atoi(am.GetValidityValue())
	if err != nil {
		return nil, nil, fmt.Errorf("invalid validity value: %v", am.GetValidityValue())
	}

	if len(am.GetTitle()) == 0 {
		return nil, nil, types.ErrAutomationTitleEmpty
	}

	if len(am.GetTitle()) > 30 {
		return nil, nil, types.ErrAutomationTitleTooLong
	}

	var messageConditions []types.MessageCondition
	if err = json.Unmarshal([]byte(am.GetConditions()), &messageConditions); err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal conditions: %w", err)
	}

	conditions, err := types.ConvertConditions(messageConditions)
	if err != nil {
		return nil, nil, fmt.Errorf("could not convert conditions: %w", err)
	}

	if err = k.CheckConditions(ctx, conditions); err != nil {
		return nil, nil, fmt.Errorf("invalid conditions: %w", err)
	}

	var actions []*types.Action
	if err = json.Unmarshal([]byte(am.GetActions()), &actions); err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal actions: %w", err)
	}

	if err = k.CheckActions(ctx, address, actions); err != nil {
		return nil, nil, fmt.Errorf("invalid actions: %w", err)
	}

	if intervalLength <= 0 {
		return nil, nil, types.ErrInvalidIntervalLength
	}

	if _, err = convertIntervalLength(intervalType, intervalLength); err != nil {
		return nil, nil, err
	}

	if err = isValidValidity(validityType, validityValue); err != nil {
		return nil, nil, err
	}

	return conditions, actions, nil
}

func (k msgServer) AutomationsRemove(ctx context.Context, msg *types.MsgAutomationsRemove) (*types.Void, error) {
	address, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, types.ErrInvalidAddress
	}

	automation, has := k.automations.Get(ctx, msg.Index)
	if !has {
		return nil, types.ErrAutomationNotFound
	}

	if automation.Address != address.String() {
		return nil, types.ErrAutomationInvalidCreator
	}

	k.automations.Remove(ctx, automation.Index)

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(
		sdk.NewEvent("automation_removed",
			sdk.Attribute{Key: "address", Value: msg.Creator},
			sdk.Attribute{Key: "index", Value: strconv.Itoa(int(msg.Index))},
		),
	)

	return &types.Void{}, nil
}

func (k Keeper) getAutomationActive(ctx context.Context, index uint64) bool {
	automation, exists := k.automations.Get(ctx, index)
	if exists {
		return automation.Active
	}

	return true
}

func (k msgServer) AutomationsActive(ctx context.Context, msg *types.MsgAutomationsActive) (*types.Void, error) {
	address, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, types.ErrInvalidAddress
	}

	automation, has := k.automations.Get(ctx, msg.Index)
	if !has {
		return nil, types.ErrAutomationNotFound
	}

	if automation.Address != address.String() {
		return nil, types.ErrAutomationInvalidCreator
	}

	if msg.Active {
		if k.GetAutomationFunds(ctx, msg.Creator).LTE(math.ZeroInt()) {
			return nil, types.ErrEmptyAutomationFunds
		}

		automation.PeriodStart = sdk.UnwrapSDKContext(ctx).BlockHeight()
		automation.PeriodTimesExecuted = 0
		automation.PeriodConditionFeesConsumed = 0
		automation.PeriodActionFeesConsumed = 0
	}

	automation.Active = msg.Active
	k.SetAutomation(ctx, automation)

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(
		sdk.NewEvent("automation_activation",
			sdk.Attribute{Key: "address", Value: msg.Creator},
			sdk.Attribute{Key: "index", Value: strconv.Itoa(int(msg.Index))},
			sdk.Attribute{Key: "active", Value: strconv.FormatBool(msg.Active)},
		),
	)

	return &types.Void{}, nil
}
