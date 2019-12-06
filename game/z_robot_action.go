package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/roomStatus"
	"time"
)

// 机器人抢地主阶段操作
func RobotGetLandlordAction(room *Room, robot, nextPlayer, lastPlayer *Player) {
	step := room.Status

	num := RandNum(3, 10)
	DelaySomeTime(time.Duration(num))

	logger.Debug("机器人抢地主阶段.............")

	// todo 机器人 抢地主阶段
	if step == roomStatus.CallLandlord {
		NotCallLandlordAction(room, robot, nextPlayer)
	} else if step == roomStatus.GetLandlord {
		NotGetLandlordAction(room, robot, nextPlayer, lastPlayer)
	} else {
		logger.Error("房间状态错误 !!!incredible")
	}
}

// 机器人打牌阶段操作
func RobotPlayAction(room *Room, robot, nextPlayer, lastPlayer *Player) {
	// 机器人打牌了
	num := RandNum(3, 10)
	DelaySomeTime(time.Duration(num))
	logger.Debug("机器人打牌阶段.............")

	DoGameHosting(room, robot, nextPlayer, lastPlayer)
}
