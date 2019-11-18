package game

import (
	"sync"
	"time"
)

var userLoginChan struct {
	chanMap map[string]chan *UserLoginCallBack
	rwMutex sync.RWMutex
}

func StartCenter() {

	userLoginChan.chanMap = make(map[string]chan *UserLoginCallBack, 1000)
	//	go changeToken()
	time.Sleep(time.Second)
	ConnectCenterWs()
}

func SaveUserLoginCallBack(key string, userInfo chan *UserLoginCallBack) {
	userLoginChan.rwMutex.Lock()
	userLoginChan.chanMap[key] = userInfo
	userLoginChan.rwMutex.Unlock()

}

func GetUserLoginCallChan(key string) chan *UserLoginCallBack {
	userLoginChan.rwMutex.Lock()
	defer userLoginChan.rwMutex.Unlock()
	session := userLoginChan.chanMap[key]
	return session
}

func RemoveUserLoginCallBack(key string) {
	userLoginChan.rwMutex.Lock()
	defer userLoginChan.rwMutex.Unlock()
	_, ok := userLoginChan.chanMap[key]
	if ok {
		delete(userLoginChan.chanMap, key)
	}
}

func GetUserLoginCallBackLen() int {
	userLoginChan.rwMutex.Lock()
	defer userLoginChan.rwMutex.Unlock()
	i := len(userLoginChan.chanMap)
	return i
}
