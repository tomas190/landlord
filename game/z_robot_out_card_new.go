package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/cardConst"
	"landlord/mconst/playerAction"
	"strconv"
	"strings"
)

/*
地主首出 finish
	总手数牌-天牌数<=1
		是:     【出非炸弹的权重最小天牌】(若都是炸弹则出权重最小天牌炸弹)

		否： 1.农民是否报单
			是：【最小权重 非单牌】(若只有单牌 【权重第二小单牌】)
			否：【最小权重牌】
		否： 2.农民是否报双
			是：【最小权重 非对子】(若只有对子牌 【最小单牌】)
			否：【最小权重牌】
*/

func NewLandlordRobotOutCardMustDo(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {

	// 2020年2月21日16:41:43 如果能一首出完
	hc := robot.HandCards
	cardsType := GetCardsType(hc)
	if cardsType>=cardConst.CARD_PATTERN_SINGLE&&cardsType<=cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_PAIRS {
		OutCardsAction(room, robot, nextPlayer, hc, cardsType)
		return
	}
	// 2020年2月21日16:41:43 如果能一首出完

	comG := completeGroupCard(robot.GroupCard)
	comRe := changeGroupToReCard(comG)
	SortReCardByWightSL(comRe)
	allLen := len(comRe)
	godCard, godNum := CheckRealGodCard(comRe, nextPlayer.HandCards, lastPlayer.HandCards)

	// 如果最后两手 有一首是炸弹 并且 无人保单或者无人报双的情况下 先走非炸弹的
	cType1, b1 := checkPlayerHasLast(nextPlayer)
	cType2, b2 := checkPlayerHasLast(lastPlayer)
	if allLen == 2 {
		var tmpThrows []*Card
		tmpThrows = append(tmpThrows, nextPlayer.ThrowCards...)
		tmpThrows = append(tmpThrows, lastPlayer.ThrowCards...)
		if !b1 && !b2 && len(tmpThrows) > 0 { // 并且有人出牌过牌的情况 这里要保春天
			if comRe[0].CardType == cardConst.CARD_PATTERN_BOMB ||
				comRe[1].CardType == cardConst.CARD_PATTERN_ROCKET {
				for i := 0; i < len(godCard); i++ {
					if comRe[i].CardType != cardConst.CARD_PATTERN_BOMB {

						logger.Debug("地主首出 最后一首是炸弹+其他 并无人报牌且不用保春天 先出其他")
						outCard := comRe[i].Card
						cardType := comRe[i].CardType
						OutCardsAction(room, robot, nextPlayer, outCard, cardType)
						return
					}
				}
			}
		}
	}

	// 这里如果总首数-天牌数<=1 // 则先出非炸弹天牌
	if allLen-godNum <= 1 {
		for i := 0; i < len(godCard); i++ {
			if godCard[i].IsGodCard && !(godCard[i].RC.CardType == cardConst.CARD_PATTERN_BOMB || godCard[i].RC.CardType == cardConst.CARD_PATTERN_ROCKET) {
				outCard := godCard[i].RC.Card
				cardType := godCard[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, outCard, cardType)
				return
			}
		}

		// 这里如果总首数-天牌数<=1 如果
		for i := 0; i < len(godCard); i++ {
			if godCard[i].IsGodCard {
				outCard := godCard[i].RC.Card
				cardType := godCard[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, outCard, cardType)
				return
			}
		}
	}

	// 如果有农民报单
	if (b1 && cType1 == cardConst.CARD_PATTERN_SINGLE) || (b2 && cType2 == cardConst.CARD_PATTERN_SINGLE) {
		// 有农民保单
		for i := 0; i < len(comRe); i++ {
			if comRe[i].CardType != cardConst.CARD_PATTERN_SINGLE {
				logger.Debug("地主首出 农民保单 非类型最小权重")
				outCard := comRe[i].Card
				cardType := comRe[i].CardType
				OutCardsAction(room, robot, nextPlayer, outCard, cardType)
				return
			}
		}

		if len(comRe) >= 2 {
			logger.Debug("地主首出 只有和这种类型相似的牌  出第二大")
			OutCardsAction(room, robot, nextPlayer, SpecialHandle(comRe[1].Card), comRe[1].CardType)
			return
		}
	}

	// 如果有农民报双
	if (b1 && cType1 == cardConst.CARD_PATTERN_PAIR) || (b2 && cType2 == cardConst.CARD_PATTERN_PAIR) {
		// 有农民保单
		for i := 0; i < len(comRe); i++ {
			if comRe[i].CardType != cardConst.CARD_PATTERN_PAIR {
				logger.Debug("地主首出 农民保双 非类型最小权重")
				outCard := comRe[i].Card
				cardType := comRe[i].CardType
				OutCardsAction(room, robot, nextPlayer, SpecialHandle(outCard), cardType)
				return
			}
		}

		if len(comRe) >= 2 {
			logger.Debug("地主首出 如果全是对子 出最小单牌")
			outCard := findMinCard(robot.HandCards)
			OutCardsAction(room, robot, nextPlayer, SpecialHandle(outCard), cardConst.CARD_PATTERN_SINGLE)
			//OutCardsAction(room, robot, nextPlayer, outCard, cardConst.CARD_PATTERN_SINGLE)
			return
		}
	}

	logger.Debug("地主首出 无人最后一手牌 取最小权值牌")
	//checkCard := SpecialHandle(comRe[0].Card)
	OutCardsAction(room, robot, nextPlayer, SpecialHandle(comRe[0].Card), comRe[0].CardType)
	//OutCardsAction(room, robot, nextPlayer, comRe[0].Card, comRe[0].CardType)
}

// 二 地主机器人跟牌策虐 finish
/*
地主跟牌
	是否有玩家报单
		当前出牌类型是否单牌
			是：查看单牌张数 如果大于等于3 【最小第二权值单张 】，若等于2【最大权值单张】，等于1,检查是否有炸弹，和天牌数量。满足条件 all-god<=1【炸】
			否： 【相同类型最小跟牌 】 如果没有则寻找拆牌之后能打过的牌 如果 手数相同或者<=1则 【出能打必大最小拍,】如果找出炸弹 判断 all-god<=1则炸 
	是否有玩家报双
		当前出牌类型是否对子
			是：查看对子牌张数 如果大于等于3 【最小第二权值单张 】，若等于2【最大权值对子】，等于1,检查是否有炸弹，和天牌数量。满足条件 all-god<=1【炸】
			否： 【相同类型最小跟牌 】 如果没有则寻找拆牌之后能打过的牌 如果 手数相同或者<=1则 【出能打必大最小拍,】如果找出炸弹 判断 all-god<=1则炸 
	普通跟牌
		最小类型牌:如果没有 则寻找 有大必打跟牌 【满足条件手牌数 可跟】如果跟出炸弹 all-god<=1则炸  


*/
func NewLandlordRobotFallowCard(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	eType := room.EffectiveType
	eCards := room.EffectiveCard

	compG := completeGroupCard(robot.GroupCard)
	comRe := changeGroupToReCard(compG)
	SortReCardByWightSL(comRe)
	allLen := len(comRe)

	cType1, b1 := checkPlayerHasLast(nextPlayer)
	cType2, b2 := checkPlayerHasLast(lastPlayer)

	// 是否有玩家报单
	if (b1 && cType1 == cardConst.CARD_PATTERN_SINGLE) || (b2 && cType2 == cardConst.CARD_PATTERN_SINGLE) {
		if eType == cardConst.CARD_PATTERN_SINGLE { //当前出牌类型是否单牌

			var callSingleCard []*Card
			// 如果两家农民都保单 则取最大保单牌 比较
			if cType1 == cardConst.CARD_PATTERN_SINGLE && cType2 == cardConst.CARD_PATTERN_SINGLE {
				if CanBeat(nextPlayer.HandCards, lastPlayer.HandCards) {
					callSingleCard = lastPlayer.HandCards
				} else {
					callSingleCard = nextPlayer.HandCards
				}
			} else {
				if cType1 == cardConst.CARD_PATTERN_SINGLE {
					callSingleCard = nextPlayer.HandCards
				} else {
					callSingleCard = lastPlayer.HandCards
				}
			}

			singles := compG.Single
			SortReCardByWightSL(singles)
			singleLen := len(singles)
			// 如果自己有单牌
			if singleLen >= 2 {
				//自己的最大权重单张能打过则出最大的 则出这个 不管是否能打过保单牌了
				if CanBeat(eCards, singles[len(singles)-1].Card) {
					OutCardsAction(room, robot, nextPlayer, singles[len(singles)-1].Card, cardConst.CARD_PATTERN_SINGLE)
					return
				} else { // 如果这里都不能打过  应该就凉了 随便出把 直接走下面普通更牌就可以了

				}
				// 如果只有一张单牌
			} else if singleLen == 1 {
				// 如果能打过并且能打过保单牌 则出
				if CanBeat(eCards, singles[0].Card) && CanBeat(callSingleCard, singles[0].Card) {
					OutCardsAction(room, robot, nextPlayer, singles[0].Card, cardConst.CARD_PATTERN_SINGLE)
					return
				} else { //
					// 寻找炸弹 如果有炸弹 且满足条件 炸
					cards, b, bType := HostingBeatBomb(robot.HandCards, []*Card{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}})
					if b { // 如果有炸弹
						// 这里要不要把炸弹移除了在判断
						_, godNum := CheckRealGodCard(comRe, nextPlayer.HandCards, lastPlayer.HandCards)
						if allLen-godNum <= 2 { // 这里等于2就可以炸了
							OutCardsAction(room, robot, nextPlayer, cards, bType)
							return
						}
					}
				}
			}

			// 寻找最大单牌
			oCard := findMaxCard(robot.HandCards)
			if CanBeat(eCards, oCard) {
				OutCardsAction(room, robot, nextPlayer, oCard, cardConst.CARD_PATTERN_SINGLE)
				return
			}

		}
		// 这里有出必出 炸弹除外
		cards, b, bType := FindCanBeatCards(robot.HandCards, eCards, eType)
		if b {
			// 这里要不要把炸弹移除了在判断
			_, godNum := CheckRealGodCard(comRe, nextPlayer.HandCards, lastPlayer.HandCards)
			if allLen-godNum <= 2 || !(bType == cardConst.CARD_PATTERN_BOMB || bType == cardConst.CARD_PATTERN_ROCKET) { // 这里等于2就可以炸了
				OutCardsAction(room, robot, nextPlayer, cards, bType)
				return
			}
		}

	}

	// 是否有玩家报双
	if (b1 && cType1 == cardConst.CARD_PATTERN_PAIR) || (b2 && cType2 == cardConst.CARD_PATTERN_PAIR) {
		if eType == cardConst.CARD_PATTERN_PAIR { //当前出牌类型是否对子
			// todo
			var callDoubleCard []*Card
			// 如果两家农民都报双 则取最大保双牌 比较
			if cType1 == cardConst.CARD_PATTERN_PAIR && cType2 == cardConst.CARD_PATTERN_PAIR {
				if CanBeat(nextPlayer.HandCards, lastPlayer.HandCards) {
					callDoubleCard = lastPlayer.HandCards
				} else {
					callDoubleCard = nextPlayer.HandCards
				}
			} else {
				if cType1 == cardConst.CARD_PATTERN_PAIR {
					callDoubleCard = nextPlayer.HandCards
				} else {
					callDoubleCard = lastPlayer.HandCards
				}
			}

			// 如果自己有对子
			doubles := compG.Double
			SortReCardByWightSL(doubles)
			doubleLen := len(doubles)
			if doubleLen > 0 {
				// 取最大对子 出
				beat := CanBeat(eCards, doubles[doubleLen-1].Card) && CanBeat(callDoubleCard, doubles[doubleLen-1].Card)
				if beat {
					OutCardsAction(room, robot, nextPlayer, doubles[doubleLen-1].Card, cardConst.CARD_PATTERN_PAIR)
					return
				}
			} else {
				// 寻找炸弹 如果有炸弹 且满足条件 炸
				cards, b, bType := HostingBeatBomb(robot.HandCards, []*Card{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}})
				if b { // 如果有炸弹
					// 这里要不要把炸弹移除了在判断
					_, godNum := CheckRealGodCard(comRe, nextPlayer.HandCards, lastPlayer.HandCards)
					if allLen-godNum <= 2 { // 这里等于2就可以炸了
						OutCardsAction(room, robot, nextPlayer, cards, bType)
						return
					}
				}
			}
			// 寻找最大对子
			oCard := findMaxDouble(robot.HandCards)
			if CanBeat(eCards, oCard) {
				OutCardsAction(room, robot, nextPlayer, oCard, cardConst.CARD_PATTERN_SINGLE)
				return
			}
		}

		// 这里有出必出 炸弹除外
		cards, b, bType := FindCanBeatCards(robot.HandCards, eCards, eType)
		if b {
			// 这里要不要把炸弹移除了在判断
			_, godNum := CheckRealGodCard(comRe, nextPlayer.HandCards, lastPlayer.HandCards)
			if allLen-godNum <= 2 || !(bType == cardConst.CARD_PATTERN_BOMB || bType == cardConst.CARD_PATTERN_ROCKET) { // 这里等于2就可以炸了
				OutCardsAction(room, robot, nextPlayer, cards, bType)
				return
			}
		}

	}

	// 以下逻辑正常跟牌
	beatCards, b, bType := FindMinFollowCards(robot.HandCards, compG, eCards, eType)
	if b {
		tmpG := CreateGroupCard(removeCards(robot.HandCards, beatCards))
		tmpComG := completeGroupCard(tmpG)
		tmpComRe := changeGroupToReCard(tmpComG)
		tmpAllLen := len(tmpComRe)
		_, godNum := CheckRealGodCard(tmpComRe, nextPlayer.HandCards, lastPlayer.HandCards)
		// 如果返回了一个炸弹出牌
		if bType == cardConst.CARD_PATTERN_BOMB || bType == cardConst.CARD_PATTERN_ROCKET {
			if tmpAllLen-godNum <= 1 {
				OutCardsAction(room, robot, nextPlayer, beatCards, bType)
				return
			}
		} else { // 如果返回的不是炸弹 拆牌之后 手术牌值多一首 可以出
			if allLen+2 >= tmpAllLen {
				OutCardsAction(room, robot, nextPlayer, beatCards, bType)
				return
			}
		}

	}


	// 2020年2月21日19:31:21 这里对拆对子打单牌进行处理
	if eType==cardConst.CARD_PATTERN_SINGLE&&robot.LastAction==playerAction.NotOutCardAction{
		// 如果这首是单牌 并且上把牌地主也没出的情况下 60%的几率拆牌压
		rand:=RandNum(0,10)
		if rand>=4 {
			beatCards, b3, bType := FindCanBeatCards(robot.HandCards, eCards, eType)
			if  b3{
				OutCardsAction(room, robot, nextPlayer, beatCards, bType)
				return
			}
		}
	}
	// 2020年2月21日19:31:21 这里对拆对子打单牌进行处理

	NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
	return

}

