package keeper_test

import (
	"fmt"
	"github.com/kopi-money/kopi/constants"
	"testing"

	keepertest "github.com/kopi-money/kopi/testutil/keeper"
	"github.com/kopi-money/kopi/x/tokenfactory/types"
	"github.com/stretchr/testify/require"
)

func TestTrade1(t *testing.T) {
	k, msgServer, ctx := keepertest.SetupTokenfactoryMsgServer(t)

	factoryDenomHash, err := keepertest.CreateFactoryDenom(ctx, msgServer, keepertest.Alice, "testdenom", 6)
	require.NoError(t, err)
	require.NoError(t, keepertest.MintFactoryDenom(ctx, msgServer, keepertest.Alice, factoryDenomHash, keepertest.Alice, "2000000"))
	require.NoError(t, keepertest.CreatePool(ctx, msgServer, keepertest.Alice, factoryDenomHash, "1000000", constants.KUSD, "1000000", "0.1", 10))

	pool, _ := k.GetLiquidityPool(ctx, factoryDenomHash)
	require.Equal(t, int64(1_000_000), pool.FactoryDenomAmount.Int64())
	require.Equal(t, int64(1_000_000), pool.KCoinAmount.Int64())

	response, err := keepertest.FactoryDenomSell(ctx, msgServer, keepertest.Alice, factoryDenomHash, factoryDenomHash, constants.KUSD, "10000", "", true)
	require.NoError(t, err)
	require.Equal(t, int64(10000), response.AmountGivenGross)
	require.Equal(t, int64(8901), response.AmountReceivedNet)

	pool, _ = k.GetLiquidityPool(ctx, factoryDenomHash)
	require.Equal(t, int64(1010000), pool.FactoryDenomAmount.Int64())
	require.Equal(t, int64(992089), pool.KCoinAmount.Int64())

	paidPrice1 := float64(response.AmountGivenGross) / float64(response.AmountReceivedNet)
	maxPriceString := fmt.Sprintf("%.8f", paidPrice1)
	response, err = keepertest.FactoryDenomSell(ctx, msgServer, keepertest.Alice, factoryDenomHash, factoryDenomHash, constants.KUSD, "10000", maxPriceString, false)
	require.ErrorIs(t, err, types.ErrMarketPriceTooHigh)
	_, err = keepertest.FactoryDenomSell(ctx, msgServer, keepertest.Alice, factoryDenomHash, factoryDenomHash, constants.KUSD, "10000", maxPriceString, true)
	require.ErrorIs(t, err, types.ErrMarketPriceTooHigh)
}

func TestTrade2(t *testing.T) {
	_, msgServer, ctx := keepertest.SetupTokenfactoryMsgServer(t)

	factoryDenomHash, err := keepertest.CreateFactoryDenom(ctx, msgServer, keepertest.Alice, "testdenom", 6)
	require.NoError(t, err)
	require.NoError(t, keepertest.MintFactoryDenom(ctx, msgServer, keepertest.Alice, factoryDenomHash, keepertest.Alice, "2000000"))
	require.NoError(t, keepertest.CreatePool(ctx, msgServer, keepertest.Alice, factoryDenomHash, "1000000", constants.KUSD, "1000000", "0.1", 10))

	response, err := keepertest.FactoryDenomBuy(ctx, msgServer, keepertest.Alice, factoryDenomHash, factoryDenomHash, constants.KUSD, "10000", "", true)
	require.NoError(t, err)

	paidPrice1 := float64(response.AmountGivenGross) / float64(response.AmountReceivedNet)
	maxPriceString := fmt.Sprintf("%.8f", paidPrice1)
	response, err = keepertest.FactoryDenomBuy(ctx, msgServer, keepertest.Alice, factoryDenomHash, factoryDenomHash, constants.KUSD, "10000", maxPriceString, false)
	require.ErrorIs(t, err, types.ErrMarketPriceTooHigh)
	_, err = keepertest.FactoryDenomSell(ctx, msgServer, keepertest.Alice, factoryDenomHash, factoryDenomHash, constants.KUSD, "10000", maxPriceString, true)
	require.ErrorIs(t, err, types.ErrMarketPriceTooHigh)
}

