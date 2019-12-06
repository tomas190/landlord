package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/cardConst"
)

type GroupCard struct {
	OriginCards      []*Card // 原始手牌
	Score            int     // 原始牌得分
	WeightScore      int     // 大牌得分  - 火箭为8分，炸弹为6分，大王4分，小王3分，一个2为2分
	OriginCardsNums  int     // 原始牌多少次能出完
	CurrentHandCards []*Card // 当前手牌

	Single []*ReCard // 单张
	Double []*ReCard // 对子
	Triple []*ReCard // 三张 不安排三代一 或者三带二
	Bomb   []*ReCard // 炸弹
	Junko  []*ReCard // 顺子
	Rocket []*ReCard // 火箭
}

//
type ReCard struct {
	Wight int32   // 该牌组权重
	Card  []*Card // 卡牌组
}

func GetOriginCardsScore(handCards []*Card) int {
	var score int
	// 牌分估值

	return score

}

// 根据手牌分析
func CreateGroupCard(handCards []*Card) *GroupCard {
	var gc GroupCard
	gc.CurrentHandCards = handCards
	gc.OriginCards = handCards

	rocket, remainCards := FindRocket(handCards)
	PrintCard(remainCards)
	bombCards, remainCards := FindAllBomb(remainCards)
	PrintCard(remainCards)
	triplets, remainCards := FindAllTriplet(remainCards)
	PrintCard(remainCards)
	doubles, remainCards := FindAllDouble(remainCards)
	PrintCard(remainCards)
	singles, remainCards := FindAllSingle(remainCards)
	PrintCard(remainCards)
	gc.OriginCards = handCards
	gc.CurrentHandCards = handCards
	gc.Rocket = rocket
	gc.Bomb = bombCards
	gc.Triple = triplets
	gc.Double = doubles
	gc.Single = singles

	return &gc
}

// 找出火箭
/*
	return
		[]*ReCard 需要找的牌
		[]*Card   剩余手牌
*/
func FindRocket(handCards []*Card) ([]*ReCard, []*Card) {
	rocket, b, _ := hasRacket(handCards)
	var reCards []*ReCard
	if b {
		remainCards := removeCards(handCards, rocket)

		var reCard ReCard
		reCard.Wight = cardConst.CARD_RANK_RED_JOKER
		reCard.Card = rocket
		reCards = append(reCards, &reCard)
		return reCards, remainCards
	}
	return nil, handCards
}

// 找到所有炸弹
/*
	return
		[]*ReCard 需要找的牌
		[]*Card   剩余手牌
*/
func FindAllBomb(handCards []*Card) ([]*ReCard, []*Card) {
	bomb := getHasNumsCard(handCards, 4)
	var reCards []*ReCard
	if len(bomb) > 0 {
		remainCard := append([]*Card{}, handCards...)
		for i := 0; i < len(bomb); i++ {
			card := findThisValueCard(bomb[i], handCards, 4)
			var re ReCard
			re.Wight = card[0].Value
			re.Card = card
			reCards = append(reCards, &re)

			// 移除牌
			remainCard = removeCards(remainCard, card)
		}
		return reCards, remainCard
	}

	return nil, handCards
}

// 找到所有三张
/*
	return
		[]*ReCard 需要找的牌
		[]*Card   剩余手牌
*/
func FindAllTriplet(handCards []*Card) ([]*ReCard, []*Card) {
	bomb := getHasNumsCard(handCards, 3)
	var reCards []*ReCard
	if len(bomb) > 0 {
		remainCard := append([]*Card{}, handCards...)
		for i := 0; i < len(bomb); i++ {
			card := findThisValueCard(bomb[i], handCards, 3)
			var re ReCard
			re.Wight = card[0].Value
			re.Card = card
			reCards = append(reCards, &re)

			// 移除牌
			remainCard = removeCards(remainCard, card)
		}
		return reCards, remainCard
	}

	return nil, handCards
}

// 找到所有对子
/*
	return
		[]*ReCard 需要找的牌
		[]*Card   剩余手牌
*/
func FindAllDouble(handCards []*Card) ([]*ReCard, []*Card) {
	bomb := getHasNumsCard(handCards, 2)
	var reCards []*ReCard
	if len(bomb) > 0 {
		remainCard := append([]*Card{}, handCards...)
		for i := 0; i < len(bomb); i++ {
			card := findThisValueCard(bomb[i], handCards, 2)
			var re ReCard
			re.Wight = card[0].Value
			re.Card = card
			reCards = append(reCards, &re)

			// 移除牌
			remainCard = removeCards(remainCard, card)
		}
		return reCards, remainCard
	}
	return nil, handCards

}

