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

	return CanBeatIt(ocs, lcs)
}

func CanBeatIt(cs, eCs *CardSet) bool {
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

// tip 牌型计算这里挖的坑比较大
//func CalCardSet(cs *CardSet) int {
//
//	if len(cs.Cards) < 1 || len(cs.Cards) > 20 {
//		return cardConst.CARD_PATTERN_ERROR
//	}
//
//	sort.Slice(cs.Cards, func(i, j int) bool {
//		return cs.Cards[i].Value < cs.Cards[j].Value
//	})
//
//	if IsSingle(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_SINGLE
//		cs.Rank = cs.Cards[0].Value
//		cs.SeqCount = 1
//		cs.SubCount = 0
//	} else if IsPair(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_PAIR
//		cs.Rank = cs.Cards[0].Value
//		cs.SeqCount = 1
//		cs.SubCount = 0
//	} else if IsTriplet(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_TRIPLET
//		cs.Rank = cs.Cards[0].Value
//		cs.SeqCount = 1
//		cs.SubCount = 0
//	} else if IsTripletWithSingle(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_TRIPLET_WITH_SINGLE
//		cs.Rank = cs.Cards[1].Value
//		cs.SeqCount = 1
//		cs.SubCount = 1
//	} else if IsTripletWithPair(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_TRIPLET_WITH_PAIR
//		cs.Rank = cs.Cards[2].Value
//		cs.SeqCount = 1
//		cs.SubCount = 1
//	} else if IsSequence(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_SEQUENCE
//		cs.Rank = cs.Cards[0].Value
//		cs.SeqCount = len(cs.Cards)
//		cs.SubCount = 0
//	} else if IsSeqOfPair(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_SEQUENCE_OF_PAIRS
//		cs.Rank = cs.Cards[0].Value
//		cs.SeqCount = len(cs.Cards) / 2
//		cs.SubCount = 0
//	} else if IsSeqOfTriplet(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS
//		cs.Rank = cs.Cards[0].Value
//		cs.SeqCount = len(cs.Cards) / 3
//		cs.SubCount = 0
//		// tip 这里比较复杂了
//	} else if IsSeqOfTriWithSingles(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_SINGLES
//		cs.Rank, cs.SeqCount, cs.SubCount = RCCSeqOfTriWithSingles(cs.Cards)
//	} else if IsSeqOfTriWithPairs(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_PAIRS
//		cs.Rank, cs.SeqCount, cs.SubCount = RCCSeqOfTriWithPairs(cs.Cards)
//	} else if IsBomb(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_BOMB
//		cs.Rank = cs.Cards[0].Value
//		cs.SeqCount = 1
//		cs.SubCount = 1
//	} else if IsRocket(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_ROCKET
//		cs.Rank = cs.Cards[0].Value
//		cs.SeqCount = 0
//		cs.SubCount = 0
//	} else if IsQuadWithSingles(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_QUADPLEX_WITH_SINGLES
//		cs.Rank = cs.Cards[2].Value
//		cs.SeqCount = 0
//		cs.SubCount = 2
//	} else if IsQuadWithPairs(cs.Cards) {
//		cs.Pattern = cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS
//		{
//			if IsBomb(cs.Cards[0:4]) {
//				cs.Rank = cs.Cards[0].Value
//			} else if IsBomb(cs.Cards[4:8]) {
//				cs.Rank = cs.Cards[len(cs.Cards)-1].Value
//			} else if IsBomb(cs.Cards[2:6]) {
//				cs.Rank = cs.Cards[2].Value
//			}
//
//			cs.SeqCount = 1
//			cs.SubCount = 2
//		}
//	} else {
//		cs.Pattern = cardConst.CARD_PATTERN_ERROR
//	}
//
//	return cs.Pattern
//}
