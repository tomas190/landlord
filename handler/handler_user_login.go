package handler

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
	"landlord/game"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/userSessionStatus"
	"landlord/msg/mproto"
	"time"
)

// 登录请求
//
func ReqLogin(m *melody.Melody, session *melody.Session, data []byte) {
	logger.Debug("=== ReqLogin ===")
	req := &mproto.ReqLogin{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		game.SendErrMsg(session, msgIdConst.ReqLogin, "请求数据异常:"+err.Error())
		return
	}

	game.PrintMsg("登录请求参数:", req)
	/*==== 参数验证 =====*/

	//playerInfo, err := userLoginVerify(req.UserId, req.UserPassword)
	playerInfo, playerPkgId,err := userLoginVerify(req.UserId, req.UserPassword, req.Token)
	if err != nil {
		game.SendErrMsg(session, msgIdConst.ReqLogin, err.Error())
		return
	}

	// 转换成 proto 对象
	//protoPlayerInfo := game.ChangePlayerInfoToProto(playerInfo)
	// 重复登录 挤下线机制 (如果该账号已经登录 则断开连接并清楚map)
	userRepeatLogin(m, session, playerInfo)

	// 返回玩家信息
	var loginResp mproto.RespLogin
	loginResp.PlayerInfo = playerInfo
	bytes, _ := proto.Marshal(&loginResp)
	_ = session.WriteBinary(game.PkgMsg(msgIdConst.RespLogin, bytes))
	// 返回玩家信息

	// 推送房间分类信息
	game.PushRoomClassify(session)

	// 保存用户信息到session 并添加登录成功tag
	p := game.ChangePlayerP2S(*playerInfo)
	game.SetSessionPlayerInfo(session, &p)
	game.SetSessionIsLogin(session)
	game.SetSessionPassword(session, req.UserPassword)
	game.SaveAgent(playerInfo.PlayerId, session)
	session.Set("playerTax", game.GetPlatformTaxPercent(playerPkgId))

	// 记录用户登录
	var playerRecode game.PlayerRecode
	playerRecode.PlayerId = req.UserId
	err = playerRecode.AddPlayerIfNotExist()
	if err != nil {
		logger.Error("记录玩家登录失败:", err.Error())
	}

	logger.Info("当前连接数:", m.Len())

}

// 向中心服发送登录验证请求
func userLoginVerify(userId, password, token string) (*mproto.PlayerInfo, int, error) {

	//var pi mproto.PlayerInfo
	//pi.PlayerId = userId
	//pi.Gold = 6300
	//pi.PlayerName = userId
	//pi.PlayerImg = "2.png"
	//return &pi, nil

	loginChan := make(chan *game.UserLoginCallBack)
	game.SaveUserLoginCallBack(userId, loginChan)
	defer func() {
		close(loginChan)
		game.RemoveUserLoginCallBack(userId)
		logger.Debug("userLoginChan:", game.GetUserLoginCallBackLen())
	}()
	game.UserLogin(userId, password, token)
	select {
	case userInfo := <-loginChan:
		var ui mproto.PlayerInfo
		ui.PlayerId = userInfo.Player.PlayerId
		ui.PlayerImg = userInfo.Player.HeadImg
		ui.PlayerName = userInfo.Player.Name
		ui.Gold = userInfo.Player.Gold
		return &ui, userInfo.Player.PlayerPkgId, nil

	case <-time.After(time.Second * 2):
		//game.SendLogToCenter("ERR", "handler/handler.go", "96", "用户登录超时 中心服无返回!")
		logger.Error("登录超时 中心服无返回!")
		return nil, 0, errors.New("登录超时")
	}

}

// 你已经被挤下线
func userRepeatLogin(m *melody.Melody, currentSession *melody.Session, playerInfo *mproto.PlayerInfo) {

	var push mproto.CloseConn
	push.Code = userSessionStatus.LoginTimeOutClose
	push.Msg = "你已经被挤下线!"

	bytes, _ := proto.Marshal(&push)

	oldSession := game.GetAgent(playerInfo.PlayerId)
	if oldSession != nil {
		// 如果旧的连接在房间中游戏  则需要同步session信息
		syncSessionInfo(oldSession, currentSession, playerInfo)
		_ = oldSession.CloseWithMsg(game.PkgMsg(msgIdConst.CloseConn, bytes))
		game.RemoveAgent(playerInfo.PlayerId)
	}

}

// 同步session的房间信息
func syncSessionInfo(oldSession, session *melody.Session, info *mproto.PlayerInfo) {
	roomId := game.GetSessionRoomId(oldSession)
	if roomId != "" {

		room := game.GetRoom(roomId)
		if player, ok := room.Players[info.PlayerId]; ok {
			player.PlayerInfo.Name = info.PlayerName
			player.PlayerInfo.HeadImg = info.PlayerImg
			player.Session = session
			logger.Debug("=============================== 已经恢复连线==============================")
			player.IsCloseSession = false
			game.SetSessionCloseTag(session,false)
		}
		game.SetSessionRoomId(session, roomId)
	}

}
