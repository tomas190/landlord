package game

import "github.com/wonderivan/logger"

// 清空玩家
func emptyPlayerCardInfo(room *Room) {
	for _, v := range room.Players {
		v.IsCanDo = false
		v.LastAction = 0
		v.HandCards = nil
		v.ThrowCards = nil
	}
}

// 设置当前操作玩家
func setCurrentPlayer(room *Room, playerId string) {
	for _, v := range room.Players {
		if v.PlayerInfo.PlayerId == playerId {
			v.IsCanDo = true
		} else {
			v.IsCanDo = false
		}
	}

}

// 根据当前玩家的位置获取上一个玩家的位置
func getLastPosition(currentPosition int32) int32 {
	switch currentPosition {
	case 3:
		return 2
	case 2:
		return 1
	case 1:
		return 3
	default:
		logger.Error("getLastPosition !!!incredible", currentPosition)
		return 0
	}

}

// 根据上个玩家的位置获取当前玩家的位置
func getCurrentPosition(lastPosition int32) int32 {
	switch lastPosition {
	case 3:
		return 1
	case 2:
		return 3
	case 1:
		return 2
	default:
		logger.Error("getCurrentPosition !!!incredible！", lastPosition)
		return 0
	}
}

// 根据当前玩家的位置获取下一个玩家的位置
func getNextPosition(currentPosition int32) int32 {
	switch currentPosition {
	case 3:
		return 1
	case 2:
		return 3
	case 1:
		return 2
	default:
		logger.Error("getNextPosition !!!incredible", currentPosition)
		return 0
	}
}

// 根据位置获取玩家
func getPlayerByPosition(room *Room, position int32) *Player {
	for _, v := range room.Players {
		if v.PlayerPosition == position {
			return v
		}
	}

	logger.Debug("getPlayerByPosition !!!incredible")
	return nil
}

// 根据位置获取玩家
func getPlayerByPlayerId(room *Room, playerId string) *Player {
	player, ok := room.Players[playerId]
	if ok {
		return player
	}
	logger.Debug("getPlayerByPlayerId !!!incredible")
	return nil
}
