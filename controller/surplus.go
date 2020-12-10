package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"landlord/game"
	"landlord/mconst/sysSet"
)

func GetSurplusOne(c *gin.Context) {
	gameId := c.DefaultQuery("game_id", "-1")
	if gameId != game.Server.GameId {
		c.JSON(httpCode, NewResp(ErrCode, "err game_id", nil))
		return
	}
	var s game.SurplusPoolOne

	one, err := s.GetLastSurplusOne()
	if err != nil {
		c.JSON(httpCode, NewResp(ErrCode, err.Error(), nil))
	}
	c.JSON(httpCode, NewResp(SuccCode, "ok", one))

}

type UptSurplusConfReq struct {
	PercentageToTotalWin           float64 `json:"percentage_to_total_win" form:"percentage_to_total_win"`
	CoefficientToTotalPlayer       float64 `json:"coefficient_to_total_player" form:"coefficient_to_total_player"`
	FinalPercentage                float64 `json:"final_percentage" form:"final_percentage"`
	PlayerLoseRateAfterSurplusPool float64 `json:"player_lose_rate_after_surplus_pool" form:"player_lose_rate_after_surplus_pool"`
	DataCorrection                 float64 `json:"data_correction" form:"data_correction"`

	RandomPercentageAfterWin  float64 `json:"random_percentage_after_win" form:"random_percentage_after_win"`
	RandomCountAfterWin       float64 `json:"random_count_after_win" form:"random_count_after_win"`
	RandomPercentageAfterLose float64 `json:"random_percentage_after_lose" form:"random_percentage_after_lose"`
	RandomCountAfterLose      float64 `json:"random_count_after_lose" form:"random_count_after_lose"`
}

func UptSurplusConf(c *gin.Context) {
	gameId := c.DefaultPostForm("game_id", "-1")
	if gameId != game.Server.GameId {
		c.JSON(httpCode, NewResp(ErrCode, "err game_id", nil))
		return
	}

	var req UptSurplusConfReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(httpCode, NewResp(ErrCode, err.Error(), nil))
		return
	}

	// 注意传0 的参数  和没传的参数  默认为0的时候 不要修改
	logger.Debug("req:", req)
	percentageToTotalWin := c.DefaultPostForm("percentage_to_total_win", "-1")
	coefficientToTotalPlayer := c.DefaultPostForm("coefficient_to_total_player", "-1")
	finalPercentage := c.DefaultPostForm("final_percentage", "-1")
	playerLoseRateAfterSurplusPool := c.DefaultPostForm("player_lose_rate_after_surplus_pool", "-1")
	dataCorrection := c.DefaultPostForm("data_correction", "null") // 这个字段比较特殊可以传0

	randomPercentageAfterWin := c.DefaultPostForm("random_percentage_after_win", "null")   // 这个字段比较特殊可以传0
	randomCountAfterWin := c.DefaultPostForm("random_count_after_win", "null")             // 这个字段比较特殊可以传0
	randomPercentageAfterLose := c.DefaultPostForm("random_percentage_after_lose", "null") // 这个字段比较特殊可以传0
	randomCountAfterLose := c.DefaultPostForm("random_count_after_lose", "null")           // 这个字段比较特殊可以传0


	var paramsNum int
	if percentageToTotalWin == "-1" {
		paramsNum++
		req.PercentageToTotalWin = -1
	}
	if coefficientToTotalPlayer == "-1" {
		paramsNum++
		req.CoefficientToTotalPlayer = -1
	}
	if finalPercentage == "-1" {
		paramsNum++
		req.FinalPercentage = -1
	}
	if playerLoseRateAfterSurplusPool == "-1" {
		paramsNum++
		req.PlayerLoseRateAfterSurplusPool = -1
	}

	if dataCorrection == "null" {
		paramsNum++
		req.DataCorrection = sysSet.DATA_CORRECTION
	}

	// new
	if randomPercentageAfterWin == "null" {
		paramsNum++
		req.RandomPercentageAfterWin = sysSet.RANDOM_PERCENTAGE_AFTER_WIN
	}

	if randomCountAfterWin == "null" {
		paramsNum++
		req.RandomCountAfterWin = sysSet.RANDOM_COUNT_AFTER_WIN
	}

	if randomPercentageAfterLose == "null" {
		paramsNum++
		req.RandomPercentageAfterLose = sysSet.RANDOM_PERCENTAGE_AFTER_LOSE
	}

	if randomCountAfterLose == "null" {
		paramsNum++
		req.RandomCountAfterLose = sysSet.RANDOM_COUNT_AFTER_LOSE
	}


	// 如果都没传参数 返回当前配置
	//if paramsNum == 5 {
	if paramsNum == 9 {
		var s game.SurplusPoolOne
		one, err := s.GetLastSurplusOne()
		if err != nil {
			c.JSON(httpCode, NewResp(ErrCode, err.Error(), nil))
			return
		}
		req.PlayerLoseRateAfterSurplusPool = one.PlayerLoseRateAfterSurplusPool
		req.CoefficientToTotalPlayer = one.CoefficientToTotalPlayer
		req.FinalPercentage = one.FinalPercentage
		req.PercentageToTotalWin = one.PercentageToTotalWin
		c.JSON(httpCode, NewResp(SuccCode, "没有要修改的参数", req))
		return
	}

	game.UptSurplusConf(req.PercentageToTotalWin,
		req.PlayerLoseRateAfterSurplusPool,
		req.CoefficientToTotalPlayer,
		req.FinalPercentage,req.DataCorrection,
		req.RandomPercentageAfterWin,
		req.RandomCountAfterWin,
		req.RandomPercentageAfterLose,
		req.RandomCountAfterLose)

	req.PlayerLoseRateAfterSurplusPool = sysSet.PLAYER_LOSE_RATE_AFTER_SURPLUS_POOL
	req.CoefficientToTotalPlayer = sysSet.COEFFICIENT_TO_TOTAL_PLAYER
	req.FinalPercentage = sysSet.FINAL_PERCENTAGE
	req.PercentageToTotalWin = sysSet.PERCENTAGE_TO_TOTAL_WIN

	req.RandomPercentageAfterWin = sysSet.RANDOM_PERCENTAGE_AFTER_WIN
	req.RandomCountAfterWin = sysSet.RANDOM_COUNT_AFTER_WIN
	req.RandomPercentageAfterLose = sysSet.RANDOM_PERCENTAGE_AFTER_LOSE
	req.RandomCountAfterLose = sysSet.RANDOM_COUNT_AFTER_LOSE
	c.JSON(httpCode, NewResp(SuccCode, "ok", req))
	return
}
