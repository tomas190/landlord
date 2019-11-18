package game

import (
	"github.com/wonderivan/logger"
	"time"
)

func UserLogin(playerId, password string) {

	logger.Debug("<-------- UserLoginCenter -------->")
	baseData := &ToCenterMessage{}
	baseData.Event = msgUserLogin
	baseData.Data = &UserReq{
		ID:       playerId,
		//PassWord: password,
		GameId:  Server.GameId,
		Token:   password,
		DevName: Server.DevName,
		DevKey:  Server.DevKey}

	WriteMsgToCenter(baseData)
}

//UserLogoutCenter 用户登出
func UserLogoutCenter(userId string, password string) {
	logger.Debug("<-------- UserLogoutCenter  -------->")
	base := &ToCenterMessage{}
	base.Event = msgUserLogout
	base.Data = &UserReq{
		ID:       userId,
		PassWord: password,
		GameId:   Server.GameId,
		Token:    password,
		DevName:  Server.DevName,
		DevKey:   Server.DevKey,
	}
	// 发送消息到中心服
	WriteMsgToCenter(base)
}

//UserSyncWinScore 同步赢分
//func UserSyncWinScore(playerId string, winMoney float64, roundId, orderId string) {
//
//	logger.Debug("<-------- GenWinOrder -------->")
//	timeUnix := time.Now().Unix()
//
//	baseData := &ToCenterMessage{}
//	baseData.Event = msgUserWinScore
//	userWin := &UserChangeScore{}
//	userWin.Auth.Token = TokenOfCenter
//	userWin.Auth.DevKey = models.Server.DevKey
//	userWin.Info.CreateTime = timeUnix
//	userWin.Info.GameId = models.Server.GameId
//	userWin.Info.ID = playerId
//	//userWin.Info.LockMoney = 0
//	userWin.Info.Money = winMoney
//	userWin.Info.Order = orderId
//	userWin.Info.PayReason = "下注"
//	//userWin.Info.PreMoney = 0
//	userWin.Info.RoundId = roundId
//	baseData.Data = userWin
//
//	WriteMsgToCenter(baseData)
//}

//UserSyncWinScore 同步赢分
func UserSyncWinScore(playerId string, winMoney float64, roundId, orderId string) {

	logger.Debug("<-------- 发送赢钱指令 -------->")
	timeUnix := time.Now().Unix()

	baseData := &ToCenterMessage{}
	baseData.Event = msgUserWinScore
	userWin := &UserChangeScore{}
	// userWin.Auth.Token = TokenOfCenter
	userWin.Auth.DevName = Server.DevName
	userWin.Auth.DevKey = Server.DevKey
	userWin.Info.CreateTime = timeUnix
	userWin.Info.GameId = Server.GameId
	userWin.Info.ID = playerId
	//userWin.Info.LockMoney = 0
	userWin.Info.Money = winMoney
	userWin.Info.Order = orderId
	userWin.Info.PayReason = "下注"
	//userWin.Info.PreMoney = 0
	userWin.Info.RoundId = roundId
	baseData.Data = userWin

	logger.Debug("发送赢分指令:")
	PrintMsg("sendCenterMsg:",baseData)

	WriteMsgToCenter(baseData)
}

/*
{
    "event":"/GameServer/GameUser/winSettlement",
    "data":{
        "auth:{
            "token":"you token",
            "dev_key":"123"
        },
        "info":{
            "id":"123456",
            "create_time":1548971234,
            "pay_reason":" 下注",
            "money":12.0,
            "lock_money":120.0,
            "pre_money":12.0,
            "order":"自己创建一个唯一ID,方便之后查询",
            "game_id":"abc",
            "round_id":"唯一ID,用于识别多人是否在同一局游戏",
        }
    }
}
*/

//UserSyncWinScore 同步输分
func UserSyncLoseScore(playerId string, lossMoney float64, roundId, orderId string) {

	logger.Debug("<-------- GenLoseOrder -------->")

	timeUnix := time.Now().Unix()

	baseData := &ToCenterMessage{}
	baseData.Event = msgUserLoseScore
	userLose := &UserChangeScore{}
	// userLose.Auth.Token = TokenOfCenter
	userLose.Auth.DevName = Server.DevName
	userLose.Auth.DevKey = Server.DevKey
	userLose.Info.CreateTime = timeUnix
	userLose.Info.GameId = Server.GameId
	userLose.Info.ID = playerId
	//userLose.Info.LockMoney = 0
	userLose.Info.Money = lossMoney
	userLose.Info.Order = orderId
	userLose.Info.PayReason = "下注"
	//userLose.Info.PreMoney = 0
	userLose.Info.RoundId = roundId
	baseData.Data = userLose

	WriteMsgToCenter(baseData)
}
