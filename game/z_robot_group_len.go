package game

import "github.com/wonderivan/logger"

func groupLen2(hands []*Card, g group) ([]*ReCard, []*Card) {
	// 连续长度为二 只能组成飞机

	if g.cardLen != 6 {
		return nil, hands
	}
	var r []*ReCard

	seqArr, isContinuously, _ := howManyCardByX(g, 3)
	if isContinuously {
		re, remainCards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
		if re!=nil {
			r = append(r, re)
		}
		return r, remainCards
	}
	return nil, hands
}

func groupLen3(hands []*Card, g group) ([]*ReCard, []*Card) {
	// 连续长度为二 只能组成飞机
	cardLen := g.cardLen // 牌的总张数
	if cardLen < 6 {
		return nil, hands
	}
	// 值需要判断三连 二连  和连对
	// 1. 判断3连
	var r []*ReCard
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if isContinuously && num >= 2 {
		re, remainCards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
		if re!=nil {
			r = append(r, re)
		}
		return r, remainCards
	}

	seqArr, isContinuously, num = howManyCardByX(g, 2)
	if isContinuously && num >= 3 {
		re, remainCards := groupJunkoDouble(hands, seqArr[0], seqArr[len(seqArr)-1])
		if re!=nil {
			r = append(r, re)
		}
		return r, remainCards
	}
	return nil, hands
}

//
//
func groupLen4(hands []*Card, g group) ([]*ReCard, []*Card) {

	// 连续长度为二 只能组成飞机
	cardLen := g.cardLen // 牌的总张数
	if cardLen < 6 {
		return nil, hands
	}

	// 值需要判断三连 二连  和连对
	// 1. 判断4连
	var r []*ReCard
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if isContinuously && num >= 2 {
		re, remainCards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
		if re!=nil {
			r = append(r, re)
		}
		return r, remainCards
	}

	// 2. 判断3连
	seqNum, has := hasContinuouslyLonger(seqArr, 3)
	if has {
		re, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
		if re!=nil {
			r = append(r, re)
		}
		return r, remainCards
	}

	// 2. 判断2连
	seqNum, has = hasContinuouslyLonger(seqArr, 2)
	if has {
		re, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
		if re!=nil {
			r = append(r, re)
		}
		return r, remainCards
	}

	seqArr, isContinuously, num = howManyCardByX(g, 2)
	if isContinuously && num >= 3 {
		re, remainCards := groupJunkoDouble(hands, seqArr[0], seqArr[len(seqArr)-1])
		if re!=nil {
			r = append(r, re)
		}
		return r, remainCards
	}
	return nil, hands
}

// 牌去重张数为5的情况下
func groupLen5(hands []*Card, g group) ([]*ReCard, []*Card) {
	var r []*ReCard
	var re *ReCard
	var remainCard []*Card
	cardLen := g.cardLen // 牌的总张数
	if cardLen <= 7 { // 如果总长度小于7 则组成顺子
		re, remainCard = groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[4].cardValue))
	} else if cardLen == 8 { // 有可能连对 如果不能连对 则还是顺子
		re, remainCard = groupLen5Has8(hands, g)
	} else if cardLen == 9 { // 如果9 就不要组成顺子了  不划算
		re, remainCard = groupLen5Has9(hands, g)
	} else if cardLen == 10 { // 如果10 就不要组成顺子了  不划算
		re, remainCard = groupLen5Has10(hands, g)
	} else if cardLen == 11 { // 如果10 就不要组成顺子了  不划算
		re, remainCard = groupLen5Has11(hands, g)
	} else if cardLen == 12 {
		return groupLen5Has12(hands, g)
	} else if cardLen == 13 {
		return groupLen5Has13(hands, g)
	} else if cardLen == 14 {
		return groupLen5Has14(hands, g)
	} else if cardLen == 15 {
		return groupLen5Has15(hands, g)
	} else {
		logger.Error("！！！ 异常的组牌信息")
	}

	if re == nil {
		return nil, hands
	}
	return append(r, re), remainCard
}

// len5的时候
/*============ 5 ==============*/
// finish
func groupLen5Has8(hands []*Card, g group) (*ReCard, []*Card) {
	seqArr, isContinuously, num := howManyCardByX(g, 2)
	if num < 3 || !isContinuously { // 如果重复为2的牌小于3 或者不连续 则不能组成连队
		return groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[4].cardValue))
	} else { // 组成连队
		return groupJunkoDouble(hands, seqArr[0], seqArr[len(seqArr)-1])
	}
}

// finish
func groupLen5Has9(hands []*Card, g group) (*ReCard, []*Card) {
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if isContinuously && num == 2 { // 连续且等于2  // 飞机无疑 不带牌
		return groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
	}

	seqArr, isContinuously, num = howManyCardByX(g, 2)
	if num == 4 {
		if isContinuously { // 重复为2 连续且等于4  // 连对无疑
			return groupJunkoDouble(hands, seqArr[0], seqArr[len(seqArr)-1])
		} else {
			seqNum, has := hasContinuouslyLonger(seqArr, 3)
			if has {
				return groupJunkoDouble(hands, seqNum[0], seqNum[len(seqNum)-1])
			}
		}
	}
	// 不组合牌型 不划算
	return nil, hands
}

