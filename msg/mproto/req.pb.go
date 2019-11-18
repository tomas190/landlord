// Code generated by protoc-gen-go. DO NOT EDIT.
// source: req.proto

package mproto

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

// CLIENT REQ : MSG_ID 0
// 1.PING
type PING struct {
	Time                 int64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PING) Reset()         { *m = PING{} }
func (m *PING) String() string { return proto.CompactTextString(m) }
func (*PING) ProtoMessage()    {}
func (*PING) Descriptor() ([]byte, []int) {
	return fileDescriptor_a52bd445c288993d, []int{0}
}

func (m *PING) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PING.Unmarshal(m, b)
}
func (m *PING) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PING.Marshal(b, m, deterministic)
}
func (m *PING) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PING.Merge(m, src)
}
func (m *PING) XXX_Size() int {
	return xxx_messageInfo_PING.Size(m)
}
func (m *PING) XXX_DiscardUnknown() {
	xxx_messageInfo_PING.DiscardUnknown(m)
}

var xxx_messageInfo_PING proto.InternalMessageInfo

func (m *PING) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

// CLIENT REQ : MSG_ID 100
// 1.登录请求
type ReqLogin struct {
	UserId               string   `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	UserPassword         string   `protobuf:"bytes,2,opt,name=userPassword,proto3" json:"userPassword,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqLogin) Reset()         { *m = ReqLogin{} }
func (m *ReqLogin) String() string { return proto.CompactTextString(m) }
func (*ReqLogin) ProtoMessage()    {}
func (*ReqLogin) Descriptor() ([]byte, []int) {
	return fileDescriptor_a52bd445c288993d, []int{1}
}

func (m *ReqLogin) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqLogin.Unmarshal(m, b)
}
func (m *ReqLogin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqLogin.Marshal(b, m, deterministic)
}
func (m *ReqLogin) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqLogin.Merge(m, src)
}
func (m *ReqLogin) XXX_Size() int {
	return xxx_messageInfo_ReqLogin.Size(m)
}
func (m *ReqLogin) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqLogin.DiscardUnknown(m)
}

var xxx_messageInfo_ReqLogin proto.InternalMessageInfo

func (m *ReqLogin) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *ReqLogin) GetUserPassword() string {
	if m != nil {
		return m.UserPassword
	}
	return ""
}

// CLIENT REQ : MSG_ID 101
// 2.登录请求
type ReqEnterRoom struct {
	//
	//ExperienceField int32 = 1 // 体验场
	//LowField        int32 = 2 // 初级场
	//MidField        int32 = 3 // 中级场
	//HighField       int32 = 4 // 高级场
	RoomType             int32    `protobuf:"varint,1,opt,name=roomType,proto3" json:"roomType,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqEnterRoom) Reset()         { *m = ReqEnterRoom{} }
func (m *ReqEnterRoom) String() string { return proto.CompactTextString(m) }
func (*ReqEnterRoom) ProtoMessage()    {}
func (*ReqEnterRoom) Descriptor() ([]byte, []int) {
	return fileDescriptor_a52bd445c288993d, []int{2}
}

func (m *ReqEnterRoom) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqEnterRoom.Unmarshal(m, b)
}
func (m *ReqEnterRoom) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqEnterRoom.Marshal(b, m, deterministic)
}
func (m *ReqEnterRoom) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqEnterRoom.Merge(m, src)
}
func (m *ReqEnterRoom) XXX_Size() int {
	return xxx_messageInfo_ReqEnterRoom.Size(m)
}
func (m *ReqEnterRoom) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqEnterRoom.DiscardUnknown(m)
}

var xxx_messageInfo_ReqEnterRoom proto.InternalMessageInfo

func (m *ReqEnterRoom) GetRoomType() int32 {
	if m != nil {
		return m.RoomType
	}
	return 0
}

// MSG_DO_ACTION = 102; // c-s DoAction
// 1.打牌
// 2.碰牌
// .吃牌{31 左吃,32中吃,33右吃}
// .杠牌{41 杠(对方玩家打牌 自己能杠),42 加杠(自己碰牌后 摸到一张可以杠的牌),43 暗杠(自己手牌有三个,然后摸到一个或者手牌有四个)}
// 5,听牌
// 6,胡牌
// REQ AND RESP
type ReqDoAction struct {
	ActionCard           int32    `protobuf:"varint,1,opt,name=actionCard,proto3" json:"actionCard,omitempty"`
	ActionType           int32    `protobuf:"varint,2,opt,name=actionType,proto3" json:"actionType,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqDoAction) Reset()         { *m = ReqDoAction{} }
func (m *ReqDoAction) String() string { return proto.CompactTextString(m) }
func (*ReqDoAction) ProtoMessage()    {}
func (*ReqDoAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_a52bd445c288993d, []int{3}
}

func (m *ReqDoAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqDoAction.Unmarshal(m, b)
}
func (m *ReqDoAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqDoAction.Marshal(b, m, deterministic)
}
func (m *ReqDoAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqDoAction.Merge(m, src)
}
func (m *ReqDoAction) XXX_Size() int {
	return xxx_messageInfo_ReqDoAction.Size(m)
}
func (m *ReqDoAction) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqDoAction.DiscardUnknown(m)
}

