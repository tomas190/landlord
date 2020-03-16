package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/playerAction"
	"landlord/mconst/roomStatus"
	"landlord/mconst/sysSet"
	"landlord/msg/mproto"
	"runtime"
	"time"
)

/*
抢地主阶段
*/

// 3.1.叫地主阶段 和抢地主阶段
func CallLandlord(room *Room, playerId string) {
	actionPlayer := room.Players[playerId]
	if actionPlayer == nil {
		logger.Error("房间里无此用户...!!!incredible")
		return
	}

	// todo 每秒记录玩家的时间点用户 玩家再次阶段退出后 再次进入房间
	uptWtChin := make(chan struct{})

	go updatePlayerWaitingTime(actionPlayer, uptWtChin, sysSet.GameDelayGetLandlordTimeInt)
	// todo 每秒记录玩家的时间点用户 玩家再次阶段退出后 再次进入房间

	nextPosition := getNextPosition(actionPlayer.PlayerPosition)
	nextPlayer := getPlayerByPosition(room, nextPosition)

	lastPosition := getLastPosition(actionPlayer.PlayerPosition)
	lastPlayer := getPlayerByPosition(room, lastPosition)

	if actionPlayer.IsRobot { // 如果是机器人
		RobotGetLandlordAction(room, actionPlayer, nextPlayer, lastPlayer)
		return
	}

	// 阻塞等待当前玩家的动作 超过系统设置时间后自动处理
	select {
	case action := <-actionPlayer.ActionChan:
		go func() {
			uptWtChin <- struct{}{}
		}()
		switch action.ActionType {
		case playerAction.CallLandlord: // 叫地主动作
			CallLandlordAction(room, actionPlayer, nextPlayer)
		case playerAction.GetLandlord: // 抢地主动作
			GetLandlordAction(room, actionPlayer, nextPlayer, lastPlayer)
		case playerAction.NotCallLandlord: // 不叫
			NotCallLandlordAction(room, actionPlayer, nextPlayer)
		case playerAction.NotGetLandlord: // 不抢
			NotGetLandlordAction(room, actionPlayer, nextPlayer, lastPlayer)
		}
	case <-time.After(time.Second * sysSet.GameDelayGetLandlordTime): // 自动进行不叫或者不抢
		go func() {
			uptWtChin <- struct{}{}
		}()
		if room.Status == roomStatus.CallLandlord {
			NotCallLandlordAction(room, actionPlayer, nextPlayer) // 不叫
		} else if room.Status == roomStatus.GetLandlord {
			NotGetLandlordAction(room, actionPlayer, nextPlayer, lastPlayer) // 不抢
		}
	}
}

/* ==========================================  四大 action ===========================================*/

// 叫地主
func CallLandlordAction(room *Room, actionPlayer, nextPlayer *Player, ) {
	//playerId := actionPlayer.PlayerInfo.PlayerId
	// room.MultiAll = room.MultiAll * 2 叫地主不加倍
	actionPlayer.LastAction = playerAction.CallLandlord
	logger.Debug(actionPlayer.PlayerInfo.PlayerId, "做了一次 叫地主的动作...")
	room.Status = roomStatus.GetLandlord
	if nextPlayer.LastAction < playerAction.NoAction { // 如果下一个玩家的已做操作<0 那么他就是地主
		ensureWhoIsLandlord(room, actionPlayer, actionPlayer)
	} else { // 则让下一个玩家抢地主
		setCurrentPlayer(room, nextPlayer.PlayerInfo.PlayerId)
		pushCallLandlordHelp(room, actionPlayer, nextPlayer, playerAction.GetLandlord)
		CallLandlord(room, nextPlayer.PlayerInfo.PlayerId)
	}
}

