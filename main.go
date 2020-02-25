package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"landlord/app"
	"landlord/exitdo"
	"landlord/game"
	"os"
)

func main() {
	go app.StartServer(gin.DebugMode, "/ws")

	exitdo.Signal.ListenKill().Done(func(sig os.Signal) {
		logger.Debug("程序关闭:",sig)
		game.BackUserToHall()
	})
}