var xxx_messageInfo_ReqDoAction proto.InternalMessageInfo

func (m *ReqDoAction) GetActionCard() int32 {
	if m != nil {
		return m.ActionCard
	}
	return 0
}

func (m *ReqDoAction) GetActionType() int32 {
	if m != nil {
		return m.ActionType
	}
	return 0
}

func init() {
	proto.RegisterType((*PING)(nil), "mproto.PING")
	proto.RegisterType((*ReqLogin)(nil), "mproto.ReqLogin")
	proto.RegisterType((*ReqEnterRoom)(nil), "mproto.ReqEnterRoom")
	proto.RegisterType((*ReqDoAction)(nil), "mproto.ReqDoAction")
}

func init() { proto.RegisterFile("req.proto", fileDescriptor_a52bd445c288993d) }

var fileDescriptor_a52bd445c288993d = []byte{
	// 189 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2c, 0x4a, 0x2d, 0xd4,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcb, 0x05, 0xd3, 0x4a, 0x52, 0x5c, 0x2c, 0x01, 0x9e,
	0x7e, 0xee, 0x42, 0x42, 0x5c, 0x2c, 0x25, 0x99, 0xb9, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xcc,
	0x41, 0x60, 0xb6, 0x92, 0x1b, 0x17, 0x47, 0x50, 0x6a, 0xa1, 0x4f, 0x7e, 0x7a, 0x66, 0x9e, 0x90,
	0x18, 0x17, 0x5b, 0x69, 0x71, 0x6a, 0x91, 0x67, 0x0a, 0x58, 0x05, 0x67, 0x10, 0x94, 0x27, 0xa4,
	0xc4, 0xc5, 0x03, 0x62, 0x05, 0x24, 0x16, 0x17, 0x97, 0xe7, 0x17, 0xa5, 0x48, 0x30, 0x81, 0x65,
	0x51, 0xc4, 0x94, 0xb4, 0xb8, 0x78, 0x82, 0x52, 0x0b, 0x5d, 0xf3, 0x4a, 0x52, 0x8b, 0x82, 0xf2,
	0xf3, 0x73, 0x85, 0xa4, 0xb8, 0x38, 0x8a, 0xf2, 0xf3, 0x73, 0x43, 0x2a, 0x0b, 0x20, 0xf6, 0xb1,
	0x06, 0xc1, 0xf9, 0x4a, 0xbe, 0x5c, 0xdc, 0x41, 0xa9, 0x85, 0x2e, 0xf9, 0x8e, 0xc9, 0x25, 0x99,
	0xf9, 0x79, 0x42, 0x72, 0x5c, 0x5c, 0x89, 0x60, 0x96, 0x73, 0x62, 0x51, 0x0a, 0x54, 0x31, 0x92,
	0x08, 0x42, 0x1e, 0x6c, 0x18, 0x13, 0xb2, 0x3c, 0x48, 0x24, 0x89, 0x0d, 0xec, 0x4b, 0x63, 0x40,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xe0, 0x47, 0x41, 0x5a, 0xfa, 0x00, 0x00, 0x00,
}
