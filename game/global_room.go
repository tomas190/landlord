package game

import (
	"github.com/wonderivan/logger"
	"sync"
)

var globalRooms struct {
	roomMaps map[string]*Room
	rwMutex  sync.RWMutex
}

// 判断用户是否存在
// globalRooms
// @param key   : roomId
func IsExistRoom(roomId string) bool {
	globalRooms.rwMutex.Lock()
	_, isExist := globalRooms.roomMaps[roomId]
	globalRooms.rwMutex.Unlock()
	return isExist
}

// 存放全局的房间
// globalRooms
// @param key   :roomId
// @param Room
func SaveRoom(roomId string, room *Room) {
	globalRooms.rwMutex.Lock()
	globalRooms.roomMaps[roomId] = room
	globalRooms.rwMutex.Unlock()
	logger.Debug("新建"+roomId+"房间后 的房间数量:", len(globalRooms.roomMaps))
}

// 获取Room
// globalRooms
// @param key   : roomId
func GetRoom(roomId string) *Room {
	globalRooms.rwMutex.Lock()
	defer globalRooms.rwMutex.Unlock()
	room := globalRooms.roomMaps[roomId]
	return room
}

// 删除Room
// globalRooms
// @param key   : roomId
func RemoveRoom(roomId string) {
	globalRooms.rwMutex.Lock()
	defer globalRooms.rwMutex.Unlock()
	_, ok := globalRooms.roomMaps[roomId]
	if ok {
		delete(globalRooms.roomMaps, roomId)
	}
	logger.Debug("清除"+roomId+"房间后 的房间数量:", len(globalRooms.roomMaps))
}
