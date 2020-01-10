package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/cardConst"
)

/* ==============================  help ===========================*/

/*
 	为3张寻找带牌
	从已有牌中找出最小的带牌
	从单张和对子中寻找带牌数量

	1.首先判断单牌的张数  和 对子的数量 如果单张数量已经满足 则优先用单张里面抽取
	大部分从单张带起


	如果张数一样 先带单张 如果单张中 包含大王 并且对子数量为0  则允许 反之带最小对子

	// 如果单张和对子的数量都为0  并且刚好有连对长度与之吻合 带上
*/
func findTripleWithCards(singles, double, junkoDouble []*ReCard, length int) []*Card {
	var withCards []*Card
	// 1.优先寻找单排
	if len(singles) >= length {
		for i := 0; i < length; i++ {
			withCards = append(withCards, singles[i].Card...)
		}
		return withCards
	}

	// todo 如果此时 单张刚好满足 且包含 2 或者王这样的大牌 并且 对子也满足且很小的情况下 则还是带对子比较好
	// 2.寻找对子
	if len(double) >= length {
		for i := 0; i < length; i++ {
			withCards = append(withCards, double[i].Card...)
		}
		return withCards
	}

	// 3.把单牌和对子混合 取最小张  如果总长度大于 的话
	var mix []*Card
	for i := 0; i < len(singles); i++ {
		mix = append(mix, singles[i].Card...)
	}
	for i := 0; i < len(double); i++ {
		mix = append(mix, double[i].Card...)
	}

	if len(mix) >= length {
		SortCardSL(mix)
		for i := 0; i < length; i++ {
			withCards = append(withCards, mix[i])
		}

		return withCards
	}

	// 如果这都还没有 连对的吻合 则带连对
	if len(junkoDouble) == length {
		var final []*Card
		for i := 0; i < length; i++ {
			final = append(final, junkoDouble[i].Card...)
		}

		return final
	}
	return nil
}

/*
	todo 待优化
	判断农民是否报单或者报双
	场上有单 就返报单类型
	双单  单
	一单一双 单
	双对  双

*/
func checkLandlordHasLast(landlord *Player) ([]*Card, int32, bool) {

	cardsType1 := GetCardsType(landlord.HandCards)
	if cardsType1 == cardConst.CARD_PATTERN_PAIR || cardsType1 == cardConst.CARD_PATTERN_SINGLE {
		return landlord.HandCards, cardsType1, true
	}

	return nil, 0, false
}

// 计算牌型权重
func countCardWight(cards []*Card, cType int32) int32 {
	switch cType {
	case cardConst.CARD_PATTERN_SINGLE:
		return cards[0].Value
	case cardConst.CARD_PATTERN_PAIR:
		return cards[0].Value
	case cardConst.CARD_PATTERN_TRIPLET:
		return cards[0].Value
	case cardConst.CARD_PATTERN_TRIPLET_WITH_SINGLE:
		return getCardManyOne(cards)
	case cardConst.CARD_PATTERN_TRIPLET_WITH_PAIR:
		return getCardManyOne(cards)
	case cardConst.CARD_PATTERN_BOMB:
		return getCardManyOne(cards)
	case cardConst.CARD_PATTERN_QUADPLEX_WITH_SINGLES: // 四代二单
		return getCardManyOne(cards)
	case cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS: // 四代二对
		return getCardManyOne(cards)
	case cardConst.CARD_PATTERN_SEQUENCE:
		SortCard(cards)
		return cards[len(cards)-1].Value
	case cardConst.CARD_PATTERN_SEQUENCE_OF_PAIRS: // 连对
		SortCard(cards)
		return cards[len(cards)-1].Value
	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS: // 飞机不带
		SortCard(cards)
		return cards[len(cards)-1].Value
	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_SINGLES: // 飞机带单
		return getTripletWeight(cards)
	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_PAIRS: // 飞机带对
		return getTripletWeight(cards)
	case cardConst.CARD_PATTERN_ROCKET: // 飞机带对
		return 1000
	}

	return 0

}

