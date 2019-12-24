package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/cardConst"
	"sort"
	"strconv"
	"strings"
)

type group struct {
	cardLen   int32
	cardGroup []cardInfo
}

type cardInfo struct {
	cardValue int32
	cardNum   int32
}

func GroupHandsCard(hands []*Card) GroupCard {
	var tmpGroup GroupCard
	remainCards := hands
	countGroup := continuouslyCountGroup(hands)
	for i := 0; i < len(countGroup); i++ {
		var reCards []*ReCard
		reCards, remainCards = groupGroup(remainCards, countGroup[i])
		if len(reCards) > 0 {
			for j := 0; j < len(reCards); j++ {
				if reCards[j] == nil {
					continue
				}

				switch reCards[j].CardType {
				case cardConst.CARD_PATTERN_SEQUENCE:
					tmpGroup.Junko = append(tmpGroup.Junko, reCards[j])
				case cardConst.CARD_PATTERN_SEQUENCE_OF_PAIRS:
					tmpGroup.JunkoDouble = append(tmpGroup.Junko, reCards[j])
				case cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS:
					tmpGroup.JunkTriple = append(tmpGroup.Junko, reCards[j])
				default:
					logger.Debug(reCards[j].CardType)
					PrintCard(reCards[j].Card)
					logger.Error("!!!无此组牌类型")
				}
			}
		}
	}

	// 组剩余牌
	groupCards := CreateGroupCard(remainCards)
	groupCards.JunkoDouble = tmpGroup.JunkoDouble
	groupCards.Junko = tmpGroup.Junko
	groupCards.JunkTriple = tmpGroup.JunkTriple
	return groupCards

}

/*

return
	[]*Card 组合好的牌
	int32	组合好的牌的类型
	[]*Card 剩余牌
*/
// 根据group进行拆分牌
func groupGroup(hands []*Card, g group) ([]*ReCard, []*Card) {
	groupLen := len(g.cardGroup) // 牌的去重张数

	if groupLen < 2 {
		return nil, hands
	} else if groupLen == 2 {
		return groupLen2(hands, g)
	} else if groupLen == 3 {
		return groupLen3(hands, g)
	} else if groupLen == 4 {
		return groupLen4(hands, g)
	} else if groupLen == 5 { // 去重后张数为5  groupLen 最大为15 及每张牌为3张的情况下
		return groupLen5(hands, g)
	} else if groupLen == 6 { // 去重后张数为6  groupLen 最大为18 及每张牌为3张的情况下
		return groupLen6(hands, g)
	} else if groupLen == 7 { // 去重后张数为7  groupLen 最大为21 取20 及每张牌为3张的情况下
		return groupLen7(hands, g)
	}

	logger.Debug("==============  超出预测范围 =============")
	logger.Debug("group:", g)
	return nil, hands
}

// 组成顺子
func groupJunko(hands []*Card, start, end int) (*ReCard, []*Card) {
	var junko []*Card
	for i := start; i <= end; i++ {
		ct := findThisValueCard(i, hands, 1)
		junko = append(junko, ct...)
		hands = removeCards(hands, ct)
	}
	var re ReCard
	re.Card = junko
	re.Wight = junko[len(junko)-1].Value
	re.CardType = cardConst.CARD_PATTERN_SEQUENCE
	return &re, hands
}

// 组成三不带飞机
func groupJunkoTriple(hands []*Card, start, end int) (*ReCard, []*Card) {
	var junko []*Card
	for i := start; i <= end; i++ {
		ct := findThisValueCard(i, hands, 3)
		junko = append(junko, ct...)
		hands = removeCards(hands, ct)
	}

	var re ReCard
	re.Card = junko
	re.Wight = junko[len(junko)-1].Value
	re.CardType = cardConst.CARD_PATTERN_SEQUENCE_OF_TRIPLETS

	return &re, hands
}

// 组成连队
func groupJunkoDouble(hands []*Card, start, end int) (*ReCard, []*Card) {
	var junko []*Card
	for i := start; i <= end; i++ {
		ct := findThisValueCard(i, hands, 2)
		junko = append(junko, ct...)
		hands = removeCards(hands, ct)
	}
	var re ReCard
	re.Card = junko
	re.Wight = junko[len(junko)-1].Value
	re.CardType = cardConst.CARD_PATTERN_SEQUENCE_OF_PAIRS
	return &re, hands
}

/*=============================== help ===============================*/

