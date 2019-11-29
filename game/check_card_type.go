package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/cardConst"
	"sort"
)

/*
	荷官——牌型信息识别
*/

// 初值和错误值
const (
	SET_COUNT_TODO  = -1
	SET_COUNT_ERROR = -2
	SET_RANK_TODO   = -3
	SET_RANK_ERROR  = -4
)

type CardSet struct {
	Pattern  int   // 这组牌的类型
	SeqCount int   // 主要数量
	SubCount int   // 带的数量
	Rank     int32 // 决定这组牌大小
	Cards    []*Card
}

func NewCardSet(cards []*Card) *CardSet {
	return &CardSet{
		Pattern:  cardConst.CARD_PATTERN_TODO,
		SeqCount: SET_COUNT_TODO,
		SubCount: SET_COUNT_TODO,
		Rank:     SET_RANK_TODO,
		Cards:    cards,
	}
}

func GetCardsType(cards []*Card) int32 {
	if len(cards) <= 0 {
		return 0
	}
	set := NewCardSet(cards)
	pattern := CalPattern(set)
	return int32(pattern)
}

// 荷官计算牌型并赋值
// tip 牌型计算这里挖的坑比较大
func CalPattern(cs *CardSet) int {
	if len(cs.Cards) < 1 || len(cs.Cards) > 20 {
		return cardConst.CARD_PATTERN_ERROR
	}

	sort.Slice(cs.Cards, func(i, j int) bool {
		return cs.Cards[i].Value < cs.Cards[j].Value
	})

	if IsSingle(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_SINGLE
		cs.Rank = cs.Cards[0].Value
		cs.SeqCount = 1
		cs.SubCount = 0
	} else if IsPair(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_PAIR
		cs.Rank = cs.Cards[0].Value
		cs.SeqCount = 1
		cs.SubCount = 0
	} else if IsTriplet(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_TRIPLET
		cs.Rank = cs.Cards[0].Value
		cs.SeqCount = 1
		cs.SubCount = 0
	} else if IsTripletWithSingle(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_TRIPLET_WITH_SINGLE
		cs.Rank = cs.Cards[1].Value
		cs.SeqCount = 1
		cs.SubCount = 1
	} else if IsTripletWithPair(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_TRIPLET_WITH_PAIR
		cs.Rank = cs.Cards[2].Value
		cs.SeqCount = 1
		cs.SubCount = 1
	} else if IsSequence(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_SEQUENCE
		cs.Rank = cs.Cards[0].Value
		cs.SeqCount = len(cs.Cards)
		cs.SubCount = 0
	} else if IsSeqOfPair(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_SEQUENCE_OF_PAIRS
		cs.Rank = cs.Cards[0].Value
		cs.SeqCount = len(cs.Cards) / 2
		cs.SubCount = 0
	} else if IsSeqOfTriplet(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS
		cs.Rank = cs.Cards[0].Value
		cs.SeqCount = len(cs.Cards) / 3
		cs.SubCount = 0
		// tip 这里比较复杂了
	//} else if IsSeqOfTriWithSingles(cs.Cards) {
	} else if IsSeqOfTriWithSinglesFix(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_SINGLES
		cs.Rank, cs.SeqCount, cs.SubCount = RCCSeqOfTriWithSingles(cs.Cards)
	//} else if IsSeqOfTriWithPairs(cs.Cards) {
	} else if IsSeqOfTriWithPairsFix(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS_WITH_ATTACHED_PAIRS
		cs.Rank, cs.SeqCount, cs.SubCount = RCCSeqOfTriWithPairs(cs.Cards)
	} else if IsBomb(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_BOMB
		cs.Rank = cs.Cards[0].Value
		cs.SeqCount = 1
		cs.SubCount = 1
	} else if IsRocket(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_ROCKET
		cs.Rank = cs.Cards[0].Value
		cs.SeqCount = 0
		cs.SubCount = 0
	} else if IsQuadWithSingles(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_QUADPLEX_WITH_SINGLES
		cs.Rank = cs.Cards[2].Value
		cs.SeqCount = 0
		cs.SubCount = 2
	} else if IsQuadWithPairs(cs.Cards) {
		cs.Pattern = cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS
		{
			if IsBomb(cs.Cards[0:4]) {
				cs.Rank = cs.Cards[0].Value
			} else if IsBomb(cs.Cards[4:8]) {
				cs.Rank = cs.Cards[len(cs.Cards)-1].Value
			} else if IsBomb(cs.Cards[2:6]) {
				cs.Rank = cs.Cards[2].Value
			}

			cs.SeqCount = 1
			cs.SubCount = 2
		}
	} else {
		cs.Pattern = cardConst.CARD_PATTERN_ERROR
	}

	return cs.Pattern
}