/*
地主下家农民首出
	是否 all-godNum<=1
		是 【随机非炸弹天牌】
	下家是否报双
		是  手中是否有天炸
			是 手中除炸弹最小对子牌能否让报双牌打过
				是 【出天炸】
			否 手中除炸弹最小对牌能否让保双牌打过
				是 【手中最小权重对子牌】
		否 手中除炸弹最小牌能否让保双牌打过
			是 【手中最小权重对子牌】
	下家是否报单
		是  手中是否有天炸
			是 手中除炸弹最小牌能否让保单牌打过
				是 【出天炸】
			否 手中除炸弹最小牌能否让保单牌打过
				是 【手中最小权重单牌】
		否 手中除炸弹最小牌能否让保单牌打过
			是 【手中最小权重单牌】
	【最小权重牌】
	地主是否保单
		是 ：手牌中是否只有单牌
			是 ： 【最大单牌】
			否 ： 【非单牌最小权重】
	地主是否保双
		是 ：是否只有对子
			是 ：【最小对子】
			否： 【非对子最小权重】
*/
func NewRobotFarmerMustDoF1(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {

	// 2020年2月21日16:41:43 如果能一首出完
	hc := robot.HandCards
	cardsType := GetCardsType(hc)
	if cardsType>=cardConst.CARD_PATTERN_SINGLE&&cardsType<=cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_PAIRS {
		OutCardsAction(room, robot, nextPlayer, hc, cardsType)
		return
	}
	// 2020年2月21日16:41:43 如果能一首出完

	comG := completeGroupCard(robot.GroupCard)
	comRe := changeGroupToReCard(comG)
	SortReCardByWightSL(comRe)
	allLen := len(comRe) // 出牌总首数
	upReCards, godNum := CheckGodCard(comRe, nextPlayer.HandCards, lastPlayer.HandCards)
	// 如果最后两手 有一首是炸弹 并且 无人保单或者无人报双的情况下 先走非炸弹的
	nFarmerType, isNFLastOne := checkPlayerHasLast(nextPlayer)    // 农民
	landType, isLandlordLastOne := checkPlayerHasLast(lastPlayer) // 地主
	if allLen == 2 {
		if !isNFLastOne && !isLandlordLastOne {
			if comRe[0].CardType == cardConst.CARD_PATTERN_BOMB ||
				comRe[1].CardType == cardConst.CARD_PATTERN_BOMB {
				for i := 0; i < len(upReCards); i++ {
					if comRe[i].CardType != cardConst.CARD_PATTERN_BOMB {
						logger.Debug("F1首出 最后一首是炸弹+其他 并无人报牌 显出其他")
						outCard := comRe[i].Card
						cardType := comRe[i].CardType
						OutCardsAction(room, robot, nextPlayer, SpecialHandle(outCard), cardType)
						return
					}
				}
			}
		}
	}

	// 这里如果总首数-天牌数<=1 // 则先出非炸弹天牌
	if allLen-godNum <= 1 {
		for i := 0; i < len(upReCards); i++ {
			if upReCards[i].IsGodCard && !(upReCards[i].RC.CardType == cardConst.CARD_PATTERN_BOMB || upReCards[i].RC.CardType == cardConst.CARD_PATTERN_ROCKET) {
				outCard := upReCards[i].RC.Card
				cardType := upReCards[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, SpecialHandle(outCard), cardType)
				return
			}
		}

		// 这里如果总首数-天牌数<=1 如果
		for i := 0; i < len(upReCards); i++ {
			if upReCards[i].IsGodCard {
				outCard := upReCards[i].RC.Card
				cardType := upReCards[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, SpecialHandle(outCard), cardType)
				return
			}
		}
	}

	// 如果下家农民最后一首单
	if isNFLastOne && nFarmerType == cardConst.CARD_PATTERN_SINGLE {
		minCard := findMinCard(robot.HandCards)

		// 寻找天炸
		godBomb, b, bType := findGodBomb(upReCards)
		if b {
			if CanBeat(minCard, removeCards(robot.HandCards, godBomb)) { // 如果有天炸 且能让下家打过最小单牌
				OutCardsAction(room, robot, nextPlayer, godBomb, bType)
				return
			}
		} else {
			if CanBeat(minCard, robot.HandCards) { // 能让下家打过最小单牌
				OutCardsAction(room, robot, nextPlayer, minCard, cardConst.CARD_PATTERN_SINGLE)
				return
			}
		}
	}

	// 如果下家农民最后一首对
	if isNFLastOne && nFarmerType == cardConst.CARD_PATTERN_PAIR {
		minCard := findMinDoubleCard(robot.HandCards)

		// 寻找天炸
		godBomb, b, bType := findGodBomb(upReCards)
		if b {
			if CanBeat(minCard, removeCards(robot.HandCards, godBomb)) { // 如果有天炸 且能让下家打过最小对子
				OutCardsAction(room, robot, nextPlayer, godBomb, bType)
				return
			}
		} else {
			if CanBeat(minCard, robot.HandCards) { // 能让下家打过最小对子
				OutCardsAction(room, robot, nextPlayer, minCard, cardConst.CARD_PATTERN_PAIR)
				return
			}
		}
	}

	// 如果地主保单
	if isLandlordLastOne && landType == cardConst.CARD_PATTERN_SINGLE {
		// 如果地主报单
		for i := 0; i < len(comRe); i++ {
			if comRe[i].CardType != landType {
				logger.Debug("F1首出 地主保单 非类型最小权重")
				outCard := comRe[i].Card
				cardType := comRe[i].CardType
				OutCardsAction(room, robot, nextPlayer, SpecialHandle(outCard), cardType)
				return
			}
		}

		logger.Debug("F1首出 只有和这种类型相似的牌  出第最大单排")
		OutCardsAction(room, robot, nextPlayer, comRe[len(comRe)-1].Card, comRe[len(comRe)-1].CardType)
		return
	}

	// 如果地主报双
	if isLandlordLastOne && landType == cardConst.CARD_PATTERN_PAIR {
		// 如果地主报双
		for i := 0; i < len(comRe); i++ {
			if comRe[i].CardType != landType {
				logger.Debug("F1首出 地主保双 非类型最小权重")
				outCard := comRe[i].Card
				cardType := comRe[i].CardType
				OutCardsAction(room, robot, nextPlayer, SpecialHandle(outCard), cardType)
				return
			}
		}

		logger.Debug("F1首出 如果全是对子 出最小单牌")
		outCard := findMinCard(robot.HandCards)
		OutCardsAction(room, robot, nextPlayer, outCard, cardConst.CARD_PATTERN_SINGLE)
		return
	}

	logger.Debug("F1首出 无人最后一手牌 取最小权值牌")
	OutCardsAction(room, robot, nextPlayer, SpecialHandle(comRe[0].Card), comRe[0].CardType)
}

/*
地主上家农民首出
	是否 all-godNum<=1
		是 【随机非炸弹天牌】
	下家是否报双
		是  手中是否有天炸
			是 手中除炸弹最小对子牌能否让报双牌打过
				是 【出天炸】
			否 手中除炸弹最小对牌能否让保双牌打过
				是 【手中最小权重对子牌】
		否 手中除炸弹最小牌能否让保双牌打过
			是 【手中最小权重对子牌】
	下家是否报单
		是  手中是否有天炸
			是 手中除炸弹最小牌能否让保单牌打过
				是 【出天炸】
			否 手中除炸弹最小牌能否让保单牌打过
				是 【手中最小权重单牌】
		否 手中除炸弹最小牌能否让保单牌打过
			是 【手中最小权重单牌】
	【最小权重牌】
	地主是否保单
		是 ：手牌中是否只有单牌
			是 ： 【最大单牌】
			否 ： 【非单牌最小权重】
	地主是否保双
		是 ：是否只有对子
			是 ：【最小对子】
			否： 【非对子最小权重】
*/
func NewRobotFarmerMustDoF2(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {

	// 2020年2月21日16:41:43 如果能一首出完
	hc := robot.HandCards
	cardsType := GetCardsType(hc)
	if cardsType>=cardConst.CARD_PATTERN_SINGLE&&cardsType<=cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_PAIRS {
		OutCardsAction(room, robot, nextPlayer, hc, cardsType)
		return
	}
	// 2020年2月21日16:41:43 如果能一首出完

	comG := completeGroupCard(robot.GroupCard)
	comRe := changeGroupToReCard(comG)
	SortReCardByWightSL(comRe)
	allLen := len(comRe) // 出牌总首数
	upReCards, godNum := CheckGodCard(comRe, nextPlayer.HandCards, lastPlayer.HandCards)
	// 如果最后两手 有一首是炸弹 并且 无人保单或者无人报双的情况下 先走非炸弹的
	landType, isLandlordLastOne := checkPlayerHasLast(nextPlayer) // 地主
	if allLen == 2 {
		if !isLandlordLastOne {
			if comRe[0].CardType == cardConst.CARD_PATTERN_BOMB ||
				comRe[1].CardType == cardConst.CARD_PATTERN_BOMB {
				for i := 0; i < len(upReCards); i++ {
					if comRe[i].CardType != cardConst.CARD_PATTERN_BOMB {
						logger.Debug("F2首出 最后一首是炸弹+其他 并无人报牌 显出其他")
						outCard := comRe[i].Card
						cardType := comRe[i].CardType
						OutCardsAction(room, robot, nextPlayer, SpecialHandle(outCard), cardType)
						return
					}
				}
			}
		}
	}

	// 这里如果总首数-天牌数<=1 // 则先出非炸弹天牌
	if allLen-godNum <= 1 {
		for i := 0; i < len(upReCards); i++ {
			if upReCards[i].IsGodCard && !(upReCards[i].RC.CardType == cardConst.CARD_PATTERN_BOMB || upReCards[i].RC.CardType == cardConst.CARD_PATTERN_ROCKET) {
				outCard := upReCards[i].RC.Card
				cardType := upReCards[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer, SpecialHandle(outCard), cardType)
				return
			}
		}

		// 这里如果总首数-天牌数<=1 如果
		for i := 0; i < len(upReCards); i++ {
			if upReCards[i].IsGodCard {
				outCard := upReCards[i].RC.Card
				cardType := upReCards[i].RC.CardType
				OutCardsAction(room, robot, nextPlayer,  SpecialHandle(outCard), cardType)
				return
			}
		}
	}

	// 如果地主保单
	if isLandlordLastOne && landType == cardConst.CARD_PATTERN_SINGLE {
		// 如果地主报单
		for i := 0; i < len(comRe); i++ {
			if comRe[i].CardType != landType {
				logger.Debug("F1首出 地主保单 非类型最小权重")
				outCard := comRe[i].Card
				cardType := comRe[i].CardType
				OutCardsAction(room, robot, nextPlayer,  SpecialHandle(outCard), cardType)
				return
			}
		}

		logger.Debug("F1首出 只有和这种类型相似的牌  出第最大单排")
		OutCardsAction(room, robot, nextPlayer, comRe[len(comRe)-1].Card, comRe[len(comRe)-1].CardType)
		return
	}

	// 如果地主报双
	if isLandlordLastOne && landType == cardConst.CARD_PATTERN_PAIR {
		// 如果地主报双
		for i := 0; i < len(comRe); i++ {
			if comRe[i].CardType != landType {
				logger.Debug("F1首出 地主保双 非类型最小权重")
				outCard := comRe[i].Card
				cardType := comRe[i].CardType
				OutCardsAction(room, robot, nextPlayer,  SpecialHandle(outCard), cardType)
				return
			}
		}

		logger.Debug("F1首出 如果全是对子 出最小单牌")
		outCard := findMinCard(robot.HandCards)
		OutCardsAction(room, robot, nextPlayer, outCard, cardConst.CARD_PATTERN_SINGLE)
		return
	}

	logger.Debug("F1首出 无人最后一手牌 取最小权值牌")
	OutCardsAction(room, robot, nextPlayer,  SpecialHandle(comRe[0].Card), comRe[0].CardType)
}

/*
	F1跟牌
	nextPlayer 农民
	lastPlayer 地主

谁出的牌
	农民1号
		是否保单
			是否有天炸
				有： 是否有最小牌能让下家农民打过
					是：【天炸】
					否：【不出】
				否：【不出】
			否，判断一下条件
				如果当前权重超过11 【不出】
				强行跟牌 首数增多 【不出】
				当前首牌多于一号农民【不出】
				【最小跟牌】
	地主出的
		地主是否报单
			是：有出必出
			否 最小跟牌 



*/
func farmerFallowF1(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	logger.Debug("F1跟牌 ====================")
	ePlayer := room.Players[room.EffectivePlayerId]
	eType := room.EffectiveType
	eCards := room.EffectiveCard

	compG := completeGroupCard(robot.GroupCard)
	comRe := changeGroupToReCard(compG)
	SortReCardByWightSL(comRe)
	upReCards, godNum := CheckRealGodCard(comRe, nextPlayer.HandCards, lastPlayer.HandCards)
	allLen := len(comRe)
	followCards, canBeat, oType := FindMinFollowCards(robot.HandCards, compG, eCards, eType)
	if canBeat {
		logger.Debug("F1跟牌  有能大过的跟牌====================")
		// 1. 天派数量大于
		if allLen-godNum <= 1 {
			OutCardsAction(room, robot, nextPlayer, followCards, oType)
			return
		}

		// 如果农民保单
		lastType, b := checkPlayerHasLast(nextPlayer) // 判断下架农民是否最后一首牌
		if b && lastType == cardConst.CARD_PATTERN_SINGLE { // 是否保单
			// 寻找天炸
			logger.Debug("F1跟牌  有能大过的跟牌  上家是出牌是农民 农民保单====================")
			godBomb, hasGodBomb, bType := findGodBomb(upReCards)
			if hasGodBomb {
				logger.Debug("F1跟牌  有能大过的跟牌  上家是出牌是农民 有天炸====================")
				minCard := findMinCard(removeCards(robot.HandCards, godBomb))
				if CanBeat(minCard, nextPlayer.HandCards) {
					OutCardsAction(room, robot, nextPlayer, godBomb, bType)
					return
				}

			}
		}

		// 2.出牌者是农民
		if !ePlayer.IsLandlord == true { // 如果是下家农民出的
			logger.Debug("F1跟牌  有能大过的跟牌  上家是出牌是农民====================")
			//nextCompG := completeGroupCard(lastPlayer.GroupCard)
			//nextComRe := changeGroupToReCard(nextCompG)
			if !(countCardWight(followCards, oType) >= cardConst.CARD_RANK_ACE){//||
				//len(nextComRe)+1 >= allLen) {
				tmpG := CreateGroupCard(removeCards(robot.HandCards, followCards))
				tmpComG := completeGroupCard(tmpG)
				tmpComRe := changeGroupToReCard(tmpComG)
				tmpAllLen := len(tmpComRe)
				_, godNum := CheckRealGodCard(tmpComRe, nextPlayer.HandCards, lastPlayer.HandCards)
				// 如果返回了一个炸弹出牌
				if oType == cardConst.CARD_PATTERN_BOMB || oType == cardConst.CARD_PATTERN_ROCKET {
					if tmpAllLen-godNum <= 1 {
						OutCardsAction(room, robot, nextPlayer, followCards, oType)
						return
					}
				} else { // 如果返回的不是炸弹 拆牌之后 手术牌值多一首 可以出
					if allLen > tmpAllLen {
						OutCardsAction(room, robot, nextPlayer, followCards, oType)
						return
					}
				}
			}

		} else { // 上首牌是地主出的
			// 以下逻辑正常跟牌
			tmpG := CreateGroupCard(removeCards(robot.HandCards, followCards))
			tmpComG := completeGroupCard(tmpG)
			tmpComRe := changeGroupToReCard(tmpComG)
			tmpAllLen := len(tmpComRe)
			_, godNum := CheckRealGodCard(tmpComRe, nextPlayer.HandCards, lastPlayer.HandCards)
			// 如果返回了一个炸弹出牌
			if oType == cardConst.CARD_PATTERN_BOMB || oType == cardConst.CARD_PATTERN_ROCKET {
				if tmpAllLen-godNum <= 1 {
					OutCardsAction(room, robot, nextPlayer, followCards, oType)
					return
				}
			} else if allLen >= tmpAllLen { // 如果返回的不是炸弹 拆牌之后 手术牌值多一首 可以出
				OutCardsAction(room, robot, nextPlayer, followCards, oType)
				return
			} else {
				OutCardsAction(room, robot, nextPlayer, followCards, oType)
				return
			}
		}
	}

	logger.Debug("F1跟牌 不出")
	NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
	return
}

func farmerFallowF2(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	logger.Debug("F2跟牌====================")
	ePlayer := room.Players[room.EffectivePlayerId]
	eType := room.EffectiveType
	eCards := room.EffectiveCard

	compG := completeGroupCard(robot.GroupCard)
	comRe := changeGroupToReCard(compG)
	SortReCardByWightSL(comRe)
	_, godNum := CheckRealGodCard(comRe, nextPlayer.HandCards, lastPlayer.HandCards)
	allLen := len(comRe)
	followCards, canBeat, oType := FindMinFollowCards(robot.HandCards, compG, eCards, eType)

	if canBeat {
		if allLen-godNum <= 1 {
			OutCardsAction(room, robot, nextPlayer, followCards, oType)
			return
		}
		if ePlayer.IsLandlord == true { // 如果有效牌是地主出的
			lastType, b := checkPlayerHasLast(nextPlayer) // 判断地主否最后一首牌
			if b && lastType == cardConst.CARD_PATTERN_SINGLE {
				logger.Debug("F2跟牌  有能大过的跟牌  上家是出牌是地主 地主保单")
				bomb, hasBomb, bType := FindCanBeatCards(robot.HandCards,
					[]*Card{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
					cardConst.CARD_PATTERN_BOMB)
				if hasBomb {
					single := checkHowManyBeatSingle(compG.Single, nextPlayer.HandCards)
					if single <= 1 {
						OutCardsAction(room, robot, nextPlayer, bomb, bType)
						return
					}

				} else { // 寻找最大单牌 单牌中的最大牌
					// todo 这里要顶牌啊
					OutCardsAction(room, robot, nextPlayer, followCards, oType)
					return
				}
			} else {
				OutCardsAction(room, robot, nextPlayer, followCards, oType)
				return
			}

		} else { // 如果是下家农民出的
			lastType, b := checkPlayerHasLast(nextPlayer) // 判断地主否最后一首牌
			if b && lastType == cardConst.CARD_PATTERN_SINGLE { // 是否保单
				// 寻找天炸
				logger.Debug("F2跟牌")
				bomb, hasBomb, bType := FindCanBeatCards(robot.HandCards,
					[]*Card{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
					cardConst.CARD_PATTERN_BOMB)
				if hasBomb {
					single := checkHowManyBeatSingle(compG.Single, nextPlayer.HandCards)
					if single <= 1 {
						OutCardsAction(room, robot, nextPlayer, bomb, bType)
						return
					}

				} else { // 寻找最大单牌 单牌中的最大牌
					//todo 有瑕疵
					maxSingle := findMaxSingle(robot.HandCards)
					if CanBeat(eCards, maxSingle) {
						OutCardsAction(room, robot, nextPlayer, maxSingle, cardConst.CARD_PATTERN_SINGLE)
						return
					}
				}
			}

			logger.Debug("F2跟牌  有能大过的跟牌  上家是出牌是农民====================")
			//lastCompG := completeGroupCard(lastPlayer.GroupCard)
			//nextComRe := changeGroupToReCard(lastCompG)
			if !(countCardWight(followCards, oType) >= cardConst.CARD_RANK_ACE ){
				//len(nextComRe)+1 >= allLen) || nextPlayer.LastAction == playerAction.OutCardAction {
				// 最小跟拍

				tmpG := CreateGroupCard(removeCards(robot.HandCards, followCards))
				tmpComG := completeGroupCard(tmpG)
				tmpComRe := changeGroupToReCard(tmpComG)
				tmpAllLen := len(tmpComRe)
				_, godNum := CheckRealGodCard(tmpComRe, nextPlayer.HandCards, lastPlayer.HandCards)
				// 如果返回了一个炸弹出牌
				if oType == cardConst.CARD_PATTERN_BOMB || oType == cardConst.CARD_PATTERN_ROCKET {
					if tmpAllLen-godNum <= 1 {
						OutCardsAction(room, robot, nextPlayer, followCards, oType)
						return
					}
				} else { // 如果返回的不是炸弹 拆牌之后 手术牌值多一首 可以出

						OutCardsAction(room, robot, nextPlayer, followCards, oType)
						return
				}
			}
		}
	}
	logger.Debug("F2 跟牌 不出")
	NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
	return

}



// bug特殊处理
/*
	3456788大鬼能出牌 123456615
	3456788大鬼出牌 123456614
*/
func SpecialHandle(outCard []*Card)[]*Card  {
	if len(outCard)!=8 {
		return outCard
	}
	var exp string
	SortCardSL(outCard)
	for i:=0; i<len(outCard);i++  {
		exp+=strconv.Itoa(int(outCard[i].Value))
	}

	var res []*Card
	if strings.Contains(exp,"123456")&&GetCardsType(outCard)==cardConst.CARD_PATTERN_ERROR{
		for i:=1;i<=6 ;  i++{
			card:=findThisValueCard(i,outCard,1)
			res = append(res, card...)
		}
		return res
	}else {
		return outCard
	}
}