package game

import (
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
