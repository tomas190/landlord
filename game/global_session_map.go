package game

import (
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
	"sync"
)

// 存放登录的用户 key userid
var globalLoginAgents struct {
	sessionMaps map[string]*melody.Session
	rwMutex     sync.RWMutex
}

func init() {
	globalLoginAgents.sessionMaps = make(map[string]*melody.Session, 100)
}

// 判断用户是否存在
// globalAgents
// @param key   : playerID
func IsExistAgent(key string) bool {
	globalLoginAgents.rwMutex.Lock()
	_, isExist := globalLoginAgents.sessionMaps[key]
	globalLoginAgents.rwMutex.Unlock()
	return isExist
}

// 存放全局的连接
// globalAgents
// @param key   :playerID
// @param agent
func SaveAgent(key string, session *melody.Session) {
	globalLoginAgents.rwMutex.Lock()
	globalLoginAgents.sessionMaps[key] = session
	//logger.Debug("新连接未登录:", len(globalAgents))
	globalLoginAgents.rwMutex.Unlock()

}

// 获取agent
// globalAgents
// @param key   : playerID
func GetAgent(key string) *melody.Session {
	globalLoginAgents.rwMutex.Lock()
	defer globalLoginAgents.rwMutex.Unlock()
	session := globalLoginAgents.sessionMaps[key]
	return session
}

// 删除agent
// globalAgents
// @param key   : playerID
func RemoveAgent(key string) {
	globalLoginAgents.rwMutex.Lock()
	defer globalLoginAgents.rwMutex.Unlock()
	_, ok := globalLoginAgents.sessionMaps[key]
	if ok {
		delete(globalLoginAgents.sessionMaps, key)
	}

	logger.Debug("清除"+key+"连接后:", len(globalLoginAgents.sessionMaps))
}

// 删除agent
// globalAgents
// @param key   :
func GetConnLen() int {
	globalLoginAgents.rwMutex.Lock()
	defer globalLoginAgents.rwMutex.Unlock()
	i := len(globalLoginAgents.sessionMaps)
	return i
}

func BackUserToHall() {
	globalLoginAgents.rwMutex.Lock()
	defer globalLoginAgents.rwMutex.Unlock()
	for playerId, s := range globalLoginAgents.sessionMaps {
		UserLogoutCenter(playerId,GetSessionPassword(s))
	}
}
