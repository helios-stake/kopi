package arbitrage

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/kopi-money/kopi/api/kopi/strategies"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "ArbitrageDeposit",
					Use:       "arbitrage-deposit [denom] [amount]",
					Short:     "Deposit funds into the arbitrage strategy.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "denom",
						},
						{
							ProtoField: "amount",
						},
					},
				},
				{
					RpcMethod: "ArbitrageRedeem",
					Use:       "arbitrage-redeem [denom] [amount]",
					Short:     "Redeem arbitrage coins from the arbitrage strategy.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "denom",
						},
						{
							ProtoField: "amount",
						},
					},
				},
				{
					RpcMethod: "AutomationsAdd",
					Use:       "automations-add [title] [conditions] [actions] [interval_type] [interval_length] [validity_type] [validity_value]",
					Short:     "Create a new Automation.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "title",
						},
						{
							ProtoField: "conditions",
						},
						{
							ProtoField: "actions",
						},
						{
							ProtoField: "interval_type",
						},
						{
							ProtoField: "interval_length",
						},
						{
							ProtoField: "validity_type",
						},
						{
							ProtoField: "validity_value",
						},
					},
				},
				{
					RpcMethod: "AutomationsUpdate",
					Use:       "automations-update [index] [title] [interval_type] [interval_length] [conditions] [actions]",
					Short:     "Update an existing Automation.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "index",
						},
						{
							ProtoField: "title",
						},
						{
							ProtoField: "interval_type",
						},
						{
							ProtoField: "interval_length",
						},
						{
							ProtoField: "conditions",
						},
						{
							ProtoField: "actions",
						},
					},
				},
				{
					RpcMethod: "AutomationsRemove",
					Use:       "automations-update [index]",
					Short:     "Remove an existing Automation.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "index",
						},
					},
				},
				{
					RpcMethod: "AutomationsActive",
					Use:       "automations-activate [index] [active]",
					Short:     "Activate/Deactivate an Automation.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "index",
						},
						{
							ProtoField: "active",
						},
					},
				},
				{
					RpcMethod: "AutomationsAddFunds",
					Use:       "automations-add-funds [amount]",
					Short:     "Adds funds to be used for executing Automations.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "amount",
						},
					},
				},
				{
					RpcMethod: "AutomationsWithdrawFunds",
					Use:       "automations-remove-funds [amount]",
					Short:     "Removes funds to be used for executing Automations.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "amount",
						},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
