package game

import (
	"gopkg.in/olahol/melody.v1"
	"time"
)

func DealPlayerEnterRoomWithRobot(session *melody.Session, playerInfo PlayerInfo, roomType int32) {

	num := RandNum(3, 15)
	DelaySomeTime(time.Duration(num))

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
	go PlayGame(room)

}