// 检查牌型
// tip 牌型检查还可以用有限状态机来实现，之后可以尝试

func IsSingle(set []*Card) bool {
	return len(set) == 1
}

func IsPair(set []*Card) bool {
	if len(set) == 2 {
		return set[0].Value == set[1].Value
	}

	return false
}

func IsTriplet(set []*Card) bool {
	if len(set) == 3 {
		return set[0].Value == set[1].Value && set[0].Value == set[2].Value
	}

	return false
}

func IsTripletWithSingle(set []*Card) bool {
	if len(set) != 4 || IsBomb(set) {
		return false
	}

	if set[0].Value == set[1].Value { // 单在后
		return IsTriplet(set[0:3])
	} else if set[2].Value == set[3].Value { // 单在前
		return IsTriplet(set[1:4])
	}

	return false
}

func IsTripletWithPair(set []*Card) bool {
	if len(set) != 5 {
		return false
	}

	if set[2].Value == set[4].Value && IsPair(set[0:2]) && IsTriplet(set[2:5]) ||
		set[2].Value == set[0].Value && IsPair(set[3:5]) && IsTriplet(set[0:3]) {
		return true
		// 带王对的情况
	} else if set[2].Value == set[4].Value && IsRocket(set[0:2]) && IsTriplet(set[2:5]) ||
		set[2].Value == set[0].Value && IsRocket(set[3:5]) && IsTriplet(set[0:3]) {
		return true
	}

	return false
}

func IsSequence(set []*Card) bool {
	if len(set) < 5 {
		return false
	}

	lastRank := int32(-1)
	// 因为是按照rank排好序的所以直接比较即可
	for i, c := range set {
		// 如果包含2或者王肯定不是顺子
		if c.Value == cardConst.CARD_RANK_TWO ||
			c.Value == cardConst.CARD_RANK_RED_JOKER ||
			c.Value == cardConst.CARD_RANK_BLACK_JOKER {
			return false
		}

		// 迭代
		if i == 0 || lastRank == c.Value-1 {
			lastRank = c.Value
			// 判断
		} else if lastRank != c.Value-1 {
			return false
		}
	}

	return true
}

func IsSeqOfPair(set []*Card) bool {
	// 张数小于6或者是单数
	if len(set) < 6 || len(set)%2 != 0 {
		return false
	}

	// 有2或者王
	for _, c := range set {
		if c.Value == cardConst.CARD_RANK_TWO ||
			c.Value == cardConst.CARD_RANK_RED_JOKER ||
			c.Value == cardConst.CARD_RANK_BLACK_JOKER {
			return false
		}
	}

	lastRank := int32(-1)
	for i := 0; i < len(set); i += 2 {
		// 前后两张不是对
		if !IsPair(set[i:i+2]) {
			return false
		}

		if i == 0 || lastRank == set[i].Value-1 {
			lastRank = set[i].Value
		} else if lastRank != set[i].Value-1 {
			return false
		}
	}

	return true
}

func IsBomb(set []*Card) bool {
	if len(set) != 4 {
		return false
	}

	if set[0].Value == set[1].Value &&
		set[0].Value == set[2].Value &&
		set[0].Value == set[3].Value {
		return true
	}

	return false
}

func IsRocket(set []*Card) bool {
	if len(set) != 2 {
		return false
	}

	if set[0].Value == cardConst.CARD_RANK_BLACK_JOKER && set[1].Value == cardConst.CARD_RANK_RED_JOKER ||
		set[0].Value == cardConst.CARD_RANK_RED_JOKER && set[1].Value == cardConst.CARD_RANK_BLACK_JOKER {
		return true
	}

	return false
}

// 四带二
func IsQuadWithSingles(set []*Card) bool {
	if len(set) != 6 {
		return false
	}

	// ****--
	if set[0].Value == set[2].Value && IsBomb(set[0:4]) {
		return true
	}
	// --****
	if set[5].Value == set[2].Value && IsBomb(set[2:6]) {
		return true
		// -****-
	}
	// --****
	if IsBomb(set[1:5]) {
		return true
	}

	return false
}

