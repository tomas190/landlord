package game

import (
	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2/bson"
	"landlord/msg/mproto"
	"time"
)

const playCardRecodeName = "landlord.player_card_recode"

// 牌局记录
type PlayCardRecode struct {
	StartTime int64 `json:"start_time" bson:"start_time"` // 确认谁是地主之后开始
	//StartTimeFmt string                `json:"start_time_fmt" bson:"start_time_fmt"`
	Players     map[string]RoomPlayer `json:"players" bson:"players"`
	PlayerIds   string                `json:"player_ids" bson:"player_ids"`
	RoundId     string                `json:"round_id" bson:"round_id"`
	RoomId      string                `json:"room_id" bson:"room_id"`
	BottomCard  []*Card               `json:"bottom_card" bson:"bottom_card"`
	Settlement  []SettlementInfo      `json:"settlement" bson:"settlement"`
	EndTime     int64                 `json:"end_time" bson:"end_time"`
	GameTaxRate float64               `json:"game_tax_rate"`
	//EndTimeFmt   string                `json:"end_time_fmt" bson:"end_time_fmt"`
}

type RoomPlayer struct {
	PlayerId   string  `json:"player_id" bson:"player_id"`
	IsRobot    bool    `json:"is_robot" bson:"is_robot"`
	IsLandlord bool    `json:"is_landlord" bson:"is_landlord"`
	Cards      []*Card `json:"cards" bson:"cards"` // 原始手牌 // 确认地主之后开始
}

type SettlementInfo struct {
	//Order         string
	PlayerId   string //玩家Id
	PlayerName string //玩家名
	IsLandlord int32  //是否地主
	//IsRobot       bool    // 是否机器人
	Multiple      int32   //倍数信息
	WinLossGold   string  // 结算金币
	CurrentGold   float64 // 结算之后的当前金币
	WinOrFail     int32   // 赢还是输  1 :赢  -1:输
	MinSettlement bool    // 是否最小结算
	RemainCards   []*Card // 剩余牌
}

func (p *PlayCardRecode) AddPlayCardRecode() {
	session, c := GetDBConn(Server.MongoDBName, playCardRecodeName)
	defer session.Close()

	if p.IsExistRoundId() {
		logger.Error("牌局异常重复roundId:", p.RoundId)
	}

	err := c.Insert(p)
	if err != nil {
		logger.Error("插入牌局记录失败!")
	}
}

func (p *PlayCardRecode) UptPlayCardRecode() {
	session, c := GetDBConn(Server.MongoDBName, playCardRecodeName)
	defer session.Close()
	err := c.Update(bson.M{"round_id": p.RoundId}, p)
	if err != nil {
		logger.Debug("更新记录失败:randId:", p.RoundId)
	}
}

func (p *PlayCardRecode) GetPlayCardRecodeByRoundId() (PlayCardRecode, error) {
	session, c := GetDBConn(Server.MongoDBName, playCardRecodeName)
	defer session.Close()
	var item PlayCardRecode
	err := c.Find(bson.M{"round_id": p.RoundId}).One(&item)
	if err != nil {
		return item, err
	}
	return item, nil
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

func DBRecode(room *Room) {
	var r PlayCardRecode
	r.RoundId = room.RoundId
	r.RoomId = room.RoomId
	r.Players = getPlayers(room)
	r.PlayerIds = getPlayerIds(room)
	r.BottomCard = room.BottomCards
	r.StartTime = time.Now().Unix()

	r.AddPlayCardRecode()
}

func DBUptRecode(room *Room, s mproto.PushSettlement) {
	var rd PlayCardRecode
	rd.RoundId = room.RoundId
	recode, err := rd.GetPlayCardRecodeByRoundId()
	if err != nil {
		logger.Error("更新结算记录失败:", err.Error())
		logger.Error("roundId:", room.RoomId)
		return
	}
	recode.EndTime = time.Now().Unix()
	recode.GameTaxRate = Server.GameTaxRate

	var sts []SettlementInfo

	for i := 0; i < len(s.Settlement); i++ {
		sti := s.Settlement[i]
		var st SettlementInfo
		st.IsLandlord = sti.IsLandlord
		st.RemainCards = append([]*Card{}, ChangeProtoToCard(sti.RemainCards)...)
		st.CurrentGold = sti.CurrentGold
		st.WinOrFail = sti.WinOrFail
		st.PlayerId = sti.PlayerId
		st.PlayerName = sti.PlayerName
		st.MinSettlement = sti.MinSettlement
		st.Multiple = sti.Multiple
		st.WinLossGold = sti.WinLossGold
		sts = append(sts, st)
	}
	recode.Settlement = sts
	recode.UptPlayCardRecode()

}

func (p *PlayCardRecode) GetPlayCardRecodeList(skip, limit int, selector bson.M, sortBy string) ([]PlayCardRecode, int, error) {
	session, c := GetDBConn(Server.MongoDBName, playCardRecodeName)
	defer session.Close()

	var wts []PlayCardRecode

	n, err := c.Find(selector).Count()
	if err != nil {
		return nil, 0, err
	}
	err = c.Find(selector).Sort(sortBy).Skip(skip).Limit(limit).All(&wts)
	if err != nil {
		return nil, 0, err
	}

	return wts, n, nil
}

func getPlayers(room *Room) map[string]RoomPlayer {
	rp := make(map[string]RoomPlayer, 3)
	for k, v := range room.Players {
		var r RoomPlayer
		r.PlayerId = v.PlayerInfo.PlayerId
		r.Cards = append([]*Card{}, v.HandCards...)
		r.IsLandlord = v.IsLandlord
		r.IsRobot = v.IsRobot
		rp[k] = r
	}
	return rp
}

func getPlayerIds(room *Room) string {
	var ids string
	for k, _ := range room.Players {
		ids += k
	}
	return ids
}
