package game

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
	"landlord/mconst/cardConst"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/playerAction"
	"landlord/mconst/roomType"
	"landlord/msg/mproto"
)

// 進入房間
func ReqEnterRoom(session *melody.Session, data []byte) {
	logger.Debug("=== ReqEnterRoom ===")
	req := &mproto.ReqEnterRoom{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		SendErrMsg(session, msgIdConst.ReqEnterRoom, "请求数据异常:"+err.Error())
		return
	}

	PrintMsg("ReqEnterRoom:", req)
	/*==== 参数验证 =====*/

	playerInfo, err := GetSessionPlayerInfo(session)
	if err != nil {
		SendErrMsg(session, msgIdConst.ReqEnterRoom, "无用户信息")
		return
	}

	if playerInfo.Gold < GetRoomClassifyBottomEnterPoint(req.RoomType) {
		SendErrMsg(session, msgIdConst.ReqEnterRoom, "金币不足!")
		return
	}

	// todo 进入房间之前 判断用户是否在房间中...
	if GetSessionRoomId(session) != "" {
		// todo 恢复房间信息

	}

	switch req.RoomType {
	case roomType.ExperienceField: // 如果是体验场
		DealPlayerEnterExpField(session, *playerInfo)
	case roomType.LowField:
		DealPlayerEnterLowField(session, *playerInfo)
	case roomType.MidField:
		DealPlayerEnterLowField(session, *playerInfo)
	case roomType.HighField:
		DealPlayerEnterLowField(session, *playerInfo)
	default:
		logger.Error("进入房间失败:无此房间类型", req.RoomType)
	}

}

// 抢地主操作
func ReqGetLandlordDo(session *melody.Session, data []byte) {
	logger.Debug("=== ReqGetLandlordDo ===")
	req := &mproto.ReqGetLandlordDo{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		SendErrMsg(session, msgIdConst.ReqGetLandlordDo, "请求数据异常:"+err.Error())
		return
	}

	info, err := GetSessionPlayerInfo(session)
	if err != nil {
		logger.Error("ReqDoAction:此session无用户信息", info)
		SendErrMsg(session, msgIdConst.ReqGetLandlordDo, "无用户信息:"+err.Error())
		return
	}

	PrintMsg("ReqGetLandlordDo:"+info.PlayerId, req)
	/*==== 参数验证 =====*/

	roomId := GetSessionRoomId(session)
	room := GetRoom(roomId)
	if room == nil {
		logger.Error("ReqDoAction:无room信息", roomId)
		SendErrMsg(session, msgIdConst.ReqGetLandlordDo, "无room信息:"+roomId)
		return
	}

	actionPlayer := room.Players[info.PlayerId]
	if actionPlayer == nil {
		logger.Error("ReqDoAction:无room信息", roomId)
		SendErrMsg(session, msgIdConst.ReqGetLandlordDo, "room无用户信息:"+roomId+"--"+info.PlayerId)
		return
	}

	var actionChan PlayerActionChan
	actionChan.ActionType = req.Action
	actionPlayer.ActionChan <- actionChan

} // 抢地主操作

// 出牌打牌操作
func ReqOutCardDo(session *melody.Session, data []byte) {
	logger.Debug("=== ReqOutCardDo ===")
	req := &mproto.ReqOutCardDo{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		SendErrMsg(session, msgIdConst.ReqOutCardDo, "请求数据异常:"+err.Error())
		return
	}

	info, err := GetSessionPlayerInfo(session)
	if err != nil {
		logger.Error("ReqOutCardDo:此session无用户信息", info)
		SendErrMsg(session, msgIdConst.ReqOutCardDo, "无用户信息:"+err.Error())
		return
	}

	PrintMsg("ReqOutCardDo:"+info.PlayerId, req)
	/*==== 参数验证 =====*/

	roomId := GetSessionRoomId(session)
	room := GetRoom(roomId)
	if room == nil {
		logger.Error("ReqOutCardDo:无room信息", roomId)
		SendErrMsg(session, msgIdConst.ReqOutCardDo, "无room信息:"+roomId)
		return
	}

	actionPlayer := room.Players[info.PlayerId]
	if actionPlayer == nil {
		logger.Error("ReqOutCardDo:无room信息", roomId)
		SendErrMsg(session, msgIdConst.ReqOutCardDo, "room无用户信息:"+roomId+"--"+info.PlayerId)
		return
	}

	outCards := ChangeProtoToCard(req.Cards)

	cardType, err := verifyOutCard(room, actionPlayer, outCards)
	if err != nil {
		logger.Error("ReqOutCardDo:出牌错误", roomId)
		SendErrMsg(session, msgIdConst.ReqOutCardDo, err.Error())
		return
	}

	var actionChan PlayerActionChan
	if len(req.Cards) <= 0 {
		actionChan.ActionType = playerAction.NotOutCardAction
	} else {
		actionChan.ActionCards = ChangeProtoToCard(req.Cards)
		actionChan.ActionType = playerAction.OutCardAction
		actionChan.CardsType = cardType
	}
	actionPlayer.ActionChan <- actionChan

}

