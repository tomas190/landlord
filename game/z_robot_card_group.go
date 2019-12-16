package game

import (
	"landlord/mconst/cardConst"
	"strconv"
	"strings"
)

type GroupCard struct {
	OriginCards      []*Card // 原始手牌
	Score            int     // 原始牌得分
	WeightScore      int     // 得分
	OriginCardsNums  int     // 原始牌多少次能出完
	CurrentHandCards []*Card // 当前手牌

	Single      []*ReCard // 单张
	Double      []*ReCard // 对子
	Triple      []*ReCard // 三张 不安排三代一 或者三带二
	Bomb        []*ReCard // 炸弹
	Junko       []*ReCard // 顺子
	JunkoDouble []*ReCard // 连对
	Rocket      []*ReCard // 火箭
}

//
type ReCard struct {
	Wight    int32   // 该牌组权重
	CardType int32   // 该牌类型
	Card     []*Card // 卡牌组
}

/*
	机器人不会主动出四带二
	todo 所有的组牌策略 均用不拆炸弹原则
*/

// 1.顺子优先组牌策略
// todo 分成保三张 和不保 三张拆法
func PriorityJunko(handCards []*Card) GroupCard {
	// 先无限取最小5连 后续合并
	junkos, remainCards := unlimitedJunko(handCards)

	// 1.在从剩余牌中找出所有炸弹
	bombs, remainCards := FindAllBomb(remainCards)

	// 2.单牌看能不能组成更大的顺子
	junkoWithSingles, remainCards := mergeJunkoWithSingle(junkos, remainCards)

	// 3.顺子于顺子合并
	mergeJunkos := mergeJunkoWithJunko(junkoWithSingles)

	// 分析剩余牌型
	group := CreateGroupCard(remainCards)
	group.Bomb = bombs
	group.Junko = mergeJunkos
	return group

}

// 1.顺子优先组牌策略
// todo 分成保三张 和不保 三张拆法
func PriorityJunkoProtectTriple(handCards []*Card) GroupCard {

	triples, remainCards := FindAllTriplet(handCards)

	// 先无限取最小5连 后续合并
	junkos, remainCards := unlimitedJunko(remainCards)

	// 1.在从剩余牌中找出所有炸弹
	bombs, remainCards := FindAllBomb(remainCards)

	// 2.单牌看能不能组成更大的顺子
	junkoWithSingles, remainCards := mergeJunkoWithSingle(junkos, remainCards)

	// 3.顺子于顺子合并
	mergeJunkos := mergeJunkoWithJunko(junkoWithSingles)

	// 分析剩余牌型
	group := CreateGroupCard(remainCards)
	group.Bomb = bombs
	group.Junko = mergeJunkos
	group.Triple = triples
	return group

}

// 1.连对 优先组牌策略
func PriorityJunkoDouble(handCards []*Card) GroupCard {
	//	doubles, remainCards := FindAllDouble(handCards)

	group := CreateGroupCard(handCards)
	return group

}

// 1.三张 优先组牌策略
func PriorityTriple(handCards []*Card) GroupCard {

	group := CreateGroupCard(handCards)
	return group

}

// 4.尽量少单的拆法
func possiblyLessSingle(handCards []*Card) GroupCard {

	group := CreateGroupCard(handCards)
	return group

}

