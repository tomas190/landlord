package game

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {

	cards := initOriginalCard()
	PrintCard(cards)
	OutOfCard(cards)
	PrintCard(cards)
	SortCard(cards)
	PrintCard(cards)

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
