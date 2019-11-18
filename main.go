package main

import (
	"github.com/gin-gonic/gin"
	"landlord/app"
)

func main() {
	app.StartServer(gin.DebugMode, "/ws")
}
