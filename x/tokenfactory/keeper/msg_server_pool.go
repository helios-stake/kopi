package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/x/tokenfactory/types"
)

func (k msgServer) CreatePool(ctx context.Context, msg *types.MsgCreatePool) (*types.Void, error) {
	factoryDenom, has := k.GetDenomByFullName(ctx, msg.FullFactoryDenomName)
	if !has {
		return nil, types.ErrDenomDoesntExists
	}

	if factoryDenom.Admin != msg.Creator {
		return nil, types.ErrIncorrectAdmin
	}

	pool, has := k.liquidityPools.Get(ctx, factoryDenom.FullName)
	if has {
		return nil, types.ErrPoolAlreadyExists
	}

	if msg.KCoin == "" {
		return nil, types.ErrEmptyKCoin
	}

	if !k.DenomKeeper.IsKCoin(ctx, msg.KCoin) {
		return nil, types.ErrNoKCoin
	}

	amountFactory, ok := math.NewIntFromString(msg.FactoryDenomAmount)
	if !ok {
		return nil, fmt.Errorf("invalid factory denom amount: %v", msg.FactoryDenomAmount)
	}

	kCoinAmount, ok := math.NewIntFromString(msg.KCoinAmount)
	if !ok {
		return nil, fmt.Errorf("invalid other kcoin amount: %v", msg.KCoinAmount)
	}

	poolFee, err := math.LegacyNewDecFromStr(msg.PoolFee)
	if err != nil {
		return nil, types.ErrInvalidFeeFormat
	}

	if poolFee.GT(math.LegacyOneDec()) {
		return nil, types.ErrPoolFeeToLarge
	}

	pool = types.LiquidityPool{
		FactoryDenomAmount: amountFactory,
		KCoin:              msg.KCoin,
		KCoinAmount:        kCoinAmount,
		PoolFee:            poolFee,
		UnlockBlocks:       msg.UnlockBlocks,
	}

	k.liquidityPools.Set(ctx, factoryDenom.FullName, pool)

	coins := sdk.NewCoins(
		sdk.NewCoin(factoryDenom.FullName, amountFactory),
		sdk.NewCoin(msg.KCoin, kCoinAmount),
	)

	adminAcc, _ := sdk.AccAddressFromBech32(msg.Creator)
	if err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, adminAcc, types.PoolFactoryLiquidity, coins); err != nil {
		return nil, err
	}

	provider := types.ProviderShare{Share: math.LegacyOneDec()}
	k.liquidityProviderShares.Set(ctx, factoryDenom.FullName, msg.Creator, provider)

	return &types.Void{}, nil
}

func (k Keeper) checkLiquidity(ctx context.Context, creator, fullName, amount string) (types.FactoryDenom, types.LiquidityPool, math.Int, math.Int, error) {
	factoryDenom, has := k.GetDenomByFullName(ctx, fullName)
	if !has {
		return types.FactoryDenom{}, types.LiquidityPool{}, math.Int{}, math.Int{}, types.ErrDenomDoesntExists
	}

	if factoryDenom.Admin != creator {
		return types.FactoryDenom{}, types.LiquidityPool{}, math.Int{}, math.Int{}, types.ErrIncorrectAdmin
	}

	pool, has := k.liquidityPools.Get(ctx, factoryDenom.FullName)
	if !has {
		return types.FactoryDenom{}, types.LiquidityPool{}, math.Int{}, math.Int{}, types.ErrPoolDoesNotExist
	}

	amountFactory, ok := math.NewIntFromString(amount)
	if !ok {
		return types.FactoryDenom{}, types.LiquidityPool{}, math.Int{}, math.Int{}, fmt.Errorf("invalid factory denom amount: %v", amount)
	}

	poolRatio := getPoolRatio(pool)
	amountOtherDenom := amountFactory.ToLegacyDec().Mul(poolRatio).TruncateInt()

	return factoryDenom, pool, amountFactory, amountOtherDenom, nil
}

