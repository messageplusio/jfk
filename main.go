package main

import (
	"fmt"
	"net/http"

	"time"

	"github.com/google/uuid"
)

func htmx() string {
	return `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Current Time</title>
		<!-- Import HTMX -->
		<script src="https://unpkg.com/htmx.org"></script>
		<style>
			/* Fullscreen styling for the div */
			body, html {
				margin: 0;
				padding: 0;
				height: 100%;
			}
			div#fullpage {
				display: flex;
				justify-content: center;
				align-items: center;
				width: 100vw;
				height: 100vh;
				background: black; /* Dark background for better visibility */
				color: white; /* White text color */
				font-size: 4em; /* Large font size for better visibility */
				font-family: 'Arial', sans-serif; /* Stylish, readable font */
				overflow: hidden; /* Hide overflow */
			}
		</style>
	</head>
	<body>
		<div id="fullpage" hx-get="/time" hx-trigger="every 1s" hx-swap="innerHTML">
			<!-- Initial content can go here if needed -->
		</div>
	</body>
	</html>
	`
}

func main() {
	// Serve the HTML page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, htmx(), uuid.New().String())
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
