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
出牌阶段 特殊延时时间
*/
func getPlayingDelayTime(room *Room, actionPlayer *Player) (time.Duration, int32) {
	t := sysSet.GameDelayTime
	ti := sysSet.GameDelayTimeInt
	if actionPlayer.IsMustDo {
		return t, ti
	}

	if room.EffectiveType == cardConst.CARD_PATTERN_ROCKET {
		logger.Debug("满足上手是炸弹 延时3秒...")
		t = 3
		ti = 3
	}
	if len(room.EffectiveCard) > 1 && len(actionPlayer.HandCards) == 1 {
		logger.Debug("满足手牌只有一张 但上手是非单张 延时3秒...")
		t = 3
		ti = 3
	}
	return t, ti
}

func PlayingGame(room *Room, actionPlayerId string) {
	//checkRoom := GetRoom(room.RoomId)
	//if checkRoom == nil {
	//	runtime.Goexit()
	//	return
	//}

	actionPlayer := room.Players[actionPlayerId]
	if actionPlayer == nil {
		logger.Error("房间里无此用户...!!!incredible")
		//return
	}

	// todo 每秒记录玩家的时间点用户 玩家再次阶段退出后 再次进入房间
	//uptWtChin := make(chan struct{})

	// 2020年2月25日15:21:31
	delayTime, delayTimeInt := getPlayingDelayTime(room, actionPlayer)
	// 2020年2月25日15:21:31
	//go updatePlayerWaitingTime(actionPlayer, uptWtChin, sysSet.GameDelayTimeInt)
	go updatePlayingWaitingTime(actionPlayer, delayTimeInt)
	// todo 每秒记录玩家的时间点用户 玩家再次阶段退出后 再次进入房间

	nextPosition := getNextPosition(actionPlayer.PlayerPosition)
	nextPlayer := getPlayerByPosition(room, nextPosition)

	lastPosition := getLastPosition(actionPlayer.PlayerPosition)
	lastPlayer := getPlayerByPosition(room, lastPosition)
	// 阻塞等待当前玩家的动作 超过系统设置时间后自动处理

	// todo 用户托管动作
	if actionPlayer.IsGameHosting {
		DoGameHosting(room, actionPlayer, nextPlayer, lastPlayer)
		// todo 如果机器人假装断线托管 根据70%的几率恢复
		return
	}

	if actionPlayer.IsRobot {
		actionPlayer.GroupCard = GroupHandsCard(actionPlayer.HandCards)
		RobotPlayAction(room, actionPlayer, nextPlayer, lastPlayer)
		return
	}

	oNum := room.OutNum
	go func(delayTime time.Duration, outNum int32, acPlayer *Player, r *Room) {
		DelaySomeTime(delayTime)
		if outNum == r.OutNum && acPlayer.IsCanDo {
			if delayTime != 3 {
				// 如果是不能出的就不托管
				acPlayer.IsGameHosting = true
				RespGameHosting(r, playerStatus.GameHosting, acPlayer.PlayerPosition, acPlayer.PlayerInfo.PlayerId)
			}

			if acPlayer.IsMustDo {
				//DoGameHosting(room, acPlayer, nextPlayer, lastPlayer) // 走托管逻辑
				cards, cType := FindMustBeOutCards(acPlayer.HandCards)
				go OutCardsAction(r, acPlayer, nextPlayer, cards, cType)
			} else {
				go NotOutCardsAction(r, acPlayer, lastPlayer, nextPlayer) // 走不出逻辑
			}
		}
	}(delayTime, oNum, actionPlayer, room)

	//actionPlayer.ActionChan = make(chan PlayerActionChan)
	//select {
	//case action := <-actionPlayer.ActionChan:
	//	go func() {
	//		uptWtChin <- struct{}{}
	//	}()
	//	switch action.ActionType {
	//	case playerAction.OutCardAction: // 出牌
	//		OutCardsAction(room, actionPlayer, nextPlayer, action.ActionCards, action.CardsType)
	//	case playerAction.NotOutCardAction: // 不出
	//		NotOutCardsAction(room, actionPlayer, lastPlayer, nextPlayer)
	//	}
	////case <-time.After(time.Second * sysSet.GameDelayTime): // 自动不出
	//case <-time.After(time.Second * delayTime): // 自动不出
	//	go func() {
	//		uptWtChin <- struct{}{}
	//	}()
	//	// todo 进入托管
	//	if delayTime != 3 {
	//		// 如果是不能出的就不托管
	//		actionPlayer.IsGameHosting = true
	//		RespGameHosting(room, playerStatus.GameHosting, actionPlayer.PlayerPosition, actionPlayer.PlayerInfo.PlayerId)
	//	}
	//
	//	if actionPlayer.IsMustDo {
	//		//DoGameHosting(room, actionPlayer, nextPlayer, lastPlayer) // 走托管逻辑
	//		cards, cType := FindMustBeOutCards(actionPlayer.HandCards)
	//		OutCardsAction(room, actionPlayer, nextPlayer, cards, cType)
	//	} else {
	//		NotOutCardsAction(room, actionPlayer, lastPlayer, nextPlayer) // 走不出逻辑
	//	}
	//}
}

