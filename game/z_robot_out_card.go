package game

import "C"
import (
	"github.com/wonderivan/logger"
	"landlord/mconst/cardConst"
)

// 机器人出牌
func robotOutCard(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	if robot.IsLandlord {
		if robot.IsMustDo {
			// todo 地主机器人首出  注意炸弹顺序
			//landlordRobotOutCardMustDo(room, robot, nextPlayer, lastPlayer)
			NewLandlordRobotOutCardMustDo(room, robot, nextPlayer, lastPlayer)

		} else {
			// todo // 地主机器人跟牌 注意农民保单的时候 要顶牌
			//landlordRobotFallowCard(room, robot, nextPlayer, lastPlayer)
			NewLandlordRobotFallowCard(room, robot, nextPlayer, lastPlayer)
		}
	} else {
		// todo
		// 农民首出
		if robot.IsMustDo {
			farmerRobotOutCardMustDo(room, robot, nextPlayer, lastPlayer)
			//landlordRobotOutCardMustDo(room, robot, nextPlayer, lastPlayer)
		} else {
			// 农民玩家跟牌
			farmerRobotFallowCard(room, robot, nextPlayer, lastPlayer)
		}
	}

}

/*
	农民玩家首出
   F1 ->F2 ->landlord
   ALL
*/
func farmerRobotOutCardMustDo(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	efficId := room.EffectivePlayerId
	efficP := room.Players[efficId]
	if efficP.IsLandlord && !nextPlayer.IsLandlord {
		// 农民下家是农民出首牌策略
		//farmerRobotOutCardMustDoF1(room, robot, nextPlayer, lastPlayer)
		NewRobotFarmerMustDoF1(room, robot, nextPlayer, lastPlayer)
	} else {
		// 农民下家是地主出首牌策略
		//farmerRobotOutCardMustDoF2(room, robot, nextPlayer, lastPlayer)
		NewRobotFarmerMustDoF2(room, robot, nextPlayer, lastPlayer)
	}
}


/*s 地主出牌策略 还有很大的优化空间*/
// 一 地主机器人首出牌策虐
/*
判断地主当前几首牌,则判断手牌中的天牌数量有几首
i.	总首数-天牌数==1  ：随机出弹外的天牌 最后一首出 普通手牌
ii.	查看是否有对手玩家保单或者双
没有：则权值最小的牌: 如果权值最小的牌为单牌或者对子：则看有没有三代 或者 飞机 如果有三代或者飞机 则让飞机三带带上这个牌出牌
有：则先去除该类型的出牌，的权值最小牌
			有：但是手牌中只有该类型的牌 则取权重倒数第二张

*/
func landlordRobotOutCardMustDo(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	logger.Debug("地主首出")
	// 1.总首数-天牌数<=1  ：随机出弹外的天牌 最后一首出 普通手牌
	completeGroups := completeGroupCard(robot.GroupCard)
	completeReCards := changeGroupToReCard(completeGroups)
	outNums := len(completeReCards) // 出牌总首数
	godCard, godCardNum := CheckGodCard(completeReCards, nextPlayer.HandCards, lastPlayer.HandCards)

	// 如果最后两手 有一首是炸弹 并且 无人保单或者无人报双的情况下 先走非炸弹的
	cType1, b1 := checkPlayerHasLast(nextPlayer)
	cType2, b2 := checkPlayerHasLast(lastPlayer)
	if outNums == 2 {
		if !b1 && !b2 {
			if completeReCards[0].CardType == cardConst.CARD_PATTERN_BOMB ||
				completeReCards[1].CardType == cardConst.CARD_PATTERN_ROCKET {
				for i := 0; i < len(godCard); i++ {
					if completeReCards[i].CardType != cardConst.CARD_PATTERN_BOMB {
						logger.Debug("地主首出 最后一首是炸弹+其他 并无人报牌 显出其他")
						outCard := completeReCards[i].Card
						cardType := completeReCards[i].CardType
						OutCardsAction(room, robot, nextPlayer, outCard, cardType)
						return
					}
				}
			}
		}
	}

	if outNums-godCardNum <= 1 { // 这里如果总首数-天牌数<=1 // 则先出天牌
		for i := 0; i < len(godCard); i++ {
			if godCard[i].IsGodCard {
				outCard := godCard[i].RC.Card
				cardType := godCard[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, outCard, cardType)
				return
			}
		}
	}

	// 判断第一位是否最后一手牌
	if b1 {
		logger.Debug("地主首出 下一个有人最后一手牌", cType1)
		for i := 0; i < len(godCard); i++ {
			if godCard[i].RC.CardType != cType1 {
				outCard := godCard[i].RC.Card
				cardType := godCard[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, outCard, cardType)
				return
			}
		}
		if len(completeReCards) >= 2 {
			logger.Debug("地主首出 只有和这种类型相似的牌  出第二大")
			OutCardsAction(room, robot, nextPlayer, completeReCards[1].Card, completeReCards[1].CardType)
			return
		}
	}

	// 判断第一位是否最后一手牌
	if b2 {
		logger.Debug("地主首出 上一个玩家有最后一手牌", cType2)
		for i := 0; i < len(godCard); i++ {
			if godCard[i].RC.CardType != cType2 {
				outCard := godCard[i].RC.Card
				cardType := godCard[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, outCard, cardType)
				return
			}
		}
		if len(completeReCards) >= 2 {
			logger.Debug("地主首出 只有和这种类型相似的牌  出第二大")
			OutCardsAction(room, robot, nextPlayer, completeReCards[1].Card, completeReCards[1].CardType)
			return
		}
	}

	SortReCardByWightSL(completeReCards)
	logger.Debug("地主首出 无人最后一手牌 取最小权值牌")
	OutCardsAction(room, robot, nextPlayer, completeReCards[0].Card, completeReCards[0].CardType)
}

