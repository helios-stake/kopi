package keeper_test

import (
	"cosmossdk.io/math"
	keepertest "github.com/kopi-money/kopi/testutil/keeper"
	denomtypes "github.com/kopi-money/kopi/x/denominations/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRatios1(t *testing.T) {
	k, denomMsg, _, ctx := keepertest.SetupDexDenomMsgServer(t)

	// Add BTC with a price 1 BTC = 1000 kUSD
	require.NoError(t, keepertest.AddDexDenom(ctx, denomMsg, &denomtypes.MsgDexAddDenom{
		Authority:    k.DenomKeeper.GetAuthority(),
		Name:         "bitcoin",
		Factor:       "1000ukusd",
		MinLiquidity: "1000000",
		MinOrderSize: "1000000",
		Exponent:     8,
	}))

	ratio, err := k.GetRatio(ctx, "bitcoin")
	require.NoError(t, err)

	// 1 BTC = 1000 kUSD = 4000 XKP
	// 1 / 4000 = 0.00025
	require.Equal(t, math.LegacyNewDecWithPrec(25, 5), ratio.Ratio)
}
