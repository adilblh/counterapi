package main

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type WindowCounter struct {
	requets    map[int64]int
	mu         sync.RWMutex
	windowSize time.Duration
	repo       CounterRepository
	stopChan   chan struct{}
	logger     *slog.Logger
}

func NewWindowCounter(repository CounterRepository, ws time.Duration, logger *slog.Logger) *WindowCounter {

	savedRequests, err := repository.Get()
	if err != nil {
		savedRequests = make(map[int64]int)
	}

	wc := &WindowCounter{
		repo:       repository,
		requets:    savedRequests,
		windowSize: ws,
		stopChan:   make(chan struct{}),
		logger:     logger,
	}

	// clean old requests every 60s & update the file storage
	go wc.periodicCleanAndSave()

	return wc
}

func (wc *WindowCounter) periodicCleanAndSave() {
	ticker := time.NewTicker(wc.windowSize)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			wc.mu.Lock()
			wc.cleanOldRequests()
			wc.saveToFile()
			wc.mu.Unlock()

		case <-wc.stopChan:
			return
		}
	}
}

func (wc *WindowCounter) cleanOldRequests() {
	reqs, err := wc.repo.Get()

	// logging the map to verify data consistency
	mapAsString := fmt.Sprintf("%v", reqs)
	wc.logger.Info("clean triggered", slog.String("reqs map", mapAsString))

	if err != nil {
		wc.logger.Error(err.Error())
	}

	for timeStamp := range reqs {
		if time.Unix(timeStamp, 0).Add(wc.windowSize).Before(time.Now()) {
			delete(wc.requets, timeStamp)
		}
	}
}

func (wc *WindowCounter) IncrementCount() {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	currentTime := time.Now().Unix()
	wc.requets[currentTime]++
	wc.saveToFile()
}

func (wc *WindowCounter) Count() (int, error) {
	wc.mu.RLock()
	defer wc.mu.RUnlock()

	reqs, err := wc.repo.Get()
	if err != nil {
		return -1, err
	}

	reqCount := 0
	for _, count := range reqs {
		reqCount += count
	}

	return reqCount, nil
}

func (wc *WindowCounter) saveToFile() {
	if err := wc.repo.Save(wc.requets); err != nil {
		wc.logger.Error(err.Error())
	}
}

func (wc *WindowCounter) Close() {
	close(wc.stopChan)
	if err := wc.repo.Save(wc.requets); err != nil {
		wc.logger.Error(err.Error())
	}
}
