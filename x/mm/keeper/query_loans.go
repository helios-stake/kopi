package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	"github.com/kopi-money/kopi/x/mm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetLoansByDenom(ctx context.Context, req *types.GetLoansByDenomQuery) (*types.GetLoansResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	cAsset, err := k.DenomKeeper.GetCAssetByBaseName(ctx, req.Denom)
	if err != nil {
		return nil, err
	}

	referenceDenom, err := k.DexKeeper.GetHighestUSDReference(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get reference denom: %w", err)
	}

	var (
		loanSum      = k.GetLoanSumWithDefault(ctx, req.Denom)
		utilityRate  = k.getUtilityRate(ctx, cAsset)
		interestRate = k.calculateInterestRate(ctx, utilityRate)
		loans        = []*types.DenomLoan{}

		amountBorrowedUSD math.LegacyDec
	)

	iterator := k.LoanIterator(ctx, cAsset.BaseDexDenom)
	for iterator.Valid() {
		keyValue := iterator.GetNextKeyValue()
		loan := keyValue.Value().Value()
		loanValue := k.getLoanValue(loanSum, *loan)

		amountBorrowedUSD, err = k.DexKeeper.GetValueIn(ctx, cAsset.BaseDexDenom, referenceDenom, loanValue)
		if err != nil {
			return nil, err
		}

		loans = append(loans, &types.DenomLoan{
			LoanIndex:         loan.Index,
			Address:           keyValue.Key(),
			AmountBorrowed:    loanValue.String(),
			AmountBorrowedUsd: amountBorrowedUSD.String(),
		})
	}

	return &types.GetLoansResponse{
		Loans:        loans,
		InterestRate: interestRate.String(),
	}, nil
}

func (k Keeper) GetLoansStats(ctx context.Context, _ *types.GetLoanStatsQuery) (*types.GetLoanStatsResponse, error) {
	referenceDenom, err := k.DexKeeper.GetHighestUSDReference(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get reference denom: %w", err)
	}

	var (
		acc                       = k.AccountKeeper.GetModuleAccount(ctx, types.PoolVault)
		vault                     = k.BankKeeper.SpendableCoins(ctx, acc.GetAddress())
		totalAvailableToBorrowUSD = math.LegacyZeroDec()
		totalLoanSumUSD           = math.LegacyZeroDec()
		loanStats                 = []*types.DenomLoanStat{}

		amountAvailable    math.Int
		amountAvailableUSD math.LegacyDec
		loanSumUSD         math.LegacyDec
	)

	for _, cAsset := range k.DenomKeeper.GetCAssets(ctx) {
		utilityRate := k.getUtilityRate(ctx, cAsset)
		interestRate := k.calculateInterestRate(ctx, utilityRate)

		amountAvailable = vault.AmountOf(cAsset.BaseDexDenom)
		amountAvailableUSD, err = k.DexKeeper.GetValueIn(ctx, cAsset.BaseDexDenom, referenceDenom, amountAvailable.ToLegacyDec())
		if err != nil {
			return nil, err
		}

		loanSum := k.GetLoanSumWithDefault(ctx, cAsset.BaseDexDenom).LoanSum
		loanSumUSD, err = k.DexKeeper.GetValueIn(ctx, cAsset.BaseDexDenom, referenceDenom, loanSum)
		if err != nil {
			return nil, err
		}

		totalAvailableToBorrowUSD = totalAvailableToBorrowUSD.Add(amountAvailableUSD)
		totalLoanSumUSD = totalAvailableToBorrowUSD.Add(loanSumUSD)

		loanStats = append(loanStats, &types.DenomLoanStat{
			Denom:                cAsset.BaseDexDenom,
			UtilityRate:          utilityRate.String(),
			InterestRate:         interestRate.String(),
			AvailableToBorrow:    amountAvailable.String(),
			AvailableToBorrowUsd: amountAvailableUSD.String(),
			LoanSum:              loanSum.String(),
			LoanSumUsd:           loanSumUSD.String(),
		})
	}

	return &types.GetLoanStatsResponse{
		LoanStats:                 loanStats,
		TotalAvailableToBorrowUsd: totalAvailableToBorrowUSD.String(),
		TotalLoanSumUsd:           totalLoanSumUSD.String(),
	}, nil
}

