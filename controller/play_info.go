package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2/bson"
	"landlord/game"
	"strconv"
)

/*
	start_time
	end_time
	user_id
	room_type
*/

type PlayInfoReq struct {
	Id        string `form:"id" json:"id"`
	StartTime int64  `form:"start_time" json:"start_time"`
	EndTime   int64  `form:"end_time" json:"end_time"`
	RoomType  int32  `form:"room_type" json:"room_type"`
}

type PlayInfoResp struct {
	GameFlow float64 `json:"game_flow"`
	WinNum   int     `json:"win_num"`
	WinGold  float64 `json:"win_gold"`

	LoseNum  int     `json:"lose_num"`
	LoseGold float64 `json:"lose_gold"`
}

func GetPlayInfo(c *gin.Context) {
	logger.Debug(" ====== GetPlayInfo ======")
	var req PlayInfoReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(httpCode, NewResp(ErrCode, err.Error(), nil))
		return
	}

	logger.Debug("req:", req)
	if req.Id == "" {
		c.JSON(httpCode, NewResp(ErrCode, "id is nil", nil))
		return
	}
	if req.StartTime == 0 {
		c.JSON(httpCode, NewResp(ErrCode, "start_time is nil", nil))
		return
	}

	if req.EndTime == 0 {
		c.JSON(httpCode, NewResp(ErrCode, "end_time is nil", nil))
		return
	}

	//if req.RoomType < 0 || req.RoomType > 4 {
	//	c.JSON(httpCode, NewResp(ErrCode, "room_type is not between 1 to 4", nil))
	//	return
	//}

	data, err := GetPlayInfoHep(req)
	if err != nil {
		c.JSON(httpCode, NewResp(ErrCode, err.Error(), nil))
		return
	}
	c.JSON(httpCode, NewResp(SuccCode, "ok", data))
}

func GetPlayInfoHep(req PlayInfoReq) (PlayInfoResp, error) {
	selector := bson.M{}
	pattern := ".*" + req.Id + ".*"
	selector["player_ids"] = bson.M{"$regex": bson.RegEx{Pattern: pattern, Options: "im"}}
	selector["end_time"] = bson.M{"$gte": req.StartTime, "$lte": req.EndTime}
	if req.RoomType != 0 {
		selector["room_type"] = req.RoomType
	}

	var resp PlayInfoResp
	var item game.PlayCardRecode
	list, err := item.GetPlayInfoList(selector)
	if err != nil {
		return resp, err
	}

	for i := 0; i < len(list); i++ {
		pc := list[i]

		isWin, winLoseGold := isPlayerWin(req.Id, pc.Settlement)
		if isWin {
			resp.WinNum += 1
			resp.WinGold += winLoseGold
		} else {
			resp.LoseNum += 1
			resp.LoseGold += winLoseGold
		}
		resp.GameFlow += winLoseGold
	}
	resp.GameFlow = Decimal2(resp.GameFlow)
	resp.WinGold = Decimal2(resp.WinGold)
	resp.LoseGold = Decimal2(resp.LoseGold)
	return resp, nil
}

func isPlayerWin(id string, info []game.SettlementInfo) (bool, float64) {
	for i := 0; i < len(info); i++ {
		if info[i].PlayerId == id {
			winLoseGold, _ := strconv.ParseFloat(info[i].WinLossGold, 64)
			if info[i].WinOrFail == 1 {
				// 税前
				return true, winLoseGold / 0.95
			} else {
				return false, -winLoseGold
			}
		}
	}
	logger.Debug("id:", id)
	logger.Debug("info:", info)
	return false, 0
}

func Decimal2(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
