package main

import (
	"sync"
	"time"
)

const (
	USER_EXPIRE = 2 * time.Hour
	COOKIE_NAME = "session"
)

type TagState struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	State bool   `json:"state"`
}

type TagMap struct {
	value map[string]*TagState
	lock  sync.RWMutex
}

func NewTagMap() *TagMap {
	var tm TagMap
	tm.value = make(map[string]*TagState)
	return &tm
}

func (t *TagMap) Off(tags ...string) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	for i := 0; i < len(tags); i++ {
		ts, exist := t.value[tags[i]]
		if exist {
			ts.State = false
		}
	}
	return false
}

func (t *TagMap) On(tags ...string) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	for i := 0; i < len(tags); i++ {
		ts, exist := t.value[tags[i]]
		if exist {
			ts.State = true
		}
	}
	return false
}

func (t *TagMap) Update(mc map[string]int) {
	t.lock.Lock()
	defer t.lock.Unlock()
	for k, v := range mc {
		ts, exist := t.value[k]
		if !exist {
			t.value[k] = &TagState{
				k, v, true,
			}
		} else {
			ts.Count = v
		}
	}
	var names []string
	for k, _ := range t.value {
		_, exist := mc[k]
		if !exist {
			names = append(names, k)
		}
	}
	for _, name := range names {
		delete(t.value, name)
	}
}

func (t *TagMap) List() []string {
	t.lock.RLock()
	defer t.lock.RUnlock()
	ret := make([]string, len(t.value))
	var i = 0
	for k, _ := range t.value {
		ret[i] = k
		i++
	}
	return ret
}

func (t *TagMap) Get(tag string) *TagState {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.value[tag]
}

func (t *TagMap) Slice() []*TagState {
	t.lock.RLock()
	defer t.lock.RUnlock()

	ret := make([]*TagState, len(t.value))
	var i = 0
	for _, v := range t.value {
		nt := new(TagState)
		*nt = *v
		ret[i] = nt
		i++
	}
	return ret
}
