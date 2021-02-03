package game

import "landlord/msg/mproto"

/*
	struct to proto
*/

func ChangeRoomClassifyToProto(item *RoomClassify) *mproto.RoomClassify {

	var result mproto.RoomClassify
	result.RoomType = item.RoomType
	result.BottomEnterPoint = item.BottomEnterPoint
	result.BottomPoint = item.BottomPoint

	return &result
}

func ChangePlayerInfoToProto(item *PlayerInfo) *mproto.PlayerInfo {

	var result mproto.PlayerInfo
	result.Gold = item.Gold
	result.PlayerName = item.Name
	result.PlayerId = item.PlayerId
	result.PlayerImg = item.HeadImg

	return &result
}

func ChangePlayerToRoomPlayerProto(player *Player) *mproto.RoomPlayer {
	var result mproto.RoomPlayer
	result.Players = ChangePlayerInfoToProto(player.PlayerInfo)
	result.Position = player.PlayerPosition
	return &result

}

func ChangeArrPlayerToRoomPlayerProto(players map[string]*Player) []*mproto.RoomPlayer {

	var result []*mproto.RoomPlayer

	for _, v := range players {
		p := ChangePlayerToRoomPlayerProto(v)
		result = append(result, p)
	}

	return result
}

func ChangeCardToProto(cards []*Card) []*mproto.Card {
	var result []*mproto.Card

	for i := 0; i < len(cards); i++ {
		var mc mproto.Card
		mc.Value = cards[i].Value
		mc.Suit = cards[i].Suit
		result = append(result, &mc)
	}
	return result
}

func ChangeProtoToCard(cards []*mproto.Card) []*Card {
	var result []*Card

	for i := 0; i < len(cards); i++ {
		var mc Card
		mc.Value = cards[i].Value
		mc.Suit = cards[i].Suit
		result = append(result, &mc)
	}
	return result
}

func ChangePlayerToRecoverPlayer(mPlayers map[string]*Player, rPlayerId string) []*mproto.RecoverPlayer {
	var resp []*mproto.RecoverPlayer

	for _, v := range mPlayers {
		var rp mproto.RecoverPlayer
		rp.Player = ChangePlayerToRoomPlayerProto(v)
		rp.IsGameHosting = v.IsGameHosting
		rp.IsMustPlay = v.IsMustDo
		rp.IsAction = v.IsCanDo
		rp.LastOutCard = ChangeCardToProto(v.LastOutCard)
		rp.LastAction = v.LastAction
		rp.CardsLen = int32(len(v.HandCards))
		if rPlayerId == v.PlayerInfo.PlayerId { // 如果是恢复玩家则将手牌传值
			rp.Cards = ChangeCardToProto(v.HandCards)
		}
		resp = append(resp, &rp)

	}
	return resp
}

// ==================== proto to struct ====================

func ChangePlayerP2S(item mproto.PlayerInfo) PlayerInfo {

	return PlayerInfo{
		PlayerId: item.PlayerId,
		Name:     item.PlayerName,
		HeadImg:  item.PlayerImg,
		Gold:     item.Gold,
	}
}
