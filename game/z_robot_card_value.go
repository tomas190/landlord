package game

import (
	"landlord/mconst/cardConst"
	"strconv"
	"strings"
)

// todo  计算完成之后根据 还要根据玩家还有多少首数牌 再次计分
// 根据组牌牌型计算分值
func countCardGroupValue(groupCard GroupCard) float64 {

	var groupValue float64
	single := groupCard.Single
	for i := 0; i < len(single); i++ {
		groupValue += countCardTypeValue(single[i], cardConst.CARD_PATTERN_SINGLE)
	}

	double := groupCard.Double
	for i := 0; i < len(double); i++ {
		groupValue += countCardTypeValue(double[i], cardConst.CARD_PATTERN_PAIR)
	}

	triple := groupCard.Triple
	for i := 0; i < len(triple); i++ {
		groupValue += countCardTypeValue(triple[i], cardConst.CARD_PATTERN_TRIPLET)
	}

	bomb := groupCard.Bomb
	for i := 0; i < len(bomb); i++ {
		groupValue += countCardTypeValue(bomb[i], cardConst.CARD_PATTERN_BOMB)
	}

	rocket := groupCard.Rocket
	for i := 0; i < len(rocket); i++ {
		groupValue += countCardTypeValue(rocket[i], cardConst.CARD_PATTERN_ROCKET)
	}

	Junko := groupCard.Junko
	for i := 0; i < len(Junko); i++ {
		groupValue += countCardTypeValue(Junko[i], cardConst.CARD_PATTERN_SEQUENCE)
	}

	return groupValue

}

// 计算牌组价值是否叫分和抢地主
/*
	2 		1.75
	小王 	2.5
	大王 	3

	炸弹 	5+base*0.1
	王炸    8
*/
func CountCardValue(cards []*Card) float32 {
	var value float32

	for i := 0; i < len(cards); i++ {
		if cards[i].Value == cardConst.CARD_RANK_TWO {
			value += 1.75
		} else if cards[i].Value == cardConst.CARD_RANK_BLACK_JOKER {
			value += 2.5
		} else if cards[i].Value == cardConst.CARD_RANK_RED_JOKER {
			value += 3
		}
	}

	reCards, _ := FindAllBomb(cards)
	for i := 0; i < len(reCards); i++ {
		if reCards[i].Wight == cardConst.CARD_RANK_TWO {
			value += 7
			value = value - 1.75*4
		} else {
			value += 4 + float32(reCards[i].Wight)*0.1
		}
	}

	_, b, _ := hasRacket(cards)
	if b {
		value += 8
		value = value - 5.5
	}
	return value
}

func CountCardValue2(cards []*Card) float32 {
	var value float32

	for i := 0; i < len(cards); i++ {
		if cards[i].Value == cardConst.CARD_RANK_TWO {
			value += 1.75
		} else if cards[i].Value == cardConst.CARD_RANK_BLACK_JOKER {
			value += 2.5
		} else if cards[i].Value == cardConst.CARD_RANK_RED_JOKER {
			value += 3
		}
	}

	reCards, _ := FindAllBomb(cards)
	for i := 0; i < len(reCards); i++ {
		if reCards[i].Wight == cardConst.CARD_RANK_TWO {
			value += 5
		} else {
			value += 4
		}
	}

	_, b, _ := hasRacket(cards)
	if b {
		value += 6
	}
	return value
}

// 根据组牌牌型计算分值 todo 很大的优化空间
func countCardTypeValue(card *ReCard, cardType int) float64 {

	/*
			这里暂时根据  2 大小王 炸弹的数量对手牌进行评分

		    后续跟拍原则 根据出牌后是否 能减少首数

	*/

	var value float64
	baseValue := float64(card.Wight - 10)
	switch cardType {
	case cardConst.CARD_PATTERN_SINGLE: // 单张价值 = 基础价值
		value = baseValue
	case cardConst.CARD_PATTERN_PAIR: // 对子价值 = 基础价值+0.05
		value = baseValue + 0.1
	case cardConst.CARD_PATTERN_TRIPLET: // 三张价值 = 基础价值+0.1
		value = baseValue + 0.2
	case cardConst.CARD_PATTERN_SEQUENCE: // 顺子  = 基础价值+0.2+ (len)*0.03
		value = baseValue + 0.2 + float64(len(card.Card)/2-4)*0.01
	case cardConst.CARD_PATTERN_SEQUENCE_OF_PAIRS: // 连对 = 基础价值+0.2+ (len)*0.05
		value = baseValue + 0.2 + float64(len(card.Card)/2-2)*0.05
	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS: // 飞机  = 基础价值+0.3+ (len)*0.1
		value = baseValue + 0.3 + float64(len(card.Card)/2-1)*0.5
	case cardConst.CARD_PATTERN_BOMB: // 炸弹  = 基础价值+7
		value = baseValue + 7
	case cardConst.CARD_PATTERN_ROCKET: // 火箭 = 20
		return 20
	default:
		return 0
	}

	return value
}