// 二 地主机器人跟牌策虐
/*
地主跟牌：//
	先从组牌中取出最小跟牌，
	1.有：出牌
	2.没有 在寻找有出必出
	2.1 有：不是炸弹:
			则重新组牌看下首数是否多过先前组牌 如果多出2 不出
		2.1 有：则判断是否炸弹
			2.11是 判断是否玩家最后一首
			2.111 是炸
			2.112 不出(隐忍一波)
		2.2 没有：不出

*/
func landlordRobotFallowCard(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	eType := room.EffectiveType
	eCards := room.EffectiveCard
	ePlayer := room.Players[room.EffectivePlayerId]
	var otherFarmer *Player
	for _, v := range room.Players {
		if v.PlayerInfo.PlayerId != ePlayer.PlayerInfo.PlayerId &&
			v.PlayerInfo.PlayerId != robot.PlayerInfo.PlayerId {
			otherFarmer = v
		}
	}
	if otherFarmer == nil {
		logger.Debug("！！！ !impossible 不出")
		NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
		return
		//panic("！！！ impossible")
	}

	// todo 判断是否有玩家报牌最后一首 需要顶牌

	// 最有组牌能否打过
	// canBeat := CanBeat(eCards, comGroup[i].Card)
	//beatCards, canBeat, beatType := FindCanBeatCards(robot.HandCards, eCards, eType)
	beatCards, canBeat, beatType := FindMinFollowCards(robot.HandCards, completeGroupCard(robot.GroupCard), eCards, eType)
	if canBeat { // 如果能打过
		//logger.Debug("fallow 11111111111111111111111")
		// 判断是否炸弹打过  如果是炸弹打过
		if beatType == cardConst.CARD_PATTERN_BOMB ||
			beatType == cardConst.CARD_PATTERN_ROCKET {

			eg := GroupHandsCard(ePlayer.HandCards)
			eCompleteGroups := completeGroupCard(eg)
			completeReCards := changeGroupToReCard(eCompleteGroups)
			outNums := len(completeReCards) // 出牌总首数

			// 这里如果自己炸了之后剩下的都是天牌  则可以炸
			selfG := GroupHandsCard(removeCards(robot.HandCards, beatCards))
			selfCompleteGroups := completeGroupCard(selfG)
			selfCompleteReCards := changeGroupToReCard(selfCompleteGroups)
			selfOutNums := len(selfCompleteReCards) // 出牌总首数
			_, selfGodNums := CheckRealGodCard(selfCompleteReCards, nextPlayer.HandCards, lastPlayer.HandCards)
			if selfOutNums-selfGodNums <= 1 {
				logger.Debug("地主跟牌 出炸 自己剩余所有天牌 !")
				OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
				return
			}

			// 这里如果自己炸了之后剩下的都是天牌  则可以炸

			if outNums == 1 {
				// 如果对方玩家最后一首是炸弹不出了
				if completeReCards[0].CardType == cardConst.CARD_PATTERN_BOMB ||
					completeReCards[0].CardType == cardConst.CARD_PATTERN_ROCKET {
					//logger.Debug("fallow 3333333333333333333")
					//logger.Debug("地主跟牌 出炸 对方最后一首是炸弹 不出!")
					NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
					return
				}
				//logger.Debug("地主跟牌 出炸 对方最后一首非炸弹牌 !")
				OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
				return
			}

			//
			logger.Debug("地主跟牌 不炸先隐忍一波!")
			NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
			return
		} else { // 如果不是炸弹打过 则重新判断先后组牌 的首数量
			// 先取出最优化组牌
			g := completeGroupCard(robot.GroupCard)
			comGroup := changeGroupToReCard(g)
			oldLen := len(comGroup)

			ifOutNewHands := removeCards(robot.HandCards, beatCards)

			newG := GroupHandsCard(ifOutNewHands)
			newComGroup := completeGroupCard(newG)
			newComRe := changeGroupToReCard(newComGroup)
			newLen := len(newComRe)

			if newLen-oldLen >= 2 { // 强拆打牌不划算
				logger.Debug("地主跟牌 不出 强拆不划算:", beatCards)
				NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
				return
			}
			logger.Debug("地主跟牌 普通跟牌,跟拍类型:", beatCards)
			OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
			return
		}
	}
	logger.Debug("!地主跟牌 没有打过的牌 不出:", eType)
	NotOutCardsAction(room, robot, lastPlayer, nextPlayer)

}



