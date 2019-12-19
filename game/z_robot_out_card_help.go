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
	从组牌中寻找权重最小的组牌
	// 注意顺子连对飞机 则权重要减去本身长度
*/
func findMinWeightCards(group GroupCard) *ReCard {
	var min ReCard
	min.Wight = 1000 // 初始值
	var all []*ReCard

	// 目前权重一样 根据以下顺序优先出
	all = append(all, group.Rocket...)
	all = append(all, group.Single...)
	all = append(all, group.Double...)
	all = append(all, group.Triple...)
	all = append(all, group.Junko...)
	all = append(all, group.JunkoDouble...)
	all = append(all, group.JunkTriple...)
	// all=append(all, group.Bomb...) 炸弹和火箭不考虑

	for i := 0; i < len(all); i++ {
		tmp := all[i]
		if tmp == nil {
			continue
		}
		tmpWight := tmp.CardType
		if tmp.CardType == cardConst.CARD_PATTERN_SEQUENCE ||
			tmp.CardType == cardConst.CARD_PATTERN_SEQUENCE_OF_PAIRS ||
			tmp.CardType == cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS {
			tmpWight = tmpWight - int32(len(tmp.Card))
		}

		if min.Wight < tmpWight {
			min = * tmp
		}
	}

	if min.Wight == 1000 {
		return nil
	}

	if min.CardType == cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS {
		findTripleWithCards(group.Single, group.Double, group.JunkoDouble, len(min.Card)/3)
	} else if min.CardType == cardConst.CARD_PATTERN_TRIPLET {
		findTripleWithCards(group.Single, group.Double, group.JunkoDouble, 1)
	}
	return &min

}

