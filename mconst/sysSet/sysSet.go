package sysSet

/*
	系统设置常量包
*/

import "time"

const (
	GameName = "landlord"

	GameDelayTime    time.Duration = 30 // 玩牌等待时间
	GameDelayTimeInt int32         = 30

	GameDelayGetLandlordTime    time.Duration = 20 // 叫地主等待时间
	GameDelayGetLandlordTimeInt int32         = 20

	GameDelayReadyTimeInt int32 = 10
)

/*
盈余池 = （该游戏全部实际的玩家历史总输 - （该游戏全部实际的玩家历史总赢 * 100%）- （该游戏的历史实际的玩家总数 * 0））* 50%，当盈余池小于0的时候，玩家70%的机率为输
*/
var (
	PERCENTAGE_TO_TOTAL_WIN             float64 = 1   // 100% 历史总赢乘的百分比字段（100%那个值
	PLAYER_LOSE_RATE_AFTER_SURPLUS_POOL float64 = 0.7 // 70%  盈余池后的玩家输百分比（70%那个值）
	COEFFICIENT_TO_TOTAL_PLAYER         float64 = 0   // 0    玩家总数所剩的系数（0那个值)
	FINAL_PERCENTAGE                    float64 = 0.5 // 50%  最后百分比（50%那个值）
	DATA_CORRECTION                     float64 = 0   // 异常数据修正
)

func InitSurplusConf(percentageToTotalWin,
	playerLoseRateAfterSurplusPool,
	coefficientToTotalPlayer,
	finalPercentage ,dataCorrection float64) {

	PERCENTAGE_TO_TOTAL_WIN = percentageToTotalWin
	PLAYER_LOSE_RATE_AFTER_SURPLUS_POOL = playerLoseRateAfterSurplusPool
	COEFFICIENT_TO_TOTAL_PLAYER = coefficientToTotalPlayer
	FINAL_PERCENTAGE = finalPercentage
	DATA_CORRECTION = dataCorrection
}
