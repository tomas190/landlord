package game

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/sysSet"
	"landlord/msg/mproto"
	"time"
)

func pushSpring(room *Room) {
	var push mproto.PushSpring
	bytes, _ := proto.Marshal(&push)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushSpring, bytes))

}

// 结算
// 最小金额计算 玩家只有这么多金币 则 只能输或者赢这么多
func Settlement(room *Room, winPlayer *Player) {
	// 1. 计算基本倍数

	mult := room.MultiAll
	settlementGold := room.RoomClass.BottomPoint * float64(mult)

	landPlayer, fp1, fp2 := getPlayerClass(room)
	roundId := fmt.Sprintf("room-%d-%d", room.RoomClass.RoomType, time.Now().Unix())

	var sPush mproto.PushSettlement

	// 如果赢家是地主
	if winPlayer.IsLandlord == true {
		var landRealWinGold float64 // 地主实际赢钱 税前
		if fp1.PlayerInfo.Gold < settlementGold { // 如果玩家1 的钱不够开
			showWinLossGold := fmt.Sprintf("-%.2f", fp1.PlayerInfo.Gold)
			landRealWinGold += fp1.PlayerInfo.Gold
			syncLossGold(fp1, fp1.PlayerInfo.Gold, roundId) // 同步金币 到中心服务 session

			ss := getSelfSettlement(room, fp1, -1, showWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, ss)
		} else {
			landRealWinGold += settlementGold
			syncLossGold(fp1, settlementGold, roundId) // 同步金币 到中心服务 session

			showWinLossGold := fmt.Sprintf("-%.2f", settlementGold)
			ss := getSelfSettlement(room, fp1, -1, showWinLossGold, false)
			sPush.Settlement = append(sPush.Settlement, ss)
		}

		if fp2.PlayerInfo.Gold < settlementGold { // 如果玩家2 的钱不够开
			showWinLossGold := fmt.Sprintf("-%.2f", fp2.PlayerInfo.Gold)
			landRealWinGold += fp2.PlayerInfo.Gold
			syncLossGold(fp2, fp2.PlayerInfo.Gold, roundId) // 同步金币 到中心服务 session

			ss := getSelfSettlement(room, fp2, -1, showWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, ss)
		} else {
			landRealWinGold += settlementGold
			syncLossGold(fp2, settlementGold, roundId) // 同步金币 到中心服务 session

			showWinLossGold := fmt.Sprintf("-%.2f", settlementGold)
			ss := getSelfSettlement(room, fp2, -1, showWinLossGold, false)
			sPush.Settlement = append(sPush.Settlement, ss)
		}

		landRealWinGoldPay := landRealWinGold * (1 - Server.GameTaxRate)      // 地主实际赢钱 税后
		syncWinGold(landPlayer, landRealWinGold, landRealWinGoldPay, roundId) // 同步金币 到中心服务 session

		showWinLossGold := fmt.Sprintf("%.2f", landRealWinGoldPay)
		ss := getSelfSettlement(room, landPlayer, 1, showWinLossGold, landRealWinGold < settlementGold*2)
		sPush.Settlement = append(sPush.Settlement, ss)

	} else { // 如果玩家不是地主
		// 1. 判断地主金币是否够开
		if landPlayer.PlayerInfo.Gold/2 < settlementGold {
			landShowWinLossGold := fmt.Sprintf("-%.2f", landPlayer.PlayerInfo.Gold)

			farmerRealWinGold := landPlayer.PlayerInfo.Gold / 2
			farmerRealWinGoldPay := farmerRealWinGold * (1 - Server.GameTaxRate)

			syncWinGold(fp1, settlementGold, farmerRealWinGoldPay, roundId)
			syncWinGold(fp2, settlementGold, farmerRealWinGoldPay, roundId)
			syncLossGold(landPlayer, landPlayer.PlayerInfo.Gold, roundId)
			//
			logger.Debug("地主玩家输钱不够开", landPlayer.PlayerInfo.Gold)
			logger.Debug("结算金额基*1", settlementGold)
			logger.Debug("结算金额基*2", settlementGold*2)

			showWinLossGold := fmt.Sprintf("%.2f", farmerRealWinGoldPay)
			fs1 := getSelfSettlement(room, fp1, 1, showWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, fs1)

			fs2 := getSelfSettlement(room, fp2, 1, showWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, fs2)

			ls := getSelfSettlement(room, landPlayer, -1, landShowWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, ls)

		} else {
			// 正常结算
			winGoldPay := settlementGold * (1 - Server.GameTaxRate)
			syncWinGold(fp1, settlementGold, winGoldPay, roundId)
			syncWinGold(fp2, settlementGold, winGoldPay, roundId)
			syncLossGold(landPlayer, settlementGold*2, roundId)

			showWinLossGold := fmt.Sprintf("%.2f", winGoldPay)
			fs1 := getSelfSettlement(room, fp1, 1, showWinLossGold, false)
			sPush.Settlement = append(sPush.Settlement, fs1)

			fs2 := getSelfSettlement(room, fp2, 1, showWinLossGold, false)
			sPush.Settlement = append(sPush.Settlement, fs2)

			landShowWinLossGold := fmt.Sprintf("-%.2f", settlementGold*2)
			ls := getSelfSettlement(room, landPlayer, -1, landShowWinLossGold, false)
			sPush.Settlement = append(sPush.Settlement, ls)
		}

	}

	var multiInfo mproto.MultipleInfo
	multiInfo.FightLandlord = fmt.Sprintf("×%d", room.MultiGetLandlord)
	multiInfo.Boom = fmt.Sprintf("×%d", room.MultiBoom)
	multiInfo.Spring = fmt.Sprintf("×%d", room.MultiSpring)
	sPush.MultipleInfo = &multiInfo
	sPush.WaitTime = sysSet.GameDelayReadyTimeInt

	bytes, _ := proto.Marshal(&sPush)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushSettlement, bytes))

}

