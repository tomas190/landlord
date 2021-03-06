package app

import (
	"fmt"
	"landlord/controller"
	"landlord/game"
	"landlord/handler"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
)

// Init server with mode
// @mode gin framework mode
func StartServer(mode, wsPath string) {
	game.HttpPostToTelegram(fmt.Sprintf("鬥地主2 游戏服务器启动成功\n时间 : %v\n版本号 : %v", time.Now().Format("2006-01-02 15:04:05"), game.VersionCode))

	gin.SetMode(mode)
	r := gin.New()
	m := melody.New()

	//r.GET("/", func(c *gin.Context) {
	//	http.ServeFile(c.Writer, c.Request, "./view/index.html")
	//})

	//r.GET("/im/ws", func(c *gin.Context) {
	r.GET(wsPath, func(c *gin.Context) {
		err := m.HandleRequestWithKeys(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error("HandlerRequest Fail:", err.Error())
		}
	})

	//m.Config.PongWait = time.Second * 100
	mController(r)
	mHandler(m)
	err := r.Run(":" + game.Server.Port)
	if err != nil {
		logger.Painc("Start server fail:", err.Error())
	}
}

func mHandler(m *melody.Melody) {
	// onMessage
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		handler.OnMessage(m, s, msg)
	})

	// onClose
	m.HandleClose(func(session *melody.Session, i int, s string) error {
		return handler.OnClose(m, session, i, s)
	})

	// onConn
	m.HandleConnect(func(session *melody.Session) {
		handler.OnConnect(m, session)
	})

	// onDisconnect
	m.HandleDisconnect(func(session *melody.Session) {
		handler.OnDisconnect(m, session)
	})

	// onErr
	m.HandleError(func(session *melody.Session, e error) {
		handler.OnErr(session, e)
	})

	// onPong
	m.HandlePong(func(session *melody.Session) {
		handler.OnPong(session)
	})

	// onSentMessage
	m.HandleSentMessage(func(session *melody.Session, bytes []byte) {
		handler.OnSentMessage(session, bytes)
	})

	// onMessageBinary
	m.HandleMessageBinary(func(session *melody.Session, bytes []byte) {
		handler.OnMessageBinary(m, session, bytes)
	})

	// onSentMessageBinary
	m.HandleSentMessageBinary(func(session *melody.Session, bytes []byte) {
		handler.OnSentMessageBinary(session, bytes)
	})

}

func init() {
	game.InitConfig()
	game.StartCenter()
}

// api
// 当玩家卡死采用接口踢出玩家
func mController(r *gin.Engine) {
	r.GET("/api/kickRoomPlayer", controller.KickRoomPlayer)
	r.GET("/api/getGameData", controller.GetLandlordData)
	r.GET("/api/version", controller.Version)
	r.GET("/api/getLog", controller.GetLog)

	r.GET("/api/getSurplusOne", controller.GetSurplusOne)
	r.POST("/api/uptSurplusConf", controller.UptSurplusConf)

	r.GET("/api/getPlayInfo", controller.GetPlayInfo)
	r.POST("/api/UptServer", controller.UptServer)

	r.GET("/api/getStatementTotal", controller.GetStatementTotal)
	r.GET("/api/getOnlineTotal", controller.GetOnlineTotal)
}
