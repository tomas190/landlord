package game

import (
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/roomStatus"
	"landlord/mconst/roomType"
	"landlord/msg/mproto"
)

// 1.给玩家和机器人发牌
func PushPlayerStartGameWithRobot(room *Room) {
	cards := CreateBrokenCard()
	//cards := CreateSortCard()

	//player, r1, r2 := getPlayersWithRobot(room)
	//// todo 玩家发牌策略
	//player.HandCards = append([]*Card{}, cards[:17]...)
	//var push mproto.PushStartGame
	//push.Cards = ChangeCardToProto(player.HandCards)
	//bytes, _ := proto.Marshal(&push)
	//PlayerSendMsg(player, PkgMsg(msgIdConst.PushStartGame, bytes))
	//
	//r1.HandCards = append([]*Card{}, cards[17:34]...)
	//r2.HandCards = append([]*Card{}, cards[34:51]...)
	//room.BottomCards = append([]*Card{}, cards[51:]...)
	//logger.Debug("底牌:")
	//PrintCard(room.BottomCards)

	// 随机发牌
	for _, v := range room.Players {
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
	// 随机发牌

	// 随机叫地主写在发牌里面 是因为三个玩家如果都不叫 则可以直接调用 PushPlayerStartGameWithRobot 重新开始发牌逻辑
	DelaySomeTime(4)
	// 3.随机叫地主
	actionPlayerId := pushFirstCallLandlord(room)
	CallLandlord(room, actionPlayerId)
}

// 1.给玩家和机器人发牌
func PushPlayerStartGameWithRobot2(room *Room) {
	//cards := CreateBrokenCard()
	//p1 := append([]*Card{}, cards[:17]...)
	//p2 := append([]*Card{}, cards[17:34]...)
	//p3 := append([]*Card{}, cards[34:51]...)
	//bottomCard := append([]*Card{}, cards[51:]...)
	//p1,p2,p3,bottomCard := CreateCardsNew()

	//p3 := []*Card{
	//	{5, 1}, {5, 1}, {5, 1},{5, 1},
	//	{6, 1}, {6, 1},// {6, 1},
	//	{7, 1}, //{7, 2}, {7, 3},
	//	{8, 1}, {8, 1}, //{8, 1}, //{8, 1},
	//
	//	{11, 1},
	//	{12, 1}, {12, 1},{12, 1},
	//	{13, 2}, {13, 2},
	//	{15, 2}, {14, 2},
	//}
	//
	//p2 := []*Card{
	//	{5, 1}, {5, 1}, {5, 1},{5, 1},
	//	{6, 1}, {6, 1},// {6, 1},
	//	{7, 1}, //{7, 2}, {7, 3},
	//	{8, 1}, {8, 1}, //{8, 1}, //{8, 1},
	//
	//	{11, 1},
	//	{12, 1}, {12, 1},{12, 1},
	//	{13, 2}, {13, 2},
	//	{15, 2}, {14, 2},
	//	//{15, 2},
	//}
	//p1 := []*Card{
	//	{5, 1}, {5, 1}, {5, 1},{5, 1},
	//	{6, 1}, {6, 1},// {6, 1},
	//	{7, 1}, //{7, 2}, {7, 3},
	//	{8, 1}, {8, 1}, //{8, 1}, //{8, 1},
	//
	//	{11, 1},
	//	{12, 1}, {12, 1},{12, 1},
	//	{13, 2}, {13, 2},
	//	{15, 2}, {14, 2},
	//}
	//bottomCard := []*Card{
	//	{12, 2}, {11, 2}, {8, 2},
	//}

	var p1, p2, p3, bottomCard []*Card
	var s SurplusPool
	surplus := s.GetLastSurplus()
	logger.Debug("当前盈余池:", surplus.CurrentSurplus)
	if surplus.CurrentSurplus <= 500 {
		logger.Debug("盈余池小于0 发好牌")
		num := RandNum(0, 10)
		if num >= 5 {
			p3, p2, p1, bottomCard = CreateGodCards()
		} else {
			p2, p1, p3, bottomCard = CreateGodCards()
		}
	} else {
		// 2020年2月24日15:39:27 如果是体验场 玩家都会有好牌
		if room.RoomClass.RoomType == roomType.ExperienceField {
			logger.Debug("体验场都发好牌...")
			level := RandNum(35, 42)
			p1, p2, p3, bottomCard = CreateGoodCard(level)
		} else { // 反之则正常发牌
			p1, p2, p3, bottomCard = CreateCardsNew()
		}
	}

	player, r1, r2 := getPlayersWithRobot(room)
	// todo 玩家发牌策略
	player.HandCards = append([]*Card{}, p1...)
	var push mproto.PushStartGame
	push.Cards = ChangeCardToProto(player.HandCards)
	bytes, _ := proto.Marshal(&push)
	PlayerSendMsg(player, PkgMsg(msgIdConst.PushStartGame, bytes))

	r1.HandCards = append([]*Card{}, p2...)
	r2.HandCards = append([]*Card{}, p3...)
	room.BottomCards = append([]*Card{}, bottomCard...)

	logger.Debug("底牌:")
	PrintCard(room.BottomCards)

	PrintCard(bottomCard)
	room.Status = roomStatus.CallLandlord

	// 组
	CountRobotCardValue(r1, r2)

	// 随机叫地主写在发牌里面 是因为三个玩家如果都不叫 则可以直接调用 PushPlayerStartGameWithRobot 重新开始发牌逻辑
	DelaySomeTime(4)
	// 3.随机叫地主
	actionPlayerId := pushFirstCallLandlord(room)
	CallLandlord(room, actionPlayerId)
}

// 计算机器人手牌分数
func CountRobotCardValue(r1, r2 *Player) {
	v1 := CountCardValue(r1.HandCards)
	r1.HandsValue = v1

	v2 := CountCardValue(r2.HandCards)
	r2.HandsValue = v2

	// 将牌分组

	groupCard1 := GroupHandsCard(r1.HandCards)
	groupCard2 := GroupHandsCard(r2.HandCards)

	r1.GroupCard = groupCard1
	r2.GroupCard = groupCard2
}

// 第一首是天牌牌
func CreateGodCards() ([]*Card, []*Card, []*Card, []*Card) {

	godCard, rCards := getGodCard()
	logger.Debug("===========:", len(rCards))

	p1card := godCard
	p2card := append([]*Card{}, rCards[:13]...)
	p3card := append([]*Card{}, rCards[13:26]...)
	boCard := append([]*Card{}, rCards[26:]...)

	i, i2, i3 := stick()

	p1card = append(p1card, i...)
	p2card = append(p2card, i2...)
	p3card = append(p3card, i3...)

	return p1card, p2card, p3card, boCard

}

// 获取13张好牌
func getGodCard() ([]*Card, []*Card) {
	cards, _ := CreateBroken8910Card()
	OutOfCardNotDeep42(cards, 5)
	ci := getGodCardIndex()

	tmp := append([]*Card{}, cards...)

	var godCard []*Card
	for i := 0; i < len(ci); i++ {
		c := findThisValueCard(ci[i], tmp, 1)
		godCard = append(godCard, c...)
		tmp = removeCards(tmp, c)
	}

	return godCard, tmp
}

// 获取好牌的value
func getGodCardIndex() []int {
	arr := []int{15, 14, 13, 13, 13, 13, 12, 12, 12, 12, 11, 11, 11, 11, 10, 10, 10, 10, 9, 9, 5, 5}
	for i := len(arr) - 1; i > 0; i-- {
		num := RandNum(0, len(arr)-1)
		arr[i], arr[num] = arr[num], arr[i]
	}

	return arr[:13]
}
