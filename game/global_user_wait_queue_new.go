package game
//
//import (
//	"github.com/wonderivan/logger"
//	"gopkg.in/olahol/melody.v1"
//	"landlord/mconst/roomType"
//	"sync"
//	"time"
//)
//
//type WaitUser struct {
//	EnterRoomTime int64
//	Player        *PlayerInfo
//	Session       *melody.Session
//}
//
//var ExpFieldWaitUser struct {
//	WaitUsers map[string]*WaitUser
//	mu        sync.RWMutex
//}
//var LowFieldWaitUser struct {
//	WaitUsers map[string]*WaitUser // key userId
//	mu        sync.RWMutex
//}
//var MidFieldWaitUser struct {
//	WaitUsers map[string]*WaitUser
//	mu        sync.RWMutex
//}
//var HighFieldWaitUser struct {
//	WaitUsers map[string]*WaitUser
//	mu        sync.RWMutex
//}
//
//func init() {
//	ExpFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
//	LowFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
//	MidFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
//	HighFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
//	globalRooms.roomMaps = make(map[string]*Room, 10)
//	go matchRobotRoutine()
//}
//
//func matchRobotRoutine() {
//	for {
//		matchExpField()
//		matchLowField()
//		matchMidField()
//		matchHighField()
//	}
//
//}
//
//func RemoveWaitUser(playerId string) {
//	RemoveExpFieldWaitUser(playerId)
//	//RemoveLowFieldWaitUser(playerId)
//	//RemoveExpFieldWaitUser(playerId)
//	//RemoveHighFieldWaitUser(playerId)
//}
//
///*=================  体验场 ===============*/
//
//// 体验场 处理玩家进入体验场
//func DealPlayerEnterExpField(session *melody.Session, playerInfo PlayerInfo) {
//	isExist := IsPlayerInExpField(playerInfo.PlayerId)
//	if isExist {
//		logger.Debug(playerInfo.PlayerId, "已经在体验场等待队列:")
//		return
//	}
//	// 添加到体验场等待
//	AddExpFieldWaitUser(session, playerInfo)
//}
//
//// 体验场 用户进入等待队列
//func AddExpFieldWaitUser(session *melody.Session, player PlayerInfo) {
//	ExpFieldWaitUser.mu.Lock()
//	defer ExpFieldWaitUser.mu.Unlock()
//	var wu WaitUser
//	wu.EnterRoomTime = time.Now().Unix()
//	wu.Player = &player
//	wu.Session = session
//	ExpFieldWaitUser.WaitUsers[player.PlayerId] = &wu
//	logger.Debug("玩家 ", player.PlayerId, " 进入体验场等待队列")
//
//	if len(ExpFieldWaitUser.WaitUsers) == 3 {
//		// 匹配三个真人玩家
//		var position int32
//		room := NewRoom(roomType.ExperienceField, nil)
//		players := make(map[string]*Player)
//		for id, wPlayer := range ExpFieldWaitUser.WaitUsers {
//			position++
//			var p Player
//			p.PlayerInfo = wPlayer.Player
//			p.PlayerPosition = position
//			p.Session = wPlayer.Session
//			//p.IsReady = true
//			p.ActionChan = make(chan PlayerActionChan)
//			players[id] = &p
//			SetSessionRoomId(wPlayer.Session, room.RoomId)
//		}
//		room.Players = players
//		SaveRoom(room.RoomId, room)
//		EmptyExpFieldWaitUser()
//		go PlayGame(room)
//	}
//}
//
//// 体验场 用户退出等待队列
//func RemoveExpFieldWaitUser(playerId string) {
//	ExpFieldWaitUser.mu.Lock()
//	defer ExpFieldWaitUser.mu.Unlock()
//
//	_, ok := ExpFieldWaitUser.WaitUsers[playerId]
//	if ok {
//		delete(ExpFieldWaitUser.WaitUsers, playerId)
//		logger.Debug("玩家 ", playerId, " 退出等待队列")
//		return
//	}
//	logger.Debug("玩家 ", playerId, " 不在体验场等待队列中")
//
//}
//
//// 体验场 置空等待队列
//func EmptyExpFieldWaitUser() {
//	ExpFieldWaitUser.mu.Lock()
//	defer ExpFieldWaitUser.mu.Unlock()
//	ExpFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
//	logger.Debug("等待队列置空...")
//}
//
//// 用户是否在等待队列
//func IsPlayerInExpField(playerId string) bool {
//	ExpFieldWaitUser.mu.Lock()
//	defer ExpFieldWaitUser.mu.Unlock()
//
//	if _, ok := ExpFieldWaitUser.WaitUsers[playerId]; ok {
//		return true
//	}
//	return false
//}
//
////机器人 匹配队列
//func matchExpField() {
//	d := time.Duration(RandNum(2, 6))
//	DelaySomeTime(d)
//	ExpFieldWaitUser.mu.Lock()
//	defer ExpFieldWaitUser.mu.Unlock()
//	if len(ExpFieldWaitUser.WaitUsers) == 0 {
//		logger.Debug("当前体验场等待队列无用户")
//		return
//	}
//
//	if len(ExpFieldWaitUser.WaitUsers) > 0 {
//		for _, p := range ExpFieldWaitUser.WaitUsers {
//			PlayWithRobot(p.Session, *p.Player, roomType.ExperienceField)
//		}
//		logger.Debug("ExpField玩家匹配的机器人")
//	}
//	ExpFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
//}
//
///*=================  体验场 ===============*/
//
///*=================  低级场 ===============*/
//
//// 体验场 处理玩家进入体验场
//func DealPlayerEnterLowField(session *melody.Session, playerInfo PlayerInfo) {
//	isExist := IsPlayerInLowField(playerInfo.PlayerId)
//	if isExist {
//		logger.Debug(playerInfo.PlayerId, "已经在体验场等待队列:")
//		return
//	}
//	// 添加到体验场等待
//	AddLowFieldWaitUser(session, playerInfo)
//}
//
//// 体验场 用户进入等待队列
//func AddLowFieldWaitUser(session *melody.Session, player PlayerInfo) {
//	LowFieldWaitUser.mu.Lock()
//	defer LowFieldWaitUser.mu.Unlock()
//	var wu WaitUser
//	wu.EnterRoomTime = time.Now().Unix()
//	wu.Player = &player
//	wu.Session = session
//	LowFieldWaitUser.WaitUsers[player.PlayerId] = &wu
//	logger.Debug("玩家 ", player.PlayerId, " 进入体验场等待队列")
//
//	if len(LowFieldWaitUser.WaitUsers) == 3 {
//		// 匹配三个真人玩家
//		var position int32
//		room := NewRoom(roomType.ExperienceField, nil)
//		players := make(map[string]*Player)
//		for id, wPlayer := range LowFieldWaitUser.WaitUsers {
//			position++
//			var p Player
//			p.PlayerInfo = wPlayer.Player
//			p.PlayerPosition = position
//			p.Session = wPlayer.Session
//			//p.IsReady = true
//			p.ActionChan = make(chan PlayerActionChan)
//			players[id] = &p
//			SetSessionRoomId(wPlayer.Session, room.RoomId)
//		}
//		room.Players = players
//		SaveRoom(room.RoomId, room)
//		EmptyLowFieldWaitUser()
//		go PlayGame(room)
//	}
//}
//
//// 体验场 用户退出等待队列
//func RemoveLowFieldWaitUser(playerId string) {
//	LowFieldWaitUser.mu.Lock()
//	defer LowFieldWaitUser.mu.Unlock()
//
//	_, ok := LowFieldWaitUser.WaitUsers[playerId]
//	if ok {
//		delete(LowFieldWaitUser.WaitUsers, playerId)
//		logger.Debug("玩家 ", playerId, " 退出等待队列")
//		return
//	}
//	logger.Debug("玩家 ", playerId, " 不在体验场等待队列中")
//
//}
//
//// 体验场 置空等待队列
//func EmptyLowFieldWaitUser() {
//	LowFieldWaitUser.mu.Lock()
//	defer LowFieldWaitUser.mu.Unlock()
//	LowFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
//	logger.Debug("等待队列置空...")
//}
//
//// 用户是否在等待队列
//func IsPlayerInLowField(playerId string) bool {
//	LowFieldWaitUser.mu.Lock()
//	defer LowFieldWaitUser.mu.Unlock()
//
//	if _, ok := LowFieldWaitUser.WaitUsers[playerId]; ok {
//		return true
//	}
//	return false
//}
//
////机器人 匹配队列
//func matchLowField() {
//	d := time.Duration(RandNum(2, 6))
//	DelaySomeTime(d)
//	LowFieldWaitUser.mu.Lock()
//	defer LowFieldWaitUser.mu.Unlock()
//
//	if len(LowFieldWaitUser.WaitUsers) == 0 {
//		logger.Debug("当前低级场等待队列无用户")
//		return
//	}
//
//	if len(LowFieldWaitUser.WaitUsers) > 0 {
//		for _, p := range LowFieldWaitUser.WaitUsers {
//			PlayWithRobot(p.Session, *p.Player, roomType.LowField)
//		}
//		logger.Debug("LowField玩家匹配的机器人")
//	}
//	LowFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
//}
//
///*=================  低级场 ===============*/
//
///*=================  中级场 ===============*/
//
//// 体验场 处理玩家进入体验场
//func DealPlayerEnterMidField(session *melody.Session, playerInfo PlayerInfo) {
//	isExist := IsPlayerInMidField(playerInfo.PlayerId)
//	if isExist {
//		logger.Debug(playerInfo.PlayerId, "已经在体验场等待队列:")
//		return
//	}
//	// 添加到体验场等待
//	AddMidFieldWaitUser(session, playerInfo)
//}
//
//// 体验场 用户进入等待队列
//func AddMidFieldWaitUser(session *melody.Session, player PlayerInfo) {
//	MidFieldWaitUser.mu.Lock()
//	defer MidFieldWaitUser.mu.Unlock()
//	var wu WaitUser
//	wu.EnterRoomTime = time.Now().Unix()
//	wu.Player = &player
//	wu.Session = session
//	MidFieldWaitUser.WaitUsers[player.PlayerId] = &wu
//	logger.Debug("玩家 ", player.PlayerId, " 进入体验场等待队列")
//
//	if len(MidFieldWaitUser.WaitUsers) == 3 {
//		// 匹配三个真人玩家
//		var position int32
//		room := NewRoom(roomType.MidField, nil)
//		players := make(map[string]*Player)
//		for id, wPlayer := range MidFieldWaitUser.WaitUsers {
//			position++
//			var p Player
//			p.PlayerInfo = wPlayer.Player
//			p.PlayerPosition = position
//			p.Session = wPlayer.Session
//			//p.IsReady = true
//			p.ActionChan = make(chan PlayerActionChan)
//			players[id] = &p
//			SetSessionRoomId(wPlayer.Session, room.RoomId)
//		}
//		room.Players = players
//		SaveRoom(room.RoomId, room)
//		EmptyMidFieldWaitUser()
//		go PlayGame(room)
//	}
//}
//
//// 体验场 用户退出等待队列
//func RemoveMidFieldWaitUser(playerId string) {
//	MidFieldWaitUser.mu.Lock()
//	defer MidFieldWaitUser.mu.Unlock()
//
//	_, ok := MidFieldWaitUser.WaitUsers[playerId]
//	if ok {
//		delete(MidFieldWaitUser.WaitUsers, playerId)
//		logger.Debug("玩家 ", playerId, " 退出等待队列")
//		return
//	}
//	logger.Debug("玩家 ", playerId, " 不在体验场等待队列中")
//
//}
//
//// 体验场 置空等待队列
//func EmptyMidFieldWaitUser() {
//	MidFieldWaitUser.mu.Lock()
//	defer MidFieldWaitUser.mu.Unlock()
//	MidFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
//	logger.Debug("等待队列置空...")
//}
//
//// 用户是否在等待队列
//func IsPlayerInMidField(playerId string) bool {
//	MidFieldWaitUser.mu.Lock()
//	defer MidFieldWaitUser.mu.Unlock()
//
//	if _, ok := MidFieldWaitUser.WaitUsers[playerId]; ok {
//		return true
//	}
//	return false
//}
//
////机器人 匹配队列
//func matchMidField() {
//	d := time.Duration(RandNum(2, 6))
//	DelaySomeTime(d)
//	MidFieldWaitUser.mu.Lock()
//	defer MidFieldWaitUser.mu.Unlock()
//
//	if len(MidFieldWaitUser.WaitUsers) == 0 {
//		logger.Debug("当前中级场等待队列无用户")
//		return
//	}
//
//	if len(MidFieldWaitUser.WaitUsers) > 0 {
//		for _, p := range MidFieldWaitUser.WaitUsers {
//			PlayWithRobot(p.Session, *p.Player, roomType.MidField)
//		}
//		logger.Debug("MidField玩家匹配的机器人")
//	}
//	MidFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
//}
//
///*=================  中级场 ===============*/
//
///*=================  高级场 ===============*/
//
//// 体验场 处理玩家进入体验场
//func DealPlayerEnterHighField(session *melody.Session, playerInfo PlayerInfo) {
//	isExist := IsPlayerInHighField(playerInfo.PlayerId)
//	if isExist {
//		logger.Debug(playerInfo.PlayerId, "已经在体验场等待队列:")
//		return
//	}
//	// 添加到体验场等待
//	AddHighFieldWaitUser(session, playerInfo)
//}
//
//// 体验场 用户进入等待队列
//func AddHighFieldWaitUser(session *melody.Session, player PlayerInfo) {
//	HighFieldWaitUser.mu.Lock()
//	defer HighFieldWaitUser.mu.Unlock()
//	var wu WaitUser
//	wu.EnterRoomTime = time.Now().Unix()
//	wu.Player = &player
//	wu.Session = session
//	HighFieldWaitUser.WaitUsers[player.PlayerId] = &wu
//	logger.Debug("玩家 ", player.PlayerId, " 进入体验场等待队列")
//
//	if len(HighFieldWaitUser.WaitUsers) == 3 {
//		// 匹配三个真人玩家
//		var position int32
//		room := NewRoom(roomType.HighField, nil)
//		players := make(map[string]*Player)
//		for id, wPlayer := range HighFieldWaitUser.WaitUsers {
//			position++
//			var p Player
//			p.PlayerInfo = wPlayer.Player
//			p.PlayerPosition = position
//			p.Session = wPlayer.Session
//			//p.IsReady = true
//			p.ActionChan = make(chan PlayerActionChan)
//			players[id] = &p
//			SetSessionRoomId(wPlayer.Session, room.RoomId)
//		}
//		room.Players = players
//		SaveRoom(room.RoomId, room)
//		EmptyHighFieldWaitUser()
//		go PlayGame(room)
//	}
//}
//
//// 体验场 用户退出等待队列
//func RemoveHighFieldWaitUser(playerId string) {
//	HighFieldWaitUser.mu.Lock()
//	defer HighFieldWaitUser.mu.Unlock()
//
//	_, ok := HighFieldWaitUser.WaitUsers[playerId]
//	if ok {
//		delete(HighFieldWaitUser.WaitUsers, playerId)
//		logger.Debug("玩家 ", playerId, " 退出等待队列")
//		return
//	}
//	logger.Debug("玩家 ", playerId, " 不在体验场等待队列中")
//
//}
//
//// 体验场 置空等待队列
//func EmptyHighFieldWaitUser() {
//	HighFieldWaitUser.mu.Lock()
//	defer HighFieldWaitUser.mu.Unlock()
//	HighFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
//	logger.Debug("等待队列置空...")
//}
//
//// 用户是否在等待队列
//func IsPlayerInHighField(playerId string) bool {
//	HighFieldWaitUser.mu.Lock()
//	defer HighFieldWaitUser.mu.Unlock()
//
//	if _, ok := HighFieldWaitUser.WaitUsers[playerId]; ok {
//		return true
//	}
//	return false
//}
//
////机器人 匹配队列
//func matchHighField() {
//	d := time.Duration(RandNum(2, 6))
//	DelaySomeTime(d)
//	HighFieldWaitUser.mu.Lock()
//	defer HighFieldWaitUser.mu.Unlock()
//
//	if len(MidFieldWaitUser.WaitUsers) == 0 {
//		logger.Debug("当前高级场等待队列无用户")
//		return
//	}
//
//	if len(HighFieldWaitUser.WaitUsers) > 0 {
//		for _, p := range HighFieldWaitUser.WaitUsers {
//			PlayWithRobot(p.Session, *p.Player, roomType.HighField)
//		}
//		logger.Debug("HighField玩家匹配的机器人")
//	}
//	HighFieldWaitUser.WaitUsers = make(map[string]*WaitUser, 3)
//}
//
///*=================  高级场 ===============*/
