package main

import (
	"fmt"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("receiving request")
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	fmt.Println("returning count")
}

func main() {
	http.HandleFunc("GET /", handleRequest)
	http.HandleFunc("GET /count", handleCount)

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
