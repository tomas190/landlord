package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
	"landlord/mconst/cardConst"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/playerAction"
	"landlord/mconst/playerStatus"
	"landlord/mconst/sysSet"
	"landlord/msg/mproto"
	"time"
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

	// todo 用户托管动作
	if actionPlayer.IsGameHosting {
		DoGameHosting(room, actionPlayer, nextPlayer, lastPlayer)
		return
	}

	select {
	case action := <-actionPlayer.ActionChan:
		switch action.ActionType {
		case playerAction.OutCardAction: // 出牌
			OutCardsAction(room, actionPlayer, nextPlayer, action.ActionCards, action.CardsType)
		case playerAction.NotOutCardAction: // 不出
			NotOutCardsAction(room, actionPlayer, lastPlayer, nextPlayer)
		}
	case <-time.After(time.Second * sysSet.GameDelayTime): // 自动不出
		// todo 进入托管
		actionPlayer.IsGameHosting = true
		RespGameHosting(room, playerStatus.GameHosting, actionPlayer.PlayerPosition, actionPlayer.PlayerInfo.PlayerId)
		DoGameHosting(room, actionPlayer, nextPlayer, lastPlayer) // 走托管逻辑
	}
}

// 出牌逻辑 // 在接受消息前判断牌是否符合规则 和 是否能打过上家 这里不做处理
// 逻辑能到这一步  是确保能正常操作的
func OutCardsAction(room *Room, actionPlayer, nextPlayer *Player, cards []*Card, cardsType int32) {

	if actionPlayer.IsLandlord {
		room.LandlordOutNum++
	}

	room.ThrowCards = append(room.ThrowCards, cards...)
	/*after*/
	after := actionPlayer.HandCards
	/*after*/
	actionPlayer.LastOutCard = cards
	actionPlayer.LastAction = playerAction.OutCardAction
	actionPlayer.HandCards = append([]*Card{}, removeCards(actionPlayer.HandCards, cards)...)
	actionPlayer.ThrowCards = append(actionPlayer.ThrowCards, cards[:]...)
	room.EffectiveCard = cards
	room.EffectiveType = cardsType

	// 出牌日志
	logger.Debug("出的牌:")
	PrintCard(cards)
	logger.Debug("出牌前:")
	PrintCard(after)
	logger.Debug("出牌后:")
	PrintCard(actionPlayer.HandCards)

	// 出牌日志

	if len(actionPlayer.HandCards) == 0 {
		//pushOutCardHelp(room, nil, actionPlayer, playerAction.NotOutCardAction, false, cards, cardsType)
		// 判断是否春天
		CheckSpring(room, actionPlayer)
		pushLastOutCard(room, actionPlayer, cards, cardsType)
		logger.Debug("玩家胜利:", actionPlayer.PlayerInfo.PlayerId)

		// 结算
		Settlement(room, actionPlayer)

		// 移除房间
		clearRoomAndPlayer(room)
		return
	}
	setCurrentPlayerOut(room, nextPlayer.PlayerInfo.PlayerId, false)
	pushOutCardHelp(room, nextPlayer, actionPlayer, playerAction.NotOutCardAction, false, cards, cardsType)
	// 推送记牌器
	pushCardCount(room)
	PlayingGame(room, nextPlayer.PlayerInfo.PlayerId)
}

// 不出逻辑
func NotOutCardsAction(room *Room, actionPlayer, lastPlayer, nextPlayer *Player, ) {
	actionPlayer.LastAction = playerAction.NotOutCardAction
	if lastPlayer.LastAction == playerAction.NotOutCardAction { // 如果上一个玩家不出 则又下一个玩家重新出牌
		reSetOutRoomToOut(room, nextPlayer.PlayerInfo.PlayerId)
		setCurrentPlayerOut(room, nextPlayer.PlayerInfo.PlayerId, true)
		pushMustOutCard(room, nextPlayer.PlayerInfo.PlayerId)
	} else { // 则由下一个玩家出牌
		setCurrentPlayerOut(room, nextPlayer.PlayerInfo.PlayerId, false)
		pushOutCardHelp(room, nextPlayer, actionPlayer, playerAction.NotOutCardAction, false, nil, -3)
	}
	PlayingGame(room, nextPlayer.PlayerInfo.PlayerId)
}

// todo 这里要修改  能出必出
func DoGameHosting(room *Room, actionPlayer, nextPlayer, lastPlayer *Player) {
	if actionPlayer.IsMustDo {
		SortCard(actionPlayer.HandCards)
		// todo最后一张
		OutCardsAction(room, actionPlayer, nextPlayer, actionPlayer.HandCards[len(actionPlayer.HandCards)-1:], cardConst.CARD_PATTERN_SINGLE)
	} else {
		// 自动不出
		NotOutCardsAction(room, actionPlayer, lastPlayer, nextPlayer)
	}
}

/*================== help func ===============*/

