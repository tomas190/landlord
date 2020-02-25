package main

import (
	"github.com/gin-gonic/gin"
	"landlord/app"
)

func main() {
	app.StartServer(gin.DebugMode, "/ws")
	//go app.StartServer(gin.DebugMode, "/ws")
	//logger.Debug("程序启动 监听信号量...")
	//exitdo.Signal.ListenKill().Done(func(sig os.Signal) {
	//	logger.Debug("程序关闭:",sig)
	//	game.BackUserToHall()
	//})
}
