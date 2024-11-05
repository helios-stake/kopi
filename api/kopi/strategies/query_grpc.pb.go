// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: kopi/strategies/query.proto

package strategies

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Query_Params_FullMethodName                         = "/kopi.strategies.Query/Params"
	Query_ArbitrageDenomBalance_FullMethodName          = "/kopi.strategies.Query/ArbitrageDenomBalance"
	Query_ArbitrageBalance_FullMethodName               = "/kopi.strategies.Query/ArbitrageBalance"
	Query_ArbitrageSimulateDepositBase_FullMethodName   = "/kopi.strategies.Query/ArbitrageSimulateDepositBase"
	Query_ArbitrageSimulateDepositCAsset_FullMethodName = "/kopi.strategies.Query/ArbitrageSimulateDepositCAsset"
	Query_ArbitrageSimulateRedemption_FullMethodName    = "/kopi.strategies.Query/ArbitrageSimulateRedemption"
	Query_ArbitrageBalanceAddress_FullMethodName        = "/kopi.strategies.Query/ArbitrageBalanceAddress"
	Query_AutomationsAll_FullMethodName                 = "/kopi.strategies.Query/AutomationsAll"
	Query_AutomationsStats_FullMethodName               = "/kopi.strategies.Query/AutomationsStats"
	Query_AutomationsFunds_FullMethodName               = "/kopi.strategies.Query/AutomationsFunds"
	Query_AutomationsAddress_FullMethodName             = "/kopi.strategies.Query/AutomationsAddress"
	Query_AutomationsIndex_FullMethodName               = "/kopi.strategies.Query/AutomationsIndex"
	Query_AutomationInterval_FullMethodName             = "/kopi.strategies.Query/AutomationInterval"
	Query_AutomationsAddressFunds_FullMethodName        = "/kopi.strategies.Query/AutomationsAddressFunds"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	ArbitrageDenomBalance(ctx context.Context, in *QueryArbitrageDenomBalanceRequest, opts ...grpc.CallOption) (*QueryArbitrageDenomBalanceResponse, error)
	ArbitrageBalance(ctx context.Context, in *QueryArbitrageBalancesRequest, opts ...grpc.CallOption) (*QueryArbitrageBalancesResponse, error)
	ArbitrageSimulateDepositBase(ctx context.Context, in *ArbitrageSimulateDepositRequest, opts ...grpc.CallOption) (*ArbitrageSimulateDepositResponse, error)
	ArbitrageSimulateDepositCAsset(ctx context.Context, in *ArbitrageSimulateDepositRequest, opts ...grpc.CallOption) (*ArbitrageSimulateDepositResponse, error)
	ArbitrageSimulateRedemption(ctx context.Context, in *ArbitrageSimulateRedemptionRequest, opts ...grpc.CallOption) (*ArbitrageSimulateRedemptionResponse, error)
	ArbitrageBalanceAddress(ctx context.Context, in *QueryArbitrageBalancesAddressRequest, opts ...grpc.CallOption) (*QueryArbitrageBalancesAddressResponse, error)
	AutomationsAll(ctx context.Context, in *QueryAutomationsAllRequest, opts ...grpc.CallOption) (*QueryAutomationsResponse, error)
	AutomationsStats(ctx context.Context, in *QueryAutomationsStatsRequest, opts ...grpc.CallOption) (*QueryAutomationsStatsResponse, error)
	AutomationsFunds(ctx context.Context, in *QueryAutomationsFundsRequest, opts ...grpc.CallOption) (*QueryAutomationsFundsResponse, error)
	AutomationsAddress(ctx context.Context, in *QueryAutomationsAddressRequest, opts ...grpc.CallOption) (*QueryAutomationsResponse, error)
	AutomationsIndex(ctx context.Context, in *QueryAutomationsByIndex, opts ...grpc.CallOption) (*Automation, error)
	AutomationInterval(ctx context.Context, in *QueryAutomationsByIndex, opts ...grpc.CallOption) (*QueryAutomationIntervalResponse, error)
	AutomationsAddressFunds(ctx context.Context, in *QueryAutomationsAddressFundsRequest, opts ...grpc.CallOption) (*QueryAutomationsAddressFundsResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ArbitrageDenomBalance(ctx context.Context, in *QueryArbitrageDenomBalanceRequest, opts ...grpc.CallOption) (*QueryArbitrageDenomBalanceResponse, error) {
	out := new(QueryArbitrageDenomBalanceResponse)
	err := c.cc.Invoke(ctx, Query_ArbitrageDenomBalance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ArbitrageBalance(ctx context.Context, in *QueryArbitrageBalancesRequest, opts ...grpc.CallOption) (*QueryArbitrageBalancesResponse, error) {
	out := new(QueryArbitrageBalancesResponse)
	err := c.cc.Invoke(ctx, Query_ArbitrageBalance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ArbitrageSimulateDepositBase(ctx context.Context, in *ArbitrageSimulateDepositRequest, opts ...grpc.CallOption) (*ArbitrageSimulateDepositResponse, error) {
	out := new(ArbitrageSimulateDepositResponse)
	err := c.cc.Invoke(ctx, Query_ArbitrageSimulateDepositBase_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ArbitrageSimulateDepositCAsset(ctx context.Context, in *ArbitrageSimulateDepositRequest, opts ...grpc.CallOption) (*ArbitrageSimulateDepositResponse, error) {
	out := new(ArbitrageSimulateDepositResponse)
	err := c.cc.Invoke(ctx, Query_ArbitrageSimulateDepositCAsset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ArbitrageSimulateRedemption(ctx context.Context, in *ArbitrageSimulateRedemptionRequest, opts ...grpc.CallOption) (*ArbitrageSimulateRedemptionResponse, error) {
	out := new(ArbitrageSimulateRedemptionResponse)
	err := c.cc.Invoke(ctx, Query_ArbitrageSimulateRedemption_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ArbitrageBalanceAddress(ctx context.Context, in *QueryArbitrageBalancesAddressRequest, opts ...grpc.CallOption) (*QueryArbitrageBalancesAddressResponse, error) {
	out := new(QueryArbitrageBalancesAddressResponse)
	err := c.cc.Invoke(ctx, Query_ArbitrageBalanceAddress_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AutomationsAll(ctx context.Context, in *QueryAutomationsAllRequest, opts ...grpc.CallOption) (*QueryAutomationsResponse, error) {
	out := new(QueryAutomationsResponse)
	err := c.cc.Invoke(ctx, Query_AutomationsAll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AutomationsStats(ctx context.Context, in *QueryAutomationsStatsRequest, opts ...grpc.CallOption) (*QueryAutomationsStatsResponse, error) {
	out := new(QueryAutomationsStatsResponse)
	err := c.cc.Invoke(ctx, Query_AutomationsStats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AutomationsFunds(ctx context.Context, in *QueryAutomationsFundsRequest, opts ...grpc.CallOption) (*QueryAutomationsFundsResponse, error) {
	out := new(QueryAutomationsFundsResponse)
	err := c.cc.Invoke(ctx, Query_AutomationsFunds_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AutomationsAddress(ctx context.Context, in *QueryAutomationsAddressRequest, opts ...grpc.CallOption) (*QueryAutomationsResponse, error) {
	out := new(QueryAutomationsResponse)
	err := c.cc.Invoke(ctx, Query_AutomationsAddress_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AutomationsIndex(ctx context.Context, in *QueryAutomationsByIndex, opts ...grpc.CallOption) (*Automation, error) {
	out := new(Automation)
	err := c.cc.Invoke(ctx, Query_AutomationsIndex_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AutomationInterval(ctx context.Context, in *QueryAutomationsByIndex, opts ...grpc.CallOption) (*QueryAutomationIntervalResponse, error) {
	out := new(QueryAutomationIntervalResponse)
	err := c.cc.Invoke(ctx, Query_AutomationInterval_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AutomationsAddressFunds(ctx context.Context, in *QueryAutomationsAddressFundsRequest, opts ...grpc.CallOption) (*QueryAutomationsAddressFundsResponse, error) {
	out := new(QueryAutomationsAddressFundsResponse)
	err := c.cc.Invoke(ctx, Query_AutomationsAddressFunds_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	ArbitrageDenomBalance(context.Context, *QueryArbitrageDenomBalanceRequest) (*QueryArbitrageDenomBalanceResponse, error)
	ArbitrageBalance(context.Context, *QueryArbitrageBalancesRequest) (*QueryArbitrageBalancesResponse, error)
	ArbitrageSimulateDepositBase(context.Context, *ArbitrageSimulateDepositRequest) (*ArbitrageSimulateDepositResponse, error)
	ArbitrageSimulateDepositCAsset(context.Context, *ArbitrageSimulateDepositRequest) (*ArbitrageSimulateDepositResponse, error)
	ArbitrageSimulateRedemption(context.Context, *ArbitrageSimulateRedemptionRequest) (*ArbitrageSimulateRedemptionResponse, error)
	ArbitrageBalanceAddress(context.Context, *QueryArbitrageBalancesAddressRequest) (*QueryArbitrageBalancesAddressResponse, error)
	AutomationsAll(context.Context, *QueryAutomationsAllRequest) (*QueryAutomationsResponse, error)
	AutomationsStats(context.Context, *QueryAutomationsStatsRequest) (*QueryAutomationsStatsResponse, error)
	AutomationsFunds(context.Context, *QueryAutomationsFundsRequest) (*QueryAutomationsFundsResponse, error)
	AutomationsAddress(context.Context, *QueryAutomationsAddressRequest) (*QueryAutomationsResponse, error)
	AutomationsIndex(context.Context, *QueryAutomationsByIndex) (*Automation, error)
	AutomationInterval(context.Context, *QueryAutomationsByIndex) (*QueryAutomationIntervalResponse, error)
	AutomationsAddressFunds(context.Context, *QueryAutomationsAddressFundsRequest) (*QueryAutomationsAddressFundsResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) ArbitrageDenomBalance(context.Context, *QueryArbitrageDenomBalanceRequest) (*QueryArbitrageDenomBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ArbitrageDenomBalance not implemented")
}
func (UnimplementedQueryServer) ArbitrageBalance(context.Context, *QueryArbitrageBalancesRequest) (*QueryArbitrageBalancesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ArbitrageBalance not implemented")
}
func (UnimplementedQueryServer) ArbitrageSimulateDepositBase(context.Context, *ArbitrageSimulateDepositRequest) (*ArbitrageSimulateDepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ArbitrageSimulateDepositBase not implemented")
}
func (UnimplementedQueryServer) ArbitrageSimulateDepositCAsset(context.Context, *ArbitrageSimulateDepositRequest) (*ArbitrageSimulateDepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ArbitrageSimulateDepositCAsset not implemented")
}
func (UnimplementedQueryServer) ArbitrageSimulateRedemption(context.Context, *ArbitrageSimulateRedemptionRequest) (*ArbitrageSimulateRedemptionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ArbitrageSimulateRedemption not implemented")
}
func (UnimplementedQueryServer) ArbitrageBalanceAddress(context.Context, *QueryArbitrageBalancesAddressRequest) (*QueryArbitrageBalancesAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ArbitrageBalanceAddress not implemented")
}
func (UnimplementedQueryServer) AutomationsAll(context.Context, *QueryAutomationsAllRequest) (*QueryAutomationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AutomationsAll not implemented")
}
func (UnimplementedQueryServer) AutomationsStats(context.Context, *QueryAutomationsStatsRequest) (*QueryAutomationsStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AutomationsStats not implemented")
}
func (UnimplementedQueryServer) AutomationsFunds(context.Context, *QueryAutomationsFundsRequest) (*QueryAutomationsFundsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AutomationsFunds not implemented")
}
func (UnimplementedQueryServer) AutomationsAddress(context.Context, *QueryAutomationsAddressRequest) (*QueryAutomationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AutomationsAddress not implemented")
}
func (UnimplementedQueryServer) AutomationsIndex(context.Context, *QueryAutomationsByIndex) (*Automation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AutomationsIndex not implemented")
}
func (UnimplementedQueryServer) AutomationInterval(context.Context, *QueryAutomationsByIndex) (*QueryAutomationIntervalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AutomationInterval not implemented")
}
func (UnimplementedQueryServer) AutomationsAddressFunds(context.Context, *QueryAutomationsAddressFundsRequest) (*QueryAutomationsAddressFundsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AutomationsAddressFunds not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Params_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ArbitrageDenomBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryArbitrageDenomBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ArbitrageDenomBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ArbitrageDenomBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ArbitrageDenomBalance(ctx, req.(*QueryArbitrageDenomBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ArbitrageBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryArbitrageBalancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ArbitrageBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ArbitrageBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ArbitrageBalance(ctx, req.(*QueryArbitrageBalancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ArbitrageSimulateDepositBase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArbitrageSimulateDepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ArbitrageSimulateDepositBase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ArbitrageSimulateDepositBase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ArbitrageSimulateDepositBase(ctx, req.(*ArbitrageSimulateDepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ArbitrageSimulateDepositCAsset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArbitrageSimulateDepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ArbitrageSimulateDepositCAsset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ArbitrageSimulateDepositCAsset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ArbitrageSimulateDepositCAsset(ctx, req.(*ArbitrageSimulateDepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ArbitrageSimulateRedemption_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArbitrageSimulateRedemptionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ArbitrageSimulateRedemption(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ArbitrageSimulateRedemption_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ArbitrageSimulateRedemption(ctx, req.(*ArbitrageSimulateRedemptionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ArbitrageBalanceAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryArbitrageBalancesAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ArbitrageBalanceAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ArbitrageBalanceAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ArbitrageBalanceAddress(ctx, req.(*QueryArbitrageBalancesAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AutomationsAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAutomationsAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AutomationsAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_AutomationsAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AutomationsAll(ctx, req.(*QueryAutomationsAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AutomationsStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAutomationsStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AutomationsStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_AutomationsStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AutomationsStats(ctx, req.(*QueryAutomationsStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AutomationsFunds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAutomationsFundsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AutomationsFunds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_AutomationsFunds_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AutomationsFunds(ctx, req.(*QueryAutomationsFundsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AutomationsAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAutomationsAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AutomationsAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_AutomationsAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AutomationsAddress(ctx, req.(*QueryAutomationsAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AutomationsIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAutomationsByIndex)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AutomationsIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_AutomationsIndex_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AutomationsIndex(ctx, req.(*QueryAutomationsByIndex))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AutomationInterval_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAutomationsByIndex)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AutomationInterval(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_AutomationInterval_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AutomationInterval(ctx, req.(*QueryAutomationsByIndex))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AutomationsAddressFunds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAutomationsAddressFundsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AutomationsAddressFunds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_AutomationsAddressFunds_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AutomationsAddressFunds(ctx, req.(*QueryAutomationsAddressFundsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kopi.strategies.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "ArbitrageDenomBalance",
			Handler:    _Query_ArbitrageDenomBalance_Handler,
		},
		{
			MethodName: "ArbitrageBalance",
			Handler:    _Query_ArbitrageBalance_Handler,
		},
		{
			MethodName: "ArbitrageSimulateDepositBase",
			Handler:    _Query_ArbitrageSimulateDepositBase_Handler,
		},
		{
			MethodName: "ArbitrageSimulateDepositCAsset",
			Handler:    _Query_ArbitrageSimulateDepositCAsset_Handler,
		},
		{
			MethodName: "ArbitrageSimulateRedemption",
			Handler:    _Query_ArbitrageSimulateRedemption_Handler,
		},
		{
			MethodName: "ArbitrageBalanceAddress",
			Handler:    _Query_ArbitrageBalanceAddress_Handler,
		},
		{
			MethodName: "AutomationsAll",
			Handler:    _Query_AutomationsAll_Handler,
		},
		{
			MethodName: "AutomationsStats",
			Handler:    _Query_AutomationsStats_Handler,
		},
		{
			MethodName: "AutomationsFunds",
			Handler:    _Query_AutomationsFunds_Handler,
		},
		{
			MethodName: "AutomationsAddress",
			Handler:    _Query_AutomationsAddress_Handler,
		},
		{
			MethodName: "AutomationsIndex",
			Handler:    _Query_AutomationsIndex_Handler,
		},
		{
			MethodName: "AutomationInterval",
			Handler:    _Query_AutomationInterval_Handler,
		},
		{
			MethodName: "AutomationsAddressFunds",
			Handler:    _Query_AutomationsAddressFunds_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kopi/strategies/query.proto",
}