// finish
func groupLen5Has10(hands []*Card, g group) (*ReCard, []*Card) {
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if isContinuously && num == 2 { // 连续且等于2  // 飞机无疑 不带牌
		return groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
	}

	seqArr, isContinuously, num = howManyCardByX(g, 2)
	// 10张牌能连续的
	if num == 5 { // 重复为2 连续且等于5  // 连对无疑
		if isContinuously {
			return groupJunkoDouble(hands, seqArr[0], seqArr[len(seqArr)-1])
		}
	} else if num == 3 {
		// 如果等于三的情况下 则判断加上有三张的牌 如果计算连续 则拆三张组合成连对
		MoreXSeqArr, isMoreXContinuously, MoreXNum := howManyCardMoreX(g, 2)
		if MoreXNum == 4 && isMoreXContinuously {
			return groupJunkoDouble(hands, MoreXSeqArr[0], MoreXSeqArr[len(MoreXSeqArr)-1])
		}
		if isContinuously {
			return groupJunkoDouble(hands, seqArr[0], seqArr[len(seqArr)-1])
		}

	}

	// 不组合牌型 不划算
	return nil, hands
}

// finish
func groupLen5Has11(hands []*Card, g group) (*ReCard, []*Card) {
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 3 { //
		if isContinuously {
			return groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
		}
		// 寻找二连飞机
		seqNum, has := hasContinuouslyLonger(seqArr, 2)
		if has {
			return groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
		}
	} else if num == 2 {
		if isContinuously {
			return groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
		}
	} else if num == 1 { // 从四连对开始找
		/*
			{2, 1}, {2, 2}, {2, 3},
			{3, 1}, {3, 2},
			{4, 1}, {4, 2},
			{5, 1},	{5, 3},
			{6, 1}, {3, 3},
		*/
		seqNum, has, num := howManyCardByX(g, 2) // bum必定等于4
		if has && num == 4 {
			return groupJunkoDouble(hands, seqNum[0], seqNum[len(seqNum)-1])
		}

		// 找三连对
		lSeqNum, has := hasContinuouslyLonger(seqNum, 3)
		if has {
			return groupJunkoDouble(hands, lSeqNum[0], lSeqNum[len(lSeqNum)-1])
		}

	}

	// 不组合牌型 不划算
	return nil, hands
}

// finish
func groupLen5Has12(hands []*Card, g group) ([]*ReCard, []*Card) {
	var r []*ReCard
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 3 { //
		if isContinuously {
			reCard, remainCards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			return append(r, reCard), remainCards
		}
		// 寻找二连飞机
		seqNum, has := hasContinuouslyLonger(seqArr, 2)
		if has {
			reCard, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			return append(r, reCard), remainCards
		}
	} else if num == 2 {
		/*
			{2, 1}, {2, 2}, {2, 3},
			{3, 1}, {3, 2},	{3,	2}
			{4, 1}, {4, 2},
			{5, 1},	{5, 3},
			{6, 1}, {3, 3},
		*/
		var re *ReCard
		remainCards := hands
		if isContinuously {
			re, remainCards = groupJunkoTriple(remainCards, seqArr[0], seqArr[len(seqArr)-1])
			r = append(r, re)
		}
		// 寻找三连对
		seqNum, has, _ := howManyCardByX(g, 2) // 必定是三
		if has {
			re, remainCards = groupJunkoDouble(remainCards, seqNum[0], seqNum[len(seqNum)-1])
			r = append(r, re)
		}

		return r, remainCards
	} // else if num == 1 { }// 不可能存在这种情况 }

	// 不组合牌型 不划算
	return nil, hands
}

// finish
func groupLen5Has13(hands []*Card, g group) ([]*ReCard, []*Card) {
	var r []*ReCard
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 4 { // 最大四连飞
		if isContinuously {
			reCard, remainCards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			return append(r, reCard), remainCards
		}
		// 寻找三连飞机
		seqNum, has := hasContinuouslyLonger(seqArr, 3)
		if has {
			reCard, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			return append(r, reCard), remainCards
		}

		// 寻找二连飞
		seqNum, has = hasContinuouslyLonger(seqArr, 2)
		if has {
			re, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			r = append(r, re)
			// 有可能还有一个双飞
			arrNum := removeArrNum(seqArr, seqNum)
			seqNum, has = hasContinuouslyLonger(arrNum, 2)
			if has { // 一个双连飞
				re, remainCards = groupJunkoTriple(remainCards, seqNum[0], seqNum[len(seqNum)-1])
				r = append(r, re)
				// 有可能还有一个双飞
			}
			return r, remainCards
		}
	} else if num == 3 {
		/*
			{1, 1}, {1, 2},
			{2, 1}, {2, 2}, {2, 3},
			{3, 1}, {3, 3},
			{4, 1}, {4, 2}, {4, 3},
			{5, 1}, {5, 2}, {5, 3},
		*/
		if isContinuously {
			reCard, remainCards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			return append(r, reCard), remainCards
		}

		// 寻找二连飞机
		seqNum, has := hasContinuouslyLonger(seqArr, 2)
		if has {
			reCard, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			return append(r, reCard), remainCards
		}

	} //else if num == 2 { // 不存在这种情况 }

	// 不组合牌型 不划算
	return nil, hands
}

