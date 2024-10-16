// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kopi/blockspeed/blockspeed.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
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

type blockspeed struct {
	PreviousTimestamp int64                       `protobuf:"varint,1,opt,name=previous_timestamp,json=previousTimestamp,proto3" json:"previous_timestamp,omitempty"`
	AverageTime       cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=average_time,json=averageTime,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"average_time"`
}

func (m *blockspeed) Reset()         { *m = blockspeed{} }
func (m *blockspeed) String() string { return proto.CompactTextString(m) }
func (*blockspeed) ProtoMessage()    {}
func (*blockspeed) Descriptor() ([]byte, []int) {
	return fileDescriptor_63a91c914deed21f, []int{0}
}
func (m *blockspeed) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *blockspeed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_blockspeed.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *blockspeed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_blockspeed.Merge(m, src)
}
func (m *blockspeed) XXX_Size() int {
	return m.Size()
}
func (m *blockspeed) XXX_DiscardUnknown() {
	xxx_messageInfo_blockspeed.DiscardUnknown(m)
}

var xxx_messageInfo_blockspeed proto.InternalMessageInfo

func (m *blockspeed) GetPreviousTimestamp() int64 {
	if m != nil {
		return m.PreviousTimestamp
	}
	return 0
}

func init() {
	proto.RegisterType((*blockspeed)(nil), "kopi.blockspeed.blockspeed")
}

func init() { proto.RegisterFile("kopi/blockspeed/blockspeed.proto", fileDescriptor_63a91c914deed21f) }

var fileDescriptor_63a91c914deed21f = []byte{
	// 236 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xcb, 0xce, 0x2f, 0xc8,
	0xd4, 0x4f, 0xca, 0xc9, 0x4f, 0xce, 0x2e, 0xc9, 0xcc, 0x4d, 0x45, 0xb0, 0xf4, 0x0a, 0x8a, 0xf2,
	0x4b, 0xf2, 0x85, 0xf8, 0x40, 0xf2, 0x7a, 0x70, 0x51, 0x29, 0x91, 0xf4, 0xfc, 0xf4, 0x7c, 0xb0,
	0x94, 0x3e, 0x88, 0x05, 0x51, 0xa5, 0xd4, 0xc4, 0xc8, 0xc5, 0xe9, 0x04, 0x53, 0x23, 0xa4, 0xcb,
	0x25, 0x54, 0x50, 0x94, 0x5a, 0x96, 0x99, 0x5f, 0x5a, 0x1c, 0x0f, 0x12, 0x28, 0x2e, 0x49, 0xcc,
	0x2d, 0x90, 0x60, 0x54, 0x60, 0xd4, 0x60, 0x0e, 0x12, 0x84, 0xc9, 0x84, 0xc0, 0x24, 0x84, 0xdc,
	0xb8, 0x78, 0x12, 0xcb, 0x52, 0x8b, 0x12, 0xd3, 0x53, 0xc1, 0xaa, 0x25, 0x98, 0x14, 0x18, 0x35,
	0x78, 0x9c, 0x94, 0x4f, 0xdc, 0x93, 0x67, 0xb8, 0x75, 0x4f, 0x5e, 0x3a, 0x39, 0xbf, 0x38, 0x37,
	0xbf, 0xb8, 0x38, 0x25, 0x5b, 0x2f, 0x33, 0x5f, 0x3f, 0x37, 0xb1, 0x24, 0x43, 0xcf, 0x27, 0x35,
	0x3d, 0x31, 0xb9, 0xd2, 0x25, 0x35, 0x39, 0x88, 0x1b, 0xaa, 0x11, 0x64, 0x98, 0x93, 0xdb, 0x89,
	0x47, 0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3,
	0x85, 0xc7, 0x72, 0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x44, 0xe9, 0xa4, 0x67, 0x96, 0x64, 0x94, 0x26,
	0xe9, 0x25, 0xe7, 0xe7, 0xea, 0x83, 0xfc, 0xa3, 0x9b, 0x9b, 0x9f, 0x97, 0x5a, 0x09, 0x66, 0xea,
	0x57, 0x20, 0x79, 0xbe, 0xa4, 0xb2, 0x20, 0xb5, 0x38, 0x89, 0x0d, 0xec, 0x27, 0x63, 0x40, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x45, 0x52, 0x92, 0xad, 0x1b, 0x01, 0x00, 0x00,
}

func (m *blockspeed) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *blockspeed) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *blockspeed) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.AverageTime.Size()
		i -= size
		if _, err := m.AverageTime.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintblockspeed(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.PreviousTimestamp != 0 {
		i = encodeVarintblockspeed(dAtA, i, uint64(m.PreviousTimestamp))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintblockspeed(dAtA []byte, offset int, v uint64) int {
	offset -= sovblockspeed(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *blockspeed) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PreviousTimestamp != 0 {
		n += 1 + sovblockspeed(uint64(m.PreviousTimestamp))
	}
	l = m.AverageTime.Size()
	n += 1 + l + sovblockspeed(uint64(l))
	return n
}

func sovblockspeed(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozblockspeed(x uint64) (n int) {
	return sovblockspeed(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *blockspeed) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowblockspeed
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
			return fmt.Errorf("proto: blockspeed: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: blockspeed: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PreviousTimestamp", wireType)
			}
			m.PreviousTimestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowblockspeed
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PreviousTimestamp |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AverageTime", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowblockspeed
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthblockspeed
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthblockspeed
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AverageTime.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipblockspeed(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthblockspeed
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
func skipblockspeed(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowblockspeed
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
					return 0, ErrIntOverflowblockspeed
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
					return 0, ErrIntOverflowblockspeed
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
				return 0, ErrInvalidLengthblockspeed
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupblockspeed
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthblockspeed
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthblockspeed        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowblockspeed          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupblockspeed = fmt.Errorf("proto: unexpected end of group")
)
