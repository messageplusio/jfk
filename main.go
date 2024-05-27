package main

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
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

//go:embed templates/*
var htmlroot embed.FS

func CreateFileWithBase64(fileName string, encodedData string) error {
	// Create a new file
	file, err := os.Create(fileName)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer file.Close()

	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	// Write the certificate to the file
	_, err = file.Write(decodedData)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func CreatePemFile() error {
	return CreateFileWithBase64("cert.pem", os.Getenv("CERT"))
}

func CreateKeyFile() error {
	return CreateFileWithBase64("privkey.pem", os.Getenv("KEY"))
}

func main() {
	slog.Info("Starting the server")
	defer slog.Info("Server stopped")
	// Serve the HTML page
	http.HandleFunc("/jokes", func(w http.ResponseWriter, r *http.Request) {
		joke := Jokes[time.Now().Second()%len(Jokes)]
		jokes, err := htmlroot.ReadFile("templates/jokes.html")
		if err != nil {
			slog.Error(err.Error())
			return
		}
		fmt.Fprintf(w, string(jokes), uuid.New().String(), fmt.Sprintf("%s<br>%s", joke.Part1, joke.Part2))
	})

	// Serve the HTML page
	http.HandleFunc("/gotemplate", func(w http.ResponseWriter, r *http.Request) {
		gotemplates, err := htmlroot.ReadFile("templates/gotemplates.html")
		if err != nil {
			slog.Error(err.Error())
			return
		}
		fmt.Fprint(w, string(gotemplates))
	})

	http.HandleFunc("/render", handleRender)

	// Endpoint to fetch the current time
	http.HandleFunc("/joke", func(w http.ResponseWriter, r *http.Request) {
		// Get the current time
		joke := Jokes[time.Now().Second()%len(Jokes)]
		// Print jokes in two lines
		fmt.Fprintf(w, "%s<br>%s", joke.Part1, joke.Part2)
	})

	if err := CreateKeyFile(); err != nil {
		slog.Error(err.Error())

	}
	if err := CreatePemFile(); err != nil {
		slog.Error(err.Error())
	}

	if err := http.ListenAndServeTLS(":443", "cert.pem", "privkey.pem", nil); err != nil {
		slog.Error(err.Error())
	}
}

func handleRender(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	templateStr := r.FormValue("template")
	jsonStr := r.FormValue("jsonData")

	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	tmpl, err := template.New("result").Parse(templateStr)
	if err != nil {
		http.Error(w, "Invalid template", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