// finish
func groupLen5Has14(hands []*Card, g group) ([]*ReCard, []*Card) {
	var r []*ReCard
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 4 { // 最大四连飞
		if isContinuously {
			reCard, remainCards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			return append(r, reCard), remainCards
		}
		// 寻找三连飞机
		seqNum, has := hasContinuouslyLonger(seqArr, 3)
		if has {
			reCard, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			return append(r, reCard), remainCards
		}

		// 寻找二连飞
		seqNum, has = hasContinuouslyLonger(seqArr, 2)
		if has {
			re, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			r = append(r, re)
			// 有可能还有一个双飞
			arrNum := removeArrNum(seqArr, seqNum)
			seqNum, has = hasContinuouslyLonger(arrNum, 2)
			if has { // 一个双连飞
				re, remainCards = groupJunkoTriple(remainCards, seqNum[0], seqNum[len(seqNum)-1])
				r = append(r, re)
				// 有可能还有一个双飞
			}
			return r, remainCards
		}
	} //else if num == 3 { // 没有这种可能
	//	/*
	//		{1, 1}, {1, 2},
	//		{2, 1}, {2, 2}, {2, 3},
	//		{3, 1}, {3, 3},
	//		{4, 1}, {4, 2}, {4, 3},
	//		{5, 1}, {5, 2}, {5, 3},
	//	*/
	//	if isContinuously {
	//		reCard, remainCards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
	//		return append(r, reCard), remainCards
	//	}
	//
	//	// 寻找二连飞机
	//	seqNum, has := hasContinuouslyLonger(seqArr, 2)
	//	if has {
	//		reCard, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
	//		return append(r, reCard), remainCards
	//	}
	//
	//} //else if num == 2 { // 不存在这种情况 }

	// 不组合牌型 不划算
	return nil, hands
}

// finish
func groupLen5Has15(hands []*Card, g group) ([]*ReCard, []*Card) {
	var r []*ReCard
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 5 { // 最大5连飞
		if isContinuously {
			reCard, remainCards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			return append(r, reCard), remainCards
		}
	}
	logger.Error(" !!!incredible")
	return nil, hands
}

/*============ 5 ==============*/

// 牌去重张数为6的情况下
// finish
func groupLen6(hands []*Card, g group) ([]*ReCard, []*Card) {

	var r []*ReCard
	cardLen := g.cardLen // 牌的总张数
	if cardLen <= 9 { // 如果总长度小于9 则组成顺子
		cards, remainCards := groupLen6less9(hands, g)
		if cards!=nil {
			r = append(r, cards)
		}
		//r = append(r, cards)
		return r, remainCards
	} else if cardLen == 10 { // 如果总长度等于10
		cards, remainCards := groupLen6Has10(hands, g)
		if cards!=nil {
			r = append(r, cards)
		}
		//r = append(r, cards)
		return r, remainCards
	} else if cardLen == 11 { // 如果总长度等于10
		cards, remainCards := groupLen6Has11(hands, g)
		if cards!=nil {
			r = append(r, cards)
		}
		return r, remainCards
	} else if cardLen == 12 {
		cards, remainCards := groupLen6Has12(hands, g)
		if cards!=nil {
			r = append(r, cards)
		}
		//r = append(r, cards)
		return r, remainCards
	} else if cardLen == 13 {
		return groupLen6Has13(hands, g)
	} else if cardLen == 14 {
		return groupLen6Has14(hands, g)
	} else if cardLen == 15 {
		return groupLen6Has15(hands, g)
	} else if cardLen == 16 {
		return groupLen6Has16(hands, g)
	} else if cardLen == 17 {
		return groupLen6Has17(hands, g)
	} else if cardLen == 18 {
		return groupLen6Has18(hands, g)
	}
	return nil, hands
}

/*============ 6 ==============*/

func groupLen6less9(hands []*Card, g group) (*ReCard, []*Card) {
	seqNum, _, num := howManyCardByX(g, 3) // 有则num必定为1
	if num == 1 { // 如果有三张的情况 要判断是否在两边 是则移除
		if seqNum[0] == int(g.cardGroup[0].cardValue) {
			cards, remainCards := groupJunko(hands, int(g.cardGroup[1].cardValue), int(g.cardGroup[5].cardValue))
			return cards, remainCards
		} else if seqNum[0] == int(g.cardGroup[5].cardValue) {
			cards, remainCards := groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[4].cardValue))
			return cards, remainCards
		}
	}
	cards, remainCards := groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))
	return cards, remainCards
}

// finish
func groupLen6Has10(hands []*Card, g group) (*ReCard, []*Card) {
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 2 { // 两个三张 能连续 则飞机 不连续 则顺子
		if isContinuously {
			return groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
		} else { // 如果不连续则判断是否这个三张在两边 优先保住大的三带
			if seqArr[1] == int(g.cardGroup[5].cardValue) {
				return groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[4].cardValue))
			} else if seqArr[0] == int(g.cardGroup[0].cardValue) {
				return groupJunko(hands, int(g.cardGroup[1].cardValue), int(g.cardGroup[5].cardValue))
			}
			return groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))
		}
	} else if num == 1 {
		if seqArr[0] == int(g.cardGroup[5].cardValue) {
			return groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[4].cardValue))
		} else if seqArr[0] == int(g.cardGroup[0].cardValue) {
			return groupJunko(hands, int(g.cardGroup[1].cardValue), int(g.cardGroup[5].cardValue))
		}
		// 组成顺子
		return groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))

	}

	//如能组成连对则组成 反之不拆 todo 其实拆也可以 单少住不拆
	seqArr, isContinuously, num = howManyCardByX(g, 2)
	if num >= 3 {
		// todo num == 4  先判断双顺
		if isContinuously {
			return groupJunkoDouble(hands, seqArr[0], seqArr[len(seqArr)-1])
		}
		seqNum, has := hasContinuouslyLonger(seqArr, 3)
		if has {
			return groupJunkoDouble(hands, seqNum[0], seqNum[len(seqNum)-1])
		}
	}
	return nil, hands
}