// 四带两对
func IsQuadWithPairs(set []*Card) bool {
	if len(set) != 8 {
		return false
	}

	if IsBomb(set[0:4]) {
		if IsPair(set[4:6]) && IsPair(set[6:8]) {
			return true
		}

		if IsPair(set[4:6]) && IsRocket(set[6:8]) {
			return true
		}
	}

	if IsBomb(set[4:8]) {
		if IsPair(set[0:2]) && IsPair(set[2:4]) {
			return true
		}
	}

	if IsBomb(set[2:6]) {
		if IsPair(set[0:2]) && IsPair(set[6:8]) {
			return true
		}

		if IsPair(set[0:2]) && IsRocket(set[6:8]) {
			return true
		}
	}

	return false
}

func IsSeqOfTriplet(set []*Card) bool {
	// 张数小于6或者张数不是3的倍数
	if len(set) < 6 || len(set)%3 != 0 {
		return false
	}

	// 有2或者王
	for _, c := range set {
		if c.Value == cardConst.CARD_RANK_TWO ||
			c.Value == cardConst.CARD_RANK_RED_JOKER ||
			c.Value == cardConst.CARD_RANK_BLACK_JOKER {
			return false
		}
	}

	lastRank := int32(-1)
	for i := 0; i < len(set); i += 3 {
		// 前后三张不是trip
		if !IsTriplet(set[i:i+3]) {
			return false
		}

		if i == 0 || lastRank == set[i].Value-1 &&
			lastRank == set[i+1].Value-1 &&
			lastRank == set[i+2].Value-1 {
			lastRank = set[i].Value
		} else if lastRank != set[i].Value-1 {
			return false
		}
	}

	return true
}

/*

	比较复杂的牌型

*/
func IsSeqOfTriWithSingles(set []*Card) bool {
	// 张数小于8或者张数不是4的倍数
	if len(set) < 8 || len(set)%4 != 0 {
		return false
	}

	// 组数
	pCount := len(set) / 4
	result := make(map[int32]int, pCount)

	for _, c := range set {
		if r, ok := result[c.Value]; ok {
			result[c.Value] = r + 1
		} else {
			result[c.Value] = 1
		}
	}

	var checkSeq []int32
	for k, r := range result {
		if r >= 3 {
			checkSeq = append(checkSeq, k)
		}
	}
	sort.Slice(checkSeq, func(i, j int) bool {
		return checkSeq[i] < checkSeq[j]
	})

	// 如果三个的数量比组数小或者里面有2
	if len(checkSeq) < pCount || Include(checkSeq, cardConst.CARD_RANK_TWO) {
		return false
	} else if len(checkSeq) == 2 {
		return IsSeq(checkSeq)
	} else if len(checkSeq) == 3 {
		return IsSeq(checkSeq)
	} else if len(checkSeq) == 4 {
		return IsSeq(checkSeq[0:len(checkSeq)-1]) || IsSeq(checkSeq[1:])
	} else if len(checkSeq) == 5 {
		return IsSeq(checkSeq[0:len(checkSeq)-1]) || IsSeq(checkSeq[1:])
	}

	return false
}

func IsSeqOfTriWithPairs(set []*Card) bool {
	// 张数小于10或者张数不是5的倍数
	if len(set) < 10 || len(set)%5 != 0 {
		return false
	}

	pCount := len(set) / 5
	result := make(map[int32]int, pCount)

	for _, c := range set {
		if r, ok := result[c.Value]; ok {
			result[c.Value] = r + 1
		} else {
			result[c.Value] = 1
		}
	}

	var checkSeq []int32
	var pairCount int
	for k, r := range result {
		if r >= 3 {
			checkSeq = append(checkSeq, k)
		}
		// 后面要处理王对
		if r == 2 {
			pairCount++
		}
	}
	sort.Slice(checkSeq, func(i, j int) bool {
		return checkSeq[i] < checkSeq[j]
	})

	// 如果三个的数量比组数小或者里面有2 注意4个2可以是带牌
	if len(checkSeq) < pCount || (Include(checkSeq, cardConst.CARD_RANK_TWO) && result[cardConst.CARD_RANK_TWO] != 4) {
		return false
	} else if len(checkSeq) == 2 {
		return IsSeq(checkSeq) && (pairCount == len(checkSeq)) ||
			(pairCount == (len(checkSeq)-1) &&
				result[cardConst.CARD_RANK_RED_JOKER] == 1 && result[cardConst.CARD_RANK_BLACK_JOKER] == 1)
	} else if len(checkSeq) == 3 {
		if IsSeq(checkSeq) && (pairCount == len(checkSeq)) ||
			(pairCount == (len(checkSeq)-1) &&
				result[cardConst.CARD_RANK_RED_JOKER] == 1 && result[cardConst.CARD_RANK_BLACK_JOKER] == 1) {
			return true
		}
		if IsSeq(checkSeq[0 : len(checkSeq)-1]) {
			return result[checkSeq[len(checkSeq)-1]] == 4
		}
		if IsSeq(checkSeq[1:]) {
			return result[checkSeq[0]] == 4
		}
	} else if len(checkSeq) == 4 {
		if IsSeq(checkSeq) && (pairCount == len(checkSeq)) ||
			(pairCount == (len(checkSeq)-1) &&
				result[cardConst.CARD_RANK_RED_JOKER] == 1 && result[cardConst.CARD_RANK_BLACK_JOKER] == 1) {
			return true
		}
		if IsSeq(checkSeq[0 : len(checkSeq)-1]) {
			return result[checkSeq[len(checkSeq)-1]] == 4 && pairCount == 2
		}
		if IsSeq(checkSeq[1:]) {
			return result[checkSeq[0]] == 4 && pairCount == 2
		}
	}

	return false
}

