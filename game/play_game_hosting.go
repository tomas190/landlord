package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/cardConst"
	"sort"
)

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
	singleCards := getHasNumsCard(handCards, 1)

	if len(singleCards) == 0 {
		return nil, false
	}

	sort.Ints(singleCards)
	var result []*Card

	for i := 0; i < len(handCards); i++ {
		if int(handCards[i].Value) == singleCards[0] {
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
	doubleArr := getHasNumsCard(handCards, 2)
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
	triple := getHasNumsCard(handCards, 3)
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

// 统计[]*card中 相同nums数量的牌有哪些
/*
	eg [{3},{3},{2},{2},{6}]
	nums 2  返回  [3,2]  ,3和2 都有两张
  , nums 1  返回  [6]
*/
func getHasNumsCard(handCards []*Card, nums int) []int {
	// 先统计牌的张输
	var counts []int
	cardCount := make(map[int32]int, len(handCards))
	for i := 0; i < len(handCards); i++ {
		cardCount[handCards[i].Value] = cardCount[handCards[i].Value] + 1
	}

	for k, v := range cardCount {
		if v == nums {
			counts = append(counts, int(k))
		}
	}
	sort.Ints(counts)
	return counts
}

/* ================================= 托管必出牌抽取 ============================*/

// 判断是否有王炸 有王炸
func hasRacket(handCards []*Card) ([]*Card, bool) {
	if len(handCards) < 2 {
		return nil, false
	}

	var flag int
	var rackets []*Card
	for i := 0; i < len(handCards); i++ {
		if handCards[i].Value >= cardConst.CARD_RANK_BLACK_JOKER {
			flag++
			rackets = append(rackets, handCards[i])
		}
	}

	if flag == 2 {
		return rackets, true
	}

	return nil, false

}

/*
从手牌中 找到这张牌
*/
func findThisValueCard(value int, handCards []*Card, cardNum int) []*Card {
	var countCardNum int
	var result []*Card
	for i := 0; i < len(handCards); i++ {
		if int(handCards[i].Value) == value {
			result = append(result, handCards[i])
			countCardNum++
			if countCardNum == cardNum {
				break
			}
		}
	}
	return result
}

/*
  返回手牌中的最小的炸弹
*/
func findMinBoom(handCards []*Card) ([]*Card, bool) {
	// 先统计牌的张输
	boom := getHasNumsCard(handCards, 4)
	if len(boom) == 0 {
		return nil, false
	}

	sort.Ints(boom)
	var result []*Card

	for i := 0; i < len(handCards); i++ {
		if int(handCards[i].Value) == boom[0] {
			result = append(result, handCards[i])
		}
		if len(result) == 4 {
			break
		}
	}
	return result, true
}

/* ================================= 托管上家牌有出必出牌抽取 ============================*/

/*
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
	CARD_PATTERN_QUADPLEX_WITH_PAIRS                         // 16四带两对            J-J-J-J-9-9-Q-Q.
)
*/

/*

  一.上家牌 【单张】 下家托管出牌规则
  1.有单张的情况下出能打过上家牌的最小单张
  2.无单张的情况下先拆对 在拆三
  3.如以上请款均无 有炸弹 则出最小炸
  4.大小王不拆 直接炸

*/
// 单张
func HostingBeatSingle(handCards, eCards []*Card) ([]*Card, bool) {
	if len(eCards) != 1 || len(handCards) <= 0 {
		logger.Error("无效牌值 !!!incredible")
		return nil, false
	}

	// 如果最后一张牌是大王  则改成和小王同级
	if handCards[0].Value == cardConst.CARD_RANK_RED_JOKER {
		handCards[0].Value = cardConst.CARD_RANK_BLACK_JOKER
		defer func() {
			handCards[0].Value = cardConst.CARD_RANK_RED_JOKER
		}()
	}

	// 1.获取所有单张
	numSingle := getHasNumsCard(handCards, 1)
	for i := 0; i < len(numSingle); i++ {
		if numSingle[i] > int(eCards[0].Value) {
			// 如果这张牌 和下一长牌都是鬼  不拆炸弹原则
			if numSingle[i] >= cardConst.CARD_RANK_BLACK_JOKER {
				if i+1 == len(numSingle)-1 && numSingle[i+1] == cardConst.CARD_RANK_RED_JOKER {
					break
				}
			}
			result := findThisValueCard(numSingle[i], handCards, 1)
			return result, true
		}
	}

	// 2.如果没有单张 拆对
	numDouble := getHasNumsCard(handCards, 2)
	for i := 0; i < len(numDouble); i++ {
		if numDouble[i] > int(eCards[0].Value) {
			result := findThisValueCard(numDouble[i], handCards, 1)
			return result, true
		}
	}

	// 3.如果没有单张 和对 拆三
	numTriple := getHasNumsCard(handCards, 3)
	for i := 0; i < len(numTriple); i++ {
		if numTriple[i] > int(eCards[0].Value) {
			result := findThisValueCard(numTriple[i], handCards, 1)
			return result, true
		}
	}

	// 4.如果没有单张 和对 三 找炸弹直接炸
	bomb, b := findMinBoom(handCards)
	if b {
		return bomb, true
	}

	// 如果没有炸弹 找王炸
	return hasRacket(handCards)

}

/*

  一.上家牌 【对子】 下家托管出牌规则
  1.有 对子 的情况下出能打过上家牌的最小 对子
  2.无 对子 的情况下拆三
  3.如以上请款均无 有炸弹 则出最小炸
  4.大小王不拆 直接炸

*/
// 对子
func HostingBeatDouble(handCards, eCards []*Card) ([]*Card, bool) {
	if len(eCards) != 2 || len(handCards) <= 1 {
		logger.Error("无效牌值 !!!incredible")
		return nil, false
	}

	// 2.获取所有对
	numDouble := getHasNumsCard(handCards, 2)
	for i := 0; i < len(numDouble); i++ {
		if int(numDouble[i]) > int(eCards[0].Value) {
			result := findThisValueCard(numDouble[i], handCards, 2)
			return result, true
		}
	}

	// 3.如果没有单张 和对 拆三
	numTriple := getHasNumsCard(handCards, 3)
	for i := 0; i < len(numTriple); i++ {
		if numTriple[i] > int(eCards[0].Value) {
			result := findThisValueCard(numTriple[i], handCards, 2)
			return result, true
		}
	}

	// 4.如果没有单张 和对 三 找炸弹直接炸
	bomb, b := findMinBoom(handCards)
	if b {
		return bomb, true
	}

	// 如果没有炸弹 找王炸
	return hasRacket(handCards)
}

/*

  一.上家牌 【三张】 下家托管出牌规则
  1.有 三张 的情况下出能打过上家牌的最小 三张
  3.如以上请款均无 有炸弹 则出最小炸
  4.大小王不拆 直接炸

*/
// 三张
func HostingBeatTriple(handCards, eCards []*Card) ([]*Card, bool) {
	if len(eCards) != 3 {
		logger.Error("无效牌值 !!!incredible")
		return nil, false
	}

	// 3.获取所有三张
	numTriple := getHasNumsCard(handCards, 3)
	for i := 0; i < len(numTriple); i++ {
		if numTriple[i] > int(eCards[0].Value) {
			result := findThisValueCard(numTriple[i], handCards, 3)
			return result, true
		}
	}

	// 4.如果没有三张 和对 三 找炸弹直接炸
	bomb, b := findMinBoom(handCards)
	if b {
		return bomb, true
	}

	// 如果没有炸弹 找王炸
	return hasRacket(handCards)
}

// 找出一副牌中出现次数最多的牌值
// 如果都一样 则返回第一个
func getCardManyOne(cards []*Card) int32 {

	if len(cards) <= 0 {
		logger.Error("建议不要这样调用 !!!incredible")
		return 0
	}

	var mostCountKey int32
	cardCount := make(map[int32]int, len(cards))
	for i := 0; i < len(cards); i++ {
		cardCount[cards[i].Value] = cardCount[cards[i].Value] + 1
	}
	var mostNum int
	for l, v := range cardCount {
		if v > mostNum {
			mostCountKey = l
		}
	}
	logger.Debug("主要值", mostCountKey)
	return mostCountKey
}

/*

  一.上家牌 【三带一】 下家托管出牌规则
  1.有 三张 的情况下出能打过上家牌的最小 三张
  3.如以上请款均无 有炸弹 则出最小炸
  4.大小王不拆 直接炸

*/
// 三带一
func HostingBeatTripleWithSingle(handCards, eCards []*Card) ([]*Card, bool) {
	cards, hasRacket := hasRacket(handCards)
	if len(eCards) != 4 || (!hasRacket && len(handCards) <= 3) {
		logger.Error("牌数量不满足检测条件...")
		return nil, false
	}

	// 3.获取所有三张
	numTriple := getHasNumsCard(handCards, 3)
	for i := 0; i < len(numTriple); i++ {
		if numTriple[i] > int(getCardManyOne(eCards)) { // 代表已经找到了三张
			logger.Debug("已经找到了能打过的三张")
			result := findThisValueCard(numTriple[i], handCards, 3)
			// 寻找一张最小单牌
			tmpCards := removeCards(handCards, result) // 先移除掉已经找到的三张牌
			cards, b := HostingBeatSingle(tmpCards, []*Card{{0, 0,}})
			if !b {
				return nil, false
			}
			result = append(result, cards[0])
			return result, true
		}
	}

	// 4.如果没有三张 和对 三 找炸弹直接炸
	bomb, b := findMinBoom(handCards)
	if b {
		return bomb, true
	}

	// 如果没有炸弹 找王炸
	return cards, hasRacket
}

/*

  一.上家牌 【三张对】 下家托管出牌规则
  1.有 三张 的情况下出能打过上家牌的最小 三张
  3.如以上请款均无 有炸弹 则出最小炸
  4.大小王不拆 直接炸

*/
// 三带对
func HostingBeatTripleWithDouble(handCards, eCards []*Card) ([]*Card, bool) {
	cards, containRacket := hasRacket(handCards)
	if len(eCards) != 5 || (!containRacket && len(handCards) <= 4) {
		logger.Error("牌数量不满足检测条件...")
		return nil, false
	}

	mainKey := getCardManyOne(eCards)
	// 3.获取所有三张
	numTriple := getHasNumsCard(handCards, 3)
	for i := 0; i < len(numTriple); i++ {
		if numTriple[i] > int(mainKey) { // 代表已经找到了三张
			//logger.Debug("已经找到了能打过的三张")
			result := findThisValueCard(numTriple[i], handCards, 3)
			// 寻找一张最小单牌
			tmpCards := removeCards(handCards, result) // 先移除掉已经找到的三张牌
			cards, b := HostingBeatDouble(tmpCards, []*Card{{0, 0,}, {0, 0,}})
			PrintCard(tmpCards)

			if !b {
				return nil, false
			} else if _, b := hasRacket(cards); b { // 如果返回的是王炸
				goto bomb
			}

			result = append(result, cards[:2]...)
			return result, true
		}
	}

bomb:
	// 4.如果没有三张 和对 三 找炸弹直接炸
	bomb, b := findMinBoom(handCards)
	if b {
		return bomb, true
	}

	// 如果没有炸弹 找王炸
	return cards, containRacket
}

/*

  一.上家牌 【顺子】 下家托管出牌规则
  1.有 顺子 的情况下出能打过上家牌的最小 顺子 // 对子 三张 都可以拆  炸弹不能拆
  3.如以上情况均无 有炸弹 则出最小炸
  4.直接炸

*/
// 顺子
func HostingBeatJunko(handCards, eCards []*Card) ([]*Card, bool) {
	cards, hasRacket := hasRacket(handCards)
	junkoLen := len(eCards)
	if junkoLen < 5 || (!hasRacket && len(handCards) < junkoLen) {
		logger.Error("牌数量不满足检测条件...")
		return nil, false
	}

	// 1.先去掉手牌中的重复值 要去掉 2 以上大的牌 (2 以上大的牌不能组成顺子)
	singleCards := junkoHelpRemove(handCards)
	logger.Debug(singleCards)
	// 有机会组成大过的顺子
	if len(singleCards) > junkoLen {
		// logger.Debug("step1....")
		SortCardSL(eCards) // 将比过牌 从小到大排序
		startValue := eCards[0].Value
		endValue := eCards[junkoLen-1].Value

		for i := 0; i < len(singleCards); i++ {
			if singleCards[i] > int(startValue) { // 如果优质大于顺子开始值
				// logger.Debug("step2....")
				if i+junkoLen <= len(singleCards) && int(startValue)-singleCards[i] == int(endValue)-singleCards[i+junkoLen-1] {
					// logger.Debug("step3....")
					logger.Debug("胜利在望")
					var result []*Card
					for t := i; t <= i+junkoLen-1; t++ {
						tmpCard := findThisValueCard(singleCards[t], handCards, 1)
						result = append(result, tmpCard...)
					}
					return result, true
				}
			}
		}
	}
	//
	bomb, b := findMinBoom(handCards)
	if b {
		return bomb, true
	}
	// 如果没有炸弹 找王炸
	return cards, hasRacket
}

// 去掉重复值 如果有四个一样的则直接排除 2 大王 小王
// 因为2 大王小王 是组不成顺子的  炸弹不能拆
func junkoHelpRemove(cards []*Card) []int {
	// 先统计牌的张输
	cardCount := make(map[int32]int, len(cards))
	for i := 0; i < len(cards); i++ {
		if !(cards[i].Value >= cardConst.CARD_RANK_TWO) {
			cardCount[cards[i].Value] = cardCount[cards[i].Value] + 1
		}
	}

	var result []int
	for k, v := range cardCount {
		if v != 4 {
			result = append(result, int(k))
		}
	}

	// 排好序
	sort.Ints(result)
	return result
}

/*

  一.上家牌 【连对】 下家托管出牌规则
  1.有 连对 的情况下出能打过上家牌的最小 顺子 // 对子 三张 都可以拆  炸弹不能拆
  3.如以上情况均无 有炸弹 则出最小炸
  4.直接炸

*/
// 连对
func HostingBeatContinuouslyDouble(handCards, eCards []*Card) ([]*Card, bool) {
	cards, hasRacket := hasRacket(handCards)
	cDoubleLen := len(eCards)
	cDoubleJunkoLen := len(eCards) / 2
	if cDoubleLen < 6 || (!hasRacket && len(handCards) < cDoubleLen) {
		logger.Error("牌数量不满足检测条件...")
		return nil, false
	}

	// 1.先去掉手牌中的重复值 要去掉 2 以上大的牌 (2 以上大的牌不能组成 连队)
	canDoubleCards := continuouslyDoubleHelpRemove(handCards)
	logger.Debug(canDoubleCards)

	// 有机会组成大过的 连队
	if len(canDoubleCards) >= cDoubleJunkoLen {
		//logger.Debug("step1.........")
		SortCardSL(eCards) // 将比过牌 从小到大排序
		startValue := eCards[0].Value
		endValue := eCards[cDoubleLen-1].Value

		for i := 0; i < len(canDoubleCards); i++ {
			if canDoubleCards[i] > int(startValue) { //
				//logger.Debug("step2.......")
				if i+cDoubleJunkoLen <= len(canDoubleCards) && int(startValue)-canDoubleCards[i] == int(endValue)-canDoubleCards[i+cDoubleJunkoLen-1] {
					//logger.Debug("step3.......")
					var result []*Card
					for t := i; t <= i+cDoubleJunkoLen-1; t++ {
						tmpCard := findThisValueCard(canDoubleCards[t], handCards, 2)
						result = append(result, tmpCard...)
					}
					return result, true
				}
			}
		}
	}

	//
	bomb, b := findMinBoom(handCards)
	if b {
		return bomb, true
	}
	// 如果没有炸弹 找王炸
	return cards, hasRacket
}

// 连队移除帮助
// 寻找所有对子和三个 移除炸弹 2 大王 小王
// 因为2 大王小王 是组不成连队的  炸弹不能拆
func continuouslyDoubleHelpRemove(cards []*Card) []int {
	// 先统计牌的张输
	cardCount := make(map[int32]int, len(cards))
	for i := 0; i < len(cards); i++ {
		if !(cards[i].Value >= cardConst.CARD_RANK_TWO) {
			cardCount[cards[i].Value] = cardCount[cards[i].Value] + 1
		}
	}

	var result []int
	for k, v := range cardCount {
		if v == 2 || v == 3 {
			result = append(result, int(k))
		}
	}

	// 排好序
	sort.Ints(result)
	return result
}

// CARD_PATTERN_SEQUENCE_OF_TRIPLETS                        // 10飞机不带翅膀        4-4-4-5-5-5.
/*

  一.上家牌 【飞机】 下家托管出牌规则
  1.有 飞机 的情况下出能打过上家牌的最小飞机 炸弹不能拆
  3.如以上情况均无 有炸弹 则出最小炸
  4.直接炸

*/
// 飞机 不带牌
func HostingBeatTriplets(handCards, eCards []*Card) ([]*Card, bool) {
	cards, hasRacket := hasRacket(handCards)
	cTripletsLen := len(eCards)
	cTripletsJunkoLen := len(eCards) / 3
	if cTripletsLen < 6 || (!hasRacket && len(handCards) < cTripletsLen) {
		logger.Error("牌数量不满足检测条件...")
		return nil, false
	}

	// 1.先去掉手牌中的重复值 要去掉 2 以上大的牌 (2 以上大的牌不能组成 飞机)
	canTripletsCards := tripletsHelpRemove(handCards)
	logger.Debug(canTripletsCards)

	// 有机会组成大过的 飞机
	if len(canTripletsCards) >= cTripletsJunkoLen {
		//logger.Debug("step1.........")
		SortCardSL(eCards) // 将比过牌 从小到大排序
		startValue := eCards[0].Value
		endValue := eCards[cTripletsLen-1].Value

		for i := 0; i < len(canTripletsCards); i++ {
			if canTripletsCards[i] > int(startValue) { //
				//logger.Debug("step2.......")
				if i+cTripletsJunkoLen <= len(canTripletsCards) && int(startValue)-canTripletsCards[i] == int(endValue)-canTripletsCards[i+cTripletsJunkoLen-1] {
					//logger.Debug("step3.......")
					var result []*Card
					for t := i; t <= i+cTripletsJunkoLen-1; t++ {
						tmpCard := findThisValueCard(canTripletsCards[t], handCards, 3)
						result = append(result, tmpCard...)
					}
					return result, true
				}
			}
		}
	}

	//
	bomb, b := findMinBoom(handCards)
	if b {
		return bomb, true
	}
	// 如果没有炸弹 找王炸
	return cards, hasRacket
}

func tripletsHelpRemove(cards []*Card) []int {
	// 先统计牌的张输
	cardCount := make(map[int32]int, len(cards))
	for i := 0; i < len(cards); i++ {
		if !(cards[i].Value >= cardConst.CARD_RANK_TWO) {
			cardCount[cards[i].Value] = cardCount[cards[i].Value] + 1
		}
	}

	var result []int
	for k, v := range cardCount {
		if v == 3 {
			result = append(result, int(k))
		}
	}

	// 排好序
	sort.Ints(result)
	return result
}

// CARD_PATTERN_SEQUENCE_OF_TRIPLETS                        // 10飞机不带翅膀        4-4-4-5-5-5.
/*

  一.上家牌 【飞机】 下家托管出牌规则
  1.有 飞机 的情况下出能打过上家牌的最小飞机 炸弹不能拆
  3.如以上情况均无 有炸弹 则出最小炸
  4.直接炸

*/
// 飞机 带单
func HostingBeatTripletsWithSingle(handCards, eCards []*Card) ([]*Card, bool) {
	cards, hasRacket := hasRacket(handCards)
	cTripletsLen := len(eCards)
	cTripletsJunkoLen := len(eCards) / 4
	if cTripletsLen < 6 || (!hasRacket && len(handCards) < cTripletsLen) {
		logger.Error("牌数量不满足检测条件...")
		return nil, false
	}

	eChangeCards := tripletsHelpRemove(eCards)

	// 1.先去掉手牌中的重复值 要去掉 2 以上大的牌 (2 以上大的牌不能组成 飞机)
	canTripletsCards := tripletsHelpRemove(handCards)
	logger.Debug(canTripletsCards)

	// 有机会组成大过的 飞机
	if len(canTripletsCards) >= cTripletsJunkoLen {
		logger.Debug("step1.........")
		SortCardSL(eCards) // 将比过牌 从小到大排序
		startValue := eChangeCards[0]
		endValue := eChangeCards[len(eChangeCards)-1]

		for i := 0; i < len(canTripletsCards); i++ {
			if canTripletsCards[i] > int(startValue) { //
				logger.Debug("step2.......")
				if i+cTripletsJunkoLen <= len(canTripletsCards) && int(startValue)-canTripletsCards[i] == int(endValue)-canTripletsCards[i+cTripletsJunkoLen-1] {
					logger.Debug("step3.......")
					var result []*Card
					for t := i; t <= i+cTripletsJunkoLen-1; t++ {
						tmpCard := findThisValueCard(canTripletsCards[t], handCards, 3)
						result = append(result, tmpCard...)
					}
					tmpCards := removeCards(handCards, result)
					// 在寻找飞机数量的单牌
					for i := 0; i < cTripletsJunkoLen; i++ {
						cards, b := HostingBeatSingle(tmpCards, []*Card{{0, 0,}})
						if !b {
							return nil, false
						}
						result = append(result, cards[0])
						tmpCards = removeCards(handCards, result)
					}
					return result, true
				}
			}
		}
	}

	//
	bomb, b := findMinBoom(handCards)
	if b {
		return bomb, true
	}
	// 如果没有炸弹 找王炸
	return cards, hasRacket
}

// 飞机 带一对
func HostingBeatTripletsWithDouble(handCards, eCards []*Card) ([]*Card, bool) {
	cards, containRacket := hasRacket(handCards)
	cTripletsLen := len(eCards)
	cTripletsJunkoLen := len(eCards) / 5
	if cTripletsLen < 6 || (!containRacket && len(handCards) < cTripletsLen) {
		logger.Error("牌数量不满足检测条件...")
		return nil, false
	}

	eChangeCards := tripletsHelpRemove(eCards)

	// 1.先去掉手牌中的重复值 要去掉 2 以上大的牌 (2 以上大的牌不能组成 飞机)
	canTripletsCards := tripletsHelpRemove(handCards)
	logger.Debug(canTripletsCards)

	// 有机会组成大过的 飞机
	if len(canTripletsCards) >= cTripletsJunkoLen {
		logger.Debug("step1.........")
		SortCardSL(eCards) // 将比过牌 从小到大排序
		startValue := eChangeCards[0]
		endValue := eChangeCards[len(eChangeCards)-1]

		for i := 0; i < len(canTripletsCards); i++ {
			if canTripletsCards[i] > int(startValue) { //
				logger.Debug("step2.......")
				if i+cTripletsJunkoLen <= len(canTripletsCards) && int(startValue)-canTripletsCards[i] == int(endValue)-canTripletsCards[i+cTripletsJunkoLen-1] {
					logger.Debug("step3.......")
					var result []*Card
					for t := i; t <= i+cTripletsJunkoLen-1; t++ {
						tmpCard := findThisValueCard(canTripletsCards[t], handCards, 3)
						result = append(result, tmpCard...)
					}
					tmpCards := removeCards(handCards, result)
					// 在寻找飞机数量的单牌
					for i := 0; i < cTripletsJunkoLen; i++ {
						cards, b := HostingBeatDouble(tmpCards, []*Card{{0, 0,}, {0, 0,}})
						if !b {
							return nil, false
						} else if _, b := hasRacket(cards); b { // 如果返回的是王炸
							goto bomb
						}

						result = append(result, cards[:2]...)
						tmpCards = removeCards(handCards, result)
					}
					return result, true
				}
			}
		}
	}

bomb:
	//
	bomb, b := findMinBoom(handCards)
	if b {
		return bomb, true
	}
	// 如果没有炸弹 找王炸
	return cards, containRacket
}

/*

  一.上家牌 【炸弹】 下家托管出牌规则
  3.如以上请款均无 有炸弹 则出最小炸
  4.大小王不拆 直接炸

*/
// 炸弹
func HostingBeatBomb(handCards, eCards []*Card) ([]*Card, bool) {
	if len(eCards) != 4 {
		logger.Error("无效牌值 !!!incredible")
		return nil, false
	}

	// 3.获取所有四张
	numTriple := getHasNumsCard(handCards, 4)
	for i := 0; i < len(numTriple); i++ {
		if numTriple[i] > int(eCards[0].Value) {
			result := findThisValueCard(numTriple[i], handCards, 4)
			return result, true
		}
	}

	// 如果没有炸弹 找王炸
	return hasRacket(handCards)
}

/*

  一.上家牌 【四带二单】 下家托管出牌规则
  3.如以上请款均无 有炸弹 则出最小炸
  4.大小王不拆 直接炸

*/
// 四带二单
func HostingBeatBombWithSingles(handCards, eCards []*Card) ([]*Card, bool) {
	if len(eCards) != 6 {
		logger.Error("无效牌值 !!!incredible")
		return nil, false
	}

	// 3.获取所有四张
	numTriple := getHasNumsCard(handCards, 4)
	for i := 0; i < len(numTriple); i++ {
		if numTriple[i] > int(eCards[0].Value) {
			result := findThisValueCard(numTriple[i], handCards, 4)
			tmpCards := removeCards(handCards, result)
			// 继续找两张单牌
			for i := 0; i < 2; i++ {
				cards, b := HostingBeatSingle(tmpCards, []*Card{{0, 0,}})
				if !b {
					return nil, false
				}
				result = append(result, cards[0])
				tmpCards = removeCards(handCards, result)
			}

			return result, true
		}
	}

	
	bomb, b := findMinBoom(handCards)
	if b {
		return bomb, true
	}

	// 如果没有炸弹 找王炸
	return hasRacket(handCards)
}

/*

  一.上家牌 【四带二对】 下家托管出牌规则
  3.如以上请款均无 有炸弹 则出最小炸
  4.大小王不拆 直接炸

*/
// 四带二对
func HostingBeatBombWithDouble(handCards, eCards []*Card) ([]*Card, bool) {
	if len(eCards) != 8 {
		logger.Error("无效牌值 !!!incredible")
		return nil, false
	}

	// 3.获取所有四张
	numTriple := getHasNumsCard(handCards, 4)
	for i := 0; i < len(numTriple); i++ {
		if numTriple[i] > int(eCards[0].Value) {
			result := findThisValueCard(numTriple[i], handCards, 4)
			tmpCards := removeCards(handCards, result)
			// 继续找两张单牌
			for i := 0; i < 2; i++ {
				cards, b := HostingBeatDouble(tmpCards, []*Card{{0, 0,}, {0, 0,}})
				if !b {
					return nil, false
				} else if _, b := hasRacket(cards); b { // 如果返回的是王炸
					goto bomb
				}
				result = append(result, cards[:2]...)
				tmpCards = removeCards(handCards, result)
			}
			return result, true
		}
	}
bomb:
	bomb, b := findMinBoom(handCards)
	if b {
		return bomb, true
	}
	// 如果没有炸弹 找王炸
	return hasRacket(handCards)
}
