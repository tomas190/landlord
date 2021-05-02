package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/cardConst"
	"landlord/mconst/sysSet"
)

// 根据盈余池获取机器人的胜率
func GetRobotWinPercentage(surplusPool float64) float64 {
	logger.Debug("当前盈余池:", surplusPool)
	if surplusPool > 500 { // 盈余池为正 取盈余池为正的机器人胜率 为确保盈利 为正条件先盈利500
		logger.Debug("盈余池为正的机器人胜率:", sysSet.RANDOM_PERCENTAGE_AFTER_WIN)
		return sysSet.RANDOM_PERCENTAGE_AFTER_WIN
	} else { // 盈余池为负 取盈余池为负的机器人胜率
		logger.Debug("盈余池为负的机器人胜率:", sysSet.RANDOM_PERCENTAGE_AFTER_WIN)
		return sysSet.RANDOM_PERCENTAGE_AFTER_LOSE
	}

}

// 根据机器人的胜率 确定输赢
// @return true : 机器人赢(机器人必赢好牌)
// @return false : 机器人输(随机发牌)
func RobotWinOrLoseByRobotWinPercentage(robotWinPercentage float64) bool {
	robotWinPercentageNum := robotWinPercentage * 100
	num := RandNum(1, 100)
	if float64(num) > robotWinPercentageNum {
		logger.Debug("结果 机器人必赢策虐")
		return true
	} else {
		logger.Debug("结果 机器输,随机发牌")
		return false
	}

}

// 获取好牌  第一首是好牌
func MustWinCard() ([]*Card, []*Card, []*Card, []*Card) {

	winCardType := RandNum(1, 5)

	switch winCardType {
	case 1:
		return MustWinCardType1()
	case 2:
		return MustWinCardType2()
	case 3:
		return MustWinCardType3()
	case 4:
		return MustWinCardType4()
	case 5:
		return MustWinCardType4()
	default:
		return MustWinCardType1()
	}

}

/*  =============   好牌类型  =============*/

// 必赢好牌类型 1
// @return
// 1 :好牌 2:随机牌 3:随机牌 4:底牌
// 一炸 加2连飞机
func MustWinCardType1() ([]*Card, []*Card, []*Card, []*Card) {

	rCard := CreateBrokenCard()

	var winCard []*Card

	threeC, remainCard := GetThreeBigCard(rCard)

	bomb, remainCard := GetBombInCards(remainCard)

	airplane2, remainCard := Get2AirplaneInCards(remainCard)

	winCard = append(winCard, threeC...)    // 3
	winCard = append(winCard, bomb...)      // 4
	winCard = append(winCard, airplane2...) // 6 13

	// 补四张牌
	buCard := remainCard[:4]
	winCard = append(winCard, buCard...)

	p1 := append([]*Card{}, remainCard[4:21]...)
	p2 := append([]*Card{}, remainCard[21:38]...)

	bottom := append([]*Card{}, remainCard[38:]...)
	return winCard, p1, p2, bottom

}

// 必赢好牌类型 2
// @return
// 1 :好牌 2:随机牌 3:随机牌 4:底牌
// 王炸 加一炸  加2连飞机
func MustWinCardType2() ([]*Card, []*Card, []*Card, []*Card) {

	rCard := CreateBrokenCard()

	var winCard []*Card

	rocket, remainCard := GetRocketInCards(rCard)

	bomb, remainCard := GetBombInCards(remainCard)

	airplane2, remainCard := Get2AirplaneInCards(remainCard)

	winCard = append(winCard, rocket...)    // 2
	winCard = append(winCard, bomb...)      // 4
	winCard = append(winCard, airplane2...) // 6 12

	// 补6张牌
	buCard := remainCard[:5]
	winCard = append(winCard, buCard...)

	p1 := append([]*Card{}, remainCard[5:22]...)
	p2 := append([]*Card{}, remainCard[22:39]...)

	bottom := append([]*Card{}, remainCard[39:]...)
	return winCard, p1, p2, bottom

}

