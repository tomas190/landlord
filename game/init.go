package game

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"landlord/mconst/sysSet"
	"landlord/msg/mproto"
	"strings"
	"time"
)

// 服务器配置
var Server struct {
	Port         string
	MaxConn      int
	CenterToken  string // 中心服API token
	CenterDomain string // 中心服域名
	CenterPort   string // 中心服端口
	APIGetToken  string // 中心服获取端口Api

	DevKey        string
	DevName       string
	GameId        string
	GameTaxRate   float64
	WinGoldNotice float64
	UseRobot      bool // 启用机器人

	MongoDBAddr string
	MongoDBUser string
	MongoDBPwd  string
	MongoDBAuth string
	MongoDBName string

	UrlSendLog string
}

var globalSession *mgo.Session
var roomClassify *mproto.PushRoomClassify

func InitConfig() {
	initLogger()
	initServerConf()
	initMongoDb()
	initRoomClassify()

	initSurplusConf()
	//cornMatch()
}

// 加载日志适配器
func initLogger() {
	//
	bytes, err := ioutil.ReadFile("conf/log.json")
	if err != nil {
		panic("[server/init.go:115] 加载日志配置失败:" + err.Error())
		return
	}

	confJson, err := simplejson.NewJson(bytes)
	if err != nil {
		panic("[server/init.go:121] 加载日志配置失败(not json):" + err.Error())
		return
	}

	fileFullName, err := confJson.Get("File").Get("filename").String()

	suffix := strings.LastIndex(fileFullName, `/`)
	if suffix != -1 {
		// filePath := fileFullName[:suffix]
		// err := common.CheckAndMakePath(filePath)
		if err != nil {
			// panic("[server/init.go:121] 加载日志配置失败(自动创建日志文件夹失败):" + err.Error())
			return
		}
	}

	err = logger.SetLogger("conf/log.json")
	//err := logger.SetLogger("D:/work/im-version1/im-server/conf/log.json")
	if err != nil {
		panic("[server/init.go] 加载日志配置失败:" + err.Error())
		return
	}
	logger.Info("日志加载成功...")
}

// 加载服务器配置
func initServerConf() {
	data, err := ioutil.ReadFile("conf/server.json")

	//data, err := ioutil.ReadFile("D:/work/im-version1/im-server/conf/server.json")
	if err != nil {
		logger.Painc("加载服务器配置失败!", err.Error())
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		logger.Painc("加载服务器配置失败!", err.Error())
	}
	logger.Info("加载服务器配置成功,正在启动服务...")

	bytes, _ := json.Marshal(Server)

	logger.Info("配置信息:", string(bytes))
}

// 初始化mongoDB
func initMongoDb() {

	//logger.Debug(".MongoDBAddr:", Server.MongoDBAddr)
	//logger.Debug(".MongoDBAuth,", Server.MongoDBAuth)
	//logger.Debug(".MongoDBUser:", Server.MongoDBUser)
	//logger.Debug(".MongoDBPwd:", Server.MongoDBPwd)
	//logger.Debug(".MongoDBName:", Server.MongoDBName)
	var err error
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{Server.MongoDBAddr},
		Timeout:  60 * time.Second,
		Database: Server.MongoDBAuth,
		Username: Server.MongoDBUser,
		Password: Server.MongoDBPwd,
	}
	globalSession, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic("Mongodb 连接失败:" + err.Error())
	}
	globalSession.SetMode(mgo.Monotonic, true)
	logger.Info("连接mongo数据库成功 address:", Server.MongoDBAddr)
	InitSurplusPool()
	UptSurplusPoolOne()

	// 同步盈余池 每隔两秒执行
	go func() {
		for true {
			DelaySomeTime(2)
			UptSurplusPoolOne()
		}
	}()
}

// 初始化房间分类信息
func initRoomClassify() {
	var resp mproto.PushRoomClassify
	var i int32
	for i = 1; i <= 4; i++ {
		var roomClassify mproto.RoomClassify
		roomClassify.RoomType = i
		roomClassify.BottomPoint = GetRoomClassifyBottomPoint(i)
		roomClassify.BottomEnterPoint = GetRoomClassifyBottomEnterPoint(i)
		resp.RoomClassify = append(resp.RoomClassify, &roomClassify)
	}
	roomClassify = &resp
}

func GetDBConn(dbName, cName string) (*mgo.Session, *mgo.Collection) {
	s := globalSession.Copy()
	c := s.DB(dbName).C(cName)
	return s, c
}

// 原始的卡牌张数
func originalCardNum() map[int32]int32 {
	m := make(map[int32]int32, 15)
	var i int32
	for i = 1; i <= 15; i++ {
		if i >= 14 {
			m[i] = 1
		} else {
			m[i] = 4
		}
	}
	return m
}

func initSurplusConf() {
	var s SurplusPoolOne

	one, err := s.GetLastSurplusOne()
	if err != nil {
		panic("初始化盈余池配置失败!" + err.Error())
	}
	sysSet.InitSurplusConf(
		one.PercentageToTotalWin,
		one.PlayerLoseRateAfterSurplusPool,
		one.CoefficientToTotalPlayer,
		one.FinalPercentage)
}
