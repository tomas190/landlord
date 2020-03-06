package game

import (
	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2/bson"
	"landlord/mconst/sysSet"
	"time"
)

const SurplusPoolName = "surplus_pool"

/*
盈余池 = （该游戏全部实际的玩家历史总输 - （该游戏全部实际的玩家历史总赢 * 100%）- （该游戏的历史实际的玩家总数 * 0））* 50%，当盈余池小于0的时候，玩家70%的机率为输
盈余池表名：surplus-pool  1
玩家历史总输字段：player-total-lose  2
玩家历史总赢字段：player-total-win  3
历史实际的玩家总数字段：total-player 4

历史总赢乘的百分比字段（100%那个值）：percentage-to-total-win  5
玩家总数所剩的系数（0那个值）：coefficient-to-total-player 6
最后百分比（50%那个值）：final-percentage 7
盈余池后的玩家输百分比（70%那个值）：player-lose-rate-after-surplus-pool 8

*/

// 盈余池
type SurplusPool struct {
	GameName string       `json:"game_name" bson:"game_name"` // 游戏名称
	RoomType RoomClassify `json:"room_type" bson:"room_type"` // 房间类型

	PlayerAllLoss float64 `json:"player_total_lose" bson:"player_total_lose"` // 所有玩家输的总和
	PlayerAllWin  float64 `json:"player_total_win" bson:"player_total_win"`   // 所有玩家赢的总和 (计算的税后)
	PlayerCount   int     `json:"total_player" bson:"total_player"`           // 当前玩家人数计算

	PercentageToTotalWin           float64 `json:"percentage_to_total_win" bson:"percentage_to_total_win"`
	CoefficientToTotalPlayer       float64 `json:"coefficient_to_total_player" bson:"coefficient_to_total_player"`
	FinalPercentage                float64 `json:"final_percentage" bson:"final_percentage	"`
	PlayerLoseRateAfterSurplusPool float64 `json:"player_lose_rate_after_surplus_pool"`

	CurrentPlayerWin  float64 `json:"current_player_win" bson:"current_player_win"`   // 当局玩家总赢
	CurrentPlayerLoss float64 `json:"current_player_loss" bson:"current_player_loss"` // 当局玩家总输
	CurrentSurplus    float64 `json:"current_surplus" bson:"current_surplus"`         // 最新盈余

	RecodeTime    int64  `json:"recode_time" bson:"recode_time"`         // 结算时间
	RecodeTimeFmt string `json:"recode_time_fmt" bson:"recode_time_fmt"` // 结算时间字符串格式化
}

// 插入最新盈余
func (s *SurplusPool) InsertSurplus() {
	session, c := GetDBConn(Server.MongoDBName, SurplusPoolName)
	defer session.Close()

	lastSurplus := s.GetLastSurplus()

	lastSurplus.RoomType = s.RoomType
	lastSurplus.CurrentPlayerLoss = s.CurrentPlayerLoss
	lastSurplus.CurrentPlayerWin = s.CurrentPlayerWin

	lastSurplus.PlayerAllLoss += s.CurrentPlayerLoss
	lastSurplus.PlayerAllWin += s.CurrentPlayerWin

	var p PlayerRecode
	playersCount := p.CountPlayers()

	/*
		盈余池 = (玩家总输 - (玩家总赢 * 103%) - (玩家数量 * 6)) * 50%
		改成
		盈余池 =((玩家总输 - (玩家总赢 * 100%) - (玩家数量 * 0)) * 50%
	*/
	//surplus := lastSurplus.PlayerAllLoss - lastSurplus.PlayerAllWin*1.03 - float64(playersCount)*6
	surplus := (lastSurplus.PlayerAllLoss -
		lastSurplus.PlayerAllWin*sysSet.PERCENTAGE_TO_TOTAL_WIN -
		float64(playersCount)*sysSet.COEFFICIENT_TO_TOTAL_PLAYER) *
		sysSet.FINAL_PERCENTAGE
	lastSurplus.CurrentSurplus = surplus
	lastSurplus.PlayerLoseRateAfterSurplusPool = sysSet.PLAYER_LOSE_RATE_AFTER_SURPLUS_POOL
	lastSurplus.FinalPercentage = sysSet.FINAL_PERCENTAGE
	lastSurplus.PlayerCount = playersCount

	now := time.Now()
	lastSurplus.RecodeTime = now.Unix()
	lastSurplus.RecodeTimeFmt = now.Format("2006-01-02 15:04:05")

	err := c.Insert(lastSurplus)
	if err != nil {
		logger.Debug("记录盈余池失败:", err.Error())
	}
}

