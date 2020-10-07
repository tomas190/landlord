package controller

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"landlord/game"
	"landlord/mconst/roomType"
	"os"
	"strconv"
	"time"
)

type GameDataReq struct {
	GameId    string `form:"game_id" json:"game_id"`
	Id        string `form:"id" json:"id"`
	StartTime int64  `form:"start_time" json:"start_time"`
	EndTime   int64  `form:"end_time" json:"end_time"`
	RoundId   string `form:"round_id" json:"round_id"`
	Page      int    `form:"page" json:"page"`
	Limit     int    `form:"limit" json:"limit"`
	Roundres  int    `form:"roundres" json:"roundres"`
	RoomId    string `form:"room_id" json:"room_id"`
	//Skip      int    `form:"skip" json:"skip"`
}

type GameData struct {
	StartTime    int64  `json:"start_time"`
	StartTimeFmt string `json:"start_time_fmt"`
	EndTime      int64  `json:"end_time"`
	EndTimeFmt   string `json:"end_time_fmt"`

	PlayerId   string      `json:"player_id"`
	RoundId    string      `json:"round_id"`
	RoomId     string      `json:"room_id"`
	RoomType   int32       `json:"room_type"`
	BottomCard interface{} `json:"bottom_card"`
	TaxRate    float64     `json:"tax_rate"`
	PlayerInfo interface{} `json:"player_info"`
	Settlement interface{} `json:"settlement"` // 结算信息 输赢结果

	CreatedAt       int64       `json:"created_at"` // 创建时间
	RoomName        string      `json:"room_name"`
	SettlementFunds float64     `json:"settlement_funds"` // 结算信资金 税后
	Card            interface{} `json:"card"`
	SpareCash       float64     `json:"spare_cash"` // 玩家剩余资金
}

func GetLandlordData(c *gin.Context) {
	var req GameDataReq

	err := c.Bind(&req)
	if err != nil {
		c.JSON(httpCode, NewResp(ErrCode, err.Error(), nil))
		return
	}

	data, winCount, err := HelpGetLandlordData(req)
	if err != nil {
		c.JSON(httpCode, NewResp(ErrCode, err.Error(), nil))
		return
	}

	if req.Roundres == 1 {
		r := simplejson.New()
		r.Set("total", data.Total)
		r.Set("winround", winCount)
		c.JSON(httpCode, NewResp(SuccCode, "ok", r))
		return
	}
	c.JSON(httpCode, NewResp(SuccCode, "ok", data))
}

type pageData struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

func HelpGetLandlordData(req GameDataReq) (*pageData, int, error) {

	// verify param
	if req.GameId != game.Server.GameId {
		return nil, 0, errors.New("auth fail")
	}

	selector := bson.M{}
	//if req.Id == "" {
	//	return nil, 0, errors.New("err params")
	//}

	playerId := req.Id
	roundId := req.RoundId
	startTime := req.StartTime
	endTime := req.EndTime
	page := req.Page
	limit := req.Limit
	roomId:=req.RoomId
	if page == 0 {
		page = 1
	}
	skip := limit * (page - 1)
	//selector["player_id"] = playerId

	if playerId!= "" {
		pattern := ".*" + playerId + ".*"
		selector["player_ids"] = bson.M{"$regex": bson.RegEx{Pattern: pattern, Options: "im"}}
	}
	if roundId != "" {
		selector["round_id"] = roundId
	}
	if roomId != "" {
		//selector["room_id"] = roomId
		ri, err := strconv.Atoi(roomId)
		if err!=nil {
			return nil, 0, errors.New("room_id不正确")
		}
		selector["room_type"] = int32(ri)
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
	//recodes, count, winCount, err := item.GetPlayCardRecodeList(skip, limit, selector, "-down_bet_time", req.Roundres, playerId)
	recodes, count, winCount, err := item.GetPlayCardRecodeList(skip, limit, selector, "-start_time", req.Roundres, playerId)
	if err != nil {
		return nil, 0, err
	}

	var gameDatas []GameData
	for i := 0; i < len(recodes); i++ {
		var gd GameData
		pr := recodes[i]
		gd.StartTime = pr.StartTime
		gd.StartTimeFmt = FormatTime(pr.StartTime, "2006-01-02 15:04:05")
		gd.PlayerId = playerId
		//gd.RoomId = pr.RoomId
		gd.RoomId = strconv.Itoa(int(pr.RoomType))
		gd.RoundId = pr.RoundId
		gd.RoomType = pr.RoomType
		gd.BottomCard = pr.BottomCard
		gd.PlayerInfo = pr.Players
		gd.Settlement = pr.Settlement
		gd.TaxRate = pr.GameTaxRate
		gd.EndTime = pr.EndTime
		gd.EndTimeFmt = FormatTime(pr.EndTime, "2006-01-02 15:04:05")

		gd.CreatedAt = pr.StartTime
		settlementFunds, spareCash := GetSettlementFundsAndSpare(pr.Settlement, playerId)
		gd.SpareCash = spareCash // todo
		gd.SettlementFunds = settlementFunds
		gd.RoomName = GetRoomName(pr.RoomType)
		gd.Card = pr.Players
		gameDatas = append(gameDatas, gd)
	}

	var result pageData
	result.Total = count
	result.List = gameDatas
	return &result, winCount, nil

}

func GetSettlementFundsAndSpare(ss []game.SettlementInfo, playerId string) (float64, float64) {
	var res, res2 float64
	for i := 0; i < len(ss); i++ {
		if ss[i].PlayerId == playerId {
			f, _ := strconv.ParseFloat(ss[i].WinLossGold, 64)
			f2, _ := ss[i].CurrentGold, 64
			res = f
			res2 = f2
			break
		}
	}
	return res, res2
}

func GetRoomName(rType int32) string {
	var roomName string
	switch rType {
	case roomType.ExperienceField:
		roomName = "体验场"
	case roomType.LowField:
		roomName = "低级场"
	case roomType.MidField:
		roomName = "中级场"
	case roomType.HighField:
		roomName = "高级场"
	}
	return roomName
}

func FormatTime(timeUnix int64, layout string) string {
	if timeUnix == 0 {
		return ""
	}
	format := time.Unix(timeUnix, 0).Format(layout)
	return format
}

func Version(c *gin.Context) {
	c.JSON(httpCode, NewResp(SuccCode, "version 2020324", "2020年3月24日16:05:07"))
}

func GetLog(c *gin.Context) {
	//http下载地址 csv
	path := "out.log"
	if checkFileIsExist(path) {
		c.File(path)
		return
	}
	c.JSON(httpCode, NewResp(ErrCode, "不存在", nil))
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist

}
