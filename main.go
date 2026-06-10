package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from k3s Go CI/CD pipeline - V6.")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server running on :80")
	http.ListenAndServe(":80", nil)
}

