package game

import (
	"github.com/google/uuid"
	"github.com/wonderivan/logger"
	"landlord/mconst/roomType"
)

type Room struct {
	RoomClass        *RoomClassify      // 房间分类
	RoomId           string             // 房间ID
	Players          map[string]*Player // 当前玩家
	LandlordPlayerId string             // 地主玩家Id
	BottomCards      []*Card            // 地主牌(及最后三张牌)
	EffectiveCard    []*Card            // 有效牌
	EffectiveType    int32              // 有效牌牌型
	MultiAll         int32              // 当局游戏总倍数
	MultiGetLandlord int32              // 当局抢地主倍数
	MultiBoom        int32              // 炸弹倍数
	MultiSpring      int32              // 是否春天
	LandlordOutNum   int32              // 地主出了多少首牌 用于计算春天
	Status           int32              // 房间状态 0 等待中 1叫地主 2.抢地主, 3正在玩
}

type RoomClassify struct {
	BottomPoint      float64 // 底分
	RoomType         int32
	BottomEnterPoint float64 // 最低入场金币
}

// 创建一个新的房间
func NewRoom(rType int32, players map[string]*Player) *Room {
	//if len(cards) != 3 {
	//	logger.Debug("无法创建房间,房间人数异常:", len(cards))
	//	return nil
	//}

	var room Room
	room.RoomId = uuid.New().String()
	room.MultiAll = 3 // 初始倍数是3
	room.MultiGetLandlord = 3
	room.Players = players
	room.RoomClass = NewRoomClassify(rType)
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