func TestTrade3(t *testing.T) {
	k, msgServer, ctx := keepertest.SetupTokenfactoryMsgServer(t)

	factoryDenomHash, err := keepertest.CreateFactoryDenom(ctx, msgServer, keepertest.Alice, "testdenom", 6)
	require.NoError(t, err)
	require.NoError(t, keepertest.MintFactoryDenom(ctx, msgServer, keepertest.Alice, factoryDenomHash, keepertest.Alice, "2000000"))
	require.NoError(t, keepertest.CreatePool(ctx, msgServer, keepertest.Alice, factoryDenomHash, "1000000", constants.KUSD, "1000000", "0.1", 10))

	pool, _ := k.GetLiquidityPool(ctx, factoryDenomHash)
	require.Equal(t, int64(1000000), pool.FactoryDenomAmount.Int64())
	require.Equal(t, int64(1000000), pool.KCoinAmount.Int64())

	response, err := keepertest.FactoryDenomSell(ctx, msgServer, keepertest.Alice, factoryDenomHash, factoryDenomHash, constants.KUSD, "1000", "1.11", true)
	require.ErrorIs(t, err, types.ErrMarketPriceTooHigh)

	response, err = keepertest.FactoryDenomSell(ctx, msgServer, keepertest.Alice, factoryDenomHash, factoryDenomHash, constants.KUSD, "10000", "1.12", true)
	require.NoError(t, err, types.ErrEmptyTrade)
	require.True(t, response.AmountGivenGross < int64(10000))
	require.True(t, response.AmountReceivedNet < int64(8910))
}

func TestTrade4(t *testing.T) {
	k, msgServer, ctx := keepertest.SetupTokenfactoryMsgServer(t)

	factoryDenomHash, err := keepertest.CreateFactoryDenom(ctx, msgServer, keepertest.Alice, "testdenom", 6)
	require.NoError(t, err)

	require.NoError(t, keepertest.MintFactoryDenom(ctx, msgServer, keepertest.Alice, factoryDenomHash, keepertest.Alice, "1000000"))
	require.NoError(t, keepertest.CreatePool(ctx, msgServer, keepertest.Alice, factoryDenomHash, "1000000", constants.KUSD, "1000000", "0.1", 10))

	poolAcc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolFactoryLiquidity)
	poolBalance1 := k.BankKeeper.SpendableCoins(ctx, poolAcc.GetAddress())
	poolBalanceFactory1 := poolBalance1.AmountOf(factoryDenomHash).Int64()
	poolBalanceKCoin1 := poolBalance1.AmountOf(constants.KUSD).Int64()

	require.NoError(t, keepertest.MintFactoryDenom(ctx, msgServer, keepertest.Alice, factoryDenomHash, keepertest.Bob, "10000"))
	response, err := keepertest.FactoryDenomSell(ctx, msgServer, keepertest.Bob, factoryDenomHash, constants.KUSD, factoryDenomHash, "10000", "", true)
	require.NoError(t, err)
	require.Equal(t, int64(10000), response.AmountGivenGross)
	require.Equal(t, int64(8990), response.AmountGivenNet)
	require.Equal(t, int64(8901), response.AmountReceivedGross)
	require.Equal(t, int64(8901), response.AmountReceivedNet)
	require.Equal(t, int64(1000), response.FeePool)
	require.Equal(t, int64(10), response.FeeReserve)

	pool, _ := k.GetLiquidityPool(ctx, factoryDenomHash)
	require.Equal(t, int64(1_009_990), pool.KCoinAmount.Int64())
	require.Equal(t, int64(991_099), pool.FactoryDenomAmount.Int64())

	poolBalance2 := k.BankKeeper.SpendableCoins(ctx, poolAcc.GetAddress())
	poolBalanceFactory2 := poolBalance2.AmountOf(factoryDenomHash).Int64()
	poolBalanceKCoin2 := poolBalance2.AmountOf(constants.KUSD).Int64()

	require.Equal(t, poolBalanceKCoin2-poolBalanceKCoin1, response.AmountGivenGross)
	require.Equal(t, poolBalanceFactory1-poolBalanceFactory2, response.AmountReceivedNet)
}

