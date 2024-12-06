package app

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/kopi-money/kopi/app/upgrades/v0_6_1"
	"github.com/kopi-money/kopi/app/upgrades/v0_6_2"
	"github.com/kopi-money/kopi/app/upgrades/v0_6_3"
	"github.com/kopi-money/kopi/app/upgrades/v0_6_4"
	"github.com/kopi-money/kopi/app/upgrades/v0_6_5_1"
	"github.com/kopi-money/kopi/app/upgrades/v0_6_5_2"
)

func (app *App) setupUpgradeHandlers(appOpts servertypes.AppOptions) {
	// set upgrade v061
	app.UpgradeKeeper.SetUpgradeHandler(
		"v0_6_1",
		v0_6_1.CreateUpgradeHandler(app.ModuleManager, app.Configurator()),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		"v0_6_2",
		v0_6_2.CreateUpgradeHandler(app.ModuleManager, app.Configurator()),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		"v0_6_3",
		v0_6_3.CreateUpgradeHandler(app.ModuleManager, app.Configurator()),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		"v0_6_4",
		v0_6_4.CreateUpgradeHandler(app.ModuleManager, app.Configurator()),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		"v0_6_5_1",
		v0_6_5_1.CreateUpgradeHandler(app.ModuleManager, app.Configurator()),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		v0_6_5_2.UpgradeName,
		v0_6_5_2.CreateUpgradeHandler(app.ModuleManager, app.Configurator(), app.WasmKeeper),
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Errorf("Failed to read upgrade info from disk: %w", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}
	var storeUpgrades *storetypes.StoreUpgrades
	switch upgradeInfo.Name {
	case v0_6_5_2.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{wasmtypes.ModuleName},
		}
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
