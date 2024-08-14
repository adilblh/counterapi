package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	windowSize = time.Second * 60
	dataFile   = "counter_storage.gob"
)

type Application struct {
	wc     *WindowCounter
	logger *slog.Logger
}

func NewApplication(wincowCounter *WindowCounter, logger *slog.Logger) *Application {
	return &Application{wc: wincowCounter, logger: logger}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repository := NewFileStorage(dataFile)

	windowCounter := NewWindowCounter(repository, windowSize, logger)
	defer windowCounter.Close()

	app := NewApplication(windowCounter, logger)

	http.HandleFunc("GET /", app.HandleRequest)
	http.HandleFunc("GET /count", app.HandleCount)

	logger.Info("Server running on port :8080")
	http.ListenAndServe(":8080", nil)
}