// finish
func groupLen6Has11(hands []*Card, g group) (*ReCard, []*Card) {
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 2 { // 两个三张 能连续 则飞机 不连续 则最优顺子
		if isContinuously {
			return groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
		} else { // 如果不连续则判断是否这个三张在两边 优先保住大的三带
			if seqArr[1] == int(g.cardGroup[5].cardValue) {
				return groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[4].cardValue))
			} else if seqArr[0] == int(g.cardGroup[0].cardValue) {
				return groupJunko(hands, int(g.cardGroup[1].cardValue), int(g.cardGroup[5].cardValue))
			} else {
				/*  todo 这种情况 这里感觉不拆比较好
				{3, 1},
				{4, 1},
				{5, 1}, {5, 2}, {5, 3},
				{6, 1}, {6, 2},
				{7, 1}, {7, 2}, {7, 3},
				{8, 1},
				*/
				return nil, hands
				//return groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))
			}
		}
		// 这里很多中情况  后面感觉越来越难写了 炸
	} else if num == 1 { // 1个三张
		if seqArr[0] == int(g.cardGroup[5].cardValue) { // 如果三张在最大的那边
			newG := g
			newG.cardLen = 5
			newG.cardGroup = newG.cardGroup[:5]
			return groupLen5Has8(hands, newG)
		} else if seqArr[0] == int(g.cardGroup[0].cardValue) { // 如果三张在最小的那边
			newG := g
			newG.cardLen = 5
			newG.cardGroup = newG.cardGroup[1:]
			return groupLen5Has8(hands, newG)
		} else { // 及三张在中间的情况
			/*  todo 及这种情况  能组成3连对就组合 不拆三张
			{3, 1},					 {3, 1},{3, 1},{3, 1},		{3, 1},						{3, 1}, {3, 1},
			{4, 1},					 {4, 1},					{4, 1},	{4, 1},				{4, 1},
			{5, 1}, {5, 2}, {5, 3},	 {5, 1},					{5, 1},	{5, 1},{5, 1},		{5, 1},
			{6, 1}, {6, 2},			 {6, 1}, {6, 2},			{6, 1},	{6, 1},				{6, 1}, {6, 2},	{6,2}
			{7, 1}, {7, 2},			 {7, 1}, {7, 2},			{7, 1}, {7, 2},				{7, 1}, {7, 2},
			{8, 1}, {8, 1},			 {8, 1}, {8, 1},			{8, 1}, 					{8, 1}, {8, 1},
																todo 这种情况就不组成连对了
			*/
			seqArr, isContinuously, num = howManyCardByX(g, 2)
			if num == 3 && isContinuously {
				return groupJunkoDouble(hands, seqArr[0], seqArr[len(seqArr)-1])
			}

		}
		return nil, hands
	} else { // 没有三张的情况
		/*
				注意这种情况
			{3, 1},
			{4, 1},	{4,	1}
			{5, 1}, {5, 2},
			{6, 1}, {6, 2},
			{7, 1}, {7, 2},
			{8, 1}, {8, 1},
			todo fix 按理说这种情况应该拆成两个顺子 这里直接组成了连对
		*/
		seqArr, isContinuously, num = howManyCardByX(g, 2)
		if num >= 3 {
			if isContinuously {
				return groupJunkoDouble(hands, seqArr[0], seqArr[len(seqArr)-1])
			}
			seqNum, has := hasContinuouslyLonger(seqArr, 4)
			if has {
				return groupJunkoDouble(hands, seqNum[0], seqNum[len(seqNum)-1])
			}
			seqNum, has = hasContinuouslyLonger(seqArr, 3)
			if has {
				return groupJunkoDouble(hands, seqNum[0], seqNum[len(seqNum)-1])
			}
		}
	}
	return nil, hands
}

// finish
func groupLen6Has12(hands []*Card, g group) (*ReCard, []*Card) {
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 3 { // 三个三张 能连续 则飞机
		if isContinuously { // 三连飞
			return groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
		} else { // 如果不连续
			seqNum, has := hasContinuouslyLonger(seqArr, 2)
			if has { // 双连飞
				return groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			}
			return nil, hands
		}
		// 这里很多中情况
	} else if num == 2 { // 2个三张 不用组
		if isContinuously { // 三连飞
			return groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
		}
	} else if num == 1 { // 1个三张
		// 此时num 必定为4
		seqArr, isContinuously, _ := howManyCardByX(g, 2)
		if isContinuously { // 四连对
			return groupJunkoDouble(hands, seqArr[0], seqArr[len(seqArr)-1])
		}

		seqArr, isContinuously, _ = howManyCardMoreX(g, 2)
		if isContinuously { // 拆三代一5连对 之后五连对
			return groupJunkoDouble(hands, seqArr[0], seqArr[len(seqArr)-1])
		}

		seqNum, has := hasContinuouslyLonger(seqArr, 3)
		if has { // 三连对

			return groupJunkoDouble(hands, seqNum[0], seqNum[len(seqNum)-1])
		}
	} else { // 这里不用说  必定5连对
		return groupJunkoDouble(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))
	}
	return nil, hands
}

