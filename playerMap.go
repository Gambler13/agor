package main

import "sync"

type RegularIntMap struct {
	sync.RWMutex
	internal map[string]*Player
}

func NewRegularIntMap() *RegularIntMap {
	return &RegularIntMap{
		internal: make(map[string]*Player),
	}
}

func (rm *RegularIntMap) Load(key string) (value *Player, ok bool) {
	rm.RLock()
	result, ok := rm.internal[key]
	rm.RUnlock()
	return result, ok
}

func (rm *RegularIntMap) Delete(key string) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}

func (rm *RegularIntMap) Store(key string, value *Player) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}

func (rm *RegularIntMap) ForEach(f func(key string, value *Player)) {
	rm.Lock()
	for k, v := range rm.internal {
		f(k, v)
	}
	rm.Unlock()
}
