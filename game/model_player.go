package game

import "gopkg.in/olahol/melody.v1"

// 玩家
type Player struct {
	Session        *melody.Session // 玩家session
	PlayerInfo     *PlayerInfo     // 玩家基本信息
	PlayerPosition int32           // 玩家当前座位
	IsRobot        bool            // 是否机器人
	HandCards      []*Card         // 手牌
	ThrowCards     []*Card         // 打初的牌
	IsLandlord     bool            // 是否地主
	IsGameHosting  bool            // 是否游戏托管
	IsCanDo        bool            // 是否可操作
	IsReady        bool            // 打完一局之后是否准备
	ActionChan     chan PlayerActionChan
}

type PlayerInfo struct {
	PlayerId string
	Name     string
	HeadImg  string
	Gold     float64
}

// 玩家动作
type PlayerActionChan struct {
}
