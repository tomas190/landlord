// Code generated by protoc-gen-go. DO NOT EDIT.
// source: resp.proto

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

// =================================  SERVER 收到指令返回 200 pong 消息为1 开始 ===========
// SERVER RESP : MSG_ID 1
// 1.PONG
type PONG struct {
	Time                 int64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PONG) Reset()         { *m = PONG{} }
func (m *PONG) String() string { return proto.CompactTextString(m) }
func (*PONG) ProtoMessage()    {}
func (*PONG) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{0}
}

func (m *PONG) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PONG.Unmarshal(m, b)
}
func (m *PONG) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PONG.Marshal(b, m, deterministic)
}
func (m *PONG) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PONG.Merge(m, src)
}
func (m *PONG) XXX_Size() int {
	return xxx_messageInfo_PONG.Size(m)
}
func (m *PONG) XXX_DiscardUnknown() {
	xxx_messageInfo_PONG.DiscardUnknown(m)
}

var xxx_messageInfo_PONG proto.InternalMessageInfo

func (m *PONG) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

// SERVER RESP : MSG_ID 200
// 1.登录返回
type RespLogin struct {
	PlayerInfo           *PlayerInfo `protobuf:"bytes,1,opt,name=playerInfo,proto3" json:"playerInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *RespLogin) Reset()         { *m = RespLogin{} }
func (m *RespLogin) String() string { return proto.CompactTextString(m) }
func (*RespLogin) ProtoMessage()    {}
func (*RespLogin) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{1}
}

func (m *RespLogin) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RespLogin.Unmarshal(m, b)
}
func (m *RespLogin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RespLogin.Marshal(b, m, deterministic)
}
func (m *RespLogin) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RespLogin.Merge(m, src)
}
func (m *RespLogin) XXX_Size() int {
	return xxx_messageInfo_RespLogin.Size(m)
}
func (m *RespLogin) XXX_DiscardUnknown() {
	xxx_messageInfo_RespLogin.DiscardUnknown(m)
}

var xxx_messageInfo_RespLogin proto.InternalMessageInfo

func (m *RespLogin) GetPlayerInfo() *PlayerInfo {
	if m != nil {
		return m.PlayerInfo
	}
	return nil
}

// SERVER PUSH : MSG_ID 300  // 推送房间类型列表
// 1.房间分类信息推送
type PushRoomClassify struct {
	RoomClassify         []*RoomClassify `protobuf:"bytes,1,rep,name=roomClassify,proto3" json:"roomClassify,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *PushRoomClassify) Reset()         { *m = PushRoomClassify{} }
func (m *PushRoomClassify) String() string { return proto.CompactTextString(m) }
func (*PushRoomClassify) ProtoMessage()    {}
func (*PushRoomClassify) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{2}
}

func (m *PushRoomClassify) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushRoomClassify.Unmarshal(m, b)
}
func (m *PushRoomClassify) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushRoomClassify.Marshal(b, m, deterministic)
}
func (m *PushRoomClassify) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushRoomClassify.Merge(m, src)
}
func (m *PushRoomClassify) XXX_Size() int {
	return xxx_messageInfo_PushRoomClassify.Size(m)
}
func (m *PushRoomClassify) XXX_DiscardUnknown() {
	xxx_messageInfo_PushRoomClassify.DiscardUnknown(m)
}

var xxx_messageInfo_PushRoomClassify proto.InternalMessageInfo

func (m *PushRoomClassify) GetRoomClassify() []*RoomClassify {
	if m != nil {
		return m.RoomClassify
	}
	return nil
}

