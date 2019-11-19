package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"landlord/mconst/msgIdConst"
	"landlord/msg/mproto"
)

func PlayGame(room *Room) {

	//1. 玩家进入房间如果有玩家正待等待则与之开始游戏
	PushPlayerEnterRoom(room)

}

// 1. 玩家进入房间如果有玩家正待等待则与之开始游戏
func PushPlayerEnterRoom(room *Room) {
	if len(room.Players) != 3 {
		logger.Error("异常房间")
		return
	}

	var push mproto.PushRoomPlayer
	push.Players = ChangeArrPlayerToRoomPlayerProto(room.Players)
	bytes, _ := proto.Marshal(&push)

	// 推送房间双方玩家信息
	players := room.Players
	MapPlayersSendMsg(players, PkgMsg(msgIdConst.PushRoomPlayer, bytes))
}
