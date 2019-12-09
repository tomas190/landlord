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
func FindCanBeatCards(handsCard, eCards []*Card, eCardType int32) ([]*Card, bool, int32) {
	switch eCardType {
	case cardConst.CARD_PATTERN_SINGLE: // 3
		return HostingBeatSingle(handsCard, eCards)

	case cardConst.CARD_PATTERN_PAIR: // 3
		return HostingBeatDouble(handsCard, eCards)

	case cardConst.CARD_PATTERN_TRIPLET: // 3
		return HostingBeatTriple(handsCard, eCards)

	case cardConst.CARD_PATTERN_TRIPLET_WITH_SINGLE: // 3
		return HostingBeatTripleWithSingle(handsCard, eCards)

	case cardConst.CARD_PATTERN_TRIPLET_WITH_PAIR: // 3
		return HostingBeatTripleWithDouble(handsCard, eCards)

	case cardConst.CARD_PATTERN_SEQUENCE: // 3
		return HostingBeatJunko(handsCard, eCards)

	case cardConst.CARD_PATTERN_SEQUENCE_OF_PAIRS: // 3
		return HostingBeatContinuouslyDouble(handsCard, eCards)

	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS: // 飞机
		return HostingBeatTriplets(handsCard, eCards)

	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_SINGLES: // 3
		return HostingBeatTripletsWithSingle(handsCard, eCards)

	case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_PAIRS: // 3
		return HostingBeatTripletsWithDouble(handsCard, eCards)

	case cardConst.CARD_PATTERN_BOMB: // 3
		return HostingBeatBomb(handsCard, eCards)

	case cardConst.CARD_PATTERN_QUADPLEX_WITH_SINGLES: // 3
		return HostingBeatBombWithSingles(handsCard, eCards)

	case cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS: // 3
		return HostingBeatBombWithDouble(handsCard, eCards)

	case cardConst.CARD_PATTERN_ROCKET: //  // 火箭直接返回不能大过
		return nil, false, cardConst.CARD_PATTERN_TODO

	}

	return nil, false, cardConst.CARD_PATTERN_TODO
}