// finish
func groupLen6Has13(hands []*Card, g group, ) ([]*ReCard, []*Card) {
	var r []*ReCard
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 3 { // 三个三张 能连续 则飞机
		if isContinuously { // 三连飞
			re, cards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			return append(r, re), cards

		} else { // 如果不连续
			seqNum, has := hasContinuouslyLonger(seqArr, 2)
			if has { // 双连飞
				re, cards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
				return append(r, re), cards
			}
			return nil, hands
		}
		// 这里很多中情况
	} else if num == 2 { // 2个三张 不用组
		if isContinuously { // 双连飞 todo 还要组成连对
			/*
				{1, 1},
				{2, 1}, {2, 2}, {2, 3},
				{3, 1}, {3, 2}, {3, 3},
				{4, 1}, {4, 2},
				{5, 1}, {5, 2},
				{6, 1}, {6, 2},
			*/
			reJt, remainCards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			r = append(r, reJt)
			// 双联组成了 判断是否还存在连对
			seqArr, isContinuously, num := howManyCardByX(g, 2)
			if isContinuously && num == 3 {
				reJt, remainCards = groupJunkoDouble(remainCards, seqArr[0], seqArr[len(seqArr)-1])
				r = append(r, reJt)
			}
			return r, remainCards
		}
	} else if num == 1 { // 1个三张
		// 必定是这种情况
		/*
			{1, 1}, {1, 2}				{1, 1}, {1, 2}				{1, 1}, {1, 2}
			{2, 1}, {2, 2},				{2, 1}, {2, 2},	{2,3}		{2, 1}, {2, 2},
			{3, 1}, {3, 2}, {3, 3},		{3, 1}, {3, 2},				{3, 1}, {3, 2},
			{4, 1}, {4, 2},				{4, 1}, {4, 2},				{4, 1}, {4, 2},
			{5, 1}, {5, 2},				{5, 1}, {5, 2},				{5, 1}, {5, 2},
			{6, 1}, {6, 2},				{6, 1}, {6, 2},				{6, 1}, {6, 2},{6，3}
		*/
		// 不拆三带
		seqArr, isContinuously, num := howManyCardByX(g, 2)
		if num == 5 { // 有可能连续
			if isContinuously {
				reJt, remainCards := groupJunkoDouble(hands, seqArr[0], seqArr[len(seqArr)-1])
				r = append(r, reJt)
				return r, remainCards
			} else {
				seqNum, has := hasContinuouslyLonger(seqArr, 4)
				if has {
					reJt, remainCards := groupJunkoDouble(hands, seqNum[0], seqNum[len(seqNum)-1])
					r = append(r, reJt)
					return r, remainCards
				}

				seqNum, has = hasContinuouslyLonger(seqArr, 3)
				if has {
					reJt, remainCards := groupJunkoDouble(hands, seqNum[0], seqNum[len(seqNum)-1])
					r = append(r, reJt)
					return r, remainCards
				}
			}
		}
	}
	return nil, hands
}

// finish
func groupLen6Has14(hands []*Card, g group, ) ([]*ReCard, []*Card) {
	var r []*ReCard
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 4 { // 4个三张 能连续 则飞机
		if isContinuously { // 4连飞
			re, cards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			return append(r, re), cards
		} else { // 如果不连续 // 有可能一个三连飞  有可能 两个双连飞 至少存在一个双连飞
			seqNum, has := hasContinuouslyLonger(seqArr, 3)
			if has { // 好一个三连飞
				re, cards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
				return append(r, re), cards
			}
			seqNum, has = hasContinuouslyLonger(seqArr, 2)
			if has { // 好一个双连飞
				re, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
				r = append(r, re)
				// 有可能还有一个双飞
				arrNum := removeArrNum(seqArr, seqNum)
				seqNum, has = hasContinuouslyLonger(arrNum, 2)
				if has { // 一个双连飞
					re, remainCards = groupJunkoTriple(remainCards, seqNum[0], seqNum[len(seqNum)-1])
					r = append(r, re)
					// 有可能还有一个双飞
				}
				return r, remainCards
			}
		}
		// 这里很多中情况
	} else if num == 3 { // 3个三张 不用组
		/*
			{1, 1},
			{2, 1},	{2, 2}
			{3, 1}, {3, 2},
			{4, 1}, {4, 2},	{4,	2},
			{5, 1}, {5, 2}, {5, 3},
			{6, 1}, {6, 2}, {6, 3},
		*/
		if isContinuously { // 3连飞
			reJt, remainCards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			r = append(r, reJt)
			return r, remainCards
		}
		// 如果不连续 找连续的双飞
		seqNum, has := hasContinuouslyLonger(seqArr, 2)
		if has { // 好一个三连飞
			re, cards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			return append(r, re), cards
		}

		// 双飞都不连  直接返回不拆了
		return nil, hands
	} else if num == 2 { // 2个三张
		/*
			{1, 1}, {1, 2},
			{2, 1}, {2, 2},
			{3, 1}, {3, 2},
			{4, 1}, {4, 2},
			{5, 1}, {5, 2}, {5, 3},
			{6, 1}, {6, 2}, {6, 3},
		*/
		if isContinuously { // 2连飞
			reJt, remainCards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			r = append(r, reJt)
			// 连续了之后判断是否还有连对
			seqNum, has, _ := howManyCardByX(g, 2) // 这个地方不用说 num 为4
			if has { // 四连对
				reJt, remainCards = groupJunkoDouble(remainCards, seqNum[0], seqNum[len(seqNum)-1])
				r = append(r, reJt)
			} else { // 四连对没有 找三连对
				seqNum, has := hasContinuouslyLonger(seqNum, 3)
				if has { // 四连对
					reJt, remainCards = groupJunkoDouble(remainCards, seqNum[0], seqNum[len(seqNum)-1])
					r = append(r, reJt)
				}
			}

			return r, remainCards
		} else { // 双飞不连续
			seqNum, has, _ := howManyCardByX(g, 2) // 这个地方不用说 num 为4
			if has { // 四连对
				reJt, remainCards := groupJunkoDouble(hands, seqNum[0], seqNum[len(seqNum)-1])
				r = append(r, reJt)
				return r, remainCards
			} else { // 四连对没有 找三连对
				seqNum, has := hasContinuouslyLonger(seqNum, 3)
				if has { //
					reJt, remainCards := groupJunkoDouble(hands, seqNum[0], seqNum[len(seqNum)-1])
					r = append(r, reJt)
					return r, remainCards
				}
			}

		}
	}
	// 至此不可能存在num=1 的情况
	return nil, hands
}

