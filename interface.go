package anyfunc

import "context"

type ClientInterface interface {
	Call(ctx context.Context, prompt Prompt) (map[string]any, error)
}
