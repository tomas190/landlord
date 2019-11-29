package game

import (
	"fmt"
	"github.com/wonderivan/logger"
	"testing"
)

/*
// 牌型
const (
	_CARD_PATTERN                                           = iota
	CARD_PATTERN_TODO                                        // 1待判定类型
	CARD_PATTERN_ERROR                                       // 2非法类型
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
func TestCalPattern(t *testing.T) {
	lastCard := []*Card{
		//{1, 1}, {1, 2}, {1, 3},// {1, 4},
		//{2, 1}, {2, 2}, {2, 3}, {2, 4},
		//{3, 1}, {3, 2}, //{3, 3}, {3, 4},
		//{4, 1}, {4, 2}, {4, 3}, //{4, 4},
		{5, 1}, {5, 2}, {5, 3}, //{5, 4},
		{6, 1}, {6, 2}, {6, 3}, //{6, 4},
		{7, 1}, {7, 2}, {7, 3}, //{7, 4},
		//{8, 1}, {8, 2}, {8, 3}, {8, 4},
		//{9, 1}, {9, 2}, {9, 3}, {9, 4},
		//{10, 1}, //{10, 2}, {10, 3}, {10, 4},
		//{11, 1}, {11, 2}, {11, 3}, {11, 4},
		{12, 1}, {12, 2}, {12, 3},// {12, 4},
		//{13, 1}, {13, 2}, {13, 3}, {13, 4},
	}
	//outCard := []*Card{{14, 1}, {14, 2}, {14, 4}, {14, 3}}

	//logger.Debug(CanBeat(lastCard, outCard))
	cardsType := GetCardsType(lastCard)
	//cardsType := IsSeqOfTriWithPairsFix(lastCard)
	cardsType := IsSeqOfTriWithSinglesFix(lastCard)
	fmt.Println("===========:", cardsType)

}

func TestCardCount(t *testing.T) {
	hands := []*Card{{12, 1}, {12, 2}, {12, 2}}
	roomThrows := []*Card{{12, 1},}

	cards := countCards(hands, roomThrows)

	m := cards.CardCount

	for k, v := range m {
		fmt.Printf("%d -- %d", k, v)
		fmt.Println()
	}

}
func TestFindMinSingle(t *testing.T) {
	hands := []*Card{{11, 1}, {12, 1}, {12, 2}, {12, 2}, {14, 1}, {15, 1},}

	SortCard(hands)
	PrintCard(hands)

	cards, b := findMinSingle(hands)
	fmt.Println("是否有单张:", b)
	if b {
		fmt.Println(cards[0].Value)
	}
	PrintCard(hands)

}

func TestFindMinDouble(t *testing.T) {
	hands := []*Card{{10, 1}, {3, 1}, {3, 1}, {5, 1}, {3, 1}, {3, 1},
		{12, 1}, {12, 1}, {12, 1}, {13, 3}, {13, 1}, {15, 1},}

	PrintCard(hands)

	cards, b := findMinTriple(hands)
	//fmt.Println("是否有最小的对子:", b)
	fmt.Println("是否有最小的三张:", b) //findMinTriple
	//fmt.Println("是否有最小的炸弹:", b) // findMinBoom
	if b {
		PrintCard(cards)
	}
	PrintCard(hands)

}

func TestInclude(t *testing.T) {
	a := []*Card{
		//{1, 1}, {1, 2}, {1, 3},// {1, 4},
		//{2, 1}, {2, 2}, {2, 3}, {2, 4},
		//{3, 1}, {3, 2}, {3, 3}, {3, 4},
		//{4, 1}, //{4, 2}, {4, 3}, {4, 4},
		//{5, 1}, //{5, 2}, {5, 3},{5, 4},
		{6, 1}, //{6, 2}, {6, 3}, {6, 4},
	}
	logger.Debug(getAllDoubleNum(a))
}
