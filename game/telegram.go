package game

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/wonderivan/logger"
)

var (
	OrderIDToOrderInfo = sync.Map{} // order 暫存，用於中心服回傳失敗時，發送訊息與登出玩家
	VersionCode        = "1.0.17 sp"
)

type OrderInfo struct {
	PlayerId string
	Event    string
	// Session  *melody.Session
}

func HttpPostToTelegram(msg string) {

	verMsg := "\n環境 : "
	switch Server.CenterDomain {
	case "161.117.178.174:12345":
		verMsg += "DEV"
		// fmt.Println("httpPostToTelegram :\n", msg+verMsg)
		logger.Debug("httpPostToTelegram :\n", msg+verMsg)
		return
	case "172.16.1.41:9502":
		verMsg += "PRE"
	default:
		verMsg += "OL"
	}
	msg += verMsg
	// 發送消息給 Telegram group LogAlerts
	chat_id := "-521977907"
	resp, err := http.PostForm("https://api.telegram.org/bot1726462670:AAEmwMgpIpxk0akDE3k-MuQCQ3rZm3NWGFU/sendMessage", url.Values{"chat_id": {chat_id}, "text": {msg}})
	if err != nil {
		// fmt.Println("http.PostForm error :", err)
		logger.Debug("http.PostForm error :", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println("ioutil.ReadAll error :", err)
		logger.Debug("ioutil.ReadAll error :", err)
	}
	// fmt.Println("Post to Telegram : ", string(body))
	logger.Debug("Post to Telegram : ", string(body))
}
