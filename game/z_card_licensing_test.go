package game

import (
	"testing"
)

func TestGetGoodCards(t *testing.T) {

	cs := CreateBrokenCard()

	cards, remainCs := GetGoodCards(cs)

	PrintCard(cards)

	SortCard(cards)
	SortCard(remainCs)

	PrintCard(cards)
	PrintCard(remainCs)

}