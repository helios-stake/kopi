package keeper_test

import (
	"testing"

	keepertest "github.com/kopi-money/kopi/testutil/keeper"
	"github.com/kopi-money/kopi/x/tokenfactory/keeper"
	"github.com/stretchr/testify/require"
)

func TestCreateDenom1(t *testing.T) {
	k, msgServer, ctx := keepertest.SetupTokenfactoryMsgServer(t)

	_, has := k.GetDenomByDisplayName(ctx, "testdenom")
	require.False(t, has)

	factoryDenomHash, err := keepertest.CreateFactoryDenom(ctx, msgServer, keepertest.Alice, "testdenom", 6)
	require.NoError(t, err)

	_, has = k.GetDenomByDisplayName(ctx, "testdenom")
	require.True(t, has)
	_, has = k.GetDenomByDisplayName(ctx, factoryDenomHash)
	require.True(t, has)

	_, has = k.GetDenomByDisplayName(ctx, "testdenom2")
	require.False(t, has)
	_, has = k.GetDenomByFullName(ctx, keeper.ToFullName("testdenom2"))
	require.False(t, has)
}

func TestCreateDenom2(t *testing.T) {
	k, msgServer, ctx := keepertest.SetupTokenfactoryMsgServer(t)

	_, has := k.GetDenomByDisplayName(ctx, "testdenom")
	require.False(t, has)

	factoryDenomHash, err := keepertest.CreateFactoryDenom(ctx, msgServer, keepertest.Alice, "testdenom", 6)
	require.NoError(t, err)

	_, has = k.GetDenomByDisplayName(ctx, "testdenom")
	require.True(t, has)
	_, has = k.GetDenomByFullName(ctx, factoryDenomHash)
	require.True(t, has)

	_, has = k.GetDenomByDisplayName(ctx, "testdenom2")
	require.False(t, has)
	_, has = k.GetDenomByFullName(ctx, keeper.ToFullName("testdenom2"))
	require.False(t, has)
}