//
//// 最小跟牌
///*
//	最小跟牌 及判断 自己是否有能打过上家与之匹配的牌型 不出炸弹
//
//*/
//func minFollowCard(actionPlayer *Player, eCards []*Card, eType int32) ([]*Card, bool) {
//
//	eWight := countCardWight(eCards, eType)
//	group := actionPlayer.GroupCard
//
//	switch eType {
//	case cardConst.CARD_PATTERN_SINGLE:
//		g := group.Single
//		if g != nil {
//			for i := 0; i < len(g); i++ {
//				if g[i].Wight > eWight {
//					return g[i].Card, true
//				}
//			}
//		}
//		return nil, false
//	case cardConst.CARD_PATTERN_PAIR:
//		g := group.Double
//		if g != nil {
//			for i := 0; i < len(g); i++ {
//				if g[i].Wight > eWight {
//					return g[i].Card, true
//				}
//			}
//		}
//		return nil, false
//
//	case cardConst.CARD_PATTERN_TRIPLET:
//		g := group.Triple
//		if g != nil {
//			for i := 0; i < len(g); i++ {
//				if g[i].Wight > eWight {
//					return g[i].Card, true
//				}
//			}
//		}
//		//cards, b, _ := HostingBeatTriple(actionPlayer.HandCards, eCards)
//
//		return nil, false
//
//	case cardConst.CARD_PATTERN_TRIPLET_WITH_SINGLE:
//		g := group.Triple
//		if g != nil {
//			for i := 0; i < len(g); i++ {
//				if g[i].Wight > eWight {
//					singles := group.Single
//					if len(singles) <= 0 { //||len { //
//						return nil, false
//					}
//					var r []*Card
//					r = append(r, singles[0].Card...)
//					r = append(r, g[i].Card...)
//					return r, true
//				}
//			}
//		}
//		return nil, false
//
//	case cardConst.CARD_PATTERN_TRIPLET_WITH_PAIR:
//		g := group.Triple
//		if g != nil {
//			for i := 0; i < len(g); i++ {
//				if g[i].Wight > eWight {
//					Double := group.Double
//					if len(Double) <= 0 { //||len { //
//						return nil, false
//					}
//					var r []*Card
//					r = append(r, Double[0].Card...)
//					r = append(r, g[i].Card...)
//					return r, true
//				}
//			}
//		}
//		return nil, false
//
//	case cardConst.CARD_PATTERN_BOMB:
//		// todo 这里先返回nil
//		return nil, false
//
//	case cardConst.CARD_PATTERN_QUADPLEX_WITH_SINGLES: // 四代二单
//		// todo 这里先返回nil
//		return nil, false
//
//	case cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS: // 四代二对
//		// todo 这里先返回nil
//		return nil, false
//
//	case cardConst.CARD_PATTERN_SEQUENCE:
//		g := group.Junko
//		if g != nil {
//			for i := 0; i < len(g); i++ {
//				if g[i].Wight > eWight && len(g[i].Card) == len(eCards) {
//					return g[i].Card, true
//				}
//			}
//		}
//		return nil, false
//
//	case cardConst.CARD_PATTERN_SEQUENCE_OF_PAIRS: // 连对
//		g := group.JunkoDouble
//		if g != nil {
//			for i := 0; i < len(g); i++ {
//				if g[i].Wight > eWight && len(g[i].Card) == len(eCards) {
//					return g[i].Card, true
//				}
//			}
//		}
//		return nil, false
//
//	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS: // 飞机不带
//		g := group.JunkTriple
//		if g != nil {
//			for i := 0; i < len(g); i++ {
//				if g[i].Wight > eWight && len(g[i].Card) == len(eCards) {
//					return g[i].Card, true
//				}
//			}
//		}
//		return nil, false
//
//	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_SINGLES: // 飞机带单
//		g := group.Triple
//		if g != nil {
//			for i := 0; i < len(g); i++ {
//				if g[i].Wight > eWight && len(g[i].Card) == len(eCards) {
//					singles := group.Single
//					if len(singles) < len(eCards)/4 { //||len { //
//						return nil, false
//					}
//					var r []*Card
//					for j := 0; j < len(eCards)/4; j++ {
//						r = append(r, singles[i].Card...)
//					}
//					r = append(r, g[i].Card...)
//					return r, true
//				}
//			}
//		}
//		return nil, false
//
//	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_PAIRS: // 飞机带对
//		g := group.Triple
//		if g != nil {
//			for i := 0; i < len(g); i++ {
//				if g[i].Wight > eWight && len(g[i].Card) == len(eCards) {
//					Double := group.Double
//					if len(Double) < len(eCards)/4 { //||len { //
//						return nil, false
//					}
//					var r []*Card
//					for j := 0; j < len(eCards)/4; j++ {
//						r = append(r, Double[i].Card...)
//					}
//					r = append(r, g[i].Card...)
//					return r, true
//				}
//			}
//		}
//		return nil, false
//
//	case cardConst.CARD_PATTERN_ROCKET: // 火箭
//		return nil, false
//	}
//	return nil, false
//}

