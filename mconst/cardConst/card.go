package cardConst

// 定义扑克级别
const (
	Card_RANK_THREE = iota + 1
	CARD_RANK_FOUR
	CARD_RANK_FIVE
	CARD_RANK_SIX
	CARD_RANK_SEVEN
	CARD_RANK_EIGHT
	CARD_RANK_NINE
	CARD_RANK_TEN
	CARD_RANK_JACK
	CARD_RANK_QUEEN
	CARD_RANK_KING
	CARD_RANK_ACE
	CARD_RANK_TWO
	CARD_RANK_BLACK_JOKER
	CARD_RANK_RED_JOKER
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
	CARD_SUIT_SPADE              // 黑桃
	CARD_SUIT_CLUB               // 梅花
	CARD_SUIT_JOKER              // 大小王无花色
)
