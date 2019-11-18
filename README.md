###### **一、部署**

**1.克隆代码到本地**
`git clone http://git.0717996.com/fido/baccarat.git`

**2.进入baccarat文件夹**
`cd baccarat`

**3.编译**
`./build`

**4.运行**
`./start`

**5.查看是否运行成功**
`tail -f out.log`

如果看到最后是这样的日志代表成功启动   
5:20PM [DEBG] [/data/server/baccarat/game/game_control.go:18] StartGame>>>>>>>>>>   
5:20PM [DEBG] [/data/server/baccarat/game/center_receive_msg.go:17] 登录中心服成功   
5:20PM [DEBG] [/data/server/baccarat/game/center_receive_msg.go:18] 用户登录中心服成功   




###### **二、配置文件及相关说明**
    
**1.项目所需套件**
    1.go语言环境   go version go1.13 linux/amd64
    2.mongo数据库  MongoDB server version: 4.0.12

**2. 配置文件说明及日志(dev为例)**
    1.配置文件路径：conf/server.json   (需添加jenkins忽略文件)
    2.日志文件路径及名称：out.log 和编译好的可执行文件同级
    3.参数说明：如下
{    
  "Port": "1227",                           百家乐项目启动端口     
  "CenterToken": "963258",                  最大连接数    
  "CenterDomain": "172.16.100.2:9502",      中心服域名    
  "CenterPort": "9502",                     中心服端口    
  "APIGetToken": "/Token/getToken",         获取中心服token API    
  "DevKey": "new_game_22",                  devKey    
  "DevName": "新游戏开发",                   devName    
  "GameId": "5b1f3a3cb76a591e7f251721",     gameId    
  "GameTaxRate": 0.05,                      玩家赢钱税收比例    
    
  "MongoDBAddr": "172.16.100.5:27017",      mongo数据库连接地址    
  "MongoDBUser": "bjl",                     mongo连接用户名    
  "MongoDBPwd": "123456",                   mongo连接密码    
  "MongoDBAuth": "",                        mongo认证(可不填默认admin)    
  "MongoDBName": "BACCARAT-Game"            mongo连接游戏库名    
}    

