package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2/bson"
	"landlord/game"
	"math"
	"strconv"
)

/*
start_time	是	int	开始时间
end_time	是	int	结束时间
id	否	int	玩家ID
package_id	是	int	品牌ID与玩家ID不会同时传的
*/
type StatementTotalReq struct {
	StartTime int `form:"start_time"`
	EndTime   int `form:"end_time"`
	Id        int `form:"id"`
	PackageId int `form:"package_id"`
}

/*
 "lose_statement_total": -18767,
 "win_statement_total": 1822971,
 "game_id": "5b1f3a3cb76a591e7f251730",
 "game_name": "百家乐",
 "count": [614005170,614005355],
 "bet_money"：500 ，
*/
type StatementTotalResp struct {
	LoseStatementTotal float64
	WinStatementTotal  float64
	GameId             string
	GameName           string
	Count              []string
	BetMoney           float64
}

func GetStatementTotal(c *gin.Context) {
	var req StatementTotalReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(httpCode, NewResp(ErrCode, err.Error(), nil))
		return
	}
	logger.Debug("req:", req)

	if req.StartTime == 0 {
		c.JSON(httpCode, NewResp(ErrCode, "start_time is empty", nil))
		return
	}

	if req.EndTime == 0 {
		c.JSON(httpCode, NewResp(ErrCode, "end_time is empty", nil))
		return
	}

	if req.Id == 0 && req.PackageId == 0 {
		c.JSON(httpCode, NewResp(ErrCode, "id and package_id are empty", nil))
		return
	}

	selector := bson.M{}
	selector["end_time"] = bson.M{"$gte": req.StartTime, "$lte": req.EndTime}

	// 统计玩家id数据
	if req.Id != 0 {
		pattern := ".*" + strconv.Itoa(req.Id) + ".*"
		selector["player_ids"] = bson.M{"$regex": bson.RegEx{Pattern: pattern, Options: "im"}}
		resp, err := getDataById(selector, strconv.Itoa(req.Id))
		if err != nil {
			c.JSON(httpCode, NewResp(ErrCode, err.Error(), nil))
			return
		}
		c.JSON(httpCode, NewResp(SuccCode, "ok", resp))
		return

	}

	// 统计pkg数据
	if req.PackageId != 0 {
		selector["settlement"] = bson.M{"$elemMatch": bson.M{"PlayerPkgId": req.PackageId}}
		resp, err := getDataByPkgId(selector, req.PackageId)
		if err != nil {
			c.JSON(httpCode, NewResp(ErrCode, err.Error(), nil))
			return
		}
		c.JSON(httpCode, NewResp(SuccCode, "ok", resp))
		return
	}

}

// 统计玩家id数据
func getDataById(selector bson.M, pId string) (StatementTotalResp, error) {
	resp := NewStatementResp()
	var item game.PlayCardRecode
	recodes, err := item.GetPlayInfoList(selector)
	if err != nil {
		return resp, err
	}

	for i := 0; i < len(recodes); i++ {
		obj := recodes[i]
		winLoseMoney := getPlayerStatementInfo(obj.Settlement, pId)
		resp.BetMoney += math.Abs(winLoseMoney)
		resp.Count = append(resp.Count, pId)
		if winLoseMoney > 0 {
			resp.WinStatementTotal += winLoseMoney
		} else {
			resp.LoseStatementTotal += winLoseMoney
		}

	}

	return resp, nil
}

// 统计pkg数据
func getDataByPkgId(selector bson.M, pkdId int) (StatementTotalResp, error) {
	resp := NewStatementResp()
	var item game.PlayCardRecode
	recodes, err := item.GetPlayInfoList(selector)
	if err != nil {
		return resp, err
	}

	var counts = make(map[string]struct{})
	for i := 0; i < len(recodes); i++ {
		obj := recodes[i]
		winLoseMoney := getPlayerStatementInfoByPkgId(obj.Settlement, pkdId, counts)
		resp.BetMoney += math.Abs(winLoseMoney)
		if winLoseMoney > 0 {
			resp.WinStatementTotal += winLoseMoney
		} else {
			resp.LoseStatementTotal += winLoseMoney
		}
	}
	resp.Count = changeCountsMapToArr(counts)
	return resp, nil
}

func NewStatementResp() StatementTotalResp {
	return StatementTotalResp{
		GameId:   "5b1f3a3cb76a591e7f251721",
		GameName: "百家乐",
	}
}

func changeCountsMapToArr(m map[string]struct{}) []string {
	var res []string
	for k := range m {
		res = append(res, k)
	}
	return res
}

func getPlayerStatementInfo(s []game.SettlementInfo, playerId string) float64 {
	var res float64
	for i := 0; i < len(s); i++ {
		if s[i].PlayerId == playerId {
			f, _ := strconv.ParseFloat(s[i].WinLossGold, 64)
			res = f
		}
	}
	return res
}

func getPlayerStatementInfoByPkgId(s []game.SettlementInfo, pkgId int, m map[string]struct{}) float64 {
	var res float64
	for i := 0; i < len(s); i++ {
		if s[i].PlayerPkgId == pkgId {
			f, _ := strconv.ParseFloat(s[i].WinLossGold, 64)
			res += f
			m[s[i].PlayerId] = struct{}{}
		}
	}
	return res
}
