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
	DidAction      int32           // 抢地主的操作 0 没有进行过抢地主 1 / -1 叫了地主  2/-2 抢了地主
	IsLandlord     bool            // 是否地主
	IsGameHosting  bool            // 是否游戏托管
	IsCanDo        bool            // 是该当前玩家操作
	IsMustDo       bool            // 是否必须出牌
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
	ActionType  int32   // 1叫地主 2抢地主 3打牌  -1 不叫 -2 步枪 -3 不要
	ActionCards []*Card // 操作牌
	CardsType   int32   // 牌的类型
}