// 获取飞机的权值
func getTripletWeight(cards []*Card) int32 {
	group := CreateGroupCard(cards)
	if len(group.Triple) <= 0 || group.Triple == nil {
		logger.Error("！！！！ 错误的飞机类型！！！")
		//	return 13
	}
	return group.Triple[0].Wight
}

/*
	将玩家的组牌加上带牌重新组牌
*/
func completeGroupCard(g GroupCard) GroupCard {
	singles := g.Single
	doubles := g.Double

	var sAndDCards []*Card
	for i := 0; i < len(singles); i++ {
		sAndDCards = append(sAndDCards, singles[i].Card...)
	}

	for i := 0; i < len(doubles); i++ {
		sAndDCards = append(sAndDCards, doubles[i].Card...)
	}

	SortCardSL(sAndDCards)

	// 飞机
	triples := g.JunkTriple
	for i := 0; i < len(triples); i++ {
		tLen := len(triples[i].Card) / 3
		if len(singles) >= tLen {
			var withCards []*Card
			for j := 0; j < tLen; j++ {
				withCards = append(withCards, singles[j].Card...)
			}
			triples[i].Card = append(triples[i].Card, withCards...)
			triples[i].CardType = cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_SINGLES
			sAndDCards = removeCards(sAndDCards, withCards)
			singles = singles[tLen:]
			continue
		}
		if len(doubles) >= tLen {
			var withCards []*Card
			for j := 0; j < tLen; j++ {
				withCards = append(withCards, doubles[j].Card...)
			}
			triples[i].Card = append(triples[i].Card, withCards...)
			triples[i].CardType = cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_PAIRS
			sAndDCards = removeCards(sAndDCards, withCards)
			doubles = doubles[tLen:]
			continue
		}

		// 不够组合 取最小
		if len(sAndDCards) >= tLen {
			logger.Debug("合并去min")
			triples[i].Card = append(triples[i].Card, sAndDCards[:tLen]...)
			triples[i].CardType = cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_SINGLES
			sAndDCards = sAndDCards[tLen:]

			tmpS, _ := FindAllSingle(sAndDCards)
			singles = tmpS

			tmpD, _ := FindAllDouble(sAndDCards)
			doubles = tmpD
			//
			//for i := 0; i < len(doubles); i++ {
			//	PrintCard(doubles[i].Card)
			//}
		}

	}

	triple := g.Triple
	sAndDReCards := append(singles, doubles...)
	SortReCardByWightSL(sAndDReCards)
	for i := 0; i < len(triple); i++ {
		// 取单张和对子的最小权重
		if len(sAndDReCards) >= 1 {
			withCards := sAndDReCards[0].Card
			triple[i].CardType = sAndDReCards[0].CardType
			triple[i].Card = append(triple[i].Card, withCards...)
			//logger.Debug("此处type:", sAndDReCards[0].CardType)
			//logger.Debug("此处type:", triple[i].CardType)

			sAndDCards = removeCards(sAndDCards, withCards)
			sAndDReCards = sAndDReCards[1:]
			tmpS, _ := FindAllSingle(sAndDCards)
			singles = tmpS
			tmpD, _ := FindAllDouble(sAndDCards)
			doubles = tmpD
		}
	}

	g.Single = singles
	g.Double = doubles
	return g
}

