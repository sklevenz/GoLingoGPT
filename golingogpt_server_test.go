package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sklevenz/GoLingoGPT/grammar"
)

func init() {
	grammerCorrector = grammar.MockGrammarCorrector{}
}

func TestDetermineLanguageDE(t *testing.T) {
	req, err := http.NewRequest("GET", "/correctText?text=example", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Language", "de")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRequest)
	handler.ServeHTTP(rr, req)

	if language := rr.Header().Get("Content-Language"); language != "de" {
		t.Errorf("wrong language: %v", language)
	}
}
func TestDetermineLanguageEn(t *testing.T) {
	req, err := http.NewRequest("GET", "/correctText?text=example", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Language", "en")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRequest)
	handler.ServeHTTP(rr, req)

	if language := rr.Header().Get("Content-Language"); language != "en" {
		t.Errorf("wrong language: %v", language)
	}
}
func TestDetermineLanguageDefault(t *testing.T) {
	req, err := http.NewRequest("GET", "/correctText?text=example", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRequest)
	handler.ServeHTTP(rr, req)

	if language := rr.Header().Get("Content-Language"); language != "en" {
		t.Errorf("wrong language: %v", language)
	}
}
func TestDetermineLanguageUnsupported(t *testing.T) {
	req, err := http.NewRequest("GET", "/correctText?text=example", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Language", "fr")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRequest)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}

func TestHandleGetRequest(t *testing.T) {
	// Set up the request and response recorder
	req, err := http.NewRequest("GET", "/correctText?text=example", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Language", "en")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRequest)
	handler.ServeHTTP(rr, req)

	// Check the status code and response body
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := `corrected: example`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandlePostRequest(t *testing.T) {
	// Create a request with example text
	requestBody := bytes.NewBufferString("example")
	req, err := http.NewRequest("POST", "/correctTextPost", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Language", "en")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRequest)

	// Perform the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "corrected: example"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
