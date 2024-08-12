package counter

import (
	"sync"
	"time"
)

type WindowCounter struct {
	requets    map[int64]int
	mu         sync.Mutex
	windowSize time.Duration
}

func NewWindowCounter(ws time.Duration) *WindowCounter {
	return &WindowCounter{requets: make(map[int64]int), windowSize: ws}
}

func (wc *WindowCounter) IncrementCount() {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	currentTime := time.Now().Unix()
	wc.requets[currentTime]++

	// clean old requests
	wc.cleanOldreqs()
}

func (wc *WindowCounter) cleanOldreqs() {
	for timeStamp := range wc.requets {
		if time.Unix(timeStamp, 0).Add(wc.windowSize).Before(time.Now()) {
			delete(wc.requets, timeStamp)
		}
	}
}

func (wc *WindowCounter) Count() int {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	reqCount := 0
	for timeStamp, count := range wc.requets {
		if time.Unix(timeStamp, 0).Add(wc.windowSize).After(time.Now()) {
			reqCount += count
		}
	}

	return reqCount
}
