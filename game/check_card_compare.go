package game

import (
	"landlord/mconst/cardConst"
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
*/
func FindCanBeatCards(handsCard, eCards []*Card) (bool, []*Card, int32) {
	// todo

	return false, nil, 0
}

/*
	必出牌 获取必出牌
*/
func FindMustBeOutCards(handsCard []*Card) ([]*Card, int32) {
	// todo 判断最后能否出完
	SortCard(handsCard)
	return handsCard[len(handsCard)-1:], cardConst.CARD_PATTERN_SINGLE

}

/*

 */
