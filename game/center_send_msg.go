package game

import (
	"fmt"
	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2/bson"
	"math"
	"strconv"
	"time"
)

func UserLogin(playerId, password, token string) {

	logger.Debug("<-------- UserLoginCenter -------->")
	baseData := &ToCenterMessage{}
	baseData.Event = msgUserLogin

	pId, err := strconv.Atoi(playerId)
	if err!=nil {
		logger.Debug("非法用户id:",playerId)
	}

	// 要改成判断token
	if token != "" {
		baseData.Data = &UserReqToken{
			Id:      pId,
			Token:   token,
			GameId:  Server.GameId,
			DevName: Server.DevName,
			DevKey:  Server.DevKey}
	} else {
		baseData.Data = &UserReqPassword{
			Id:       pId,
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
	//DelaySomeTime(1)
	pId, err := strconv.Atoi(userId)
	if err!=nil {
		logger.Debug("非法用户id:",userId)
	}
	logger.Debug("<-------- UserLogoutCenter  -------->")
	base := &ToCenterMessage{}
	base.Event = msgUserLogout
	base.Data = &UserReq{
		//ID:       userId,
		ID:       pId,
		Password: password,
		GameId:   Server.GameId,
		Token:    password,
		DevName:  Server.DevName,
		DevKey:   Server.DevKey,
	}

	//DelaySomeTime(1)
	WriteMsgToCenter(base)
	RemoveAgent(userId)
	//	var num int
	//LoginOut:
	//	can := canLoginOut(userId)
	//	if can || num >=5 {
	//		logger.Debug("loginOut normal.", num)
	//		WriteMsgToCenter(base)
	//		RemoveAgent(userId)
	//	} else {
	//		t := time.Tick(time.Millisecond*300)
	//		<-t
	//		num++
	//		goto LoginOut
	//	}

	// 发送消息到中心服
	// 延时1秒后发送退出中心服消息
	//logger.Debug("loginOut delay.")
	//WriteMsgToCenter(base)
	//global.RemoveAgent(userId)
}

//UserSyncWinScore 同步赢分
func UserSyncWinScore(playerId string, winMoney float64, roundId string) {
	pId, err := strconv.Atoi(playerId)
	if err!=nil {
		logger.Debug("非法用户id:",playerId)
	}
	//	addPlayerMsgNum(playerId) // 增加消息值   // 收到中心服务的时候减少值
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
	userWin.Info.ID = pId
	//userWin.Info.LockMoney = 0
	userWin.Info.Money = winMoney
	//userWin.Info.Order = orderId
	userWin.Info.Order = bson.NewObjectId().Hex()
	userWin.Info.PayReason = "对局"
	//userWin.Info.PreMoney = 0
	userWin.Info.RoundId = roundId
	baseData.Data = userWin

	logger.Debug("发送赢分指令:")
	PrintMsg("sendCenterMsg:", baseData)

	WriteMsgToCenter(baseData)
}

//UserSyncWinScore 同步输分
func UserSyncLoseScore(playerId string, lossMoney float64, roundId string) {
	// addPlayerMsgNum(playerId) // 增加消息值   // 收到中心服务的时候减少值
	pId, err := strconv.Atoi(playerId)
	if err!=nil {
		logger.Debug("非法用户id:",playerId)
	}
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
	userLose.Info.ID = pId
	//userLose.Info.LockMoney = 0
	userLose.Info.Money = lossMoney
	//userLose.Info.Order = orderId
	userLose.Info.Order = bson.NewObjectId().Hex()
	userLose.Info.PayReason = "对局"
	//userLose.Info.PreMoney = 0
	userLose.Info.RoundId = roundId
	userLose.Info.BetMoney = math.Abs(lossMoney)
	baseData.Data = userLose

	WriteMsgToCenter(baseData)
}

func NoticeWinMoreThan(playerId, playerName string, winGold float64) {
	logger.Debug("<-------- NoticeWinMoreThan  -------->")
	pId, err := strconv.Atoi(playerId)
	if err!=nil {
		logger.Debug("非法用户id:",playerId)
	}
	//style := `<size=20><color=YELLOW>恭喜!</color><color=orange>888888888</color><color=YELLOW>在</color></><color=WHITE><b><size=25>龙虎斗</color></b></><color=YELLOW><size=20>中一把赢了</color></><color=YELLOW><b><size=35>8888.88</color></b></><color=YELLOW><size=20>金币！</color></>`

	//msg := fmt.Sprintf("<size=20><color=YELLOW>恭喜!</color><color=orange>%s</color><color=YELLOW>在</color></><color=WHITE><b><size=25>斗地主</color></b></><color=YELLOW><size=20>中一把赢了</color></><color=YELLOW><b><size=35>%.2f</color></b></><color=YELLOW><size=20>金币！</color></>", playerName, winGold)
	msg := fmt.Sprintf("<size=20><color=yellow>恭喜!</color><color=orange>%s</color><color=yellow>在</color></><color=orange><size=25>斗地主</color></><color=yellow><size=20>中一把赢了</color></><color=yellow><size=30>%2.f</color></><color=yellow><size=25>金币！</color></>", playerName, winGold)
	base := &ToCenterMessage{}
	base.Event = msgWinMoreThanNotice
	base.Data = &Notice{
		DevName: Server.DevName,
		DevKey:  Server.DevKey,
		ID:      pId,
		GameId:  Server.GameId,
		Type:    2000,
		Message: msg,
		Topic:   "系统提示",
	}
	WriteMsgToCenter(base)
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
