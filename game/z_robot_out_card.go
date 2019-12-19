package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/cardConst"
)

// 机器人出牌
func robotOutCard(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	if robot.IsLandlord {
		if robot.IsMustDo {
			// todo 地主机器人首出
			landlordRobotOutCardMustDo(room, robot, nextPlayer, lastPlayer)
		} else {
			// todo // 地主机器人跟牌
			landlordRobotFallowCardF1(room, robot, nextPlayer, lastPlayer)
		}
	} else {
		// todo

		if robot.IsMustDo {
			// 农民首出
			landlordRobotOutCardMustDo(room, robot, nextPlayer, lastPlayer)

		} else {
			// 农民玩家一号跟牌
			farmerRobotFallowCard(room, robot, nextPlayer, lastPlayer)
		}

		// DoGameHosting(room, robot, nextPlayer, lastPlayer)
	}

}

/*s 地主出牌策略 还有很大的优化空间*/

// 一 地主机器人首出牌策虐
/*
	0.判断能否一首出完 // todo 注意 三代把王打出去了
	1.检测 对方玩家是否 报单或者报双 或者最后一手牌	(todo 另写方法出牌)

	取牌组中权重值最小的牌
*/
func landlordRobotOutCardMustDo(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {

	//1. 先判断能否一首出完 	// 检测牌型正确代表能一首出完
	cardType := GetCardsType(robot.HandCards)
	if !(cardType > cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS || cardType < cardConst.CARD_PATTERN_SINGLE) {
		/*
				todo 这个牌型是否三带中包含对王  如果包含对王则 是否地主春天... 很多考虑因素
			&{5 2},&{5 4},
			&{6 1},&{6 2},&{6 3},&{6 4},
			&{7 1},&{7 3},&{7 2},
			&{8 1},&{8 2},&{8 4},
			&{9 1},&{9 2},&{9 4},
			&{12 4},
		*/
		OutCardsAction(room, robot, nextPlayer, robot.HandCards, cardType)
		return
	}

	endCards, finalType := checkNotLandlordHasLast(nextPlayer, lastPlayer) // 这里的都是农民
	if finalType == 0 { // 无人最后一手牌
		landlordRobotOutCardMustDoNormal(room, robot, nextPlayer, lastPlayer)
	} else {
		//landlordRobotOutCardMustDoNormal(room, robot, nextPlayer, lastPlayer)
		landlordRobotOutCardMustDoHasLast(room, robot, nextPlayer, lastPlayer, finalType, endCards)
	}

}

/*
	一 地主机器人首出牌策虐
	1.如果敌对方玩家有最后一手牌的情况下
*/
func landlordRobotOutCardMustDoHasLast(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player,
	endCardType int32, endCard []*Card) {
	var outCard *ReCard
	if endCardType == cardConst.CARD_PATTERN_SINGLE { // 如果有玩家报单
		outCard = findMinWeightCardsExpectSome(robot.GroupCard, true, false)
		if outCard == nil { // 如果都没有牌了
			// 判断自己的牌的张数是否
			outCard = findBestSingleCard(robot.GroupCard.Single)
			if outCard == nil {
				cards, cardType := FindMustBeOutCards(robot.HandCards)
				OutCardsAction(room, robot, nextPlayer, cards, cardType)
				return
			}
		}
	} else if endCardType == cardConst.CARD_PATTERN_PAIR { // 如果有玩家报双
		outCard = findMinWeightCardsExpectSome(robot.GroupCard, false, true)
		if outCard == nil { // 除了炸弹都没牌 只有对子的情况下
			outCard = findMinWeightCards(robot.GroupCard)
			if outCard.CardType == cardConst.CARD_PATTERN_PAIR {
				if outCard.Wight >= endCard[0].Value {
					OutCardsAction(room, robot, nextPlayer, outCard.Card, outCard.CardType)
					return
				}
			} else {
				OutCardsAction(room, robot, nextPlayer, append([]*Card{}, outCard.Card[0]), cardConst.CARD_PATTERN_SINGLE)
				return
			}

		}
		OutCardsAction(room, robot, nextPlayer, outCard.Card, outCard.CardType)
	}

}

/*
	// todo 注意 三代把王打出去了
	一 地主机器人首出牌策虐
	1.如果敌对方玩家有最后一手牌的情况下
	常规的地主首牌策略
	取各组牌权重值最小的牌
*/
func landlordRobotOutCardMustDoNormal(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {

	reCard := findMinWeightCards(robot.GroupCard)
	if reCard == nil {
		logger.Debug("机器人首出 没有找到最小权重.....!!!!!!!!!!!!!!!!!!!")
		outCards, ct := FindMustBeOutCards(robot.HandCards)
		OutCardsAction(room, robot, nextPlayer, outCards, ct)
		return
	}
	OutCardsAction(room, robot, nextPlayer, reCard.Card, reCard.CardType)
	//OutCardsAction(room, robot, nextPlayer, outCard, cardType)
}

// 二 地主机器人跟牌策虐
/*
	出牌顺序  f1 f2 landlord
	f1 跟牌


	// 1. 判断玩家是否打完这首牌之后 只剩最后一首牌
*/
func landlordRobotFallowCardF1(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	// 暂时使用托管规则
	DoGameHosting(room, robot, nextPlayer, lastPlayer)

	//lastCard := room.EffectiveCard
	//lastType := room.EffectiveType
	//
	//switch lastType {
	//case cardConst.CARD_PATTERN_SINGLE:  // 跟单张
	//	
	//	
	//}

}

/*e  农民出牌策略 还有很大的优化空间*/

/*
	一 农民玩家 手出策略
	F1 F2 LANDLORD
*/

func farmerRobotFallowCard(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	efficId := room.EffectivePlayerId
	efficP := room.Players[efficId]
	if efficP.IsLandlord && !nextPlayer.IsLandlord { // 我的下家不是地主
		farmerRobotFallowCardF1(room, robot, nextPlayer, lastPlayer)
	} else { // 下家是地主
		//farmerRobotFallowCardF1(room, robot, nextPlayer, lastPlayer)
		farmerRobotFallowCardF2(room, robot, nextPlayer, lastPlayer,!efficP.IsLandlord)
	}

}

/*e  农民出牌策略 还有很大的优化空间*/
/*
	一 农民玩家 手出策略
	F1 -> F2 -> LANDLORD
	F1
*/

func farmerRobotFallowCardF1(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	efficType := room.EffectiveType
	efficCard := room.EffectiveCard

	// 1. 判断地主是否最后一首牌
	// checkLandlordHasLast()

	// 采用最小跟牌
	cards, b ,cardType:= minFollowCard(robot, efficCard, efficType)
	if !b { // todo  这里要多做判断
		NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
	}
	// 待修复
	OutCardsAction(room, robot, nextPlayer, cards, cardType)
}

/*e  农民出牌策略 还有很大的优化空间*/
/*
	下家是地主
	一 农民玩家 手出策略
	F1 -> F2 -> LANDLORD
	F2
*/

func farmerRobotFallowCardF2(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player, isFriendOut bool) {
	efficType := room.EffectiveType
	efficCard := room.EffectiveCard

	// 1. 判断地主是否最后一首牌 有出必出
	_, _, b := checkLandlordHasLast(room.Players[room.LandlordPlayerId])
	if b {
		DoGameHosting(room, robot, nextPlayer, lastPlayer)
		return
	}

	// 采用最小跟牌
	cards, b := minFollowCard(robot, efficCard, efficType)
	if !b { // todo  这里要多做判断
		NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
	}
	// 待修复
	OutCardsAction(room, robot, nextPlayer, cards, efficType)

}

//
