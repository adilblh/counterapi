package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	c "github.com/adilblh/counterapi/counter"
)

const windowSize = time.Second * 60

type Server struct {
	windowCounter *c.WindowCounter
}

func NewServer(windowCounter *c.WindowCounter) *Server {
	return &Server{windowCounter: windowCounter}
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request received ...")
	s.windowCounter.IncrementCount()
}

func (s *Server) handleCount(w http.ResponseWriter, r *http.Request) {
	reqCount := s.windowCounter.Count()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"count": reqCount})
}

func main() {
	windowCounter := c.NewWindowCounter(windowSize)
	s := NewServer(windowCounter)

	http.HandleFunc("GET /", s.handleRequest)
	http.HandleFunc("GET /count", s.handleCount)

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
