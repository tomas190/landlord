package game

import (
	"github.com/wonderivan/logger"
	"landlord/mconst/roomStatus"
	"landlord/mconst/roomType"
)

// 机器人抢地主阶段操作
func RobotGetLandlordAction(room *Room, robot, nextPlayer, lastPlayer *Player) {
	step := room.Status
	DelaySomeTime(getWaitTimeCallLandlord())
	logger.Debug("机器人抢地主阶段.............")

	// 如果实在最低级的房间 让机器人变得爱抢地主
	if room.RoomClass.RoomType == roomType.ExperienceField {
		num := RandNum(0, 10)
		if num <= 7 {
			if step == roomStatus.CallLandlord {
				CallLandlordAction(room, robot, nextPlayer)
			} else if step == roomStatus.GetLandlord {
				GetLandlordAction(room, robot, nextPlayer, lastPlayer)
			}
		}else {
			if step == roomStatus.CallLandlord {
				NotCallLandlordAction(room, robot, nextPlayer)
			} else if step == roomStatus.GetLandlord {
				NotGetLandlordAction(room, robot, nextPlayer, lastPlayer)
			}
		}
		return
	}

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