// 必赢好牌类型 3
// @return
// 1 :好牌 2:随机牌 3:随机牌 4:底牌
// 三连飞机
func MustWinCardType3() ([]*Card, []*Card, []*Card, []*Card) {

	rCard := CreateBrokenCard()

	var winCard []*Card

	threeC, remainCard := GetThreeBigCard(rCard)

	//bomb, remainCard := GetBombInCards(remainCard)

	airplane3, remainCard := Get3AirplaneInCards(remainCard)

	winCard = append(winCard, threeC...) // 3
	//winCard = append(winCard, bomb...)      //
	winCard = append(winCard, airplane3...) // 9 12

	// 补6张牌
	buCard := remainCard[:5]
	winCard = append(winCard, buCard...)

	p1 := append([]*Card{}, remainCard[5:22]...)
	p2 := append([]*Card{}, remainCard[22:39]...)

	bottom := append([]*Card{}, remainCard[39:]...)
	return winCard, p1, p2, bottom

}

// 必赢好牌类型 4
// @return
// 1 :好牌 2:随机牌 3:随机牌 4:底牌
// 两炸 加1连飞机
func MustWinCardType4() ([]*Card, []*Card, []*Card, []*Card) {

	rCard := CreateBrokenCard()

	var winCard []*Card

	threeC, remainCard := GetThreeBigCard(rCard)

	bomb, remainCard := GetBombInCards(remainCard)
	bomb2, remainCard := GetBombInCards(remainCard)
	threeSame, remainCard := Get1AirplaneInCards(remainCard)

	winCard = append(winCard, threeC...) // 3
	//winCard = append(winCard, bomb...)      //
	winCard = append(winCard, bomb...) // 4 12
	winCard = append(winCard, bomb2...) // 4 11
	winCard = append(winCard, threeSame...) // 3 14

	// 补3张牌
	buCard := remainCard[:3]
	winCard = append(winCard, buCard...)

	p1 := append([]*Card{}, remainCard[3:20]...)
	p2 := append([]*Card{}, remainCard[20:37]...)

	bottom := append([]*Card{}, remainCard[37:]...)
	return winCard, p1, p2, bottom

}


// 必赢好牌类型 5
// @return
// 1 :好牌 2:随机牌 3:随机牌 4:底牌
// 两炸
func MustWinCardType5() ([]*Card, []*Card, []*Card, []*Card) {

	rCard := CreateBrokenCard()

	var winCard []*Card

	threeC, remainCard := GetThreeBigCard(rCard)

	bomb, remainCard := GetBombInCards(remainCard)
	bomb2, remainCard := GetBombInCards(remainCard)

	winCard = append(winCard, threeC...) // 3
	//winCard = append(winCard, bomb...)      //
	winCard = append(winCard, bomb...) // 4 12
	winCard = append(winCard, bomb2...) // 4 11

	// 补5张牌
	buCard := remainCard[:6]
	winCard = append(winCard, buCard...)

	p1 := append([]*Card{}, remainCard[6:23]...)
	p2 := append([]*Card{}, remainCard[23:40]...)

	bottom := append([]*Card{}, remainCard[40:]...)
	return winCard, p1, p2, bottom

}
/*  =============   好牌类型  =============*/

/*========= help func ===========*/

// 从牌中获取王炸
func GetRocketInCards(c []*Card) ([]*Card, []*Card) {
	rocket, b, _ := hasRacket(c)
	var reCards []*Card
	if b {
		remainCards := removeCards(c, rocket)
		var reCard Card
		reCards = append(reCards, &reCard)
		return rocket, remainCards
	}
	return nil, c
}

// 从牌中获取一个随机炸弹
func GetBombInCards(c []*Card) ([]*Card, []*Card) {
	boom := getHasNumsCard(c, 4)
	var result []*Card
	if len(boom) > 0 {
		randNum := RandNum(0, len(boom)-1)

		for i := 0; i < len(c); i++ {
			if int(c[i].Value) == boom[randNum] {
				result = append(result, c[i])
			}
			if len(result) == 4 {
				break
			}
		}
		remainCards := removeCards(c, result)
		return result, remainCards
	}
	return nil, c
}

