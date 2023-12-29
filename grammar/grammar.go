package grammar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// OpenAI default api url as variable to allow mocking of the api.
var OPENAI_API_URL = "https://api.openai.com/v1/chat/completions"

// some defaults
const OPENAI_API_KEY = "OPENAI_API_KEY"
const OPENAI_API_MODEL = "gpt-4"
const OPENAI_API_ROLE = "user"
const OPENAI_PROMPT_EN = "Correct the grammar of the following text: "
const OPENAI_PROMPT_DE = "Korrigiere die Grammatik des folgenden Textes: "
const OPENAI_MOCK = "OPENAI_MOCK"

// support promt in German and English language
const (
	DE Language = iota
	EN
)

// the language type (DE, EN)
type Language int

// GrammarCorrector is the interface that wraps the CorrectGrammar method.
type GrammarCorrector interface {
	CorrectGrammar(apiKey string, lang Language, text string) (string, error)
}

// OpenAI api structure
type OpenAIChoice struct {
	Index        int              `json:"index"`
	Message      OpenAIMessage    `json:"message"`
	Logprobs     *json.RawMessage `json:"logprobs"` // Using RawMessage for potential null value
	FinishReason string           `json:"finish_reason"`
}

// OpenAI api structure
type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// OpenAI response api structure
type OpenAIResponse struct {
	ID                string           `json:"id"`
	Object            string           `json:"object"`
	Created           int64            `json:"created"`
	Model             string           `json:"model"`
	Choices           []OpenAIChoice   `json:"choices"`
	Usage             OpenAIUsage      `json:"usage"`
	SystemFingerprint *json.RawMessage `json:"system_fingerprint"` // Using RawMessage for potential null value
}

// OpenAI api structure
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAI request api structure
type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

// GPT implementation:
type MockGrammarCorrector struct{}

// CorrectGrammar takes an API key, a language and text, and returns the grammar-corrected text
func (gc MockGrammarCorrector) CorrectGrammar(apiKey string, lang Language, text string) (string, error) {
	switch lang {
	case DE:
		return "korrigiert: " + text, nil
	case EN:
		return "corrected: " + text, nil
	default:
		return "", fmt.Errorf("language not supported: %v", lang)
	}
}

// GPT implementation:
type GPTGrammarCorrector struct{}

// CorrectGrammar takes an API key, a language and text, and returns the grammar-corrected text
func (gc GPTGrammarCorrector) CorrectGrammar(apiKey string, lang Language, text string) (string, error) {
	var prompt string

	// Build a ChatGPT prompt
	switch lang {
	case EN:
		prompt = OPENAI_PROMPT_EN
	case DE:
		prompt = OPENAI_PROMPT_DE
	default:
		return "", fmt.Errorf("language not supported: %v", lang)
	}

	// Construct the request body
	request := OpenAIRequest{
		Model: OPENAI_API_MODEL,
		Messages: []OpenAIMessage{
			{
				Role:    OPENAI_API_ROLE,
				Content: prompt,
			},
			{
				Role:    OPENAI_API_ROLE,
				Content: text,
			},
		},
	}

	log.Printf("prompt: %v", request)

	byteData, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	reader := bytes.NewReader(byteData)

	// Create an HTTP request
	req, err := http.NewRequest("POST", OPENAI_API_URL, reader)
	if err != nil {
		return "", err
	}

	// Set the required headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("request: %v", req)
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the response
	var chatResponse OpenAIResponse
	if err := json.Unmarshal(body, &chatResponse); err != nil {
		return "", err
	}
	log.Printf("response: %v", chatResponse)

	// Check HTTP response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status code: %s (%v)", resp.Status, chatResponse)
	}

	defer resp.Body.Close()

	if len(chatResponse.Choices) > 0 {
		return chatResponse.Choices[0].Message.Content, nil
	}

	return "", nil
}
