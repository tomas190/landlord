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

	// 如果实在最低级的房间 让机器人变得爱抢地主
	//if room.RoomClass.RoomType == roomType.ExperienceField {
	//	num := RandNum(0, 10)
	//	if num <= 7 {
	//		if step == roomStatus.CallLandlord {
	//			CallLandlordAction(room, robot, nextPlayer)
	//		} else if step == roomStatus.GetLandlord {
	//			GetLandlordAction(room, robot, nextPlayer, lastPlayer)
	//		}
	//	}else {
	//		if step == roomStatus.CallLandlord {
	//			NotCallLandlordAction(room, robot, nextPlayer)
	//		} else if step == roomStatus.GetLandlord {
	//			NotGetLandlordAction(room, robot, nextPlayer, lastPlayer)
	//		}
	//	}
	//	return
	//}

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
		// 异常叫抢地主操作
		// 如果正在玩 说明多余的机器人叫抢地主操作
		if step != roomStatus.Playing {
			// 异常不抢处理
			logger.Error("房间状态错误 !!!incredible",step)
			NotGetLandlordAction(room, robot, nextPlayer, lastPlayer)
		}

	}
}
