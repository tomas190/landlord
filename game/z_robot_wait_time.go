package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/sysSet"
	"time"
)

// 为让用户感觉更加真实 增加等带时间计算
/*
	60%的机率 让用户在1~5秒之间成功匹配
	35%的机率 让用户在5~10秒之间匹配
	5%的机率 让用户在10~15之间匹配
*/
// 玩家进入房间匹配机器人等待时间
func getWaitTimePlayerEnterRoom() time.Duration {
	destiny := RandNum(1, 100)
	if destiny <= 60 {
		logger.Debug("============ 正常几率 ============ ", destiny)
		delayTime := RandNum(1, 5)
		return time.Second * time.Duration(delayTime)
	}

	if destiny <= 95 {
		logger.Debug("============  35几率 ============ ", destiny)
		delayTime := RandNum(5, 10)
		return time.Second * time.Duration(delayTime)
	}
	logger.Debug("============  5几率  ============ ", destiny)
	delayTime := RandNum(10, 15)
	return time.Second * time.Duration(delayTime)
}

// 为让用户感觉更加真实 增加等待时间计算
/*
	70%的机率 让机器人在1~3秒之间
	20%的机率 让机器人在3~5秒之间
	8%的机率 让机器人在5~8之间
	2%的机率 让机器人在10-20之间

*/
// 机器人叫抢地主阶段决策等待时间
func getWaitTimeCallLandlord() time.Duration {
	destiny := RandNum(1, 100)
	if destiny <= 70 {
		delayTime := RandNum(1, 3)
		logger.Debug("============ 1叫抢地主阶段决策等待时间 ============ ", delayTime)
		return time.Duration(delayTime)
	}

	if destiny <= 90 {
		delayTime := RandNum(3, 5)
		logger.Debug("============ 2叫抢地主阶段决策等待时间 ============ ", delayTime)
		return time.Duration(delayTime)
	}

	if destiny <= 98 {
		delayTime := RandNum(5, 10)
		logger.Debug("============ 3叫抢地主阶段决策等待时间 ============ ", delayTime)
		return time.Duration(delayTime)
	}
	delayTime := RandNum(10, 20)
	logger.Debug("============ 4叫抢地主阶段决策等待时间 ============ ", delayTime)
	return time.Duration(delayTime)
}

/*
// 机器人出牌等待时间 概率
*/

// 建议概率 70%
// 为让用户感觉更加真实 增加等待时间计算
// 玩家一般玩牌 在正常情况会 在1到2秒 做出选择
// 机器人正常出牌速度
func getWaitTimeOutCardFast() time.Duration {
	delayTime := RandNum(1, 2)
	return time.Duration(delayTime)
}


// 建议概率 15%
// 为让用户感觉更加真实 增加等待时间计算
// 玩家一般玩牌 在正常情况会 在3到5秒 做出选择
// 机器人中等出牌速度
func getWaitTimeOutCardMedium() time.Duration {
	delayTime := RandNum(3, 5)
	return time.Duration(delayTime)
}

// 建议概率 7%
// 为让用户感觉更加真实 增加等待时间计算
// 玩家一般玩牌 在正常情况会 在6到15秒 做出选择
// 机器人慢出牌速度
func getWaitTimeOutCardSlowy() time.Duration {
	delayTime := RandNum(6, 15)
	return time.Duration(delayTime)
}

// 建议概率 6%
// 为让用户感觉更加真实 增加等待时间计算
// 玩家一般玩牌 在正常情况会 在5到10秒 做出选择
// 机器人极慢出牌速度
func getWaitTimeOutCardSoSlowy() time.Duration {
	delayTime := RandNum(16, 29)
	return time.Duration(delayTime)
}

// 建议概率 2%
// 为让用户感觉更加真实 增加等待时间计算
// 玩家一般玩牌 在正常情况会 在5到10秒 做出选择
// 机器人假装掉线  这个地方一地要配合 不出操作
func getWaitTimeOutCardFakerToDisconnection() time.Duration {
	return sysSet.GameDelayTime
}
