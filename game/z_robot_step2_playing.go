package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/playerStatus"
)

// 机器人打牌阶段操作
func RobotPlayAction(room *Room, robot, nextPlayer, lastPlayer *Player) {
	// 机器人打牌了
	isFakerDisconnection := delayDestiny()
	if isFakerDisconnection && !robot.IsMustDo { // 如果概率出现了 假装掉线 则配合不出操作 并且机器人以后走托管流程
		logger.Debug("机器人打牌阶段 中2%概率掉线托管...")
		robot.IsGameHosting = true
		RespGameHosting(room, playerStatus.GameHosting, robot.PlayerPosition, robot.PlayerInfo.PlayerId)
		NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
		return
	}
	logger.Debug("机器人打牌阶段...")
	DoGameHosting(room, robot, nextPlayer, lastPlayer)
}


// 机器人出牌决策等待时间
func delayDestiny() bool {
	destiny := RandNum(1, 100)
	if destiny <= 70 {
		logger.Debug("机器人打牌阶段决策时间:快速 1-2s")
		DelaySomeTime(getWaitTimeOutCardFast())
		return false
	}

	if destiny <= 85 {
		logger.Debug("机器人打牌阶段决策时间:中速 3-5s")
		DelaySomeTime(getWaitTimeOutCardMedium())
		return false
	}

	if destiny <= 92 {
		logger.Debug("机器人打牌阶段决策时间:慢速 6-15s")
		DelaySomeTime(getWaitTimeOutCardSlowy())
		return false
	}

	if destiny <= 98 {
		logger.Debug("机器人打牌阶段决策时间:超慢 15-29s")
		DelaySomeTime(getWaitTimeOutCardSoSlowy())
		return false
	}
	logger.Debug("机器人打牌阶段决策时间:假装短线 30s")
	DelaySomeTime(getWaitTimeOutCardFakerToDisconnection())
	return true
}
