package game

import "sync"

var opMap orderPlayerMap

type orderPlayerMap struct {
	opMap map[string]string
	mu    sync.Mutex
}

func (o orderPlayerMap) Set(order, playerId string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.opMap[order] = playerId
}

func (o orderPlayerMap) Get(order string)string {
	o.mu.Lock()
	defer o.mu.Unlock()
	playerId:=o.opMap[order]
	delete(o.opMap,order)
	return playerId
}