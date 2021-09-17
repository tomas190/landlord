package game

import "gopkg.in/mgo.v2/bson"

const playerRecodeName = "landlord.player_recode"

// 玩家的记录
type PlayerRecode struct {
	PlayerId string `json:"player_id" bson:"player_id"` // 玩家Id
}

// 如果用户没有登录过则添加 并同时新增盈余池记录
func (p *PlayerRecode) AddPlayerIfNotExist() error {
	session, c := GetDBConn(Server.MongoDBName, playerRecodeName)
	defer session.Close()

	isExist := p.IsExistPlayerId()
	if isExist {
		return nil
	} else {
		err := c.Insert(p)
		if err == nil {
			// var s SurplusPool
			// s.InsertSurplusNewUser()
		}

		return err
	}
}

// 判断名称是否存在
func (p *PlayerRecode) IsExistPlayerId() bool {
	session, c := GetDBConn(Server.MongoDBName, playerRecodeName)
	defer session.Close()
	var item PlayerRecode
	err := c.Find(bson.M{"player_id": p.PlayerId}).One(&item)
	if err != nil {
		return false
	}
	return true
}

// 计算玩家总人数
func (p *PlayerRecode) CountPlayers() int {
	session, c := GetDBConn(Server.MongoDBName, playerRecodeName)
	defer session.Close()

	n, err := c.Find(bson.M{}).Count()
	if err != nil {
		return 0
	}
	return n
}
