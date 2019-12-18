package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/roomStatus"
)

// 机器人抢地主阶段操作
func RobotGetLandlordAction(room *Room, robot, nextPlayer, lastPlayer *Player) {
	step := room.Status
	DelaySomeTime(getWaitTimeCallLandlord())
	logger.Debug("机器人抢地主阶段.............")

	// todo 机器人 抢地主阶段
	if step == roomStatus.CallLandlord {
		if robot.HandsValue >= 6 {
			CallLandlordAction(room, robot, nextPlayer)
		} else {
			NotCallLandlordAction(room, robot, nextPlayer)
		}
	} else if step == roomStatus.GetLandlord {
		if robot.HandsValue >= 8 {
			GetLandlordAction(room, robot, nextPlayer, lastPlayer)
		} else {
			NotGetLandlordAction(room, robot, nextPlayer, lastPlayer)
		}
	} else {
		logger.Error("房间状态错误 !!!incredible")
	}
}
