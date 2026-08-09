package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type GroqClient struct {
	client *openai.Client
	model  string
}

func NewGroqClient() (*GroqClient, error) {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Env loading err!")
	}

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Fatal("Groq api key not set!")
	}

	model := os.Getenv("GROQ_MODEL")
	if model == "" {
		model = "openai/gpt-oss-20b"
	}

	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("GROQ_API_KEY")),
		option.WithBaseURL("https://api.groq.com/openai/v1/"),
	)

	return &GroqClient{
		client: &client,
		model:  model,
	}, nil
}

func (g *GroqClient) GenerateResponse(history []Message, tools []ToolDef) (LLMRes, error) {

	var sdkTools []openai.ChatCompletionToolUnionParam

	for _, t := range tools {
		properties := make(map[string]interface{})

		for pName, pDetail := range t.Params {
			properties[pName] = map[string]interface{}{
				"type":        "string",
				"description": pDetail.Description,
			}
		}

		params := openai.FunctionParameters{
			"type":       "object",
			"properties": properties,
		}

		if len(t.Required) > 0 {
			params["required"] = t.Required
		}

		sdkTools = append(sdkTools, openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        t.Name,
			Description: openai.String(t.Description),
			Parameters:  params,
		}))
	}

	var sdkMessages []openai.ChatCompletionMessageParamUnion

	for _, msg := range history {
		switch msg.Role {
		case "user":
			sdkMessages = append(sdkMessages, openai.UserMessage(msg.Content))
		case "assistant":
			if len(msg.ToolCalls) > 0 {
				// Build assistant message with tool calls so the API sees them
				var sdkToolCalls []openai.ChatCompletionMessageToolCallUnionParam
				for _, tc := range msg.ToolCalls {
					argsJSON, _ := json.Marshal(tc.Arguments)
					fnCall := openai.ChatCompletionMessageFunctionToolCallParam{
						ID: tc.ID,
						Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{
							Name:      tc.Function,
							Arguments: string(argsJSON),
						},
					}
					sdkToolCalls = append(sdkToolCalls, openai.ChatCompletionMessageToolCallUnionParam{
						OfFunction: &fnCall,
					})
				}
				assistantMsg := openai.ChatCompletionAssistantMessageParam{
					ToolCalls: sdkToolCalls,
				}
				sdkMessages = append(sdkMessages, openai.ChatCompletionMessageParamUnion{
					OfAssistant: &assistantMsg,
				})
			} else {
				sdkMessages = append(sdkMessages, openai.AssistantMessage(msg.Content))
			}
		case "system":
			sdkMessages = append(sdkMessages, openai.SystemMessage(msg.Content))
		case "tool":
			toolCallID := ""
			if len(msg.ToolCalls) > 0 {
				toolCallID = msg.ToolCalls[0].ID
			}
			sdkMessages = append(sdkMessages, openai.ToolMessage(msg.Content, toolCallID))
		}
	}

	completion, err := g.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: sdkMessages,
		Tools:    sdkTools,
		Model:    openai.ChatModel(g.model),
	})
	if err != nil {
		return LLMRes{}, fmt.Errorf("groq api err : %v", err)
	}

	result := LLMRes{}
	if len(completion.Choices) > 0 {
		choice := completion.Choices[0]
		result.Text = choice.Message.Content

		for _, tc := range choice.Message.ToolCalls {
			var arguments map[string]any
			json.Unmarshal([]byte(tc.Function.Arguments), &arguments)

			result.ToolCalls = append(result.ToolCalls, ToolCall{
				ID:        tc.ID,
				Function:  tc.Function.Name,
				Arguments: arguments,
			})
		}
	}
	return result, nil
}
