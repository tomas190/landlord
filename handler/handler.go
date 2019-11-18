package handler

import (
	"landlord/game"
	"landlord/mconst/msgIdConst"
	"landlord/msg/mproto"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
)

// onMessage
func OnMessage(m *melody.Melody, session *melody.Session, bytes []byte) {
	logger.Error("OnMessage ! 居然不是发送二进制！")
}

// onClose
func OnClose(m *melody.Melody, session *melody.Session, i int, s string) error {

	return nil
}

// onConnection
func OnConnect(m *melody.Melody, session *melody.Session) {

	defer func() {
		logger.Info("当前连接数:", m.Len())
	}()

	if m.Len() >= game.Server.MaxConn {
		logger.Debug("连接数已经超过限制:", m.Len())
		var m mproto.ErrMsg
		m.Msg = "连接数已满"
		bytes, _ := proto.Marshal(&m)
		_ = session.CloseWithMsg(bytes)
		return
	}
	// 处理新的连接
	DealOnConnect(session)

}

// onDisconnection
func OnDisconnect(m *melody.Melody, session *melody.Session) {
	logger.Info("Handler OnDisconnect :")
}

// onErr
func OnErr(session *melody.Session, e error) {
	logger.Info("Handler OnErr :", e.Error())
}

// onPong
func OnPong(session *melody.Session) {

	// logger.Info("Handler OnPong :")
}

// onSentMsg
func OnSentMessage(session *melody.Session, bytes []byte) {
	logger.Info("Handler OnSentMessage :", string(bytes))

}

// 收到客户端发送来得消息
// OnMessageBinary
func OnMessageBinary(m *melody.Melody, session *melody.Session, bytes []byte) {
	//logger.Info("Handler OnMessageBinary original :", bytes)
	//logger.Info("Handler OnMessageBinary :", string(bytes))
	dataLen := len(bytes)
	if dataLen <= 2 {
		logger.Debug(string(bytes))
		_ = session.WriteBinary([]byte("data is nil !"))
		return
	}
	msgId := game.GetMsgId(bytes[:2])

	data := bytes[2:]

	fmt.Println("msgId:", msgId)
	fmt.Println("data:", data)

	switch msgId {
	case msgIdConst.Ping: // ping pong  // 0
		DealPingPong(session, data)
	case msgIdConst.ReqLogin: // 用户登录 // 100
		ReqLogin(m, session, data)
	case msgIdConst.ReqEnterRoom: // 进入房间  101  // push 301 302
		game.ReqEnterRoom(session, data)
	case msgIdConst.ReqDoAction: // 获取玩家列表  102
		//game.ReqDoAction(session, data)
	default:
		logger.Error("未知指令")
	}

}

// OnSentMessageBinary  // debug 开发调试
func OnSentMessageBinary(session *melody.Session, bytes []byte) {

}
