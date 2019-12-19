package cardConst

// 定义扑克级别
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

// 定义扑克牌面值
const (
	CARD_NAME_THREE       = "3"
	CARD_NAME_FOUR        = "4"
	CARD_NAME_FIVE        = "5"
	CARD_NAME_SIX         = "6"
	CARD_NAME_SEVEN       = "7"
	CARD_NAME_EIGHT       = "8"
	CARD_NAME_NINE        = "9"
	CARD_NAME_TEN         = "10"
	CARD_NAME_JACK        = "J"
	CARD_NAME_QUEEN       = "Q"
	CARD_NAME_KING        = "K"
	CARD_NAME_ACE         = "A"
	CARD_NAME_TWO         = "2"
	CARD_NAME_BLACK_JOKER = "BlackJoker"
	CARD_NAME_RED_JOKER   = "RedJoker"
)

// 定义扑克牌花色
const (
	// CARD_SUIT_        = iota
	CARD_SUIT_DIAMOND = iota + 1 // 黑桃
	CARD_SUIT_HEART              // 红桃
	CARD_SUIT_SPADE              // 樱花
	CARD_SUIT_CLUB               // 方片
	CARD_SUIT_JOKER              // 大小王无花色
)

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
