package handler

import (
	"landlord/game"
	"landlord/mconst/msgIdConst"
	"landlord/mconst/userSessionStatus"
	"landlord/msg/mproto"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
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

	playerInfo, err := userLoginVerify(req.UserId, req.UserPassword)
	if err != nil {
		game.SendErrMsg(session, msgIdConst.ReqLogin, err.Error())
		return
	}

	// 转换成 proto 对象
	protoPlayerInfo := game.ChangePlayerInfoToProto(playerInfo)
	// 重复登录 挤下线机制 (如果该账号已经登录 则断开连接并清楚map)
	userRepeatLogin(m, req.UserId)

	// 返回玩家信息
	bytes, _ := proto.Marshal(protoPlayerInfo)
	_ = session.WriteBinary(game.PkgMsg(msgIdConst.RespLogin, bytes))
	// 返回玩家信息

	// 推送房间分类信息
	game.PushRoomClassify(session)

	// 保存用户信息到session 并添加登录成功tag
	game.SetSessionPlayerInfo(session, playerInfo)
	game.SetSessionIsLogin(session)

}

// 向中心服发送登录验证请求
func userLoginVerify(userId, password string) (*game.PlayerInfo, error) {

	loginChan := make(chan *game.UserLoginCallBack)
	game.SaveUserLoginCallBack(userId, loginChan)
	defer func() {
		close(loginChan)
		game.RemoveUserLoginCallBack(userId)
		logger.Debug("userLogin:", game.GetUserLoginCallBackLen())
	}()
	game.UserLogin(userId, password)
	select {
	case userLogin := <-loginChan:
		userInfo := userLogin.Player

		return &userInfo, nil

	case <-time.After(time.Second * 2):
		logger.Error("登录超时 中心服无返回!")
		return nil, errors.New("登录超时")
	}

}

// 你已经被挤下线
func userRepeatLogin(m *melody.Melody, userId string) {

	var push mproto.CloseConn
	push.Code = userSessionStatus.LoginTimeOutClose
	push.Msg = "你已经被挤下线!"

	bytes, _ := proto.Marshal(&push)

	agent := game.GetAgent(userId)
	if agent != nil {
		_ = agent.CloseWithMsg(game.PkgMsg(msgIdConst.CloseConn, bytes))
		game.RemoveAgent(userId)
	}

}
