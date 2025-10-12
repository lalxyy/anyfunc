package main

import (
	"context"
	"log/slog"

	"github.com/lalxyy/anyfunc"
)

const apiKey = "YOUR_API_KEY"

func main() {
	ctx := context.Background()
	client, err := anyfunc.NewClient(anyfunc.BackendGemini, apiKey)
	if err != nil {
		slog.Error("Error creating client", "error", err)
		return
	}
	prompt := anyfunc.Prompt{
		Description: "Return the greatest common factor of given two numbers `num1` and `num2`.",
		Parameters: map[string]any{
			"num1": 45,
			"num2": 60,
		},
	}
	response, err := client.Call(ctx, prompt)
	slog.Info("Response", "data", response, "error", err)
}
