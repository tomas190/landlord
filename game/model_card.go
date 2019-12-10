package game

import (
	"fmt"
	"landlord/mconst/cardConst"
	"sort"
)

// 定义扑克牌
type Card struct {
	Value int32 // card值用于排序比较
	Suit  int32 // card花色
}

// 创建一副有序的牌
func CreateSortCard() []*Card {
	result := initOriginalCard()
	SortCard(result)
	return result
}

// 创建一副乱序的牌
func CreateBrokenCard() []*Card {
	var result []*Card
	destiny := RandNum(0, 90)
	if destiny >= 0 || destiny <= 30 {
		result = initOriginalCard()
	} else if destiny >= 31 || destiny <= 60 {
		result = initOriginalCard2()
	} else {
		result = initOriginalCard3()
	}

	dSort := RandNum(25, 35)
	OutOfCardNotDeep(result, dSort)
	return result
}

func initOriginalCard() []*Card {
	var result []*Card
	for j := 1; j <= 13; j++ {
		for i := 1; i <= 4; i++ {
			var card Card
			card.Value = int32(j)
			card.Suit = int32(i)
			result = append(result, &card)
		}
	}
	var bigCard Card
	var smlCard Card
	bigCard.Value = cardConst.CARD_RANK_RED_JOKER
	bigCard.Suit = cardConst.CARD_SUIT_JOKER
	smlCard.Value = cardConst.CARD_RANK_BLACK_JOKER
	smlCard.Suit = cardConst.CARD_SUIT_JOKER

	result = append(result, &bigCard, &smlCard)
	return result
}

func initOriginalCard2() []*Card {
	var result []*Card

	for i := 1; i <= 4; i++ {
		for j := 1; j <= 13; j++ {
			var card Card
			card.Value = int32(j)
			card.Suit = int32(i)
			result = append(result, &card)
		}
	}
	var bigCard Card
	var smlCard Card
	bigCard.Value = cardConst.CARD_RANK_RED_JOKER
	bigCard.Suit = cardConst.CARD_SUIT_JOKER
	smlCard.Value = cardConst.CARD_RANK_BLACK_JOKER
	smlCard.Suit = cardConst.CARD_SUIT_JOKER

	result = append(result, &bigCard, &smlCard)
	return result
}

func initOriginalCard3() []*Card {
	var result []*Card
	var bigCard Card
	var smlCard Card
	bigCard.Value = cardConst.CARD_RANK_RED_JOKER
	bigCard.Suit = cardConst.CARD_SUIT_JOKER
	smlCard.Value = cardConst.CARD_RANK_BLACK_JOKER
	smlCard.Suit = cardConst.CARD_SUIT_JOKER
	for i := 1; i <= 4; i++ {
		for j := 7; j <= 13; j++ {
			var card Card
			card.Value = int32(j)
			card.Suit = int32(i)
			result = append(result, &card)
		}
	}

	for i := 1; i <= 4; i++ {
		for j := 1; j <= 6; j++ {
			var card Card
			card.Value = int32(j)
			card.Suit = int32(i)
			result = append(result, &card)
		}
	}
	result = append(result, &bigCard, &smlCard)
	return result
}

// 随机乱序
func OutOfCard(arr []*Card) {
	for i := len(arr) - 1; i > 0; i-- {
		num := RandNum(0, 53)
		arr[i], arr[num] = arr[num], arr[i]
	}
}

// 随机乱序
func OutOfCardNotDeep(arr []*Card, deepLevel int) {
	for i := len(arr) - 1; i > deepLevel; i-- {
		num := RandNum(0, 53)
		arr[i], arr[num] = arr[num], arr[i]
	}
}

// 排序  从大到小
func SortCard(cards []*Card) {
	v := func(c1, c2 *Card) bool {
		return c1.Value > c2.Value
	}

	s := func(c1, c2 *Card) bool {
		return c1.Suit > c2.Suit
	}
	OrderedBy(v, s).Sort(cards)
}

// 排序从小到大
func SortCardSL(cards []*Card) {
	v := func(c1, c2 *Card) bool {
		return c1.Value < c2.Value
	}

	s := func(c1, c2 *Card) bool {
		return c1.Suit > c2.Suit
	}
	OrderedBy(v, s).Sort(cards)
}

func PrintCard(cards []*Card) {
	for i := 0; i < len(cards); i++ {
		fmt.Print(cards[i])
		fmt.Print(",")
	}
	fmt.Println()
}

// =======================  sort =======================
type lessFunc func(p, p1 *Card) bool

type multiSort struct {
	cards []*Card
	less  []lessFunc
}

func (ms *multiSort) Sort(changes []*Card) {
	ms.cards = changes
	sort.Sort(ms)
}

func OrderedBy(less ...lessFunc) *multiSort {
	return &multiSort{
		less: less,
	}
}

func (ms *multiSort) Len() int {
	return len(ms.cards)
}

func (ms *multiSort) Swap(i, j int) {
	ms.cards[i], ms.cards[j] = ms.cards[j], ms.cards[i]
}

func (ms *multiSort) Less(i, j int) bool {
	p, q := &ms.cards[i], &ms.cards[j]
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(*p, *q):
			// p < q, so we have a decision.
			return true
		case less(*q, *p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](*p, *q)
}
