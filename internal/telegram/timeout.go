package telegram

import (
	"sync"
	"time"
)

type TimeStamp struct {
	mx    sync.RWMutex
	stamp map[int64]int64
}

func MakeStamp() *TimeStamp {
	return &TimeStamp{
		stamp: make(map[int64]int64),
	}
}

func (t *TimeStamp) SetStamp(key, val int64) {
	t.mx.Lock()
	t.stamp[key] = val
	t.mx.Unlock()
	//log.Printf("SetStamp %d", val)
}

func (t *TimeStamp) SetStampNow(key int64) {
	t.mx.Lock()
	t.stamp[key] = time.Now().Unix()
	t.mx.Unlock()
}

func (t *TimeStamp) GetStamp(key int64) (int64, bool) {
	t.mx.RLock()
	defer t.mx.RUnlock()
	val, ok := t.stamp[key]
	return val, ok
}

func (t *TimeStamp) CheckTimeout(id int64) bool {
	curr_time := time.Now().Unix()
	if tstamp, ok := t.GetStamp(id); ok && 0 < tstamp {
		//log.Printf("CheckTimeout %d, pass %d", tstamp, (curr_time - tstamp))
		if (curr_time - tstamp) < 10 {
			return true
		}
	}
	t.SetStamp(id, curr_time)
	return false
}

func (t *TimeStamp) DeleteTimeouts() {
	for {
		curr_time := time.Now().Unix()
		userIDs := make([]int64, 0)

		t.mx.RLock()
		for id, tstamp := range t.stamp {
			if 0 < tstamp && 20 < (curr_time-tstamp) {
				userIDs = append(userIDs, id)
			}
		}
		t.mx.RUnlock()

		for _, id := range userIDs {
			if 0 < id {
				t.mx.Lock()
				delete(t.stamp, id)
				t.mx.Unlock()
				//log.Printf("Delete stamp id %d", id)
			}
		}
		userIDs = nil

		time.Sleep(10 * time.Second)
	}
}
