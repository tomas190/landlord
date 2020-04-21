package game

import (
	"github.com/wonderivan/logger"
)

// 机器人打牌阶段操作
func RobotPlayAction(room *Room, robot, nextPlayer, lastPlayer *Player) {
	// 机器人打牌了
	logger.Debug("机器人打牌阶段...")

	if robot.WaitingTime > 3 {
		delayDestiny()
	} else {
		// 如果特殊情况机器人的等待时间是3秒 则快速出牌
		DelaySomeTime(getWaitTimeOutCardFast())
	}
	//go func() {
	//	uptWtChin <- struct{}{}
	//}()
	robotOutCard(room, robot, nextPlayer, lastPlayer)
}
