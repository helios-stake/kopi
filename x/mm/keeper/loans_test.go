package keeper_test

import (
	"github.com/kopi-money/kopi/constants"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/kopi-money/kopi/testutil/keeper"
	"github.com/kopi-money/kopi/x/mm/types"
	"github.com/stretchr/testify/require"
)

func TestLoans1(t *testing.T) {
	_, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	require.Error(t, keepertest.Borrow(ctx, msg, &types.MsgBorrow{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "100",
	}))
}

func TestLoans2(t *testing.T) {
	_, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	_, _ = msg.AddDeposit(ctx, &types.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "100",
	})

	require.Error(t, keepertest.Borrow(ctx, msg, &types.MsgBorrow{
		Creator: keepertest.Bob,
		Denom:   constants.KUSD,
		Amount:  "100",
	}))
}

func TestLoans3(t *testing.T) {
	_, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	_, _ = msg.AddDeposit(ctx, &types.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "100",
	})

	require.NoError(t, keepertest.AddCollateral(ctx, msg, &types.MsgAddCollateral{
		Creator: keepertest.Bob,
		Denom:   constants.BaseCurrency,
		Amount:  "1000",
	}))

	require.Error(t, keepertest.Borrow(ctx, msg, &types.MsgBorrow{
		Creator: keepertest.Bob,
		Denom:   constants.KUSD,
		Amount:  "101",
	}))
}

func TestLoans4(t *testing.T) {
	k, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	acc, _ := sdk.AccAddressFromBech32(keepertest.Bob)
	coins := k.BankKeeper.SpendableCoins(ctx, acc)
	found, balance1 := coins.Find(constants.KUSD)
	require.True(t, found)

	_, _ = msg.AddDeposit(ctx, &types.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "10000",
	})

	require.NoError(t, keepertest.AddCollateral(ctx, msg, &types.MsgAddCollateral{
		Creator: keepertest.Bob,
		Denom:   constants.BaseCurrency,
		Amount:  "100000",
	}))

	require.NoError(t, keepertest.Borrow(ctx, msg, &types.MsgBorrow{
		Creator: keepertest.Bob,
		Denom:   constants.KUSD,
		Amount:  "1000",
	}))

	coins = k.BankKeeper.SpendableCoins(ctx, acc)
	_, balance2 := coins.Find(constants.KUSD)
	diff := balance2.Amount.Sub(balance1.Amount)
	require.Equal(t, math.NewInt(1000), diff)

	iterator := k.LoanIterator(ctx, constants.KUSD)
	loans := iterator.GetAll()
	require.Equal(t, 1, len(loans))
	require.Equal(t, keepertest.Bob, loans[0].Address)

	loanValue := k.GetLoanValue(ctx, constants.KUSD, keepertest.Bob)
	require.Equal(t, math.LegacyNewDec(1000), loanValue)
}

func TestLoans5(t *testing.T) {
	k, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	acc, _ := sdk.AccAddressFromBech32(keepertest.Bob)
	balance1 := k.BankKeeper.SpendableCoin(ctx, acc, constants.KUSD).Amount

	_, err := msg.AddDeposit(ctx, &types.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "10000",
	})
	require.NoError(t, err)

	require.NoError(t, keepertest.AddCollateral(ctx, msg, &types.MsgAddCollateral{
		Creator: keepertest.Bob,
		Denom:   constants.BaseCurrency,
		Amount:  "100000",
	}))

	require.NoError(t, keepertest.Borrow(ctx, msg, &types.MsgBorrow{
		Creator: keepertest.Bob,
		Denom:   constants.KUSD,
		Amount:  "1000",
	}))

	balance2 := k.BankKeeper.SpendableCoin(ctx, acc, constants.KUSD).Amount

	require.NoError(t, keepertest.PartiallyRepayLoan(ctx, msg, &types.MsgPartiallyRepayLoan{
		Creator: keepertest.Bob,
		Denom:   constants.KUSD,
		Amount:  "1",
	}))

	balance3 := k.BankKeeper.SpendableCoin(ctx, acc, constants.KUSD).Amount
	require.True(t, balance3.LT(balance2))

	diff := balance3.Sub(balance1)
	require.Equal(t, math.NewInt(999), diff)

	iterator := k.LoanIterator(ctx, constants.KUSD)
	loans := iterator.GetAll()
	require.Equal(t, 1, len(loans))
	require.Equal(t, keepertest.Bob, loans[0].Address)

	loanValue := k.GetLoanValue(ctx, constants.KUSD, keepertest.Bob)
	require.Equal(t, math.LegacyNewDec(999), loanValue)
}