// 当有新玩家插入最新盈余
func (s *SurplusPool) InsertSurplusNewUser() {
	session, c := GetDBConn(Server.MongoDBName, SurplusPoolName)
	defer session.Close()

	lastSurplus := s.GetLastSurplus()
	lastSurplus.CurrentPlayerWin = 0
	lastSurplus.CurrentPlayerLoss = 0

	var p PlayerRecode
	playersCount := p.CountPlayers()

	/*
		盈余池 = (玩家总输 - (玩家总赢 * 103%) - (玩家数量 * 6)) * 50%
		改成
		盈余池 =((玩家总输 - (玩家总赢 * 100%) - (玩家数量 * 0)) * 50%
	*/
	//surplus := lastSurplus.PlayerAllLoss - lastSurplus.PlayerAllWin*1.03 - float64(playersCount)*6
	surplus := (lastSurplus.PlayerAllLoss -
		lastSurplus.PlayerAllWin*sysSet.PERCENTAGE_TO_TOTAL_WIN -
		float64(playersCount)*sysSet.COEFFICIENT_TO_TOTAL_PLAYER) *
		sysSet.FINAL_PERCENTAGE
	lastSurplus.CurrentSurplus = surplus

	lastSurplus.PlayerLoseRateAfterSurplusPool = sysSet.PLAYER_LOSE_RATE_AFTER_SURPLUS_POOL
	lastSurplus.FinalPercentage = sysSet.FINAL_PERCENTAGE
	lastSurplus.PercentageToTotalWin = sysSet.PERCENTAGE_TO_TOTAL_WIN
	lastSurplus.CoefficientToTotalPlayer = sysSet.COEFFICIENT_TO_TOTAL_PLAYER

	lastSurplus.CurrentSurplus = surplus
	lastSurplus.PlayerCount = playersCount

	now := time.Now()
	lastSurplus.RecodeTime = now.Unix()
	lastSurplus.RecodeTimeFmt = now.Format("2006-01-02 15:04:05")

	err := c.Insert(lastSurplus)
	if err != nil {
		logger.Debug("记录盈余池失败:", err.Error())
	}
}

// 获取最新盈余
func (s *SurplusPool) GetLastSurplus() *SurplusPool {
	session, c := GetDBConn(Server.MongoDBName, SurplusPoolName)
	defer session.Close()

	var surplus SurplusPool
	err := c.Find(nil).Sort("-recode_time").One(&surplus)
	if err != nil {
		logger.Error("获取盈余池失败:", err.Error())
		return &surplus
	}
	return &surplus
}

// 初始化一条空的数据
func InitSurplusPool() {
	session, c := GetDBConn(Server.MongoDBName, SurplusPoolName)
	defer session.Close()

	isExist := IsExistSurplusPool()
	if isExist {
		return
	} else {
		now := time.Now()
		var initItem SurplusPool
		initItem.GameName = sysSet.GameName
		initItem.PlayerAllWin = 0
		initItem.PlayerAllLoss = 0
		initItem.CurrentSurplus = 0
		initItem.CurrentPlayerLoss = 0
		initItem.CurrentPlayerWin = 0
		initItem.RecodeTime = now.Unix()
		initItem.RecodeTimeFmt = now.Format("2006-01-02 15:04:05")
		err := c.Insert(initItem)
		if err != nil {
			logger.Error("初始化盈余池数据失败:", err.Error())
			panic("初始化盈余池数据失败:" + err.Error())
		}
	}

}

//// 判断名称是否存在
func IsExistSurplusPool() bool {
	session, c := GetDBConn(Server.MongoDBName, SurplusPoolName)
	defer session.Close()
	var item PlayerRecode
	err := c.Find(bson.M{"game_name": sysSet.GameName}).One(&item)
	if err != nil {
		return false
	}
	return true
}
