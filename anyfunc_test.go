package anyfunc

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

var apiKey string

func TestBasicFunctionality(t *testing.T) {
	ctx := context.Background()

	client := NewClient(apiKey)
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
	slog.Info("TestBasicFunctionality passed")
}

func TestMain(m *testing.M) {
	// Try to read API key from file if it exists.
	apiKeyBytes, err := os.ReadFile("api_key")
	if err != nil {
		slog.Info("Error reading api_key file, skipping", "error", err)
	} else {
		apiKey = string(apiKeyBytes)
	}
	// If API key is still not set, try to read from environment variable.
	if apiKey == "" {
		apiKey = os.Getenv("API_KEY")
	}
	if apiKey == "" {
		slog.Error("API key not set. Please set the API_KEY environment variable or create an api_key file.")
		os.Exit(1)
	}
	os.Exit(m.Run())
}
