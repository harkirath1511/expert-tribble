package llm

type ToolParam struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"desc"`
}

type ToolDef struct {
	Name        string               `json:"name"`
	Description string               `json:"desc"`
	Params      map[string]ToolParam `json:"param"`
	Required    []string             `json:"required,omitempty"`
}

type LLMClient interface {
	GenerateResponse(history []Message, tools []ToolDef) (LLMRes, error)
}
