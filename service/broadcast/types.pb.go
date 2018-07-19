// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: types.proto

package broadcast

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type OP int32

const (
	OP_UNKNOWN    OP = 0
	OP_BLOCK_LIST OP = 1
	OP_GET_BLOCKS OP = 2
	OP_GET_HASHES OP = 3
	OP_HASH_LIST  OP = 4
)

var OP_name = map[int32]string{
	0: "UNKNOWN",
	1: "BLOCK_LIST",
	2: "GET_BLOCKS",
	3: "GET_HASHES",
	4: "HASH_LIST",
}
var OP_value = map[string]int32{
	"UNKNOWN":    0,
	"BLOCK_LIST": 1,
	"GET_BLOCKS": 2,
	"GET_HASHES": 3,
	"HASH_LIST":  4,
}

func (x OP) String() string {
	return proto.EnumName(OP_name, int32(x))
}
func (OP) EnumDescriptor() ([]byte, []int) { return fileDescriptorTypes, []int{0} }

type Block struct {
	Node         uint64             `protobuf:"varint,1,opt,name=node,proto3" json:"node,omitempty"`
	Previous     []byte             `protobuf:"bytes,2,opt,name=previous,proto3" json:"previous,omitempty"`
	References   []*SignedData      `protobuf:"bytes,3,rep,name=references" json:"references,omitempty"`
	Round        uint64             `protobuf:"varint,4,opt,name=round,proto3" json:"round,omitempty"`
	Transactions []*TransactionData `protobuf:"bytes,5,rep,name=transactions" json:"transactions,omitempty"`
}

func (m *Block) Reset()                    { *m = Block{} }
func (m *Block) String() string            { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()               {}
func (*Block) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{0} }

func (m *Block) GetNode() uint64 {
	if m != nil {
		return m.Node
	}
	return 0
}

func (m *Block) GetPrevious() []byte {
	if m != nil {
		return m.Previous
	}
	return nil
}

func (m *Block) GetReferences() []*SignedData {
	if m != nil {
		return m.References
	}
	return nil
}

func (m *Block) GetRound() uint64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *Block) GetTransactions() []*TransactionData {
	if m != nil {
		return m.Transactions
	}
	return nil
}

type BlockList struct {
	Blocks []*SignedData `protobuf:"bytes,1,rep,name=blocks" json:"blocks,omitempty"`
}

func (m *BlockList) Reset()                    { *m = BlockList{} }
func (m *BlockList) String() string            { return proto.CompactTextString(m) }
func (*BlockList) ProtoMessage()               {}
func (*BlockList) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{1} }

func (m *BlockList) GetBlocks() []*SignedData {
	if m != nil {
		return m.Blocks
	}
	return nil
}

type BlockReference struct {
	Hash  []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Node  uint64 `protobuf:"varint,2,opt,name=node,proto3" json:"node,omitempty"`
	Round uint64 `protobuf:"varint,3,opt,name=round,proto3" json:"round,omitempty"`
}

func (m *BlockReference) Reset()                    { *m = BlockReference{} }
func (m *BlockReference) String() string            { return proto.CompactTextString(m) }
func (*BlockReference) ProtoMessage()               {}
func (*BlockReference) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{2} }

func (m *BlockReference) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *BlockReference) GetNode() uint64 {
	if m != nil {
		return m.Node
	}
	return 0
}

func (m *BlockReference) GetRound() uint64 {
	if m != nil {
		return m.Round
	}
	return 0
}

type GetBlocks struct {
	Blocks []*BlockReference `protobuf:"bytes,1,rep,name=blocks" json:"blocks,omitempty"`
}

func (m *GetBlocks) Reset()                    { *m = GetBlocks{} }
func (m *GetBlocks) String() string            { return proto.CompactTextString(m) }
func (*GetBlocks) ProtoMessage()               {}
func (*GetBlocks) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{3} }

func (m *GetBlocks) GetBlocks() []*BlockReference {
	if m != nil {
		return m.Blocks
	}
	return nil
}

type GetHashes struct {
	Since uint64 `protobuf:"varint,1,opt,name=since,proto3" json:"since,omitempty"`
}

func (m *GetHashes) Reset()                    { *m = GetHashes{} }
func (m *GetHashes) String() string            { return proto.CompactTextString(m) }
func (*GetHashes) ProtoMessage()               {}
func (*GetHashes) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{4} }

