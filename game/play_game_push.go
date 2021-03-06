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

	var countDown int32
	countDown = 19
	for _, v := range room.Players {
		if v.IsCanDo {
			logger.Debug("=============== 当前正在操作的玩家:", v.PlayerInfo.PlayerId, "====================")
			countDown = v.WaitingTime
		}
	}

	var resp mproto.PushRoomRecover
	resp.Players = ChangePlayerToRecoverPlayer(room.Players, playerId)
	resp.BottomPoint = room.RoomClass.BottomPoint
	resp.Multi = room.MultiAll
	resp.Countdown = countDown
	resp.LandLordPlayerId = room.LandlordPlayerId
	resp.EffectiveCard = ChangeCardToProto(room.EffectiveCard)
	resp.EffectiveCardType = room.EffectiveType
	if room.LandlordPlayerId != "" {
		landPlayer := room.Players[room.LandlordPlayerId]
		if landPlayer != nil {
			resp.LandlordPosition = landPlayer.PlayerPosition
			resp.BottomCard = ChangeCardToProto(room.BottomCards)

		} else {
			logger.Error("该房间无玩家信息 !!!incredible")
		}
	}
	resp.RoomStatus = room.Status
	bytes, _ := proto.Marshal(&resp)
	_ = session.WriteBinary(PkgMsg(msgIdConst.PushRoomRecover, bytes))
}

// 推送记牌器
func pushCardCount(room *Room) {

	players := room.Players
	throwCards := room.ThrowCards
	for _, v := range players {
		result := countCards(v.HandCards, throwCards)
		bytes, _ := proto.Marshal(result)
		PlayerSendMsg(v, PkgMsg(msgIdConst.PushCardCount, bytes))
	}

}

// 计算记牌器
func countCards(handCards, roomThrowCards []*Card) *mproto.PushCardCount {
	m := originalCardNum()
	throwCards := append(roomThrowCards, handCards...)
	for i := 0; i < len(throwCards); i++ {
		m[throwCards[i].Value] = m[throwCards[i].Value] - 1
	}

	var result mproto.PushCardCount
	result.CardCount = m
	return &result
}
