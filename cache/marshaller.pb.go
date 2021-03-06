// Code generated by protoc-gen-go. DO NOT EDIT.
// source: marshaller.proto

package cache

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ProtoMarshaller struct {
	ThreadId             string   `protobuf:"bytes,1,opt,name=thread_id,json=threadId,proto3" json:"thread_id,omitempty"`
	Content              string   `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	SentAt               int64    `protobuf:"varint,3,opt,name=sent_at,json=sentAt,proto3" json:"sent_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProtoMarshaller) Reset()         { *m = ProtoMarshaller{} }
func (m *ProtoMarshaller) String() string { return proto.CompactTextString(m) }
func (*ProtoMarshaller) ProtoMessage()    {}
func (*ProtoMarshaller) Descriptor() ([]byte, []int) {
	return fileDescriptor_0eb4777a4738c610, []int{0}
}

func (m *ProtoMarshaller) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProtoMarshaller.Unmarshal(m, b)
}
func (m *ProtoMarshaller) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProtoMarshaller.Marshal(b, m, deterministic)
}
func (m *ProtoMarshaller) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtoMarshaller.Merge(m, src)
}
func (m *ProtoMarshaller) XXX_Size() int {
	return xxx_messageInfo_ProtoMarshaller.Size(m)
}
func (m *ProtoMarshaller) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtoMarshaller.DiscardUnknown(m)
}

var xxx_messageInfo_ProtoMarshaller proto.InternalMessageInfo

func (m *ProtoMarshaller) GetThreadId() string {
	if m != nil {
		return m.ThreadId
	}
	return ""
}

func (m *ProtoMarshaller) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *ProtoMarshaller) GetSentAt() int64 {
	if m != nil {
		return m.SentAt
	}
	return 0
}

func init() {
	proto.RegisterType((*ProtoMarshaller)(nil), "messages.ProtoMarshaller")
}

func init() { proto.RegisterFile("marshaller.proto", fileDescriptor_0eb4777a4738c610) }

var fileDescriptor_0eb4777a4738c610 = []byte{
	// 131 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc8, 0x4d, 0x2c, 0x2a,
	0xce, 0x48, 0xcc, 0xc9, 0x49, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xc8, 0x4d,
	0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x2d, 0x56, 0x4a, 0xe4, 0xe2, 0x0f, 0x00, 0x09, 0xf9, 0xc2, 0x95,
	0x08, 0x49, 0x73, 0x71, 0x96, 0x64, 0x14, 0xa5, 0x26, 0xa6, 0xc4, 0x67, 0xa6, 0x48, 0x30, 0x2a,
	0x30, 0x6a, 0x70, 0x06, 0x71, 0x40, 0x04, 0x3c, 0x53, 0x84, 0x24, 0xb8, 0xd8, 0x93, 0xf3, 0xf3,
	0x4a, 0x52, 0xf3, 0x4a, 0x24, 0x98, 0xc0, 0x52, 0x30, 0xae, 0x90, 0x38, 0x17, 0x7b, 0x71, 0x6a,
	0x5e, 0x49, 0x7c, 0x62, 0x89, 0x04, 0xb3, 0x02, 0xa3, 0x06, 0x73, 0x10, 0x1b, 0x88, 0xeb, 0x58,
	0x92, 0xc4, 0x06, 0xb6, 0xd3, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xb1, 0x3b, 0xf6, 0x10, 0x87,
	0x00, 0x00, 0x00,
}
