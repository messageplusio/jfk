package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		name := os.Getenv("NAME")
		fmt.Fprintf(w, "[%s]:Hello, you've requested: %s\n Name:%s", id, r.URL.Path, name)
	})

	fmt.Println("Server is starting...")
	http.ListenAndServe(":8080", nil)
}
