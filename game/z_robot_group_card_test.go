package game

import (
	"fmt"
	"github.com/wonderivan/logger"
	"testing"
)

func TestUserSy(t *testing.T) {
	hands := []*Card{

		{7, 1}, {7, 2},
		{8, 2}, {8, 1},
		{9, 3}, {9, 2}, {9, 1},
		{10, 3},
		{11, 3},
		{12, 2}, {12, 1},
	}

	res := continuouslyCountGroup(hands)
	logger.Debug(res[0].cardLen)
	logger.Debug(len(res[0].cardGroup))

	cards, remainCards1 := groupGroup(hands, res[0])
	for i := 0; i < len(cards); i++ {
		logger.Debug("类型", cards[i].CardType)
		PrintCard(cards[i].Card)
	}

	logger.Debug("remain")
	PrintCard(remainCards1)

}

func TestG(t *testing.T) {
	// cards := CreateBrokenCard()

	p1, _, _, _ := CreateCardsNew()

	logger.Debug("手牌")
	hands := p1
	SortCardSL(hands)
	PrintCard(hands)
	gc := GroupHandsCard(hands)

	value := CountCardValue(hands)
	logger.Debug("手牌分数：", value)

	logger.Debug("单张")
	rcs := gc.Single
	for i := 0; i < len(rcs); i++ {
		fmt.Print("weight:", rcs[i].Wight, "  ")
		PrintCard(rcs[i].Card)
	}

	logger.Debug("对子")
	rcd := gc.Double
	for i := 0; i < len(rcd); i++ {
		fmt.Print("weight:", rcd[i].Wight, "  ")
		PrintCard(rcd[i].Card)
	}

	logger.Debug("三张")
	rct := gc.Triple
	for i := 0; i < len(rct); i++ {
		fmt.Print("weight:", rct[i].Wight, "  ")
		PrintCard(rct[i].Card)
	}

	logger.Debug("炸弹")
	rcb := gc.Bomb
	for i := 0; i < len(rcb); i++ {
		fmt.Print("weight:", rcb[i].Wight, "  ")
		PrintCard(rcb[i].Card)
	}

	logger.Debug("火箭")
	rcr := gc.Rocket
	for i := 0; i < len(rcr); i++ {
		fmt.Print("weight:", rcr[i].Wight, "  ")
		PrintCard(rcr[i].Card)
	}

	logger.Debug("顺子")
	rcj := gc.Junko
	for i := 0; i < len(rcj); i++ {
		fmt.Print("weight:", rcj[i].Wight, "  ")
		PrintCard(rcj[i].Card)
	}

	logger.Debug("连对")
	rcjd := gc.JunkoDouble
	for i := 0; i < len(rcjd); i++ {
		fmt.Print("weight:", rcjd[i].Wight, "  ")
		PrintCard(rcjd[i].Card)
	}

	logger.Debug("飞机")
	rcjt := gc.junkTriple
	for i := 0; i < len(rcjt); i++ {
		fmt.Print("weight:", rcjt[i].Wight, "  ")
		PrintCard(rcjt[i].Card)
	}
}

/*
&{1 2},&{1 1},&{2 2},&{3 2},&{3 1},&{4 4},&{4 1},&{5 4},&{6 1},&{7 4},&{7 1},&{9 4},&{9 3},&{12 1},&{13 3},&{13 1},&{14 5},
*/

func TestAd(t *testing.T) {

	hands := []*Card{
		{5,2},{5,2},{5,2},
		{6,2},{6,2},{6,2},
		{7,2},{7,2},{7,2},
		{8,2},{8,2},{8,2},{8,2},
		{10,2},
		{12,2},{12,2},{12,2},
		{13,2},
		{14,2},
	}

	gc := GroupHandsCard(hands)
	logger.Debug("单张")
	rcs := gc.Single
	for i := 0; i < len(rcs); i++ {
		fmt.Print("weight:", rcs[i].Wight, "  ")
		PrintCard(rcs[i].Card)
	}

	logger.Debug("对子")
	rcd := gc.Double
	for i := 0; i < len(rcd); i++ {
		fmt.Print("weight:", rcd[i].Wight, "  ")
		PrintCard(rcd[i].Card)
	}

	logger.Debug("三张")
	rct := gc.Triple
	for i := 0; i < len(rct); i++ {
		fmt.Print("weight:", rct[i].Wight, "  ")
		PrintCard(rct[i].Card)
	}

	logger.Debug("炸弹")
	rcb := gc.Bomb
	for i := 0; i < len(rcb); i++ {
		fmt.Print("weight:", rcb[i].Wight, "  ")
		PrintCard(rcb[i].Card)
	}

	logger.Debug("火箭")
	rcr := gc.Rocket
	for i := 0; i < len(rcr); i++ {
		fmt.Print("weight:", rcr[i].Wight, "  ")
		PrintCard(rcr[i].Card)
	}

	logger.Debug("顺子")
	rcj := gc.Junko
	for i := 0; i < len(rcj); i++ {
		fmt.Print("weight:", rcj[i].Wight, "  ")
		PrintCard(rcj[i].Card)
	}

	logger.Debug("连对")
	rcjd := gc.JunkoDouble
	for i := 0; i < len(rcjd); i++ {
		fmt.Print("weight:", rcjd[i].Wight, "  ")
		PrintCard(rcjd[i].Card)
	}

	logger.Debug("飞机")
	rcjt := gc.junkTriple
	for i := 0; i < len(rcjt); i++ {
		fmt.Print("weight:", rcjt[i].Wight, "  ")
		PrintCard(rcjt[i].Card)
	}

}