// 从牌中获取一个随机1连飞机 tip
func Get1AirplaneInCards(c []*Card) ([]*Card, []*Card) {
	item := getHasMoreNumsCard(c, 3)
	var result []*Card
	if len(item) >= 2 {
		randNum := RandNum(0, len(item)-2)

		for i := 0; i < len(c); i++ {
			if int(c[i].Value) == item[randNum] {
				result = append(result, c[i])
			}
			if len(result) == 3 {
				break
			}
		}


		remainCards := removeCards(c, result)
		return result, remainCards
	}
	return nil, c
}

// 从牌中获取一个随机2连飞机 tip 三个2 三个A 也可
func Get2AirplaneInCards(c []*Card) ([]*Card, []*Card) {
	item := getHasMoreNumsCard(c, 3)
	var result []*Card
	if len(item) >= 2 {
		randNum := RandNum(0, len(item)-2)

		for i := 0; i < len(c); i++ {
			if int(c[i].Value) == item[randNum] {
				result = append(result, c[i])
			}
			if len(result) == 3 {
				break
			}
		}

		for i := 0; i < len(c); i++ {
			if int(c[i].Value) == item[randNum+1] {
				result = append(result, c[i])
			}
			if len(result) == 6 {
				break
			}
		}

		remainCards := removeCards(c, result)
		return result, remainCards
	}
	return nil, c
}

// 从牌中获取一个随机3连飞机 tip 三个2 三个A 也可
func Get3AirplaneInCards(c []*Card) ([]*Card, []*Card) {
	item := getHasMoreNumsCard(c, 3)
	var result []*Card
	if len(item) >= 3 {
		randNum := RandNum(0, len(item)-3)

		for i := 0; i < len(c); i++ {
			if int(c[i].Value) == item[randNum] {
				result = append(result, c[i])
			}
			if len(result) == 3 {
				break
			}
		}

		for i := 0; i < len(c); i++ {
			if int(c[i].Value) == item[randNum+1] {
				result = append(result, c[i])
			}
			if len(result) == 6 {
				break
			}
		}

		for i := 0; i < len(c); i++ {
			if int(c[i].Value) == item[randNum+2] {
				result = append(result, c[i])
			}
			if len(result) == 9 {
				break
			}
		}

		remainCards := removeCards(c, result)
		return result, remainCards
	}
	return nil, c
}

// 从牌中获取一个随机4连飞机 tip 三个2 三个A 也可
func Get4AirplaneInCards(c []*Card) ([]*Card, []*Card) {
	item := getHasMoreNumsCard(c, 3)
	var result []*Card
	if len(item) >= 4 {
		randNum := RandNum(0, len(item)-4)

		for i := 0; i < len(c); i++ {
			if int(c[i].Value) == item[randNum] {
				result = append(result, c[i])
			}
			if len(result) == 3 {
				break
			}
		}

		for i := 0; i < len(c); i++ {
			if int(c[i].Value) == item[randNum+1] {
				result = append(result, c[i])
			}
			if len(result) == 6 {
				break
			}
		}

		for i := 0; i < len(c); i++ {
			if int(c[i].Value) == item[randNum+2] {
				result = append(result, c[i])
			}
			if len(result) == 9 {
				break
			}
		}

		for i := 0; i < len(c); i++ {
			if int(c[i].Value) == item[randNum+3] {
				result = append(result, c[i])
			}
			if len(result) == 12 {
				break
			}
		}

		remainCards := removeCards(c, result)
		return result, remainCards
	}
	return nil, c
}

// 随机取大鬼 小鬼 或者2 共计三张
func GetThreeBigCard(c []*Card) ([]*Card, []*Card) {
	var result []*Card
	for i := 0; i < len(c); i++ {
		if c[i].Value == cardConst.CARD_RANK_RED_JOKER ||
			c[i].Value == cardConst.CARD_RANK_BLACK_JOKER ||
			c[i].Value == cardConst.CARD_RANK_TWO {
			result = append(result, c[i])
		}
		if len(result) >= 3 {
			break
		}
	}

	remainCards := removeCards(c, result)
	return result, remainCards

}
