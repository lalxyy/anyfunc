package anythingfunction

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const defaultSystemPrompt = `You are an AI assistant that helps people finish math
calculation and / or text transformation tasks. You are given a description
of a task and a set of parameters in JSON format.
You should generate a response that fulfills the task using the provided
parameters. Your response should be in JSON format and should only contain the
result of the task, along with a field called "successful" that is set to 
boolean value true. If the result field is not clearly named, put the result
in the "result" field. If you cannot complete the task, respond with an
appropriate error message in JSON format setting the "successful" value to
false, and "error" field containing the error message.`

type Client struct {
	openAIClient openai.Client
	systemPrompt string
}

// Ensure Client implements ClientInterface.
var _ ClientInterface = &Client{}

// Prompt represents a structured prompt with a description and parameters.
type Prompt struct {
	Description string
	Parameters  map[string]any
}

// NewClient initializes and returns a new Client with the provided API key.
func NewClient(apiKey string) *Client {
	openAIClient := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &Client{
		openAIClient: openAIClient,
		systemPrompt: defaultSystemPrompt,
	}
}

// Run processes the given prompt using the OpenAI API and returns the
// response.
func (c *Client) Run(ctx context.Context, prompt Prompt) (map[string]any, error) {
	parameterJSON, err := json.MarshalIndent(prompt.Parameters, "", "  ")
	if err != nil {
		return nil, err
	}
	slog.Debug("Given parameters", "parameters", string(parameterJSON))

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(c.systemPrompt),
			openai.UserMessage(prompt.Description + "\n\n" + string(parameterJSON)),
		},
		Model: openai.ChatModel("gpt-5"),
	}

	chatCompletion, err := c.openAIClient.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, err
	}
	response := chatCompletion.Choices[0].Message.Content
	slog.Debug("Raw response", "response", response)

	var result map[string]any
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, err
	}
	return result, nil
}
