package agent

import (
	"fmt"

	"github.com/harkirath1511/docker-cli/internal/docker"
	"github.com/harkirath1511/docker-cli/internal/llm"
	"github.com/moby/moby/client"
)

const maxObservationLength = 4000

type Agent struct {
	dockerClient *client.Client
	llm          llm.LLMClient
	toolsDef     []llm.ToolDef

	History     []llm.Message
	CurrentStep int
	MaxSteps    int
	IsDone      bool
}

func NewAgent(dockerClient *client.Client, llmClient llm.LLMClient, tools []llm.ToolDef) *Agent {
	return &Agent{
		dockerClient: dockerClient,
		llm:          llmClient,
		toolsDef:     tools,
		MaxSteps:     10,
		History: []llm.Message{
			{Role: "system", Content: ReActSystemPrompt},
		},
	}
}

// StartNewGoal resets the loop state for a new user request
func (a *Agent) StartNewGoal(userGoal string) {
	a.CurrentStep = 0
	a.IsDone = false
	a.History = append(a.History, llm.Message{
		Role:    "user",
		Content: userGoal,
	})
}

// NextStep performs one cycle of Thought -> Action decision
func (a *Agent) NextStep() StepResult {
	if a.IsDone {
		return StepResult{IsDone: true}
	}
	a.CurrentStep++

	// 1. Safety Guard: Max steps check
	if a.CurrentStep > a.MaxSteps {
		a.IsDone = true
		return StepResult{
			StepNumber: a.CurrentStep,
			IsDone:     true,
			Error:      fmt.Errorf("agent exceeded safety limit of %d steps without finishing", a.MaxSteps),
		}
	}

	// 2. Query LLM with current history and tool definitions
	resp, err := a.llm.GenerateResponse(a.History, a.toolsDef)
	if err != nil {
		return StepResult{
			StepNumber: a.CurrentStep,
			Error:      fmt.Errorf("LLM error at step %d: %w", a.CurrentStep, err),
		}
	}

	// 3. Scenario A: No tools requested -> Agent has finished the goal
	if len(resp.ToolCalls) == 0 {
		a.IsDone = true
		a.History = append(a.History, llm.Message{
			Role:    "assistant",
			Content: resp.Text,
		})
		return StepResult{
			StepNumber: a.CurrentStep,
			FinalText:  resp.Text,
			IsDone:     true,
		}
	}

	// 4. Scenario B: Agent generated Thought + Tool Call(s) (Action)
	a.History = append(a.History, llm.Message{
		Role:      "assistant",
		Content:   resp.Text, // Contains the LLM's 'Thought'
		ToolCalls: resp.ToolCalls,
	})

	return StepResult{
		StepNumber: a.CurrentStep,
		Thought:    resp.Text,
		ToolCalls:  resp.ToolCalls,
		IsDone:     false,
	}
}

// ApplyObservation receives a single tool output and updates the agent's memory
func (a *Agent) ApplyObservation(execResult ToolExecutionResult) {
	content := execResult.Result
	if execResult.Error != nil {
		content = fmt.Sprintf("Action failed with error: %v. Please diagnose the cause and try an alternative step.", execResult.Error)
	}

	// Truncate overly long observations to prevent token context explosion
	if len(content) > maxObservationLength {
		content = content[:maxObservationLength] + "\n...[truncated output]..."
	}

	a.History = append(a.History, llm.Message{
		Role:    "tool",
		Content: content,
		ToolCalls: []llm.ToolCall{
			{ID: execResult.ToolCallID},
		},
	})
}

// ApplyDenial records that the user denied execution of sensitive operations
func (a *Agent) ApplyDenial(toolCalls []llm.ToolCall) {
	for _, tc := range toolCalls {
		a.History = append(a.History, llm.Message{
			Role:    "tool",
			Content: "Action was denied by the user. Do not retry this operation. Explain to the user or find a safe alternative.",
			ToolCalls: []llm.ToolCall{
				{ID: tc.ID},
			},
		})
	}
}

// ExecuteTools runs a batch of Docker tool calls synchronously and returns results
func (a *Agent) ExecuteTools(tcs []llm.ToolCall) BatchToolExecutionResult {
	var batch BatchToolExecutionResult
	for _, tc := range tcs {
		result, err := docker.Execute(a.dockerClient, tc.Function, tc.Arguments)
		batch.Results = append(batch.Results, ToolExecutionResult{
			ToolCallID: tc.ID,
			Function:   tc.Function,
			Result:     result,
			Error:      err,
		})
	}
	return batch
}
