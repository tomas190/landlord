package msgIdConst

/*
enum MessageID {
    MSG_PING = 0; // c-s 发送ping
    MSG_PONG = 1; // s-c 发送 pong

    MSG_LOGIN = 100; // c-s 登录
    MSG_LOGIN_RESP = 200; // s-c 登录返回

    MSG_ENTER_ROOM = 101; // c-s 进入房间  返回见 301 302

    MSG_PLAYER_WITH_BET_LIST = 102; // c-s 获取房间类用户列表 // 当用户点击用户列表时候发送
    MSG_PLAYER_WITH_BET_LIST_RESP = 202; // s-c  返回用户列表 // 会按照每10条数据一直推送给客户端

  	// SERVER PUSH  MSG - FROM 300
	MSG_ROOM_LIST_PUSH = 300; // s-c 服务器推送房间列表 // 用户登录过后推送
	MSG_ROOM_INFO_PUSH = 301; // s-c 推送房间信息 // 用户进入房间后推送
	MSG_PLAYER_LIST_PUSH = 302; // s-c 推送房间内用户列表 // 用户进入房间后推送(注：进入房间后推送的用户信息只有6个用于展示 // 每局游戏结算完成后会推送改指令)

	// SERVER PUSH ERROR MSG - FROM 500
	MSG_CLOSE_CONN_PUSH = 500; // s-c 服务器主动断开
	MSG_ERR_MSG_PUSH = 501; // s-c 错误信息推送
}

*/
// 指令请求
const (
	Ping             uint16 = 0
	ReqLogin         uint16 = 100 // 用户登录请求
	ReqEnterRoom     uint16 = 101 // 用户进入房间
	ReqGetLandlordDo uint16 = 102 // 争夺地主的操作
	ReqOutCardDo     uint16 = 103 // 出牌操作

)

// 指令返回
const (
	Pong      uint16 = 1
	RespLogin uint16 = 200 // 用户登录返回
)

// 主动推送
const (
	PushRoomClassify  uint16 = 300 // 推送房间分类信息  用户登录认证成功之后 推送
	PushRoomPlayer           = 301 // 推送房间信息  用户当玩家成功匹配时候 推送
	PushStartGame            = 302 // 推送开始游戏 发牌
	PushCallLandlord         = 303 // 推送叫地主阶段消息
	PushWhoIsLandlord        = 304 // 推送叫地主阶段消息
	PushOutCard              = 305 // 推送出牌消息
)

// 主动推送 Error
const (
	CloseConn uint16 = 500
	ErrMsg    uint16 = 501
)
