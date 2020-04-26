package game

import (
	"strconv"
	"landlord/mconst/roomType"
)

func CreateRobot(roomType int32) *Player {
	var player Player
	var pi PlayerInfo
	pi.PlayerId = strconv.Itoa(RandNum(123115480, 956855780))
	pi.Gold = getRobotGold(roomType)
	pi.Name = pi.PlayerId
	pi.HeadImg = getRobotImg()
	player.PlayerInfo = &pi
	player.IsRobot = true
	//player.ActionChan = make(chan PlayerActionChan)
	return &player
}

func getRobotImg() string {
	img := []string{"1.png", "2.png", "3.png", "4.png", "5.png", "6.png", "7.png", "8.png", "9.png", "10.png",
		"11.png", "12.png", "13.png", "14.png", "15.png", "16.png", "17.png", "18.png", "19.png", "20.png"}
	return img[RandNum(0, len(img)-1)]
}


/*
	ExperienceField int32 = 1 // 体验场
	LowField        int32 = 2 // 初级场
	MidField        int32 = 3 // 中级场
	HighField       int32 = 4 // 高级场
*/
func getRobotGold(rt int32)float64{
	switch rt {
	case roomType.ExperienceField:
		return float64(RandNum(6, 30))+float64(RandNum(6, 30))/100
	case roomType.LowField:
		return float64(RandNum(20, 200))+float64(RandNum(6, 30))/100
	case roomType.MidField:
		return float64(RandNum(30, 300))+float64(RandNum(6, 30))/100
	case roomType.HighField:
		return float64(RandNum(101, 300))+float64(RandNum(6, 30))/100
	}
	return float64(RandNum(101, 300))+float64(RandNum(6, 30))/100
}

