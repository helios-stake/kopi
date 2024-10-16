package keeper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/x/mm/types"
)

func (k msgServer) Borrow(ctx context.Context, msg *types.MsgBorrow) (*types.Void, error) {
	amountStr := strings.ReplaceAll(msg.Amount, ",", "")
	amount, ok := math.NewIntFromString(amountStr)
	if !ok {
		return nil, types.ErrInvalidAmountFormat
	}

	if amount.LT(math.ZeroInt()) {
		return nil, types.ErrNegativeAmount
	}

	if amount.IsZero() {
		return nil, types.ErrZeroAmount
	}

	address, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, types.ErrInvalidAddress
	}

	_, _, err = k.Keeper.Borrow(ctx, address, msg.Denom, amount)
	if err != nil {
		return nil, fmt.Errorf("could not execute borrow: %w", err)
	}

	return &types.Void{}, nil
}

func (k Keeper) Borrow(ctx context.Context, address sdk.AccAddress, denom string, borrowAmount math.Int) (math.Int, math.Int, error) {
	cAsset, err := k.DenomKeeper.GetCAssetByBaseName(ctx, denom)
	if err != nil {
		return math.Int{}, math.Int{}, types.ErrInvalidDepositDenom
	}

	if borrowAmount.LT(math.ZeroInt()) {
		return math.Int{}, math.Int{}, types.ErrNegativeAmount
	}

	if borrowAmount.IsZero() {
		return math.Int{}, math.Int{}, types.ErrZeroAmount
	}

	vaultAcc := k.AccountKeeper.GetModuleAccount(ctx, types.PoolVault)
	vaultBalance := k.BankKeeper.SpendableCoin(ctx, vaultAcc.GetAddress(), cAsset.BaseDexDenom).Amount

	if vaultBalance.LT(borrowAmount) {
		k.Logger().Error(fmt.Sprintf("%v < %v %v", vaultBalance.String(), borrowAmount.String(), cAsset.BaseDexDenom))
		return math.Int{}, math.Int{}, types.ErrNotEnoughFundsInVault
	}

	borrowableAmount, err := k.CalculateBorrowableAmount(ctx, address.String(), denom)
	if err != nil {
		return math.Int{}, math.Int{}, err
	}

	if borrowableAmount.TruncateInt().LT(borrowAmount) {
		return math.Int{}, math.Int{}, types.ErrCollateralBorrowLimitExceeded
	}

	if cAsset.MinimumLoanSize.GT(math.ZeroInt()) && borrowAmount.LT(cAsset.MinimumLoanSize) {
		return math.Int{}, math.Int{}, types.ErrLoanSizeTooSmall
	}

	if k.checkBorrowLimitExceeded(ctx, cAsset, borrowAmount) {
		return math.Int{}, math.Int{}, types.ErrBorrowLimitExceeded
	}

	loanIndex, _ := k.updateLoan(ctx, denom, address.String(), borrowAmount.ToLegacyDec())

	coins := sdk.NewCoins(sdk.NewCoin(denom, borrowAmount))
	if err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.PoolVault, address, coins); err != nil {
		return math.Int{}, math.Int{}, err
	}

	loanValue := k.GetLoanValue(ctx, denom, address.String()).TruncateInt()

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(
		sdk.NewEvent("funds_borrowed",
			sdk.Attribute{Key: "address", Value: address.String()},
			sdk.Attribute{Key: "denom", Value: denom},
			sdk.Attribute{Key: "amount", Value: borrowAmount.String()},
			sdk.Attribute{Key: "index", Value: strconv.Itoa(int(loanIndex))},
			sdk.Attribute{Key: "borrowed_amount", Value: loanValue.String()},
		),
	)

	return borrowAmount, loanValue, nil
}

func (k msgServer) RepayLoan(ctx context.Context, msg *types.MsgRepayLoan) (*types.Void, error) {
	loanValue := k.GetLoanValue(ctx, msg.Denom, msg.Creator)
	if loanValue.IsZero() {
		return nil, types.ErrNoLoanFound
	}

	if err := k.Repay(ctx, msg.Denom, msg.Creator, loanValue.Ceil().TruncateInt()); err != nil {
		return nil, err
	}

	return &types.Void{}, nil
}

func (k msgServer) PartiallyRepayLoan(ctx context.Context, msg *types.MsgPartiallyRepayLoan) (*types.Void, error) {
	if _, err := k.DenomKeeper.GetCAssetByBaseName(ctx, msg.Denom); err != nil {
		return nil, types.ErrInvalidDepositDenom
	}

	_, found := k.loans.Get(ctx, msg.Denom, msg.Creator)
	if !found {
		return nil, types.ErrNoLoanFound
	}

	amountStr := strings.ReplaceAll(msg.Amount, ",", "")
	repayAmount, ok := math.NewIntFromString(amountStr)
	if !ok {
		return nil, types.ErrInvalidAmountFormat
	}

	if repayAmount.LT(math.ZeroInt()) {
		return nil, types.ErrNegativeAmount
	}

	if repayAmount.IsZero() {
		return nil, types.ErrZeroAmount
	}

	if err := k.Repay(ctx, msg.Denom, msg.Creator, repayAmount); err != nil {
		return nil, err
	}

	return &types.Void{}, nil
}

func (k Keeper) Repay(ctx context.Context, denom, address string, repayAmount math.Int) error {
	acc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return types.ErrInvalidAddress
	}

	loanValue := k.GetLoanValue(ctx, denom, address)
	repayAmount = math.MinInt(loanValue.Ceil().TruncateInt(), repayAmount)

	if k.BankKeeper.SpendableCoin(ctx, acc, denom).Amount.LT(repayAmount) {
		return types.ErrNotEnoughFunds
	}

	loanIndex, removed := k.updateLoan(ctx, denom, address, repayAmount.ToLegacyDec().Neg())

	coins := sdk.NewCoins(sdk.NewCoin(denom, repayAmount))
	if err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, acc, types.PoolVault, coins); err != nil {
		return err
	}

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(
		sdk.NewEvent("loan_repaid",
			sdk.Attribute{Key: "address", Value: address},
			sdk.Attribute{Key: "denom", Value: denom},
			sdk.Attribute{Key: "index", Value: strconv.Itoa(int(loanIndex))},
			sdk.Attribute{Key: "amount", Value: repayAmount.String()},
		),
	)

	if removed {
		sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(
			sdk.NewEvent("loan_removed",
				sdk.Attribute{Key: "address", Value: address},
				sdk.Attribute{Key: "denom", Value: denom},
				sdk.Attribute{Key: "index", Value: strconv.Itoa(int(loanIndex))},
			),
		)
	}

	return nil
}
