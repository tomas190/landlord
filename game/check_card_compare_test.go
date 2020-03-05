package game

import (
	"fmt"
	"github.com/wonderivan/logger"
	"landlord/mconst/cardConst"
	"testing"
)

func TestHostingBeatSingle(t *testing.T) {

	//eCard := []*Card{{2, 2}}
	eCard := []*Card{
		{3, 1}, {3, 2}, //{4, 3},
		//{8, 2}, {8, 3},
		//{7, 2}, {7, 3},
		//{6, 2}, {6, 3},
		//{5, 3}, {5, 1}, //{5, 3},
		{13, 1}, {13, 2}, {13, 1}, {13, 2},
	}
	//eCard := []*Card{{3, 1}, {3, 2},{3, 3},{4, 4},{4, 4}}

	hands := []*Card{
		//{3, 4},
		//{4, 1},
		//	{5, 1}, {5, 2},{5, 2},
		//	{6, 1}, {6, 2}, {6, 3},
		//{7, 1}, {7, 2},
		{8, 2}, {8, 1}, {8, 3}, {8, 4},
		//{9, 1}, {9, 2}, {9, 3},
		{10, 1}, {10, 2},// {10, 3}, {10, 4},

		//{11, 1}, {11, 2}, {11, 3},{11, 4},
		//{14, 1}, {15, 2},
	}

	//cards, b ,ctype:= HostingBeatBombWithDouble(hands, eCard)
	gCard := GroupHandsCard(hands)
	cards, b ,ctype:= FindMinFollowCards(hands,gCard, eCard,cardConst.CARD_PATTERN_QUADPLEX_WITH_SINGLES)
	//b:= CanBeat(eCard, hands)
	//b:= CanBeat(hands,eCard)

	if b {
		PrintCard(cards)
		fmt.Println("能打過",ctype)
		return
	}
	fmt.Println("打不过")
}

func TestHostingBeatDouble(t *testing.T) {

	//eCard := []*Card{{2, 2}}
	//eCard := []*Card{{9 ,1},{9 ,2},{9, 3},{10, 1},{10, 2},}
	//eCard := []*Card{{3, 1}, {3, 2},{3, 3},{4, 4},{4, 4}}
	hands := []*Card{

		{7, 1}, {8, 1}, {8, 1}, {8, 1}, {8, 2}, {11, 1}, {11, 2}, {12, 1}, {12, 1}, {12, 1}, {12, 2}, {13, 2}, {13, 2}, {14, 2}, {15, 2},
	}

	cards, b, _ := FindCanBeatCards(hands, []*Card{{7, 3}, {7, 1}, {7, 2}, {5, 1}, {5, 1}}, cardConst.CARD_PATTERN_TRIPLET_WITH_PAIR)
	if b {
		PrintCard(cards)
		return
	}
	fmt.Println("打不过")
}

func TestOne(t *testing.T) {
	hands := []*Card{
		{6, 1}, {6, 2}, {6, 3}, {3, 4},
		{4, 1}, {5, 1},
		{7, 1},
		{8, 2}, {8, 1},
		{10, 1}, {10, 2}, {10, 3},
		{9, 1}, {9, 2}, {9, 3},
		{11, 1}, {11, 2}, {11, 3}, {11, 4},
		{14, 1}, {15, 2},
	}

	logger.Debug(junkoHelpRemove(hands))

	SortCardSL(hands)
	PrintCard(hands)

}

func TestCanBeat(t *testing.T) {
	eCards := []*Card{{3, cardConst.CARD_SUIT_SPADE}, {3, cardConst.CARD_SUIT_CLUB}}
	cards := []*Card{{7, 2}, {7, 1}}
	logger.Debug(CanBeat(eCards, cards))
}
