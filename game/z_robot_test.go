package game

import (
	"fmt"
	"github.com/wonderivan/logger"
	"landlord/mconst/roomType"
	"testing"
)

func TestCreateGroupCard(t *testing.T) {
	hands := []*Card{
		{3, 4},
		{4, 1},
		{5, 1},
		{6, 1}, {6, 2}, {6, 3},
		{7, 1},
		{8, 2}, {8, 1},
		{9, 1}, {9, 2}, {9, 3},
		{10, 1}, {10, 2}, {10, 3},
		{11, 1}, {11, 2}, {11, 3}, {11, 4},
		//{14, 1},
		{15, 2},
	}

	gc := CreateGroupCard(hands)
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
	rcb, _ := FindAllBomb(hands)
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

	//logger.Debug("单张")
	//rcs, _ := FindAllSingle(hands)
	//for i := 0; i < len(rcs); i++ {
	//	PrintCard(rcs[i].Card)
	//}
	//
	//logger.Debug("对子")
	//rcd, _ := FindAllDouble(hands)
	//for i := 0; i < len(rcd); i++ {
	//	PrintCard(rcd[i].Card)
	//}
	//
	//logger.Debug("三张")
	//rct, _ := FindAllTriplet(hands)
	//for i := 0; i < len(rct); i++ {
	//	PrintCard(rct[i].Card)
	//}
	//
	//logger.Debug("炸弹")
	//rcb, remainCards := FindAllBomb(hands)
	//for i := 0; i < len(rcb); i++ {
	//	PrintCard(rcb[i].Card)
	//}
	//
	//logger.Debug("剩余:")
	//PrintCard(remainCards)

	//logger.Debug("火箭")
	//rcr, _ := FindRocket(hands)
	//for i := 0; i < len(rcr); i++ {
	//	PrintCard(rcr[i].Card)
	//}

}

func TestJunko(t *testing.T) {
	hands := []*Card{
		//{1, 1}, {1, 2}, //{1, 3},// {1, 4},
		//{2, 1}, {2, 2}, //{2, 3}, {2, 4},
		//{3, 1}, {3, 2}, //{3, 3}, {3, 4},
		{4, 1}, {4, 2}, //{4, 3}, // {4, 4},
		{5, 1}, {5, 2}, //{5, 3}, {5, 4},

		{6, 1}, {6, 2}, //{6, 3}, //{6, 4},
		{7, 1}, {7, 2}, //{7, 3}, //{7, 4},
		{8, 1}, {8, 2}, //{8, 3}, {8, 4},
		{9, 1}, {9, 2}, // {9, 3}, //{9, 4},
		//{10, 1}, {10, 2}, //{10, 3}, {10, 4},

		{11, 1}, {11, 2},          //{11, 3}, {11, 4},
		{12, 1}, {12, 2}, {12, 3}, // {12, 4},
		{13, 1}, {13, 2}, {13, 3}, //{13, 4},
		{14, 1}, {15, 2},
	}

	single := []*Card{
		//{1, 1}, {1, 2}, {1, 3},// {1, 4},
		//{2, 1}, {2, 2}, {2, 3}, //{2, 4},
		{3, 1},                 //{3, 2}, //{3, 3}, {3, 4},
		{4, 1},                 //{4, 2}, //{4, 3}, // {4, 4},
		{5, 1}, {5, 2}, {5, 3}, //{5, 4},

		{6, 1},  //{6, 2},   //{6, 3}, //{6, 4},
		{7, 1},  //{7, 2},   //{7, 3}, //{7, 4},
		{8, 1},  //{8, 2},   //{8, 3}, {8, 4},
		{9, 1},  //{9, 2},   // {9, 3}, //{9, 4},
		{10, 1}, //{10, 2}, //{10, 3}, {10, 4},

		{11, 1}, {11, 2}, //{11, 3}, {11, 4},
		{12, 1},          // {12, 2}, {12, 3}, // {12, 4},
		{13, 1},          //{13, 2}, {13, 3}, //{13, 4},
		{14, 1}, {15, 2},
	}

	//rc, remainCards := FindPossibleLongSingleJunko(hands)
	rc, remainCards := unlimitedJunko(hands)
	logger.Debug("junkos or======================")
	for i := 0; i < len(rc); i++ {
		fmt.Print("weight:", rc[i].Wight, "  ")
		PrintCard(rc[i].Card)
	}
	logger.Debug("remain or======================")
	PrintCard(remainCards)

	//junko := mergeJunkoWithJunko(rc)
	//logger.Debug("merge or======================")
	//for i := 0; i < len(junko); i++ {
	//	fmt.Print("weight:", junko[i].Wight, "  ")
	//	PrintCard(junko[i].Card)
	//}

	cards, _ := FindAllSingle(single)
	logger.Debug("merge or single======================")
	logger.Debug("单牌：")
	for i := 0; i < len(cards); i++ {
		fmt.Print("weight:", cards[i].Wight, "  ")
		PrintCard(cards[i].Card)
	}

	logger.Debug("rcrcrcrcr======================")
	for i := 0; i < len(rc); i++ {
		fmt.Print("weight:", rc[i].Wight, "  ")
		PrintCard(rc[i].Card)
	}

}

