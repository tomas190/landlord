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
	LastAction     int32           // 上一首的操作 0 没有进行过抢地主 1 / -1 叫了地主  2/-2 抢了地主  3/-3 出牌/没出
	LastOutCard    []*Card         // 上一首出的牌 过
	IsLandlord     bool            // 是否地主
	IsGameHosting  bool            // 是否游戏托管
	IsCanDo        bool            // 是该当前玩家操作
	IsMustDo       bool            // 是否必须出牌
	IsCloseSession bool            // session 是否断开
	IsExitRoom     bool            // 是否退出房间
	WaitingTime    int32           // 等待时间 // 当改玩家操作的时候 不托管的情况下 从系统设置时间开始
	// 每秒-1 来防止用户在叫地主阶段退出 用于恢复时间点
	///ActionChan chan PlayerActionChan
	GroupCard  GroupCard // 组牌信息 只有机器人才有
	HandsValue float32   // 手牌评分
}

type PlayerInfo struct {
	PlayerId    string
	Name        string
	HeadImg     string
	Gold        float64
	PlayerPkgId int
	IsOnClear   bool
}

// 玩家动作
type PlayerActionChan struct {
	ActionType  int32   // 1叫地主 2抢地主 3打牌  -1 不叫 -2 步枪 -3 不要
	ActionCards []*Card // 操作牌
	CardsType   int32   // 牌的类型
}
