package keeper

import (
	"github.com/kopi-money/kopi/x/strategies/types"
)

var _ types.QueryServer = Keeper{}
