package game

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/wonderivan/logger"
	"math/rand"
	"time"
)

func GetMsgId(mid []byte) uint16 {
	u := binary.BigEndian.Uint16(mid)
	return u
}

func PkgMsg(msgId uint16, data []byte) []byte {

	bytes := make([]byte, 2+len(data))

	binary.BigEndian.PutUint16(bytes[:2], msgId)

	copy(bytes[2:], data)

	return bytes

}

func PrintMsg(msgInfo string, data interface{}) {
	bytes, _ := json.Marshal(data)
	logger.Debug(msgInfo, string(bytes))

}

func GetAuthKey(userId string) string {
	key := uuid.New().String() + userId
	authKey := MakeAuthKey(key)
	return authKey
}

func MakeAuthKey(key string) string {
	data := []byte(key)
	sum := md5.Sum(data)
	code := fmt.Sprintf("%x", sum)
	return code
}

var count int64

// 获取一个随机数 [startNum:endNum]
func RandNum(startNum, endNum int) int {
	count++
	if count >= 1<<16 {
		count = 0
	}
	rand.Seed(time.Now().UnixNano() + count)
	rnd := rand.Intn(endNum - startNum + 1)
	return rnd + startNum
}

func DelaySomeTime(seconds time.Duration) {
	// 延时1秒后推送开局信息
	delay := time.NewTicker(seconds * time.Second)
	<-delay.C
}

func changToArrInt32(arr []byte) []int32 {
	var result []int32
	for i := 0; i < len(arr); i++ {
		result = append(result, int32(arr[i]))
	}
	return result
}

func printGroup(gc GroupCard) {
	logger.Debug("单张")
	rcs := gc.Single
	for i := 0; i < len(rcs); i++ {
		fmt.Print("weight:", rcs[i].Wight, "  ")
		PrintCard(rcs[i].Card)
	}

	logger.Debug("对子")
	rcd := gc.Double
	for i := 0; i < len(rcd); i++ {
		fmt.Print("weight:", rcd[i].Wight, "  ")
		PrintCard(rcd[i].Card)
	}

	logger.Debug("三张")
	rct := gc.Triple
	for i := 0; i < len(rct); i++ {
		fmt.Print("weight:", rct[i].Wight, "  ")
		PrintCard(rct[i].Card)
	}

	logger.Debug("炸弹")
	rcb := gc.Bomb
	for i := 0; i < len(rcb); i++ {
		fmt.Print("weight:", rcb[i].Wight, "  ")
		PrintCard(rcb[i].Card)
	}

	logger.Debug("火箭")
	rcr := gc.Rocket
	for i := 0; i < len(rcr); i++ {
		fmt.Print("weight:", rcr[i].Wight, "  ")
		PrintCard(rcr[i].Card)
	}

	logger.Debug("顺子")
	rcj := gc.Junko
	for i := 0; i < len(rcj); i++ {
		fmt.Print("weight:", rcj[i].Wight, "  ")
		PrintCard(rcj[i].Card)
	}

	logger.Debug("连对")
	rcjd := gc.JunkoDouble
	for i := 0; i < len(rcjd); i++ {
		fmt.Print("weight:", rcjd[i].Wight, "  ")
		PrintCard(rcjd[i].Card)
	}

	logger.Debug("飞机")
	rcjt := gc.JunkTriple
	for i := 0; i < len(rcjt); i++ {
		fmt.Print("weight:", rcjt[i].Wight, "  ")
		PrintCard(rcjt[i].Card)
	}
}
