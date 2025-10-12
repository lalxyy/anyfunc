# Anything Function

Anything Function is a Go library that allows calling LLMs (like ChatGPT, Gemini, Claude, ...) like functions and RPCs.

This project aims to make it easier for developers to integrate LLM capabilities with existing code in Go.

## Examples

```go
import "github.com/lalxyy/anyfunc"

// Create an API key from OpenAI platform and replace the placeholder.
client := anyfunc.NewClient(anyfunc.BackendGemini, "YOUR_API_KEY") // Or `anyfunc.BackendOpenAI`
prompt := anyfunc.Prompt{
  Description: "Return the greatest common factor of given two numbers `num1` and `num2`.",
  Parameters: map[string]any{
    "num1": 45,
    "num2": 60,
  },
}
response, err := client.Call(ctx, prompt)
// Response would be a map like this:
// map[string]any{
//   "result": 15,
// }
```
