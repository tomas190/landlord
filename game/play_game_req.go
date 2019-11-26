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

	roomId := GetSessionRoomId(session)
	if roomId != "" {
		room := GetRoom(roomId)
		if room.RoomClass.RoomType != req.RoomType { // 如果跟请求的type 不一样则推送原有房间type
			//todo  用户waitTime 和上一个操作 上一个牌 待处理
			PushRecoverRoom(session, room, playerInfo.PlayerId)
		}
		return
	}

	switch req.RoomType {
	case roomType.ExperienceField: // 如果是体验场
		DealPlayerEnterExpField(session, *playerInfo)
	case roomType.LowField:
		DealPlayerEnterLowField(session, *playerInfo)
	case roomType.MidField:
		DealPlayerEnterMidField(session, *playerInfo)
	case roomType.HighField:
		DealPlayerEnterHighField(session, *playerInfo)
	default:
		logger.Error("进入房间失败:无此房间类型", req.RoomType)
	}

} // 進入房間

func ReqEnterRoomCheck(session *melody.Session, data []byte) {
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

	roomId := GetSessionRoomId(session)
	if roomId != "" {
		room := GetRoom(roomId)
		if room.RoomClass.RoomType != req.RoomType { // 如果跟请求的type 不一样则推送原有房间type
			//todo  用户waitTime 和上一个操作 上一个牌 待处理
			RespEnterRoomCheck(session, room.RoomClass.RoomType)
		} else {
			RespEnterRoomCheck(session, req.RoomType)
		}
	} else {
		RespEnterRoomCheck(session, req.RoomType)
	}

}

// 进入房间返回
func RespEnterRoomCheck(session *melody.Session, roomType int32) {
	var resp mproto.RespEnterRoomCheck
	resp.RoomType = roomType
	bytes, _ := proto.Marshal(&resp)
	_ = session.WriteBinary(PkgMsg(msgIdConst.RespEnterRoomCheck, bytes))
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
		logger.Error("ReqOutCardDo:出牌错误", err.Error())
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

	if roomId!=req.RoomId { // 如果请求的roomId 和 自己的roomId 不一样 ze
		SendErrMsg(session, msgIdConst.ReqExitRoom, "roomId不一致！")
		return
	}

	// 2. 如果玩家在游戏中 则设置退出房间标记
	room := GetRoom(roomId)
	if room == nil {
		logger.Error("ReqOutCardDo:无room信息", roomId)
		SendErrMsg(session, msgIdConst.ReqExitRoom, "无room信息:"+roomId)
		return
	}
	// 设置退出房间标记
	// 并设置托管操作
	player := room.Players[info.PlayerId]
	if player != nil {
		player.IsExitRoom = true
		player.IsGameHosting = true
	} else {
		logger.Error("该房间无玩家信息 !!!incredible")
	}

} // 退出房间

// 托管
func ReqGameHosting(session *melody.Session, data []byte) {
	logger.Debug("=== ReqGameHosting ===")
	req := &mproto.ReqGameHosting{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		SendErrMsg(session, msgIdConst.ReqGameHosting, "请求数据异常:"+err.Error())
		return
	}

	info, err := GetSessionPlayerInfo(session)
	if err != nil {
		logger.Error("ReqGameHosting:此session无用户信息", info)
		SendErrMsg(session, msgIdConst.ReqGameHosting, "无用户信息:"+err.Error())
		return
	}

	PrintMsg("ReqGameHosting:"+info.PlayerId, req)
	/*==== 参数验证 =====*/

	roomId := GetSessionRoomId(session)
	if roomId == "" {
		logger.Debug(info.PlayerId, "ReqGameHosting 玩家不在房间")
		SendErrMsg(session, msgIdConst.ReqGameHosting, "托管失败:玩家不在房间中...")
		return
	}

	room := GetRoom(roomId)
	if room == nil {
		logger.Error("ReqGameHosting:无room信息", roomId)
		SendErrMsg(session, msgIdConst.ReqGameHosting, "无room信息:"+roomId)
		return
	}

	// todo 如果玩家在自己出牌阶段没出牌 点击了托管 则根据是否必出 进行托管托管逻辑出牌


	// 设置退出房间标记
	// 并设置托管操作
	player := room.Players[info.PlayerId]
	if player != nil {
		if req.GameHostType == 1 {
			// 托管
			player.IsGameHosting = true
		} else if req.GameHostType == -1 {
			// 取消托管
			player.IsGameHosting = false
		}

		//var resp mproto.RespGameHosting
		//resp.GameHostType = req.GameHostType
		//resp.PlayerId = player.PlayerInfo.PlayerId
		//resp.Position = player.PlayerPosition
		//bytes, _ := proto.Marshal(&resp) // 广播给房间的人
		//MapPlayersSendMsg(room.Players, bytes)
		// _ = session.WriteBinary(PkgMsg(msgIdConst.RespGameHosting, bytes))
		RespGameHosting(room, req.GameHostType, player.PlayerPosition, player.PlayerInfo.PlayerId)

	} else {
		logger.Error("ReqGameHosting 该房间无玩家信息 !!!incredible")
	}

}

func RespGameHosting(room *Room, ghType, position int32, PlayerId string) {
	var resp mproto.RespGameHosting
	resp.GameHostType = ghType
	resp.PlayerId = PlayerId
	resp.Position = position
	bytes, _ := proto.Marshal(&resp) // 广播给房间的人
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.RespGameHosting, bytes))
}

// 发送消息
func ReqSendMsg(session *melody.Session, data []byte) {
	logger.Debug("=== ReqSendMsg ===")
	req := &mproto.ReqSendMsg{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		SendErrMsg(session, msgIdConst.ReqSendMsg, "请求数据异常:"+err.Error())
		return
	}

	info, err := GetSessionPlayerInfo(session)
	if err != nil {
		logger.Error("ReqSendMsg:此session无用户信息", info)
		SendErrMsg(session, msgIdConst.ReqSendMsg, "无用户信息:"+err.Error())
		return
	}

	PrintMsg("ReqSendMsg:"+info.PlayerId, req)
	/*==== 参数验证 =====*/

	roomId := GetSessionRoomId(session)
	if roomId == "" {
		logger.Debug(info.PlayerId, "ReqSendMsg 玩家不在房间")
		SendErrMsg(session, msgIdConst.ReqSendMsg, "发送消息失败:玩家不在房间中...")
		return
	}

	room := GetRoom(roomId)
	if room == nil {
		logger.Error("ReqSendMsg:无room信息", roomId)
		SendErrMsg(session, msgIdConst.ReqGameHosting, "无room信息:"+roomId)
		return
	}

	var resp mproto.RespSendMsg
	resp.Msg = req.Msg
	bytes, _ := proto.Marshal(&resp)
	MapPlayersSendMsg(room.Players, PkgMsg(msgIdConst.RespSendMsg, bytes))

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
		return errors.New("你手中没有这样的牌")
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
