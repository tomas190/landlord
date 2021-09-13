package game

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/google/uuid"
	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2/bson"
)

func dealServerLogin(data *simplejson.Json) {

	code := data.Get("code").MustInt()

	if code != 200 {
		panic("服务器登录中心服失败！")
		return
	}

	bytes, _ := json.Marshal(data)
	logger.Debug("serverLoginResp:", string(bytes))
	// 设置各平台税收
	pTaxPercent := data.Get("msg").Get("globals")
	SetPlatformTaxPercent(pTaxPercent)
	// 设置各平台税收

	logger.Debug("登录中心服成功")

}

func dealUserLogin(data *simplejson.Json) {

	code := data.Get("code").MustInt()

	if code != 200 {
		logger.Debug("用户登录中心服失败！")
		return
	}
	logger.Debug("用户登录中心服成功")

	logger.Debug("<----------- 登录成功的用户信息 ----------->")
	userInfo := data.Get("msg").Get("game_user")
	userAccount := data.Get("msg").Get("game_account")
	logger.Debug(" 用户名称->", userInfo.Get("game_nick").MustString())
	logger.Debug(" 用户头像->", userInfo.Get("game_img").MustString())
	logger.Debug(" 用户金币->", userAccount.Get("balance").MustFloat64())
	logger.Debug(" 用户已鎖金币->", userAccount.Get("lock_balance").MustFloat64())
	logger.Debug(" 用户pkgId->", userInfo.Get("package_id").MustInt())

	var userLogin UserLoginCallBack
	var user PlayerInfo
	user.PlayerId = strconv.Itoa(userInfo.Get("id").MustInt())
	user.Name = userInfo.Get("game_nick").MustString()
	user.HeadImg = userInfo.Get("game_img").MustString()
	user.Gold = userAccount.Get("balance").MustFloat64()
	user.PlayerPkgId = userInfo.Get("package_id").MustInt()
	userLogin.Player = user
	userLogin.LoginStatus = true
	callChan := GetUserLoginCallChan(user.PlayerId)

	if callChan != nil {
		go func() {
			callChan <- &userLogin
		}()
	}

	lockGold := userAccount.Get("lock_balance").MustFloat64()
	needLock := user.Gold - lockGold
	msg := "user login lock all money"
	if lockGold > 0 {
		msg = "user login lock more money"
	}
	if needLock > 0 {
		order := bson.NewObjectId().Hex()
		OrderIDToOrderInfo.Store(order, OrderInfo{
			PlayerId: user.PlayerId,
			Event:    msg,
		})
		UserLockMoney(user.PlayerId, needLock, uuid.New().String(), msg, order)
	} else {
		logger.Debug("Login but user %v never need lock money", user.PlayerId)
	}
	// UserLockMoney(user.PlayerId, user.Gold, uuid.New().String(), "user login lock all money", order)

}

func dealWinSocer(data *simplejson.Json) {

	code := data.Get("code").MustInt()

	if code != 200 {
		//SendLogToCenter("ERR", "game/center_receive_msg.go", "62", "同步中心服赢钱失败:"+ObjToString(data))
		logger.Debug("dealWinSoc！", data)
		logger.Debug("同步中心服赢钱错误!")
		return
	}
	bytes, _ := json.Marshal(data)
	//playerId := data.Get("msg").Get("id").MustInt()
	//reducePlayerMsgNum(strconv.Itoa(playerId))
	fmt.Println("赢钱成功返回:", string(bytes))
	checkLoginOut(bytes)
}

func dealLossSocer(data *simplejson.Json) {

	code := data.Get("code").MustInt()

	if code != 200 {
		errorDealLossSocer(data)
		//SendLogToCenter("ERR", "game/center_receive_msg.go", "76", "同步中心服输钱失败:"+ObjToString(data))
		logger.Debug("dealLossSoc！", data)
		logger.Debug("同步中心服输钱错误!")
		return
	}

	// 刪除成功的输钱訂單
	order := data.Get("msg").Get("order").MustString()
	_, ok := OrderIDToOrderInfo.Load(order)
	if ok {
		OrderIDToOrderInfo.Delete(order)
	} else {
		logger.Debug("未找到符合的輸錢訂單 order=%v", order)
	}

	bytes, _ := json.Marshal(data)
	//playerId := data.Get("msg").Get("id").MustInt()
	//reducePlayerMsgNum(strconv.Itoa(playerId))
	fmt.Println("输钱成功返回:", string(bytes))
	checkLoginOut(bytes)
}

func dealUserLoginOutCenter(json *simplejson.Json) {

	code := json.Get("code").MustInt()
	if code != 200 {
		logger.Debug("用户登出失败！")
		logger.Debug("dealLoginOut！", json)
	}
}

func checkLoginOut(stByte []byte) {
	s, err := simplejson.NewJson(stByte)
	if err != nil {
		logger.Error("检查异常:", err.Error())
		logger.Error("检查异常:", string(stByte))
		return
	}

	idInt := s.Get("msg").Get("id").MustInt()
	id := strconv.Itoa(idInt)
	logger.Debug("玩家id:", id)

	session := GetAgent(id)
	if session == nil {
		logger.Error("获取session异常")
		return
	}

	isClose := GetSessionCloseTag(session)

	if isClose {
		logger.Debug("玩家已经断线：", id)
		ClearClosePlayer(session)
	}
	logger.Debug("玩家没有离线:", id)

}