func (k msgServer) AddLiquidity(ctx context.Context, msg *types.MsgAddLiquidity) (*types.Void, error) {
	factoryDenom, pool, amountFactory, amountKCoin, err := k.checkLiquidity(ctx, msg.Creator, msg.FullFactoryDenomName, msg.FactoryDenomAmount)
	if err != nil {
		return nil, err
	}

	pool.FactoryDenomAmount = pool.FactoryDenomAmount.Add(amountFactory)
	pool.KCoinAmount = pool.KCoinAmount.Add(amountKCoin)
	k.liquidityPools.Set(ctx, factoryDenom.FullName, pool)

	acc, _ := sdk.AccAddressFromBech32(msg.Creator)
	coins := sdk.NewCoins(
		sdk.NewCoin(factoryDenom.FullName, amountFactory),
		sdk.NewCoin(pool.KCoin, amountKCoin),
	)

	if err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, acc, types.PoolFactoryLiquidity, coins); err != nil {
		return nil, err
	}

	_ = k.updateLiquidityShare(ctx, factoryDenom, pool, amountFactory, msg.Creator)

	return &types.Void{}, nil
}

func (k msgServer) UnlockLiquidity(ctx context.Context, msg *types.MsgUnlockLiquidity) (*types.Void, error) {
	factoryDenom, pool, amountFactory, amountKCoin, err := k.checkLiquidity(ctx, msg.Creator, msg.FullFactoryDenomName, msg.FactoryDenomAmount)
	if err != nil {
		return nil, err
	}

	pool.FactoryDenomAmount = pool.FactoryDenomAmount.Sub(amountFactory)
	pool.KCoinAmount = pool.KCoinAmount.Sub(amountKCoin)
	k.liquidityPools.Set(ctx, factoryDenom.FullName, pool)

	if pool.FactoryDenomAmount.IsNegative() || pool.KCoinAmount.IsNegative() {
		return nil, types.ErrNegativeLiquidity
	}

	coins := sdk.NewCoins(
		sdk.NewCoin(factoryDenom.FullName, amountFactory),
		sdk.NewCoin(pool.KCoin, amountKCoin),
	)

	if err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.PoolFactoryLiquidity, types.PoolUnlocking, coins); err != nil {
		return nil, err
	}

	if err = k.updateLiquidityShare(ctx, factoryDenom, pool, amountFactory, msg.Creator); err != nil {
		return nil, err
	}

	k.SetLiquidityUnlocking(ctx, types.LiquidityUnlocking{
		Index:              0,
		Address:            msg.Creator,
		FactoryDenomHash:   factoryDenom.FullName,
		FactoryDenomAmount: amountFactory,
		KCoin:              pool.KCoin,
		KCoinAmount:        amountKCoin,
		CreatedAt:          sdk.UnwrapSDKContext(ctx).BlockHeight(),
	})

	return &types.Void{}, nil
}

func (k msgServer) UpdateLiquidityPoolSettings(ctx context.Context, msg *types.MsgUpdateLiquidityPoolSettings) (*types.Void, error) {
	factoryDenom, has := k.GetDenomByFullName(ctx, msg.FullFactoryDenomName)
	if !has {
		return nil, types.ErrDenomDoesntExists
	}

	if factoryDenom.Admin != msg.Creator {
		return nil, types.ErrIncorrectAdmin
	}

	pool, has := k.liquidityPools.Get(ctx, factoryDenom.FullName)
	if !has {
		return nil, types.ErrPoolDoesNotExist
	}

	var err error
	pool.PoolFee, err = math.LegacyNewDecFromStr(msg.PoolFee)
	if err != nil {
		return nil, types.ErrInvalidFeeFormat
	}

	if pool.PoolFee.LT(math.LegacyZeroDec()) {
		return nil, types.ErrInvalidNegativeFee
	}

	if pool.PoolFee.GT(math.LegacyOneDec()) {
		return nil, types.ErrPoolFeeToLarge
	}

	pool.UnlockBlocks = msg.UnlockBlocks

	k.SetLiquidityPool(ctx, factoryDenom.FullName, pool)
	return &types.Void{}, nil
}
