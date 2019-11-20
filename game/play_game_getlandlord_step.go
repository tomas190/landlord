package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/playerAction"
	"landlord/mconst/roomStatus"
	"landlord/mconst/sysSet"
	"landlord/msg/mproto"
)

/* ==========================================  四大 action ===========================================*/

// 叫地主
func CallLandlordAction(room *Room, actionPlayer, nextPlayer *Player, ) {
	//playerId := actionPlayer.PlayerInfo.PlayerId
	actionPlayer.DidAction = playerAction.CallLandlord
	logger.Debug(actionPlayer.PlayerInfo.PlayerId, "做了一次 叫地主的动作...")
	room.Status = roomStatus.GetLandlord
	if nextPlayer.DidAction < playerAction.NoAction { // 如果下一个玩家的已做操作<0 那么他就是地主
		ensureWhoIsLandlord(room, actionPlayer, actionPlayer)
	} else { // 则让下一个玩家抢地主
		setCurrentPlayer(room, nextPlayer.PlayerInfo.PlayerId)
		pushCallLandlordHelp(room, actionPlayer, nextPlayer, playerAction.GetLandlord)
		CallLandlord(room, nextPlayer.PlayerInfo.PlayerId)
	}
}

// 抢地主
func GetLandlordAction(room *Room, actionPlayer, nextPlayer, lastPlayer *Player, ) {
	lastAction := actionPlayer.DidAction
	actionPlayer.DidAction = playerAction.GetLandlord
	logger.Debug(actionPlayer.PlayerInfo.PlayerId, "做了一次 抢地主动作...")
	if lastAction == playerAction.CallLandlord { // 如果玩家抢了地主 又已经叫过地主的情况下 那他就是地主
		ensureWhoIsLandlord(room, actionPlayer, actionPlayer)
	} else {
		// 如果下一个玩家不叫或者不抢  上一个玩家叫了地主 则该上一个玩家抢地主
		if nextPlayer.DidAction < playerAction.NoAction && lastPlayer.DidAction == playerAction.CallLandlord {
			setCurrentPlayer(room, nextPlayer.PlayerInfo.PlayerId)
			pushCallLandlordHelp(room, actionPlayer, lastPlayer, playerAction.GetLandlord)
			CallLandlord(room, lastPlayer.PlayerInfo.PlayerId)
		} else { // 则让下一个玩家抢地主
			setCurrentPlayer(room, nextPlayer.PlayerInfo.PlayerId)
			pushCallLandlordHelp(room, actionPlayer, nextPlayer, playerAction.GetLandlord)
			CallLandlord(room, nextPlayer.PlayerInfo.PlayerId)
		}

	}
}

// 不叫
func NotCallLandlordAction(room *Room, actionPlayer, nextPlayer *Player, ) {
	actionPlayer.DidAction = playerAction.NotCallLandlord
	logger.Debug(actionPlayer.PlayerInfo.PlayerId, "做了一次 不叫...")
	if nextPlayer.DidAction == playerAction.NotCallLandlord { // 如果下一个玩家已经做了不叫的动作 重新发牌
		logger.Debug("三个玩家都不叫 重新发牌")
		emptyPlayerCardInfo(room) // 清空数据
		PushPlayerStartGame(room)
	} else { // 则让下一个玩家叫地主
		setCurrentPlayer(room, nextPlayer.PlayerInfo.PlayerId)
		pushCallLandlordHelp(room, actionPlayer, nextPlayer, playerAction.CallLandlord)
		CallLandlord(room, nextPlayer.PlayerInfo.PlayerId)
	}
}