// finish
func groupLen6Has15(hands []*Card, g group, ) ([]*ReCard, []*Card) {
	var r []*ReCard
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 4 { // 4个三张 能连续 则飞机
		if isContinuously { // 4连飞
			re, cards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			return append(r, re), cards
		} else { // 如果不连续 // 有可能一个三连飞  有可能 两个双连飞 至少存在一个双连飞
			seqNum, has := hasContinuouslyLonger(seqArr, 3)
			if has { // 好一个三连飞
				re, cards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
				return append(r, re), cards
			}
			seqNum, has = hasContinuouslyLonger(seqArr, 2)
			if has { // 好一个双连飞
				re, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
				r = append(r, re)
				// 有可能还有一个双飞
				arrNum := removeArrNum(seqArr, seqNum)
				seqNum, has = hasContinuouslyLonger(arrNum, 2)
				if has { // 一个双连飞
					re, remainCards = groupJunkoTriple(remainCards, seqNum[0], seqNum[len(seqNum)-1])
					r = append(r, re)
					// 有可能还有一个双飞
				}
				return r, remainCards
			}
		}
		// 这里很多中情况
	} else if num == 3 { // 3个三张 不用组
		var reJt *ReCard
		remainCards := hands
		/*
			{1, 1}, {1,	2}
			{2, 1},	{2, 2}
			{3, 1}, {3, 2},
			{4, 1}, {4, 2},	{4,	2},
			{5, 1}, {5, 2}, {5, 3},
			{6, 1}, {6, 2}, {6, 3},
		*/
		if isContinuously { // 3连飞
			reJt, remainCards = groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			r = append(r, reJt)

			//三连对
			seqNum, has, _ := howManyCardByX(g, 2) // 这里num 必定3
			if has {
				reJt, remainCards = groupJunkoDouble(remainCards, seqNum[0], seqNum[len(seqNum)-1])
				r = append(r, reJt)
			}
			return r, remainCards
		}
		// 如果不连续 找连续的双飞
		seqNum, has := hasContinuouslyLonger(seqArr, 2)
		if has { // 好一个三连飞
			reJt, remainCards = groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			r = append(r, reJt)
		}
		// 三连对
		seqNum, has, _ = howManyCardByX(g, 2) // 这里num 必定3
		if has {
			reJt, remainCards = groupJunkoDouble(remainCards, seqNum[0], seqNum[len(seqNum)-1])
			r = append(r, reJt)
		}
		return r, remainCards

		// 双飞都不连  直接返回不拆了
	} else if num == 2 { // 不可能存在两个三张的情况
	}
	// 至此不可能存在num=1 的情况
	return nil, hands
}

// finish
func groupLen6Has16(hands []*Card, g group, ) ([]*ReCard, []*Card) {
	var r []*ReCard
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 5 {
		if isContinuously {
			re, cards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			return append(r, re), cards
		}
		// 寻找4连飞
		seqNum, has := hasContinuouslyLonger(seqArr, 4)
		if has {
			re, cards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			return append(r, re), cards
		}

		//寻找三连
		seqNum, has = hasContinuouslyLonger(seqArr, 3)
		if has { // 可能还存在着二连
			re, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			r = append(r, re)

			arrNum := removeArrNum(seqArr, seqNum)
			seqNum, has = hasContinuouslyLonger(arrNum, 2)
			if has {
				re, remainCards = groupJunkoTriple(remainCards, seqNum[0], seqNum[len(seqNum)-1])
				r = append(r, re)
			}
			return r, remainCards
		}

	} else if num == 4 { // 4个三张 能连续 则飞机
		if isContinuously { // 4连飞
			re, cards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			return append(r, re), cards
		} else { // 如果不连续 // 有可能一个三连飞  有可能 两个双连飞 至少存在一个双连飞
			seqNum, has := hasContinuouslyLonger(seqArr, 3)
			if has { // 好一个三连飞
				re, cards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
				return append(r, re), cards
			}
			seqNum, has = hasContinuouslyLonger(seqArr, 2)
			if has { // 好一个双连飞
				re, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
				r = append(r, re)
				// 有可能还有一个双飞
				arrNum := removeArrNum(seqArr, seqNum)
				seqNum, has = hasContinuouslyLonger(arrNum, 2)
				if has { // 一个双连飞
					re, remainCards = groupJunkoTriple(remainCards, seqNum[0], seqNum[len(seqNum)-1])
					r = append(r, re)
					// 有可能还有一个双飞
				}
				return r, remainCards
			}
		}
		// 这里很多中情况
	} //else if num == 3 { // 3个三张 不用组 不逊在等于3一下的情况

	//} else if num == 2 { // 不可能存在两个三张的情况
	//}
	// 至此不可能存在num=1 的情况
	return nil, hands
}

