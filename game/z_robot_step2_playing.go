package game

import (
	"github.com/wonderivan/logger"
)

// 机器人打牌阶段操作
func RobotPlayAction(room *Room, robot, nextPlayer, lastPlayer *Player) {
	// 机器人打牌了
	//isFakerDisconnection := delayDestiny()
	//delayDestiny()
	logger.Debug("机器人打牌阶段...")
	//if isFakerDisconnection { // 如果概率出现了 假装掉线 则配合不出操作 并且机器人以后走托管流程
	//	logger.Debug("机器人打牌阶段 中0.001%概率掉线托管...")
	//	robot.IsGameHosting = true
	//	RespGameHosting(room, playerStatus.GameHosting, robot.PlayerPosition, robot.PlayerInfo.PlayerId)
	//	if robot.IsMustDo {
	//		DoGameHosting(room, robot, nextPlayer, lastPlayer)
	//	} else {
	//		NotOutCardsAction(room, robot, lastPlayer, nextPlayer)
	//	}
	//	return
	//}
	//DoGameHosting(room, robot, nextPlayer, lastPlayer)
	if robot.WaitingTime >= 3 {
		// 如果特殊情况机器人的等待时间是3秒 则快速出牌
		DelaySomeTime(getWaitTimeOutCardFast())
	} else {
		delayDestiny()
	}
	robotOutCard(room, robot, nextPlayer, lastPlayer)
}
