package timeout

import (
	"sync"
	"time"
)

const userTimeout = time.Duration(10) * time.Second

type Timeout struct {
	mx     sync.RWMutex
	waiter map[int64]struct{}
}

func Make() *Timeout {
	return &Timeout{
		waiter: make(map[int64]struct{}),
	}
}

func (t *Timeout) Exist(id int64) bool {
	t.mx.RLock()
	defer t.mx.RUnlock()
	_, ok := t.waiter[id]
	return ok
}

func (t *Timeout) SetTimer(id int64) {
	t.mx.Lock()
	t.waiter[id] = struct{}{}
	t.mx.Unlock()

	time.AfterFunc(userTimeout, func() {
		t.mx.Lock()
		delete(t.waiter, id)
		t.mx.Unlock()
		//log.Printf("tout del (%d)", id)
	})

	//log.Printf("tout set (%d)", id)
}
