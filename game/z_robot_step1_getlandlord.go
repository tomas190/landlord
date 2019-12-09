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
		NotCallLandlordAction(room, robot, nextPlayer)
	} else if step == roomStatus.GetLandlord {
		NotGetLandlordAction(room, robot, nextPlayer, lastPlayer)
	} else {
		logger.Error("房间状态错误 !!!incredible")
	}
}