// finish
func groupLen6Has17(hands []*Card, g group, ) ([]*ReCard, []*Card) {
	var r []*ReCard
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 5 {
		if isContinuously {
			re, cards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			return append(r, re), cards
		}
		// 寻找4连飞
		seqNum, has := hasContinuouslyLonger(seqArr, 4)
		if has {
			re, cards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			return append(r, re), cards
		}

		//寻找三连
		seqNum, has = hasContinuouslyLonger(seqArr, 3)
		if has { // 可能还存在着二连
			re, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			r = append(r, re)

			arrNum := removeArrNum(seqArr, seqNum)
			lSeqNum, has := hasContinuouslyLonger(arrNum, 2)
			if has {
				//	logger.Debug(lSeqNum)
				re, remainCards = groupJunkoTriple(remainCards, lSeqNum[0], lSeqNum[len(lSeqNum)-1])
				r = append(r, re)
			}
			return r, remainCards
		}
	} // 不可能存在 5一下的 3三张
	return nil, hands
}

// finish
func groupLen6Has18(hands []*Card, g group, ) ([]*ReCard, []*Card) {
	var r []*ReCard
	seqArr, isContinuously, num := howManyCardByX(g, 3)
	if num == 6 {
		if isContinuously {
			re, cards := groupJunkoTriple(hands, seqArr[0], seqArr[len(seqArr)-1])
			return append(r, re), cards
		}
	} // 不可能存在 6以下的 3三张
	return nil, hands
}

/*============ 6 ==============*/

/*============ 7 ==============*/
// 牌去重张数为6的情况下
// finish
func groupLen7(hands []*Card, g group) ([]*ReCard, []*Card) {
	var r []*ReCard
	cardLen := g.cardLen // 牌的总张数
	if cardLen <= 8 { // 如果总长度小于9 则组成顺子
		cards, remainCards := groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))
		r = append(r, cards)
		return r, remainCards
	} else if cardLen == 9 {
		cards, remainCards := groupLen7Has9(hands, g)
		r = append(r, cards)
		return r, remainCards
	} else if cardLen == 10 { // 如果总长度等于10
		return groupLen7Has10(hands, g)
	} else if cardLen == 11 {
		return groupLen7Has11(hands, g)
	} else if cardLen == 12 {
		return groupLen7Has12(hands, g)
	} else if cardLen == 13 {
		return groupLen6Has13(hands, g)
	} else if cardLen == 14 {
		return groupLen6Has14(hands, g)
	} else if cardLen == 15 {
		return groupLen6Has15(hands, g)
	} else if cardLen == 16 {
		return groupLen6Has16(hands, g)
	} else if cardLen == 17 {
		return groupLen6Has17(hands, g)
	} else if cardLen == 18 {
		return groupLen6Has18(hands, g)
	}
	return nil, hands
}

/*============ 7 ==============*/
func groupLen7Has9(hands []*Card, g group) (*ReCard, []*Card) {
	seqNum, _, num := howManyCardByX(g, 3) // 有则num必定为1
	if num == 1 { // 如果有三张的情况 要判断是否在两边 是则移除
		if seqNum[0] == int(g.cardGroup[0].cardValue) {
			cards, remainCards := groupJunko(hands, int(g.cardGroup[1].cardValue), int(g.cardGroup[6].cardValue))

			return cards, remainCards
		} else if seqNum[0] == int(g.cardGroup[6].cardValue) {
			cards, remainCards := groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))
			return cards, remainCards
		}
	}
	cards, remainCards := groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[6].cardValue))
	return cards, remainCards
}

func groupLen7Has10(hands []*Card, g group) ([]*ReCard, []*Card) {
	/* // 这种情况可以组成两个顺子
	{1, 1},				{1, 1},
	{2, 1},				{2, 1},
	{3, 1}, {3, 2}		{3, 1},
	{4, 1}, {4, 2},		{4, 1}, {4, 2},
	{5, 1}, {5, 2},		{5, 1},
	{6, 1},				{6, 1},
	{7, 1},				{7, 1},{7, 2},{7, 2},
	*/

	// 1.还是判断是否有三张在两边
	seqNum, _, num := howManyCardByX(g, 3) // 有则num必定为1
	if num == 1 { // 如果有三张的情况 要判断是否在两边 是则移除
		if seqNum[0] == int(g.cardGroup[0].cardValue) {
			cards, remainCards := groupJunko(hands, int(g.cardGroup[1].cardValue), int(g.cardGroup[6].cardValue))
			return append([]*ReCard{}, cards), remainCards
		} else if seqNum[0] == int(g.cardGroup[6].cardValue) {
			cards, remainCards := groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))
			return append([]*ReCard{}, cards), remainCards
		}
	}

	// 然后在找出所有顺子
	// 2. 是否能找出两个顺子
	cards, remainCards := unlimitedJunko(hands)
	reCards, remainCards := mergeJunkoWithSingle(cards, remainCards)
	return reCards, remainCards
}

