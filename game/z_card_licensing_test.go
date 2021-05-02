package game

import (
	"fmt"
	"github.com/wonderivan/logger"
	"strconv"
	"testing"
)

func TestGetGoodCards(t *testing.T) {
	//test push code

	cs := CreateBrokenCard()

	cards, remainCs := GetGoodCards(cs)

	PrintCard(cards)

	SortCard(cards)
	SortCard(remainCs)

	PrintCard(cards)
	PrintCard(remainCs)

}

func TestRobotWinOrLoseByRobotWinPercentage(t *testing.T) {
	var robotWinNum, robotLoseNum int

	for j := 1; j < 10; j++ {
		for i := 1; i <= 1000; i++ {
			if RobotWinOrLoseByRobotWinPercentage(0.5) {
				robotWinNum++
			} else {
				robotLoseNum++
			}
		}
		fmt.Println("100局 机器人赢的局数:", robotWinNum, "/n机器人输的局数:", robotLoseNum)
		robotWinNum = 0
		robotLoseNum = 0
	}

}

func TestGetAndRemoveCard(t *testing.T) {

	//PrintCard(rCard)
	for i := 0; i < 100; i++ {
		rCard := CreateBrokenCard()
		result, remainCards := GetThreeBigCard(rCard)

		logger.Debug(len(result))
		logger.Debug(len(remainCards))

		PrintCard(result)
		SortCard(remainCards)
		PrintCard(remainCards)

		logger.Debug("===============================")
	}
}

func TestMustWinCardType1(t *testing.T) {
	goodCard, p1, p2, bottom := MustWinCardType4()

	fmt.Println(len(goodCard))
	fmt.Println(len(p1))
	fmt.Println(len(p2))
	fmt.Println(len(bottom))

	SortCard(goodCard)
	SortCard(p1)
	SortCard(p2)
	SortCard(bottom)

	PrintCard1(goodCard)
	PrintCard1(p1)
	PrintCard1(p2)
	PrintCard1(bottom)

	var rC []*Card
	rC = append(rC, goodCard...)
	rC = append(rC, p1...)
	rC = append(rC, p2...)
	rC = append(rC, bottom...)

	fmt.Println(len(rC))
	SortCard(rC)
	PrintCard1(rC)

}

func PrintCard1(cards []*Card) {
	for i := 0; i < len(cards); i++ {
		suit := printSuit(cards[i].Suit)
		value := printValue(cards[i].Value)
		fmt.Print(suit,value," ")
	}
	fmt.Println()
}

func printSuit(i int32) interface{} {
	switch i {
	case 1:
		return "黑桃"
	case 2:
		return "黑桃"
	case 3:
		return "梅花"
	case 4:
		return "方片"
	default:
		return ""
	}

}

func printValue(value int32) string {

	if value == 15 {
		return "大王"
	} else if value == 14 {
		return "小王"
	} else if value == 13 {
		return "2"
	} else if value == 12 {
		return "A"
	} else if value == 11 {
		return "K"
	} else if value == 10 {
		return "Q"
	} else if value == 9 {
		return "J"
	} else {
		return strconv.Itoa(int(value) + 2)
	}

}