func syncWinGold(player *Player, gold, goldPay float64, roundId string) float64 {
	orderId := fmt.Sprintf("%s-%s-win", roundId, player.PlayerInfo.PlayerId)
	player.PlayerInfo.Gold = player.PlayerInfo.Gold + goldPay // 同步到房间id
	err := SetSessionGold(player.Session, goldPay)            // 同步到session
	if err != nil {
		logger.Error("同步进步到session失败: !!!incredible")
	}
	UserSyncWinScore(player.PlayerInfo.PlayerId, gold, roundId, orderId) // 同步到中心服务

	// 赢钱超过设定值发送 跑马灯
	if !player.IsRobot && goldPay > Server.WinGoldNotice {
		NoticeWinMoreThan(player.PlayerInfo.PlayerId, player.PlayerInfo.Name, goldPay)
	}

	return player.PlayerInfo.Gold
}

func syncLossGold(player *Player, gold float64, roundId string) float64 {
	orderId := fmt.Sprintf("%s-%s-loss", roundId, player.PlayerInfo.PlayerId)
	player.PlayerInfo.Gold = player.PlayerInfo.Gold - gold
	err := SetSessionGold(player.Session, -gold) // 同步到session
	if err != nil {
		logger.Error("同步进步到session失败: !!!incredible")
	}
	UserSyncLoseScore(player.PlayerInfo.PlayerId, -gold, roundId, orderId)
	return player.PlayerInfo.Gold
}

func getSelfSettlement(room *Room, player *Player, winOrFail int32, winOrLossGold string, isMinSett bool) *mproto.Settlement {
	var result mproto.Settlement

	if player.IsLandlord {
		result.IsLandlord = 1
		result.Multiple = room.MultiAll * 2
	} else {
		result.IsLandlord = -1
		result.Multiple = room.MultiAll
	}
	result.PlayerId = player.PlayerInfo.PlayerId
	result.Position = player.PlayerPosition
	result.CurrentGold = player.PlayerInfo.Gold
	result.PlayerName = player.PlayerInfo.Name
	result.WinOrFail = winOrFail
	result.WinLossGold = winOrLossGold
	result.RemainCards = ChangeCardToProto(player.HandCards)
	result.MinSettlement = isMinSett
	return &result
}