/*
	判断农民是否报单或者报双
	场上有单 就返报单类型
	双单  单
	一单一双 单
	双对  双

*/
func checkNotLandlordHasLast(farmer1, farmer2 *Player) ([]*Card, int32) {

	var result int32
	var resultCard []*Card
	//var resultHands []*Card
	cardsType1 := GetCardsType(farmer1.HandCards)
	if cardsType1 == cardConst.CARD_PATTERN_PAIR || cardsType1 == cardConst.CARD_PATTERN_SINGLE {
		result += result
		resultCard = farmer1.HandCards
	}

	cardsType2 := GetCardsType(farmer2.HandCards)
	if cardsType2 == cardConst.CARD_PATTERN_PAIR || cardsType2 == cardConst.CARD_PATTERN_SINGLE {
		result += result
		if cardsType1 == cardsType2 {
			if farmer2.HandCards[0].Value >= farmer1.HandCards[0].Value {
				resultCard = farmer2.HandCards
			}
		} else if cardsType1 == cardConst.CARD_PATTERN_PAIR && cardsType2 == cardConst.CARD_PATTERN_SINGLE {
			resultCard = farmer2.HandCards
		}
	}

	if result == 7 || result == 3 || result == 6 {
		result = cardConst.CARD_PATTERN_SINGLE

	} else if result == 8 || result == 4 {
		result = cardConst.CARD_PATTERN_PAIR
	}

	return resultCard, result
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

/*
	有人报单找牌
*/
func findMinWeightCardsExpectSome(group GroupCard, expSingle, expDouble bool) *ReCard {
	var min ReCard
	min.Wight = 1000 // 初始值
	var all []*ReCard

	// 目前权重一样 根据以下顺序优先出
	//all = append(all, group.Rocket...)

	if !expSingle {
		all = append(all, group.Single...)
	}
	if !expDouble {
		all = append(all, group.Double...)
	}
	all = append(all, group.Triple...)
	all = append(all, group.Junko...)
	all = append(all, group.JunkoDouble...)
	all = append(all, group.JunkTriple...)
	// all=append(all, group.Bomb...) 炸弹和火箭不考虑

	for i := 0; i < len(all); i++ {
		tmp := all[i]
		if tmp == nil {
			continue
		}
		tmpWight := tmp.CardType
		if tmp.CardType == cardConst.CARD_PATTERN_SEQUENCE ||
			tmp.CardType == cardConst.CARD_PATTERN_SEQUENCE_OF_PAIRS ||
			tmp.CardType == cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS {
			tmpWight = tmpWight - int32(len(tmp.Card))
		}

		if min.Wight < tmpWight {
			min = * tmp
		}
	}

	if min.Wight == 1000 {
		return nil
	}
	return &min
}

/*
	当有人报单的时候的单牌出牌
*/

func findBestSingleCard(re []*ReCard) *ReCard {
	// 根据长度取中间的单张
	sLen := len(re)
	if sLen == 0 {
		return nil
	}

	var index int

	if sLen%2 == 0 {
		index = (sLen - 1) / 2
	} else {
		index = sLen/2 - 1
	}

	var result ReCard
	result.CardType = cardConst.CARD_PATTERN_SINGLE
	result.Wight = re[index].Card[0].Value
	result.Card = append(result.Card, re[index].Card...)
	return &result
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

// 最小跟牌
func minFollowCard(actionPlayer *Player, eCards []*Card, eType int32) ([]*Card, bool) {

	eWight := countCardWight(eCards, eType)
	group := actionPlayer.GroupCard

	switch eType {
	case cardConst.CARD_PATTERN_SINGLE:
		g := group.Single
		if g != nil {
			for i := 0; i < len(g); i++ {
				if g[i].Wight > eWight {
					return g[i].Card, true
				}
			}
		}
		return nil, false
	case cardConst.CARD_PATTERN_PAIR:
		g := group.Double
		if g != nil {
			for i := 0; i < len(g); i++ {
				if g[i].Wight > eWight {
					return g[i].Card, true
				}
			}
		}
		return nil, false

	case cardConst.CARD_PATTERN_TRIPLET:
		//g := group.Triple
		//if g != nil {
		//	for i := 0; i < len(g); i++ {
		//		if g[i].Wight > eWight {
		//			return g[i].Card, true
		//		}
		//	}
		//}
		cards, b, _ := HostingBeatTriple(actionPlayer.HandCards, eCards)

		return cards, b

	case cardConst.CARD_PATTERN_TRIPLET_WITH_SINGLE:
		g := group.Triple
		if g != nil {
			for i := 0; i < len(g); i++ {
				if g[i].Wight > eWight {
					singles := group.Single
					if len(singles) <= 0 { //||len { //
						return nil, false
					}
					var r []*Card
					r = append(r, singles[0].Card...)
					r = append(r, g[i].Card...)
					return r, true
				}
			}
		}
		return nil, false

	case cardConst.CARD_PATTERN_TRIPLET_WITH_PAIR:
		g := group.Triple
		if g != nil {
			for i := 0; i < len(g); i++ {
				if g[i].Wight > eWight {
					return g[i].Card, true
				}
			}
		}
		return nil, false

	case cardConst.CARD_PATTERN_BOMB:
		// todo 这里先返回nil
		return nil, false

	case cardConst.CARD_PATTERN_QUADPLEX_WITH_SINGLES: // 四代二单
		// todo 这里先返回nil
		return nil, false

	case cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS: // 四代二对
		// todo 这里先返回nil
		return nil, false

	case cardConst.CARD_PATTERN_SEQUENCE:
		g := group.Junko
		if g != nil {
			for i := 0; i < len(g); i++ {
				if g[i].Wight > eWight && len(g[i].Card) == len(eCards) {
					return g[i].Card, true
				}
			}
		}
		return nil, false

	case cardConst.CARD_PATTERN_SEQUENCE_OF_PAIRS: // 连对
		g := group.JunkoDouble
		if g != nil {
			for i := 0; i < len(g); i++ {
				if g[i].Wight > eWight && len(g[i].Card) == len(eCards) {
					return g[i].Card, true
				}
			}
		}
		return nil, false

	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS: // 飞机不带
		g := group.JunkTriple
		if g != nil {
			for i := 0; i < len(g); i++ {
				if g[i].Wight > eWight && len(g[i].Card) == len(eCards) {
					return g[i].Card, true
				}
			}
		}
		return nil, false

	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_SINGLES: // 飞机带单
		g := group.Triple
		if g != nil {
			for i := 0; i < len(g); i++ {
				if g[i].Wight > eWight&& len(g[i].Card) == len(eCards) {
					singles := group.Single
					if len(singles) < len(eCards)/4 { //||len { //
						return nil, false
					}
					var r []*Card
					for j := 0; j < len(eCards)/4; j++ {
						r = append(r, singles[i].Card...)
					}
					r = append(r, g[i].Card...)
					return r, true
				}
			}
		}
		return nil, false

	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_PAIRS: // 飞机带对
		g := group.Triple
		if g != nil {
			for i := 0; i < len(g); i++ {
				if g[i].Wight > eWight && len(g[i].Card) == len(eCards){
					Double := group.Double
					if len(Double) < len(eCards)/4 { //||len { //
						return nil, false
					}
					var r []*Card
					for j := 0; j < len(eCards)/4; j++ {
						r = append(r, Double[i].Card...)
					}
					r = append(r, g[i].Card...)
					return r, true
				}
			}
		}
		return nil, false

	case cardConst.CARD_PATTERN_ROCKET: // 火箭
		return nil, false
	}
	return nil, false
}

// 获取飞机的权值
func getTripletWeight(cards []*Card) int32 {
	group := CreateGroupCard(cards)
	if len(group.Triple) <= 0 || group.Triple == nil {
		logger.Error("！！！！ 错误的飞机类型！！！")
		//	return 13
	}
	return group.Triple[0].Wight

}
