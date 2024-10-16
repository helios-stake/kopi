// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kopi/blockspeed/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type QueryBlockspeedRequest struct {
}

func (m *QueryBlockspeedRequest) Reset()         { *m = QueryBlockspeedRequest{} }
func (m *QueryBlockspeedRequest) String() string { return proto.CompactTextString(m) }
func (*QueryBlockspeedRequest) ProtoMessage()    {}
func (*QueryBlockspeedRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ff3ae169c344afa, []int{0}
}
func (m *QueryBlockspeedRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryBlockspeedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryBlockspeedRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryBlockspeedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBlockspeedRequest.Merge(m, src)
}
func (m *QueryBlockspeedRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryBlockspeedRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBlockspeedRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBlockspeedRequest proto.InternalMessageInfo

type QueryBlockspeedResponse struct {
	BlocksPerSecond string `protobuf:"bytes,1,opt,name=blocks_per_second,json=blocksPerSecond,proto3" json:"blocks_per_second,omitempty"`
	SecondsPerBlock string `protobuf:"bytes,2,opt,name=seconds_per_block,json=secondsPerBlock,proto3" json:"seconds_per_block,omitempty"`
}

func (m *QueryBlockspeedResponse) Reset()         { *m = QueryBlockspeedResponse{} }
func (m *QueryBlockspeedResponse) String() string { return proto.CompactTextString(m) }
func (*QueryBlockspeedResponse) ProtoMessage()    {}
func (*QueryBlockspeedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ff3ae169c344afa, []int{1}
}
func (m *QueryBlockspeedResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryBlockspeedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryBlockspeedResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryBlockspeedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBlockspeedResponse.Merge(m, src)
}
func (m *QueryBlockspeedResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryBlockspeedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBlockspeedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBlockspeedResponse proto.InternalMessageInfo

func (m *QueryBlockspeedResponse) GetBlocksPerSecond() string {
	if m != nil {
		return m.BlocksPerSecond
	}
	return ""
}

func (m *QueryBlockspeedResponse) GetSecondsPerBlock() string {
	if m != nil {
		return m.SecondsPerBlock
	}
	return ""
}

type QueryParamsRequest struct {
}

func (m *QueryParamsRequest) Reset()         { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryParamsRequest) ProtoMessage()    {}
func (*QueryParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ff3ae169c344afa, []int{2}
}
func (m *QueryParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsRequest.Merge(m, src)
}
func (m *QueryParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsRequest proto.InternalMessageInfo

type QueryParamsResponse struct {
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryParamsResponse) Reset()         { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryParamsResponse) ProtoMessage()    {}
func (*QueryParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ff3ae169c344afa, []int{3}
}
func (m *QueryParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsResponse.Merge(m, src)
}
func (m *QueryParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsResponse proto.InternalMessageInfo

func (m *QueryParamsResponse) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func init() {
	proto.RegisterType((*QueryBlockspeedRequest)(nil), "kopi.blockspeed.QueryBlockspeedRequest")
	proto.RegisterType((*QueryBlockspeedResponse)(nil), "kopi.blockspeed.QueryBlockspeedResponse")
	proto.RegisterType((*QueryParamsRequest)(nil), "kopi.blockspeed.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "kopi.blockspeed.QueryParamsResponse")
}

func init() { proto.RegisterFile("kopi/blockspeed/query.proto", fileDescriptor_1ff3ae169c344afa) }

var fileDescriptor_1ff3ae169c344afa = []byte{
	// 407 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xc1, 0xae, 0xd2, 0x40,
	0x14, 0x86, 0x5b, 0x12, 0x49, 0x18, 0x17, 0x84, 0x91, 0x08, 0x16, 0x2d, 0xa6, 0x98, 0x48, 0x48,
	0xe8, 0x04, 0xdc, 0xb9, 0x64, 0xe3, 0x16, 0x70, 0xe7, 0x86, 0x4c, 0xcb, 0xa4, 0x36, 0xd0, 0x39,
	0x43, 0x67, 0x30, 0xb2, 0x35, 0x3c, 0x80, 0x89, 0x2f, 0xe1, 0xd2, 0xc7, 0x60, 0x49, 0xe2, 0xc6,
	0x95, 0x31, 0x60, 0xe2, 0x6b, 0x98, 0xce, 0x94, 0xcb, 0xbd, 0x94, 0xdc, 0xbb, 0x69, 0x4e, 0xce,
	0xff, 0x9f, 0xff, 0x7c, 0x3d, 0x2d, 0x6a, 0x2d, 0x40, 0xc4, 0x24, 0x58, 0x42, 0xb8, 0x90, 0x82,
	0xb1, 0x39, 0x59, 0xad, 0x59, 0xba, 0xf1, 0x45, 0x0a, 0x0a, 0x70, 0x35, 0x13, 0xfd, 0xb3, 0xe8,
	0xd4, 0x68, 0x12, 0x73, 0x20, 0xfa, 0x69, 0x3c, 0x4e, 0x3d, 0x82, 0x08, 0x74, 0x49, 0xb2, 0x2a,
	0xef, 0x3e, 0x8f, 0x00, 0xa2, 0x25, 0x23, 0x54, 0xc4, 0x84, 0x72, 0x0e, 0x8a, 0xaa, 0x18, 0xb8,
	0xcc, 0xd5, 0x5e, 0x08, 0x32, 0x01, 0x49, 0x02, 0x2a, 0x99, 0x59, 0x48, 0x3e, 0x0d, 0x02, 0xa6,
	0xe8, 0x80, 0x08, 0x1a, 0xc5, 0x5c, 0x9b, 0x4f, 0x49, 0x97, 0x80, 0x82, 0xa6, 0x34, 0xc9, 0x93,
	0xbc, 0x26, 0x7a, 0x3a, 0xc9, 0xe6, 0x47, 0x37, 0xfa, 0x94, 0xad, 0xd6, 0x4c, 0x2a, 0x6f, 0x85,
	0x1a, 0x05, 0x45, 0x0a, 0xe0, 0x92, 0xe1, 0x1e, 0xaa, 0x99, 0xbc, 0x99, 0x60, 0xe9, 0x4c, 0xb2,
	0x10, 0xf8, 0xbc, 0x69, 0xbf, 0xb4, 0xbb, 0x95, 0x69, 0xd5, 0x08, 0x63, 0x96, 0xbe, 0xd7, 0xed,
	0xcc, 0x6b, 0x0c, 0xc6, 0xac, 0xe5, 0x66, 0xc9, 0x78, 0x73, 0x61, 0xcc, 0x52, 0xbd, 0xc4, 0xab,
	0x23, 0xac, 0x57, 0x8e, 0x35, 0xe1, 0x09, 0x64, 0x82, 0x9e, 0xdc, 0xe9, 0xe6, 0x10, 0x6f, 0x51,
	0xd9, 0xbc, 0x89, 0xde, 0xfc, 0x78, 0xd8, 0xf0, 0x2f, 0x8e, 0xed, 0x9b, 0x81, 0x51, 0x65, 0xf7,
	0xbb, 0x6d, 0x7d, 0xff, 0xf7, 0xa3, 0x67, 0x4f, 0xf3, 0x89, 0xe1, 0xb6, 0x84, 0x1e, 0xe9, 0x4c,
	0xac, 0x50, 0xd9, 0xd8, 0x70, 0xa7, 0x30, 0x5f, 0x64, 0x71, 0x5e, 0xdd, 0x6f, 0x32, 0x68, 0x5e,
	0xfb, 0xcb, 0xcf, 0xbf, 0xdf, 0x4a, 0xcf, 0x70, 0x83, 0x5c, 0xbf, 0x3d, 0xde, 0xda, 0x08, 0x9d,
	0xef, 0x8a, 0x5f, 0x5f, 0x4f, 0x2d, 0x7c, 0x13, 0xa7, 0xfb, 0xb0, 0x31, 0x47, 0xe8, 0x68, 0x84,
	0x17, 0xb8, 0x55, 0x40, 0x38, 0x97, 0xa3, 0x77, 0xbb, 0x83, 0x6b, 0xef, 0x0f, 0xae, 0xfd, 0xe7,
	0xe0, 0xda, 0x5f, 0x8f, 0xae, 0xb5, 0x3f, 0xba, 0xd6, 0xaf, 0xa3, 0x6b, 0x7d, 0xe8, 0x47, 0xb1,
	0xfa, 0xb8, 0x0e, 0xfc, 0x10, 0x12, 0x1d, 0xd0, 0x4f, 0x80, 0xb3, 0x8d, 0xc9, 0xfa, 0x7c, 0x3b,
	0x4d, 0x6d, 0x04, 0x93, 0x41, 0x59, 0xff, 0x4c, 0x6f, 0xfe, 0x07, 0x00, 0x00, 0xff, 0xff, 0x9d,
	0x01, 0xf4, 0x44, 0x0d, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	Blockspeed(ctx context.Context, in *QueryBlockspeedRequest, opts ...grpc.CallOption) (*QueryBlockspeedResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/kopi.blockspeed.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Blockspeed(ctx context.Context, in *QueryBlockspeedRequest, opts ...grpc.CallOption) (*QueryBlockspeedResponse, error) {
	out := new(QueryBlockspeedResponse)
	err := c.cc.Invoke(ctx, "/kopi.blockspeed.Query/Blockspeed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	Blockspeed(context.Context, *QueryBlockspeedRequest) (*QueryBlockspeedResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) Blockspeed(ctx context.Context, req *QueryBlockspeedRequest) (*QueryBlockspeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Blockspeed not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
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
		FullMethod: "/kopi.blockspeed.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Blockspeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBlockspeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Blockspeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kopi.blockspeed.Query/Blockspeed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Blockspeed(ctx, req.(*QueryBlockspeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kopi.blockspeed.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "Blockspeed",
			Handler:    _Query_Blockspeed_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kopi/blockspeed/query.proto",
}

func (m *QueryBlockspeedRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryBlockspeedRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryBlockspeedRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryBlockspeedResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryBlockspeedResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryBlockspeedResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.SecondsPerBlock) > 0 {
		i -= len(m.SecondsPerBlock)
		copy(dAtA[i:], m.SecondsPerBlock)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.SecondsPerBlock)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.BlocksPerSecond) > 0 {
		i -= len(m.BlocksPerSecond)
		copy(dAtA[i:], m.BlocksPerSecond)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.BlocksPerSecond)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryBlockspeedRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryBlockspeedResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.BlocksPerSecond)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	l = len(m.SecondsPerBlock)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryBlockspeedRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryBlockspeedRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryBlockspeedRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryBlockspeedResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryBlockspeedResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryBlockspeedResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlocksPerSecond", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BlocksPerSecond = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SecondsPerBlock", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SecondsPerBlock = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
