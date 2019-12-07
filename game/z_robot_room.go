package game

import (
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
	"time"
)

// 处理玩家进入房间 和机器人玩
func DealPlayerEnterRoomWithRobot(session *melody.Session, playerInfo PlayerInfo, roomType int32, waitChan chan struct{}) {

	delayTime := RandNum(10, 15)

	select {
	case <-waitChan:
		logger.Debug("============= 玩家不想玩了 =============")
		close(waitChan)
		// do nothing
	case <-time.After(time.Second * time.Duration(delayTime)):
		logger.Debug("============= ... =============")
		PlayWithRobot(session, playerInfo, roomType)
	}

}

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
