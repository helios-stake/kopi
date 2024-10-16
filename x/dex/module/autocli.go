package dex

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/kopi-money/kopi/api/kopi/dex"
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
				{
					RpcMethod: "QuerySimulateSell",
					Use:       "simulate-trade [denom_giving] [denom_receiving] [amount]",
					Short:     "Simulates a sell without executing it",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "denom_giving",
						},
						{
							ProtoField: "denom_receiving",
						},
						{
							ProtoField: "amount",
						},
					},
				},
				{
					RpcMethod: "QuerySimulateBuy",
					Use:       "simulate-buy [denom_giving] [denom_receiving] [amount]",
					Short:     "Simulates a buy without executing it",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "denom_giving",
						},
						{
							ProtoField: "denom_receiving",
						},
						{
							ProtoField: "amount",
						},
					},
				},
				{
					RpcMethod: "Liquidity",
					Use:       "liquidity [denom]",
					Short:     "Return the DEX liquidity of the given denom",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "denom",
						},
					},
				},
				{
					RpcMethod: "LiquidityQueue",
					Use:       "liquidity-queue [denom]",
					Short:     "Return the DEX liquidity queue of the given denom",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "denom",
						},
					},
				},
				{
					RpcMethod: "LiquidityPair",
					Use:       "liquidity-pair [denom]",
					Short:     "Return the DEX liquidity pair of the given denom",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "denom",
						},
					},
				},
				{
					RpcMethod: "OrdersAddress",
					Use:       "orders-address [address]",
					Short:     "Returns all open orders for a given address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "address",
						},
					},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "AddLiquidity",
					Use:       "add-liquidity [denom] [amount]",
					Short:     "Send a AddLiquidity tx",
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
					RpcMethod: "Sell",
					Use:       "sell [denom_giving] [denom_receiving] [amount] [max_price] [minimum_trade_amount]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "denom_giving",
						},
						{
							ProtoField: "denom_receiving",
						},
						{
							ProtoField: "amount",
						},
						{
							ProtoField: "minimum_trade_amount",
						},
						{
							ProtoField: "max_price",
							Optional:   true,
						},
					},
				},
				{
					RpcMethod: "Buy",
					Use:       "buy [denom_giving] [denom_receiving] [amount] [max_price] [minimum_trade_amount]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "denom_giving",
						},
						{
							ProtoField: "denom_receiving",
						},
						{
							ProtoField: "amount",
						},
						{
							ProtoField: "minimum_trade_amount",
						},
						{
							ProtoField: "max_price",
							Optional:   true,
						},
					},
				},
				{
					RpcMethod: "AddOrder",
					Use:       "add-order [denom_giving] [denom_receiving] [amount] [trade_amount] [max_price] [blocks] [interval] [allow_incomplete]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "denom_giving",
						},
						{
							ProtoField: "denom_receiving",
						},
						{
							ProtoField: "amount",
						},
						{
							ProtoField: "trade_amount",
						},
						{
							ProtoField: "max_price",
						},
						{
							ProtoField: "blocks",
						},
						{
							ProtoField: "interval",
						},
						{
							ProtoField: "allow_incomplete",
						},
					},
				},
				{
					RpcMethod: "RemoveOrder",
					Use:       "remove-order [index]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "index",
						},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
