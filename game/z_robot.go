package game

import "strconv"

func CreateRobot() *Player {
	var player Player
	var pi PlayerInfo
	pi.PlayerId = strconv.Itoa(RandNum(12311548, 95685578))
	pi.Gold = float64(RandNum(300, 1983))
	pi.Name = pi.PlayerId
	pi.HeadImg = getRobotImg()
	player.PlayerInfo = &pi
	player.IsRobot = true
	player.ActionChan = make(chan PlayerActionChan)
	return &player
}

func getRobotImg() string {
	img := []string{"1.png", "2.png", "3.png", "4.png", "5.png", "6.png", "7.png", "8.png", "9.png", "10.png",
		"11.png", "12.png", "13.png", "14.png", "15.png", "16.png", "17.png", "18.png", "19.png", "20.png"}
	return img[RandNum(0, len(img)-1)]
}
