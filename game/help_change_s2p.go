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

func ChangeArrPlayerToGamePlayer(items []*Player) []*mproto.GamePlayer {

	//result.IsRobot = item.IsRobot
	var result []*mproto.GamePlayer
	for i := 0; i < len(items); i++ {
		item := items[i]
		var gPlayer mproto.GamePlayer
		gPlayer.PlayerInfo = ChangePlayerInfoToProto(item.PlayerInfo)
		gPlayer.HandMahjongs = item.HandMahjongs
		gPlayer.DoubleMahjongs = item.PengMahjongs
		gPlayer.EatMahjongs = item.ChiMahjongs
		gPlayer.ThreeMahjongs = item.GangMahjongs
		gPlayer.DeepThreeMahjongs = item.AnGangMahjongs
		gPlayer.IsVillage = item.IsVillage
		gPlayer.IsListen = item.IsListen
		result = append(result, &gPlayer)
	}
	return result
}

func ChangeArrPlayerToRoomPlayer(items []*Player) []*mproto.RoomPlayerInfo {

	//result.IsRobot = item.IsRobot
	var result []*mproto.RoomPlayerInfo
	for i := 0; i < len(items); i++ {
		item := items[i]
		var rPlayer mproto.RoomPlayerInfo
		rPlayer.PlayerInfo = ChangePlayerInfoToProto(item.PlayerInfo)
		rPlayer.WindDirection = item.Direction
		result = append(result, &rPlayer)
	}
	return result
}

func ChangePlayerToGamePlayer(item *Player) *mproto.GamePlayer {
	var result mproto.GamePlayer
	result.PlayerInfo = ChangePlayerInfoToProto(item.PlayerInfo)
	//result.IsRobot = item.IsRobot
	result.HandMahjongs = item.HandMahjongs
	result.DoubleMahjongs = item.PengMahjongs
	result.EatMahjongs = item.ChiMahjongs
	result.ThreeMahjongs = item.GangMahjongs
	result.DeepThreeMahjongs = item.AnGangMahjongs
	result.IsVillage = item.IsVillage
	result.IsListen = item.IsListen
	return &result
}

func ChangePlayerInfoToProto(item *PlayerInfo) *mproto.PlayerInfo {

	var result mproto.PlayerInfo
	result.Gold = item.Gold
	result.PlayerName = item.Name
	result.PlayerId = item.PlayerId
	result.PlayerImg = item.HeadImg

	return &result
}
