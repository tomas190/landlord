package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/roomStatus"
	"landlord/msg/mproto"
)

func PlayGame(room *Room) {

	// 1. 玩家进入房间如果有玩家正待等待则与之开始游戏
	PushPlayerEnterRoom(room)
	DelaySomeTime(1)

	// 2.给玩家发牌
	PushPlayerStartGame(room)
	DelaySomeTime(4)

	// 3.随机叫地主

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

// 2.给玩家发牌
func PushPlayerStartGame(room *Room) {
	cards := CreateBrokenCard()
	players := room.Players
	for _, v := range players {
		v.HandCards = cards[:17]
		logger.Debug("玩家" + v.PlayerInfo.PlayerId + "的牌：")
		PrintCard(v.HandCards)
		cards = cards[17:]
		var push mproto.PushStartGame
		push.Cards = ChangeCardToProto(v.HandCards)
		bytes, _ := proto.Marshal(&push)
		PlayerSendMsg(v, PkgMsg(msgIdConst.PushStartGame, bytes))
	}
	room.bottomCards = cards
	logger.Debug("底牌:")
	PrintCard(cards)
	room.Status = roomStatus.GetLandlord
}

// 3.叫地主阶段
func PushGetLandlord(room *Room) {




}
