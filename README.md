一、部署
1.克隆代码到本地 git clone http://git.0717996.com/fido/landlord.git

2.进入baccarat文件夹 cd landlord

3.将build.sh 和 start.sh修改权限为可执行 chmod 777 build.sh start.sh

4.编译 ./build.sh

5.运行 ./start.sh

6.查看是否运行成功 tail -f out.log

如果看到最后是这样的日志代表成功启动
5:20PM [DEBG] [/data/server/baccarat/game/game_control.go:18] StartGame>>>>>>>>>>
5:20PM [DEBG] [/data/server/baccarat/game/center_receive_msg.go:17] 登录中心服成功

二、配置文件及相关说明
1.项目所需套件

1.go语言环境   go version go1.13 linux/amd64
2.mongo数据库  MongoDB server version: 4.0.12
2. 配置文件说明及日志(dev为例)

1.配置文件路径：conf/server.json   (需添加jenkins忽略文件)
2.日志配置文件路径：conf/log.json   (默认不需要更改)(需添加jenkins忽略文件)
3.日志文件路径及名称：out.log 和编译好的可执行文件同级
服务配置说明 {

{
  "Port": "1221",                               //     斗地主项目启动端口
  "MaxConn": 1200,                              //     最大连接数
  "CenterToken": "963258",                      //     中心服Token
  "CenterDomain": "172.16.100.2:9502",          //     中心服域名
  "CenterPort": "9502",                         //     中心服端口
  "DevKey": "new_game_4",                       //     devKey
  "DevName": "新游戏开发",                       //     devName
  "GameId": "5b1f3a3cb76a591e7f251711",         //     gameId  
  "GameTaxRate": 0.05,                          //     玩家赢钱税收比例
  "WinGoldNotice": 100,                         //     玩家一局赢钱超过此设定将发送跑马灯通知
  "MongoDBAddr": "172.16.100.5:27017",          //     mongo数据库连接地址
  "MongoDBUser": "LANDLORD-game",               //     mongo连接用户名
  "MongoDBPwd": "123456",                       //     mongo连接密码
  "MongoDBAuth": "admin",                       //     mongo认证(可不填默认admin)
  "MongoDBName": "LANDLORD-Game",               //     mongo连接游戏库名

  "UrlSendLog": "http://172.16.100.7:4151/pub?topic=game-server"
}

日志配置说明 {

"TimeFormat":"2006-01-02 15:04:05", // 输出日志开头时间格式(详见**时间格式**)
"Console": {            // 控制台日志配置
    "level": "TRAC",    // 控制台日志输出等级(详见**日志等级**)
    "color": true       // 控制台日志颜色开关 
},
"File": {                   // 文件日志配置
    "filename": "app.log",  // 初始日志文件名
    "level": "TRAC",        // 日志文件日志输出等级
    "daily": true,          // 跨天后是否创建新日志文件，当append=true时有效
    "maxlines": 1000000,    // 日志文件最大行数，当append=true时有效
    "maxsize": 1,           // 日志文件最大大小，当append=true时有效
    "maxdays": 3,           // 日志文件有效期
    "append": true,         // 是否支持日志追加
    "permit": "0660"        // 新创建的日志文件权限属性
},
"Conn": {                       // 网络日志配置
    "net":"tcp",                // 日志传输模式
    "addr":"0.0.0.0",           // 日志接收服务器
    "level": "Warn",            // 网络日志输出等级
    "reconnect":true,           // 网络断开后是否重连
    "reconnectOnMsg":false,     // 发送完每条消息后是否断开网络
}
}

日志等级 等级 配置 释义 控制台颜色
0 EMER 系统级紧急，比如磁盘出错，内存异常，网络不可用等 红色底
1 ALRT 系统级警告，比如数据库访问异常，配置文件出错等 紫色
2 CRIT 系统级危险，比如权限出错，访问异常等 蓝色
3 EROR 用户级错误 红色
4 WARN 用户级警告 黄色
5 INFO 用户级重要 天蓝色
6 DEBG 用户级调试 绿色
7 TRAC 用户级基本输出 绿色

时间格式 ANSIC "Mon Jan _2 15:04:05 2006"
UnixDate "Mon Jan _2 15:04:05 MST 2006"
RubyDate "Mon Jan 02 15:04:05 -0700 2006"
RFC822 "02 Jan 06 15:04 MST"
RFC822Z "02 Jan 06 15:04 -0700"
RFC850 "Monday, 02-Jan-06 15:04:05 MST"
RFC1123 "Mon, 02 Jan 2006 15:04:05 MST"
RFC1123Z "Mon, 02 Jan 2006 15:04:05 -0700"
RFC3339 "2006-01-02T15:04:05Z07:00"
RFC3339Nano "2006-01-02T15:04:05.999999999Z07:00"
Kitchen "3:04PM"
Stamp "Jan _2 15:04:05"
StampMilli "Jan _2 15:04:05.000"
StampMicro "Jan _2 15:04:05.000000"
StampNano "Jan _2 15:04:05.000000000"
RFC3339Nano1 "2006-01-02 15:04:05.999999999 -0700 MST"
DEFAULT "2006-01-02 15:04:05"