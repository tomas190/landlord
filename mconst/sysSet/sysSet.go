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
