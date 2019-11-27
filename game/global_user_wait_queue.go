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

func RemoveWaitUser(playerId string) {
	RemoveExpFieldWaitUser(playerId)
	RemoveLowFieldWaitUser(playerId)
	RemoveMidFieldWaitUser(playerId)
	RemoveHighFieldWaitUser(playerId)
}

/*=================  体验场 ===============*/

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

	logger.Debug("玩家 ", player.PlayerId, " 进入体验场等待队列")
}

// 体验场 用户退出等待队列
func RemoveExpFieldWaitUser(playerId string) {
	ExpFieldWaitUser.mu.Lock()
	defer ExpFieldWaitUser.mu.Unlock()

	_, ok := ExpFieldWaitUser.WaitUsers[playerId]
	if ok {
		delete(ExpFieldWaitUser.WaitUsers, playerId)
		logger.Debug("玩家 ", playerId, " 退出等待队列")
		return
	}
	logger.Debug("玩家 ", playerId, " 不在体验场等待队列中")

}

// 体验场 置空等待队列
func EmptyExpFieldWaitUser() {
	ExpFieldWaitUser.mu.Lock()
	defer ExpFieldWaitUser.mu.Unlock()
	ExpFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
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
			position++
			var p Player
			p.PlayerInfo = wPlayer.Player
			p.PlayerPosition = position
			p.Session = wPlayer.Session
			//p.IsReady = true
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
		//p.IsReady = true
		p.ActionChan = make(chan PlayerActionChan)
		players[playerInfo.PlayerId] = &p
		room.Players = players

		// 等待用户分配房间后清空 等待队列
		ExpFieldWaitUser.mu.Unlock()
		EmptyExpFieldWaitUser()
		// 设置用户全局房间Id
		SetSessionRoomId(session, room.RoomId)
		// 保存房间
		SaveRoom(room.RoomId, room)

		// 开启线程 游戏开始
		go PlayGame(room)
	} else {
		logger.Debug("这不可能！ 如果出现这个消息 代码需要调整")
	}

	logger.Debug("体验场等待的用户:", ExpFieldWaitUser.WaitUsers)

}

/*=================  体验场 ===============*/

/*=================  低级场 ===============*/

// 低级场 用户进入等待队列
func AddLowFieldWaitUser(session *melody.Session, player PlayerInfo) {
	LowFieldWaitUser.mu.Lock()
	defer LowFieldWaitUser.mu.Unlock()

	var wu WaitUser
	wu.EnterRoomTime = time.Now().Unix()
	wu.Player = &player
	wu.Session = session
	LowFieldWaitUser.WaitUsers[player.PlayerId] = &wu

	// todo  需要3到5秒后如果没有玩家与之进行匹配 则分配一个机器人

	logger.Debug("玩家 ", player.PlayerId, " 进入低级场等待队列")
}

// 低级场 用户退出等待队列
func RemoveLowFieldWaitUser(playerId string) {
	LowFieldWaitUser.mu.Lock()
	defer LowFieldWaitUser.mu.Unlock()

	_, ok := LowFieldWaitUser.WaitUsers[playerId]
	if ok {
		delete(LowFieldWaitUser.WaitUsers, playerId)
		logger.Debug("玩家 ", playerId, " 退出等待队列")
		return
	}
	logger.Debug("玩家 ", playerId, " 不在低级场等待队列中")

}

// 低级场 置空等待队列
func EmptyLowFieldWaitUser() {
	LowFieldWaitUser.mu.Lock()
	defer LowFieldWaitUser.mu.Unlock()
	LowFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
	logger.Debug("低级场等待队列置空...")

}