// 根据手牌分析
// 将手牌分成 单张 对子 三张(不带任何) 炸弹 火箭
// 这个方法不组成顺子 连对 飞机等牌型
func CreateGroupCard(handCards []*Card) GroupCard {
	var gc GroupCard
	gc.CurrentHandCards = handCards
	gc.OriginCards = handCards

	rocket, remainCards := FindRocket(handCards)
	//	PrintCard(remainCards)
	bombCards, remainCards := FindAllBomb(remainCards)
	//PrintCard(remainCards)
	triplets, remainCards := FindAllTriplet(remainCards)
	//	PrintCard(remainCards)
	doubles, remainCards := FindAllDouble(remainCards)
	//PrintCard(remainCards)
	singles, remainCards := FindAllSingle(remainCards)
	//PrintCard(remainCards)
	gc.OriginCards = handCards
	gc.CurrentHandCards = handCards
	gc.Rocket = rocket
	gc.Bomb = bombCards
	gc.Triple = triplets
	gc.Double = doubles
	gc.Single = singles

	return gc
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

	singles := getHasNumsCard(handCards, 1)

	// 移除同时存在大小鬼的单张
	// 将同时又大小王的情况下 将其移除
	if len(singles) >= 2 {
		if singles[len(singles)-1] == cardConst.CARD_RANK_RED_JOKER &&
			singles[len(singles)-2] == cardConst.CARD_RANK_BLACK_JOKER {
			singles = append([]int{}, singles[:len(singles)-2]...)
		}
	}

	var reCards []*ReCard
	if len(singles) > 0 {
		remainCard := append([]*Card{}, handCards...)
		for i := 0; i < len(singles); i++ {
			card := findThisValueCard(singles[i], handCards, 1)
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

// 找到所有顺子组合
/*
	return
		[]*ReCard 需要找的牌
		[]*Card   剩余手牌
*/
func FindAllJunkos(hands []*Card) ([]*ReCard, []*Card) {
	remainCards := hands
	var junkos []*ReCard
	for {
		var junko *ReCard
		junko, remainCards = FindPossibleLongSingleJunko(remainCards)
		if junko == nil {
			break
		}
		junkos = append(junkos, junko)
	}

	return junkos, remainCards
}

// 组牌策略 顺子开始

// 寻找尽可能长的顺子 炸弹不会在此方法里面计算
/*
	return
		[]*ReCard 需要找的牌
		[]*Card   剩余手牌
*/
func FindPossibleLongSingleJunko(handCards []*Card) (*ReCard, []*Card) {

	// 1.先去掉手牌中的重复值 要去掉 2 以上大的牌 (2 以上大的牌不能组成顺子) 炸弹的不包含在内
	singleCards := junkoHelpRemove(handCards)
	if len(singleCards) < 5 {
		return nil, handCards
	}
	//	logger.Debug("去重:", singleCards)

	var psj []int
	hLen := len(singleCards)
	for i := 0; i < hLen; i++ {
		for j := hLen; j >= i+5; j-- {
			// 最先组成的一定是最长的
			if isJunko(singleCards[i:j]) {
				// 已经找到
				//	logger.Debug("find:", singleCards[i:j])
				psj = singleCards[i:j]
				goto dre
			}
		}
	}

dre:
	if len(psj) < 5 {
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

	return &rc, result

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

// ============================

// 无限取最小5连
// 并从权值从从小到达追加
func unlimitedJunko(handCards []*Card) ([]*ReCard, []*Card) {
	//min := []*Card{{1, 1}, {2, 1}, {3, 1}, {4, 1}, {5, 1}}
	min := []*Card{{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}}

	var result []*ReCard
	tmpHands := append([]*Card{}, handCards...)
	for {
		//logger.Debug("tmphands:")
		//PrintCard(tmpHands)
		if len(tmpHands) < 5 {
			break
		}
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
	保准junkos是顺子
*/
func mergeJunkoWithJunko(junkos []*ReCard) []*ReCard {

	if len(junkos) <= 1 {
		return junkos
	}
	copyJunkos := append([]*ReCard{}, junkos...)

	var mergeIndexs string
	var mergeJunkos []*ReCard
	for i := 0; i < len(copyJunkos)-1; i++ {
		for j := i + 1; j < len(copyJunkos); j++ {
			if strings.Contains(mergeIndexs, strconv.Itoa(i)) ||
				strings.Contains(mergeIndexs, strconv.Itoa(j)) {
				continue
			}
			card, b := canMergeJunkos(copyJunkos[i], copyJunkos[j])
			if b {
				//添加index标记
				mergeIndexs = mergeIndexs + strconv.Itoa(i) + strconv.Itoa(j)
				mergeJunkos = append(mergeJunkos, card)
				break
			}
		}
	}

	//	logger.Debug("合并的index", mergeIndexs)
	for i := 0; i < len(junkos); i++ {
		if !strings.Contains(mergeIndexs, strconv.Itoa(i)) {
			mergeJunkos = append(mergeJunkos, junkos[i])
		}
	}
	return mergeJunkos

}

// 两个顺子是否可以合并
// 如果可以合并则组合好返回
func canMergeJunkos(card1, card2 *ReCard) (*ReCard, bool) {
	if len(card1.Card) < 5 ||
		len(card2.Card) < 5 {
		return nil, false
	}

	if card1.Card[0].Value-card2.Card[len(card2.Card)-1].Value == 1 ||
		card2.Card[0].Value-card1.Card[len(card1.Card)-1].Value == 1 {
		card1.Card = append(card1.Card, card2.Card...)
		SortCardSL(card1.Card)
		card1.Wight = card1.Card[len(card1.Card)-1].Value
		return card1, true
	}
	return nil, false
}

// 顺子于单张合并
/*
	保准junkos 顺子
	return :更大的顺子
			剩余的牌
*/
func mergeJunkoWithSingle(junkos []*ReCard, singles []*Card) ([]*ReCard, []*Card) {
	var mergeIndexs string
	//logger.Debug("last I:==================len:", len(singles))
	for i := 0; i < len(junkos); i++ {
		for j := 0; j < len(singles); j++ {
			if strings.Contains(mergeIndexs, strconv.Itoa(j)) || singles[j].Value >= cardConst.CARD_RANK_TWO {
				continue
			}
			//logger.Debug("i===", i)
			card, b := canMergeJunkoWithSingle(junkos[i], singles[j])
			if b {
				mergeIndexs = mergeIndexs + "," + strconv.Itoa(j) + ","
				junkos[i] = card
			}
		}
	}

	var remain []*Card
	for i := 0; i < len(singles); i++ {
		if !strings.Contains(mergeIndexs, strconv.Itoa(i)) {
			remain = append(remain, singles[i])
		}
	}

	return junkos, remain
}

// 顺子和单排是否可以合并
// 如果可以合并则组合好返回
func canMergeJunkoWithSingle(card1 *ReCard, single *Card) (*ReCard, bool) {

	if single.Value >= cardConst.CARD_RANK_TWO ||
		single.Value < cardConst.CARD_RANK_THREE ||
		len(card1.Card) < 5 {
		return nil, false
	}

	if card1.Card[0].Value-1 == single.Value ||
		card1.Card[len(card1.Card)-1].Value+1 == single.Value {
		card1.Card = append(card1.Card, single)
		SortCardSL(card1.Card)
		return card1, true
	}

	return nil, false
}

// 将所有牌装换成类型的单牌
func changeCardToReCard(cards []*Card) []*ReCard {
	var result []*ReCard
	for i := 0; i < len(cards); i++ {
		var r ReCard
		r.Card = append(r.Card, cards[i])
		result = append(result, &r)
	}

	return result
}