/*
	将group 装换成 []*Recard
*/

func changeGroupToReCard(g GroupCard) []*ReCard {
	var result []*ReCard
	result = append(result, g.Single...)
	result = append(result, g.Double...)
	result = append(result, g.Triple...)
	result = append(result, g.Bomb...)
	result = append(result, g.Junko...)
	result = append(result, g.JunkoDouble...)
	result = append(result, g.JunkTriple...)
	result = append(result, g.Rocket...)

	return result
}

// 判断玩家是否保单或者报双

/*
	判断农民是否报单或者报双
	场上有单 就返报单类型
	双单  单
	一单一双 单
	双对  双

*/
func checkPlayerHasLast(player *Player) (int32, bool) {

	cardsType1 := GetCardsType(player.HandCards)
	if cardsType1 == cardConst.CARD_PATTERN_PAIR || cardsType1 == cardConst.CARD_PATTERN_SINGLE {
		return cardsType1, true
	}
	return 0, false
}

// 寻找牌中最小的一张牌
func findMinCard(hands []*Card) []*Card {
	if len(hands) <= 0 {
		return nil
	}
	SortCardSL(hands)
	return append([]*Card{}, hands[0])
}

// 寻找牌中最大的一张牌
func findMaxCard(hands []*Card) []*Card {
	if len(hands) <= 0 {
		return nil
	}
	SortCardSL(hands)
	return append([]*Card{}, hands[len(hands)-1])
}

// 寻找牌中最大的一张单牌
func findMaxSingle(hands []*Card) []*Card {
	if len(hands) <= 0 {
		return nil
	}
	SortCardSL(hands)
	return append([]*Card{}, hands[len(hands)-1])
}


// 寻找牌中最大的对子
// 有可能为空
func findMaxDouble(hands []*Card) []*Card {
	if len(hands) <= 0 {
		return nil
	}

	numsCard := getHasMoreNumsCard(hands, 2)
	if len(numsCard) > 0 {
		card := findThisValueCard(numsCard[len(numsCard)-1], hands, 2)
		return card
	}
	return nil
}

// 寻找牌中最小的一对牌 包括拆三带
// 这个方法和 findMinDouble() 不一样
func findMinDoubleCard(hands []*Card) []*Card {
	if len(hands) <= 0 {
		return nil
	}

	numsCard := getHasMoreNumsCard(hands, 2)
	if len(numsCard) > 0 {
		card := findThisValueCard(numsCard[0], hands, 2)
		return card
	}
	return nil
}

