package ui

// ChatRole identifies who sent a chat entry
type ChatRole int

const (
	RoleUser ChatRole = iota
	RoleAI
	RoleTool
	RoleError
)

// ChatEntry is one rendered line/block in the chat history
type ChatEntry struct {
	Role    ChatRole
	Content string // rendered text (may be multi-line)
}

// --- tea.Msg types sent between goroutines and the Update loop ---

// AIResponseMsg is sent when the LLM returns a result
type AIResponseMsg struct {
	Text      string
	ToolCalls []ToolCallMsg
}

// ToolCallMsg holds one tool call from the AI
type ToolCallMsg struct {
	ID       string
	Function string
	Args     map[string]any
}

// ToolResultMsg is sent after a docker tool executes
type ToolResultMsg struct {
	ToolCallID string
	Function   string
	Result     string
	Err        error
}

// AIErrorMsg is sent when the LLM returns an error
type AIErrorMsg struct {
	Err error
}