// SERVER PUSH : MSG_ID 301  // 推送房间玩家信息
// 1.推送房间玩家信息
type PushRoomPlayer struct {
	Players              []*RoomPlayer `protobuf:"bytes,1,rep,name=players,proto3" json:"players,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *PushRoomPlayer) Reset()         { *m = PushRoomPlayer{} }
func (m *PushRoomPlayer) String() string { return proto.CompactTextString(m) }
func (*PushRoomPlayer) ProtoMessage()    {}
func (*PushRoomPlayer) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{3}
}

func (m *PushRoomPlayer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushRoomPlayer.Unmarshal(m, b)
}
func (m *PushRoomPlayer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushRoomPlayer.Marshal(b, m, deterministic)
}
func (m *PushRoomPlayer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushRoomPlayer.Merge(m, src)
}
func (m *PushRoomPlayer) XXX_Size() int {
	return xxx_messageInfo_PushRoomPlayer.Size(m)
}
func (m *PushRoomPlayer) XXX_DiscardUnknown() {
	xxx_messageInfo_PushRoomPlayer.DiscardUnknown(m)
}

var xxx_messageInfo_PushRoomPlayer proto.InternalMessageInfo

func (m *PushRoomPlayer) GetPlayers() []*RoomPlayer {
	if m != nil {
		return m.Players
	}
	return nil
}

// SERVER PUSH : MSG_ID 302  // 发牌
// 2.发牌
type PushStartGame struct {
	Cards                []*Card  `protobuf:"bytes,2,rep,name=cards,proto3" json:"cards,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushStartGame) Reset()         { *m = PushStartGame{} }
func (m *PushStartGame) String() string { return proto.CompactTextString(m) }
func (*PushStartGame) ProtoMessage()    {}
func (*PushStartGame) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{4}
}

func (m *PushStartGame) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushStartGame.Unmarshal(m, b)
}
func (m *PushStartGame) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushStartGame.Marshal(b, m, deterministic)
}
func (m *PushStartGame) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushStartGame.Merge(m, src)
}
func (m *PushStartGame) XXX_Size() int {
	return xxx_messageInfo_PushStartGame.Size(m)
}
func (m *PushStartGame) XXX_DiscardUnknown() {
	xxx_messageInfo_PushStartGame.DiscardUnknown(m)
}

var xxx_messageInfo_PushStartGame proto.InternalMessageInfo

func (m *PushStartGame) GetCards() []*Card {
	if m != nil {
		return m.Cards
	}
	return nil
}

// SERVER PUSH : MSG_ID 303  // 抢地主阶段
// 3.抢地主阶段
type PushGetLandlord struct {
	LastPlayerId         string   `protobuf:"bytes,1,opt,name=lastPlayerId,proto3" json:"lastPlayerId,omitempty"`
	LastPlayerPosition   int32    `protobuf:"varint,2,opt,name=lastPlayerPosition,proto3" json:"lastPlayerPosition,omitempty"`
	LastPlayerAction     int32    `protobuf:"varint,3,opt,name=lastPlayerAction,proto3" json:"lastPlayerAction,omitempty"`
	PlayerId             string   `protobuf:"bytes,4,opt,name=playerId,proto3" json:"playerId,omitempty"`
	PlayerPosition       int32    `protobuf:"varint,5,opt,name=playerPosition,proto3" json:"playerPosition,omitempty"`
	Action               int32    `protobuf:"varint,6,opt,name=action,proto3" json:"action,omitempty"`
	Multi                int32    `protobuf:"varint,7,opt,name=multi,proto3" json:"multi,omitempty"`
	Countdown            int32    `protobuf:"varint,8,opt,name=Countdown,proto3" json:"Countdown,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushGetLandlord) Reset()         { *m = PushGetLandlord{} }
func (m *PushGetLandlord) String() string { return proto.CompactTextString(m) }
func (*PushGetLandlord) ProtoMessage()    {}
func (*PushGetLandlord) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{5}
}

func (m *PushGetLandlord) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushGetLandlord.Unmarshal(m, b)
}
func (m *PushGetLandlord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushGetLandlord.Marshal(b, m, deterministic)
}
func (m *PushGetLandlord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushGetLandlord.Merge(m, src)
}
func (m *PushGetLandlord) XXX_Size() int {
	return xxx_messageInfo_PushGetLandlord.Size(m)
}
func (m *PushGetLandlord) XXX_DiscardUnknown() {
	xxx_messageInfo_PushGetLandlord.DiscardUnknown(m)
}

var xxx_messageInfo_PushGetLandlord proto.InternalMessageInfo

func (m *PushGetLandlord) GetLastPlayerId() string {
	if m != nil {
		return m.LastPlayerId
	}
	return ""
}

func (m *PushGetLandlord) GetLastPlayerPosition() int32 {
	if m != nil {
		return m.LastPlayerPosition
	}
	return 0
}

func (m *PushGetLandlord) GetLastPlayerAction() int32 {
	if m != nil {
		return m.LastPlayerAction
	}
	return 0
}

func (m *PushGetLandlord) GetPlayerId() string {
	if m != nil {
		return m.PlayerId
	}
	return ""
}

func (m *PushGetLandlord) GetPlayerPosition() int32 {
	if m != nil {
		return m.PlayerPosition
	}
	return 0
}

func (m *PushGetLandlord) GetAction() int32 {
	if m != nil {
		return m.Action
	}
	return 0
}

func (m *PushGetLandlord) GetMulti() int32 {
	if m != nil {
		return m.Multi
	}
	return 0
}

func (m *PushGetLandlord) GetCountdown() int32 {
	if m != nil {
		return m.Countdown
	}
	return 0
}

// SERVER PUSH : MSG_ID 304  // 地主推送
// 4.地主推送
type PushLandlord struct {
	LandlordId           string   `protobuf:"bytes,1,opt,name=landlordId,proto3" json:"landlordId,omitempty"`
	Position             int32    `protobuf:"varint,2,opt,name=position,proto3" json:"position,omitempty"`
	Cards                []*Card  `protobuf:"bytes,3,rep,name=cards,proto3" json:"cards,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushLandlord) Reset()         { *m = PushLandlord{} }
func (m *PushLandlord) String() string { return proto.CompactTextString(m) }
func (*PushLandlord) ProtoMessage()    {}
func (*PushLandlord) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{6}
}

