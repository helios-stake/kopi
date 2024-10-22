package upgrades

import "github.com/kopi-money/kopi/app/upgrades/v0_6_1"

func UpgradeHandlers() Upgrades {
	return Upgrades{
		{
			UpgradeName:          "v0_6_1",
			CreateUpgradeHandler: v0_6_1.CreateUpgradeHandler,
		},
	}
}
