// Code generated by protoc-gen-go. DO NOT EDIT.
// source: transaction.proto

package corepb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Data struct {
	Type    string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *Data) Reset()                    { *m = Data{} }
func (m *Data) String() string            { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()               {}
func (*Data) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *Data) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Data) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type Transaction struct {
	Hash    []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	From    []byte `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To      []byte `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	Value   []byte `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	Data    *Data  `protobuf:"bytes,5,opt,name=data" json:"data,omitempty"`
	Nonce   uint64 `protobuf:"varint,6,opt,name=nonce" json:"nonce,omitempty"`
	ChainId uint32 `protobuf:"varint,7,opt,name=chain_id,json=chainId" json:"chain_id,omitempty"`
	Alg     uint32 `protobuf:"varint,8,opt,name=alg" json:"alg,omitempty"`
	Sign    []byte `protobuf:"bytes,9,opt,name=sign,proto3" json:"sign,omitempty"`
}

func (m *Transaction) Reset()                    { *m = Transaction{} }
func (m *Transaction) String() string            { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()               {}
func (*Transaction) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *Transaction) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *Transaction) GetFrom() []byte {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Transaction) GetTo() []byte {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *Transaction) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Transaction) GetData() *Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Transaction) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *Transaction) GetChainId() uint32 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

func (m *Transaction) GetAlg() uint32 {
	if m != nil {
		return m.Alg
	}
	return 0
}

func (m *Transaction) GetSign() []byte {
	if m != nil {
		return m.Sign
	}
	return nil
}

func init() {
	proto.RegisterType((*Data)(nil), "corepb.Data")
	proto.RegisterType((*Transaction)(nil), "corepb.Transaction")
}

func init() { proto.RegisterFile("transaction.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 231 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x86, 0xe5, 0x34, 0x4d, 0xda, 0x6b, 0x40, 0x70, 0x62, 0x38, 0x36, 0xab, 0x93, 0xa7, 0x0c,
	0xc0, 0x23, 0xb0, 0xb0, 0x5a, 0xec, 0xe8, 0x9a, 0x84, 0x26, 0x52, 0xf0, 0x45, 0x89, 0x41, 0xea,
	0x7b, 0xf2, 0x40, 0xc8, 0x17, 0x10, 0xdb, 0xf7, 0xff, 0xd6, 0x6f, 0x7d, 0x3a, 0xb8, 0x8d, 0x33,
	0x87, 0x85, 0x9b, 0x38, 0x48, 0xa8, 0xa7, 0x59, 0xa2, 0x60, 0xd1, 0xc8, 0xdc, 0x4d, 0xa7, 0xe3,
	0x13, 0xe4, 0xcf, 0x1c, 0x19, 0x11, 0xf2, 0x78, 0x99, 0x3a, 0x32, 0xd6, 0xb8, 0xbd, 0x57, 0x46,
	0x82, 0x72, 0xe2, 0xcb, 0x28, 0xdc, 0x52, 0x66, 0x8d, 0xab, 0xfc, 0x5f, 0x3c, 0x7e, 0x1b, 0x38,
	0xbc, 0xfe, 0xff, 0x99, 0xd6, 0x3d, 0x2f, 0xbd, 0xae, 0x2b, 0xaf, 0x9c, 0xba, 0xf7, 0x59, 0x3e,
	0x7e, 0xa7, 0xca, 0x78, 0x0d, 0x59, 0x14, 0xda, 0x68, 0x93, 0x45, 0xc1, 0x3b, 0xd8, 0x7e, 0xf1,
	0xf8, 0xd9, 0x51, 0xae, 0xd5, 0x1a, 0xd0, 0x42, 0xde, 0x72, 0x64, 0xda, 0x5a, 0xe3, 0x0e, 0x0f,
	0x55, 0xbd, 0xaa, 0xd6, 0xc9, 0xd3, 0xeb, 0x4b, 0xda, 0x05, 0x09, 0x4d, 0x47, 0x85, 0x35, 0x2e,
	0xf7, 0x6b, 0xc0, 0x7b, 0xd8, 0x35, 0x3d, 0x0f, 0xe1, 0x6d, 0x68, 0xa9, 0xb4, 0xc6, 0x5d, 0xf9,
	0x52, 0xf3, 0x4b, 0x8b, 0x37, 0xb0, 0xe1, 0xf1, 0x4c, 0x3b, 0x6d, 0x13, 0x26, 0xbd, 0x65, 0x38,
	0x07, 0xda, 0xaf, 0x7a, 0x89, 0x4f, 0x85, 0xde, 0xe6, 0xf1, 0x27, 0x00, 0x00, 0xff, 0xff, 0x7e,
	0xb9, 0x8c, 0x41, 0x30, 0x01, 0x00, 0x00,
}
