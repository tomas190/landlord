package game

import (
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
	"landlord/mconst/roomType"
	"sync"
	"time"
)

/*
const (
	ExperienceField = 1 // 体验场
	LowField        = 2 // 初级场
	MidField        = 3 // 中级场
	HighField       = 4 // 高级场
)
*/

type WaitUser struct {
	EnterRoomTime int64
	Player        *PlayerInfo
	Session       *melody.Session
}

var ExpFieldWaitUser struct {
	WaitUsers map[string]*WaitUser
	mu        sync.RWMutex
}
var LowFieldWaitUser struct {
	WaitUsers map[string]*WaitUser // key userId
	mu        sync.RWMutex
}
var MidFieldWaitUser struct {
	WaitUsers map[string]*WaitUser
	mu        sync.RWMutex
}
var HighFieldWaitUser struct {
	WaitUsers map[string]*WaitUser
	mu        sync.RWMutex
}

func init() {
	ExpFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
	LowFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
	MidFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
	HighFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)

	globalRooms.roomMaps = make(map[string]*Room)
}

// 体验场 等待队列是否为空
func IsExpFieldWaitUserNil() bool {
	ExpFieldWaitUser.mu.Lock()
	defer ExpFieldWaitUser.mu.Unlock()
	if len(ExpFieldWaitUser.WaitUsers) == 0 {
		return true
	}
	return false
}

// 体验场 用户进入等待队列
func AddExpFieldWaitUser(session *melody.Session, player PlayerInfo) {
	ExpFieldWaitUser.mu.Lock()
	defer ExpFieldWaitUser.mu.Unlock()

	var wu WaitUser
	wu.EnterRoomTime = time.Now().Unix()
	wu.Player = &player
	wu.Session = session
	ExpFieldWaitUser.WaitUsers[player.PlayerId] = &wu

	// todo  需要3到5秒后如果没有玩家与之进行匹配 则分配一个机器人

	logger.Debug("玩家 ", player.PlayerId, " 进入等待队列")
}

// 体验场 置空等待队列
func EmptyExpFieldWaitUser(session *melody.Session, player PlayerInfo) {
	ExpFieldWaitUser.mu.Lock()
	defer ExpFieldWaitUser.mu.Unlock()

	ExpFieldWaitUser.WaitUsers = nil

	logger.Debug("玩家 ", player.PlayerId, " 已经开始游戏")
	logger.Debug("等待队列置空...")

}

// 体验场 处理玩家进入体验场
func DealPlayerEnterExpField(session *melody.Session, playerInfo PlayerInfo) {
	if len(ExpFieldWaitUser.WaitUsers) < 2 { // 斗地主需要三个人才能玩
		AddExpFieldWaitUser(session, playerInfo)
	} else if len(ExpFieldWaitUser.WaitUsers) == 2 {
		ExpFieldWaitUser.mu.Lock()
		players := make(map[string]*Player)
		room := NewRoom(roomType.ExperienceField, nil)
		var position int32
		for id, wPlayer := range ExpFieldWaitUser.WaitUsers {
			// todo 这一这里是否
			position++
			var p Player
			p.PlayerInfo = wPlayer.Player
			p.PlayerPosition = position
			p.Session = wPlayer.Session
			p.IsReady = true
			p.ActionChan = make(chan PlayerActionChan)
			players[id] = &p

			//设置用户全局房间Id
			logger.Debug("wp:", wPlayer.Session)
			logger.Debug("rd:", room.RoomId)
			SetSessionRoomId(wPlayer.Session, room.RoomId)
		}
		var p Player
		p.PlayerInfo = &playerInfo
		p.PlayerPosition = 3
		p.Session = session
		p.IsReady = true
		p.ActionChan = make(chan PlayerActionChan)
		players[playerInfo.PlayerId] = &p
		room.Players = players

		// 等待用户分配房间后清空 等待队列
		ExpFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)

		ExpFieldWaitUser.mu.Unlock()

		// 设置用户全局房间Id
		SetSessionRoomId(session, room.RoomId)
		// 保存房间
		SaveRoom(room.RoomId, room)

		// 开启线程 游戏开始
		go PlayGame(room)
	} else {
		logger.Debug("这不可能！ 如果出现这个消息 代码需要调整")
	}

	logger.Debug("等待的人儿:",ExpFieldWaitUser.WaitUsers)

}
