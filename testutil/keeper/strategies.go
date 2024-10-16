package keeper

import (
	"context"
	"cosmossdk.io/core/address"
	"github.com/kopi-money/kopi/constants"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	dexkeeper "github.com/kopi-money/kopi/x/dex/keeper"
	dextypes "github.com/kopi-money/kopi/x/dex/types"
	mmkeeper "github.com/kopi-money/kopi/x/mm/keeper"
	mmtypes "github.com/kopi-money/kopi/x/mm/types"
	"github.com/kopi-money/kopi/x/strategies/types"
	"github.com/stretchr/testify/require"

	"github.com/kopi-money/kopi/cache"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kopi-money/kopi/x/strategies/keeper"
)

type DummyBlockspeedKeeper struct{}

func (d DummyBlockspeedKeeper) GetSecondsPerBlock(ctx context.Context) math.LegacyDec {
	//TODO implement me
	panic("implement me")
}

func (d DummyBlockspeedKeeper) BlocksPerYear(_ context.Context) (math.LegacyDec, error) {
	return math.LegacyNewDec(constants.SecondsPerYear).Quo(math.LegacyNewDec(2)), nil
}

type DummyDistrubtionKeeper struct {
}

func (d DummyDistrubtionKeeper) CalculateDelegationRewards(ctx context.Context, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins, err error) {
	//TODO implement me
	panic("implement me")
}

func (d DummyDistrubtionKeeper) IncrementValidatorPeriod(ctx context.Context, val stakingtypes.ValidatorI) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (d DummyDistrubtionKeeper) WithdrawDelegationRewards(_ context.Context, _ sdk.AccAddress, _ sdk.ValAddress) (sdk.Coins, error) {
	panic("implement me")
}

type DummyStakingKeeper struct {
}

func (d DummyStakingKeeper) GetDelegatorDelegations(ctx context.Context, delegator sdk.AccAddress, maxRetrieve uint16) (delegations []stakingtypes.Delegation, err error) {
	//TODO implement me
	panic("implement me")
}

func (d DummyStakingKeeper) Validator(ctx context.Context, address sdk.ValAddress) (stakingtypes.ValidatorI, error) {
	//TODO implement me
	panic("implement me")
}

func (d DummyStakingKeeper) ValidatorAddressCodec() address.Codec {
	//TODO implement me
	panic("implement me")
}

func (d DummyStakingKeeper) IterateDelegations(_ context.Context, _ sdk.AccAddress, _ func(index int64, delegation stakingtypes.DelegationI) (stop bool)) error {
	panic("implement me")
}

func (d DummyStakingKeeper) Delegate(_ context.Context, _ sdk.AccAddress, _ math.Int, _ stakingtypes.BondStatus, _ stakingtypes.Validator, _ bool) (math.LegacyDec, error) {
	panic("implement me")
}

func (d DummyStakingKeeper) GetBondedValidatorsByPower(_ context.Context) ([]stakingtypes.Validator, error) {
	panic("implement me")
}

func StrategiesKeeper(t *testing.T) (keeper.Keeper, dexkeeper.Keeper, mmkeeper.Keeper, context.Context) {
	dexKeeper, mmKeeper, ctx, keys := MmKeeperKeys(t)

	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	strategiesKeeper := keeper.NewKeeper(
		keys.cdc,
		runtime.NewKVStoreService(keys.str),
		log.NewNopLogger(),

		mmKeeper.AccountKeeper,
		mmKeeper.BankKeeper,
		DummyBlockspeedKeeper{},
		DummyDistrubtionKeeper{},
		DummyStakingKeeper{},

		mmKeeper.DenomKeeper.(types.DenomKeeper),
		mmKeeper.DexKeeper.(types.DexKeeper),
		mmKeeper,

		authority.String(),
	)
	cache.AddCache(strategiesKeeper)

	require.NoError(t, cache.Transact(ctx, func(innerCtx context.Context) error {
		return strategiesKeeper.SetParams(innerCtx, StrategiesTestingParams())
	}))

	return strategiesKeeper, dexKeeper, mmKeeper, ctx
}

func StrategiesTestingParams() types.Params {
	return types.Params{
		AutomationFeeCondition: 1,
		AutomationFeeAction:    1,
	}
}

func SetupStrategiesMsgServer(t *testing.T) (keeper.Keeper, types.MsgServer, dextypes.MsgServer, mmtypes.MsgServer, context.Context) {
	k, dexKeeper, mmKeeper, ctx := StrategiesKeeper(t)
	addFunds(ctx, k.BankKeeper, t)
	return k, keeper.NewMsgServerImpl(k), dexkeeper.NewMsgServerImpl(dexKeeper), mmkeeper.NewMsgServerImpl(mmKeeper), ctx
}

func AddArbitrageDeposit(ctx context.Context, k types.MsgServer, msg *types.MsgArbitrageDeposit) error {
	return cache.Transact(ctx, func(innerCtx context.Context) error {
		_, err := k.ArbitrageDeposit(innerCtx, msg)
		return err
	})
}

func Redeem(ctx context.Context, k types.MsgServer, msg *types.MsgArbitrageRedeem) error {
	return cache.Transact(ctx, func(innerCtx context.Context) error {
		_, err := k.ArbitrageRedeem(innerCtx, msg)
		return err
	})
}

func AddAutomation(ctx context.Context, k types.MsgServer, address, title, conditionsString, actionsString, intervalType, intervalLength, validityType, validityValue string) error {
	return AddAutomationMsg(ctx, k, &types.MsgAutomationsAdd{
		Creator:        address,
		Title:          title,
		IntervalType:   intervalType,
		IntervalLength: intervalLength,
		ValidityType:   validityType,
		ValidityValue:  validityValue,
		Conditions:     conditionsString,
		Actions:        actionsString,
	})
}

func AddAutomationMsg(ctx context.Context, k types.MsgServer, msg *types.MsgAutomationsAdd) error {
	return cache.Transact(ctx, func(innerCtx context.Context) error {
		_, err := k.AutomationsAdd(innerCtx, msg)
		return err
	})
}

func AddAutomationFunds(ctx context.Context, k types.MsgServer, address, amount string) error {
	return cache.Transact(ctx, func(innerCtx context.Context) error {
		_, err := k.AutomationsAddFunds(innerCtx, &types.MsgAutomationsAddFunds{
			Creator: address,
			Amount:  amount,
		})
		return err
	})
}
