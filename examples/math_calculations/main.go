package main

import (
	"context"
	"log/slog"

	anythingfunction "github.com/lalxyy/anything-function"
)

func main() {
	ctx := context.Background()
	client := anythingfunction.NewClient("YOUR_API_KEY")
	prompt := anythingfunction.Prompt{
		Description: "Return the greatest common factor of given two numbers `num1` and `num2`.",
		Parameters: map[string]any{
			"num1": 45,
			"num2": 60,
		},
	}
	response, err := client.Run(ctx, prompt)
	slog.Info("Response", "data", response, "error", err)
}
