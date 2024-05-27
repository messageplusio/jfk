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
	http.HandleFunc("/", serveFiles)
	http.HandleFunc("/render", handleTemplateRender)
	http.HandleFunc("/joke", handleJoke)

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

func handleTemplateRender(w http.ResponseWriter, r *http.Request) {
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

func handleJoke(w http.ResponseWriter, r *http.Request) {
	// Get the current time
	joke := Jokes[time.Now().Second()%len(Jokes)]
	// Print jokes in two lines
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%s<br>%s", joke.Part1, joke.Part2)
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "" {
		http.ServeFileFS(w, r, htmlroot, "templates/index.html")
		return
	}
	http.ServeFileFS(w, r, htmlroot, "templates/"+r.URL.Path[1:]+".html")
}