// 抢地主
func GetLandlordAction(room *Room, actionPlayer, nextPlayer, lastPlayer *Player, ) {
	room.MultiAll = room.MultiAll * 2
	room.MultiGetLandlord = room.MultiGetLandlord * 2
	lastAction := actionPlayer.LastAction
	actionPlayer.LastAction = playerAction.GetLandlord
	logger.Debug(actionPlayer.PlayerInfo.PlayerId, "做了一次 抢地主动作...")
	if lastAction == playerAction.CallLandlord { // 如果玩家抢了地主 又已经叫过地主的情况下 那他就是地主
		ensureWhoIsLandlord(room, actionPlayer, actionPlayer)
	} else {
		// 如果下一个玩家不叫或者不抢  上一个玩家叫了地主 则该上一个玩家抢地主
		if nextPlayer.LastAction < playerAction.NoAction && lastPlayer.LastAction == playerAction.CallLandlord {
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
	actionPlayer.LastAction = playerAction.NotCallLandlord
	logger.Debug(actionPlayer.PlayerInfo.PlayerId, "做了一次 不叫...")
	if nextPlayer.LastAction == playerAction.NotCallLandlord { // 如果下一个玩家已经做了不叫的动作 重新发牌
		if room.reStartNum >= 1 {
			room.reStartNum = 0
			// 随机以为地主
			landPosition := int32(RandNum(1, 3))
			landPlayer := getPlayerByPosition(room, landPosition)
			ensureWhoIsLandlord(room, landPlayer, actionPlayer)
			// 2020年3月3日19:34:36 这里没加return可能会造成玩家卡死
			return
			// 2020年3月3日19:34:36 这里没加return可能会造成玩家卡死
		} else {
			room.reStartNum = room.reStartNum + 1
			logger.Debug("三个玩家都不叫 重新发牌")
			pushLastCallLandlord(room, actionPlayer)
			emptyPlayerCardInfo(room) // 清空数据

			DelaySomeTime(1)
			// 2020年2月22日16:02:23 bug控制针异常
			var truePlayerNum int
			//
			for _, p := range room.Players {
				if !p.IsRobot{
					truePlayerNum++
				}
			}

			if truePlayerNum==3 {
				PushPlayerStartGame(room)
			}else {
				PushPlayerStartGameWithRobot2(room)
			}

			//if Server.UseRobot {
			//	// PushPlayerStartGameWithRobot(room)
			//
			//	PushPlayerStartGameWithRobot2(room)
			//} else {
			//	PushPlayerStartGame(room)
			//}
		}
	} else { // 则让下一个玩家叫地主
		setCurrentPlayer(room, nextPlayer.PlayerInfo.PlayerId)
		pushCallLandlordHelp(room, actionPlayer, nextPlayer, playerAction.CallLandlord)
		CallLandlord(room, nextPlayer.PlayerInfo.PlayerId)
	}
}

// 不抢
func NotGetLandlordAction(room *Room, actionPlayer, nextPlayer, lastPlayer *Player, ) {
	actionPlayer.LastAction = playerAction.NotGetLandlord
	logger.Debug(actionPlayer.PlayerInfo.PlayerId, "做了一次 不抢...")
	if lastPlayer.LastAction < playerAction.NoAction { // 如果上一个玩家已经做了不抢的动作  那么下一个玩家就是地主
		ensureWhoIsLandlord(room, nextPlayer, actionPlayer)
	} else if nextPlayer.LastAction < playerAction.NoAction { // 如果下一个玩家已经做了不抢的动作  那么上一个玩家就是地主 1
		ensureWhoIsLandlord(room, lastPlayer, actionPlayer)
	} else if lastPlayer.LastAction == playerAction.GetLandlord && // 如果上一个玩家抢了地主 并且下一个玩家做了不叫或者不抢
		nextPlayer.LastAction < playerAction.NoAction { // 那么上一个玩家就是地主
		ensureWhoIsLandlord(room, lastPlayer, actionPlayer)
	} else if nextPlayer.LastAction == playerAction.GetLandlord &&
		lastPlayer.LastAction == playerAction.GetLandlord { // 如果上下两个玩家都抢了地主 那上一个玩家就是地主
		ensureWhoIsLandlord(room, lastPlayer, actionPlayer)
	} else {
		setCurrentPlayer(room, nextPlayer.PlayerInfo.PlayerId)
		pushCallLandlordHelp(room, actionPlayer, nextPlayer, playerAction.GetLandlord)
		CallLandlord(room, nextPlayer.PlayerInfo.PlayerId)
	}
}

/* ==========================================  四大 action ===========================================*/

/* ==================== 动作action 的消息推送  ==========================*/
// 3.第一次开始叫地主
func pushFirstCallLandlord(room *Room) string {
	lastPosition := int32(RandNum(1, 3))
	lastPlayer := getPlayerByPosition(room, lastPosition)

	actionPosition := getNextPosition(lastPosition)
	actionPlayer := getPlayerByPosition(room, actionPosition)
	setCurrentPlayer(room, actionPlayer.PlayerInfo.PlayerId)
	pushCallLandlordHelp(room, lastPlayer, actionPlayer, playerAction.CallLandlord)
	return actionPlayer.PlayerInfo.PlayerId
}

// 抢地主阶段辅助推送
/*
最后一个玩家的动作决定了谁是地主但是要显示这个玩家发出的动作
*/
func pushLastCallLandlord(room *Room, lastPlayer *Player) {
	var push mproto.PushGetLandlord
	push.Action = room.Status
	push.LastPlayerPosition = lastPlayer.PlayerPosition
	push.LastPlayerId = lastPlayer.PlayerInfo.PlayerId
	push.LastPlayerAction = lastPlayer.LastAction
	push.Countdown = sysSet.GameDelayGetLandlordTimeInt
	push.Multi = room.MultiAll

	bytes, _ := proto.Marshal(&push)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushCallLandlord, bytes))

}

