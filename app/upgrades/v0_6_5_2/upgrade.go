package v0_6_5_2

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator, k wasmkeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		// Set CosmWasm params
		wasmParams := wasmtypes.DefaultParams()
		wasmParams.CodeUploadAccess = wasmtypes.AllowNobody
		wasmParams.InstantiateDefaultPermission = wasmtypes.AccessTypeAnyOfAddresses
		if err := k.SetParams(ctx, wasmParams); err != nil {
			return vm, fmt.Errorf("unable to set CosmWasm params")
		}

		return vm, nil
	}
}