// 找到所有单张
/*
	return
		[]*ReCard 需要找的牌
		[]*Card   剩余手牌
*/
func FindAllSingle(handCards []*Card) ([]*ReCard, []*Card) {

	if len(handCards) <= 0 {
		return nil, handCards
	}
	SortCard(handCards)
	// 如果最后一张牌是大王  则改成和小王同级
	if _, has, _ := hasRacket(handCards); has {
		handCards[0].Value = cardConst.CARD_RANK_BLACK_JOKER
		defer func() {
			handCards[0].Value = cardConst.CARD_RANK_RED_JOKER
		}()
	}

	bomb := getHasNumsCard(handCards, 1)
	var reCards []*ReCard
	if len(bomb) > 0 {
		remainCard := append([]*Card{}, handCards...)
		for i := 0; i < len(bomb); i++ {
			card := findThisValueCard(bomb[i], handCards, 1)
			var re ReCard
			re.Wight = card[0].Value
			re.Card = card
			reCards = append(reCards, &re)

			// 移除牌
			remainCard = removeCards(remainCard, card)
		}
		return reCards, remainCard
	}
	return nil, handCards

}

// 组牌策略 顺子开始

// 寻找尽可能长的顺子
/*
	return
		[]*ReCard 需要找的牌
		[]*Card   剩余手牌
*/
func FindPossibleLongSingleJunko(handCards []*Card) ([]*ReCard, []*Card) {

	// 1.先去掉手牌中的重复值 要去掉 2 以上大的牌 (2 以上大的牌不能组成顺子) 炸弹的不包含在内
	singleCards := junkoHelpRemove(handCards)
	if len(singleCards) < 5 {
		return nil, handCards
	}
	logger.Debug("去重:", singleCards)

	var psj []int
	hLen := len(singleCards)
	for i := 0; i < hLen; i++ {
		for j := hLen; j >= i+5; j-- {
			// 最先组成的一定是最长的
			if isJunko(singleCards[i:j]) {
				// 已经找到
				logger.Debug("find:", singleCards[i:j])
				psj = singleCards[i:j]
				goto dre
			}
		}
	}

dre:
	if len(psj) <= 5 {
		return nil, handCards
	}

	var junko []*Card
	for i := 0; i < len(psj); i++ {
		card := findThisValueCard(psj[i], handCards, 1)
		junko = append(junko, card...)
	}

	var rc ReCard
	rc.Card = junko
	rc.Wight = junko[len(junko)-1].Value

	result := removeCards(handCards, junko)
	return append([]*ReCard{}, &rc), result

}

// 简单计算是否顺子
/*
	arr保证是去重并且是排好序的数字*/
func isJunko(arr []int) bool {
	if len(arr) < 5 {
		return false
	}
	if arr[0]+len(arr) == arr[len(arr)-1]+1 {
		return true
	}
	return false
}

// 顺子优先组牌策略
//func priorityJunko(handCards []*Card) []*ReCard {
//	junkos, remainCards := unlimitedJunko(handCards)
//
//}

// ============================

// 无限取最小5连
func unlimitedJunko(handCards []*Card) ([]*ReCard, []*Card) {
	//min := []*Card{{1, 1}, {2, 1}, {3, 1}, {4, 1}, {5, 1}}
	min := []*Card{{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}}

	var result []*ReCard
	tmpHands := append([]*Card{}, handCards...)
	for {
		//logger.Debug("tmphands:")
		//PrintCard(tmpHands)
		var rc ReCard
		cards, b, i := HostingBeatJunko(tmpHands, min)
		if !b || i != cardConst.CARD_PATTERN_SEQUENCE || len(cards) < 5 {
			break
		}
		rc.Card = cards
		rc.Wight = cards[len(cards)-1].Value
		result = append(result, &rc)
		tmpHands = removeCards(tmpHands, cards)

	}
	return result, tmpHands
}

// 顺子于顺子合并
/*
	保准junkos 是顺子 并且从小到大依次排序
*/
func mergeJunko(junkos []*ReCard) []*ReCard {

	if len(junkos) <= 1 {
		return junkos
	}
	var newJunkos []*ReCard
	for i := 0; i < len(junkos); i++ {
		var newJunko ReCard
		if i == len(junkos)-1 {
			break
		}
		cJunko := junkos[i].Card
		nJunko := junkos[i+1].Card
		if cJunko[len(junkos[i].Card)-1].Value+1 == nJunko[0].Value {
			newJunko.Card = append(newJunko.Card, cJunko...)
			newJunko.Card = append(newJunko.Card, nJunko...)
			newJunko.Wight = nJunko[len(nJunko)-1].Value
			newJunkos = append(newJunkos, &newJunko)
		}
	}

	return newJunkos
}
