package game

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/wonderivan/logger"
	"strconv"
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
	logger.Debug(" 用户pkgId->", userAccount.Get("package_id").MustInt())

	var userLogin UserLoginCallBack
	var user PlayerInfo
	user.PlayerId = strconv.Itoa(userInfo.Get("id").MustInt())
	user.Name = userInfo.Get("game_nick").MustString()
	user.HeadImg = userInfo.Get("game_img").MustString()
	user.Gold = userAccount.Get("balance").MustFloat64()
	user.PlayerPkgId = userAccount.Get("package_id").MustInt()
	userLogin.Player = user
	userLogin.LoginStatus = true
	callChan := GetUserLoginCallChan(user.PlayerId)

	if callChan != nil {
		go func() {
			callChan <- &userLogin
		}()
	}

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
		//SendLogToCenter("ERR", "game/center_receive_msg.go", "76", "同步中心服输钱失败:"+ObjToString(data))
		logger.Debug("dealLossSoc！", data)
		logger.Debug("同步中心服输钱错误!")
		return
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

	isClose := GetSessionCloseTag(session)

	if isClose {
		logger.Debug("玩家已经断线：",id)
		ClearClosePlayer(session)
	}
	logger.Debug("玩家没有离线:", id)

}