func (k Keeper) GetUserLoans(ctx context.Context, req *types.GetUserLoansQuery) (*types.GetUserLoansResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	referenceDenom, err := k.DexKeeper.GetHighestUSDReference(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get reference denom: %w", err)
	}

	var (
		addr      = k.AccountKeeper.GetModuleAccount(ctx, types.PoolVault)
		vault     = k.BankKeeper.SpendableCoins(ctx, addr.GetAddress())
		userLoans = []*types.UserLoanStat{}

		amountAvailable    math.LegacyDec
		amountAvailableUSD math.LegacyDec
		amountBorrowedUSD  math.LegacyDec
	)

	for _, cAsset := range k.DenomKeeper.GetCAssets(ctx) {
		utilityRate := k.getUtilityRate(ctx, cAsset)
		interestRate := k.calculateInterestRate(ctx, utilityRate)

		loanSum := k.GetLoanSumWithDefault(ctx, cAsset.BaseDexDenom)

		loan, has := k.loans.Get(ctx, cAsset.BaseDexDenom, req.Address)
		loanValue := math.LegacyZeroDec()
		if has {
			loanValue = k.getLoanValue(loanSum, loan)
		}

		vaultAmount := vault.AmountOf(cAsset.BaseDexDenom)
		amountAvailable, err = k.CalculateBorrowableAmount(ctx, req.Address, cAsset.BaseDexDenom)
		if err != nil {
			return nil, err
		}

		amountAvailable = math.LegacyMinDec(vaultAmount.ToLegacyDec(), amountAvailable)

		amountAvailableUSD, err = k.DexKeeper.GetValueIn(ctx, cAsset.BaseDexDenom, referenceDenom, amountAvailable)
		if err != nil {
			return nil, err
		}

		amountBorrowedUSD, err = k.DexKeeper.GetValueIn(ctx, cAsset.BaseDexDenom, referenceDenom, loanValue)
		if err != nil {
			return nil, err
		}

		userLoans = append(userLoans, &types.UserLoanStat{
			Denom:              cAsset.BaseDexDenom,
			AmountBorrowed:     loanValue.String(),
			AmountBorrowedUsd:  amountBorrowedUSD.String(),
			AmountAvailable:    amountAvailable.String(),
			AmountAvailableUsd: amountAvailableUSD.String(),
			InterestRate:       interestRate.String(),
		})
	}

	return &types.GetUserLoansResponse{
		UserLoans: userLoans,
	}, nil
}

func (k Keeper) GetUserDenomLoan(ctx context.Context, req *types.GetUserDenomLoanQuery) (*types.GetUserDenomLoanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	cAsset, err := k.DenomKeeper.GetCAssetByBaseName(ctx, req.Denom)
	if err != nil {
		return nil, err
	}

	loanValue := k.GetLoanValue(ctx, cAsset.BaseDexDenom, req.Address)
	amountUSD, err := k.DexKeeper.GetValueInUSD(ctx, cAsset.BaseDexDenom, loanValue)
	if err != nil {
		return nil, err
	}

	return &types.GetUserDenomLoanResponse{
		Amount:    loanValue.String(),
		AmountUsd: amountUSD.String(),
	}, nil
}

func (k Keeper) GetNumLoans(ctx context.Context, _ *types.GetNumLoansQuery) (*types.GetNumLoansResponse, error) {
	return &types.GetNumLoansResponse{
		Num: int64(k.GetLoansNum(ctx)),
	}, nil
}

func (k Keeper) GetValueLoans(ctx context.Context, _ *types.GetValueLoansQuery) (*types.GetValueLoansResponse, error) {
	referenceDenom, err := k.DexKeeper.GetHighestUSDReference(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get reference denom: %w", err)
	}

	valueUSD := math.LegacyZeroDec()
	for _, cAsset := range k.DenomKeeper.GetCAssets(ctx) {
		loanSum := k.GetLoanSumWithDefault(ctx, cAsset.BaseDexDenom)

		var value math.LegacyDec
		value, err = k.DexKeeper.GetValueIn(ctx, cAsset.BaseDexDenom, referenceDenom, loanSum.LoanSum)
		if err != nil {
			return nil, fmt.Errorf("could not get value in usd: %w", err)
		}

		valueUSD = valueUSD.Add(value)
	}

	return &types.GetValueLoansResponse{
		Value: valueUSD.String(),
	}, nil
}

func (k Keeper) GetNumAddressLoans(ctx context.Context, req *types.GetNumAddressLoansQuery) (*types.GetNumAddressLoansResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.GetNumAddressLoansResponse{
		Amount: int64(k.GetLoansNumForAddress(ctx, req.Address)),
	}, nil
}

func (k Keeper) GetAvailableToBorrow(ctx context.Context, req *types.GetAvailableToBorrowRequest) (*types.GetAvailableToBorrowResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	amount, err := k.CalcAvailableToBorrow(ctx, req.Address, req.Denom)
	if err != nil {
		return nil, fmt.Errorf("could not calculate available amount to borrow: %w", err)
	}

	amountUSD, err := k.DexKeeper.GetValueInUSD(ctx, req.Denom, amount.ToLegacyDec())
	if err != nil {
		return nil, fmt.Errorf("could not convert amount to usd: %w", err)
	}

	return &types.GetAvailableToBorrowResponse{
		Amount:    amount.String(),
		AmountUsd: amountUSD.String(),
	}, nil
}