// 取绝对值
func absSelf(num int32) int32 {
	if num > 0 {
		return num
	}
	return -num
}

/*=====================================================================================*/

/*
以下是不合理的拆法仅做参考

拆牌

如果有炸弹，找出来（炸弹不拆）；
如果有飞机，找出来（飞机也不拆）；(345677788899910JQK) 这种情况拆不拆？
找出所有的顺子，每个顺子尽量长；
处理顺子
顺子分裂，现有一个顺子345678910，发现手上还剩下一张6和7，那么顺子分裂成34567和678910；
顺子拆出连对，现有一个顺子345678910，发现手上还有345，变成334455和678910更好，就把顺子里面的345放回去；
顺子拆出三张，现有一个顺子345678，发现手上还有88，那么变成34567和888更好，就把顺子里面的8放回去；
顺子如果盖住对子、三张、连对，如果发现打散牌组数更少，则打散，比如7789910JJJQKK，拆散了更好；
顺子拆出两头的对子，现有一个顺子67890JQ，发现手上还有Q和6，那么把孙子里面的Q和6放回去；
反复进行1）到5），直到没有进一步变化；
剩余的牌里面找出所有的连对；
查看所有的连对，如果长度超过3，且一端有三条，把三条拆出来；
剩余的牌里面找出所有的三张；
延长顺子：如果一个顺子的两端外面有一个对子，如果这个对子小于10，则并入顺子，比如34567+88,那么变成345678+8;
合并顺子：相同的顺子变成连对，首尾相接的顺子连成一个；
剩余的牌里面找出所有对子和单张。

*/

// 相同权值和长度的顺子合并成连对
/*
	保准junkos是顺子
	return : 连对
             剩余的顺子
*/
func mergeSameJunkoAsSeqDouble(junkos []*ReCard) ([]*ReCard, []*ReCard) {

	if len(junkos) <= 1 {
		return nil, junkos
	}
	copyJunkos := append([]*ReCard{}, junkos...)

	var mergeIndexs string
	var mergeDoubles []*ReCard
	var remainSingles []*ReCard
	for i := 0; i < len(copyJunkos)-1; i++ {
		for j := i + 1; j < len(copyJunkos); j++ {
			if strings.Contains(mergeIndexs, strconv.Itoa(i)) ||
				strings.Contains(mergeIndexs, strconv.Itoa(j)) {
				continue
			}
			card, b := isSameJunko(copyJunkos[i], copyJunkos[j])
			if b {
				//添加index标记
				mergeIndexs = mergeIndexs + strconv.Itoa(i) + strconv.Itoa(j)
				mergeDoubles = append(mergeDoubles, card)
				break
			}
		}
	}

	//	logger.Debug("合并的index", mergeIndexs)
	for i := 0; i < len(junkos); i++ {
		if !strings.Contains(mergeIndexs, strconv.Itoa(i)) {
			remainSingles = append(remainSingles, junkos[i])
		}
	}
	return mergeDoubles, remainSingles

}

func isSameJunko(card *ReCard, card2 *ReCard) (*ReCard, bool) {
	if card == nil || card2 == nil {
		return nil, false
	}
	if len(card.Card) <= 0 || len(card2.Card) <= 0 {
		return nil, false
	}

	if card.Wight == card2.Wight && len(card.Card) == len(card2.Card) {
		card.Card = append(card.Card, card2.Card...)
		SortCardSL(card.Card)
		return card, true
	}
	return nil, false
}

/*
	param
		hands []*Card : 手牌

	return
		[]*ReCard	:相同顺子组成的连对
					:顺子
		[]*Card		:剩余手牌
*/
func FindAllJunko(hands []*Card) ([]*ReCard, []*ReCard, []*Card) {

	minJunko, remainCards := unlimitedJunko(hands)
	// 1.在从剩余牌中找出所有炸弹
	_, remainCards = FindAllBomb(remainCards)
	mJunkoWs, remainCards := mergeJunkoWithSingle(minJunko, remainCards)
	junkos := mergeJunkoWithJunko(mJunkoWs)
	seqDouble, remainJunkos := mergeSameJunkoAsSeqDouble(junkos)
	return seqDouble, remainJunkos, remainCards

}