// 出牌逻辑 // 在接受消息前判断牌是否符合规则 和 是否能打过上家 这里不做处理
// 逻辑能到这一步  是确保能正常操作的
func OutCardsAction(room *Room, actionPlayer, nextPlayer *Player, cards []*Card, cardsType int32) {

	// 机器人出牌补丁
	if actionPlayer.IsRobot {
		// 机器人出牌全面检测 是否符合出牌规则
		if actionPlayer.IsMustDo {
			logger.Debug("机器人首出牌检测...")
			cards, cardsType = OutCardCheck(cards, cardsType)
		} else {
			// 当前牌是否大过上家 打不过就不出
			logger.Debug("机器人跟牌检测...")
			eCard := room.EffectiveCard
			can := CanBeat(eCard, cards, )
			if !can {
				logger.Debug("检测不通过...重新跟牌")
				logger.Debug("上手牌:", room.EffectiveCard)
				logger.Debug("非法牌:", cards)
				logger.Debug("非法牌型:", cardsType)
				beatCards, b, bcType := FindCanBeatCards(actionPlayer.HandCards, room.EffectiveCard, room.EffectiveType)
				if !b {
					// 检测出没有打过的牌 则不出
					lastPosition := getLastPosition(actionPlayer.PlayerPosition)
					lastPlayer := getPlayerByPosition(room, lastPosition)
					go NotOutCardsAction(room, actionPlayer, lastPlayer, nextPlayer)
					return
				}
				cards = append([]*Card{}, beatCards...)
				cardsType = bcType
				logger.Debug("修正牌:", cards)
				logger.Debug("修正牌牌型:", cardsType)
			}
			logger.Debug("检测通过...")
		}
	}

	if actionPlayer.IsLandlord {
		room.LandlordOutNum++
	}
	room.OutNum++

	// 炸弹计算翻倍
	if cardsType == cardConst.CARD_PATTERN_BOMB {
		room.MultiAll = room.MultiAll * 2
		if room.MultiBoom == 0 {
			room.MultiBoom = 2
		} else {
			room.MultiBoom = room.MultiBoom * 2
		}
	}
	// 火箭倍数
	if cardsType == cardConst.CARD_PATTERN_ROCKET {
		room.MultiAll = room.MultiAll * 2
		room.MultiRocket = 2
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
	room.EffectivePlayerId = actionPlayer.PlayerInfo.PlayerId

	// 出牌日志
	logger.Debug("出的牌:")
	PrintCard(cards)
	logger.Debug("出牌前:")
	PrintCard(after)
	logger.Debug("出牌后:")
	PrintCard(actionPlayer.HandCards)

	// 出牌日志

	if len(actionPlayer.HandCards) == 0 {
		// 有人打完代表这局结束
		//pushOutCardHelp(room, nil, actionPlayer, playerAction.NotOutCardAction, false, cards, cardsType)
		// 判断是否春天
		isSpring := CheckSpring(room, actionPlayer)
		pushLastOutCard(room, actionPlayer, cards, cardsType)
		logger.Debug("玩家胜利:", actionPlayer.PlayerInfo.PlayerId)
		pushCardCount(room)
		if isSpring {
			//DelaySomeTime(1)
			pushSpring(room)
			//DelaySomeTime(1)
		}
		//
		// 结算
		Settlement(room, actionPlayer)

		//
		//if isSpring {
		//	DelaySomeTime(1)
		//}
		//DelaySomeTime(2)
		//// 移除房间
		clearRoomAndPlayer(room)
		//runtime.Goexit()
		return
	}
	if actionPlayer.IsRobot { // 重新组排
		actionPlayer.GroupCard = GroupHandsCard(actionPlayer.HandCards)
	}
	setCurrentPlayerOut(room, nextPlayer.PlayerInfo.PlayerId, false)
	_, delayTimeInt := getPlayingDelayTime(room, nextPlayer)
	pushOutCardHelp(room, nextPlayer, actionPlayer, playerAction.NotOutCardAction, false, cards, cardsType, delayTimeInt)
	// 推送记牌器
	pushCardCount(room)
	nextPlayer.IsCanDo = true
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
		_, delayTimeInt := getPlayingDelayTime(room, nextPlayer)
		pushOutCardHelp(room, nextPlayer, actionPlayer, playerAction.NotOutCardAction, false, nil, -3, delayTimeInt)
	}
	room.OutNum++
	PlayingGame(room, nextPlayer.PlayerInfo.PlayerId)
}

