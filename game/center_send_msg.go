package game

import (
	"github.com/wonderivan/logger"
	"time"
)

func UserLogin(playerId, password, token string) {

	logger.Debug("<-------- UserLoginCenter -------->")
	baseData := &ToCenterMessage{}
	baseData.Event = msgUserLogin

	// 要改成判断token
	if token != "" {
		baseData.Data = &UserReqToken{
			ID:      playerId,
			Token:   token,
			GameId:  Server.GameId,
			DevName: Server.DevName,
			DevKey:  Server.DevKey}
	} else {
		baseData.Data = &UserReqPassword{
			ID:       playerId,
			PassWord: password,
			GameId:   Server.GameId,
			DevName:  Server.DevName,
			DevKey:   Server.DevKey}
	}

	WriteMsgToCenter(baseData)

	//加入待处理map，等待处理
	//c4c.waitUser[userId] = &UserCallback{}
	//c4c.waitUser[userId].Data.ID = userId
	//c4c.waitUser[userId].Callback = callback
}

//UserLogoutCenter 用户登出
func UserLogoutCenter(userId string, password string) {
	logger.Debug("<-------- UserLogoutCenter  -------->")
	base := &ToCenterMessage{}
	base.Event = msgUserLogout
	base.Data = &UserReq{
		ID:       userId,
		Password: password,
		GameId:   Server.GameId,
		Token:    password,
		DevName:  Server.DevName,
		DevKey:   Server.DevKey,
	}

	var num int
LoginOut:
	can := canLoginOut(userId)
	if can || num == 3 {
		logger.Debug("loginOut normal.", num)
		WriteMsgToCenter(base)
		RemoveAgent(userId)
	} else {
		t := time.Tick(time.Second * 1)
		<-t
		num++
		goto LoginOut
	}

	// 发送消息到中心服
	// 延时1秒后发送退出中心服消息
	//logger.Debug("loginOut delay.")
	//WriteMsgToCenter(base)
	//global.RemoveAgent(userId)
}



//UserSyncWinScore 同步赢分
func UserSyncWinScore(playerId string, winMoney float64, roundId, orderId string) {
	addPlayerMsgNum(playerId) // 增加消息值   // 收到中心服务的时候减少值
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


//UserSyncWinScore 同步输分
func UserSyncLoseScore(playerId string, lossMoney float64, roundId, orderId string) {
	addPlayerMsgNum(playerId) // 增加消息值   // 收到中心服务的时候减少值
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

func canLoginOut(userId string) bool {
	agent := GetAgent(userId)
	if agent == nil {
		logger.Debug("已经移除玩家信息")
		return true
	}
	value, exists := agent.Get("msgNum")
	if exists {
		if value.(int) == 0 {
			return true
		} else {
			return false
		}
	}
	return true
}


func addPlayerMsgNum(playerId string) {
	agent := GetAgent(playerId)
	if agent == nil {
		logger.Error("无", playerId, "的session信息")
		return
	}

	num, ok := agent.Get("msgNum")
	if ok {
		agent.Set("msgNum", num.(int)+1)
	} else {
		agent.Set("msgNum", 1)
	}

}

func reducePlayerMsgNum(playerId string) {

	agent := GetAgent(playerId)
	if agent == nil {
		logger.Error("无", playerId, "的session信息")
		return
	}

	num, ok := agent.Get("msgNum")
	if ok {
		agent.Set("msgNum", num.(int)-1)
	} else {
		logger.Error("不正常的操作流程")
		// agent.Set("msgNum", 1)
	}
}