// 退出房间
func ReqExitRoom(session *melody.Session, data []byte) {
	logger.Debug("=== ReqExitRoom ===")
	req := &mproto.ReqExitRoom{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		SendErrMsg(session, msgIdConst.ReqExitRoom, "请求数据异常:"+err.Error())
		return
	}

	info, err := GetSessionPlayerInfo(session)
	if err != nil {
		logger.Error("ReqExitRoom:此session无用户信息", info)
		SendErrMsg(session, msgIdConst.ReqExitRoom, "无用户信息:"+err.Error())
		return
	}

	PrintMsg("ReqExitRoom:"+info.PlayerId, req)
	/*==== 参数验证 =====*/

	roomId := GetSessionRoomId(session)
	// 1. 如果roomId为空代表玩家是在等待队列 则移除等待队列
	if roomId == "" {
		logger.Debug(info.PlayerId, "当前在等待队列中..")
		RemoveWaitUser(info.PlayerId)
		return
	}

	// 2. 如果玩家在游戏中 则设置退出房间标记
	room := GetRoom(roomId)
	if room == nil {
		logger.Error("ReqOutCardDo:无room信息", roomId)
		SendErrMsg(session, msgIdConst.ReqExitRoom, "无room信息:"+roomId)
		return
	}
	player := room.Players[info.PlayerId]
	if player != nil {
		player.IsExitRoom = true
	} else {
		logger.Error("改房间无玩家信息 !!!incredible")
	}

}

/*=================== help func ===================*/
// 检测出牌是否合理
func verifyOutCard(room *Room, actionPlayer *Player, outCards []*Card) (int32, error) {
	if len(outCards) <= 0 && !actionPlayer.IsMustDo {
		return room.EffectiveType, nil
	}

	// 1. 判断是否该这个玩家出牌
	if !actionPlayer.IsCanDo {
		return 0, errors.New("当前不该你出牌")
	}
	cardType := GetCardsType(outCards)
	// 2.检测牌型是否正确
	if cardType > cardConst.CARD_PATTERN_QUADPLEX_WITH_PAIRS || cardType < cardConst.CARD_PATTERN_SINGLE {
		return 0, errors.New("你出牌不符合规则")
	}

	// 2. 必须出牌检测
	if actionPlayer.IsMustDo {
		err := verifyMustOutCard(actionPlayer, outCards, cardType)
		if err != nil {
			return 0, err
		}
	} else {
		// 3. 跟牌检测
		err := verifyFollowOutCard(room, actionPlayer, outCards)
		if err != nil {
			return 0, err
		}
	}

	// 14火箭
	if cardType == cardConst.CARD_PATTERN_BOMB || cardType == cardConst.CARD_PATTERN_ROCKET {
		room.MultiAll = room.MultiAll * 2
		if room.MultiBoom == 0 {
			room.MultiBoom = 2
		} else {
			room.MultiBoom = room.MultiBoom * 2
		}

	}

	room.EffectiveCard = outCards
	room.EffectiveType = cardType

	return cardType, nil
}

// 必须出牌检测
func verifyMustOutCard(actionPlayer *Player, cards []*Card, cardType int32) error {
	// 1.判断玩家的手牌中是否存在这样的牌
	exist := checkCardsIsExist(actionPlayer.HandCards, cards)
	if !exist {
		return errors.New("===你手中没有这样的牌===")
	}

	return nil
}

// 跟牌出牌检测
func verifyFollowOutCard(room *Room, actionPlayer *Player, cards []*Card) error {
	// 检测是否打过上家
	can := CanBeat(room.EffectiveCard, cards)
	if !can {
		return errors.New("出牌不能打过上家")
	}
	return nil
}

// 判断玩家手牌中是否存在这些牌
func checkCardsIsExist(handCards []*Card, cards []*Card) bool {
	if len(cards) > len(handCards) {
		return false
	}
	for i := 0; i < len(cards); i++ {
		var flag bool
		for j := 0; j < len(handCards); j++ {
			if cards[i].Value == handCards[j].Value && cards[i].Suit == handCards[j].Suit {

				flag = true
				break
			}
		}
		if !flag {
			logger.Debug("玩家手牌:")
			PrintCard(handCards)
			logger.Debug("玩家出牌:")
			PrintCard(cards)
			logger.Debug("不存在这张牌:", cards[i])
			return false
		}
	}
	return true
}
