package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
	"landlord/mconst/msgIdConst"
	"landlord/msg/mproto"
)

// 推送房间 分类消息
func PushRoomClassify(session *melody.Session) {
	resp := roomClassify
	bytes, _ := proto.Marshal(resp)
	_ = session.WriteBinary(PkgMsg(msgIdConst.PushRoomClassify, bytes))

}

// 推送恢复房间
func PushRecoverRoom(session *melody.Session, room *Room, playerId string) {

	player := room.Players[playerId]
	if player == nil {
		logger.Error("该房间无玩家信息 !!!incredible")
		SendErrMsg(session, msgIdConst.ReqEnterRoom, "恢复房间信息失败,无用户信息")
		return
	}

	var resp mproto.PushRoomRecover
	resp.Players = ChangePlayerToRecoverPlayer(room.Players, playerId)
	resp.BottomPoint = room.RoomClass.BottomPoint
	resp.Multi = room.MultiAll
	resp.Countdown = player.WaitingTime

	bytes, _ := proto.Marshal(&resp)
	_ = session.WriteBinary(PkgMsg(msgIdConst.PushRoomRecover, bytes))
}
