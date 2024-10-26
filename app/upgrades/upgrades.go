package upgrades

import (
	"github.com/kopi-money/kopi/app/upgrades/v0_6_1"
	"github.com/kopi-money/kopi/app/upgrades/v0_6_2"
	"github.com/kopi-money/kopi/app/upgrades/v0_6_3"
)

func UpgradeHandlers() Upgrades {
	return Upgrades{
		{
			UpgradeName:          "v0_6_1",
			CreateUpgradeHandler: v0_6_1.CreateUpgradeHandler,
		},
		{
			UpgradeName:          "v0_6_2",
			CreateUpgradeHandler: v0_6_2.CreateUpgradeHandler,
		},
		{
			UpgradeName:          "v0_6_3",
			CreateUpgradeHandler: v0_6_3.CreateUpgradeHandler,
		},
	}
}
