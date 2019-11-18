package landlordLog

import (
	"encoding/json"
	"landlord/mconst/logType"
	"github.com/wonderivan/logger"
)

type mlog struct {
	Mark string      `json:"mark"`
	Log  interface{} `json:"log"`
}

func RecodeLog(lType int, logInfo interface{}, mark string) {

	var l mlog
	l.Mark = mark
	l.Log = logInfo
	bytes, _ := json.Marshal(l)
	log := string(bytes)

	switch lType {
	case logType.Info:
		logger.Info(log)
	case logType.Debug:
		logger.Debug(log)
	case logType.Error:
		logger.Error(log)
	default:
		logger.Debug("位置日志类型:", lType)
		logger.Debug("log:", lType)
	}

}
