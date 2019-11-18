package game

import "gopkg.in/olahol/melody.v1"

// 玩家
type Player struct {
	Session        *melody.Session // 玩家session
	PlayerInfo     *PlayerInfo
	IsRobot        bool    // 是否机器人
	HandMahjongs   []int32 // 手牌
	FlowerMahjongs []int32 // 花牌
	PengMahjongs   []int32 // 碰牌
	ChiMahjongs    []int32 // 吃牌 [11,13,15,] 则代表吃牌 11,12,13 ; 13,14,15 ; 15,16,17 {记录吃牌最小的数字}
	GangMahjongs   []int32 // 杠牌
	AnGangMahjongs []int32 // 暗杠牌
	ThrowMahjongs  []int32 // 弃牌堆
	Direction      string  // 风向
	IsVillage      bool    // 是否是庄
	IsListen       bool    // 是否报听
	IsGameHosting  bool    // 是否游戏托管
	IsAction       bool    // 是否可操作 如打牌 杠牌 吃牌...
	IsReady        bool    // 打完一局之后是否准备
	//ActionChan     chan PlayerActionChan
}

type PlayerInfo struct {
	PlayerId string
	Name     string
	HeadImg  string
	Gold     float64
}
