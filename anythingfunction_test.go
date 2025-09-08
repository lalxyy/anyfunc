package anythingfunction

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

func TestBasicFunctionality(t *testing.T) {
	ctx := context.Background()

	apiKeyBytes, err := os.ReadFile("api_key")
	if err != nil {
		t.Fatalf("Error reading API key: %v", err)
	}
	apiKey := string(apiKeyBytes)

	client := NewClient(apiKey, nil)
	prompt := Prompt{
		Description: "Return the greatest common factor of given two numbers `num1` and `num2`.",
		Parameters: map[string]any{
			"num1": 45,
			"num2": 60,
		},
	}
	response, err := client.Run(ctx, prompt)
	if err != nil {
		t.Fatalf("Error running prompt: %v", err)
	}
	slog.Info("Response", "data", response)
}
