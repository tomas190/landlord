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
	<-time.After(time.Second * seconds)
}

func changToArrInt32(arr []byte) []int32 {
	var result []int32
	for i := 0; i < len(arr); i++ {
		result = append(result, int32(arr[i]))
	}
	return result
}
