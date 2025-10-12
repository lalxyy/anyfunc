package anyfunc

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

// Struct to hold the API key for tests.
type apiKeyBundle struct {
	OpenAIKey string `yaml:"openAI"`
	GeminiKey string `yaml:"gemini"`
}

var apiKeys apiKeyBundle

func TestBasicFunctionality(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		backend Backend
		apiKey  string
	}{
		{
			name:    "OpenAI Backend",
			backend: BackendOpenAI,
			apiKey:  apiKeys.OpenAIKey,
		},
		{
			name:    "Gemini Backend",
			backend: BackendGemini,
			apiKey:  apiKeys.GeminiKey,
		},
		// Add more backends as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.backend, tt.apiKey)
			if err != nil {
				t.Fatalf("Error creating client: %v", err)
			}
			prompt := Prompt{
				Description: "Return the greatest common factor of given two numbers `num1` and `num2`.",
				Parameters: map[string]any{
					"num1": 45,
					"num2": 60,
				},
			}
			response, err := client.Call(ctx, prompt)
			if err != nil {
				t.Fatalf("Error running prompt: %v", err)
			}
			returnValue, ok := response["result"].(float64)
			if !ok {
				t.Fatalf("Expected result to be an integer, got %T", response["result"])
			}
			expectedValue := 15.0
			if returnValue != expectedValue {
				t.Fatalf("Expected result to be %f, got %f", expectedValue, returnValue)
			}
		})
	}

	slog.Info("TestBasicFunctionality passed")
}

func TestMain(m *testing.M) {
	// Try to read API key from file if it exists.
	apiKeyBytes, err := os.ReadFile("api_key.yaml")
	if err != nil {
		slog.Info("Error reading api_key file, skipping", "error", err)
	} else {
		err := yaml.Unmarshal(apiKeyBytes, &apiKeys)
		if err != nil {
			slog.Info("Error parsing api_key file, skipping", "error", err)
		}
	}
	// If API key is still not set, try to read from environment variable.
	if apiKeys.OpenAIKey == "" {
		apiKeys.OpenAIKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKeys.GeminiKey == "" {
		apiKeys.GeminiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKeys.OpenAIKey == "" || apiKeys.GeminiKey == "" {
		slog.Error("API key not set. Please set the API_KEY environment variable or create an api_key file.")
		os.Exit(1)
	}
	os.Exit(m.Run())
}