// 低级场 处理玩家进入体验场
func DealPlayerEnterLowField(session *melody.Session, playerInfo PlayerInfo) {
	if len(LowFieldWaitUser.WaitUsers) < 2 { // 斗地主需要三个人才能玩
		AddLowFieldWaitUser(session, playerInfo)
	} else if len(LowFieldWaitUser.WaitUsers) == 2 {
		LowFieldWaitUser.mu.Lock()
		players := make(map[string]*Player)
		room := NewRoom(roomType.LowField, nil)
		var position int32
		for id, wPlayer := range LowFieldWaitUser.WaitUsers {
			position++
			var p Player
			p.PlayerInfo = wPlayer.Player
			p.PlayerPosition = position
			p.Session = wPlayer.Session
			//p.IsReady = true
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
		//p.IsReady = true
		p.ActionChan = make(chan PlayerActionChan)
		players[playerInfo.PlayerId] = &p
		room.Players = players

		// 等待用户分配房间后清空 等待队列
		LowFieldWaitUser.mu.Unlock()
		EmptyLowFieldWaitUser()
		// 设置用户全局房间Id
		SetSessionRoomId(session, room.RoomId)
		// 保存房间
		SaveRoom(room.RoomId, room)

		// 开启线程 游戏开始
		go PlayGame(room)
	} else {
		logger.Debug("这不可能！ 如果出现这个消息 代码需要调整")
	}

	logger.Debug("低级场等待的用户:", LowFieldWaitUser.WaitUsers)

}

/*=================  低级场 ===============*/

/*=================  中级场 ===============*/

// 中级场 用户进入等待队列
func AddMidFieldWaitUser(session *melody.Session, player PlayerInfo) {
	MidFieldWaitUser.mu.Lock()
	defer MidFieldWaitUser.mu.Unlock()

	var wu WaitUser
	wu.EnterRoomTime = time.Now().Unix()
	wu.Player = &player
	wu.Session = session
	MidFieldWaitUser.WaitUsers[player.PlayerId] = &wu

	// todo  需要3到5秒后如果没有玩家与之进行匹配 则分配一个机器人

	logger.Debug("玩家 ", player.PlayerId, " 进入中级场等待队列")
}

// 中级场 用户退出等待队列
func RemoveMidFieldWaitUser(playerId string) {
	MidFieldWaitUser.mu.Lock()
	defer MidFieldWaitUser.mu.Unlock()

	_, ok := MidFieldWaitUser.WaitUsers[playerId]
	if ok {
		delete(MidFieldWaitUser.WaitUsers, playerId)
		logger.Debug("玩家 ", playerId, " 退出等待队列")
		return
	}
	logger.Debug("玩家 ", playerId, " 不在中级场等待队列中")

}

// 中级场 置空等待队列
func EmptyMidFieldWaitUser() {
	MidFieldWaitUser.mu.Lock()
	defer MidFieldWaitUser.mu.Unlock()
	MidFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
	logger.Debug("等待队列置空...")

}

// 中级场 处理玩家进入中级场
func DealPlayerEnterMidField(session *melody.Session, playerInfo PlayerInfo) {
	if len(MidFieldWaitUser.WaitUsers) < 2 { // 斗地主需要三个人才能玩
		AddMidFieldWaitUser(session, playerInfo)
	} else if len(MidFieldWaitUser.WaitUsers) == 2 {
		MidFieldWaitUser.mu.Lock()
		players := make(map[string]*Player)
		room := NewRoom(roomType.MidField, nil)
		var position int32
		for id, wPlayer := range MidFieldWaitUser.WaitUsers {
			position++
			var p Player
			p.PlayerInfo = wPlayer.Player
			p.PlayerPosition = position
			p.Session = wPlayer.Session
			//p.IsReady = true
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
		//p.IsReady = true
		p.ActionChan = make(chan PlayerActionChan)
		players[playerInfo.PlayerId] = &p
		room.Players = players

		// 等待用户分配房间后清空 等待队列
		MidFieldWaitUser.mu.Unlock()
		EmptyMidFieldWaitUser()
		// 设置用户全局房间Id
		SetSessionRoomId(session, room.RoomId)
		// 保存房间
		SaveRoom(room.RoomId, room)

		// 开启线程 游戏开始
		go PlayGame(room)
	} else {
		logger.Debug("这不可能！ 如果出现这个消息 代码需要调整")
	}

	logger.Debug("中级场等待的用户:", MidFieldWaitUser.WaitUsers)

}

/*=================  中级场 ===============*/

/*=================  高级场 ===============*/

// 高级场 用户进入等待队列
func AddHighFieldWaitUser(session *melody.Session, player PlayerInfo) {
	HighFieldWaitUser.mu.Lock()
	defer HighFieldWaitUser.mu.Unlock()

	var wu WaitUser
	wu.EnterRoomTime = time.Now().Unix()
	wu.Player = &player
	wu.Session = session
	HighFieldWaitUser.WaitUsers[player.PlayerId] = &wu

	// todo  需要3到5秒后如果没有玩家与之进行匹配 则分配一个机器人

	logger.Debug("玩家 ", player.PlayerId, " 进入高级场等待队列")
}

// 高级场 用户退出等待队列
func RemoveHighFieldWaitUser(playerId string) {
	HighFieldWaitUser.mu.Lock()
	defer HighFieldWaitUser.mu.Unlock()

	_, ok := HighFieldWaitUser.WaitUsers[playerId]
	if ok {
		delete(HighFieldWaitUser.WaitUsers, playerId)
		logger.Debug("玩家 ", playerId, " 退出等待队列")
		return
	}
	logger.Debug("玩家 ", playerId, " 不在高级场等待队列中")

}

// 高级场 置空等待队列
func EmptyHighFieldWaitUser() {
	HighFieldWaitUser.mu.Lock()
	defer HighFieldWaitUser.mu.Unlock()
	HighFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
	logger.Debug("等待队列置空...")

}

// 高级场 处理玩家进入高级场
func DealPlayerEnterHighField(session *melody.Session, playerInfo PlayerInfo) {
	// 1.判断


	if len(HighFieldWaitUser.WaitUsers) < 2 { // 斗地主需要三个人才能玩
		AddHighFieldWaitUser(session, playerInfo)
	} else if len(HighFieldWaitUser.WaitUsers) == 2 {
		HighFieldWaitUser.mu.Lock()


		players := make(map[string]*Player)
		room := NewRoom(roomType.HighField, nil)
		var position int32
		for id, wPlayer := range HighFieldWaitUser.WaitUsers {
			position++
			var p Player
			p.PlayerInfo = wPlayer.Player
			p.PlayerPosition = position
			p.Session = wPlayer.Session
			//p.IsReady = true
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
		//p.IsReady = true
		p.ActionChan = make(chan PlayerActionChan)
		players[playerInfo.PlayerId] = &p
		room.Players = players

		// 等待用户分配房间后清空 等待队列
		HighFieldWaitUser.mu.Unlock()
		EmptyHighFieldWaitUser()
		// 设置用户全局房间Id
		SetSessionRoomId(session, room.RoomId)
		// 保存房间
		SaveRoom(room.RoomId, room)

		// 开启线程 游戏开始
		go PlayGame(room)
	} else {
		logger.Debug("这不可能！ 如果出现这个消息 代码需要调整")
	}

	logger.Debug("高级场等待的用户:", HighFieldWaitUser.WaitUsers)

}

/*=================  高级场 ===============*/
