package agent

import "github.com/harkirath1511/docker-cli/internal/llm"

type StepResult struct {
	StepNumber int
	Thought    string
	ToolCalls  []llm.ToolCall
	FinalText  string
	IsDone     bool
	Error      error
}

type ToolExecutionResult struct {
	ToolCallID string
	Function   string
	Result     string
	Error      error
}

type BatchToolExecutionResult struct {
	Results []ToolExecutionResult
}
