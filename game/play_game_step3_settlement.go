package game

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/sysSet"
	"landlord/msg/mproto"
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
	origiSettlementGold := settlementGold

	landPlayer, fp1, fp2 := getPlayerClass(room)
	//roundId := fmt.Sprintf("room-%d-%d", room.RoomClass.RoomType, time.Now().Unix())
	roundId := room.RoundId
	//order := bson.NewObjectId().String()

	var sPush mproto.PushSettlement

	// 如果赢家是地主
	if winPlayer.IsLandlord == true {
		var landRealWinGold float64 // 地主实际赢钱 税前
		// 1.如果地主的钱小于赢钱  自己本身多少钱 就只能赢这么多钱
		if landPlayer.PlayerInfo.Gold < settlementGold*2 {
			settlementGold = landPlayer.PlayerInfo.Gold / 2
		}

		// 2.如果玩家一的钱不够开
		if fp1.PlayerInfo.Gold < settlementGold { // 如果玩家1 的钱不够开
			showWinLossGold := fmt.Sprintf("-%.2f", fp1.PlayerInfo.Gold)
			landRealWinGold += fp1.PlayerInfo.Gold
			syncLossGold(fp1, fp1.PlayerInfo.Gold, roundId, *room.RoomClass) // 同步金币 到中心服务 session

			ss := getSelfSettlement(room, fp1, -1, showWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, ss)
		} else {
			landRealWinGold += settlementGold
			syncLossGold(fp1, settlementGold, roundId, *room.RoomClass) // 同步金币 到中心服务 session

			showWinLossGold := fmt.Sprintf("-%.2f", settlementGold)
			ss := getSelfSettlement(room, fp1, -1, showWinLossGold, false)
			sPush.Settlement = append(sPush.Settlement, ss)
		}

		// 如果玩家2的钱不够开
		if fp2.PlayerInfo.Gold < settlementGold { // 如果玩家2 的钱不够开
			showWinLossGold := fmt.Sprintf("-%.2f", fp2.PlayerInfo.Gold)
			landRealWinGold += fp2.PlayerInfo.Gold
			syncLossGold(fp2, fp2.PlayerInfo.Gold, roundId, *room.RoomClass) // 同步金币 到中心服务 session

			ss := getSelfSettlement(room, fp2, -1, showWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, ss)
		} else {
			landRealWinGold += settlementGold
			syncLossGold(fp2, settlementGold, roundId, *room.RoomClass) // 同步金币 到中心服务 session

			showWinLossGold := fmt.Sprintf("-%.2f", settlementGold)
			ss := getSelfSettlement(room, fp2, -1, showWinLossGold, false)
			sPush.Settlement = append(sPush.Settlement, ss)
		}

		landRealWinGoldPay := landRealWinGold * (1 - Server.GameTaxRate)                       // 地主实际赢钱 税后
		syncWinGold(landPlayer, landRealWinGold, landRealWinGoldPay, roundId, *room.RoomClass) // 同步金币 到中心服务 session

		showWinLossGold := fmt.Sprintf("%.2f", landRealWinGoldPay)
		ss := getSelfSettlement(room, landPlayer, 1, showWinLossGold, landRealWinGold < origiSettlementGold*2)
		sPush.Settlement = append(sPush.Settlement, ss)

	} else { // 如果玩家不是地主
		// 1.如果农民的钱小于结算赢钱
		fp1S := settlementGold
		fp2S := settlementGold
		if fp1.PlayerInfo.Gold < settlementGold || fp2.PlayerInfo.Gold < settlementGold {
			if fp1.PlayerInfo.Gold < settlementGold {
				fp1S = fp1.PlayerInfo.Gold
			}

			if fp2.PlayerInfo.Gold < settlementGold {
				fp2S = fp2.PlayerInfo.Gold
			}
		}

		// 2. 判断地主金币是否够开
		if landPlayer.PlayerInfo.Gold < fp1S+fp2S {
			//farmerRealWinGold := landPlayer.PlayerInfo.Gold / 2
			//farmerRealWinGoldPay := farmerRealWinGold * (1 - Server.GameTaxRate)
			farmerWinGold := landPlayer.PlayerInfo.Gold / 2
			f1RealWin := farmerWinGold
			f2RealWin := farmerWinGold
			var landRealLoss float64
			var f1RealWinPay, f2RealWinPay float64

			if fp1S < farmerWinGold {
				f1RealWin = fp1S
			}
			if fp2S < farmerWinGold {
				f2RealWin = fp2S
			}

			landRealLoss = f1RealWin + f2RealWin

			f1RealWinPay = f1RealWin * (1 - Server.GameTaxRate)
			f2RealWinPay = f2RealWin * (1 - Server.GameTaxRate)

			syncWinGold(fp1, f1RealWin, f1RealWinPay, roundId, *room.RoomClass)
			syncWinGold(fp2, f2RealWin, f2RealWinPay, roundId, *room.RoomClass)
			syncLossGold(landPlayer, landRealLoss, roundId, *room.RoomClass)
			//
			//logger.Debug("地主玩家输钱不够开", landPlayer.PlayerInfo.Gold)
			//logger.Debug("结算金额基*1", settlementGold)
			//logger.Debug("结算金额基*2", settlementGold*2)

			f1ShowWinLossGold := fmt.Sprintf("%.2f", f1RealWinPay)
			fs1 := getSelfSettlement(room, fp1, 1, f1ShowWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, fs1)

			f2ShowWinLossGold := fmt.Sprintf("%.2f", f2RealWinPay)
			fs2 := getSelfSettlement(room, fp2, 1, f2ShowWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, fs2)

			landShowWinLossGold := fmt.Sprintf("-%.2f", landRealLoss)
			ls := getSelfSettlement(room, landPlayer, -1, landShowWinLossGold, true)
			sPush.Settlement = append(sPush.Settlement, ls)

		} else {
			// 正常结算
			fp1WinGoldPay := fp1S * (1 - Server.GameTaxRate)
			syncWinGold(fp1, fp1S, fp1WinGoldPay, roundId, *room.RoomClass)

			fp2WinGoldPay := fp2S * (1 - Server.GameTaxRate)
			syncWinGold(fp2, fp2S, fp2WinGoldPay, roundId, *room.RoomClass)

			syncLossGold(landPlayer, fp1S+fp2S, roundId, *room.RoomClass)

			fp1ShowWinLossGold := fmt.Sprintf("%.2f", fp1WinGoldPay)
			fs1 := getSelfSettlement(room, fp1, 1, fp1ShowWinLossGold, fp1S < settlementGold)
			sPush.Settlement = append(sPush.Settlement, fs1)

			fp2ShowWinLossGold := fmt.Sprintf("%.2f", fp2WinGoldPay)
			fs2 := getSelfSettlement(room, fp2, 1, fp2ShowWinLossGold, fp2S < settlementGold)
			sPush.Settlement = append(sPush.Settlement, fs2)

			landShowWinLossGold := fmt.Sprintf("-%.2f", fp1S+fp2S)
			ls := getSelfSettlement(room, landPlayer, -1, landShowWinLossGold, (fp1S+fp2S) < settlementGold*2)
			sPush.Settlement = append(sPush.Settlement, ls)
		}

	}

	var multiInfo mproto.MultipleInfo
	multiInfo.FightLandlord = fmt.Sprintf("×%d", room.MultiGetLandlord)
	multiInfo.Boom = fmt.Sprintf("×%d", room.MultiBoom)
	multiInfo.Spring = fmt.Sprintf("×%d", room.MultiSpring)
	multiInfo.Rocket = fmt.Sprintf("×%d", room.MultiRocket)
	sPush.MultipleInfo = &multiInfo
	sPush.WaitTime = sysSet.GameDelayReadyTimeInt

	// 更新就
	go DBUptRecode(room, sPush)

	logger.Debug("结算信息:", sPush)
	logger.Debug(fmt.Println(sPush))
	bytes, _ := proto.Marshal(&sPush)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.PushSettlement, bytes))

}

