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
		condition, err := convertCondition(&messageCondition)
		if err != nil {
			return nil, fmt.Errorf("[%d]: %w", index, err)
		}

		conditions = append(conditions, condition)
	}

	return conditions, nil
}

type MessageIn interface {
	GetValue() string
	GetConditionType() int64
	GetComparison() string
	GetReferencePrice() string
	GetString1() string
	GetString2() string
}

func convertCondition(messageCondition MessageIn) (*Condition, error) {
	value, err := math.LegacyNewDecFromStr(messageCondition.GetValue())
	if err != nil {
		return nil, fmt.Errorf("could not parse value: %w", err)
	}

	var referencePrice *math.LegacyDec
	if IsPriceChangeCondition(messageCondition.GetConditionType()) {
		var rp math.LegacyDec
		rp, err = math.LegacyNewDecFromStr(messageCondition.GetReferencePrice())
		if err != nil {
			return nil, fmt.Errorf("could not parse reference price: %w", err)
		}

		referencePrice = &rp
	}

	if err = checkAutomationString(messageCondition.GetString1()); err != nil {
		return nil, fmt.Errorf("invalid string1: %w", err)
	}

	if err = checkAutomationString(messageCondition.GetString2()); err != nil {
		return nil, fmt.Errorf("invalid string2: %w", err)
	}

	if !IsValidComparison(messageCondition.GetConditionType(), messageCondition.GetComparison()) {
		return nil, fmt.Errorf("invalid comparison: %v", messageCondition.GetComparison())
	}

	return &Condition{
		ConditionType:  messageCondition.GetConditionType(),
		String1:        messageCondition.GetString1(),
		String2:        messageCondition.GetString2(),
		Comparison:     messageCondition.GetComparison(),
		Value:          value,
		ReferencePrice: referencePrice,
	}, nil
}