// 农民下家是农名出手牌策略
/*
	0.自己能否一首走完
	1.判断下家农民是否报单或者报双 如果是 则先出天炸 然后 最小对子或者最小单牌出牌

*/
func farmerRobotOutCardMustDoF1(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player, ) {
	logger.Debug("农民玩家一号首出")
	completeGroups := completeGroupCard(robot.GroupCard)
	completeReCards := changeGroupToReCard(completeGroups)
	outNums := len(completeReCards) // 出牌总首数
	godCard, godCardNum := CheckGodCard(completeReCards, nextPlayer.HandCards, lastPlayer.HandCards)
	// 如果最后两手 有一首是炸弹 并且 无人保单或者无人报双的情况下 先走非炸弹的
	nFarmerType, b1 := checkPlayerHasLast(nextPlayer) // 农民
	landType, b2 := checkPlayerHasLast(lastPlayer)    // 地主
	if outNums == 2 {
		if !b1 && !b2 {
			if completeReCards[0].CardType == cardConst.CARD_PATTERN_BOMB ||
				completeReCards[1].CardType == cardConst.CARD_PATTERN_BOMB {
				for i := 0; i < len(godCard); i++ {
					if completeReCards[i].CardType != cardConst.CARD_PATTERN_BOMB {
						logger.Debug("地主首出 最后一首是炸弹+其他 并无人报牌 显出其他")
						outCard := completeReCards[i].Card
						cardType := completeReCards[i].CardType
						OutCardsAction(room, robot, nextPlayer, outCard, cardType)
						return
					}
				}
			}
		}
	}

	if outNums-godCardNum <= 1 { // 这里如果总首数-天牌数<=1 // 则先出天牌
		for i := 0; i < len(godCard); i++ {
			if godCard[i].IsGodCard {
				outCard := godCard[i].RC.Card
				cardType := godCard[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, outCard, cardType)
				return
			}
		}
	}

	//判断下家是否报单或者双
	// 1.如果下家农民是炸弹或者火箭 如果自己有炸弹或者火箭 则先出炸弹或者火箭  没有则出权重最小组牌
	if nFarmerType == cardConst.CARD_PATTERN_BOMB || nFarmerType == cardConst.CARD_PATTERN_ROCKET {
		if len(robot.GroupCard.Bomb) >= 1 {
			OutCardsAction(room, robot, nextPlayer, robot.GroupCard.Bomb[0].Card, robot.GroupCard.Bomb[0].CardType)
			return
		} else if len(robot.GroupCard.Rocket) >= 1 {
			OutCardsAction(room, robot, nextPlayer, robot.GroupCard.Bomb[0].Card, robot.GroupCard.Bomb[0].CardType)
			return
		} else { // 随便出什么牌 这里出权重最小的牌
			SortReCardByWightSL(completeReCards)
			OutCardsAction(room, robot, nextPlayer, completeReCards[0].Card, completeReCards[0].CardType)
			return
		}
		// 如果下家玩家只剩一张单张 并且自己有天炸的情况下 先出天炸 再出单张
	} else if nFarmerType == cardConst.CARD_PATTERN_SINGLE {
		// todo 天炸
		// 1. 先判断自己最小的牌能否被下家打过 不能则取最小权重出牌
		card := findMinCard(robot.HandCards)
		if CanBeat(nextPlayer.HandCards, card) { // 能打过
			// 自己是否有天炸
			OutCardsAction(room, robot, nextPlayer, card, cardConst.CARD_PATTERN_SINGLE)
			return
		} else { // 最小权重出牌
			SortReCardByWightSL(completeReCards)
			OutCardsAction(room, robot, nextPlayer, completeReCards[0].Card, completeReCards[0].CardType)
			return
		}
	} else if nFarmerType == cardConst.CARD_PATTERN_PAIR {
		// 去最小的两对子 这个方法和 findMinDouble()
		card := findMinDoubleCard(robot.HandCards)
		if CanBeat(nextPlayer.HandCards, card) { // 能打过
			OutCardsAction(room, robot, nextPlayer, card, cardConst.CARD_PATTERN_SINGLE)
			return
		}

		// 最小权重出牌
		SortReCardByWightSL(completeReCards)
		OutCardsAction(room, robot, nextPlayer, completeReCards[0].Card, completeReCards[0].CardType)
		return
	}

	// 这里很复杂
	/*===========  地主报单 ===========*/
	if landType == cardConst.CARD_PATTERN_SINGLE {
		logger.Debug("F1首出 地主玩家有最后一手单牌", cardConst.CARD_PATTERN_SINGLE)
		for i := 0; i < len(godCard); i++ {
			if godCard[i].RC.CardType != cardConst.CARD_PATTERN_SINGLE {
				outCard := godCard[i].RC.Card
				cardType := godCard[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, outCard, cardType)
				return
			}
		}
		if len(completeReCards) >= 2 {
			logger.Debug("F1首出 地主玩家有最后一单牌 只有和这种单牌  出最大单牌")
			outCard := completeReCards[len(completeReCards)-1].Card
			cardType := completeReCards[len(completeReCards)-1].CardType
			OutCardsAction(room, robot, nextPlayer, outCard, cardType)
			return
		}
	}

	/*===========  地主报双 ===========*/
	if landType == cardConst.CARD_PATTERN_PAIR {
		logger.Debug("F1首出 地主玩家有最后一手牌", cardConst.CARD_PATTERN_PAIR)
		for i := 0; i < len(godCard); i++ {
			if godCard[i].RC.CardType != cardConst.CARD_PATTERN_PAIR {
				outCard := godCard[i].RC.Card
				cardType := godCard[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, outCard, cardType)
				return
			}
		}
		if len(completeReCards) >= 2 {
			logger.Debug("F1首出 地主玩家有最后一对子 只有和这种单牌  出最小单牌")
			cards := findMinCard(robot.HandCards)
			OutCardsAction(room, robot, nextPlayer, cards, cardConst.CARD_PATTERN_SINGLE)
			return
		}
	}

	SortReCardByWightSL(completeReCards)
	logger.Debug("F1首出 无人最后一手牌 取最小权值牌")
	OutCardsAction(room, robot, nextPlayer, completeReCards[0].Card, completeReCards[0].CardType)
}

