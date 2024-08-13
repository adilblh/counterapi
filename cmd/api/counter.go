package main

import (
	"sync"
	"time"
)

type WindowCounter struct {
	requets    map[int64]int
	mu         sync.Mutex
	windowSize time.Duration
	repo       CounterRepository
}

func NewWindowCounter(repository CounterRepository, ws time.Duration) *WindowCounter {
	return &WindowCounter{repo: repository, requets: make(map[int64]int), windowSize: ws}
}

func (wc *WindowCounter) IncrementCount() {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	currentTime := time.Now().Unix()
	wc.requets[currentTime]++

	// clean old requests
	wc.cleanOldreqs()
	wc.repo.Save(wc.requets)
}

func (wc *WindowCounter) cleanOldreqs() {
	for timeStamp := range wc.requets {
		if time.Unix(timeStamp, 0).Add(wc.windowSize).Before(time.Now()) {
			delete(wc.requets, timeStamp)
		}
	}
}

func (wc *WindowCounter) Count() (int, error) {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	reqs, err := wc.repo.Get()
	if err != nil {
		return -1, err
	}

	reqCount := 0
	for timeStamp, count := range reqs {
		if time.Unix(timeStamp, 0).Add(wc.windowSize).After(time.Now()) {
			reqCount += count
		}
	}

	return reqCount, nil
}
