package main

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

const testFileName string = "counter_storage_test.gob"

var repo *FileStorage

func TestMain(m *testing.M) {
	repo = NewFileStorage(testFileName)

	// Run tests
	code := m.Run()

	// remove temprary test file
	os.Remove(testFileName)

	// Exit with the test's exit code
	os.Exit(code)
}

func TestPeriodicCleanAndSave(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	wc := NewWindowCounter(repo, time.Second*2, logger)
	defer wc.Close()

	// increment the count
	wc.IncrementCount()
	wc.IncrementCount()

	reqCount := 0
	for _, count := range wc.requets {
		reqCount += count
	}

	if reqCount != 2 {
		// There should be exactly 2 request in the requests map
		t.Fatalf("expected 2 request, got %d", reqCount)
	}

	// Adding a delay of  second
	time.Sleep(3 * time.Second)

	count, err := wc.Count()
	if count != 0 {
		t.Fatalf("expected 0 request, got %d", count)
	}
	if err != nil {
		t.Fatalf("unexpected error from Count: %v", err)
	}
}
