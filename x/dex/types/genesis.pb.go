// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kopi/dex/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// GenesisState defines the dex module's genesis state.
type GenesisState struct {
	Params             Params              `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	LiquidityList      []DenomLiquidity    `protobuf:"bytes,2,rep,name=liquidity_list,json=liquidityList,proto3" json:"liquidity_list"`
	OrderList          []Order             `protobuf:"bytes,3,rep,name=orderList,proto3" json:"orderList"`
	WalletTradeAmount  []WalletTradeAmount `protobuf:"bytes,4,rep,name=wallet_trade_amount,json=walletTradeAmount,proto3" json:"wallet_trade_amount"`
	LiquidityNextIndex uint64              `protobuf:"varint,5,opt,name=liquidity_next_index,json=liquidityNextIndex,proto3" json:"liquidity_next_index,omitempty"`
	OrderNextIndex     uint64              `protobuf:"varint,6,opt,name=order_next_index,json=orderNextIndex,proto3" json:"order_next_index,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_8564f0e5ae5a7c5b, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetLiquidityList() []DenomLiquidity {
	if m != nil {
		return m.LiquidityList
	}
	return nil
}

func (m *GenesisState) GetOrderList() []Order {
	if m != nil {
		return m.OrderList
	}
	return nil
}

func (m *GenesisState) GetWalletTradeAmount() []WalletTradeAmount {
	if m != nil {
		return m.WalletTradeAmount
	}
	return nil
}

func (m *GenesisState) GetLiquidityNextIndex() uint64 {
	if m != nil {
		return m.LiquidityNextIndex
	}
	return 0
}

func (m *GenesisState) GetOrderNextIndex() uint64 {
	if m != nil {
		return m.OrderNextIndex
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "kopi.dex.GenesisState")
}

func init() { proto.RegisterFile("kopi/dex/genesis.proto", fileDescriptor_8564f0e5ae5a7c5b) }

