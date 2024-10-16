// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kopi/mm/blockspeed.proto

package types

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	cosmossdk_io_math "cosmossdk.io/math"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

type BlockSpeed struct {
	LatestTimestamp   string                      `protobuf:"bytes,1,opt,name=latest_timestamp,json=latestTimestamp,proto3" json:"latest_timestamp,omitempty"`
	AverageBlockSpeed cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=average_block_speed,json=averageBlockSpeed,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"average_block_speed"`
}

func (m *BlockSpeed) Reset()         { *m = BlockSpeed{} }
func (m *BlockSpeed) String() string { return proto.CompactTextString(m) }
func (*BlockSpeed) ProtoMessage()    {}
func (*BlockSpeed) Descriptor() ([]byte, []int) {
	return fileDescriptor_7f70153a1c310696, []int{0}
}
func (m *BlockSpeed) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BlockSpeed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BlockSpeed.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BlockSpeed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockSpeed.Merge(m, src)
}
func (m *BlockSpeed) XXX_Size() int {
	return m.Size()
}
func (m *BlockSpeed) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockSpeed.DiscardUnknown(m)
}

var xxx_messageInfo_BlockSpeed proto.InternalMessageInfo

func (m *BlockSpeed) GetLatestTimestamp() string {
	if m != nil {
		return m.LatestTimestamp
	}
	return ""
}

func init() {
	proto.RegisterType((*BlockSpeed)(nil), "kopi.mm.BlockSpeed")
}

func init() { proto.RegisterFile("kopi/mm/blockspeed.proto", fileDescriptor_7f70153a1c310696) }

var fileDescriptor_7f70153a1c310696 = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xc8, 0xce, 0x2f, 0xc8,
	0xd4, 0xcf, 0xcd, 0xd5, 0x4f, 0xca, 0xc9, 0x4f, 0xce, 0x2e, 0x2e, 0x48, 0x4d, 0x4d, 0xd1, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x07, 0xc9, 0xe8, 0xe5, 0xe6, 0x4a, 0x89, 0xa4, 0xe7, 0xa7,
	0xe7, 0x83, 0xc5, 0xf4, 0x41, 0x2c, 0x88, 0xb4, 0x52, 0x0f, 0x23, 0x17, 0x97, 0x13, 0x48, 0x4f,
	0x30, 0x48, 0x8f, 0x90, 0x26, 0x97, 0x40, 0x4e, 0x62, 0x49, 0x6a, 0x71, 0x49, 0x7c, 0x49, 0x66,
	0x6e, 0x6a, 0x71, 0x49, 0x62, 0x6e, 0x81, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x3f, 0x44,
	0x3c, 0x04, 0x26, 0x2c, 0x14, 0xcc, 0x25, 0x9c, 0x58, 0x96, 0x5a, 0x94, 0x98, 0x9e, 0x1a, 0x0f,
	0xb6, 0x34, 0x1e, 0x6c, 0xab, 0x04, 0x93, 0x02, 0xa3, 0x06, 0x8f, 0x93, 0xf2, 0x89, 0x7b, 0xf2,
	0x0c, 0xb7, 0xee, 0xc9, 0x4b, 0x27, 0xe7, 0x17, 0xe7, 0xe6, 0x17, 0x17, 0xa7, 0x64, 0xeb, 0x65,
	0xe6, 0xeb, 0xe7, 0x26, 0x96, 0x64, 0xe8, 0xf9, 0xa4, 0xa6, 0x27, 0x26, 0x57, 0xba, 0xa4, 0x26,
	0x07, 0x09, 0x42, 0xf5, 0x23, 0xec, 0x77, 0xb2, 0x3f, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39,
	0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63,
	0x39, 0x86, 0x28, 0xd5, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0x7d, 0x90,
	0x97, 0x74, 0x73, 0xf3, 0xf3, 0x52, 0x2b, 0xc1, 0x4c, 0xfd, 0x0a, 0x90, 0xcf, 0x4b, 0x2a, 0x0b,
	0x52, 0x8b, 0x93, 0xd8, 0xc0, 0xde, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xce, 0x3c, 0x08,
	0xfd, 0x11, 0x01, 0x00, 0x00,
}

func (m *BlockSpeed) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BlockSpeed) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BlockSpeed) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.AverageBlockSpeed.Size()
		i -= size
		if _, err := m.AverageBlockSpeed.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintBlockspeed(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.LatestTimestamp) > 0 {
		i -= len(m.LatestTimestamp)
		copy(dAtA[i:], m.LatestTimestamp)
		i = encodeVarintBlockspeed(dAtA, i, uint64(len(m.LatestTimestamp)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintBlockspeed(dAtA []byte, offset int, v uint64) int {
	offset -= sovBlockspeed(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BlockSpeed) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.LatestTimestamp)
	if l > 0 {
		n += 1 + l + sovBlockspeed(uint64(l))
	}
	l = m.AverageBlockSpeed.Size()
	n += 1 + l + sovBlockspeed(uint64(l))
	return n
}

func sovBlockspeed(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBlockspeed(x uint64) (n int) {
	return sovBlockspeed(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BlockSpeed) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBlockspeed
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
			return fmt.Errorf("proto: BlockSpeed: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BlockSpeed: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestTimestamp", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockspeed
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
				return ErrInvalidLengthBlockspeed
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBlockspeed
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LatestTimestamp = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AverageBlockSpeed", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockspeed
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
				return ErrInvalidLengthBlockspeed
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBlockspeed
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AverageBlockSpeed.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBlockspeed(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBlockspeed
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
func skipBlockspeed(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBlockspeed
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
					return 0, ErrIntOverflowBlockspeed
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
					return 0, ErrIntOverflowBlockspeed
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
				return 0, ErrInvalidLengthBlockspeed
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBlockspeed
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBlockspeed
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBlockspeed        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBlockspeed          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBlockspeed = fmt.Errorf("proto: unexpected end of group")
)