// 农民下家是地主出首牌策略
/*
	先看自己能否走完  先顾自己在顾别人

*/
func farmerRobotOutCardMustDoF2(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	//
	logger.Debug("F2首出")
	// 1.总首数-天牌数<=1  ：随机出弹外的天牌 最后一首出 普通手牌
	completeGroups := completeGroupCard(robot.GroupCard)
	completeReCards := changeGroupToReCard(completeGroups)
	outNums := len(completeReCards) // 出牌总首数
	godCard, godCardNum := CheckGodCard(completeReCards, nextPlayer.HandCards, lastPlayer.HandCards)

	// 如果最后两手 有一首是炸弹 并且 无人保单或者无人报双的情况下 先走非炸弹的
	landType, b1 := checkPlayerHasLast(nextPlayer) // 地主
	_, b2 := checkPlayerHasLast(lastPlayer)        // 农民1号
	if outNums == 2 {
		if !b1 && !b2 {
			if completeReCards[0].CardType == cardConst.CARD_PATTERN_BOMB ||
				completeReCards[1].CardType == cardConst.CARD_PATTERN_BOMB {
				for i := 0; i < len(godCard); i++ {
					if completeReCards[i].CardType != cardConst.CARD_PATTERN_BOMB {
						logger.Debug("地主首出 最后一首是炸弹+其他 并无人报牌 显出其他")
						outCard := completeReCards[i].Card
						cardType := completeReCards[i].CardType
						// todo 发送表情 根据概率发送 表情 和你合作很愉快
						OutCardsAction(room, robot, nextPlayer, outCard, cardType)
						return
					}
				}
			}
		}
	}

	if outNums-godCardNum <= 1 { // 这里如果总首数-天牌数<=1 // 则先出天牌
		for i := 0; i < len(godCard); i++ {
			if godCard[i].IsGodCard {
				outCard := godCard[i].RC.Card
				cardType := godCard[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, outCard, cardType)
				return
			}
		}
	}

	// 这里先判断地主 报单的情况
	if landType == cardConst.CARD_PATTERN_SINGLE {
		logger.Debug("F2首出 下一个有人最后一手牌", cardConst.CARD_PATTERN_SINGLE)
		for i := 0; i < len(godCard); i++ {
			if godCard[i].RC.CardType != cardConst.CARD_PATTERN_SINGLE {
				outCard := godCard[i].RC.Card
				cardType := godCard[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, outCard, cardType)
				return
			}
		}
		if len(completeReCards) >= 2 {
			logger.Debug("F2首出 地主玩家有最后一单牌 只有和这种单牌  出最大单牌")
			outCard := completeReCards[len(completeReCards)-1].Card
			cardType := completeReCards[len(completeReCards)-1].CardType
			OutCardsAction(room, robot, nextPlayer, outCard, cardType)
			return
		}
	}

	/*===========  地主报双 ===========*/
	if landType == cardConst.CARD_PATTERN_PAIR {
		logger.Debug("F2首出 地主玩家有最后一对", cardConst.CARD_PATTERN_PAIR)
		for i := 0; i < len(godCard); i++ {
			if godCard[i].RC.CardType != cardConst.CARD_PATTERN_PAIR {
				outCard := godCard[i].RC.Card
				cardType := godCard[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, outCard, cardType)
				return
			}
		}
		if len(completeReCards) >= 2 {
			logger.Debug("F2首出 地主玩家有最后一对子 只有和这种单牌  出最小的一张牌")
			cards := findMinCard(robot.HandCards)
			OutCardsAction(room, robot, nextPlayer, cards, cardConst.CARD_PATTERN_SINGLE)
			return
		}
	}

	SortReCardByWightSL(completeReCards)
	logger.Debug("F2首出 无人最后一手牌 取最小权值牌")
	OutCardsAction(room, robot, nextPlayer, completeReCards[0].Card, completeReCards[0].CardType)
}

