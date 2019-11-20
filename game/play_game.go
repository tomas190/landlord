package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/playerAction"
	"landlord/mconst/roomStatus"
	"landlord/mconst/sysSet"
	"landlord/msg/mproto"
	"time"
)

func PlayGame(room *Room) {

	// 1. 玩家进入房间如果有玩家正待等待则与之开始游戏
	PushPlayerEnterRoom(room)
	DelaySomeTime(1)

	// 2.给玩家发牌
	PushPlayerStartGame(room)

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
	room.Status = roomStatus.CallLandlord

	DelaySomeTime(4)
	// 3.随机叫地主
	actionPlayerId := PushFirstCallLandlord(room)
	CallLandlord(room, actionPlayerId)
}

// 3.第一次开始叫地主
func PushFirstCallLandlord(room *Room) string {
	lastPosition := int32(RandNum(1, 3))
	lastPlayer := getPlayerByPosition(room, lastPosition)

	actionPosition := getNextPosition(lastPosition)
	actionPlayer := getPlayerByPosition(room, actionPosition)

	pushCallLandlordHelp(room, lastPlayer, actionPlayer, playerAction.CallLandlord)
	return actionPlayer.PlayerInfo.PlayerId
}

// 3.1.叫地主阶段 和抢地主阶段
func CallLandlord(room *Room, playerId string) {
	actionPlayer := room.Players[playerId]
	if actionPlayer == nil {
		logger.Error("房间里无此用户...")
		return
	}

	nextPosition := getNextPosition(actionPlayer.PlayerPosition)
	nextPlayer := getPlayerByPosition(room, nextPosition)

	lastPosition := getLastPosition(actionPlayer.PlayerPosition)
	lastPlayer := getPlayerByPosition(room, lastPosition)

	select {
	case action := <-actionPlayer.ActionChan:
		switch action.ActionType {
		case playerAction.CallLandlord: // 叫地主动作
			CallLandlordAction(room, actionPlayer, nextPlayer)
		case playerAction.GetLandlord: // 抢地主动作
			GetLandlordAction(room, actionPlayer, nextPlayer, lastPlayer)
		case playerAction.NotCallLandlord: // 不叫
			NotCallLandlordAction(room, actionPlayer, nextPlayer)
		case playerAction.NotGetLandlord: // 不抢
			NotGetLandlordAction(room, actionPlayer, nextPlayer, lastPlayer)
		}
	case <-time.After(time.Second * sysSet.GameDelayTime): // 自动进行不叫或者不抢
		if room.Status == roomStatus.CallLandlord {
			NotCallLandlordAction(room, actionPlayer, nextPlayer) // 不叫
		} else if room.Status == roomStatus.GetLandlord {
			NotGetLandlordAction(room, actionPlayer, nextPlayer, lastPlayer) // 不抢
		}
	}

}

