package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/kopi-money/kopi/x/denominations/keeper"
	"github.com/stretchr/testify/require"
)

func TestConversion1(t *testing.T) {
	var sourceExponent uint64 = 1
	var targetExponent uint64 = 2

	value1 := math.LegacyNewDec(1)
	value2 := keeper.ConvertToExponent(value1, sourceExponent, targetExponent)

	require.Equal(t, int64(10), value2.RoundInt64())
}

func TestConversion2(t *testing.T) {
	var sourceExponent uint64 = 2
	var targetExponent uint64 = 1

	value1 := math.LegacyNewDec(10)
	value2 := keeper.ConvertToExponent(value1, sourceExponent, targetExponent)

	require.Equal(t, int64(1), value2.RoundInt64())
}

func TestConversion3(t *testing.T) {
	var sourceExponent uint64 = 1
	var targetExponent uint64 = 1

	value1 := math.LegacyNewDec(1)
	value2 := keeper.ConvertToExponent(value1, sourceExponent, targetExponent)

	require.Equal(t, int64(1), value2.RoundInt64())
}