/*e  农民出牌策略 还有很大的优化空间*/

/*
	一 农民玩家 手出策略
	F1 F2 LANDLORD
*/
func farmerRobotFallowCard(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	efficId := room.EffectivePlayerId
	efficP := room.Players[efficId]
	if efficP.IsLandlord && !nextPlayer.IsLandlord {
		// 农民下家是农民跟牌策略
		//farmerRobotFallowCardF1(room, robot, nextPlayer, lastPlayer, !efficP.IsLandlord)
		farmerFallowF1(room, robot, nextPlayer, lastPlayer)
	} else {
		// 农民下家是地主跟牌策略
		//farmerRobotFallowCardF1(room, robot, nextPlayer, lastPlayer)
		//farmerRobotFallowCardF2(room, robot, nextPlayer, lastPlayer, !efficP.IsLandlord)
		farmerFallowF2(room, robot, nextPlayer, lastPlayer)
	}
}

/*e  农民跟牌策略 还有很大的优化空间*/
/*
	一 农民玩家F1跟牌 下家是农民 手出策略
	F1 -> F2 -> LANDLORD
	F1
*/
func farmerRobotFallowCardF1(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player, isFriendOut bool) {
	eType := room.EffectiveType
	eCards := room.EffectiveCard
	ePlayer := room.Players[room.EffectivePlayerId]

	// 最有组牌能否打过
	// canBeat := CanBeat(eCards, comGroup[i].Card)
	//beatCards, canBeat, beatType := FindCanBeatCards(robot.HandCards, eCards, eType)
	beatCards, canBeat, beatType := FindMinFollowCards(robot.HandCards, completeGroupCard(robot.GroupCard), eCards, eType)
	if canBeat { // 如果能打过
		if ePlayer.IsLandlord { // 如果上首牌是地主出的
			// 判断是否炸弹打过  如果是炸弹打过
			if beatType == cardConst.CARD_PATTERN_BOMB ||
				beatType == cardConst.CARD_PATTERN_ROCKET {

				eg := GroupHandsCard(ePlayer.HandCards)
				eCompleteGroups := completeGroupCard(eg)
				completeReCards := changeGroupToReCard(eCompleteGroups)
				outNums := len(completeReCards) // 出牌总首数

				// 这里如果自己炸了之后剩下的都是天牌  则可以炸
				selfG := GroupHandsCard(removeCards(robot.HandCards, beatCards))
				selfCompleteGroups := completeGroupCard(selfG)
				selfCompleteReCards := changeGroupToReCard(selfCompleteGroups)
				selfOutNums := len(selfCompleteReCards) // 出牌总首数
				_, selfGodNums := CheckRealGodCard(selfCompleteReCards, nextPlayer.HandCards, lastPlayer.HandCards)
				if selfOutNums-selfGodNums <= 1 {
					logger.Debug("农民玩家F1跟牌地主 出炸 自己剩余所有天牌 !")
					OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
					return
				}

				// 这里如果自己炸了之后剩下的都是天牌  则可以炸

				if outNums == 1 {
					// 如果对方玩家最后一首是炸弹不出了
					if completeReCards[0].CardType == cardConst.CARD_PATTERN_BOMB ||
						completeReCards[0].CardType == cardConst.CARD_PATTERN_ROCKET {
						//logger.Debug("fallow 3333333333333333333")
						//logger.Debug("地主跟牌 出炸 对方最后一首是炸弹 不出!")
						NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
						return
					}
					//logger.Debug("地主跟牌 出炸 对方最后一首非炸弹牌 !")
					OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
					return
				}

				//
				logger.Debug("农民玩家F1跟牌地主 不炸先隐忍一波!")
				NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
				return
			} else { // 如果不是炸弹打过 则重新判断先后组牌 的首数量
				// 先取出最优化组牌
				g := completeGroupCard(robot.GroupCard)
				comGroup := changeGroupToReCard(g)
				oldLen := len(comGroup)

				ifOutNewHands := removeCards(robot.HandCards, beatCards)

				newG := GroupHandsCard(ifOutNewHands)
				newComGroup := completeGroupCard(newG)
				newComRe := changeGroupToReCard(newComGroup)
				newLen := len(newComRe)

				if newLen-oldLen > 0 { // 强拆打牌不划算 这里是0
					logger.Debug("农民玩家F1跟牌地主跟牌 不出 强拆不划算:", beatCards)
					NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
					return
				}
				logger.Debug("农民玩家F1跟牌地主跟牌 普通跟牌,跟拍类型:", beatCards)
				OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
				return
			}
		} else { // 如果跟友方玩家的牌
			lcCards, lcType, b := checkLandlordHasLast(lastPlayer) // 判断地主是否最后一首
			if b && lcType == cardConst.CARD_PATTERN_SINGLE { // 同时地主玩家报单的情况下
				logger.Debug("!农民玩家F1跟牌 跟下家 地主报单 ", eType)
				comGroup := completeGroupCard(robot.GroupCard)
				canotBeatNum := checkHowManyBeatSingle(comGroup.Single, lcCards)
				if canotBeatNum <= 1 { // 自己的手牌单牌能小于1
					OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
					return
				} else {
					friCardG := completeGroupCard(nextPlayer.GroupCard)
					cannotBeat := checkHowManyBeatSingle(friCardG.Single, lcCards)
					if cannotBeat <= 1 { // 如果友方玩家单牌小于数小于1 则不出
						NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
						return
					} else {
						OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
						return
					}
				}
			} else { // 如果地主玩家没报单的情况下 总首数-天牌数 谁小谁出
				comGroup := completeGroupCard(robot.GroupCard)
				selfComRe := changeGroupToReCard(comGroup)
				_, selfGodNum := CheckRealGodCard(selfComRe, lastPlayer.HandCards, nextPlayer.HandCards)

				friCardG := completeGroupCard(nextPlayer.GroupCard)
				friComRe := changeGroupToReCard(friCardG)
				_, friGodNum := CheckRealGodCard(friComRe, robot.HandCards, lastPlayer.HandCards)
				// 如果
				if (len(friComRe) - friGodNum) < (len(selfComRe) - selfGodNum) {
					OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
					return
				} else {
					OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
					return
				}
			}
		}
	}

	logger.Debug("!农民玩家F1跟牌 没有打过的牌 不出:", eType)
	NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
}

