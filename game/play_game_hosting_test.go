package game

import (
	"github.com/wonderivan/logger"
	"testing"
)

func TestHostingBeatTripleWithSingle(t *testing.T) {
	hands := []*Card{
		//{1, 1}, {1, 2}, {1, 3},// {1, 4},
		//{2, 1}, {2, 2}, {2, 3}, {2, 4},
		//{3, 1}, {3, 2}, //{3, 3}, {3, 4},
		//{4, 1}, {4, 2}, {4, 3},// {4, 4},
		{5, 1}, //{5, 2}, //{5, 3}, {5, 4},
		//{6, 1}, {6, 2}, {6, 3}, //{6, 4},
		{7, 1}, {7, 2}, {7, 3},// {7, 4},
		//{8, 1}, {8, 2}, {8, 3}, {8, 4},
		//{9, 1}, {9, 2}, {9, 3}, {9, 4},
		//{10, 1}, //{10, 2}, {10, 3}, {10, 4},
		//{11, 1}, {11, 2}, {11, 3}, {11, 4},
		//{12, 1}, {12, 2}, {12, 3},// {12, 4},
		//{13, 1}, {13, 2}, {13, 3}, //{13, 4},
		//{14, 1}, {15, 2},
	}
	outCard := []*Card{{13, 1}, {13, 2}, {13, 3},{1, 3}}

	// 测试三不带
	//cards, b, i := HostingBeatTriple(hands, outCard)
	cards, b, i := HostingBeatTripleWithSingle(hands, outCard)

	if b {
		logger.Debug("能打过")
		PrintCard(cards)
		logger.Debug("类型:",i)
		return
	}

	logger.Debug("打不过")

}

func TestGetAgent(t *testing.T) {
//	outCard := []*Card{{13, 1}, {13, 2}, {13, 3},{1, 3}}
	outCard2 := []*Card{{13, 1}, {13, 2}, {13, 3},{9, 3}}

	one := getCardManyOne(outCard2)

	logger.Debug(one)
}