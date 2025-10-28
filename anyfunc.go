package anyfunc

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"google.golang.org/genai"
)

const defaultSystemPrompt = `You are an AI assistant that helps people finish math
calculation and / or text transformation tasks. You are given a description
of a task and a set of parameters in JSON format.
You should generate a response that fulfills the task using the provided
parameters. Your response should be in JSON format, do not use Markdown formatting,
and should only contain the result of the task, along with a field called
"successful" that is set to boolean value true. If the result field is not
clearly named, put the result in the exact field named "result". If you cannot
complete the task, respond with an appropriate error message in JSON format setting
the "successful" value to false, and "error" field containing the error message.`

type Backend int

const (
	BackendOpenAI Backend = iota
	BackendGemini
)

type Client struct {
	backend      Backend
	openAIClient openai.Client
	geminiClient *genai.Client
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
func NewClient(backend Backend, apiKey string) (*Client, error) {
	ctx := context.Background()
	switch backend {
	case BackendOpenAI:
		// Initialize OpenAI client
		openAIClient := openai.NewClient(
			option.WithAPIKey(apiKey),
		)
		return &Client{
			backend:      BackendOpenAI,
			openAIClient: openAIClient,
			systemPrompt: defaultSystemPrompt,
		}, nil
	case BackendGemini:
		// Initialize Gemini client
		geminiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: apiKey,
		})
		if err != nil {
			return nil, err
		}
		return &Client{
			backend:      BackendGemini,
			geminiClient: geminiClient,
			systemPrompt: defaultSystemPrompt,
		}, nil
	}
	return nil, errors.New("unsupported backend")
}

// Call processes the given prompt using the OpenAI API and returns the
// response.
func (c *Client) Call(ctx context.Context, prompt Prompt) (map[string]any, error) {
	parameterJSON, err := json.MarshalIndent(prompt.Parameters, "", "  ")
	if err != nil {
		return nil, err
	}
	slog.Debug("Given parameters", "parameters", string(parameterJSON))

	var response string
	switch c.backend {
	case BackendOpenAI:
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
		response = chatCompletion.Choices[0].Message.Content
	case BackendGemini:
		result, err := c.geminiClient.Models.GenerateContent(
			ctx,
			"gemini-2.5-flash",
			genai.Text(prompt.Description+"\n\n"+string(parameterJSON)),
			&genai.GenerateContentConfig{
				SystemInstruction: genai.NewContentFromText(c.systemPrompt, genai.RoleModel),
			},
		)
		if err != nil {
			return nil, err
		}
		response = result.Text()
	}
	slog.Debug("Raw response", "response", response)

	var result map[string]any
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, err
	}
	return result, nil
}