// 抢地主阶段辅助推送
func pushCallLandlordHelp(room *Room, lastPlayer, nextPlayer *Player, showAction int32) {
	var push mproto.PushGetLandlord
	push.Action = room.Status
	push.LastPlayerPosition = lastPlayer.PlayerPosition
	push.LastPlayerId = lastPlayer.PlayerInfo.PlayerId
	push.LastPlayerAction = lastPlayer.LastAction

	push.PlayerPosition = nextPlayer.PlayerPosition
	push.PlayerId = nextPlayer.PlayerInfo.PlayerId
	push.Countdown = sysSet.GameDelayGetLandlordTimeInt
	push.Action = showAction
	push.Multi = room.MultiAll

	bytes, _ := proto.Marshal(&push)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushCallLandlord, bytes))
}

// 推送地主玩家
func pushWhoIsLandlord(room *Room, landlordPlayer *Player) {

	landlordPlayer.HandCards = append(landlordPlayer.HandCards, room.BottomCards...)
	SortCard(landlordPlayer.HandCards)
	var push mproto.PushLandlord
	push.LandlordId = landlordPlayer.PlayerInfo.PlayerId
	push.Cards = ChangeCardToProto(room.BottomCards)
	push.Position = landlordPlayer.PlayerPosition
	bytes, _ := proto.Marshal(&push)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushWhoIsLandlord, bytes))

}

/* ==================== 动作action 的消息推送  ==========================*/

// 已经确定谁是地主
func ensureWhoIsLandlord(room *Room, landlordPlayer, actionPlayer *Player) {
	setCurrentPlayer(room, landlordPlayer.PlayerInfo.PlayerId)
	landlordPlayer.IsLandlord = true
	room.Status = roomStatus.Playing

	room.LandlordPlayerId = landlordPlayer.PlayerInfo.PlayerId

	logger.Debug("=============== 玩牌开始 ===========")
	logger.Debug("地主玩家:", landlordPlayer.PlayerInfo.PlayerId)
	pushLastCallLandlord(room, actionPlayer)

	DelaySomeTime(1)
	pushWhoIsLandlord(room, landlordPlayer)

	//
	reSetOutRoomToOut(room, landlordPlayer.PlayerInfo.PlayerId)         // 清空玩家动作
	setCurrentPlayerOut(room, landlordPlayer.PlayerInfo.PlayerId, true) // 设置位当前操作玩家

	//pushMustOutCard(room, landlordPlayer.PlayerInfo.PlayerId)
	go pushFirstMustOutCard(room, landlordPlayer.PlayerInfo.PlayerId)
	go pushCardCount(room)
	go PlayingGame(room, landlordPlayer.PlayerInfo.PlayerId)

}

// 叫地主阶段 保存更新用户等待时间点
func updatePlayerWaitingTime(actionPlayer *Player, tmpChan chan struct{}, waitTime int32) {
	actionPlayer.WaitingTime = waitTime
	for {
		select {
		case <-tmpChan:
			logger.Debug("玩家已经确认操作:操作时间点:", actionPlayer.WaitingTime)
			//actionPlayer.WaitingTime = waitTime
			runtime.Goexit()
		case <-time.After(time.Second):
			if actionPlayer.WaitingTime <= 0 {
				runtime.Goexit()
			}
			actionPlayer.WaitingTime = actionPlayer.WaitingTime - 1
		}
	}
}
