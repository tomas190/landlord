package game

import (
	"fmt"
	"landlord/mconst/cardConst"
)

/*
const (
	CARD_RANK_THREE       = iota + 1 // 1
	CARD_RANK_FOUR                   // 2
	CARD_RANK_FIVE                   // 3
	CARD_RANK_SIX                    // 4
	CARD_RANK_SEVEN                  // 5
	CARD_RANK_EIGHT                  // 6
	CARD_RANK_NINE                   // 7
	CARD_RANK_TEN                    // 8
	CARD_RANK_JACK                   // 9
	CARD_RANK_QUEEN                  // 10
	CARD_RANK_KING                   // 11
	CARD_RANK_ACE                    // 12
	CARD_RANK_TWO                    // 13
	CARD_RANK_BLACK_JOKER            // 14
	CARD_RANK_RED_JOKER              // 15
)


好牌：

大小王：0.9
2：0.85
A：0.8
K：0.75
Q：0.7
J：0.65
10：0.6
9：0.55
8：0.5
7：0.45
6：0.4
5：0.35
4：0.3
3：0.25

差牌：
3：0.9
4：0.85
5：0.8
6：0.75
7：0.7
8：0.65
9：0.6
10：0.55
J：0.5
Q：0.45
K：0.4
A：0.35
2：0.3
大小王：0.25
*/

func PercentageOfGoodCards() int32 {

	/*
		大小王：0.9
		2：0.85
		A：0.8
		K：0.75
		Q：0.7
		J：0.65
		10：0.6
		9：0.55
		8：0.5
		7：0.45
		6：0.4
		5：0.35
		4：0.3
		3：0.25
	*/
	num := RandNum(1, 100)
	if num <= 90 {
		if num%2 == 0 {
			return cardConst.CARD_RANK_RED_JOKER
		} else {
			return cardConst.CARD_RANK_BLACK_JOKER
		}
	}
	if num <= 85 {
		return cardConst.CARD_RANK_TWO
	}

	if num <= 80 {
		return cardConst.CARD_RANK_ACE
	}

	if num <= 75 {
		return cardConst.CARD_RANK_KING
	}

	if num <= 70 {
		return cardConst.CARD_RANK_QUEEN
	}

	if num <= 65 {
		return cardConst.CARD_RANK_JACK
	}

	if num <= 60 {
		return cardConst.CARD_RANK_TEN
	}

	if num <= 55 {
		return cardConst.CARD_RANK_NINE
	}
	if num <= 50 {
		return cardConst.CARD_RANK_EIGHT
	}
	if num <= 45 {
		return cardConst.CARD_RANK_SEVEN
	}
	if num <= 40 {
		return cardConst.CARD_RANK_SIX
	}
	if num <= 35 {
		return cardConst.CARD_RANK_FIVE
	}
	if num <= 30 {
		return cardConst.CARD_RANK_FOUR
	}
	return cardConst.CARD_RANK_THREE
}

func GetGoodCards(cards []*Card) ([]*Card, []*Card) {

	itemCards := append([]*Card{}, cards...)

	var result []*Card

	for i := 0; ; i++ {
		if len(result) >= 17 {
			break
		}

		c := PercentageOfGoodCards()
		fmt.Println(c)

		card, temp := GetAndRemoveCard(itemCards, c)
		if card != nil {
			fmt.Println("nil=================")
			itemCards = temp
			result = append(result, card)
		}

	}
	return result, itemCards
}

// 将牌送 牌队中移除
func GetAndRemoveCard(cards []*Card, removeCardValue int32) (*Card, []*Card) {
	var result *Card
	for i := 0; i < len(cards); i++ {
		if cards[i].Value == removeCardValue {
			result = cards[i]
			cards = append(cards[:i], cards[i+1:]...)
			break
		}
	}
	return result, cards
}

func getTwoPercentage() bool {
	return false
}
