package game

import (
	"landlord/mconst/msgIdConst"
	"landlord/msg/mproto"
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
)

// 玩家发送消息
func PlayersSendMsg(players []*Player, msg []byte) {

	for i:=0;i<len(players) ;i++  {
		player:=players[i]
		if player.IsRobot {
			continue
		}

		session := player.Session

		if session == nil {
			logger.Debug("异常: player.session is nil ! id:", player.PlayerInfo.PlayerId)
			continue
		}

		_ = session.WriteBinary(msg)
	}

}

// 玩家发送消息
func PlayerSendMsg(player *Player, msg []byte) {

	if player.IsRobot {
		return
	}

	session := player.Session

	if session == nil {
		logger.Debug("异常: player.session is nil ! id:", player.PlayerInfo.PlayerId)
		return
	}

	_ = session.WriteBinary(msg)
}

// 发送错误消息
func SendErrMsg(session *melody.Session, mId uint16, msg string) {

	var errMsgPush mproto.ErrMsg
	errMsgPush.MsgId = int32(mId)
	errMsgPush.Msg = msg

	bytes, _ := proto.Marshal(&errMsgPush)

	_ = session.WriteBinary(PkgMsg(msgIdConst.ErrMsg, bytes))

}