func TestLoans6(t *testing.T) {
	k, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	acc, _ := sdk.AccAddressFromBech32(keepertest.Bob)
	coins := k.BankKeeper.SpendableCoins(ctx, acc)
	found, balance1 := coins.Find(constants.KUSD)
	require.True(t, found)

	_, _ = msg.AddDeposit(ctx, &types.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "10000",
	})

	require.NoError(t, keepertest.AddCollateral(ctx, msg, &types.MsgAddCollateral{
		Creator: keepertest.Bob,
		Denom:   constants.BaseCurrency,
		Amount:  "100000",
	}))

	require.NoError(t, keepertest.Borrow(ctx, msg, &types.MsgBorrow{
		Creator: keepertest.Bob,
		Denom:   constants.KUSD,
		Amount:  "1000",
	}))

	require.NoError(t, keepertest.PartiallyRepayLoan(ctx, msg, &types.MsgPartiallyRepayLoan{
		Creator: keepertest.Bob,
		Denom:   constants.KUSD,
		Amount:  "1001",
	}))

	coins = k.BankKeeper.SpendableCoins(ctx, acc)
	_, balance2 := coins.Find(constants.KUSD)
	require.Equal(t, balance1.Amount, balance2.Amount)

	iterator := k.LoanIterator(ctx, constants.KUSD)
	require.Equal(t, 0, len(iterator.GetAll()))
}

func TestLoans7(t *testing.T) {
	k, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	acc, _ := sdk.AccAddressFromBech32(keepertest.Bob)
	coins := k.BankKeeper.SpendableCoins(ctx, acc)
	found, balance1 := coins.Find(constants.KUSD)
	require.True(t, found)

	_, err := msg.AddDeposit(ctx, &types.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "10000",
	})

	require.NoError(t, err)

	require.NoError(t, keepertest.AddCollateral(ctx, msg, &types.MsgAddCollateral{
		Creator: keepertest.Bob,
		Denom:   constants.BaseCurrency,
		Amount:  "100000",
	}))

	require.NoError(t, err)

	require.NoError(t, keepertest.Borrow(ctx, msg, &types.MsgBorrow{
		Creator: keepertest.Bob,
		Denom:   constants.KUSD,
		Amount:  "1000",
	}))

	require.NoError(t, keepertest.RepayLoan(ctx, msg, &types.MsgRepayLoan{
		Creator: keepertest.Bob,
		Denom:   constants.KUSD,
	}))

	coins = k.BankKeeper.SpendableCoins(ctx, acc)
	_, balance2 := coins.Find(constants.KUSD)
	require.Equal(t, balance1.Amount, balance2.Amount)

	iterator := k.LoanIterator(ctx, constants.KUSD)
	require.Equal(t, 0, len(iterator.GetAll()))
}

func TestLoans8(t *testing.T) {
	k, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	_, err := msg.AddDeposit(ctx, &types.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "10000",
	})

	require.NoError(t, err)

	borrowable, err := k.CalcAvailableToBorrow(ctx, keepertest.Alice, constants.KUSD)
	require.NoError(t, err)
	require.Equal(t, int64(0), borrowable.Int64())
}

func TestLoans9(t *testing.T) {
	k, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	require.NoError(t, keepertest.AddCollateral(ctx, msg, &types.MsgAddCollateral{
		Creator: keepertest.Alice,
		Denom:   constants.BaseCurrency,
		Amount:  "10000",
	}))

	_, err := msg.AddDeposit(ctx, &types.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "10000",
	})

	require.NoError(t, err)

	borrowable, err := k.CalcAvailableToBorrow(ctx, keepertest.Alice, constants.KUSD)
	require.NoError(t, err)
	require.Equal(t, int64(1250), borrowable.Int64())
}

func TestLoans10(t *testing.T) {
	k, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	require.NoError(t, keepertest.AddCollateral(ctx, msg, &types.MsgAddCollateral{
		Creator: keepertest.Alice,
		Denom:   constants.BaseCurrency,
		Amount:  "10000000",
	}))

	_, err := msg.AddDeposit(ctx, &types.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "10000",
	})

	require.NoError(t, err)

	withdrawable, err := k.CalcAvailableToBorrow(ctx, keepertest.Alice, constants.KUSD)
	require.NoError(t, err)
	require.Equal(t, int64(10000), withdrawable.Int64())
}

func TestLoans11(t *testing.T) {
	_, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	require.NoError(t, keepertest.AddCollateral(ctx, msg, &types.MsgAddCollateral{
		Creator: keepertest.Alice,
		Denom:   constants.BaseCurrency,
		Amount:  "100000000",
	}))

	_, err := msg.AddDeposit(ctx, &types.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "10000",
	})
	require.NoError(t, err)

	require.NoError(t, keepertest.Borrow(ctx, msg, &types.MsgBorrow{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "7500",
	}))
}