// 设置房间重新出牌
// 及玩家 出的牌 下两家都不要
func reSetOutRoomToOut(room *Room, playerId string) {
	room.EffectiveCard = nil
	// 重新设置玩家位无动作
	for _, v := range room.Players {
		if v.PlayerInfo.PlayerId == playerId {
			v.IsMustDo = true
			v.IsCanDo = true
		} else {
			v.IsMustDo = false
			v.IsCanDo = false
		}
		//v.LastAction = playerAction.NoAction
	}

}

// 设置当前玩家出牌
func setCurrentPlayerOut(room *Room, playerId string, isMustDo bool) {
	for _, v := range room.Players {
		if v.PlayerInfo.PlayerId == playerId {
			v.IsCanDo = true
			v.IsMustDo = isMustDo
		} else {
			v.IsCanDo = false
			v.IsMustDo = false
		}
	}

}

// 推送必须出牌的消息 (第一次出牌开始 / 玩家出牌后 后面两个玩家都不要)
func pushMustOutCard(room *Room, playerId string) string {
	actionPlayer := getPlayerByPlayerId(room, playerId)
	lastPlayer := getPlayerByPosition(room, getLastPosition(actionPlayer.PlayerPosition))

	pushOutCardHelp(room, actionPlayer, lastPlayer, playerAction.NoAction, true, nil, 0)
	return actionPlayer.PlayerInfo.PlayerId
}

// 推送出牌辅助方法
func pushOutCardHelp(room *Room, actionPlayer, lastPlayer *Player, lastAction int32,
	isMustPlay bool, outCard []*Card, cardType int32) {
	var push mproto.PushOutCard

	push.LastPlayerId = lastPlayer.PlayerInfo.PlayerId
	push.LastPlayerPosition = lastPlayer.PlayerPosition
	push.LastAction = lastAction

	push.LastPlayerCards = ChangeCardToProto(outCard)
	push.LastPlayerCardsType = cardType
	push.LastRemainLen = int32(len(lastPlayer.HandCards))

	if actionPlayer != nil {
		push.PlayerPosition = actionPlayer.PlayerPosition
		push.PlayerId = actionPlayer.PlayerInfo.PlayerId
		push.Countdown = sysSet.GameDelayTimeInt
		push.IsMustPlay = isMustPlay
	}
	push.Multi = room.MultiAll
	bytes, _ := proto.Marshal(&push)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushOutCard, bytes))
}

// 推送玩家出的最后一首牌张牌
func pushLastOutCard(room *Room, actionPlayer *Player, lastCards []*Card, cardType int32) {
	pushOutCardHelp(room, nil, actionPlayer, playerAction.NoAction, false, lastCards, cardType)
}

// 将牌送 牌队中移除
func removeCards(cards, removeCards []*Card) []*Card {
	newCard := append([]*Card{}, cards...)
	for i := 0; i < len(removeCards); i++ {
		newCard = removeCard(newCard, removeCards[i])
		newCard = append([]*Card{}, newCard...)
	}
	return newCard
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

// 判断是否春天
func CheckSpring(room *Room, player *Player) {
	// 1. 判断玩家是否地主
	_, f1, f2 := getPlayerClass(room)
	if player.IsLandlord == true {
		farmerThrows := append(f1.ThrowCards, f2.ThrowCards...)
		if len(farmerThrows) <= 0 || farmerThrows == nil {
			logger.Debug("================ 地主春天 =================")
			room.MultiAll = room.MultiAll * 2
			room.MultiSpring = 2
		}
	} else {
		if room.LandlordOutNum == 1 {
			logger.Debug("================ 农民春天 =================")
			room.MultiAll = room.MultiAll * 2
			room.MultiSpring = 2
		}
	}
}

/*
	返回第一个地主玩家
	后面两个农民玩家
*/
func getPlayerClass(room *Room) (*Player, *Player, *Player) {
	var landPlayer *Player
	var fPlayers []*Player
	for _, p := range room.Players {
		if p.IsLandlord == true {
			landPlayer = p
		} else {
			fPlayers = append(fPlayers, p)
		}
	}

	if len(fPlayers) != 2 || fPlayers == nil {
		logger.Error("分类玩家失败: !!!incredible")
		return nil, nil, nil
	}
	return landPlayer, fPlayers[0], fPlayers[1]
}

// 清空房间 和用户 的session
func clearRoomAndPlayer(room *Room) {
	players := room.Players
	for _, player := range players {
		if player.IsCloseSession { // 如果玩家已经断线 登出中心服
			ClearClosePlayer(player.Session)
		} else {
			// 置空玩家的roomId
			SetSessionRoomId(player.Session, "")
		}
	}
	// 移除房间
	RemoveRoom(room.RoomId)
}

func ClearClosePlayer(session *melody.Session) {
	playerInfo, err := GetSessionPlayerInfo(session)
	if err != nil {
		logger.Debug("无用户session信息:", err.Error())
		return
	}
	password := GetSessionPassword(session)
	// 登出中心服
	UserLogoutCenter(playerInfo.PlayerId, password)
	RemoveAgent(playerInfo.PlayerId)
}