// 托管操作
func DoGameHosting(room *Room, actionPlayer, nextPlayer, lastPlayer *Player) {
	DelaySomeTime(1)
	if actionPlayer.IsMustDo {
		// 取牌
		cards, cType := FindMustBeOutCards(actionPlayer.HandCards)
		OutCardsAction(room, actionPlayer, nextPlayer, cards, cType)
	} else if bCards, b, bType := FindCanBeatCards(actionPlayer.HandCards, room.EffectiveCard, room.EffectiveType); b {
		//  判断出上家的牌型 如果有能大过上家的牌 则出没有则不出
		OutCardsAction(room, actionPlayer, nextPlayer, bCards, bType)
	} else {
		//  判断出上家的牌型 如果有能大过上家的牌 则出没有则不出
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
		v.WaitingTime = sysSet.GameDelayTimeInt
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

// 推送确定完必须出牌的消息 (第一次出牌开始 / 玩家出牌后 后面两个玩家都不要)
func pushFirstMustOutCard(room *Room, playerId string) string {
	actionPlayer := getPlayerByPlayerId(room, playerId)
	var push mproto.PushOutCard
	push.PlayerPosition = actionPlayer.PlayerPosition
	push.PlayerId = actionPlayer.PlayerInfo.PlayerId
	push.Countdown = sysSet.GameDelayTimeInt
	push.IsMustPlay = true
	push.Multi = room.MultiAll
	bytes, _ := proto.Marshal(&push)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushOutCard, bytes))
	return actionPlayer.PlayerInfo.PlayerId
}

// 推送必须出牌的消息 (第一次出牌开始 / 玩家出牌后 后面两个玩家都不要)
func pushMustOutCard(room *Room, playerId string) string {
	actionPlayer := getPlayerByPlayerId(room, playerId)
	lastPlayer := getPlayerByPosition(room, getLastPosition(actionPlayer.PlayerPosition))

	//pushOutCardHelp(room, actionPlayer, lastPlayer, playerAction.NoAction, true, nil, 0)
	pushOutCardHelp(room, actionPlayer, lastPlayer, lastPlayer.LastAction, true, nil, 0, sysSet.GameDelayTimeInt)
	return actionPlayer.PlayerInfo.PlayerId
}

// 推送出牌辅助方法
func pushOutCardHelp(room *Room, actionPlayer, lastPlayer *Player, lastAction int32,
	isMustPlay bool, outCard []*Card, cardType, delayTime int32) {
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
		//push.Countdown = sysSet.GameDelayTimeInt
		push.Countdown = delayTime
		push.IsMustPlay = isMustPlay
	}
	push.Multi = room.MultiAll
	bytes, _ := proto.Marshal(&push)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushOutCard, bytes))
}

// 推送玩家出的最后一首牌张牌
func pushLastOutCard(room *Room, actionPlayer *Player, lastCards []*Card, cardType int32) {
	pushOutCardHelp(room, nil, actionPlayer, actionPlayer.LastAction, false, lastCards, cardType, sysSet.GameDelayTimeInt)
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
func CheckSpring(room *Room, player *Player) bool {
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

	if room.MultiSpring == 2 {
		return true
	}
	return false
}

// 清空房间 和用户 的session
func clearRoomAndPlayer(room *Room) {
	players := room.Players
	for _, player := range players {
		if !player.IsRobot {
			if player.IsCloseSession { // 如果玩家已经断线 登出中心服
				ClearClosePlayer(player.Session)
			}
			//else {
			// 置空玩家的roomId
			//SetSessionRoomId(player.Session, "")
			//}
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
	RemoveWaitUser(playerInfo.PlayerId)
}

// 出牌前监测牌型是否正确
func OutCardCheck(outCard []*Card, cardType int32) ([]*Card, int32) {

	// 顺子监测补丁
	if cardType == cardConst.CARD_PATTERN_SEQUENCE {
		cardsType := GetCardsType(outCard)
		if cardsType == cardType {
			return outCard, cardType
		} else {
			// 组成最小的顺子
			gc := GroupHandsCard(outCard)
			gc = completeGroupCard(gc)
			if len(gc.Junko) >= 1 {
				logger.Debug("非法出牌 补丁已经修正 非法牌:")
				PrintCard(outCard)
				logger.Debug("非法出牌 补丁已经修正 修正:")
				PrintCard(gc.Junko[0].Card)
				return gc.Junko[0].Card, cardType
			}
			// 如果这里顺子也组不了就返回第一张单张
			if cardsType == cardConst.CARD_PATTERN_ERROR {
				return append([]*Card{}, outCard[0]), cardConst.CARD_PATTERN_SINGLE
			}
		}
	}

	getCardsType := GetCardsType(outCard)
	if getCardsType >= cardConst.CARD_PATTERN_SINGLE && getCardsType <= cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS {
		return outCard, cardType
	}

	return FindMustBeOutCards(outCard)

	//return outCard, cardType
}