func TestLoans12(t *testing.T) {
	k, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	require.NoError(t, keepertest.AddCollateral(ctx, msg, &types.MsgAddCollateral{
		Creator: keepertest.Alice,
		Denom:   constants.BaseCurrency,
		Amount:  "2500000",
	}))

	_, err := msg.AddDeposit(ctx, &types.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "1000000",
	})
	require.NoError(t, err)

	borrowable, err := k.CalcAvailableToBorrow(ctx, keepertest.Alice, constants.KUSD)
	require.NoError(t, err)

	require.NoError(t, keepertest.Borrow(ctx, msg, &types.MsgBorrow{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  borrowable.String(),
	}))
}

func TestLoans13(t *testing.T) {
	k, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	require.NoError(t, keepertest.AddCollateral(ctx, msg, &types.MsgAddCollateral{
		Creator: keepertest.Alice,
		Denom:   constants.BaseCurrency,
		Amount:  "10000000",
	}))

	_, err := msg.AddDeposit(ctx, &types.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "100000",
	})
	require.NoError(t, err)

	require.NoError(t, keepertest.AddCollateral(ctx, msg, &types.MsgAddCollateral{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "100000",
	}))

	borrowable1, err := k.CalcAvailableToBorrow(ctx, keepertest.Alice, constants.KUSD)
	require.NoError(t, err)

	require.NoError(t, keepertest.Borrow(ctx, msg, &types.MsgBorrow{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "1000",
	}))

	borrowable2, err := k.CalcAvailableToBorrow(ctx, keepertest.Alice, constants.KUSD)
	require.NoError(t, err)

	borrowableInt1 := borrowable1.Int64()
	borrowableInt2 := borrowable2.Int64()

	require.Less(t, borrowableInt2, borrowableInt1)
}

func TestLoans14(t *testing.T) {
	k, _, msg, ctx := keepertest.SetupMMMsgServer(t)

	require.NoError(t, keepertest.AddCollateral(ctx, msg, &types.MsgAddCollateral{
		Creator: keepertest.Alice,
		Denom:   constants.BaseCurrency,
		Amount:  "10000000",
	}))

	require.NoError(t, keepertest.AddCollateral(ctx, msg, &types.MsgAddCollateral{
		Creator: keepertest.Bob,
		Denom:   constants.BaseCurrency,
		Amount:  "10000000",
	}))

	_, err := msg.AddDeposit(ctx, &types.MsgAddDeposit{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "100000",
	})
	require.NoError(t, err)

	require.NoError(t, keepertest.Borrow(ctx, msg, &types.MsgBorrow{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
		Amount:  "9000",
	}))

	loan, found := k.LoadLoan(ctx, constants.KUSD, keepertest.Alice)
	require.True(t, found)
	require.Equal(t, math.LegacyNewDec(9000), loan.Weight)

	loanSum := k.GetLoanSumWithDefault(ctx, constants.KUSD)
	require.Equal(t, math.LegacyNewDec(9000), loanSum.LoanSum)
	require.Equal(t, math.LegacyNewDec(9000), loanSum.WeightSum)

	require.NoError(t, keepertest.Borrow(ctx, msg, &types.MsgBorrow{
		Creator: keepertest.Bob,
		Denom:   constants.KUSD,
		Amount:  "1000",
	}))

	loan, found = k.LoadLoan(ctx, constants.KUSD, keepertest.Bob)
	require.True(t, found)
	require.Equal(t, math.LegacyNewDec(1000), loan.Weight)

	loanSum = k.GetLoanSumWithDefault(ctx, constants.KUSD)
	require.Equal(t, math.LegacyNewDec(10000), loanSum.LoanSum)
	require.Equal(t, math.LegacyNewDec(10000), loanSum.WeightSum)

	require.NoError(t, keepertest.RepayLoan(ctx, msg, &types.MsgRepayLoan{
		Creator: keepertest.Bob,
		Denom:   constants.KUSD,
	}))

	_, found = k.LoadLoan(ctx, constants.KUSD, keepertest.Bob)
	require.True(t, !found)

	loanSum = k.GetLoanSumWithDefault(ctx, constants.KUSD)
	require.Equal(t, math.LegacyNewDec(9000), loanSum.LoanSum)
	require.Equal(t, math.LegacyNewDec(9000), loanSum.WeightSum)

	require.NoError(t, keepertest.RepayLoan(ctx, msg, &types.MsgRepayLoan{
		Creator: keepertest.Alice,
		Denom:   constants.KUSD,
	}))

	_, found = k.LoadLoan(ctx, constants.KUSD, keepertest.Alice)
	require.True(t, !found)

	loanSum = k.GetLoanSumWithDefault(ctx, constants.KUSD)
	require.Equal(t, math.LegacyNewDec(0), loanSum.LoanSum)
	require.Equal(t, math.LegacyNewDec(0), loanSum.WeightSum)
}