/*e  农民出牌策略 还有很大的优化空间*/
/*
	下家是地主
	一 农民玩家 手出策略
	F1 -> F2 -> LANDLORD
	F2
*/
func farmerRobotFallowCardF2(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player, isFriendOut bool) {
	eType := room.EffectiveType
	eCards := room.EffectiveCard
	ePlayer := room.Players[room.EffectivePlayerId]

	// 最有组牌能否打过
	// canBeat := CanBeat(eCards, comGroup[i].Card)
	//beatCards, canBeat, beatType := FindCanBeatCards(robot.HandCards, eCards, eType)
	beatCards, canBeat, beatType := FindMinFollowCards(robot.HandCards, completeGroupCard(robot.GroupCard), eCards, eType)
	if canBeat { // 如果能打过
		if ePlayer.IsLandlord { // 如果上首牌是地主出的
			//logger.Debug("fallow 11111111111111111111111")
			// 判断是否炸弹打过  如果是炸弹打过
			if beatType == cardConst.CARD_PATTERN_BOMB ||
				beatType == cardConst.CARD_PATTERN_ROCKET {

				eg := GroupHandsCard(ePlayer.HandCards)
				eCompleteGroups := completeGroupCard(eg)
				completeReCards := changeGroupToReCard(eCompleteGroups)
				outNums := len(completeReCards) // 出牌总首数

				// 这里如果自己炸了之后剩下的都是天牌  则可以炸
				selfG := GroupHandsCard(removeCards(robot.HandCards, beatCards))
				selfCompleteGroups := completeGroupCard(selfG)
				selfCompleteReCards := changeGroupToReCard(selfCompleteGroups)
				selfOutNums := len(selfCompleteReCards) // 出牌总首数
				_, selfGodNums := CheckRealGodCard(selfCompleteReCards, nextPlayer.HandCards, lastPlayer.HandCards)
				if selfOutNums-selfGodNums <= 1 {
					logger.Debug("农民玩家F2跟牌地主 出炸 自己剩余所有天牌 !")
					OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
					return
				}

				// 这里如果自己炸了之后剩下的都是天牌  则可以炸

				if outNums == 1 {
					// 如果对方玩家最后一首是炸弹不出了
					if completeReCards[0].CardType == cardConst.CARD_PATTERN_BOMB ||
						completeReCards[0].CardType == cardConst.CARD_PATTERN_ROCKET {
						//logger.Debug("fallow 3333333333333333333")
						//logger.Debug("地主跟牌 出炸 对方最后一首是炸弹 不出!")
						NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
						return
					}
					//logger.Debug("地主跟牌 出炸 对方最后一首非炸弹牌 !")
					OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
					return
				}

				//
				logger.Debug("农民玩家F2跟牌地主 不炸先隐忍一波!")
				NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
				return
			} else { // 如果不是炸弹打过 则重新判断先后组牌 的首数量
				// 先取出最优化组牌
				g := completeGroupCard(robot.GroupCard)
				comGroup := changeGroupToReCard(g)
				oldLen := len(comGroup)

				ifOutNewHands := removeCards(robot.HandCards, beatCards)

				newG := GroupHandsCard(ifOutNewHands)
				newComGroup := completeGroupCard(newG)
				newComRe := changeGroupToReCard(newComGroup)
				newLen := len(newComRe)

				if newLen-oldLen > 0 { // 强拆打牌不划算 这里是0
					logger.Debug("农民玩家F2跟牌地主跟牌 不出 强拆不划算:", beatCards)
					NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
					return
				}
				logger.Debug("农民玩家F2跟牌地主跟牌 普通跟牌,跟拍类型:", beatCards)
				OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
				return
			}
		} else { // 如果跟友方玩家的牌
			lcCards, lcType, b := checkLandlordHasLast(lastPlayer) // 判断地主是否最后一首
			if b && lcType == cardConst.CARD_PATTERN_SINGLE { // 同时地主玩家报单的情况下
				logger.Debug("!农民玩家F2跟牌 跟下家 地主报单 ", eType)
				comGroup := completeGroupCard(robot.GroupCard)
				canotBeatNum := checkHowManyBeatSingle(comGroup.Single, lcCards)
				if canotBeatNum <= 1 { // 自己的手牌单牌能小于1
					OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
					return
				} else {
					friCardG := completeGroupCard(nextPlayer.GroupCard)
					cannotBeat := checkHowManyBeatSingle(friCardG.Single, lcCards)
					if cannotBeat <= 1 { // 如果友方玩家单牌小于数小于1 则不出
						NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
						return
					} else {
						OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
						return
					}
				}
			} else { // 如果地主玩家没报单的情况下 总首数-天牌数 谁小谁出
				comGroup := completeGroupCard(robot.GroupCard)
				selfComRe := changeGroupToReCard(comGroup)
				_, selfGodNum := CheckRealGodCard(selfComRe, lastPlayer.HandCards, nextPlayer.HandCards)

				friCardG := completeGroupCard(nextPlayer.GroupCard)
				friComRe := changeGroupToReCard(friCardG)
				_, friGodNum := CheckRealGodCard(friComRe, robot.HandCards, lastPlayer.HandCards)
				// 如果
				if (len(friComRe) - friGodNum) < (len(selfComRe) - selfGodNum) {
					OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
					return
				} else {
					OutCardsAction(room, robot, nextPlayer, beatCards, beatType)
					return
				}
			}
		}
	}

	logger.Debug("!农民玩家F2跟牌 没有打过的牌 不出:", eType)
	NotOutCardsAction(room, robot, lastPlayer, nextPlayer)

}

