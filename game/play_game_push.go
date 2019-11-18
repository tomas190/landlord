package game

import (
	"github.com/golang/protobuf/proto"
	"gopkg.in/olahol/melody.v1"
	"landlord/mconst/msgIdConst"
	"landlord/msg/mproto"
)

func PushRoomClassify(session *melody.Session) {

	var resp mproto.PushRoomClassify
	var i int32
	for i = 1; i <= 4; i++ {
		var roomClassify mproto.RoomClassify
		roomClassify.RoomType = i
		roomClassify.BottomPoint = GetRoomClassifyBottomPoint(i)
		roomClassify.BottomEnterPoint = GetRoomClassifyBottomEnterPoint(i)
		resp.RoomClassify = append(resp.RoomClassify, &roomClassify)
	}

	bytes, _ := proto.Marshal(&resp)
	_ = session.WriteBinary(PkgMsg(msgIdConst.PushRoomClassify, bytes))

}