// 连续分组
func continuouslyCountGroup(hands []*Card) []group {
	var groups []group
	sp := continuouslyCountStr(hands)
	//logger.Debug("sp:", sp)
	gsp := strings.Split(sp, ",cut,")
	for i := 0; i < len(gsp); i++ {
		//logger.Error("send:", gsp[i])
		g := vertToGroup(gsp[i])
		groups = append(groups, g)
	}
	return groups
}

/*
	str 2-1,3-3,4-3
*/
func vertToGroup(str string) group {
	var result group
	cis := strings.Split(str, ",")
	for i := 0; i < len(cis); i++ {
		ci := strings.Split(cis[i], "-")
		if len(ci) != 2 {
			logger.Error("!!! don't do this", cis)
			logger.Error("!!! don't do this", cis[i])
			logger.Error("!!! don't do this", len(ci))
			continue
		}
		var cardIf cardInfo
		v, _ := strconv.Atoi(ci[0])
		n, _ := strconv.Atoi(ci[1])
		cardIf.cardValue = int32(v)
		cardIf.cardNum = int32(n)
		result.cardLen = result.cardLen + int32(n)
		result.cardGroup = append(result.cardGroup, cardIf)
	}
	return result
}

// 找出连续的并且统计其值的个数的组
func continuouslyCountStr(hands []*Card) string {
	// 1.先统计张数和数量
	cim := make(map[int32]int32, 11)
	var cii []int
	for i := 0; i < len(hands); i++ {
		// 把2和大王的移除
		if hands[i].Value >= cardConst.CARD_RANK_TWO {
			continue
		}
		cim[hands[i].Value] = cim[hands[i].Value] + 1
	}

	for k, v := range cim {
		if v < 4 {
			cii = append(cii, int(k))
		}
	}

	sort.Ints(cii)
	var res string
	for i := 0; i < len(cii); i++ {
		if i != 0 {
			if cii[i]-cii[i-1] != 1 {
				res = res + "cut,"
			}
		}
		if i == len(cii)-1 {
			res += strconv.Itoa(cii[i]) + "-" + strconv.Itoa(int(cim[int32(cii[i])]))
		} else {
			res += strconv.Itoa(cii[i]) + "-" + strconv.Itoa(int(cim[int32(cii[i])])) + ","
		}
	}

	return res
}

// 根据group判断里面有x张牌的有多少张 x 为该牌的重复值
// 且判断相同重复的值是否连续
func howManyCardByX(g group, x int32) ([]int, bool, int) {
	var num int
	var rp []int
	for i := 0; i < len(g.cardGroup); i++ {
		if g.cardGroup[i].cardNum == x {
			num++
			rp = append(rp, int(g.cardGroup[i].cardValue))
		}
	}
	return rp, isContinuously(rp), num

}

// 根据group判断里面有x张牌的有多少张 >=x 为该牌的重复值
// 且判断相同重复的值是否连续
/*
[]int 连续的牌的数组
bool 是否连续
int 满足条件的牌
*/

func howManyCardMoreX(g group, x int32) ([]int, bool, int) {
	var num int
	var rp []int
	for i := 0; i < len(g.cardGroup); i++ {
		if g.cardGroup[i].cardNum >= x {
			num++
			rp = append(rp, int(g.cardGroup[i].cardValue))
		}
	}
	return rp, isContinuously(rp), num

}

// 判断这个数组是否连续
func isContinuously(arr []int) bool {
	sort.Ints(arr)
	if len(arr) <= 1 {
		return true
	}
	if arr[len(arr)-1]-arr[0] == len(arr)-1 {
		return true
	}
	return false
}

// 判断这个数组是否有连续且连续长度=long
// 返回第一个长度为long 且连续的数组
func hasContinuouslyLonger(arr []int, long int) ([]int, bool) {
	//	logger.Debug(arr)
	sort.Ints(arr)
	if len(arr) < long {
		return nil, false
	}
	for i := 0; i < len(arr); i++ {
		if i > len(arr)-long {
			break
		}
		if arr[i+long-1]-arr[i] == long-1 {

			return append([]int{}, arr[i:i+long]...), true
		}
	}
	return nil, false
}

// 从nums 中移除掉 removes 的元素
func removeArrNum(nums, removes []int) []int {
	var s string
	for i := 0; i < len(removes); i++ {
		s = s + "," + strconv.Itoa(removes[i]) + ","
	}
	var result []int
	for i := 0; i < len(nums); i++ {
		tmp := "," + strconv.Itoa(nums[i]) + ","
		if !strings.Contains(s, tmp) {
			result = append(result, nums[i])
		}
	}

	return result
}
