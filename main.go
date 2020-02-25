package main

import (
	"github.com/gin-gonic/gin"
	"landlord/app"
	"landlord/exitdo"
	"landlord/game"
	"os"
)

func main() {
	go app.StartServer(gin.DebugMode, "/ws")

	exitdo.Signal.ListenKill().Done(func(sig os.Signal) {
		game.BackUserToHall()
	})
}
