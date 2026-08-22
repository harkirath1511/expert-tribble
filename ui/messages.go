package ui


type ChatRole int

const (
	RoleUser ChatRole = iota
	RoleAI
	RoleTool
	RoleError
	RoleApprovalRequest // HITL: warning card shown before execution
	RoleApprovalGiven   // HITL: user pressed Y
	RoleApprovalDenied  // HITL: user pressed N
)


type ChatEntry struct {
	Role    ChatRole
	Content string 
}




type AIResponseMsg struct {
	Text      string
	ToolCalls []ToolCallMsg
}


type ToolCallMsg struct {
	ID       string
	Function string
	Args     map[string]any
}


type ToolResultMsg struct {
	ToolCallID string
	Function   string
	Result     string
	Err        error
}


type AIErrorMsg struct {
	Err error
}
