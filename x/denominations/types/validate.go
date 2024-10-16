package types

import (
	"cosmossdk.io/math"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	coinRegex = regexp.MustCompile(`^(\d+\.\d+|\d+)([a-zA-Z]+(?:/[a-zA-Z0-9]+)?)?$`)
	hashRegex = regexp.MustCompile(`^[A-F0-9]{64}$`)
)

func IsDec(decString string, minValue math.LegacyDec) error {
	dec, err := math.LegacyNewDecFromStr(decString)
	if err != nil {
		return fmt.Errorf("invalid dec: %v", decString)
	}

	if minValue.IsNil() {
		return nil
	}

	if dec.LT(minValue) {
		return fmt.Errorf("%v less than %v", dec.String(), minValue.String())
	}

	return nil
}

func IsInt(intString string, minValue math.Int) error {
	integer, ok := math.NewIntFromString(intString)
	if !ok {
		return fmt.Errorf("invalid int: %v", intString)
	}

	if minValue.IsNil() {
		return nil
	}

	if integer.LT(minValue) {
		return fmt.Errorf("%v less than %v", integer.Int64(), minValue.Int64())
	}

	return nil
}

type NewDexDenom interface {
	GetFactor() string
	GetMinLiquidity() string
	GetMinOrderSize() string
	GetName() string
}

func validateNewDexDenom(denom NewDexDenom) error {
	if _, _, err := ExtractNumberAndString(denom.GetFactor()); err != nil {
		return err
	}

	if err := IsInt(denom.GetMinLiquidity(), math.ZeroInt()); err != nil {
		return fmt.Errorf("min_liquidity: %w", err)
	}

	if err := IsInt(denom.GetMinOrderSize(), math.ZeroInt()); err != nil {
		return fmt.Errorf("min_order_size: %w", err)
	}

	if err := ValidateDenomName(denom.GetName()); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	return nil
}

func ExtractNumberAndString(input string) (math.LegacyDec, string, error) {
	factor, err := math.LegacyNewDecFromStr(input)
	if err == nil {
		return factor, "", nil
	}

	matches := coinRegex.FindStringSubmatch(input)
	if len(matches) != 3 {
		return math.LegacyDec{}, "", fmt.Errorf("factor string has to be in coin-string format")
	}

	factor, err = math.LegacyNewDecFromStr(matches[1])
	if err != nil {
		return factor, "", err
	}

	return factor, matches[2], nil
}

// ValidateDenomName checks denom names. Names must not be longer than 16 characters. Exceptions are denoms like IBC or
// factory tokens with the pattern prefix/hash. Here the prefix must not be longer than 16 characters and the hash no
// longer than 64.
func ValidateDenomName(name string) error {
	if strings.Contains(name, "/") {
		parts := strings.Split(name, "/")
		if len(parts) > 2 {
			return fmt.Errorf("invalid number of denom parts")
		}

		if len(parts[1]) > 64 {
			return fmt.Errorf("hash part must not contain more than 64 characters")
		}

		if !hashRegex.Match([]byte(strings.ToUpper(parts[1]))) {
			return fmt.Errorf("invalid hash part")
		}

		name = parts[0]
	}

	if len(name) > 16 {
		return fmt.Errorf("name must not have more than 16 characters")
	}

	for i := 0; i < len(name); i++ {
		if name[i] > unicode.MaxASCII {
			return fmt.Errorf("invalid non-ascii character: %v", name[i])
		}
	}

	return nil
}
