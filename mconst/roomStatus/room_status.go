package roomStatus

const (
	// 房间状态 0 等待中 1叫地主  2正在玩
	Wait = iota
	GetLandlord
	Playing
)
