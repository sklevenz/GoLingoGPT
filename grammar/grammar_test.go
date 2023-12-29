package grammar

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var apiKey string
var mockServer *httptest.Server

const jsonResponse = `{
	"id": "chatcmpl-8a6BFWm1yk2eohvtBmvxMhsdslHgy",
	"object": "chat.completion",
	"created": 1703614257,
	"model": "gpt-4-0613",
	"choices": [
	  {
		"index": 0,
		"message": {
		  "role": "assistant",
		  "content": "xyz"
		},
		"logprobs": null,
		"finish_reason": "stop"
	  }
	],
	"usage": {
	  "prompt_tokens": 37,
	  "completion_tokens": 12,
	  "total_tokens": 49
	},
	"system_fingerprint": null
  }`

func init() {
	apiKey = os.Getenv(OPENAI_API_KEY)
	if apiKey == "" {
		log.Fatalf("Environment variable %s not set!", OPENAI_API_KEY)
	}

	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, jsonResponse)
	}))

	OPENAI_API_URL = mockServer.URL
}

func TestCorrectGrammarEnglish(t *testing.T) {

	text := "abc"

	correctedText, err := GPTGrammarCorrector{}.CorrectGrammar(apiKey, EN, text)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	if correctedText != "xyz" {
		t.Errorf("wrong correction, expected: 'xyz', was: '%s'", correctedText)
	}
}

func TestCorrectGrammarGerman(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, jsonResponse)
	}))
	OPENAI_API_URL = server.URL
	defer server.Close()

	text := "abc"

	correctedText, err := GPTGrammarCorrector{}.CorrectGrammar(apiKey, DE, text)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	if correctedText != "xyz" {
		t.Errorf("wrong correction, expected: 'xyz', was: '%s'", correctedText)
	}
}

func TestCorrectGrammarUnsupportedLanguage(t *testing.T) {
	text := "Unsupported Language."

	_, err := GPTGrammarCorrector{}.CorrectGrammar(apiKey, 99, text)

	if err == nil {
		t.Errorf("error expected: %v", err)
	}
	if !strings.Contains(err.Error(), "language not supported") {
		t.Errorf("wrong error message, expected='language not supported: fr' was '%v'", err.Error())
	}
}

func TestCorrectGrammarMock(t *testing.T) {
	text := "Test."

	correctedText, err := MockGrammarCorrector{}.CorrectGrammar(apiKey, EN, text)

	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	if correctedText != "corrected: Test." {
		t.Errorf("wrong correction, expected: 'corrected: Test.', was: '%s'", correctedText)
	}
}
