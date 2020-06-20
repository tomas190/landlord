package app

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
	"landlord/controller"
	"landlord/game"
	"landlord/handler"
)

// Init server with mode
// @mode gin framework mode
func StartServer(mode, wsPath string) {
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
	r.GET("/api/getSurplusOne", controller.UptSurplusConf)
}