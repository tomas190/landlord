package controller

import (
	"errors"
	"landlord/game"
	"landlord/mconst/msgIdConst"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

// 当玩家卡死在房间 异常时 提出玩家到房间
func KickRoomPlayer(c *gin.Context) {
	playerId := c.PostForm("player_id") // 玩家Id
	loginOut := c.PostForm("login_out") // 是否退出游戏到大厅(当玩家没有在斗地主游戏中时候)
	token := c.PostForm("token")

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
				//game.SendErrMsg(agent, msgIdConst.ErrMsg, "系统已将你踢出房间,请重新登录游戏")
				kickAll(room, "系统已将你踢出房间,请重新登录游戏.", true)
			} else {
				kickAll(room, "系统已将你踢出房间,请重新进入房间", false)
				//game.SendErrMsg(agent, msgIdConst.ErrMsg, "系统已将你踢出房间,请重新进入房间")
			}
		}
		c.JSON(httpCode, NewResp(SuccCode, "已经踢出玩家", nil))
		return
	}
	c.JSON(httpCode, NewResp(ErrCode, "玩家不在房间中", nil))
	return
}

func kickAll(room *game.Room, msg string, loginOut bool) {
	players := room.Players
	for _, v := range players {
		if v != nil {
			if !v.IsRobot {
				game.SetSessionRoomId(v.Session, "")
				game.SendErrMsg(v.Session, msgIdConst.ErrMsg, msg)
				if loginOut {
					//game.UserLogoutCenter(v.PlayerInfo.PlayerId, game.GetSessionPassword(v.Session))
					game.UserLogoutCenterAfterUnlockMoney(v.PlayerInfo.PlayerId, v.PlayerInfo.Gold)
				}
			}
		}
	}

}

func verifyKickRoomPlayer(token, playerId, isLoginOut string) error {
	if token != game.Server.CenterToken {
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

type UserDataByPackageID struct {
	PackageID int     `json:"packageID"`
	UserData  []int64 `json:"userData"`
}

// 取得在線玩家(分渠道)
func GetOnlineTotal(c *gin.Context) {

	packageIDS := c.DefaultQuery("package_id", "")
	packageID, errP := strconv.ParseInt(packageIDS, 10, 64)
	if packageIDS != "" && errP != nil {
		c.JSON(httpCode, NewResp(ErrCode, "参数错误：package_id为非整数", nil))
		return
	}
	// logger.Debug("packageIDS=%v, packageID%v", packageIDS, packageID)

	data := make(map[string]interface{})
	list := game.GetAllOnlineUser()

	data["game_id"] = game.Server.GameId
	data["game_name"] = "鬥地主2"
	tmplist := make([]UserDataByPackageID, 0)
	var tmp UserDataByPackageID
	if packageIDS != "" && errP == nil {
		if len(list[int(packageID)]) > 0 {
			tmp.PackageID = int(packageID)
			tmp.UserData = list[int(packageID)]
			tmplist = append(tmplist, tmp)
		}
	} else {
		for i, v := range list {
			tmp.PackageID = i
			tmp.UserData = v
			tmplist = append(tmplist, tmp)
		}
	}
	data["game_data"] = tmplist
	c.JSON(httpCode, NewResp(SuccCode, "ok", data))
	logger.Debug("getOnlineUserList request success")
}
