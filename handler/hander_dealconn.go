package handler

import (
	"landlord/game"
	"landlord/mconst/msgIdConst"

	"landlord/mconst/userSessionStatus"
	"landlord/msg/mproto"
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
	"time"
)

// 处理连接
func DealOnConnect(session *melody.Session) {
	go closeUnLoginConn(session)

}

// 处理ping pong
func DealPingPong(session *melody.Session, data []byte) {
	//logger.Debug("=== DealPingPong ===")
	req := &mproto.PING{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		game.SendErrMsg(session, msgIdConst.Ping, "请求数据异常:"+err.Error())
		return
	}

	var pong mproto.PONG
	pong.Time = time.Now().Unix()
	bytes, _ := proto.Marshal(&pong)

	_ = session.WriteBinary(game.PkgMsg(msgIdConst.Pong, bytes))
	session.Set("ping", pong.Time)

}


// 5秒之后断开没有登录的连接
func closeUnLoginConn(session *melody.Session) {
	//time.Sleep(time.Second * 5) // 5秒之后未成功登录 则断开连接
	game.DelaySomeTime(5)

	isLoginSucc := game.GetSessionIsLogin(session)

	if !isLoginSucc {
		var push mproto.CloseConn
		push.Code = userSessionStatus.LoginTimeOutClose
		push.Msg = "login delay!"

		bytes, _ := proto.Marshal(&push)

		msg := game.PkgMsg(msgIdConst.CloseConn, bytes)
		err := session.CloseWithMsg(msg)
		if err != nil {
			logger.Error("推送断开消息失败 err:", err.Error())
		}
		logger.Debug("一个连接在5秒内未通信 被断开!")
	}
}

func showConnNums(m *melody.Melody, info string) {
	time.Sleep(time.Second * 1)
	logger.Debug(info+"之后的连接数:", m.Len())
	// logger.Debug(info+"之后的连接数map:", global.GetConnLen())
}
