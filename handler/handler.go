package handler

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
	"landlord/game"
	"landlord/mconst/msgIdConst"
	"landlord/msg/mproto"
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
		//logger.Debug(string(bytes))
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
	case msgIdConst.ReqGetLandlordDo: // 玩家抢地主请求 102
		game.ReqGetLandlordDo(session, data)
	case msgIdConst.ReqOutCardDo: // 玩家出牌请求 102
		game.ReqOutCardDo(session, data)
	default:
		logger.Error("未知指令")
	}

}

// OnSentMessageBinary  // debug 开发调试
func OnSentMessageBinary(session *melody.Session, bytes []byte) {

	if len(bytes) <= 0 {
		logger.Debug("auto ")
		return
	}
	msgId := game.GetMsgId(bytes[:2])
	if msgId == msgIdConst.Pong {
		return
	}

	return

	// todo  huck
	// logger.Info("Handler OnSentMessageBinary :", )
	switch msgId {
	case msgIdConst.RespLogin:

		resp := &mproto.RespLogin{}
		err := proto.Unmarshal(bytes[2:], resp)
		if err != nil {
			logger.Debug("打印服务器发送给客户端消息失败:", err.Error())
			return
		}
		fmt.Println("msgId.RespLogin")
		game.PrintMsg("RespLogin:", resp)
		fmt.Println()
	case msgIdConst.PushRoomClassify:
		resp := &mproto.PushRoomClassify{}
		err := proto.Unmarshal(bytes[2:], resp)
		if err != nil {
			logger.Debug("打印服务器发送给客户端消息失败:", err.Error())
			return
		}
		fmt.Println("msgId.PushRoomClassify")
		game.PrintMsg("PushRoomClassify:", resp)
		fmt.Println()
	case msgIdConst.PushRoomPlayer:
		resp := &mproto.PushRoomPlayer{}
		err := proto.Unmarshal(bytes[2:], resp)
		if err != nil {
			logger.Debug("打印服务器发送给客户端消息失败:", err.Error())
			return
		}
		fmt.Println("msgId.PushRoomPlayer")
		game.PrintMsg("PushRoomPlayer:", resp)
		fmt.Println()
	case msgIdConst.PushCallLandlord:
		resp := &mproto.PushGetLandlord{}
		err := proto.Unmarshal(bytes[2:], resp)
		if err != nil {
			logger.Debug("打印服务器发送给客户端消息失败:", err.Error())
			return
		}
		fmt.Println("msgId.PushCallLandlord")
		game.PrintMsg("PushCallLandlord:", resp)
		fmt.Println()
	case msgIdConst.PushOutCard:
		resp := &mproto.PushOutCard{}
		err := proto.Unmarshal(bytes[2:], resp)
		if err != nil {
			logger.Debug("打印服务器发送给客户端消息失败:", err.Error())
			return
		}
		fmt.Println("msgId.PushOutCard")
		game.PrintMsg("PushOutCard:", resp)
		fmt.Println()

	case msgIdConst.ErrMsg:
		resp := &mproto.ErrMsg{}
		err := proto.Unmarshal(bytes[2:], resp)
		if err != nil {
			logger.Debug("打印服务器发送给客户端消息失败:", err.Error())
			return
		}
		fmt.Println("msgid.ErrMsg:")
		game.PrintMsg("errorMsg:", resp)
		fmt.Println()

	}

}
