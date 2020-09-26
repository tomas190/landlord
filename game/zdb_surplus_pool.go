package game

import (
	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2/bson"
	"landlord/mconst/sysSet"
)

/*
Lob, [07.03.20 10:51]
[Forwarded from cos]
盈余池 = （该游戏全部实际的玩家历史总输 - （该游戏全部实际的玩家历史总赢 * 100%）- （该游戏的历史实际的玩家总数 * 0））* 50%，当盈余池小于0的时候，玩家70%的机率为输
盈余池表名：surplus-pool
1玩家历史总输字段：player-total-lose
2玩家历史总赢字段：player-total-win
4历史实际的玩家总数字段：total-player
3历史总赢乘的百分比字段（100%那个值）：percentage-to-total-win
5玩家总数所剩的系数（0那个值）：coefficient-to-total-player
6最后百分比（50%那个值）：final-percentage
7盈余池 表中 增加 “玩家输赢差额：player-total-lose-win”  字段
8盈余池：surplus-pool
9盈余池后的玩家输百分比（70%那个值）：player-lose-rate-after-surplus-pool
*/
const SurplusPoolOneName = "surplus_pool"

// 此记录永远只有一条
type SurplusPoolOne struct {
	PlayerTotalLose                float64 `json:"player_total_lose" bson:"player_total_lose"`
	PlayerTotalWin                 float64 `json:"player_total_win" bson:"player_total_win"`
	TotalPlayer                    int     `json:"total_player" bson:"total_player"`
	PercentageToTotalWin           float64 `json:"percentage_to_total_win" bson:"percentage_to_total_win"`
	CoefficientToTotalPlayer       float64 `json:"coefficient_to_total_player" bson:"coefficient_to_total_player"`
	FinalPercentage                float64 `json:"final_percentage" bson:"final_percentage"`
	PlayerLoseRateAfterSurplusPool float64 `json:"player_lose_rate_after_surplus_pool" bson:"player_lose_rate_after_surplus_pool"`
	SurplusPool                    float64 `json:"surplus_pool" bson:"surplus_pool"`
	PlayerTotalLoseWin             float64 `json:"player_total_lose_win" bson:"player_total_lose_win"`
}

// 初始化盈余池数据
// 取最新盈余池 更新数据
func UptSurplusPoolOne() {
	var surplus SurplusPool
	sp := surplus.GetLastSurplus()
	var spo SurplusPoolOne
	spo.PlayerTotalLose = sp.PlayerAllLoss
	spo.PlayerTotalWin = sp.PlayerAllWin
	spo.TotalPlayer = sp.PlayerCount
	spo.PercentageToTotalWin = sysSet.PERCENTAGE_TO_TOTAL_WIN
	spo.CoefficientToTotalPlayer = sysSet.COEFFICIENT_TO_TOTAL_PLAYER
	spo.FinalPercentage = sysSet.FINAL_PERCENTAGE
	spo.PlayerLoseRateAfterSurplusPool = sysSet.PLAYER_LOSE_RATE_AFTER_SURPLUS_POOL

	if sp.CurrentSurplus == sp.PlayerAllLoss-sp.PlayerAllWin {
		var p PlayerRecode
		playersCount := p.CountPlayers()
		spo.SurplusPool = (sp.PlayerAllLoss -
			sp.PlayerAllWin*sysSet.PERCENTAGE_TO_TOTAL_WIN -
			float64(playersCount)*sysSet.COEFFICIENT_TO_TOTAL_PLAYER) *
			sysSet.FINAL_PERCENTAGE
	}
	spo.SurplusPool = sp.CurrentSurplus
	spo.PlayerTotalLoseWin = sp.PlayerAllLoss - sp.PlayerAllWin

	spo.EmptyData()
	spo.insertSurplusPoolOne()
}

func (s *SurplusPoolOne) insertSurplusPoolOne() {
	session, c := GetDBConn(Server.MongoDBName, SurplusPoolOneName)
	defer session.Close()
	err := c.Insert(s)
	if err != nil {
		logger.Error("更新盈余池失败:", err.Error())
	}

	//logger.Debug("========== 单个盈余池 =========")
	//fmt.Printf("%+v\n",*s)

}

// 清空数据
func (s *SurplusPoolOne) EmptyData() {
	session, c := GetDBConn(Server.MongoDBName, SurplusPoolOneName)
	defer session.Close()

	_, err := c.RemoveAll(bson.M{})
	if err != nil {
		logger.Debug("清空数据失败err:", err.Error())
	}
}

func (s *SurplusPoolOne) GetLastSurplusOne() (*SurplusPoolOne, error) {
	session, c := GetDBConn(Server.MongoDBName, SurplusPoolOneName)
	defer session.Close()

	var surplus SurplusPoolOne
	err := c.Find(nil).Sort("-_id").One(&surplus)
	if err != nil {
		//SendLogToCenter("ERR", "game/zdb_surplus_new.go", "62", "获取盈余池失败:"+err.Error())
		logger.Error("获取盈余池失败:", err.Error())
		return &surplus, err
	}
	return &surplus, nil
}

func UptSurplusConf(percentageToTotalWin,
	playerLoseRateAfterSurplusPool,
	coefficientToTotalPlayer,
	finalPercentage float64) {
	if percentageToTotalWin != -1 {
		sysSet.PERCENTAGE_TO_TOTAL_WIN = percentageToTotalWin
	}
	if playerLoseRateAfterSurplusPool != -1 {
		sysSet.PLAYER_LOSE_RATE_AFTER_SURPLUS_POOL = playerLoseRateAfterSurplusPool
	}

	if coefficientToTotalPlayer != -1 {
		sysSet.COEFFICIENT_TO_TOTAL_PLAYER = coefficientToTotalPlayer
	}
	if finalPercentage != -1 {
		sysSet.FINAL_PERCENTAGE = finalPercentage
	}

	UptSurplusPoolOne()
}