// 上锁用户金币信息
// 当收到这条消息返回的时候 需要登出中心服 只有退出的时候才会解锁玩家的金币
func dealUserLockScore(data *simplejson.Json) {

	code := data.Get("code").MustInt()
	if code != 200 {

		// 上锁金币失败处理
		errorDealLockFail(data)
		return
	}

	// 刪除成功的上鎖訂單
	order := data.Get("msg").Get("order").MustString()
	_, ok := OrderIDToOrderInfo.Load(order)
	if ok {
		OrderIDToOrderInfo.Delete(order)
	} else {
		logger.Debug("未找到符合的上鎖訂單 order=%v", order)
	}
}

// 解锁用户金币信息
// 当收到这条消息返回的时候 需要登出中心服 只有退出的时候才会解锁玩家的金币
func dealUserUnlockScore(data *simplejson.Json) {
	code := data.Get("code").MustInt()

	if code != 200 {

		// 解锁金币失败处理
		errorDealUnlockFail(data)
		return
	}

	bytes, _ := json.Marshal(data)
	playerId := data.Get("msg").Get("id").MustInt()
	fmt.Println("解锁金币成功:", string(bytes))
	//

	agent := GetAgent(strconv.Itoa(playerId))
	if agent == nil {
		logger.Error("获取玩家session异常:", playerId)
		return
	}

	info, err := GetSessionPlayerInfo(agent)
	if err != nil {
		logger.Error("获取玩家信息异常:", playerId)
		logger.Error("err:", err.Error())
		return
	}

	password := GetSessionPassword(agent)
	UserLogoutCenter(info.PlayerId, password)
}

func errorDealLockFail(data *simplejson.Json) {

	// 刪除成功的上鎖訂單
	var msg string
	order := data.Get("msg").Get("order").MustString()
	orderInfo, ok := OrderIDToOrderInfo.Load(order)
	if ok {
		msg = fmt.Sprintf("鬥地主 中心服返回错误\n玩家 :%v\n事件 :%v\n錯誤訊息 :%v\n时间 : %v", orderInfo.(OrderInfo).PlayerId, orderInfo.(OrderInfo).Event, data.Get("msg"), time.Now().Format("2006-01-02 15:04:05"))
		kickRoomByUserID(orderInfo.(OrderInfo).PlayerId)
		OrderIDToOrderInfo.Delete(order)
	} else {
		msg = fmt.Sprintf("鬥地主 中心服返回错误\n未找到符合的orderID %v\n訊息：%v\n时间 : %v", order, data.Get("msg"), time.Now().Format("2006-01-02 15:04:05"))
	}
	HttpPostToTelegram(msg)
}

func errorDealUnlockFail(data *simplejson.Json) {
	// 打印错误信息
	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Error("marshal err:" + err.Error())
		return
	}
	logger.Debug("unlock fail resp:" + string(bytes))

	arr := data.Get("msg").Get("data").MustArray()
	order := data.Get("msg").Get("order").MustString()
	if len(arr) == 3 {
		if arr[0].(string) == "game account lock balance is not enough" {
			userCurrentGold, err := arr[1].(json.Number).Float64()
			if err != nil {
				logger.Debug("获取当前金币异常 err:", err.Error())
				return
			}
			// 根据返回的order获取对应的玩家id
			playerId := opMap.Get(order)
			agent := GetAgent(playerId)
			if agent == nil {
				logger.Error("获取玩家session异常:", playerId)
				if order == "" || playerId == "" {
					logger.Error("通过金币获取玩家session:", userCurrentGold)
					// 根据玩家当前金币 获取玩家ID
					playerId, agent = GetAgentByUserGold(userCurrentGold)
					if agent == nil {
						logger.Error("通过金币获取玩家session err:", userCurrentGold)
						return
					} else {
						logger.Info("通过金币获取玩家session success:", playerId, userCurrentGold)
					}
				}
			}
			UserUnLockMoney(playerId, userCurrentGold, uuid.New().String(), "user unlock fail lock again")
		}
	}
}

func errorDealLossSocer(data *simplejson.Json) {
	var msg string

	order := data.Get("msg").Get("order").MustString()
	orderInfo, ok := OrderIDToOrderInfo.Load(order)
	if ok {
		msg = fmt.Sprintf("鬥地主 中心服返回错误\n玩家 :%v\n事件 :%v\n錯誤訊息 :%v\n时间 : %v", orderInfo.(OrderInfo).PlayerId, orderInfo.(OrderInfo).Event, data.Get("msg"), time.Now().Format("2006-01-02 15:04:05"))
		kickRoomByUserID(orderInfo.(OrderInfo).PlayerId)
		OrderIDToOrderInfo.Delete(order)
	} else {
		msg = fmt.Sprintf("鬥地主 中心服返回错误\n未找到符合的orderID %v\n訊息：%v\n时间 : %v", order, data.Get("msg"), time.Now().Format("2006-01-02 15:04:05"))
	}
	HttpPostToTelegram(msg)
}
