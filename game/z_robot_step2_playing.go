package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/cardConst"
	"landlord/mconst/playerStatus"
)

// 机器人打牌阶段操作
func RobotPlayAction(room *Room, robot, nextPlayer, lastPlayer *Player) {
	// 机器人打牌了
	isFakerDisconnection := delayDestiny()
	logger.Debug("机器人打牌阶段...")
	if isFakerDisconnection { // 如果概率出现了 假装掉线 则配合不出操作 并且机器人以后走托管流程
		logger.Debug("机器人打牌阶段 中0.001%概率掉线托管...")
		robot.IsGameHosting = true
		RespGameHosting(room, playerStatus.GameHosting, robot.PlayerPosition, robot.PlayerInfo.PlayerId)
		if robot.IsMustDo {
			DoGameHosting(room, robot, nextPlayer, lastPlayer)
		} else {
			NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
		}
		return
	}
	//DoGameHosting(room, robot, nextPlayer, lastPlayer)
	robotOutCard(room, robot, nextPlayer, lastPlayer)
}

// 机器人出牌
func robotOutCard(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {
	if robot.IsLandlord {
		if robot.IsMustDo {
			landlordRobotOutCardMustDo(room, robot, nextPlayer, lastPlayer)
		} else {
			// todo
			DoGameHosting(room, robot, nextPlayer, lastPlayer)
		}
	} else {
		// todo
		DoGameHosting(room, robot, nextPlayer, lastPlayer)
	}

}

// 地主机器人首出牌策虐
/*
	0.判断能否一首出完
	1.检测 对方玩家是否 报单或者报双 或者最后一手牌
	2.在判断是否有顺子
	3.先判断是否有连对
	4.在判断是否有飞机
*/
func landlordRobotOutCardMustDo(room *Room, robot *Player, nextPlayer *Player, lastPlayer *Player) {

	//1. 先判断能否一首出完 	// 检测牌型正确代表能一首出完
	cardType := GetCardsType(robot.HandCards)
	if !(cardType > cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS || cardType < cardConst.CARD_PATTERN_SINGLE) {
		OutCardsAction(room, robot, nextPlayer, robot.HandCards, cardType)
		return
	}

	group := robot.GroupCard
	var outCard []*Card
	if len(group.Junko) >= 1 { // 如果有顺子
		// todo 这里要先去最小的
		cardType = group.Junko[0].CardType
		outCard = group.Junko[0].Card
	} else if len(group.JunkoDouble) >= 1 { // 如果有连对
		// todo 这里要先去最小的
		cardType = group.JunkoDouble[0].CardType
		outCard = group.JunkoDouble[0].Card
	} else if len(group.junkTriple) >= 1 { // 飞机
		// todo 这里应该是找出  单张和对子中最小的牌
		outCard = group.junkTriple[0].Card
		tLen := len(outCard) / 3
		withCards := findTripleWithCards(group.Single, group.Double, group.JunkoDouble, tLen)
		if withCards != nil {
			outCard = append(outCard, withCards...)
		}
		cardType = group.junkTriple[0].CardType
	} else if len(group.Triple) >= 1 {
		outCard = group.Triple[0].Card
		tLen := len(outCard) / 3
		if len(group.Single) >= tLen {
			var ss []*Card
			for i := 0; i < tLen; i++ {
				ss = append(ss, group.Single[i].Card...)
			}
			outCard = append(outCard, ss...)
		} else if len(group.Double) >= tLen {
			var dd []*Card
			for i := 0; i < tLen; i++ {
				dd = append(dd, group.Double[i].Card...)
			}
			outCard = append(outCard, dd...)
		}
		cardType = group.Triple[0].CardType
	} else { // 则看是出单张还是 连对
		//todo  1.是否有人报单 或者双
		if len(group.Double) > len(group.Single) && len(group.Double) >= 1 {
			outCard = append(outCard, group.Double[0].Card...)
		} else if len(group.Double) < len(group.Single) && len(group.Single) >= 1 {
			outCard = append(outCard, group.Single[0].Card...)
		} else {
			// fuck
			outCard, cardType = FindMustBeOutCards(robot.HandCards)
		}
	}
	// newHands := removeCards(robot.HandCards, outCard)

	OutCardsAction(room, robot, nextPlayer, outCard, cardType)

}

/*
 	为3张寻找带牌
	从已有牌中找出最小的带牌
	从单张和对子中寻找带牌数量

	1.首先判断单牌的张数  和 对子的数量 如果单张数量已经满足 则优先用单张里面抽取
	大部分从单张带起


	如果张数一样 先带单张 如果单张中 包含大王 并且对子数量为0  则允许 反之带最小对子

	// 如果单张和对子的数量都为0  并且刚好有连对长度与之吻合 带上
*/
func findTripleWithCards(singles, double, junkoDouble []*ReCard, length int) []*Card {
	var withCards []*Card
	// 1.优先寻找单排
	if len(singles) >= length {
		for i := 0; i < length; i++ {
			withCards = append(withCards, singles[i].Card...)
		}
		return withCards
	}

	// todo 如果此时 单张刚好满足 且包含 2 或者王这样的打牌 并且 对子也满足且很小的情况下 则还是带对子比较好
	// 2.寻找对子
	if len(double) >= length {
		for i := 0; i < length; i++ {
			withCards = append(withCards, double[i].Card...)
		}
		return withCards
	}

	// 3.把单牌和对子混合 取最小张  如果总长度大于 的话
	var mix []*Card
	for i := 0; i < len(singles); i++ {
		mix = append(mix, singles[i].Card...)
	}
	for i := 0; i < len(double); i++ {
		mix = append(mix, double[i].Card...)
	}

	if len(mix) >= length {
		SortCardSL(mix)
		for i := 0; i < length; i++ {
			withCards = append(withCards, mix[i])
		}

		return withCards
	}

	// 如果这都还没有 连对的吻合 则带连对
	if len(junkoDouble) == length {
		var final []*Card
		for i := 0; i < length; i++ {
			final = append(final, junkoDouble[i].Card...)
		}

		return final
	}
	return nil
}
