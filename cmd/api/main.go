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
	port       = "8080"
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

	address := ":" + port
	app.logger.Info("Server running", "port", port)
	http.ListenAndServe(address, nil)
}