// finish
func groupLen7Has11(hands []*Card, g group) ([]*ReCard, []*Card) {
	// 1.还是判断是否有三张在两边
	seqNum, isContinuously, num := howManyCardByX(g, 3) // 有则num必定为1
	if num == 1 { // 如果有三张的情况 要判断是否在两边 是则移除
		if seqNum[0] == int(g.cardGroup[0].cardValue) {
			cards, remainCards := groupJunko(hands, int(g.cardGroup[1].cardValue), int(g.cardGroup[6].cardValue))
			return append([]*ReCard{}, cards), remainCards
		} else if seqNum[0] == int(g.cardGroup[6].cardValue) {
			cards, remainCards := groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))
			return append([]*ReCard{}, cards), remainCards
		}
		//cards, remainCards := groupJunko(hands, seqNum[0], seqNum[len(seqNum)])
		//return append([]*ReCard{}, cards), remainCards
	} else if num == 2 {
		if isContinuously {
			triple, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			junko, remainCards := unlimitedJunko(remainCards)
			return append(junko, triple), remainCards
		}

		if seqNum[0] == int(g.cardGroup[0].cardValue) {
			if seqNum[0] == int(g.cardGroup[0].cardValue) {
				cards, remainCards := groupJunko(hands, int(g.cardGroup[1].cardValue), int(g.cardGroup[6].cardValue))
				return append([]*ReCard{}, cards), remainCards
			} else if seqNum[0] == int(g.cardGroup[6].cardValue) {
				cards, remainCards := groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))
				return append([]*ReCard{}, cards), remainCards
			}

			if seqNum[1] == int(g.cardGroup[0].cardValue) {
				cards, remainCards := groupJunko(hands, int(g.cardGroup[1].cardValue), int(g.cardGroup[6].cardValue))
				return append([]*ReCard{}, cards), remainCards
			} else if seqNum[1] == int(g.cardGroup[6].cardValue) {
				cards, remainCards := groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))
				return append([]*ReCard{}, cards), remainCards
			}
			cards, remainCards := groupJunko(hands, seqNum[0], seqNum[len(seqNum)])
			return append([]*ReCard{}, cards), remainCards
		}

	}

	// 然后在找出所有顺子
	// 2. 是否能找出两个顺子
	cards, remainCards := unlimitedJunko(hands)
	if len(cards) == 2 {
		reCards, remainCards := mergeJunkoWithSingle(cards, remainCards)
		return reCards, remainCards
	} else {
		cards, remainCards := FindPossibleLongSingleJunko(hands)
		return append([]*ReCard{}, cards), remainCards
	}
}


// finish
func groupLen7Has12(hands []*Card, g group) ([]*ReCard, []*Card) {
	// 1.还是判断是否有三张在两边
	seqNum, isContinuously, num := howManyCardByX(g, 3) // 有则num必定为1
	if num == 1 { // 如果有三张的情况 要判断是否在两边 是则移除
		if seqNum[0] == int(g.cardGroup[0].cardValue) {
			cards, remainCards := groupJunko(hands, int(g.cardGroup[1].cardValue), int(g.cardGroup[6].cardValue))
			return append([]*ReCard{}, cards), remainCards
		} else if seqNum[0] == int(g.cardGroup[6].cardValue) {
			cards, remainCards := groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))
			return append([]*ReCard{}, cards), remainCards
		}
		//cards, remainCards := groupJunko(hands, seqNum[0], seqNum[len(seqNum)])
		//return append([]*ReCard{}, cards), remainCards
	} else if num == 2 {
		if isContinuously {
			triple, remainCards := groupJunkoTriple(hands, seqNum[0], seqNum[len(seqNum)-1])
			junko, remainCards := unlimitedJunko(remainCards)
			return append(junko, triple), remainCards
		}

		if seqNum[0] == int(g.cardGroup[0].cardValue) {
			if seqNum[0] == int(g.cardGroup[0].cardValue) {
				cards, remainCards := groupJunko(hands, int(g.cardGroup[1].cardValue), int(g.cardGroup[6].cardValue))
				return append([]*ReCard{}, cards), remainCards
			} else if seqNum[0] == int(g.cardGroup[6].cardValue) {
				cards, remainCards := groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))
				return append([]*ReCard{}, cards), remainCards
			}

			if seqNum[1] == int(g.cardGroup[0].cardValue) {
				cards, remainCards := groupJunko(hands, int(g.cardGroup[1].cardValue), int(g.cardGroup[6].cardValue))
				return append([]*ReCard{}, cards), remainCards
			} else if seqNum[1] == int(g.cardGroup[6].cardValue) {
				cards, remainCards := groupJunko(hands, int(g.cardGroup[0].cardValue), int(g.cardGroup[5].cardValue))
				return append([]*ReCard{}, cards), remainCards
			}
			cards, remainCards := groupJunko(hands, seqNum[0], seqNum[len(seqNum)])
			return append([]*ReCard{}, cards), remainCards
		}

	}

	// 然后在找出所有顺子
	// 2. 是否能找出两个顺子
	cards, remainCards := unlimitedJunko(hands)
	if len(cards) == 2 {
		reCards, remainCards := mergeJunkoWithSingle(cards, remainCards)
		return reCards, remainCards
	} else {
		cards, remainCards := FindPossibleLongSingleJunko(hands)
		return append([]*ReCard{}, cards), remainCards
	}
}