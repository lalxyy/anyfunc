package anythingfunction

import "context"

type ClientInterface interface {
	Run(ctx context.Context, prompt Prompt) (map[string]any, error)
}
