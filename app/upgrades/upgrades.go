package upgrades

import "github.com/kopi-money/kopi/app/upgrades/v0_6"

func UpgradeHandlers() Upgrades {
	return Upgrades{
		{
			UpgradeName:          "v0_6",
			CreateUpgradeHandler: v0_6.CreateUpgradeHandler,
		},
	}
}
