package main

import (
	"fmt"
	"net/http"

	"time"

	"github.com/google/uuid"
)

func main() {
	// Serve the HTML page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(w, `
            <!DOCTYPE html>
            <html lang="en">
            <head>
                <meta charset="UTF-8">
                <meta name="viewport" content="width=device-width, initial-scale=1.0">
                <title>Current Time</title>
                <!-- Import HTMX -->
                <script src="https://unpkg.com/htmx.org"></script>
            </head>
            <body>
                <h1>Current Time Example</h1>
                <div id="time" hx-get="/time" hx-trigger="every 1s" hx-swap="outerHTML">
                    %s <br> %s
                </div>
            </body>
            </html>
        `, currentTime, uuid.New().String())
	})

	// Endpoint to fetch the current time
	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprint(w, currentTime)
	})

	// Start the server
	fmt.Println("Server starting on http://0.0.0.0:8080/")
	http.ListenAndServe(":8080", nil)
}
