package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/sklevenz/GoLingoGPT/grammar"
)

var apiKey string
var grammerCorrector grammar.GrammarCorrector = grammar.GPTGrammarCorrector{}

func init() {
	apiKey = os.Getenv(grammar.OPENAI_API_KEY)
	if apiKey == "" {
		log.Fatalf("Environment variable %s not set!", grammar.OPENAI_API_KEY)
	}
	if os.Getenv(grammar.OPENAI_MOCK) == "true" {
		grammerCorrector = grammar.MockGrammarCorrector{}
		log.Println("server in mock mode")
	}

}

func determineLanguage(w http.ResponseWriter, r *http.Request) (grammar.Language, error) {
	language := r.Header.Get("Content-Language")

	// Handle text based on language
	switch language {
	case "de", "de-DE":
		w.Header().Add("Content-Language", language)
		return grammar.DE, nil
	case "", "en", "en-US", "en-GB":
		if language == "" {
			language = "en"
		}
		w.Header().Add("Content-Language", language)
		return grammar.EN, nil
	default:
		return 0, fmt.Errorf("unsupported language: %v", language)
	}
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")

	lang, err := determineLanguage(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	correctedText, err := grammerCorrector.CorrectGrammar(apiKey, lang, text)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error correcting grammar: ", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(correctedText))
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error reading request body", http.StatusInternalServerError)
		return
	}
	lang, err := determineLanguage(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	correctedText, err := grammerCorrector.CorrectGrammar(apiKey, lang, string(body))
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error correcting grammar: ", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(correctedText))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetRequest(w, r)
	case http.MethodPost:
		handlePostRequest(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/correctText", handleRequest)

	log.Printf("Start servcer on port 8080 in OPENAI_MOCK mode '%v'", os.Getenv(grammar.OPENAI_MOCK))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
