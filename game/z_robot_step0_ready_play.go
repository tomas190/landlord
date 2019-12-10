package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/roomStatus"
	"landlord/msg/mproto"
)

// 1.给玩家和机器人发牌
func PushPlayerStartGameWithRobot(room *Room) {
	cards := CreateBrokenCard()
	//cards := CreateSortCard()

	//player, r1, r2 := getPlayersWithRobot(room)
	//// todo 玩家发牌策略
	//player.HandCards = append([]*Card{}, cards[:17]...)
	//var push mproto.PushStartGame
	//push.Cards = ChangeCardToProto(player.HandCards)
	//bytes, _ := proto.Marshal(&push)
	//PlayerSendMsg(player, PkgMsg(msgIdConst.PushStartGame, bytes))
	//
	//r1.HandCards = append([]*Card{}, cards[17:34]...)
	//r2.HandCards = append([]*Card{}, cards[34:51]...)
	//room.BottomCards = append([]*Card{}, cards[51:]...)
	//logger.Debug("底牌:")
	//PrintCard(room.BottomCards)

	// 随机发牌
	for _, v := range room.Players {
		v.HandCards = append([]*Card{}, cards[:17]...)
		SortCard(v.HandCards)
		logger.Debug("玩家" + v.PlayerInfo.PlayerId + "的牌：")
		PrintCard(v.HandCards)
		cards = append([]*Card{}, cards[17:]...)
		var push mproto.PushStartGame
		push.Cards = ChangeCardToProto(v.HandCards)
		bytes, _ := proto.Marshal(&push)
		PlayerSendMsg(v, PkgMsg(msgIdConst.PushStartGame, bytes))
	}
	room.BottomCards = append([]*Card{}, cards...)
	logger.Debug("底牌:")
	PrintCard(cards)
	room.Status = roomStatus.CallLandlord
	// 随机发牌

	// 随机叫地主写在发牌里面 是因为三个玩家如果都不叫 则可以直接调用 PushPlayerStartGameWithRobot 重新开始发牌逻辑
	DelaySomeTime(4)
	// 3.随机叫地主
	actionPlayerId := pushFirstCallLandlord(room)
	CallLandlord(room, actionPlayerId)
}

