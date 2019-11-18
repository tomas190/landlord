package game

import (
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

	var userLogin UserLoginCallBack
	var user PlayerInfo
	user.PlayerId = strconv.Itoa(userInfo.Get("id").MustInt())
	user.Name = userInfo.Get("game_nick").MustString()
	user.HeadImg = userInfo.Get("game_img").MustString()
	user.Gold = userAccount.Get("balance").MustFloat64()
	userLogin.Player = user
	userLogin.LoginStatus = true
	callChan := GetUserLoginCallChan(user.PlayerId)

	if callChan != nil {
		callChan <- &userLogin
	}

}

func dealWinSocer(json *simplejson.Json) {

	code := json.Get("code").MustInt()

	if code != 200 {
		logger.Debug("dealWinSoc！", json)
		logger.Debug("同步中心服赢钱错误!")
	}

	//fmt.Println("赢钱成功返回:",json)

}

func dealLossSocer(json *simplejson.Json) {

	code := json.Get("code").MustInt()

	if code != 200 {
		logger.Debug("dealLossSoc！", json)
		logger.Debug("同步中心服输钱错误!")
	}
}

func dealUserLoginOutCenter(json *simplejson.Json) {

	code := json.Get("code").MustInt()
	if code != 200 {
		logger.Debug("用户登出失败！")
		logger.Debug("dealLoginOut！", json)
	}
}
