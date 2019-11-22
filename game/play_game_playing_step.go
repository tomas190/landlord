package game

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/playerAction"
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
	// todo 用户托管动作

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

	if actionPlayer.IsLandlord {
		room.LandlordOutNum++
	}

	/*after*/
	after := actionPlayer.HandCards
	/*after*/

	actionPlayer.DidAction = playerAction.OutCardAction
	actionPlayer.HandCards = append([]*Card{}, removeCards(actionPlayer.HandCards, cards)...)
	actionPlayer.ThrowCards = append(actionPlayer.ThrowCards, cards[:]...)
	room.EffectiveCard = cards
	room.EffectiveType = cardsType

	// 出牌日志
	logger.Debug("出的牌:")
	PrintCard(cards)
	logger.Debug("出牌前:")
	PrintCard(actionPlayer.HandCards)
	logger.Debug("出牌后:")
	PrintCard(after)

	// 出牌日志

	if len(actionPlayer.HandCards) == 0 {
		//pushOutCardHelp(room, nil, actionPlayer, playerAction.NotOutCardAction, false, cards, cardsType)
		CheckSpring(room, actionPlayer)
		pushLastOutCard(room, actionPlayer, cards, cardsType)
		logger.Debug("玩家胜利:", actionPlayer.PlayerInfo.PlayerId)
		// 判断是否春天
		//
		Settlement(room, actionPlayer)
		return
	}
	setCurrentPlayerOut(room, nextPlayer.PlayerInfo.PlayerId, false)
	pushOutCardHelp(room, nextPlayer, actionPlayer, playerAction.NotOutCardAction, false, cards, cardsType)
	PlayingGame(room, nextPlayer.PlayerInfo.PlayerId)
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

// 不出逻辑
func NotOutCardsAction(room *Room, actionPlayer, lastPlayer, nextPlayer *Player, ) {
	actionPlayer.DidAction = playerAction.NotOutCardAction
	if lastPlayer.DidAction == playerAction.NotOutCardAction { // 如果上一个玩家不出 则又下一个玩家重新出牌
		reSetOutRoomToOut(room, nextPlayer.PlayerInfo.PlayerId)
		setCurrentPlayerOut(room, nextPlayer.PlayerInfo.PlayerId, true)
		pushMustOutCard(room, nextPlayer.PlayerInfo.PlayerId)
	} else { // 则由下一个玩家出牌
		setCurrentPlayerOut(room, nextPlayer.PlayerInfo.PlayerId, false)
		pushOutCardHelp(room, nextPlayer, actionPlayer, playerAction.NotOutCardAction, false, nil, -3)
	}
	PlayingGame(room, nextPlayer.PlayerInfo.PlayerId)
}

