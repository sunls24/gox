package notifier

import (
	"sync"
	"sync/atomic"
)

var (
	m            sync.Mutex
	nextID       atomic.Uint32
	waitersByKey = map[string]map[uint32]chan any{}
)

func Wait(key string) (<-chan any, func()) {
	ch := make(chan any, 1)
	id := nextID.Add(1)

	m.Lock()
	defer m.Unlock()
	if waitersByKey[key] == nil {
		waitersByKey[key] = map[uint32]chan any{}
	}
	waitersByKey[key][id] = ch
	return ch, func() {
		m.Lock()
		defer m.Unlock()

		waiters, ok := waitersByKey[key]
		if !ok {
			return
		}
		delete(waiters, id)
		if len(waiters) == 0 {
			delete(waitersByKey, key)
		}
	}
}

func Notify(key string, value any) {
	m.Lock()
	waiters, ok := waitersByKey[key]
	if ok {
		delete(waitersByKey, key)
	}
	m.Unlock()
	if !ok {
		return
	}
	for _, ch := range waiters {
		ch <- value
		close(ch)
	}
}
