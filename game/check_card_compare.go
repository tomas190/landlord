package game

import (
	"landlord/mconst/cardConst"
	"sort"
)

// 是否能打过上家牌
/*
	cs:打出的牌
	eCs:上手要打出的牌
*/

func CanBeat(lastCards, outCards []*Card) bool {

	lcs := NewCardSet(lastCards)
	ocs := NewCardSet(outCards)

	CalPattern(ocs)
	CalPattern(lcs)

	return canBeatIt(ocs, lcs)
}

func canBeatIt(cs, eCs *CardSet) bool {
	if len(cs.Cards) == 0 {
		return false
	}

	// 如果上手牌是王炸
	if eCs.Pattern == cardConst.CARD_PATTERN_ROCKET {
		return false
	}

	// 这手牌是火箭
	if cs.Pattern == cardConst.CARD_PATTERN_ROCKET {
		return true
	}

	if cs.Pattern == cardConst.CARD_PATTERN_BOMB {
		if eCs.Pattern != cardConst.CARD_PATTERN_BOMB && eCs.Pattern != cardConst.CARD_PATTERN_ROCKET {
			return true
		}
	}

	return eCs.Pattern == cs.Pattern &&
		eCs.SeqCount == cs.SeqCount &&
		eCs.SubCount == cs.SubCount &&
		eCs.Rank < cs.Rank

}

/*
	从手牌中找出能打过eCards的一组牌

	CARD_PATTERN_SINGLE                                      // 3单张               from 3 (low) up to red joker (high)
	CARD_PATTERN_PAIR                                        // 4对子               3-3, A-A
	CARD_PATTERN_TRIPLET                                     // 5三不带             9-9-9.
	CARD_PATTERN_TRIPLET_WITH_SINGLE                         // 6三带一             9-9-9-3 beats 8-8-8-A.
	CARD_PATTERN_TRIPLET_WITH_PAIR                           // 7三带对             Q-Q-Q-6-6 beats 10-10-10-K-K.
	CARD_PATTERN_SEQUENCE                                    // 8顺子               from 3 up to ace - for example 8-9-10-J-Q. 2 and jokers cannot be used.
	CARD_PATTERN_SEQUENCE_OF_PAIRS                           // 9连对               10-10-J-J-Q-Q-K-K.
	CARD_PATTERN_SEQUENCE_OF_TRIPLETS                        // 10飞机不带翅膀        4-4-4-5-5-5.
	CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_SINGLES  // 11飞机带单翅膀        7-7-7-8-8-8-3-6.
	CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_PAIRS    // 12飞机带对翅膀        8-8-8-9-9-9-4-4-J-J.
	CARD_PATTERN_BOMB                                        // 13炸弹
	CARD_PATTERN_ROCKET                                      // 14火箭
	CARD_PATTERN_QUADPLEX_WITH_SINGLES                       // 15四带两单            6-6-6-6-8-9,
	CARD_PATTERN_QUADPLEX_WITH_PAIRS                         // 16四带两对            J-J-J-J-9-9-Q-Q.m

*/
func FindCanBeatCards(handsCard, eCards []*Card,eCardType int32) (bool, []*Card, int32) {
	// todo
	switch eCardType {
	case cardConst.CARD_PATTERN_SINGLE:


	}



	return false, nil, 0
}

/*
	桌面无牌，单张-对子-三张-炸弹出最小的牌型
	尾牌符合一次性出完原则，则自动出完，比如三带一、顺子、三代二、四带二、火箭等等

*/

/* ================================= 托管必出牌抽取 ==========================*/

func FindMustBeOutCards(handsCard []*Card) ([]*Card, int32) {

	if len(handsCard) <= 0 {
		return nil, cardConst.CARD_PATTERN_TODO
	}

	//1. 先判断能否一首出完 	// 检测牌型正确代表能一首出完
	cardType := GetCardsType(handsCard)
	if !(cardType > cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS || cardType < cardConst.CARD_PATTERN_SINGLE) {
		return handsCard, cardType
	}

	//2.老老实实找找单张 对子
	SortCard(handsCard)

	// 2.1 找单张
	if cards, b := findMinSingle(handsCard); b {
		return cards, cardConst.CARD_PATTERN_SINGLE
	}

	if cards, b := findMinDouble(handsCard); b {
		return cards, cardConst.CARD_PATTERN_PAIR
	}

	if cards, b := findMinTriple(handsCard); b {
		return cards, cardConst.CARD_PATTERN_TRIPLET
	}

	if cards, b := findMinBoom(handsCard); b {
		return cards, cardConst.CARD_PATTERN_BOMB
	}

	return handsCard[len(handsCard)-1:], cardConst.CARD_PATTERN_SINGLE
}