var fileDescriptor_8564f0e5ae5a7c5b = []byte{
	// 365 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xcd, 0x4a, 0xc3, 0x40,
	0x14, 0x85, 0x93, 0xb6, 0x16, 0x9d, 0x6a, 0xad, 0x63, 0x95, 0x50, 0x21, 0x96, 0x2e, 0x24, 0x1b,
	0x13, 0x69, 0x5f, 0x40, 0x8b, 0x22, 0x42, 0xf1, 0xa7, 0x0a, 0x82, 0x9b, 0x90, 0x9a, 0x4b, 0x1c,
	0x4c, 0x32, 0x31, 0x99, 0xd2, 0xf4, 0x2d, 0x7c, 0xac, 0x2e, 0xbb, 0x74, 0x25, 0xd2, 0xae, 0x7c,
	0x0b, 0xc9, 0x64, 0x3a, 0x11, 0x71, 0x37, 0x39, 0xe7, 0x3b, 0x37, 0xe7, 0x72, 0xd1, 0xfe, 0x2b,
	0x8d, 0x88, 0xe5, 0x42, 0x6a, 0x79, 0x10, 0x42, 0x42, 0x12, 0x33, 0x8a, 0x29, 0xa3, 0x78, 0x3d,
	0xd3, 0x4d, 0x17, 0xd2, 0x56, 0xd3, 0xa3, 0x1e, 0xe5, 0xa2, 0x95, 0xbd, 0x72, 0xbf, 0xb5, 0x27,
	0x73, 0x91, 0x13, 0x3b, 0x81, 0x88, 0xb5, 0x34, 0x29, 0xfb, 0xe4, 0x6d, 0x4c, 0x5c, 0xc2, 0xa6,
	0xc2, 0x69, 0x4a, 0x87, 0xc6, 0x2e, 0xc4, 0x42, 0xed, 0x48, 0x75, 0xe2, 0xf8, 0x3e, 0x30, 0x9b,
	0xc5, 0x8e, 0x0b, 0xb6, 0x13, 0xd0, 0x71, 0xc8, 0x72, 0xa6, 0xf3, 0x5d, 0x42, 0x9b, 0x97, 0x79,
	0xb9, 0x7b, 0xe6, 0x30, 0xc0, 0x26, 0xaa, 0xe6, 0x3f, 0xd5, 0xd4, 0xb6, 0x6a, 0xd4, 0xba, 0x0d,
	0x73, 0x55, 0xd6, 0xbc, 0xe5, 0x7a, 0xbf, 0x32, 0xfb, 0x3c, 0x54, 0x86, 0x82, 0xc2, 0x17, 0xa8,
	0x2e, 0xdb, 0xd8, 0x3e, 0x49, 0x98, 0x56, 0x6a, 0x97, 0x8d, 0x5a, 0x57, 0x2b, 0x72, 0xe7, 0x10,
	0xd2, 0x60, 0xb0, 0x82, 0x44, 0x7e, 0x4b, 0xa6, 0x06, 0x24, 0x61, 0xb8, 0x87, 0x36, 0x78, 0xf5,
	0xec, 0x43, 0x2b, 0xf3, 0x09, 0xdb, 0xc5, 0x84, 0x9b, 0xcc, 0x12, 0xc1, 0x82, 0xc3, 0x77, 0x68,
	0xf7, 0x9f, 0xcd, 0xb4, 0x0a, 0x8f, 0x1f, 0x14, 0xf1, 0x47, 0x0e, 0x3d, 0x64, 0xcc, 0x19, 0x47,
	0xc4, 0xa8, 0x9d, 0xc9, 0x5f, 0x03, 0x9f, 0xa0, 0x66, 0xb1, 0x4e, 0x08, 0x29, 0xb3, 0x49, 0xe8,
	0x42, 0xaa, 0xad, 0xb5, 0x55, 0xa3, 0x32, 0xc4, 0xd2, 0xbb, 0x86, 0x94, 0x5d, 0x65, 0x0e, 0x36,
	0x50, 0x83, 0x37, 0xfa, 0x4d, 0x57, 0x39, 0x5d, 0xe7, 0xba, 0x24, 0xfb, 0xa7, 0xb3, 0x85, 0xae,
	0xce, 0x17, 0xba, 0xfa, 0xb5, 0xd0, 0xd5, 0xf7, 0xa5, 0xae, 0xcc, 0x97, 0xba, 0xf2, 0xb1, 0xd4,
	0x95, 0xa7, 0x23, 0x8f, 0xb0, 0x97, 0xf1, 0xc8, 0x7c, 0xa6, 0x81, 0x95, 0xb5, 0x3e, 0x0e, 0x68,
	0x08, 0x53, 0xfe, 0xb4, 0x52, 0x7e, 0x41, 0x36, 0x8d, 0x20, 0x19, 0x55, 0xf9, 0xd1, 0x7a, 0x3f,
	0x01, 0x00, 0x00, 0xff, 0xff, 0xe6, 0x4a, 0x91, 0x48, 0x59, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.OrderNextIndex != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.OrderNextIndex))
		i--
		dAtA[i] = 0x30
	}
	if m.LiquidityNextIndex != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.LiquidityNextIndex))
		i--
		dAtA[i] = 0x28
	}
	if len(m.WalletTradeAmount) > 0 {
		for iNdEx := len(m.WalletTradeAmount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.WalletTradeAmount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.OrderList) > 0 {
		for iNdEx := len(m.OrderList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.OrderList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.LiquidityList) > 0 {
		for iNdEx := len(m.LiquidityList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.LiquidityList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.LiquidityList) > 0 {
		for _, e := range m.LiquidityList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.OrderList) > 0 {
		for _, e := range m.OrderList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.WalletTradeAmount) > 0 {
		for _, e := range m.WalletTradeAmount {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.LiquidityNextIndex != 0 {
		n += 1 + sovGenesis(uint64(m.LiquidityNextIndex))
	}
	if m.OrderNextIndex != 0 {
		n += 1 + sovGenesis(uint64(m.OrderNextIndex))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidityList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LiquidityList = append(m.LiquidityList, DenomLiquidity{})
			if err := m.LiquidityList[len(m.LiquidityList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OrderList = append(m.OrderList, Order{})
			if err := m.OrderList[len(m.OrderList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WalletTradeAmount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.WalletTradeAmount = append(m.WalletTradeAmount, WalletTradeAmount{})
			if err := m.WalletTradeAmount[len(m.WalletTradeAmount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidityNextIndex", wireType)
			}
			m.LiquidityNextIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LiquidityNextIndex |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderNextIndex", wireType)
			}
			m.OrderNextIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrderNextIndex |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
