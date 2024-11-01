package keeper_test

import (
	keepertest "github.com/kopi-money/kopi/testutil/keeper"
	"github.com/kopi-money/kopi/x/strategies/types"
	"github.com/stretchr/testify/require"
	"testing"
)

var importString = `
[
  {
    "title": "sell test",
    "interval_type": "6",
    "interval_length": "1",
    "validity_type": "7",
    "validity_value": "0",
    "actions": [
      {
        "action_type": 1,
        "amount": "123221000",
        "string1": "ukopi",
        "string2": "ukusd"
      }
    ],
    "conditions": [
      {
        "comparison": "GT",
        "condition_type": 0,
        "string1": "ukopi",
        "string2": "ukusd",
        "value": "33.000000000000000000"
      }
    ]
  }
]
`

func TestImport1(t *testing.T) {
	_, msgServer, _, _, ctx := keepertest.SetupStrategiesMsgServer(t)

	msg := &types.MsgAutomationsImport{
		Creator:     keepertest.Alice,
		Automations: importString,
	}

	require.NoError(t, msg.ValidateBasic())
	require.Error(t, keepertest.ImportAutomationsMsg(ctx, msgServer, msg))
}
