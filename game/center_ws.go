package game

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
)

var centerServerConn *websocket.Conn

//CreatConnect 和Center建立链接
func ConnectCenterWs() {
	wsAddr := "ws://" + Server.CenterDomain + "/"
	//c4c.centerUrl = "ws" + strings.TrimPrefix(conf.Server.CenterServer, "http") //域名生成使用
	logger.Debug("centerAddr:", wsAddr)
	conn, _, err := websocket.DefaultDialer.Dial(wsAddr, nil)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
	centerServerConn = conn

	go startCenterServer()
	onBreath()

	// 服务器登录中心服务
	loginCenterServer()

}

// 每隔5秒发送空的数据 不然连接会被断开
func onBreath() {
	go func() {
		for {
			time.Sleep(time.Second * 4)
			err := centerServerConn.WriteMessage(websocket.TextMessage, []byte(""))
			if err != nil {
				logger.Debug("发送心跳失败")
				// 尝试重连
				ConnectCenterWs()
				break
			}
		}
	}()

}

func startCenterServer() {

	defer func() {
		_ = centerServerConn.Close()
	}()
	for {
		msgType, p, err := centerServerConn.ReadMessage()
		if err != nil {
			logger.Error("read msg err:", err.Error())
			break
		}
		// logger.Info("receive msg str:", string(p))
		onReceiveCenterMsg(msgType, p)
	}

}

func loginCenterServer() {

	port, err := strconv.Atoi(Server.Port)
	if err != nil {
		logger.Error("port非法:", port)
	}

	var loginReq ServerLoginReq
	loginReq.Host = Server.CenterDomain
	//loginReq.Port = Server.Port
	loginReq.Port = port
	loginReq.DevKey = Server.DevKey
	loginReq.GameId = Server.GameId
	// 	loginReq.Token = TokenOfCenter
	loginReq.DevName = Server.DevName

	var msg ToCenterMessage
	msg.Event = msgServerLogin
	msg.Data = loginReq

	logger.Debug("发送服务器登录")
	WriteMsgToCenter(msg)

}

func WriteMsgToCenter(data interface{}) {
	bytes, _ := json.Marshal(data)
	logger.Debug("WriteMsgToCenter :", string(bytes))
	err := centerServerConn.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		logger.Error("write msg err:", err.Error())
	}

}

func onReceiveCenterMsg(messType int, msgFromCenter []byte) {
	logger.Debug("onReceiveCenterMsg:", string(msgFromCenter))
	if messType != websocket.TextMessage {
		return
	}

	msgJson, err := simplejson.NewJson(msgFromCenter)
	if err != nil {
		logger.Debug("解析中心服msg错误:", err.Error())
		logger.Debug("解析错误中心服msg:", string(msgFromCenter))
		return
	}

	event := msgJson.Get("event").MustString()
	data := msgJson.Get("data")

	switch event {
	case msgServerLogin:
		dealServerLogin(data)
	case msgUserLogin:
		dealUserLogin(data)
	case msgUserWinScore:
		dealWinSocer(data)
	case msgUserLoseScore:
		dealLossSocer(data)
	case msgUserLogout:
		dealUserLoginOutCenter(data)
	case msgUserLockScore:
		dealUserLockScore(data)
	case msgUserUnLockScore:
		dealUserUnlockScore(data)
	default:
		logger.Error("Receive a message but don't identify~")
		logger.Error(string(msgFromCenter))
	}

}
