package game

import (
	"landlord/mconst/roomType"
	"github.com/wonderivan/logger"
	"strconv"
)

type Room struct {
	RoomClass         *RoomClassify         // 房间分类
	RoomId            string                // 房间ID
	Players           []*Player             // 当前玩家
	//ControlChan       chan PlayerActionChan // 控制流
	RemainingMahjongs []int32               // 当前房间剩余麻将
	Status            int32                 // 房间状态  1 正在玩 2
}

type RoomClassify struct {
	BottomPoint      float64 // 底分
	RoomType         int32
	BottomEnterPoint float64 // 最低入场金币
}

// 创建一个新的房间
func NewRoom(rType int32, players []*Player) *Room {
	if len(players) != 2 {
		logger.Debug("无法创建房间,房间人数异常:", len(players))
		return nil
	}

	var room Room
	room.RoomId = strconv.Itoa(int(rType))
	room.Players = players
	room.RoomClass = NewRoomClassify(rType)
//	room.ControlChan = make(chan PlayerActionChan)
	return &room
}

// 创建一个房间分类
func NewRoomClassify(rType int32) *RoomClassify {
	var result RoomClassify
	result.BottomEnterPoint = GetRoomClassifyBottomEnterPoint(rType)
	result.BottomPoint = GetRoomClassifyBottomPoint(rType)
	result.RoomType = rType
	return &result
}

/*===================================  help  func =======================*/
/*
	ExperienceField = 1 // 体验场
	LowField        = 2 // 初级场
	MidField        = 3 // 中级场
	HighField       = 4 // 高级场
*/
// 获取房间底分
func GetRoomClassifyBottomPoint(rType int32) float64 {
	if rType < roomType.ExperienceField || rType > roomType.HighField {
		logger.Error("未知的房间类型")
		return 0
	}

	switch rType {
	case roomType.ExperienceField:
		return 0.01
	case roomType.LowField:
		return 0.5
	case roomType.MidField:
		return 1
	case roomType.HighField:
		return 5
	}
	return 0
}

// 获取房间最低入场分
func GetRoomClassifyBottomEnterPoint(rType int32) float64 {
	if rType < roomType.ExperienceField || rType > roomType.HighField {
		logger.Error("未知的房间类型")
		return 0
	}

	switch rType {
	case roomType.ExperienceField:
		return 0.6
	case roomType.LowField:
		return 10
	case roomType.MidField:
		return 20
	case roomType.HighField:
		return 100
	}
	return 0
}