func (m *GetHashes) GetSince() uint64 {
	if m != nil {
		return m.Since
	}
	return 0
}

type HashList struct {
	Hashes [][]byte `protobuf:"bytes,1,rep,name=hashes" json:"hashes,omitempty"`
	Latest uint64   `protobuf:"varint,2,opt,name=latest,proto3" json:"latest,omitempty"`
	Since  uint64   `protobuf:"varint,3,opt,name=since,proto3" json:"since,omitempty"`
}

func (m *HashList) Reset()                    { *m = HashList{} }
func (m *HashList) String() string            { return proto.CompactTextString(m) }
func (*HashList) ProtoMessage()               {}
func (*HashList) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{5} }

func (m *HashList) GetHashes() [][]byte {
	if m != nil {
		return m.Hashes
	}
	return nil
}

func (m *HashList) GetLatest() uint64 {
	if m != nil {
		return m.Latest
	}
	return 0
}

func (m *HashList) GetSince() uint64 {
	if m != nil {
		return m.Since
	}
	return 0
}

type SignedData struct {
	Data      []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *SignedData) Reset()                    { *m = SignedData{} }
func (m *SignedData) String() string            { return proto.CompactTextString(m) }
func (*SignedData) ProtoMessage()               {}
func (*SignedData) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{6} }

func (m *SignedData) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SignedData) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type TransactionData struct {
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Fee  uint64 `protobuf:"varint,2,opt,name=fee,proto3" json:"fee,omitempty"`
}

func (m *TransactionData) Reset()                    { *m = TransactionData{} }
func (m *TransactionData) String() string            { return proto.CompactTextString(m) }
func (*TransactionData) ProtoMessage()               {}
func (*TransactionData) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{7} }

func (m *TransactionData) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *TransactionData) GetFee() uint64 {
	if m != nil {
		return m.Fee
	}
	return 0
}

func init() {
	proto.RegisterType((*Block)(nil), "broadcast.Block")
	proto.RegisterType((*BlockList)(nil), "broadcast.BlockList")
	proto.RegisterType((*BlockReference)(nil), "broadcast.BlockReference")
	proto.RegisterType((*GetBlocks)(nil), "broadcast.GetBlocks")
	proto.RegisterType((*GetHashes)(nil), "broadcast.GetHashes")
	proto.RegisterType((*HashList)(nil), "broadcast.HashList")
	proto.RegisterType((*SignedData)(nil), "broadcast.SignedData")
	proto.RegisterType((*TransactionData)(nil), "broadcast.TransactionData")
	proto.RegisterEnum("broadcast.OP", OP_name, OP_value)
}

func init() { proto.RegisterFile("service/broadcast/types.proto", fileDescriptorTypes) }

