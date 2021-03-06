package game

import (
	"sync"
)

var userLoginChan struct {
	chanMap map[string]chan *UserLoginCallBack
	rwMutex sync.RWMutex
}

var centerWriteMutex sync.Mutex

func StartCenter() {

	opMap.opMap = make(map[string]string, 16)

	userLoginChan.chanMap = make(map[string]chan *UserLoginCallBack, 1000)
	//	go changeToken()
	DelaySomeTime(1)
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