// 判断玩家天牌数
/*
	天牌: 指该手牌 不能被其他玩家(非炸弹)压的牌

	return
		[]*ReCard  : 标记是否天牌的lg
		int: 天牌数量
*/
func CheckGodCard(lg []*ReCard, p1Hands, p2Hands []*Card) ([]UpReCard, int) {
	var urs []UpReCard
	var godNums int
	for i := 0; i < len(lg); i++ {
		cardsType := GetCardsType(lg[i].Card)
		lg[i].CardType = cardsType

		if cardsType < cardConst.CARD_PATTERN_SINGLE || cardsType > cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS {
			logger.Debug("!!! impossible....", )
			PrintCard(lg[i].Card)
		}
		var f int
		var ur UpReCard
		ur.RC = lg[i]
		_, b, bt := FindCanBeatCards(p1Hands, lg[i].Card, cardsType)
		if !b && bt != cardConst.CARD_PATTERN_BOMB && bt != cardConst.CARD_PATTERN_ROCKET {
			f++
		}
		_, b1, bt2 := FindCanBeatCards(p2Hands, lg[i].Card, cardsType)
		if !b1 && bt2 != cardConst.CARD_PATTERN_BOMB && bt2 != cardConst.CARD_PATTERN_ROCKET {
			f++

		}
		if f == 2 {
			ur.IsGodCard = true
			godNums++
		}

		urs = append(urs, ur)
	}

	return urs, godNums
}

// 判断玩家天牌数
/*
	天牌: 指该手牌 不能被其他玩家(包括炸弹)压的牌

	return
		[]*ReCard  : 标记是否天牌的lg
		int: 天牌数量
*/
func CheckRealGodCard(lg []*ReCard, p1Hands, p2Hands []*Card) ([]UpReCard, int) {
	var urs []UpReCard
	var godNums int
	for i := 0; i < len(lg); i++ {
		cardsType := GetCardsType(lg[i].Card)
		lg[i].CardType = cardsType

		if cardsType < cardConst.CARD_PATTERN_SINGLE || cardsType > cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS {
			logger.Debug("!!! impossible....", cardsType)
			PrintCard(lg[i].Card)
		}
		var f int
		var ur UpReCard
		ur.RC = lg[i]
		_, b, _ := FindCanBeatCards(p1Hands, lg[i].Card, cardsType)
		if !b /*&& bt != cardConst.CARD_PATTERN_BOMB && bt != cardConst.CARD_PATTERN_ROCKET */ {
			f++
		}
		_, b1, _ := FindCanBeatCards(p2Hands, lg[i].Card, cardsType)
		if !b1 /*&& bt2 != cardConst.CARD_PATTERN_BOMB && bt2 != cardConst.CARD_PATTERN_ROCKET*/ {
			f++

		}
		if f == 2 {
			ur.IsGodCard = true
			godNums++
		}

		urs = append(urs, ur)
	}
	return urs, godNums
}

/*



 */
func checkHowManyBeatSingle(singles []*ReCard, single []*Card) int {
	var result int
	for i := 0; i < len(singles); i++ {
		if singles[i].Card[0].Value < single[0].Value {
			result++
		}

	}
	return result
}

