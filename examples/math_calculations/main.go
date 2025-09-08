package main

import (
	"context"
	"log/slog"

	"github.com/lalxyy/anyfunc"
)

func main() {
	ctx := context.Background()
	client := anyfunc.NewClient("YOUR_API_KEY")
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
