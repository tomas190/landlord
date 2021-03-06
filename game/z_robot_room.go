package game

import (
	"time"

	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
)

//type WaitRoomChan struct {
//	IsClose  bool
//	WaitChan chan struct{}
//}

// 处理玩家进入房间 和机器人玩
//func DealPlayerEnterRoomWithRobot(session *melody.Session, playerInfo PlayerInfo, roomType int32, waitRoom *WaitRoomChan) {
//	destiny := getWaitTimePlayerEnterRoom()
//	//mp := make(map[string]*Player, 3)
//	//room := NewRoom(roomType, mp)
//	// fakerIntoRoom(room, playerInfo, session, destiny)
//	select {
//	case <-waitRoom.WaitChan:
//		logger.Debug("============= 玩家在等待过程中退出了房间 =============")
//		waitRoom.IsClose = true
//		close(waitRoom.WaitChan)
//	case <-time.After(time.Second * destiny):
//		waitRoom.IsClose = true
//		close(waitRoom.WaitChan)
//		PlayWithRobot(session, playerInfo, roomType)
//	}
//
//}

// 匹配机器人并创建房间
//func PlayWithRobot(session *melody.Session, playerInfo PlayerInfo, room *Room) {
func PlayWithRobot(session *melody.Session, playerInfo PlayerInfo, roomType int32) {
	logger.Debug("========  玩家匹配机器人 ========")

	var p Player
	p.PlayerInfo = &playerInfo
	p.Session = session
	p.PlayerPosition = 1
	//p.ActionChan = make(chan PlayerActionChan)

	robot2 := CreateRobot(roomType)
	robot2.PlayerPosition = 2

	robot3 := CreateRobot(roomType)
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
	logger.Debug("和机器人玩游戏")
	// 1. 玩家进入房间如果有玩家正待等待则与之开始游戏
	PushPlayerEnterRoom(room)
	DelaySomeTime(1)

	// 2.给玩家发牌
	// PushPlayerStartGameWithRobot2(room)
	PushPlayerStartGameWithRobot3(room)
	//PushPlayerStartGameWithRobotLast(room)
	// ..．流程控制到这里结束　发牌  抢地主  玩牌 直接由 PushPlayerStartGame 开始 且循环

}

// 用这种多线程会有很多问题
func fakerIntoRoom(room *Room, playerInfo PlayerInfo, session *melody.Session, destiny time.Duration) {
	go func() { // 假操作房间进入
		d := RandNum(0, 10)
		if d <= 7 {
			var p Player
			p.PlayerInfo = &playerInfo
			p.Session = session
			p.PlayerPosition = 1
			//	p.ActionChan = make(chan PlayerActionChan)
			room.Players[p.PlayerInfo.PlayerId] = &p

			// 第一个进入的机器人
			robot2 := CreateRobot(room.RoomClass.RoomType)
			robot2.PlayerPosition = 2
			room.Players[robot2.PlayerInfo.PlayerId] = robot2

			DelaySomeTime(destiny / 3)
			logger.Debug("===================推送第一个机器人")
			PushFakerPlayerEnterRoom(room.Players, &p)

			// 第一个进入的机器人
			robot3 := CreateRobot(room.RoomClass.RoomType)
			robot3.PlayerPosition = 3
			room.Players[robot3.PlayerInfo.PlayerId] = robot3

			DelaySomeTime(destiny / 2)
			logger.Debug("===================推送第二个机器人")
			PushFakerPlayerEnterRoom(room.Players, &p)
		} else {
			logger.Debug("ssssssssssssssssss")
			var p Player
			p.PlayerInfo = &playerInfo
			p.Session = session
			p.PlayerPosition = 1
			//p.ActionChan = make(chan PlayerActionChan)
			room.Players[p.PlayerInfo.PlayerId] = &p

			// 第一个进入的机器人
			robot2 := CreateRobot(room.RoomClass.RoomType)
			robot2.PlayerPosition = 2
			room.Players[robot2.PlayerInfo.PlayerId] = robot2

			DelaySomeTime(destiny / 4)
			logger.Debug("===================推送第一个机器人")
			PushFakerPlayerEnterRoom(room.Players, &p)
			delete(room.Players, robot2.PlayerInfo.PlayerId)

			DelaySomeTime(destiny / 5)
			PushFakerPlayerQuitRoom(&p)

			robotN2 := CreateRobot(room.RoomClass.RoomType)
			robotN2.PlayerPosition = 2
			room.Players[robotN2.PlayerInfo.PlayerId] = robotN2

			DelaySomeTime(destiny / 5)
			PushFakerPlayerEnterRoom(room.Players, &p)

			// 第一个进入的机器人
			robot3 := CreateRobot(room.RoomClass.RoomType)
			robot3.PlayerPosition = 3
			room.Players[robot3.PlayerInfo.PlayerId] = robot3

			DelaySomeTime(destiny / 4)
			logger.Debug("===================推送第二个机器人")
			PushFakerPlayerEnterRoom(room.Players, &p)
		}
	}()
}