var fileDescriptorTypes = []byte{
	// 443 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0x5d, 0x8b, 0xd3, 0x40,
	0x14, 0x35, 0x4d, 0x5b, 0x37, 0x77, 0xeb, 0x1a, 0x06, 0x95, 0x58, 0x14, 0x6a, 0x9e, 0x8a, 0xb0,
	0x2d, 0x2a, 0x22, 0xf8, 0x50, 0x70, 0x75, 0xd9, 0xca, 0x96, 0xb6, 0x24, 0x15, 0x1f, 0x97, 0x49,
	0x32, 0x4d, 0x82, 0xeb, 0x4c, 0x99, 0x99, 0x2c, 0xf8, 0xff, 0xfc, 0x61, 0x32, 0x37, 0x9f, 0x5b,
	0x64, 0xdf, 0xee, 0xb9, 0x73, 0xe6, 0xdc, 0x73, 0xcf, 0x0c, 0xbc, 0x56, 0x4c, 0xde, 0xe5, 0x31,
	0x9b, 0x47, 0x52, 0xd0, 0x24, 0xa6, 0x4a, 0xcf, 0xf5, 0x9f, 0x03, 0x53, 0xb3, 0x83, 0x14, 0x5a,
	0x10, 0xa7, 0x69, 0x8f, 0xcf, 0xd3, 0x5c, 0x67, 0x45, 0x34, 0x8b, 0xc5, 0xef, 0x79, 0x2a, 0x52,
	0x31, 0x47, 0x46, 0x54, 0xec, 0x11, 0x21, 0xc0, 0xaa, 0xbc, 0xe9, 0xff, 0xb5, 0x60, 0x70, 0x71,
	0x2b, 0xe2, 0x5f, 0x84, 0x40, 0x9f, 0x8b, 0x84, 0x79, 0xd6, 0xc4, 0x9a, 0xf6, 0x03, 0xac, 0xc9,
	0x18, 0x4e, 0x0e, 0x92, 0xdd, 0xe5, 0xa2, 0x50, 0x5e, 0x6f, 0x62, 0x4d, 0x47, 0x41, 0x83, 0xc9,
	0x47, 0x00, 0xc9, 0xf6, 0x4c, 0x32, 0x1e, 0x33, 0xe5, 0xd9, 0x13, 0x7b, 0x7a, 0xfa, 0xfe, 0xf9,
	0xac, 0x31, 0x32, 0x0b, 0xf3, 0x94, 0xb3, 0xe4, 0x1b, 0xd5, 0x34, 0xe8, 0x10, 0xc9, 0x33, 0x18,
	0x48, 0x51, 0xf0, 0xc4, 0xeb, 0xe3, 0x9c, 0x12, 0x90, 0x05, 0x8c, 0xb4, 0xa4, 0x5c, 0xd1, 0x58,
	0xe7, 0x82, 0x2b, 0x6f, 0x80, 0x72, 0xe3, 0x8e, 0xdc, 0xae, 0x3d, 0x46, 0xcd, 0x7b, 0x7c, 0xff,
	0x33, 0x38, 0xb8, 0xc5, 0x2a, 0x57, 0x9a, 0x9c, 0xc3, 0x30, 0x32, 0x40, 0x79, 0xd6, 0x43, 0xae,
	0x2a, 0x92, 0xbf, 0x86, 0x33, 0xbc, 0x1b, 0xd4, 0x26, 0x4d, 0x14, 0x19, 0x55, 0x19, 0x46, 0x31,
	0x0a, 0xb0, 0x6e, 0xe2, 0xe9, 0x75, 0xe2, 0x69, 0x76, 0xb1, 0x3b, 0xbb, 0xf8, 0x0b, 0x70, 0xae,
	0x98, 0x46, 0x49, 0x45, 0xde, 0x1d, 0x79, 0x79, 0xd9, 0xf1, 0x72, 0x7f, 0x6a, 0xe3, 0xe7, 0x0d,
	0xde, 0x5f, 0x52, 0x95, 0x95, 0x71, 0xa9, 0x9c, 0xc7, 0xf5, 0xb3, 0x94, 0xc0, 0xdf, 0xc2, 0x89,
	0x39, 0xc7, 0x6d, 0x5f, 0xc0, 0x30, 0x43, 0x2e, 0x4e, 0x18, 0x05, 0x15, 0x32, 0xfd, 0x5b, 0xaa,
	0x99, 0xd2, 0x95, 0xe5, 0x0a, 0xb5, 0x8a, 0x76, 0x57, 0x71, 0x01, 0xd0, 0x46, 0x63, 0x96, 0x4d,
	0xa8, 0xa6, 0x75, 0x00, 0xa6, 0x26, 0xaf, 0xc0, 0x51, 0x79, 0xca, 0xa9, 0x2e, 0x24, 0xab, 0x3e,
	0x43, 0xdb, 0xf0, 0x3f, 0xc1, 0xd3, 0xa3, 0x17, 0xfa, 0xaf, 0x88, 0x0b, 0xf6, 0x9e, 0xd5, 0x21,
	0x9a, 0xf2, 0xed, 0x16, 0x7a, 0x9b, 0x2d, 0x39, 0x85, 0xc7, 0x3f, 0xd6, 0xd7, 0xeb, 0xcd, 0xcf,
	0xb5, 0xfb, 0x88, 0x9c, 0x01, 0x5c, 0xac, 0x36, 0x5f, 0xaf, 0x6f, 0x56, 0xdf, 0xc3, 0x9d, 0x6b,
	0x19, 0x7c, 0x75, 0xb9, 0xbb, 0xc1, 0x5e, 0xe8, 0xf6, 0x6a, 0xbc, 0xfc, 0x12, 0x2e, 0x2f, 0x43,
	0xd7, 0x26, 0x4f, 0xc0, 0x31, 0x75, 0x49, 0xef, 0x47, 0x43, 0xfc, 0xd9, 0x1f, 0xfe, 0x05, 0x00,
	0x00, 0xff, 0xff, 0xe3, 0x06, 0x5e, 0x8f, 0x34, 0x03, 0x00, 0x00,
}
