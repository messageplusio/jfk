package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	_ "embed"
	"log/slog"

	"github.com/google/uuid"
)

//go:embed jokes.txt
var jokesText string

type Joke struct {
	Part1 string
	Part2 string
}

var Jokes []Joke

func init() {
	slog.Info("Loading jokes")
	defer slog.Info("Jokes loaded")
	// Load jokes from the embedded file
	// Split the file into lines
	// Each joke has two lines
	lines := strings.Split(jokesText, "\n")
	for i := 0; i < len(lines); i += 2 {
		if len(lines) > i+1 {
			Jokes = append(Jokes, Joke{
				Part1: lines[i],
				Part2: lines[i+1],
			})
		}
	}
}

//go:embed index.html
var indexHTML string

func main() {
	slog.Info("Starting the server")
	defer slog.Info("Server stopped")
	// Serve the HTML page
	http.HandleFunc("/jokes", func(w http.ResponseWriter, r *http.Request) {
		joke := Jokes[time.Now().Second()%len(Jokes)]
		fmt.Fprintf(w, indexHTML, uuid.New().String(), fmt.Sprintf("%s<br>%s", joke.Part1, joke.Part2))
	})

	// Endpoint to fetch the current time
	http.HandleFunc("/joke", func(w http.ResponseWriter, r *http.Request) {
		// Get the current time
		joke := Jokes[time.Now().Second()%len(Jokes)]
		// Print jokes in two lines
		fmt.Fprintf(w, "%s<br>%s", joke.Part1, joke.Part2)
	})

	http.HandleFunc("/.well-known/acme-challenge/uCQXlP5kVBZm58MdAIf5sAotGlUeZPjxobpibkG0XBk", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "uCQXlP5kVBZm58MdAIf5sAotGlUeZPjxobpibkG0XBk.d8WdiVqQsDwkW4ZhxpaCsuZCL8-cna9LHpgrNZCR0eM")
	})

	// Start the server
	fmt.Println("Server starting on http://0.0.0.0/", os.Getenv("NAME"))
	if err := http.ListenAndServe(":80", nil); err != nil {
		slog.Error(err.Error())
	}
}
