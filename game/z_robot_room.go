package game

import (
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
	"time"
)

type WaitRoomChan struct {
	IsClose  bool
	WaitChan chan struct{}
}

// 处理玩家进入房间 和机器人玩
func DealPlayerEnterRoomWithRobot(session *melody.Session, playerInfo PlayerInfo, roomType int32, waitRoom *WaitRoomChan) {
	select {
	case <-waitRoom.WaitChan:
		logger.Debug("============= 玩家在等待过程中退出了房间 =============")
		waitRoom.IsClose = true
		close(waitRoom.WaitChan)
	case <-time.After(time.Second * getWaitTimePlayerEnterRoom()):
		waitRoom.IsClose = true
		close(waitRoom.WaitChan)
		PlayWithRobot(session, playerInfo, roomType)
	}

}

// 匹配机器人并创建房间
func PlayWithRobot(session *melody.Session, playerInfo PlayerInfo, roomType int32) {
	/*
		todo 是否宰羊模式
	*/

	var p Player
	p.PlayerInfo = &playerInfo
	p.Session = session
	p.PlayerPosition = 1
	p.ActionChan = make(chan PlayerActionChan)

	robot2 := CreateRobot()
	robot2.PlayerPosition = 2

	robot3 := CreateRobot()
	robot3.PlayerPosition = 3

	mp := make(map[string]*Player, 3)
	mp[playerInfo.PlayerId] = &p
	mp[robot2.PlayerInfo.PlayerId] = robot2
	mp[robot3.PlayerInfo.PlayerId] = robot3

	room := NewRoom(roomType, mp)

	// 设置用户全局房间Id
	SetSessionRoomId(session, room.RoomId)
	// 保存房间
	SaveRoom(room.RoomId, room)

	// 开启线程 游戏开始
	go PlayGameWithRobot(room)
}

// 开始和机器人玩游戏
func PlayGameWithRobot(room *Room) {

	// 1. 玩家进入房间如果有玩家正待等待则与之开始游戏
	PushPlayerEnterRoom(room)
	DelaySomeTime(1)

	// 2.给玩家发牌
	//PushPlayerStartGameWithRobot(room)
	PushPlayerStartGameWithRobot2(room)
	// ..．流程控制到这里结束　发牌  抢地主  玩牌 直接由 PushPlayerStartGame 开始 且循环

}
