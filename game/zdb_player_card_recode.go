package game

import (
	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2/bson"
)

const playCardRecodeName = "landlord.player_card_recode"

// 牌局记录
type PlayCardRecode struct {
	StartTime    int64 // 确认谁是地主之后开始
	StartTimeFmt string
	Players      map[string]RoomPlayer
	RoundId      string
	Settlement   []SettlementInfo
	EndTime      int64
	EndTimeFmt   string
}

type RoomPlayer struct {
	playerInfo PlayerInfo
	isRobot    bool
	isLandlord bool
	Cards      []*Card // 原始手牌 // 确认地主之后开始
}

type SettlementInfo struct {
	Order         string
	PlayerId      string  //玩家Id
	PlayerName    string  //玩家名
	IsLandlord    int32   //是否地主
	Multiple      int32   //倍数信息
	WinLossGold   string  // 结算金币
	CurrentGold   float64 // 结算之后的当前金币
	WinOrFail     int32   // 赢还是输  1 :赢  -1:输
	MinSettlement bool    // 是否最小结算
	RemainCards   []*Card // 剩余牌
}

func (p *PlayCardRecode) addPlayCardRecode() {
	session, c := GetDBConn(Server.MongoDBName, playCardRecodeName)
	defer session.Close()

	if p.IsExistRoundId(){
		logger.Error("牌局异常重复roundId:",p.RoundId)
	}

	err := c.Insert(p)
	if err != nil {
		logger.Error("插入牌局记录失败!")
	}
}

// 判断名称是否存在
func (p *PlayCardRecode) IsExistRoundId() bool {
	session, c := GetDBConn(Server.MongoDBName, playCardRecodeName)
	defer session.Close()
	var item PlayCardRecode
	err := c.Find(bson.M{"round_id": p.RoundId}).One(&item)
	if err != nil {
		return false
	}
	return true
}
