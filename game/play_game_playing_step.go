package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/playerAction"
	"landlord/mconst/sysSet"
	"landlord/msg/mproto"
)

/*
出牌阶段
*/

func PlayingGame(room *Room, actionPlayerId string) {
	actionPlayer := room.Players[actionPlayerId]
	if actionPlayer == nil {
		logger.Error("房间里无此用户...!!!incredible")
		return
	}

	nextPosition := getNextPosition(actionPlayer.PlayerPosition)
	nextPlayer := getPlayerByPosition(room, nextPosition)

	lastPosition := getLastPosition(actionPlayer.PlayerPosition)
	lastPlayer := getPlayerByPosition(room, lastPosition)
	// 阻塞等待当前玩家的动作 超过系统设置时间后自动处理
	select {
	case action := <-actionPlayer.ActionChan:
		switch action.ActionType {
		case playerAction.OutCardAction: // 出牌
			OutCardsAction(room, actionPlayer, nextPlayer, action.ActionCards, action.CardsType)
		case playerAction.NotOutCardAction: // 不出
			NotOutCardsAction(room, actionPlayer, lastPlayer, nextPlayer)
		}
		//case <-time.After(time.Second * sysSet.GameDelayTime): // 自动不出
		//	NotOutCardsAction(room, actionPlayer, lastPlayer, nextPlayer)
	}
}

// 出牌逻辑 // 在接受消息前判断牌是否符合规则 和 是否能打过上家 这里不做处理
// 逻辑能到这一步  是确保能正常操作的
func OutCardsAction(room *Room, actionPlayer, nextPlayer *Player, cards []*Card, cardsType int32) {
	actionPlayer.DidAction = playerAction.OutCardAction
	actionPlayer.HandCards = removeCards(actionPlayer.HandCards, cards)
	room.EffectiveCard = cards
	room.EffectiveType = cardsType

	if len(actionPlayer.HandCards) == 0 {
		//pushOutCardHelp(room, nil, actionPlayer, playerAction.NotOutCardAction, false, cards, cardsType)
		pushLastOutCard(room, actionPlayer, cards, cardsType)
		logger.Debug("玩家胜利:", actionPlayer.PlayerInfo.PlayerId)
		// todo 当前玩家胜利  推送结算消息
		return
	}
	pushOutCardHelp(room, nextPlayer, actionPlayer, playerAction.NotOutCardAction, false, cards, cardsType)
	PlayingGame(room, nextPlayer.PlayerInfo.PlayerId)
}

// 不出逻辑
func NotOutCardsAction(room *Room, actionPlayer, lastPlayer, nextPlayer *Player, ) {
	actionPlayer.DidAction = playerAction.NotOutCardAction
	if lastPlayer.DidAction == playerAction.NotOutCardAction { // 如果上一个玩家不出 则又下一个玩家重新出牌
		reSetOutRoomToOut(room)
		pushMustOutCard(room, nextPlayer.PlayerInfo.PlayerId)
	} else { // 则由下一个玩家出牌
		pushOutCardHelp(room, nextPlayer, actionPlayer, playerAction.NotOutCardAction, false, nil, -3)
	}
	PlayingGame(room, nextPlayer.PlayerInfo.PlayerId)
}

// 设置房间重新出牌
// 及玩家 出的牌 下两家都不要
func reSetOutRoomToOut(room *Room) {
	room.EffectiveCard = nil
	// 重新设置玩家位无动作
	for _, v := range room.Players {
		v.DidAction = playerAction.NoAction
	}

}

// 推送必须出牌的消息 (第一次出牌开始 / 玩家出牌后 后面两个玩家都不要)
func pushMustOutCard(room *Room, playerId string) string {
	actionPlayer := getPlayerByPlayerId(room, playerId)
	getLastPosition(actionPlayer.PlayerPosition)
	lastPlayer := getPlayerByPosition(room, getLastPosition(actionPlayer.PlayerPosition))

	pushOutCardHelp(room, actionPlayer, lastPlayer, playerAction.NoAction, true, nil, 0)
	return actionPlayer.PlayerInfo.PlayerId
}

// 推送出牌辅助方法
func pushOutCardHelp(room *Room, actionPlayer, lastPlayer *Player, lastAction int32,
	isMustPlay bool, outCard []*Card, cardType int32) {
	var push mproto.PushOutCard
	if lastPlayer != nil {
		push.LastPlayerId = lastPlayer.PlayerInfo.PlayerId
		push.LastPlayerPosition = lastPlayer.PlayerPosition
		push.LastAction = lastAction
	}
	push.LastPlayerCards = ChangeCardToProto(outCard)
	push.LastPlayerCardsType = cardType

	if actionPlayer != nil {
		push.PlayerPosition = actionPlayer.PlayerPosition
		push.PlayerId = actionPlayer.PlayerInfo.PlayerId
		push.Countdown = sysSet.GameDelayTimeInt
		push.IsMustPlay = isMustPlay
	}

	bytes, _ := proto.Marshal(&push)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushOutCard, bytes))
}

// 推送玩家出的最后一首牌张牌
func pushLastOutCard(room *Room, actionPlayer *Player, lastCards []*Card, cardType int32) string {
	pushOutCardHelp(room, nil, actionPlayer, playerAction.NoAction, false, lastCards, cardType)
	return actionPlayer.PlayerInfo.PlayerId
}

// 将牌送 牌队中移除
func removeCards(cards, removeCards []*Card) []*Card {
	for i := 0; i < len(removeCards); i++ {
		cards = removeCard(cards, removeCards[i])
	}
	return cards
}

// 将牌送 牌队中移除
func removeCard(cards []*Card, removeCard *Card) []*Card {
	for i := 0; i < len(cards); i++ {
		if cards[i].Value == removeCard.Value && cards[i].Suit == removeCard.Suit {
			cards = append(cards[:i], cards[i+1:]...)
			break
		}
	}
	return cards
}
