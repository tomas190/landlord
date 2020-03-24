package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"landlord/game"
	"time"
)

type GameDataReq struct {
	GameId    string `form:"game_id" json:"game_id"`
	Id        string `form:"id" json:"id"`
	StartTime int64  `form:"start_time" json:"start_time"`
	EndTime   int64  `form:"end_time" json:"end_time"`
	RoundId   string `form:"round_id" json:"round_id"`
	Skip      int    `form:"skip" json:"skip"`
	Limit     int    `form:"limit" json:"limit"`
}

type GameData struct {
	StartTime    int64  `json:"start_time"`
	StartTimeFmt string `json:"start_time_fmt"`
	EndTime      int64  `json:"end_time"`
	EndTimeFmt   string `json:"end_time_fmt"`

	PlayerId   string      `json:"player_id"`
	RoundId    string      `json:"round_id"`
	RoomId     string      `json:"room_id"`
	TaxRate    float64     `json:"tax_rate"`
	PlayerInfo interface{} `json:"player_info"`
	Settlement interface{} `json:"settlement"` // 结算信息 输赢结果
}

func GetBaccaratData(c *gin.Context) {
	var req GameDataReq

	err := c.Bind(&req)
	if err != nil {
		c.JSON(httpCode, NewResp(ErrCode, err.Error(), nil))
		return
	}

	data, err := HelpGetBaccaratData(req)
	if err != nil {
		c.JSON(httpCode, NewResp(ErrCode, err.Error(), nil))
		return
	}
	c.JSON(httpCode, NewResp(SuccCode, "ok", data))
}

type pageData struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

func HelpGetBaccaratData(req GameDataReq) (*pageData, error) {

	// verify param
	if req.GameId != game.Server.GameId {
		return nil, errors.New("auth fail")
	}

	selector := bson.M{}
	if req.Id == "" {
		return nil, errors.New("err params")
	}

	playerId := req.Id
	roundId := req.RoundId
	startTime := req.StartTime
	endTime := req.EndTime
	skip := req.Skip
	limit := req.Limit

	//selector["player_id"] = playerId

	pattern := ".*" + playerId + ".*"
	selector["player_ids"] = bson.M{"$regex": bson.RegEx{Pattern: pattern, Options: "im"}}

	if roundId != "" {
		selector["rand_id"] = roundId
	}

	if startTime != 0 && endTime != 0 {
		selector["start_time"] = bson.M{"$gte": startTime, "$lte": endTime}
	}

	if startTime != 0 && endTime == 0 {
		selector["start_time"] = bson.M{"$gt": startTime}
	}

	if endTime != 0 && startTime == 0 {
		selector["start_time"] = bson.M{"$lt": endTime}
	}

	var item game.PlayCardRecode
	recodes, count, err := item.GetPlayCardRecodeList(skip, limit, selector, "down_bet_time")
	if err != nil {
		return nil, err
	}

	var gameDatas []GameData
	for i := 0; i < len(recodes); i++ {
		var gd GameData
		pr := recodes[i]
		gd.StartTime = pr.StartTime * 1000
		gd.StartTimeFmt = FormatTime(pr.StartTime, "2006-01-02 15:04:05")
		gd.PlayerId = playerId
		gd.RoomId = pr.RoomId
		gd.RoundId = pr.RoundId
		gd.PlayerInfo = pr.Players
		gd.Settlement = pr.Settlement
		gd.TaxRate = pr.GameTaxRate
		gd.EndTime = pr.EndTime
		gd.EndTimeFmt = FormatTime(pr.EndTime, "2006-01-02 15:04:05")
		gameDatas = append(gameDatas, gd)

	}

	var result pageData
	result.Total = count
	result.List = gameDatas
	return &result, nil

}

func FormatTime(timeUnix int64, layout string) string {
	if timeUnix == 0 {
		return ""
	}
	format := time.Unix(timeUnix, 0).Format(layout)
	return format
}


func Version(c *gin.Context){
	c.JSON(httpCode, NewResp(SuccCode, "version 2020324", "2020年3月24日16:05:07"))
}