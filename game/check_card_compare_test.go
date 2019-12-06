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
		{3, 1}, {3, 2}, {3, 1}, {3, 2},
		{4, 1}, {4, 2}, //{4, 3},
		//{5, 3}, {5, 1}, //{5, 3},
		//{6, 2}, {6, 3},
		//{7, 2}, {7, 3},
		//{8, 2}, {8, 3},
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
		//{10, 1}, {10, 2}, {10, 3}, {10, 4},

		//{11, 1}, {11, 2}, {11, 3},{11, 4},
		{14, 1}, {15, 2},
	}

	//cards, b := HostingBeatBombWithDouble(hands, eCard)
	cards, b, rType := HostingBeatBombWithSingles(hands, eCard)

	if b {
		PrintCard(cards)
		logger.Debug("返回类型", rType)
		return
	}
	fmt.Println("打不过")
}

func TestHostingBeatDouble(t *testing.T) {

	//eCard := []*Card{{2, 2}}
	//eCard := []*Card{{9 ,1},{9 ,2},{9, 3},{10, 1},{10, 2},}
	//eCard := []*Card{{3, 1}, {3, 2},{3, 3},{4, 4},{4, 4}}
	hands := []*Card{
		//{6, 1}, {6, 2},{6, 3},{3, 4},
		//{4, 1}, {5, 1},
		//{7, 1},
		//{8, 2}, {8, 1},
		{10, 1}, {10, 2}, {10, 3}, {10, 3},
		//	{9, 1}, {9, 2}, {9, 3},
		//{11, 1}, {11, 2}, {11, 3},{11, 4},
		//{14, 1},
		{15, 4},
	}

	cards, b, _ := FindCanBeatCards(hands, []*Card{{14, 5,}}, cardConst.CARD_PATTERN_SINGLE)
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
