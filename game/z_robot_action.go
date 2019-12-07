package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/roomStatus"
	"landlord/msg/mproto"
	"time"
)

// 开始和机器人玩游戏
func PlayGameWithRobot(room *Room) {

	// 1. 玩家进入房间如果有玩家正待等待则与之开始游戏
	PushPlayerEnterRoom(room)
	DelaySomeTime(1)

	// 2.给玩家发牌
	PushPlayerStartGameWithRobot(room)

	// ..．流程控制到这里结束　发牌  抢地主  玩牌 直接由 PushPlayerStartGame 开始 且循环

}

// 2.给玩家和机器人发牌
func PushPlayerStartGameWithRobot(room *Room) {
	cards := CreateBrokenCard()
	//cards := CreateSortCard()
	players := room.Players
	for _, v := range players {
		v.HandCards = append([]*Card{}, cards[:17]...)
		SortCard(v.HandCards)
		logger.Debug("玩家" + v.PlayerInfo.PlayerId + "的牌：")
		PrintCard(v.HandCards)
		cards = append([]*Card{}, cards[17:]...)
		var push mproto.PushStartGame
		push.Cards = ChangeCardToProto(v.HandCards)
		bytes, _ := proto.Marshal(&push)
		PlayerSendMsg(v, PkgMsg(msgIdConst.PushStartGame, bytes))
	}
	room.BottomCards = append([]*Card{}, cards...)
	logger.Debug("底牌:")
	PrintCard(cards)
	room.Status = roomStatus.CallLandlord

	// 随机叫地主写在发牌里面 是因为三个玩家如果都不叫 则可以直接调用 PushPlayerStartGame 重新开始发牌逻辑
	DelaySomeTime(4)
	// 3.随机叫地主
	actionPlayerId := pushFirstCallLandlord(room)
	CallLandlord(room, actionPlayerId)
}

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
