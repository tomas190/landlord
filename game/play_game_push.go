package game

import (
	"github.com/golang/protobuf/proto"
	"gopkg.in/olahol/melody.v1"
	"landlord/mconst/msgIdConst"
)

// 推送房间 分类消息
func PushRoomClassify(session *melody.Session) {
	resp := roomClassify
	bytes, _ := proto.Marshal(resp)
	_ = session.WriteBinary(PkgMsg(msgIdConst.PushRoomClassify, bytes))

}
