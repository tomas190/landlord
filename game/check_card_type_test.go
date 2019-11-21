package game

import (
	"fmt"
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
	card := []*Card{{12, 1}, {12, 1},{12, 1} ,
		{11, 1}, {11, 1},{11, 1},
	{9,2},{9,2}}
	set := NewCardSet(card)

	pattern := CalPattern(set)
	fmt.Println(pattern)

}