func (m *PushLandlord) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushLandlord.Unmarshal(m, b)
}
func (m *PushLandlord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushLandlord.Marshal(b, m, deterministic)
}
func (m *PushLandlord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushLandlord.Merge(m, src)
}
func (m *PushLandlord) XXX_Size() int {
	return xxx_messageInfo_PushLandlord.Size(m)
}
func (m *PushLandlord) XXX_DiscardUnknown() {
	xxx_messageInfo_PushLandlord.DiscardUnknown(m)
}

var xxx_messageInfo_PushLandlord proto.InternalMessageInfo

func (m *PushLandlord) GetLandlordId() string {
	if m != nil {
		return m.LandlordId
	}
	return ""
}

func (m *PushLandlord) GetPosition() int32 {
	if m != nil {
		return m.Position
	}
	return 0
}

func (m *PushLandlord) GetCards() []*Card {
	if m != nil {
		return m.Cards
	}
	return nil
}

// SERVER PUSH : MSG_ID 305  // 出牌阶段
// 5.出牌阶段
type PushOutCard struct {
	LastPlayerId         string   `protobuf:"bytes,1,opt,name=lastPlayerId,proto3" json:"lastPlayerId,omitempty"`
	LastPlayerPosition   int32    `protobuf:"varint,2,opt,name=lastPlayerPosition,proto3" json:"lastPlayerPosition,omitempty"`
	LastPlayerCards      []*Card  `protobuf:"bytes,3,rep,name=lastPlayerCards,proto3" json:"lastPlayerCards,omitempty"`
	LastPlayerCardsType  int32    `protobuf:"varint,4,opt,name=lastPlayerCardsType,proto3" json:"lastPlayerCardsType,omitempty"`
	LastAction           int32    `protobuf:"varint,5,opt,name=lastAction,proto3" json:"lastAction,omitempty"`
	LastRemainLen        int32    `protobuf:"varint,6,opt,name=lastRemainLen,proto3" json:"lastRemainLen,omitempty"`
	PlayerId             string   `protobuf:"bytes,7,opt,name=playerId,proto3" json:"playerId,omitempty"`
	PlayerPosition       int32    `protobuf:"varint,8,opt,name=playerPosition,proto3" json:"playerPosition,omitempty"`
	IsMustPlay           bool     `protobuf:"varint,9,opt,name=isMustPlay,proto3" json:"isMustPlay,omitempty"`
	Multi                int32    `protobuf:"varint,10,opt,name=Multi,proto3" json:"Multi,omitempty"`
	Countdown            int32    `protobuf:"varint,11,opt,name=Countdown,proto3" json:"Countdown,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushOutCard) Reset()         { *m = PushOutCard{} }
func (m *PushOutCard) String() string { return proto.CompactTextString(m) }
func (*PushOutCard) ProtoMessage()    {}
func (*PushOutCard) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{7}
}

func (m *PushOutCard) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushOutCard.Unmarshal(m, b)
}
func (m *PushOutCard) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushOutCard.Marshal(b, m, deterministic)
}
func (m *PushOutCard) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushOutCard.Merge(m, src)
}
func (m *PushOutCard) XXX_Size() int {
	return xxx_messageInfo_PushOutCard.Size(m)
}
func (m *PushOutCard) XXX_DiscardUnknown() {
	xxx_messageInfo_PushOutCard.DiscardUnknown(m)
}

var xxx_messageInfo_PushOutCard proto.InternalMessageInfo

func (m *PushOutCard) GetLastPlayerId() string {
	if m != nil {
		return m.LastPlayerId
	}
	return ""
}

func (m *PushOutCard) GetLastPlayerPosition() int32 {
	if m != nil {
		return m.LastPlayerPosition
	}
	return 0
}

func (m *PushOutCard) GetLastPlayerCards() []*Card {
	if m != nil {
		return m.LastPlayerCards
	}
	return nil
}

func (m *PushOutCard) GetLastPlayerCardsType() int32 {
	if m != nil {
		return m.LastPlayerCardsType
	}
	return 0
}

func (m *PushOutCard) GetLastAction() int32 {
	if m != nil {
		return m.LastAction
	}
	return 0
}

func (m *PushOutCard) GetLastRemainLen() int32 {
	if m != nil {
		return m.LastRemainLen
	}
	return 0
}

func (m *PushOutCard) GetPlayerId() string {
	if m != nil {
		return m.PlayerId
	}
	return ""
}

func (m *PushOutCard) GetPlayerPosition() int32 {
	if m != nil {
		return m.PlayerPosition
	}
	return 0
}

func (m *PushOutCard) GetIsMustPlay() bool {
	if m != nil {
		return m.IsMustPlay
	}
	return false
}

func (m *PushOutCard) GetMulti() int32 {
	if m != nil {
		return m.Multi
	}
	return 0
}

func (m *PushOutCard) GetCountdown() int32 {
	if m != nil {
		return m.Countdown
	}
	return 0
}

// SERVER PUSH : MSG_ID 500  // 服务器断开指令
// 1.断开连接指令
type CloseConn struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CloseConn) Reset()         { *m = CloseConn{} }
func (m *CloseConn) String() string { return proto.CompactTextString(m) }
func (*CloseConn) ProtoMessage()    {}
func (*CloseConn) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{8}
}

func (m *CloseConn) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CloseConn.Unmarshal(m, b)
}
func (m *CloseConn) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CloseConn.Marshal(b, m, deterministic)
}
func (m *CloseConn) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CloseConn.Merge(m, src)
}
func (m *CloseConn) XXX_Size() int {
	return xxx_messageInfo_CloseConn.Size(m)
}
func (m *CloseConn) XXX_DiscardUnknown() {
	xxx_messageInfo_CloseConn.DiscardUnknown(m)
}

var xxx_messageInfo_CloseConn proto.InternalMessageInfo

func (m *CloseConn) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *CloseConn) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

// SERVER PUSH : MSG_ID 501  // 服务器错误指令
// 1.服务器错误指令
type ErrMsg struct {
	MsgId                int32    `protobuf:"varint,1,opt,name=msgId,proto3" json:"msgId,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ErrMsg) Reset()         { *m = ErrMsg{} }
func (m *ErrMsg) String() string { return proto.CompactTextString(m) }
func (*ErrMsg) ProtoMessage()    {}
func (*ErrMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{9}
}

func (m *ErrMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ErrMsg.Unmarshal(m, b)
}
func (m *ErrMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ErrMsg.Marshal(b, m, deterministic)
}
func (m *ErrMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrMsg.Merge(m, src)
}
func (m *ErrMsg) XXX_Size() int {
	return xxx_messageInfo_ErrMsg.Size(m)
}
func (m *ErrMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrMsg.DiscardUnknown(m)
}

var xxx_messageInfo_ErrMsg proto.InternalMessageInfo

func (m *ErrMsg) GetMsgId() int32 {
	if m != nil {
		return m.MsgId
	}
	return 0
}

func (m *ErrMsg) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

// ============================================= 属性data =================================
// 1.PushStartGame
type RoomPlayer struct {
	Players              *PlayerInfo `protobuf:"bytes,1,opt,name=players,proto3" json:"players,omitempty"`
	Position             int32       `protobuf:"varint,2,opt,name=Position,proto3" json:"Position,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *RoomPlayer) Reset()         { *m = RoomPlayer{} }
func (m *RoomPlayer) String() string { return proto.CompactTextString(m) }
func (*RoomPlayer) ProtoMessage()    {}
func (*RoomPlayer) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{10}
}

func (m *RoomPlayer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomPlayer.Unmarshal(m, b)
}
func (m *RoomPlayer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomPlayer.Marshal(b, m, deterministic)
}
func (m *RoomPlayer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomPlayer.Merge(m, src)
}
func (m *RoomPlayer) XXX_Size() int {
	return xxx_messageInfo_RoomPlayer.Size(m)
}
func (m *RoomPlayer) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomPlayer.DiscardUnknown(m)
}

var xxx_messageInfo_RoomPlayer proto.InternalMessageInfo

func (m *RoomPlayer) GetPlayers() *PlayerInfo {
	if m != nil {
		return m.Players
	}
	return nil
}

func (m *RoomPlayer) GetPosition() int32 {
	if m != nil {
		return m.Position
	}
	return 0
}

// 1. LoginResp   2.PlayerListPush
type PlayerInfo struct {
	PlayerId             string   `protobuf:"bytes,1,opt,name=playerId,proto3" json:"playerId,omitempty"`
	PlayerName           string   `protobuf:"bytes,2,opt,name=playerName,proto3" json:"playerName,omitempty"`
	PlayerImg            string   `protobuf:"bytes,3,opt,name=playerImg,proto3" json:"playerImg,omitempty"`
	Gold                 float64  `protobuf:"fixed64,4,opt,name=gold,proto3" json:"gold,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlayerInfo) Reset()         { *m = PlayerInfo{} }
func (m *PlayerInfo) String() string { return proto.CompactTextString(m) }
func (*PlayerInfo) ProtoMessage()    {}
func (*PlayerInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{11}
}

func (m *PlayerInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerInfo.Unmarshal(m, b)
}
func (m *PlayerInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerInfo.Marshal(b, m, deterministic)
}
func (m *PlayerInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerInfo.Merge(m, src)
}
func (m *PlayerInfo) XXX_Size() int {
	return xxx_messageInfo_PlayerInfo.Size(m)
}
func (m *PlayerInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerInfo proto.InternalMessageInfo

func (m *PlayerInfo) GetPlayerId() string {
	if m != nil {
		return m.PlayerId
	}
	return ""
}

func (m *PlayerInfo) GetPlayerName() string {
	if m != nil {
		return m.PlayerName
	}
	return ""
}

func (m *PlayerInfo) GetPlayerImg() string {
	if m != nil {
		return m.PlayerImg
	}
	return ""
}

func (m *PlayerInfo) GetGold() float64 {
	if m != nil {
		return m.Gold
	}
	return 0
}

// RoomClassifyPush
type RoomClassify struct {
	RoomType             int32    `protobuf:"varint,1,opt,name=roomType,proto3" json:"roomType,omitempty"`
	BottomPoint          float64  `protobuf:"fixed64,2,opt,name=bottomPoint,proto3" json:"bottomPoint,omitempty"`
	BottomEnterPoint     float64  `protobuf:"fixed64,3,opt,name=bottomEnterPoint,proto3" json:"bottomEnterPoint,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoomClassify) Reset()         { *m = RoomClassify{} }
func (m *RoomClassify) String() string { return proto.CompactTextString(m) }
func (*RoomClassify) ProtoMessage()    {}
func (*RoomClassify) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{12}
}

func (m *RoomClassify) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomClassify.Unmarshal(m, b)
}
func (m *RoomClassify) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomClassify.Marshal(b, m, deterministic)
}
func (m *RoomClassify) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomClassify.Merge(m, src)
}
func (m *RoomClassify) XXX_Size() int {
	return xxx_messageInfo_RoomClassify.Size(m)
}
func (m *RoomClassify) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomClassify.DiscardUnknown(m)
}

var xxx_messageInfo_RoomClassify proto.InternalMessageInfo

func (m *RoomClassify) GetRoomType() int32 {
	if m != nil {
		return m.RoomType
	}
	return 0
}

func (m *RoomClassify) GetBottomPoint() float64 {
	if m != nil {
		return m.BottomPoint
	}
	return 0
}

func (m *RoomClassify) GetBottomEnterPoint() float64 {
	if m != nil {
		return m.BottomEnterPoint
	}
	return 0
}

// PushStartGame
type Card struct {
	Value                int32    `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	Suit                 int32    `protobuf:"varint,2,opt,name=suit,proto3" json:"suit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Card) Reset()         { *m = Card{} }
func (m *Card) String() string { return proto.CompactTextString(m) }
func (*Card) ProtoMessage()    {}
func (*Card) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c5365792f61ddff, []int{13}
}

func (m *Card) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Card.Unmarshal(m, b)
}
func (m *Card) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Card.Marshal(b, m, deterministic)
}
func (m *Card) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Card.Merge(m, src)
}
func (m *Card) XXX_Size() int {
	return xxx_messageInfo_Card.Size(m)
}
func (m *Card) XXX_DiscardUnknown() {
	xxx_messageInfo_Card.DiscardUnknown(m)
}

var xxx_messageInfo_Card proto.InternalMessageInfo

func (m *Card) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *Card) GetSuit() int32 {
	if m != nil {
		return m.Suit
	}
	return 0
}

func init() {
	proto.RegisterType((*PONG)(nil), "mproto.PONG")
	proto.RegisterType((*RespLogin)(nil), "mproto.RespLogin")
	proto.RegisterType((*PushRoomClassify)(nil), "mproto.PushRoomClassify")
	proto.RegisterType((*PushRoomPlayer)(nil), "mproto.PushRoomPlayer")
	proto.RegisterType((*PushStartGame)(nil), "mproto.PushStartGame")
	proto.RegisterType((*PushGetLandlord)(nil), "mproto.PushGetLandlord")
	proto.RegisterType((*PushLandlord)(nil), "mproto.PushLandlord")
	proto.RegisterType((*PushOutCard)(nil), "mproto.PushOutCard")
	proto.RegisterType((*CloseConn)(nil), "mproto.CloseConn")
	proto.RegisterType((*ErrMsg)(nil), "mproto.ErrMsg")
	proto.RegisterType((*RoomPlayer)(nil), "mproto.RoomPlayer")
	proto.RegisterType((*PlayerInfo)(nil), "mproto.PlayerInfo")
	proto.RegisterType((*RoomClassify)(nil), "mproto.RoomClassify")
	proto.RegisterType((*Card)(nil), "mproto.Card")
}

func init() { proto.RegisterFile("resp.proto", fileDescriptor_3c5365792f61ddff) }

var fileDescriptor_3c5365792f61ddff = []byte{
	// 642 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0x5f, 0x6b, 0xdb, 0x3e,
	0x14, 0xc5, 0x75, 0x9c, 0xc6, 0x37, 0xe9, 0x1f, 0xf4, 0x2b, 0x3f, 0x4c, 0x19, 0x25, 0x88, 0x31,
	0xc2, 0x18, 0xa1, 0x6b, 0x61, 0xec, 0x69, 0x63, 0x84, 0x52, 0x06, 0x49, 0x1b, 0xb4, 0xb1, 0x77,
	0xb5, 0x76, 0x33, 0x83, 0x2d, 0x19, 0x4b, 0xde, 0xd6, 0x7d, 0x9b, 0x3d, 0xee, 0x5b, 0x0e, 0x5d,
	0xf9, 0x7f, 0x13, 0xf6, 0xb4, 0xa7, 0xe8, 0x1c, 0x1d, 0xe9, 0x2a, 0xe7, 0x1e, 0x5f, 0x80, 0x3c,
	0x52, 0xd9, 0x3c, 0xcb, 0xa5, 0x96, 0x64, 0x98, 0xe2, 0x2f, 0x3d, 0x85, 0xc1, 0xfa, 0xf6, 0xe6,
	0x9a, 0x10, 0x18, 0xe8, 0x38, 0x8d, 0x02, 0x67, 0xea, 0xcc, 0x5c, 0x86, 0x6b, 0xfa, 0x1e, 0x7c,
	0x16, 0xa9, 0x6c, 0x29, 0x37, 0xb1, 0x20, 0x17, 0x00, 0x59, 0xc2, 0x1f, 0xa3, 0xfc, 0xa3, 0x78,
	0x90, 0x28, 0x1b, 0x5f, 0x90, 0xb9, 0xbd, 0x65, 0xbe, 0xae, 0x77, 0x58, 0x4b, 0x45, 0x97, 0x70,
	0xbc, 0x2e, 0xd4, 0x57, 0x26, 0x65, 0xba, 0x48, 0xb8, 0x52, 0xf1, 0xc3, 0x23, 0x79, 0x0b, 0x93,
	0xbc, 0x85, 0x03, 0x67, 0xea, 0xce, 0xc6, 0x17, 0x27, 0xd5, 0x4d, 0x6d, 0x2d, 0xeb, 0x28, 0xe9,
	0x3b, 0x38, 0xac, 0x6e, 0xb3, 0xf5, 0xc8, 0x2b, 0xd8, 0xb7, 0xd5, 0x54, 0x79, 0x0d, 0x69, 0x5f,
	0x63, 0x45, 0xac, 0x92, 0xd0, 0x4b, 0x38, 0x30, 0xe7, 0x3f, 0x69, 0x9e, 0xeb, 0x6b, 0x9e, 0x46,
	0x84, 0x82, 0x77, 0xcf, 0xf3, 0x50, 0x05, 0x7b, 0x78, 0x78, 0x52, 0x1d, 0x5e, 0xf0, 0x3c, 0x64,
	0x76, 0x8b, 0xfe, 0xda, 0x83, 0x23, 0x73, 0xea, 0x3a, 0xd2, 0x4b, 0x2e, 0xc2, 0x44, 0xe6, 0x21,
	0xa1, 0x30, 0x49, 0xb8, 0xd2, 0xe5, 0x9f, 0x0e, 0xd1, 0x0c, 0x9f, 0x75, 0x38, 0x32, 0x07, 0xd2,
	0xe0, 0xb5, 0x54, 0xb1, 0x8e, 0xa5, 0x08, 0xf6, 0xa6, 0xce, 0xcc, 0x63, 0x5b, 0x76, 0xc8, 0x4b,
	0x38, 0x6e, 0xd8, 0x0f, 0xf7, 0xa8, 0x76, 0x51, 0xfd, 0x84, 0x27, 0xa7, 0x30, 0xca, 0xaa, 0xda,
	0x03, 0xac, 0x5d, 0x63, 0xf2, 0x02, 0x0e, 0xb3, 0x6e, 0x4d, 0x0f, 0x6f, 0xe9, 0xb1, 0xe4, 0x7f,
	0x18, 0x72, 0x5b, 0x65, 0x88, 0xfb, 0x25, 0x22, 0x27, 0xe0, 0xa5, 0x45, 0xa2, 0xe3, 0x60, 0x1f,
	0x69, 0x0b, 0xc8, 0x33, 0xf0, 0x17, 0xb2, 0x10, 0x3a, 0x94, 0xdf, 0x45, 0x30, 0xc2, 0x9d, 0x86,
	0xa0, 0x02, 0x26, 0xc6, 0xa2, 0xda, 0x9f, 0x33, 0x80, 0xa4, 0x5c, 0xd7, 0xee, 0xb4, 0x18, 0x7c,
	0x7f, 0xd7, 0x91, 0x1a, 0x37, 0x3d, 0x71, 0x77, 0xf7, 0xe4, 0xb7, 0x0b, 0x63, 0x53, 0xf0, 0xb6,
	0xd0, 0x86, 0xfe, 0x27, 0xfd, 0x78, 0x03, 0x47, 0x0d, 0xbb, 0xd8, 0xf9, 0xa2, 0xbe, 0x88, 0x9c,
	0xc3, 0x7f, 0x3d, 0xea, 0xf3, 0x63, 0x16, 0x61, 0x9b, 0x3c, 0xb6, 0x6d, 0xcb, 0xba, 0xa5, 0x74,
	0xd9, 0x73, 0xdb, 0xad, 0x16, 0x43, 0x9e, 0xc3, 0x81, 0x41, 0x2c, 0x4a, 0x79, 0x2c, 0x96, 0x51,
	0xd5, 0xb0, 0x2e, 0xd9, 0xc9, 0xc4, 0xfe, 0x5f, 0x33, 0x31, 0xda, 0x9a, 0x89, 0x33, 0x80, 0x58,
	0xad, 0x0a, 0xfb, 0xc4, 0xc0, 0x9f, 0x3a, 0xb3, 0x11, 0x6b, 0x31, 0x26, 0x1b, 0x2b, 0xcc, 0x06,
	0xd8, 0x6c, 0xac, 0x9e, 0x66, 0x63, 0xdc, 0xcf, 0xc6, 0x6b, 0xf0, 0x17, 0x89, 0x54, 0xd1, 0x42,
	0x0a, 0x61, 0x86, 0xcc, 0xbd, 0x0c, 0xed, 0x90, 0xf1, 0x18, 0xae, 0xc9, 0x31, 0xb8, 0xa9, 0xda,
	0x60, 0x27, 0x7c, 0x66, 0x96, 0xf4, 0x1c, 0x86, 0x57, 0x79, 0xbe, 0x52, 0x1b, 0x0c, 0xa3, 0xda,
	0x94, 0x1d, 0x35, 0x61, 0x34, 0x60, 0xcb, 0x89, 0x2f, 0x00, 0xbb, 0xa6, 0xc2, 0xae, 0x31, 0x55,
	0x49, 0x8c, 0x71, 0xbd, 0x38, 0xd4, 0x98, 0xfe, 0x04, 0x68, 0x8e, 0x74, 0x2c, 0x76, 0x7a, 0x16,
	0x9f, 0x55, 0xd3, 0xf1, 0x86, 0xa7, 0x51, 0xf9, 0xb4, 0x16, 0x63, 0x4c, 0x2a, 0xb5, 0xe9, 0x06,
	0xbf, 0x6b, 0x9f, 0x35, 0x84, 0xf1, 0x65, 0x23, 0x13, 0xfb, 0x31, 0x3b, 0x0c, 0xd7, 0xf4, 0x07,
	0x4c, 0x3a, 0x73, 0xf3, 0x14, 0x46, 0x66, 0x1a, 0x62, 0x9a, 0xac, 0x1d, 0x35, 0x26, 0x53, 0x18,
	0xdf, 0x49, 0xad, 0x65, 0xba, 0x96, 0xb1, 0xd0, 0x58, 0xde, 0x61, 0x6d, 0xca, 0x8c, 0x17, 0x0b,
	0xaf, 0x84, 0x36, 0x1d, 0x37, 0x32, 0x17, 0x65, 0x4f, 0x78, 0x7a, 0x0e, 0x03, 0xfc, 0xac, 0x4e,
	0xc0, 0xfb, 0xc6, 0x93, 0xa2, 0x2a, 0x67, 0x81, 0x79, 0xab, 0x2a, 0x62, 0x5d, 0x7a, 0x85, 0xeb,
	0xbb, 0x21, 0xda, 0x7b, 0xf9, 0x27, 0x00, 0x00, 0xff, 0xff, 0x75, 0x25, 0x01, 0xa1, 0x61, 0x06,
	0x00, 0x00,
}
