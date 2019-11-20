package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/sysSet"
	"landlord/msg/mproto"
)

/*
出牌阶段
*/

func PlayingGame(room *Room, actionPlayerId string) {
	actionPlayer := room.Players[actionPlayerId]
	if actionPlayer == nil {
		logger.Error("PlayingGame 房间里无此用户...!!!incredible")
		return
	}

	//nextPosition := getNextPosition(actionPlayer.PlayerPosition)
	//nextPlayer := getPlayerByPosition(room, nextPosition)
	//
	//lastPosition := getLastPosition(actionPlayer.PlayerPosition)
	//lastPlayer := getPlayerByPosition(room, lastPosition)

}

// 推送必须出牌的消息 (第一次出牌开始 / 玩家出牌后 后面两个玩家都不要)
func pushMustOutCard(room *Room, playerId string) string {
	actionPlayer := getPlayerByPlayerId(room, playerId)

	return actionPlayer.PlayerInfo.PlayerId
}

func pushOutCardHelp(room *Room, actionPlayer, lastPlayer *Player, isMustPlay bool, lastAction int32, outCard []*Card) {

	var push mproto.PushPlayCard

	push.LastPlayerId = lastPlayer.PlayerInfo.PlayerId
	push.LastPlayerPosition = lastPlayer.PlayerPosition
	push.LastAction = lastAction
	push.LastPlayerCards = ChangeCardToProto(outCard)
	push.LastPlayerCardsType = getCardType(outCard)

	push.PlayerPosition = actionPlayer.PlayerPosition
	push.PlayerId = actionPlayer.PlayerInfo.PlayerId
	push.Countdown = sysSet.GameDelayTimeInt
	push.IsMustPlay = isMustPlay

}

func getCardType(cards []*Card) int32 {

	return 0
}
