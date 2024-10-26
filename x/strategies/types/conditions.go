package types

import (
	"cosmossdk.io/math"
	"fmt"
	"strings"
)

const (
	ConditionPrice = iota
	ConditionWalletAmount
	ConditionWalletValue
	ConditionCollateralAmount
	ConditionCollateralValue
	ConditionLiquidityAmount
	ConditionLiquidityValue
	ConditionInterestRate
	ConditionCreditLineUsage
	ConditionLoanAmount
	ConditionLoanValue
	ConditionAutomationFundsAmount
	ConditionStakingRewardsAmount
	ConditionBorrowableAmount
	ConditionPriceChangeAmount
	ConditionPriceChangePercentage
)

const (
	ComparisonLessThan          = "LT"
	ComparisonLessThanEquals    = "LTE"
	ComparisonGreaterThan       = "GT"
	ComparisonGreaterThanEquals = "GTE"
	ComparisonIncreasedBy       = "IC"
	ComparisonDecreasedBy       = "DC"
)

func IsPriceChangeCondition(conditionType int64) bool {
	switch conditionType {
	case ConditionPriceChangeAmount, ConditionPriceChangePercentage:
		return true
	default:
		return false
	}
}

func IsValidComparison(conditionType int64, comp string) bool {
	switch strings.ToUpper(comp) {
	case ComparisonLessThan, ComparisonLessThanEquals, ComparisonGreaterThan, ComparisonGreaterThanEquals:
		return !IsPriceChangeCondition(conditionType)
	case ComparisonIncreasedBy, ComparisonDecreasedBy:
		return IsPriceChangeCondition(conditionType)
	default:
		return false
	}
}

func ConvertConditions(messageConditions []MessageCondition) ([]*Condition, error) {
	var conditions []*Condition

	for index, messageCondition := range messageConditions {
		condition, err := convertCondition(messageCondition)
		if err != nil {
			return nil, fmt.Errorf("index %v: %w", index, err)
		}

		conditions = append(conditions, condition)
	}

	return conditions, nil
}

func convertCondition(messageCondition MessageCondition) (*Condition, error) {
	value, err := math.LegacyNewDecFromStr(messageCondition.Value)
	if err != nil {
		return nil, fmt.Errorf("could not parse value: %w", err)
	}

	var referencePrice *math.LegacyDec
	if IsPriceChangeCondition(messageCondition.ConditionType) {
		var rp math.LegacyDec
		rp, err = math.LegacyNewDecFromStr(messageCondition.ReferencePrice)
		if err != nil {
			return nil, fmt.Errorf("could not parse reference price: %w", err)
		}

		referencePrice = &rp
	}

	if err = checkAutomationString(messageCondition.String1); err != nil {
		return nil, fmt.Errorf("invalid string1: %w", err)
	}

	if err = checkAutomationString(messageCondition.String2); err != nil {
		return nil, fmt.Errorf("invalid string2: %w", err)
	}

	if !IsValidComparison(messageCondition.ConditionType, messageCondition.Comparison) {
		return nil, fmt.Errorf("invalid comparison: %v", messageCondition.Comparison)
	}

	return &Condition{
		ConditionType:  messageCondition.ConditionType,
		String1:        messageCondition.String1,
		String2:        messageCondition.String2,
		Comparison:     messageCondition.Comparison,
		Value:          value,
		ReferencePrice: referencePrice,
	}, nil
}
