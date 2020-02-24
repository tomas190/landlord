package game

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {

	for i:=0;i<10 ;i++  {
		card := CreateBrokenCard()

		PrintCard(card)
	}

}

func TestCheckCardsIsExist(t *testing.T) {
	hans := []*Card{{2, 2}, {5, 2}, {8, 2}, {9, 2}}
	card := []*Card{{2, 2}, {5, 2}, {6, 2},}

	fmt.Println(checkCardsIsExist(hans, card))
}

func TestReqEnterRoom(t *testing.T) {
	hands := []*Card{{2, 2}, {5, 2}, {8, 2}, {9, 2}}
	PrintCard(hands)
	out := []*Card{ {9, 2},{2, 2}, {5, 2}}
	result := removeCards(hands, out)
	PrintCard(hands)
	PrintCard(out)
	PrintCard(result)

}

func TestCreateGoodCard(t *testing.T)  {

	var base []int
	for i:=1;i<=13 ;i++  {
		base=append(base, i)
	}

	OutOfSort(base,0)

	dSort := RandNum(35, 48)
	a1, a2, a3, bo := CreateGoodCard(dSort)

	SortCardSL(a1)
	SortCardSL(a2)
	SortCardSL(a3)
	SortCardSL(bo)

	fmt.Println("===========")
	PrintCard(a1)
	PrintCard(a2)
	PrintCard(a3)
	PrintCard(bo)

}