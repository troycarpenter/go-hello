package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from k3s Go CI/CD pipeline - V2\n")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