/*
	寻找最小跟牌策略
	先从最有组合里面找
*/
func FindMinFollowCards(hands []*Card, completeG GroupCard, eCard []*Card, eType int32) ([]*Card, bool, int32) {
	var res []*Card
	switch eType {
	case cardConst.CARD_PATTERN_SINGLE: // 3
		SortReCardByWightSL(completeG.Single)
		for i := 0; i < len(completeG.Single); i++ {
			tmp := completeG.Single[i]
			if tmp != nil {
				if CanBeat(eCard, tmp.Card, ) {
					res = tmp.Card
					break
				}
			}
		}
	case cardConst.CARD_PATTERN_PAIR: // 4
		SortReCardByWightSL(completeG.Double)
		for i := 0; i < len(completeG.Double); i++ {
			tmp := completeG.Double[i]

			if tmp != nil {
				if CanBeat(eCard, tmp.Card, ) {
					res = tmp.Card
					break
				}
			}
		}
	case cardConst.CARD_PATTERN_TRIPLET: // 5
		SortReCardByWightSL(completeG.Triple)
		for i := 0; i < len(completeG.Triple); i++ {
			tmp := completeG.Triple[i]
			if tmp != nil {
				if CanBeat(eCard, tmp.Card, ) {
					res = tmp.Card
					break
				}
			}
		}
	case cardConst.CARD_PATTERN_TRIPLET_WITH_SINGLE: // 6
		SortReCardByWightSL(completeG.Triple)
		for i := 0; i < len(completeG.Triple); i++ {
			tmp := completeG.Triple[i]
			if tmp != nil {
				if CanBeat(eCard, tmp.Card, ) {
					res = tmp.Card
					break
				}
			}
		}
	case cardConst.CARD_PATTERN_TRIPLET_WITH_PAIR: // 7
		SortReCardByWightSL(completeG.Triple)
		for i := 0; i < len(completeG.Triple); i++ {
			tmp := completeG.Triple[i]
			if tmp != nil {
				if CanBeat(eCard, tmp.Card, ) {
					res = tmp.Card
					break
				}
			}
		}
	case cardConst.CARD_PATTERN_SEQUENCE: // 8
		SortReCardByWightSL(completeG.Junko)
		for i := 0; i < len(completeG.Junko); i++ {
			tmp := completeG.Junko[i]
			if tmp != nil {
				if CanBeat(eCard, tmp.Card, ) {
					res = tmp.Card
					break
				}
			}
		}
	case cardConst.CARD_PATTERN_SEQUENCE_OF_PAIRS: // 9
		SortReCardByWightSL(completeG.JunkoDouble)
		for i := 0; i < len(completeG.JunkoDouble); i++ {
			tmp := completeG.JunkoDouble[i]
			if tmp != nil {
				if CanBeat(eCard, tmp.Card, ) {
					res = tmp.Card
					break
				}
			}
		}

	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS: // 10
		SortReCardByWightSL(completeG.Triple)
		for i := 0; i < len(completeG.Triple); i++ {
			tmp := completeG.Triple[i]
			if tmp != nil {
				if CanBeat(eCard, tmp.Card, ) {
					res = tmp.Card
					break
				}
			}
		}

	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_SINGLES: // 11
		SortReCardByWightSL(completeG.Triple)
		for i := 0; i < len(completeG.Triple); i++ {
			tmp := completeG.Triple[i]
			if tmp != nil {
				if CanBeat(eCard, tmp.Card, ) {
					res = tmp.Card
					break
				}
			}
		}

	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_PAIRS: // 12
		SortReCardByWightSL(completeG.Triple)
		for i := 0; i < len(completeG.Triple); i++ {
			tmp := completeG.Triple[i]
			if tmp != nil {
				if CanBeat(eCard, tmp.Card, ) {
					res = tmp.Card
					break
				}
			}
		}
	case cardConst.CARD_PATTERN_BOMB: // 13
		SortReCardByWightSL(completeG.Bomb)
		for i := 0; i < len(completeG.Bomb); i++ {
			tmp := completeG.Bomb[i]
			if tmp != nil {
				if CanBeat(eCard, tmp.Card, ) {
					res = tmp.Card
					break
				}
			}
		}
	case cardConst.CARD_PATTERN_ROCKET: // 14
		return nil, false, -1

	case cardConst.CARD_PATTERN_QUADPLEX_WITH_SINGLES: // 15
		// todo 先不处理

	case cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS: // 16
		// todo 先不处理
	}

	if len(res) >= 0 && res != nil {
		return res, true, eType
	} else {
		return FindCanBeatCards(hands, eCard, eType)
	}

}

// 寻找天炸
func findGodBomb(upComRe []UpReCard) ([]*Card, bool, int32) {

	for i := 0; i < len(upComRe); i++ {
		if (upComRe[i].RC.CardType == cardConst.CARD_PATTERN_BOMB || upComRe[i].RC.CardType == cardConst.CARD_PATTERN_ROCKET) &&
			upComRe[i].IsGodCard {
			return upComRe[i].RC.Card, true, upComRe[i].RC.CardType
		}
	}
	return nil, false, 0
}

