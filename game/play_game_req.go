package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
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

	switch req.RoomType {
	case roomType.ExperienceField: // 如果是体验场
		DealPlayerEnterExpField(session, *playerInfo)
	case roomType.LowField:
	// todo
	case roomType.MidField:
	case roomType.HighField:

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

	var actionChan PlayerActionChan
	if len(req.Cards) <= 0 {
		actionChan.ActionType = playerAction.NotOutCardAction
	} else {
		actionChan.ActionType = playerAction.OutCardAction
		actionChan.ActionCards = ChangeProtoToCard(req.Cards)
		// todo 获取房间牌型和有效牌 并判断是否可以出此牌 如果否则推送错误消息
		// todo 如果出牌是炸弹 则 房间翻倍

	}
	actionPlayer.ActionChan <- actionChan

}