// 结算
// todo 最小金额计算 玩家只有这么多金币 则 只能输或者赢这么多
func Settlement(room *Room, winPlayer *Player) {
	// 1. 计算基本倍数

	mult := room.MultiAll
	settlementGold := room.RoomClass.BottomPoint * float64(mult)

	landPlayer, fp1, fp2 := getPlayerClass(room)
	roundId := fmt.Sprintf("room-%d-%d", room.RoomClass.RoomType, time.Now().Unix())

	var sPush mproto.PushSettlement

	// 如果赢家是地主
	if winPlayer.IsLandlord == true {
		var landRealWinGold float64 // 地主实际赢钱 税前
		if fp1.PlayerInfo.Gold < settlementGold { // 如果玩家1 的钱不够开
			landRealWinGold += fp1.PlayerInfo.Gold
			syncLossGold(fp1, fp1.PlayerInfo.Gold, roundId) // 同步金币 到中心服务 session

			showWinLossGold := fmt.Sprintf("-%.2f", fp1.PlayerInfo.Gold)
			ss := getSelfSettlement(room, fp1, -1, showWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, ss)
		} else {
			landRealWinGold += settlementGold
			syncLossGold(fp1, settlementGold, roundId) // 同步金币 到中心服务 session

			showWinLossGold := fmt.Sprintf("-%.2f", settlementGold)
			ss := getSelfSettlement(room, fp1, -1, showWinLossGold, false)
			sPush.Settlement = append(sPush.Settlement, ss)
		}

		if fp2.PlayerInfo.Gold < settlementGold { // 如果玩家2 的钱不够开
			landRealWinGold += fp2.PlayerInfo.Gold
			syncLossGold(fp1, fp2.PlayerInfo.Gold, roundId) // 同步金币 到中心服务 session

			showWinLossGold := fmt.Sprintf("-%.2f", fp2.PlayerInfo.Gold)
			ss := getSelfSettlement(room, fp2, -1, showWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, ss)
		} else {
			landRealWinGold += settlementGold
			syncLossGold(fp1, settlementGold, roundId) // 同步金币 到中心服务 session

			showWinLossGold := fmt.Sprintf("-%.2f", settlementGold)
			ss := getSelfSettlement(room, fp2, -1, showWinLossGold, false)
			sPush.Settlement = append(sPush.Settlement, ss)
		}

		landRealWinGoldPay := landRealWinGold * Server.GameTaxRate            // 地主实际赢钱 税后
		syncWinGold(landPlayer, landRealWinGold, landRealWinGoldPay, roundId) // 同步金币 到中心服务 session

		showWinLossGold := fmt.Sprintf("%.2f", landRealWinGoldPay)
		ss := getSelfSettlement(room, landPlayer, 1, showWinLossGold, false)
		sPush.Settlement = append(sPush.Settlement, ss)

	} else { // 如果玩家不是地主
		// 1. 判断地主金币是否够开
		if landPlayer.PlayerInfo.Gold/2 < settlementGold {

			farmerRealWinGold := landPlayer.PlayerInfo.Gold / 2
			farmerRealWinGoldPay := farmerRealWinGold * Server.GameTaxRate

			syncWinGold(fp1, settlementGold, farmerRealWinGoldPay, roundId)
			syncWinGold(fp2, settlementGold, farmerRealWinGoldPay, roundId)
			syncLossGold(landPlayer, landPlayer.PlayerInfo.Gold, roundId)
			//
			logger.Debug("地主玩家输钱不够开", landPlayer.PlayerInfo.Gold)
			logger.Debug("结算金额基*1", settlementGold)
			logger.Debug("结算金额基*2", settlementGold*2)

			showWinLossGold := fmt.Sprintf("%.2f", farmerRealWinGoldPay)
			fs1 := getSelfSettlement(room, fp1, 1, showWinLossGold, false)
			sPush.Settlement = append(sPush.Settlement, fs1)

			fs2 := getSelfSettlement(room, fp2, 1, showWinLossGold, false)
			sPush.Settlement = append(sPush.Settlement, fs2)

			landShowWinLossGold := fmt.Sprintf("-%.2f", landPlayer.PlayerInfo.Gold)
			ls := getSelfSettlement(room, landPlayer, -1, landShowWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, ls)

		} else {
			// 正常结算
			winGoldPay := settlementGold * Server.GameTaxRate
			syncWinGold(fp1, settlementGold, winGoldPay, roundId)
			syncWinGold(fp2, settlementGold, winGoldPay, roundId)
			syncLossGold(landPlayer, settlementGold*2, roundId)

			showWinLossGold := fmt.Sprintf("%.2f", winGoldPay)
			fs1 := getSelfSettlement(room, fp1, 1, showWinLossGold, false)
			sPush.Settlement = append(sPush.Settlement, fs1)

			fs2 := getSelfSettlement(room, fp2, 1, showWinLossGold, false)
			sPush.Settlement = append(sPush.Settlement, fs2)

			landShowWinLossGold := fmt.Sprintf("-%.2f", settlementGold*2)
			ls := getSelfSettlement(room, landPlayer, -1, landShowWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, ls)
		}

	}

	var mulitInfo mproto.MultipleInfo
	mulitInfo.FightLandlord = fmt.Sprintf("×%d", room.MultiBoom)
	mulitInfo.Boom = fmt.Sprintf("×%d", room.MultiBoom)
	mulitInfo.Spring = fmt.Sprintf("×%d", room.MultiSpring)
	sPush.MultipleInfo = &mulitInfo
	sPush.WaitTime = sysSet.GameDelayReadyTimeInt

	bytes, _ := proto.Marshal(&sPush)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushSettlement, bytes))

}

func syncWinGold(player *Player, gold, goldPay float64, roundId string) float64 {
	orderId := fmt.Sprintf("%s-%s-win", roundId, player.PlayerInfo.PlayerId)
	player.PlayerInfo.Gold += goldPay              // 同步到房间id
	err := SetSessionGold(player.Session, goldPay) // 同步到session
	if err != nil {
		logger.Error("同步进步到session失败: !!!incredible")
	}
	UserSyncWinScore(player.PlayerInfo.PlayerId, gold, roundId, orderId) // 同步到中心服务
	return player.PlayerInfo.Gold
}

func syncLossGold(player *Player, gold float64, roundId string) float64 {
	orderId := fmt.Sprintf("%s-%s-loss", roundId, player.PlayerInfo.PlayerId)
	player.PlayerInfo.Gold -= gold
	err := SetSessionGold(player.Session, -gold) // 同步到session
	if err != nil {
		logger.Error("同步进步到session失败: !!!incredible")
	}
	UserSyncLoseScore(player.PlayerInfo.PlayerId, -gold, roundId, orderId)
	return player.PlayerInfo.Gold
}

func getSelfSettlement(room *Room, player *Player, winOrFail int32, winOrLossGold string, isMinSettlement bool) *mproto.Settlement {
	var result mproto.Settlement

	if player.IsLandlord {
		result.IsLandlord = 1
		result.Multiple = room.MultiAll * 2
	} else {
		result.IsLandlord = -1
		result.Multiple = room.MultiAll
	}
	result.PlayerId = player.PlayerInfo.PlayerId
	result.Position = player.PlayerPosition
	result.CurrentGold = player.PlayerInfo.Gold
	result.PlayerName = player.PlayerInfo.Name
	result.WinOrFail = winOrFail
	result.WinLossGold = winOrLossGold
	result.RemainCards = ChangeCardToProto(player.HandCards)
	result.MinSettlement = isMinSettlement
	return &result
}

/*
	返回第一个地主玩家
	后面两个农民玩家
*/
func getPlayerClass(room *Room) (*Player, *Player, *Player) {
	var landPlayer *Player
	var fplayers []*Player
	for _, p := range room.Players {
		if p.IsLandlord == true {
			landPlayer = p
		} else {
			fplayers = append(fplayers, p)
		}
	}

	if len(fplayers) != 2 || fplayers == nil {
		logger.Error("分类玩家失败: !!!incredible")
		return nil, nil, nil
	}
	return landPlayer, fplayers[0], fplayers[1]
}

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
		v.DidAction = playerAction.NoAction
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
