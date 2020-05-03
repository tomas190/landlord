package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/roomStatus"
	"landlord/msg/mproto"
)

/*
	游戏流程控制
*/

func PlayGame(room *Room) {

	// 1. 玩家进入房间如果有玩家正待等待则与之开始游戏
	PushPlayerEnterRoom(room)
	DelaySomeTime(1)

	// 2.给玩家发牌
	PushPlayerStartGame(room)

	// ..．流程控制到这里结束　发牌  抢地主  玩牌 直接由 PushPlayerStartGame 开始 且循环

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

func PushFakerPlayerEnterRoom(players map[string]*Player, realPlayer *Player) {

	playerProto := ChangeArrPlayerToRoomPlayerProto(players)

	var resp mproto.PushRoomPlayer

	resp.Players = playerProto

	bytes, _ := proto.Marshal(&resp)

	_ = realPlayer.Session.WriteBinary(PkgMsg(msgIdConst.PushRoomPlayer, bytes))

}

func PushFakerPlayerQuitRoom(player *Player) {

	rPlayer := ChangePlayerToRoomPlayerProto(player)
	var resp mproto.PushRoomPlayer
	var rPs []*mproto.RoomPlayer
	rPs = append(rPs, rPlayer)
	resp.Players = rPs
	bytes, _ := proto.Marshal(&resp)
	_ = player.Session.WriteBinary(PkgMsg(msgIdConst.PushRoomPlayer, bytes))
}

// 2.给玩家发牌
func PushPlayerStartGame(room *Room) {
	//cards := CreateBrokenCard()
	level := RandNum(35, 44)
	//cards := CreateGoodCardAll(level)
	i, i2, i3, bCards := CreateGoodCard(level)
	//cards := CreateSortCard()
	//players := room.Players
	//for _, v := range players {
	//	v.HandCards = append([]*Card{}, cards[:17]...)
	//	SortCard(v.HandCards)
	//	logger.Debug("玩家" + v.PlayerInfo.PlayerId + "的牌：")
	//	PrintCard(v.HandCards)
	//	cards = append([]*Card{}, cards[17:]...)
	//	var push mproto.PushStartGame
	//	push.Cards = ChangeCardToProto(v.HandCards)
	//	bytes, _ := proto.Marshal(&push)
	//	PlayerSendMsg(v, PkgMsg(msgIdConst.PushStartGame, bytes))
	//}
	players := room.Players
	var num int
	for _, v := range players {
		if num == 0 {
			v.HandCards = append([]*Card{}, i...)
		} else if num == 1 {
			v.HandCards = append([]*Card{}, i2...)
		} else {
			v.HandCards = append([]*Card{}, i3...)
		}
		num++
		SortCard(v.HandCards)
		logger.Debug("玩家" + v.PlayerInfo.PlayerId + "的牌：")
		PrintCard(v.HandCards)
		var push mproto.PushStartGame
		push.Cards = ChangeCardToProto(v.HandCards)
		bytes, _ := proto.Marshal(&push)
		PlayerSendMsg(v, PkgMsg(msgIdConst.PushStartGame, bytes))
	}

	room.BottomCards = append([]*Card{}, bCards...)
	logger.Debug("底牌:")
	PrintCard(bCards)
	room.Status = roomStatus.CallLandlord

	// 随机叫地主写在发牌里面 是因为三个玩家如果都不叫 则可以直接调用 PushPlayerStartGame 重新开始发牌逻辑
	DelaySomeTime(4)
	// 3.随机叫地主
	actionPlayerId := pushFirstCallLandlord(room)
	CallLandlord(room, actionPlayerId)
}
