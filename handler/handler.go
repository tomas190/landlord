package handler

import (
	"fmt"
	"landlord/game"
	"landlord/mconst/msgIdConst"
	"landlord/msg/mproto"

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
	// 如果客户端断开连接
	//if game.Server.UseRobot {
	//	value, exists := session.Get("waitChan")
	//	if exists {
	//		wc := value.(*game.WaitRoomChan)
	//		if !wc.IsClose {
	//			go func() {
	//				wc.WaitChan <- struct{}{}
	//			}()
	//		}
	//	}
	//}
	dealCloseConn(session)
	logger.Info("Handler OnDisconnect :")
}

// onErr
func OnErr(session *melody.Session, e error) {
	//go func() {
	//	if !session.IsClosed() {
	//		// 异常是否在房间游戏
	//		info, err := game.GetSessionPlayerInfo(session)
	//		if err != nil {
	//			if info!=nil {
	//				game.RemoveAgent(info.PlayerId)
	//				roomId := game.GetSessionRoomId(session)
	//				game.RemoveRoom(roomId)
	//				err = session.Close()
	//				if err != nil {
	//					logger.Info("断开session异常 :", err.Error())
	//				}
	//			}
	//		}
	//	}
	//}()
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
	case msgIdConst.ReqOutCardDo: // 玩家出牌请求 103
		game.ReqOutCardDo(session, data)
	case msgIdConst.ReqExitRoom: // 玩家退出房间 104
		game.ReqExitRoom(session, data)
	case msgIdConst.ReqGameHosting: // 玩家托管 105
		game.ReqGameHosting(session, data)
	case msgIdConst.ReqEnterRoomCheck: // 进入房间检查 106
		game.ReqEnterRoomCheck(session, data)
	case msgIdConst.ReqSendMsg: // 发送消息 107
		game.ReqSendMsg(session, data)
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

	if msgId != msgIdConst.RespGameHosting {
		return
	}
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
	case msgIdConst.PushStartGame:
		resp := &mproto.PushStartGame{}
		err := proto.Unmarshal(bytes[2:], resp)
		if err != nil {
			logger.Debug("打印服务器发送给客户端消息失败:", err.Error())
			return
		}
		fmt.Println("msgId.PushStartGame")
		game.PrintMsg("PushStartGame:", resp)
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

	case msgIdConst.PushSettlement:
		resp := &mproto.PushSettlement{}
		err := proto.Unmarshal(bytes[2:], resp)
		if err != nil {
			logger.Debug("打印服务器发送给客户端消息失败:", err.Error())
			return
		}
		fmt.Println("msgId.PushSettlement")
		game.PrintMsg("PushSettlement:", resp)
		fmt.Println()
	case msgIdConst.RespSendMsg:
		resp := &mproto.RespSendMsg{}
		err := proto.Unmarshal(bytes[2:], resp)
		if err != nil {
			logger.Debug("打印服务器发送给客户端消息失败:", err.Error())
			return
		}
		fmt.Println("msgId.RespSendMsg")
		game.PrintMsg("RespSendMsg:", resp)
		fmt.Println()
	case msgIdConst.PushCardCount:
		resp := &mproto.PushCardCount{}
		err := proto.Unmarshal(bytes[2:], resp)
		if err != nil {
			logger.Debug("打印服务器发送给客户端消息失败:", err.Error())
			return
		}
		fmt.Println("msgId.PushCardCount")
		game.PrintMsg("PushCardCount:", resp)
		fmt.Println()
	case msgIdConst.PushRoomRecover:
		resp := &mproto.PushRoomRecover{}
		err := proto.Unmarshal(bytes[2:], resp)
		if err != nil {
			logger.Debug("打印服务器发送给客户端消息失败:", err.Error())
			return
		}
		fmt.Println("msgId.PushRoomRecover")
		game.PrintMsg("PushRoomRecover:", resp)
		fmt.Println()
	case msgIdConst.RespGameHosting:
		resp := &mproto.RespGameHosting{}
		err := proto.Unmarshal(bytes[2:], resp)
		if err != nil {
			logger.Debug("打印服务器发送给客户端消息失败:", err.Error())
			return
		}
		fmt.Println("msgId.PushRoomRecover")
		game.PrintMsg("PushRoomRecover:", resp)
		fmt.Println()

		// ==========================================
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

// 处理用户断开连接
func dealCloseConn(session *melody.Session) {

	info, err := game.GetSessionPlayerInfo(session)
	if err != nil {
		logger.Debug("无用户session信息")
		return
	}
	logger.Debug("dealCloseConn: user=", info.PlayerId)
	game.RemoveWaitUser(info.PlayerId)

	flag := false
	room, b := game.IsPlayerInRoom(info.PlayerId)
	if b && room != nil {
		logger.Debug("dealCloseConn: user=", info.PlayerId, "is in room ", room.RoomId)
	} else {
		game.ClearClosePlayer(session)
		flag = true
	}

	roomId := game.GetSessionRoomId(session)
	if roomId == "" { // 证明用户不在游戏中
		if !flag {
			oldSession := game.GetAgent(info.PlayerId)
			if session == oldSession {
				game.ClearClosePlayer(session)
			} else {
				logger.Debug("dealCloseConn user=", info.PlayerId, " session != oldSession")
			}
		}
	} else { // 设置清除标记
		room := game.GetRoom(roomId)
		for _, p := range room.Players {
			if p.PlayerInfo.PlayerId == info.PlayerId {
				p.IsExitRoom = true
				p.IsCloseSession = true
				game.SetSessionCloseTag(session, true)
			}
		}
	}
}