/*
	tip 因为这块挖的坑比较大，这两种比较复杂的牌型，单独打补丁
*/

func RCCSeqOfTriWithSingles(set []*Card) (rank int32, seqCount, subCount int) {

	// 组数
	pCount := len(set) / 4
	result := make(map[int32]int, pCount)

	for _, c := range set {
		if r, ok := result[c.Value]; ok {
			result[c.Value] = r + 1
		} else {
			result[c.Value] = 1
		}
	}

	var checkSeq []int32
	for k, r := range result {
		if r >= 3 {
			checkSeq = append(checkSeq, k)
		}
	}
	sort.Slice(checkSeq, func(i, j int) bool {
		return checkSeq[i] < checkSeq[j]
	})

	if (len(checkSeq) == 2 || len(checkSeq) == 3) && (pCount == len(checkSeq)) {
		return checkSeq[0], len(checkSeq), len(checkSeq)
	} else if len(checkSeq) == 4 || len(checkSeq) == 5 {
		if IsSeq(checkSeq) && pCount == len(checkSeq) {
			return checkSeq[0], len(checkSeq), len(checkSeq)
		}
		if IsSeq(checkSeq[0:len(checkSeq)-1]) && pCount == len(checkSeq)-1 {
			return checkSeq[0], len(checkSeq) - 1, len(checkSeq) - 1
		}
		if IsSeq(checkSeq[1:]) && pCount == len(checkSeq)-1 {
			return checkSeq[1], len(checkSeq) - 1, len(checkSeq) - 1
		}
	}

	return SET_RANK_ERROR, SET_COUNT_ERROR, SET_COUNT_ERROR
}

func RCCSeqOfTriWithPairs(set []*Card) (rank int32, seqCount, subCount int) {

	pCount := len(set) / 5
	result := make(map[int32]int, pCount)

	for _, c := range set {
		if r, ok := result[c.Value]; ok {
			result[c.Value] = r + 1
		} else {
			result[c.Value] = 1
		}
	}

	var checkSeq []int32
	var pairCount int
	for k, r := range result {
		if r >= 3 {
			checkSeq = append(checkSeq, k)
		}
		// tip 后面要处理王对
		if r == 2 {
			pairCount++
		}
	}
	sort.Slice(checkSeq, func(i, j int) bool {
		return checkSeq[i] < checkSeq[j]
	})

	if (len(checkSeq) == 2 || len(checkSeq) == 3) && (pCount == len(checkSeq)) {
		return checkSeq[0], len(checkSeq), len(checkSeq)
	} else if len(checkSeq) == 4 || len(checkSeq) == 5 {
		if IsSeq(checkSeq) && pCount == len(checkSeq) {
			return checkSeq[0], len(checkSeq), len(checkSeq)
		}
		if IsSeq(checkSeq[0:len(checkSeq)-1]) && pCount == len(checkSeq)-1 {
			return checkSeq[0], len(checkSeq) - 1, len(checkSeq) - 1
		}
		if IsSeq(checkSeq[1:]) && pCount == len(checkSeq)-1 {
			return checkSeq[1], len(checkSeq) - 1, len(checkSeq) - 1
		}
	}
	return SET_RANK_ERROR, SET_COUNT_ERROR, SET_COUNT_ERROR
}