func TestTrade5(t *testing.T) {
	k, msgServer, ctx := keepertest.SetupTokenfactoryMsgServer(t)

	factoryDenomHash, err := keepertest.CreateFactoryDenom(ctx, msgServer, keepertest.Alice, "testdenom", 6)
	require.NoError(t, err)
	require.NoError(t, keepertest.MintFactoryDenom(ctx, msgServer, keepertest.Alice, factoryDenomHash, keepertest.Alice, "1000000"))
	require.NoError(t, keepertest.CreatePool(ctx, msgServer, keepertest.Alice, factoryDenomHash, "1000000", constants.KUSD, "1000000", "0.1", 10))

	poolAcc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolFactoryLiquidity)
	poolBalance1 := k.BankKeeper.SpendableCoins(ctx, poolAcc.GetAddress())
	poolBalanceFactory1 := poolBalance1.AmountOf(factoryDenomHash).Int64()
	poolBalanceKCoin1 := poolBalance1.AmountOf(constants.KUSD).Int64()

	require.NoError(t, keepertest.MintFactoryDenom(ctx, msgServer, keepertest.Alice, factoryDenomHash, keepertest.Bob, "10000"))
	response, err := keepertest.FactoryDenomSell(ctx, msgServer, keepertest.Bob, factoryDenomHash, factoryDenomHash, constants.KUSD, "10000", "", true)
	require.NoError(t, err)
	require.Equal(t, response.AmountReceivedNet, response.AmountReceivedGross-response.FeeReserve-response.FeePool)
	require.Equal(t, int64(10_000), response.AmountGivenGross)
	require.Equal(t, int64(10_000), response.AmountGivenNet)
	require.Equal(t, int64(9_900), response.AmountReceivedGross)
	require.Equal(t, int64(8_901), response.AmountReceivedNet)
	require.Equal(t, int64(9), response.FeeReserve)
	require.Equal(t, int64(990), response.FeePool)

	pool, _ := k.GetLiquidityPool(ctx, factoryDenomHash)
	require.Equal(t, int64(1_010_000), pool.FactoryDenomAmount.Int64())
	require.Equal(t, int64(992_089), pool.KCoinAmount.Int64())

	poolBalance2 := k.BankKeeper.SpendableCoins(ctx, poolAcc.GetAddress())
	poolBalanceFactory2 := poolBalance2.AmountOf(factoryDenomHash).Int64()
	poolBalanceKCoin2 := poolBalance2.AmountOf(constants.KUSD).Int64()

	require.Equal(t, poolBalanceFactory2-poolBalanceFactory1, response.AmountGivenGross)
	require.Equal(t, poolBalanceKCoin1-poolBalanceKCoin2, response.AmountReceivedNet)
}

func TestTrade6(t *testing.T) {
	k, msgServer, ctx := keepertest.SetupTokenfactoryMsgServer(t)

	factoryDenomHash, err := keepertest.CreateFactoryDenom(ctx, msgServer, keepertest.Alice, "testdenom", 6)
	require.NoError(t, err)
	require.NoError(t, keepertest.MintFactoryDenom(ctx, msgServer, keepertest.Alice, factoryDenomHash, keepertest.Alice, "1000000"))
	require.NoError(t, keepertest.CreatePool(ctx, msgServer, keepertest.Alice, factoryDenomHash, "1000000", constants.KUSD, "1000000", "0.1", 10))

	poolAcc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolFactoryLiquidity)
	poolBalance1 := k.BankKeeper.SpendableCoins(ctx, poolAcc.GetAddress())
	poolBalanceFactory1 := poolBalance1.AmountOf(factoryDenomHash).Int64()
	poolBalanceKCoin1 := poolBalance1.AmountOf(constants.KUSD).Int64()

	require.NoError(t, keepertest.MintFactoryDenom(ctx, msgServer, keepertest.Alice, factoryDenomHash, keepertest.Bob, "10000"))
	response, err := keepertest.FactoryDenomBuy(ctx, msgServer, keepertest.Bob, factoryDenomHash, constants.KUSD, factoryDenomHash, "1000", "", true)
	require.NoError(t, err)
	require.Equal(t, int64(1000), response.AmountReceivedNet)
	require.Equal(t, int64(1000), response.AmountReceivedGross)
	require.Equal(t, int64(1001), response.AmountGivenNet)
	require.Equal(t, int64(1102), response.AmountGivenGross)
	require.Equal(t, int64(100), response.FeePool)
	require.Equal(t, int64(1), response.FeeReserve)

	pool, _ := k.GetLiquidityPool(ctx, factoryDenomHash)
	require.Equal(t, int64(1001202), pool.KCoinAmount.Int64())
	require.Equal(t, int64(999000), pool.FactoryDenomAmount.Int64())

	poolBalance2 := k.BankKeeper.SpendableCoins(ctx, poolAcc.GetAddress())
	poolBalanceFactory2 := poolBalance2.AmountOf(factoryDenomHash).Int64()
	poolBalanceKCoin2 := poolBalance2.AmountOf(constants.KUSD).Int64()

	require.Equal(t, poolBalanceFactory1-poolBalanceFactory2, response.AmountReceivedNet)
	require.Equal(t, poolBalanceKCoin2-poolBalanceKCoin1, response.AmountGivenGross)
}
