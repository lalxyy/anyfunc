package anyfunc

import "context"

type ClientInterface interface {
	CallForJSON(context.Context, Prompt) (string, error)
	Call(context.Context, Prompt) (map[string]any, error)
}
