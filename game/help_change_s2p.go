package game

import "landlord/msg/mproto"

/*
	struct to proto
*/

func ChangeRoomClassifyToProto(item *RoomClassify) *mproto.RoomClassify {

	var result mproto.RoomClassify
	result.RoomType = item.RoomType
	result.BottomEnterPoint = item.BottomEnterPoint
	result.BottomPoint = item.BottomPoint

	return &result
}

func ChangePlayerP2S(item mproto.PlayerInfo) PlayerInfo {

	return PlayerInfo{
		item.PlayerId,
		item.PlayerName,
		item.PlayerImg,
		item.Gold,
	}
}

func ChangePlayerInfoToProto(item *PlayerInfo) *mproto.PlayerInfo {

	var result mproto.PlayerInfo
	result.Gold = item.Gold
	result.PlayerName = item.Name
	result.PlayerId = item.PlayerId
	result.PlayerImg = item.HeadImg

	return &result
}