/*
  返回手牌中的最小的单张  // 不用管是否组成顺子
	handCards 需要时排好序的切片
*/
func findMinSingle(handCards []*Card) ([]*Card, bool) {
	if len(handCards) <= 0 {
		return nil, false
	}


	// 如果最后一张牌是大王  则改成和小王同级
	if handCards[0].Value == cardConst.CARD_RANK_RED_JOKER {
		handCards[0].Value = cardConst.CARD_RANK_BLACK_JOKER
		defer func() {
			handCards[0].Value = cardConst.CARD_RANK_RED_JOKER
		}()
	}
	// 先统计牌的张输
	var doubleArr []int
	cardCount := make(map[int32]int, len(handCards))
	for i := 0; i < len(handCards); i++ {
		cardCount[handCards[i].Value] = cardCount[handCards[i].Value] + 1
	}

	for k, v := range cardCount {
		if v == 1 {
			doubleArr = append(doubleArr, int(k))
		}
	}

	if len(doubleArr) == 0 {
		return nil, false
	}

	sort.Ints(doubleArr)
	var result []*Card

	for i := 0; i < len(handCards); i++ {
		if int(handCards[i].Value) == doubleArr[0] {
			result = append(result, handCards[i])
		}
		if len(result) == 1 {
			break
		}
	}
	return result, true
}

/*
  返回手牌中的最小的对子
  handCards 需要时排好序的切片
*/
func findMinDouble(handCards []*Card) ([]*Card, bool) {

	// 先统计牌的张输
	var doubleArr []int
	cardCount := make(map[int32]int, len(handCards))
	for i := 0; i < len(handCards); i++ {
		cardCount[handCards[i].Value] = cardCount[handCards[i].Value] + 1
	}

	for k, v := range cardCount {
		if v == 2 {
			doubleArr = append(doubleArr, int(k))
		}
	}

	if len(doubleArr) == 0 {
		return nil, false
	}

	sort.Ints(doubleArr)
	var result []*Card

	for i := 0; i < len(handCards); i++ {
		if int(handCards[i].Value) == doubleArr[0] {
			result = append(result, handCards[i])
		}
		if len(result) == 2 {
			break
		}
	}
	return result, true
}

/*
  返回手牌中的最小的三张一样的牌
*/
func findMinTriple(handCards []*Card) ([]*Card, bool) {

	// 先统计牌的张输
	var triple []int
	cardCount := make(map[int32]int, len(handCards))
	for i := 0; i < len(handCards); i++ {
		cardCount[handCards[i].Value] = cardCount[handCards[i].Value] + 1
	}

	for k, v := range cardCount {
		if v == 3 {
			triple = append(triple, int(k))
		}
	}

	if len(triple) == 0 {
		return nil, false
	}

	sort.Ints(triple)
	var result []*Card

	for i := 0; i < len(handCards); i++ {
		if int(handCards[i].Value) == triple[0] {
			result = append(result, handCards[i])
		}
		if len(result) == 3 {
			break
		}
	}
	return result, true
}

/*
  返回手牌中的最小的炸弹
*/
func findMinBoom(handCards []*Card) ([]*Card, bool) {
	// 先统计牌的张输
	var triple []int
	cardCount := make(map[int32]int, len(handCards))
	for i := 0; i < len(handCards); i++ {
		cardCount[handCards[i].Value] = cardCount[handCards[i].Value] + 1
	}

	for k, v := range cardCount {
		if v == 4 {
			triple = append(triple, int(k))
		}
	}

	if len(triple) == 0 {
		return nil, false
	}

	sort.Ints(triple)
	var result []*Card

	for i := 0; i < len(handCards); i++ {
		if int(handCards[i].Value) == triple[0] {
			result = append(result, handCards[i])
		}
		if len(result) == 4 {
			break
		}
	}
	return result, true
}

/* ================================= 托管必出牌抽取 ============================*/