// 不抢
func NotGetLandlordAction(room *Room, actionPlayer, nextPlayer, lastPlayer *Player, ) {
	actionPlayer.DidAction = playerAction.NotGetLandlord
	logger.Debug(actionPlayer.PlayerInfo.PlayerId, "做了一次 不抢...")
	if lastPlayer.DidAction < playerAction.NoAction { // 如果上一个玩家已经做了不抢的动作  那么下一个玩家就是地主
		ensureWhoIsLandlord(room, nextPlayer, actionPlayer)
	} else if nextPlayer.DidAction < playerAction.NoAction { // 如果下一个玩家已经做了不抢的动作  那么上一个玩家就是地主
		ensureWhoIsLandlord(room, lastPlayer, actionPlayer)
	} else if lastPlayer.DidAction == playerAction.GetLandlord &&// 如果上一个玩家抢了地主 并且下一个玩家做了不叫或者不抢
		nextPlayer.DidAction < playerAction.NoAction { // 那么上一个玩家就是地主
		ensureWhoIsLandlord(room, lastPlayer, actionPlayer)
	} else if nextPlayer.DidAction == playerAction.GetLandlord { // 如果下一个玩家抢了地主 那下一个玩家就是地主
		ensureWhoIsLandlord(room, nextPlayer, actionPlayer)

	} else {
		setCurrentPlayer(room, nextPlayer.PlayerInfo.PlayerId)
		pushCallLandlordHelp(room, actionPlayer, nextPlayer, playerAction.GetLandlord)
		CallLandlord(room, nextPlayer.PlayerInfo.PlayerId)
	}
}

/* ==========================================  四大 action ===========================================*/

// 已经确定谁是地主
func ensureWhoIsLandlord(room *Room, landlordPlayer, actionPlayer *Player) {
	setCurrentPlayer(room, landlordPlayer.PlayerInfo.PlayerId)
	landlordPlayer.IsLandlord = true
	logger.Debug("=============== 玩牌开始 ===========")
	logger.Debug("地主玩家:", landlordPlayer.PlayerInfo.PlayerId)
	// todo 推送地主牌 开始玩牌逻辑
	pushCallLandlordLastAction(room, actionPlayer)
	pushWhoIsLandlord(room, landlordPlayer)
}

// 抢地主阶段辅助推送
func pushCallLandlordHelp(room *Room, lastPlayer, nextPlayer *Player, showAction int32) {
	var push mproto.PushGetLandlord
	push.Action = room.Status
	push.LastPlayerPosition = lastPlayer.PlayerPosition
	push.LastPlayerId = lastPlayer.PlayerInfo.PlayerId
	push.LastPlayerAction = lastPlayer.DidAction

	push.PlayerPosition = nextPlayer.PlayerPosition
	push.PlayerId = nextPlayer.PlayerInfo.PlayerId
	push.Countdown = sysSet.GameDelayTimeInt
	push.Action = showAction

	bytes, _ := proto.Marshal(&push)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushCallLandlord, bytes))
}

// 抢地主阶段辅助推送
/*
最后一个玩家的动作决定了谁是地主但是要显示这个玩家发出的动作

*/
func pushCallLandlordLastAction(room *Room, lastPlayer *Player) {
	var push mproto.PushGetLandlord
	push.Action = room.Status
	push.LastPlayerPosition = lastPlayer.PlayerPosition
	push.LastPlayerId = lastPlayer.PlayerInfo.PlayerId
	push.LastPlayerAction = lastPlayer.DidAction
	push.Countdown = sysSet.GameDelayTimeInt

	bytes, _ := proto.Marshal(&push)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushCallLandlord, bytes))

}

// 推送地主玩家

func pushWhoIsLandlord(room *Room, landlordPlayer *Player) {

	landlordPlayer.HandCards = append(landlordPlayer.HandCards, room.bottomCards...)
	var push mproto.PushLandlord
	push.LandlordId = landlordPlayer.PlayerInfo.PlayerId
	push.Cards = ChangeCardToProto(room.bottomCards)
	push.Position = landlordPlayer.PlayerPosition
	bytes, _ := proto.Marshal(&push)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushWhoIsLandlord, bytes))

}

// 清空玩家
func emptyPlayerCardInfo(room *Room) {
	for _, v := range room.Players {
		v.IsCanDo = false
		v.DidAction = 0
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
		logger.Error("!!!!!!!!!incredible！1", currentPosition)
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
		logger.Error("!!!!!!!!!incredible！2", lastPosition)
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
		logger.Error("!!!!!!!!!incredible！3", currentPosition)
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

	logger.Debug("!!!!!!!!!incredible")
	return nil
}
