package llm

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiClient() (*GeminiClient, error) {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := genai.NewClient(context.Background(), option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal("Some err init gemini client : ", err)
	}

	model := client.GenerativeModel("gemini-3.1-flash-lite-preview")

	return &GeminiClient{client: client, model: model}, nil
}

func (g *GeminiClient) GenerateResponse(history []Message, tools []ToolDef) (LLMRes, error) {

	var declarations []*genai.FunctionDeclaration
	for _, t := range tools {

		properties := make(map[string]*genai.Schema)
		for paramName, paramDetail := range t.Params {
			properties[paramName] = &genai.Schema{
				Type:        genai.TypeString, // Most Docker args are strings
				Description: paramDetail.Description,
			}
		}

		declarations = append(declarations, &genai.FunctionDeclaration{
			Name:        t.Name,
			Description: t.Description,
			Parameters: &genai.Schema{
				Type:       genai.TypeObject,
				Properties: properties,
				Required:   t.Required,
			},
		})
	}
	g.model.Tools = []*genai.Tool{{FunctionDeclarations: declarations}}

	var genaiHistory []*genai.Content
	for i := 0; i < len(history)-1; i++ {
		genaiHistory = append(genaiHistory, &genai.Content{
			Role:  history[i].Role,
			Parts: []genai.Part{genai.Text(history[i].Content)},
		})
	}

	chat := g.model.StartChat()
	chat.History = genaiHistory

	lastMsg := history[len(history)-1]
	fmt.Println("Msg : ", lastMsg)
	resp, err := chat.SendMessage(context.Background(), genai.Text(lastMsg.Content))
	if err != nil {
		return LLMRes{}, err
	}

	res := LLMRes{}
	for _, part := range resp.Candidates[0].Content.Parts {
		switch v := part.(type) {
		case genai.Text:
			res.Text = string(v)
		case genai.FunctionCall:
			res.ToolCalls = append(res.ToolCalls, ToolCall{
				Function:  v.Name,
				Arguments: v.Args,
			})
		}
	}

	return res, nil
}

// // THIS implements the interface from client.go
// func (g *GeminiClient) GenerateResponse(history []Message, tools []ToolDef) (LLMRes, error) {
// 	// A. Map your Tools to genai.FunctionDeclarations
// 	var genaiTools []*genai.Tool
// 	for _, t := range tools {
// 		genaiTool := &genai.Tool{
// 			FunctionDeclarations: []*genai.FunctionDeclaration{
// 				{
// 					Name:        t.Name,
// 					Description: t.Description,
// 				},
// 			},
// 		}
// 		genaiTools = append(genaiTools, genaiTool)
// 	}

// 	// B. Map your History to []*genai.Content
// 	var genaiHistory []*genai.Content
// 	for _, h := range history {
// 		genaiHistory = append(genaiHistory, &genai.Content{
// 			Role:  h.Role,
// 			Parts: []genai.Part{genai.Text(h.Content)},
// 		})
// 	}

// 	// C. Call g.model.GenerateContent(...)
// 	g.model.Tools = genaiTools
// 	resp, err := g.model.GenerateContent(context.Background(), genaiHistory...)
// 	if err != nil {
// 		return LLMRes{}, err
// 	}

// 	// D. Convert the result back to your LLMResponse struct
// res := LLMRes{}
// for _, part := range resp.Candidates[0].Content.Parts {
// 	switch v := part.(type) {
// 	case genai.Text:
// 		res.Text = string(v)
// 	case genai.FunctionCall:
// 		res.ToolCalls = append(res.ToolCalls, ToolCall{
// 			ID:        "not-sure-what-to-put-here",
// 			Function:  v.Name,
// 			Arguments: v.Args,
// 		})
// 	}
// }

// 	return res, nil
// }
