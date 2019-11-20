package roomStatus

const (
	// 房间状态 0 等待中 1叫地主 2.抢地主  3.正在玩
	Wait = iota
	CallLandlord
	GetLandlord
	Playing
)
