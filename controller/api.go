package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"landlord/game"
	"landlord/mconst/msgIdConst"
)

// 当玩家卡死在房间 异常时 提出玩家到房间
func KickRoomPlayer(c *gin.Context) {
	playerId := c.PostForm("player_id") // 玩家Id
	loginOut := c.PostForm("login_out") // 是否退出游戏到大厅(当玩家没有在斗地主游戏中时候)
	token := c.PostForm("player_id")

	if err := verifyKickRoomPlayer(token, playerId, loginOut); err != nil {
		c.JSON(httpCode, NewResp(ErrCode, err.Error(), nil))
		return
	}

	room, b := game.IsPlayerInRoom(playerId)
	if b && room != nil {
		agent := game.GetAgent(playerId)
		game.RemoveRoom(room.RoomId)
		if agent != nil {
			// 清空玩家room信息
			game.SetSessionRoomId(agent, "")
			if loginOut == "yes" {
				game.UserLogoutCenter(playerId, game.GetSessionPassword(agent))
				game.SendErrMsg(agent, msgIdConst.ErrMsg, "系统已将你踢出房间,请重新登录游戏.")
			}else {
				game.SendErrMsg(agent, msgIdConst.ErrMsg, "系统已将你踢出房间,请重新进入房间")
			}
		}
		c.JSON(httpCode, NewResp(SuccCode,"已经踢出玩家", nil))
		return
	}
	c.JSON(httpCode, NewResp(ErrCode,"玩家不在房间中", nil))
	return
}

func verifyKickRoomPlayer(token, playerId, isLoginOut string) error {
	if token !=game.Server.CenterToken {
		return errors.New("验证失败")
	}

	if playerId == "" {
		return errors.New("验证失败01")
	}

	if isLoginOut == "" {
		return errors.New("验证失败02")
	}
	return nil
}
