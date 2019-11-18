package game

import (
	"landlord/mconst/direction"
	"landlord/mconst/roomType"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
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
	Player        PlayerInfo
	Session       *melody.Session
}

var ExpFieldWaitUser struct {
	WaitUsers []*WaitUser
	mu        sync.RWMutex
}
var LowFieldWaitUser struct {
	WaitUsers []*WaitUser
	mu        sync.RWMutex
}
var MidFieldWaitUser struct {
	WaitUsers []*WaitUser
	mu        sync.RWMutex
}
var HighFieldWaitUser struct {
	WaitUsers []*WaitUser
	mu        sync.RWMutex
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
	wu.Player = player
	wu.Session = session
	ExpFieldWaitUser.WaitUsers = append(ExpFieldWaitUser.WaitUsers, &wu)

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
	// todo  这里死锁  后面修改
	//ExpFieldWaitUser.mu.Lock()
	//defer ExpFieldWaitUser.mu.Unlock()
	if len(ExpFieldWaitUser.WaitUsers) <= 0 {
		AddExpFieldWaitUser(session, playerInfo)
	} else {
		wp := ExpFieldWaitUser.WaitUsers[0]
		wPlayerInfo := wp.Player

		// 置空等待队列
		EmptyExpFieldWaitUser(session, wPlayerInfo)
		//wSession := wp.Session
		var wPlayer Player
		wPlayer.PlayerInfo = &wPlayerInfo
		wPlayer.Session = wp.Session
		wPlayer.Direction = direction.East
		//wPlayer.ActionChan = make(chan PlayerActionChan)

		var cPlayer Player
		cPlayer.PlayerInfo = &playerInfo
		cPlayer.Session = session
		cPlayer.Direction = direction.West
		//cPlayer.ActionChan = make(chan PlayerActionChan)

		var players []*Player
		players = append(players, &wPlayer, &cPlayer)
		room := NewRoom(roomType.ExperienceField, players)

		// 设置用户全局房间Id
		SetSessionRoomId(wPlayer.Session, room.RoomId)
		SetSessionRoomId(cPlayer.Session, room.RoomId)
		// 保存房间
		SaveRoom(room.RoomId, room)

		// 开启线程 游戏开始
		//go PlayGame(room)
	}

}
