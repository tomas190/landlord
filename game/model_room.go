package game

import (
	"fmt"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/roomType"
	"landlord/msg/mproto"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2/bson"
)

type Room struct {
	mu                sync.Mutex
	RoomClass         *RoomClassify      // 房间分类
	RoomId            string             // 房间ID
	RoundId           string             // roundId
	Players           map[string]*Player // 当前玩家
	LandlordPlayerId  string             // 地主玩家Id
	ThrowCards        []*Card            // 弃牌堆
	BottomCards       []*Card            // 地主牌(及最后三张牌)
	reStartNum        int32              // 都不叫重发次数
	EffectiveCard     []*Card            // 有效牌
	EffectiveType     int32              // 有效牌牌型
	EffectivePlayerId string             // 有效的玩家Id
	MultiAll          int32              // 当局游戏总倍数
	MultiGetLandlord  int32              // 当局抢地主倍数
	MultiBoom         int32              // 炸弹倍数
	MultiSpring       int32              // 是否春天
	MultiRocket       int32              // 火箭倍数
	LandlordOutNum    int32              // 地主出了多少首牌 用于计算春天
	Status            int32              // 房间状态 0 等待中 1叫地主 2.抢地主, 3正在玩
	OutNum            int32              // 总共出了多少次
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
	room.RoundId = fmt.Sprintf("room-%d-%d-%s_2", rType, time.Now().Unix(), bson.NewObjectId().Hex())
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
		//return 0.6
		return 8
	case roomType.LowField:
		return 20
	case roomType.MidField:
		return 50
	case roomType.HighField:
		return 100
	}
	return 0
}

func GetRoomNum(r *Room) int32 {
	r.mu.Lock()
	defer r.mu.Unlock()
	outNum := r.OutNum
	return outNum

}

func SetRoomNum(r *Room) {
	o := GetRoomNum(r)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.OutNum = o + 1
}

func kickRoomByUserID(playerId string) {
	agent := GetAgent(playerId)

	room, b := IsPlayerInRoom(playerId)
	if b && room != nil {
		RemoveRoom(room.RoomId)
		if agent != nil {
			for _, v := range room.Players {
				if v != nil {
					if !v.IsRobot {
						SetSessionRoomId(v.Session, "")
						SendErrMsg(v.Session, msgIdConst.ErrMsg, "系统已将你踢出,请重新登录游戏.")
						// UserLogoutCenterAfterUnlockMoney(v.PlayerInfo.PlayerId, v.PlayerInfo.Gold)
					}
				}
			}
		}
	}
	if agent != nil {
		var push mproto.CloseConn
		push.Msg = "更新维护,请稍后进入!"
		bytes, _ := proto.Marshal(&push)
		msg := PkgMsg(msgIdConst.CloseConn, bytes)
		agent.CloseWithMsg(msg)
	}
}
