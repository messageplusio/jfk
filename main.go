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

func CreatePemFile() {
	// Create a new file
	file, err := os.Create("cert.pem")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer file.Close()

	// Write the certificate to the file
	_, err = file.WriteString(os.Getenv("CERT"))
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func CreateKeyFile() {
	// Create a new file
	file, err := os.Create("privkey.pem")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer file.Close()

	// Write the private key to the file
	_, err = file.WriteString(os.Getenv("KEY"))
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

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

	CreateKeyFile()
	CreatePemFile()

	if err := http.ListenAndServeTLS(":443", "cert.pem", "privkey.pem", nil); err != nil {
		slog.Error(err.Error())
	}
}