func TestPriorityJunko(t *testing.T) {
	hands := []*Card{
		//{1, 1},// {1, 2},         //{1, 3},// {1, 4},
		{2, 1},                 //{2, 2}, //{2, 3}, {2, 4},
		{3, 1},                 //{3, 2}, //{3, 3}, {3, 4},
		{4, 1},                 //{4, 2}, //{4, 3}, // {4, 4},
		{5, 1}, {5, 2}, {5, 3}, // {5, 4},

		{6, 1}, {6, 2}, {6, 3}, // {6, 4},
		{7, 1}, {7, 2}, {7, 3}, //{7, 4},
		{8, 1},                 //{8, 2},   //{8, 3}, {8, 4},
		{9, 1},                 //{9, 2},   // {9, 3}, //{9, 4},
		{10, 1},                //{10, 2}, //{10, 3}, {10, 4},

		{11, 1}, // {11, 2}, // {11, 3}, {11, 4},
		//{12, 1}, //{12, 2}, //{12, 3}, //{12, 4},
		//{13, 1}, {13, 2}, {13, 3}, //{13, 4},
		//{14, 1}, {15, 2},
	}

	gc := PriorityJunko(hands)
	//gc := PriorityJunkoProtectTriple(hands)
	logger.Debug("单张")
	rcs := gc.Single
	for i := 0; i < len(rcs); i++ {
		fmt.Print("weight:", rcs[i].Wight, "  ")
		PrintCard(rcs[i].Card)
	}

	logger.Debug("顺子")
	rcj := gc.Junko
	for i := 0; i < len(rcj); i++ {
		fmt.Print("weight:", rcj[i].Wight, "  ")
		PrintCard(rcj[i].Card)
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
	rcb, _ := FindAllBomb(hands)
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

}

func TestCr(t *testing.T) {
	hands := []*Card{
		//{1, 1},// {1, 2},         //{1, 3},// {1, 4},
		{2, 1},                 //{2, 2}, //{2, 3}, {2, 4},
		{3, 1},                 //{3, 2}, //{3, 3}, {3, 4},
		{4, 1},                 //{4, 2}, //{4, 3}, // {4, 4},
		{5, 1}, {5, 2}, {5, 3}, // {5, 4},

		{6, 1}, {6, 2}, {6, 3}, // {6, 4},
		{7, 1}, {7, 2}, {7, 3}, //{7, 4},
		{8, 1},                 //{8, 2},   //{8, 3}, {8, 4},
		{9, 1},                 //{9, 2},   // {9, 3}, //{9, 4},
		{10, 1},                //{10, 2}, //{10, 3}, {10, 4},

		{11, 1}, // {11, 2}, // {11, 3}, {11, 4},
		//{12, 1}, //{12, 2}, //{12, 3}, //{12, 4},
		//{13, 1}, {13, 2}, {13, 3}, //{13, 4},
		//{14, 1}, {15, 2},
	}
	//
	//c := CreateBrokenCard()
	//hands=c[:17]
	remainCards := hands

	logger.Debug("hands")
	SortCardSL(hands)
	PrintCard(hands)
	seqDouble, junkos, remainCard := FindAllJunko(remainCards)

	gc := CreateGroupCard(remainCard)
	gc.Junko = junkos
	gc.JunkoDouble = seqDouble
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
	rcb, _ := FindAllBomb(hands)
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
}

func TestCreateRobot(t *testing.T) {
	cards, i, i2, i3 := CreateGodCards()

	SortCard(cards)
	SortCard(i)
	SortCard(i2)
	SortCard(i3)

	fmt.Println("天牌")
	PrintCard(cards)
	fmt.Println("爛牌1")
	PrintCard(i)
	fmt.Println("爛牌2")
	PrintCard(i2)
	fmt.Println("底牌")
	PrintCard(i3)

	a := append([]*Card{}, cards...)
	a = append(a, i...)
	a = append(a, i2...)
	a = append(a, i3...)

	fmt.Println("縂:", len(a))
	SortCardSL(a)

	for i := 0; i < len(a); i++ {
		if i%4 == 0 {
			fmt.Println()
		}
		PrintCard(append([]*Card{},a[i]))
	}

}

func TestAddHighFieldWaitUser(t *testing.T) {
	fmt.Println(getRobotGold(roomType.ExperienceField))
	fmt.Println(getRobotGold(roomType.LowField))
	fmt.Println(getRobotGold(roomType.MidField))
	fmt.Println(getRobotGold(roomType.HighField))
}