func IsSeq(arr []int32) bool {
	if len(arr) < 2 {
		return false
	}
	lastI := int32(0)
	for i := range arr {
		if i == 0 || lastI == arr[i]-1 {
			lastI = arr[i]
		} else {
			return false
		}
	}

	return true
}

func Include(arr []int32, target int32) bool {
	for _, v := range arr {
		if target == v {
			return true
		}
	}

	return false
}

/*
 2019年11月29日14:42:50  优化fix  飞机监测
*/

func IsSeqOfTriWithSinglesFix(cards []*Card) bool {
	// 张数小于8或者张数不是4的倍数
	cardsLen := len(cards)
	if cardsLen < 8 || cardsLen%4 != 0 {
		return false
	}
	seqLen := cardsLen / 4
	seqNums := tripletsHelpRemoveFix(cards)
	logger.Debug("len:", len(seqNums))
	logger.Debug("seq:", seqLen)
	if len(seqNums) > seqLen {
		seqNums = removeNotSeq(seqNums)
	}
	if len(seqNums) == seqLen {
		for i := 0; i < seqLen; i++ {
			if i+1 == seqLen {
				break
			}
			if seqNums[i+1]-seqNums[i] != 1 || seqNums[i+1] >= cardConst.CARD_RANK_TWO {
				return false
			}
		}
	} else {
		return false
	}
	return true
}

// 测试
func IsSeqOfTriWithPairsFix(cards []*Card) bool {
	// 张数小于10或者张数不是5的倍数
	cardsLen := len(cards)
	if cardsLen < 10 || cardsLen%5 != 0 {
		return false
	}
	seqLen := cardsLen / 5
	seqNums := tripletsHelpRemoveFix(cards)
	//logger.Debug("len:", len(seqNums))
	//logger.Debug("seq:", seqLen)
	if len(seqNums) >= seqLen {
		if len(seqNums) == seqLen {
			if tdHelp(seqLen, seqNums, cards) {
				return true
			}
		} else {
			rReqNums := seqNums[1:]
			if tdHelp(seqLen, rReqNums, cards) {
				return true
			}

			eReqNums := seqNums[:len(seqNums)-1]
			if tdHelp(seqLen, eReqNums, cards) {
				return true
			}
		}
	}
	return false
}

// 移除一个数组和其他 数据不连续的
/*
	// 前提是排好序的数组
	eg :{4,5,6,7,8,10}    把10移除掉
   res :{4,5,6,7,8}
*/
func removeNotSeq(arr []int) []int {
	if len(arr) <= 2 {
		return arr
	}

	if arr[len(arr)-1]-arr[len(arr)-2] != 1 {
		return arr[:len(arr)-1]
	} else {
		return arr[1:]
	}

}

func tripletsHelpRemoveFix(cards []*Card) []int {
	// 先统计牌的张输
	cardCount := make(map[int32]int, len(cards))
	for i := 0; i < len(cards); i++ {
		if !(cards[i].Value >= cardConst.CARD_RANK_TWO) {
			cardCount[cards[i].Value] = cardCount[cards[i].Value] + 1
		}
	}

	var result []int
	for k, v := range cardCount {
		if v >= 3 {
			result = append(result, int(k))
		}
	}

	// 排好序
	sort.Ints(result)
	return result
}

// 获取牌中对子的数量
// [1,1,2,2,3,3,4,4] 返回4
// [1,1,2,2,4,4,5,6] 返回 3
func getAllDoubleNum(cards []*Card) int {
	//if len(cards)%2 != 0 {
	//	logger.Debug("sssssssss",len(cards))
	//	return -1
	//}

	var result int
	for i := 0; i < len(cards); i++ {
		if i+1 == len(cards) {
			break
		}

		if cards[i].Value == cards[i+1].Value {
			result++
			i++
		}
	}
	return result
}

func tdHelp(seqLen int, seqNums []int, cards []*Card) bool {
	//logger.Debug(seqNums)
	//PrintCard(cards)
	for i := 0; i < seqLen; i++ {
		if i+1 == seqLen {
			break
		}
		if seqNums[i+1]-seqNums[i] == 1 && seqNums[i+1] < cardConst.CARD_RANK_TWO {
			var threes []*Card
			for i := 0; i < len(seqNums); i++ {
				tmp := findThisValueCard(seqNums[i], cards, 3)
				threes = append(threes, tmp...)
			}
			remains := removeCards(cards, threes)
			//logger.Debug("remains:")
			PrintCard(remains)
			if getAllDoubleNum(remains) == seqLen {
				return true
			}
		}
	}
	return false
}
