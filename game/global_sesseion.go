package game

import (
	"errors"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
)

/*type SessionPlayerInfo struct {
	PlayerId      string
	Name    string
	HeadImg string
	Gold    float64
	AuthKey string
}
	RoomId string
	IsLogin bool
*/

func SetSessionPlayerInfo(session *melody.Session, playerInfo *PlayerInfo) {
	session.Set("playerInfo", playerInfo)

}

func GetSessionPlayerInfo(session *melody.Session) (*PlayerInfo, error) {
	var result *PlayerInfo
	value, exists := session.Get("playerInfo")
	if exists {
		return nil, errors.New("不存在绑定的playerInfo信息")
	}

	result = value.(*PlayerInfo)
	return result, nil
}

func SetSessionGold(session *melody.Session, gold float64) error {
	value, exists := session.Get("playerInfo")
	if !exists {
		return errors.New("SetSessionGold fail 无用户信息")
	}

	pi := value.(*PlayerInfo)

	if gold < 0 && pi.Gold < (-gold) {
		// todo 钱不够扣的情况
		logger.Debug("出现这种情况 钱不够扣")
	}

	pi.Gold = pi.Gold + gold
	return nil

}

func SetSessionIsLogin(session *melody.Session) {
	session.Set("isLogin", true)
}

// 该session是否已经成功登陆
func GetSessionIsLogin(session *melody.Session) bool {
	value, exists := session.Get("isLogin")
	if !exists {
		return false
	}
	isLogin := value.(bool)
	return isLogin
}

//
func SetSessionRoomId(session *melody.Session, roomId string) {
	session.Set("roomId", roomId)
}

//
func GetSessionRoomId(session *melody.Session) string {
	value, exists := session.Get("roomId")
	if !exists {
		return ""
	}
	return value.(string)
}
