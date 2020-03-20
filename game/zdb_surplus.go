package game

import (
	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2/bson"
	"landlord/mconst/sysSet"
	"sync"
	"time"
)

var dbMu sync.RWMutex

const SurplusPoolName = "landlord.surplus_pool"

// 盈余池
type SurplusPool struct {
	GameName          string       `json:"game_name" bson:"game_name"`                     // 游戏名称
	RecodeTime        int64        `json:"recode_time" bson:"recode_time"`                 // 结算时间
	RecodeTimeFmt     string       `json:"recode_time_fmt" bson:"recode_time_fmt"`         // 结算时间字符串格式化
	RoomType          RoomClassify `json:"room_type" bson:"room_type"`                     // 房间类型
	CurrentPlayerWin  float64      `json:"current_player_win" bson:"current_player_win"`   // 当局玩家总赢
	CurrentPlayerLoss float64      `json:"current_player_loss" bson:"current_player_loss"` // 当局玩家总输
	CurrentSurplus    float64      `json:"current_surplus" bson:"current_surplus"`         // 最新盈余
	PlayerAllLoss     float64      `json:"player_all_loss" bson:"player_all_loss"`         // 所有玩家输的总和
	PlayerAllWin      float64      `json:"player_all_win" bson:"player_all_win"`           // 所有玩家赢的总和 (计算的税后)
	PlayerCount       int          `json:"player_count" bson:"player_count"`               // 当前玩家人数计算
}

// 插入最新盈余
func (s *SurplusPool) InsertSurplus() {

	session, c := GetDBConn(Server.MongoDBName, SurplusPoolName)
	defer session.Close()
	dbMu.Lock()
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
	lastSurplus.PlayerCount = playersCount

	now := time.Now()
	lastSurplus.RecodeTime = now.Unix()
	lastSurplus.RecodeTimeFmt = now.Format("2006-01-02 15:04:05")

	err := c.Insert(lastSurplus)
	if err != nil {
		logger.Debug("记录盈余池失败:", err.Error())
	}
	dbMu.Unlock()
	// 同步更新
	UptSurplusPoolOne()
}

// 当有新玩家插入最新盈余
func (s *SurplusPool) InsertSurplusNewUser() {
	//dbMu.Lock()
	//defer dbMu.Unlock()
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
	lastSurplus.CurrentSurplus = surplus
	lastSurplus.PlayerCount = playersCount

	now := time.Now()
	lastSurplus.RecodeTime = now.Unix()
	lastSurplus.RecodeTimeFmt = now.Format("2006-01-02 15:04:05")

	err := c.Insert(lastSurplus)
	if err != nil {
		logger.Debug("记录盈余池失败:", err.Error())
	}
	// 同步更新
	UptSurplusPoolOne()
}

// 获取最新盈余
func (s *SurplusPool) GetLastSurplus() *SurplusPool {
	//dbMu.Lock()
	//defer dbMu.Unlock()
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
	var item SurplusPool
	err := c.Find(bson.M{"game_name": sysSet.GameName}).One(&item)
	if err != nil {
		return false
	}
	return true
}
