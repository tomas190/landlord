package game

import (
	"strconv"
	"sync"

	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
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

//func BackUserToHall() {
//	globalLoginAgents.rwMutex.Lock()
//	defer globalLoginAgents.rwMutex.Unlock()
//	for playerId, s := range globalLoginAgents.sessionMaps {
//		logger.Debug("退回玩家资金到大厅:",playerId)
//		UserLogoutCenter(playerId,GetSessionPassword(s))
//	}
//}

// 获取agent
// globalAgents
// @param userGold   : playerGold
func GetAgentByUserGold(userGold float64) (string, *melody.Session) {
	globalLoginAgents.rwMutex.Lock()
	defer globalLoginAgents.rwMutex.Unlock()
	var count int // 保证只有一位玩家  如果多个则返回nil
	var result *melody.Session
	var playerId string
	for pid, session := range globalLoginAgents.sessionMaps {
		value, exists := session.Get("playerInfo")
		if exists {
			p := value.(*PlayerInfo)
			if p.Gold == userGold {
				count++
				result = session
				playerId = pid
			}
		}
	}
	if count == 1 {
		return playerId, result
	}

	return "", nil
}

// 取得所有在線玩家
func GetAllOnlineUser() map[int][]int64 {

	globalLoginAgents.rwMutex.RLock()
	defer globalLoginAgents.rwMutex.RUnlock()

	list := make(map[int][]int64)
	for _, session := range globalLoginAgents.sessionMaps {
		p, err := GetSessionPlayerInfo(session)
		if err != nil {
			logger.Debug("无用户session信息:", err.Error())
		} else {
			userID, errUID := strconv.ParseInt(p.PlayerId, 10, 64)
			if errUID != nil {
			} else {
				// logger.Debug("GetAllOnlineUser userID= %v, PlayerPkgId = %v", userID, p.PlayerPkgId)
				list[p.PlayerPkgId] = append(list[p.PlayerPkgId], userID)
			}
		}
	}
	return list

}