func syncWinGold(player *Player, gold, goldPay float64, roundId string, roomType RoomClassify) float64 {
	//orderId := fmt.Sprintf("%s-%s-win", roundId, player.PlayerInfo.PlayerId)
	player.PlayerInfo.Gold = player.PlayerInfo.Gold + goldPay // 同步到房间id

	if !player.IsRobot { // 如果不是机器人则同步session信息
		err := SetSessionGold(player.Session, goldPay) // 同步到session
		if err != nil {
			logger.Error("同步进步到session失败: !!!incredible")
		}
		UserSyncWinScore(player.PlayerInfo.PlayerId, gold, roundId) // 同步到中心服务

		// 赢钱超过设定值发送 跑马灯
		if !player.IsRobot && goldPay > Server.WinGoldNotice {
			NoticeWinMoreThan(player.PlayerInfo.PlayerId, player.PlayerInfo.Name, goldPay)
		}
	}

	if !player.IsRobot {
		// 更新盈余池
		var surplus SurplusPool
		surplus.RoomType = roomType
		surplus.CurrentPlayerWin = gold
		surplus.InsertSurplus()
	}

	return player.PlayerInfo.Gold
}

func syncLossGold(player *Player, gold float64, roundId string, roomType RoomClassify) float64 {
	//orderId := fmt.Sprintf("%s-%s-loss", roundId, player.PlayerInfo.PlayerId)
	player.PlayerInfo.Gold = player.PlayerInfo.Gold - gold
	if !player.IsRobot { // 如果不是机器人则同步session信息
		err := SetSessionGold(player.Session, -gold) // 同步到session
		if err != nil {
			logger.Error("同步进步到session失败: !!!incredible")
		}
		UserSyncLoseScore(player.PlayerInfo.PlayerId, -gold, roundId)
	}

	if !player.IsRobot {
		// 更新盈余池
		var surplus SurplusPool
		surplus.RoomType = roomType
		surplus.CurrentPlayerLoss = gold
		surplus.InsertSurplus()
	}
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
