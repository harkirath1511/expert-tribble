package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/harkirath1511/docker-cli/internal/docker"
	"github.com/harkirath1511/docker-cli/internal/llm"
	"github.com/moby/moby/client"
)


func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.recalcLayout()

	
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC:
			return m, tea.Quit

		
		case tea.KeyEnter:
			text := strings.TrimSpace(m.input.Value())
			if text == "" || m.thinking {
				break
			}
			m.input.Reset()
			
			m = m.recalcLayout()

			
			m.entries = append(m.entries, ChatEntry{Role: RoleUser, Content: text})
			m.syncViewport()

			
			m.llmHistory = append(m.llmHistory, llm.Message{Role: "user", Content: text})
			m.thinking = true
			cmds = append(cmds, callAI(m.ai, m.llmHistory, m.toolsDef))

		case tea.KeyCtrlL:
			
			m.entries = nil
			m.syncViewport()
		}

	
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	
	case docker.PingResultMsg:
		if msg.Err != nil {
			m.entries = append(m.entries, ChatEntry{
				Role:    RoleError,
				Content: fmt.Sprintf("Docker unreachable: %v", msg.Err),
			})
			m.statusMsg = "Docker offline"
		} else {
			m.statusMsg = "Groq · llama-3.3-70b"
		}
		m.syncViewport()

	
	case AIResponseMsg:
		if msg.Text != "" {
			
			m.entries = append(m.entries, ChatEntry{Role: RoleAI, Content: msg.Text})
			m.llmHistory = append(m.llmHistory, llm.Message{Role: "assistant", Content: msg.Text})
			m.thinking = false
			m.syncViewport()
			break
		}

		
		for _, tc := range msg.ToolCalls {
			m.entries = append(m.entries, ChatEntry{
				Role:    RoleTool,
				Content: fmt.Sprintf("%s(%v)", tc.Function, formatArgs(tc.Args)),
			})

			m.llmHistory = append(m.llmHistory, llm.Message{
				Role:    "assistant",
				Content: "",
				ToolCalls: []llm.ToolCall{
					{ID: tc.ID, Function: tc.Function, Arguments: tc.Args},
				},
			})

			cmds = append(cmds, executeTool(m.dockerClient, tc))
		}
		m.syncViewport()

	
	case ToolResultMsg:
		if msg.Err != nil {
			m.entries = append(m.entries, ChatEntry{
				Role:    RoleError,
				Content: fmt.Sprintf("%s failed: %v", msg.Function, msg.Err),
			})
			m.llmHistory = append(m.llmHistory, llm.Message{
				Role:    "tool",
				Content: fmt.Sprintf("error: %v", msg.Err),
				ToolCalls: []llm.ToolCall{
					{ID: msg.ToolCallID},
				},
			})
		} else {
			m.entries = append(m.entries, ChatEntry{
				Role:    RoleTool,
				Content: msg.Result,
			})
			m.llmHistory = append(m.llmHistory, llm.Message{
				Role:    "tool",
				Content: msg.Result,
				ToolCalls: []llm.ToolCall{
					{ID: msg.ToolCallID},
				},
			})
		}
		m.syncViewport()

		
		cmds = append(cmds, callAI(m.ai, m.llmHistory, m.toolsDef))

	
	case AIErrorMsg:
		m.entries = append(m.entries, ChatEntry{
			Role:    RoleError,
			Content: fmt.Sprintf("AI error: %v", msg.Err),
		})
		m.thinking = false
		m.syncViewport()
	}

	
	var vpCmd, inputCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, vpCmd, inputCmd)

	
	m = m.recalcLayout()

	return m, tea.Batch(cmds...)
}



func (m Model) recalcLayout() Model {
	if m.width == 0 || m.height == 0 {
		return m
	}

	
	inputH := m.input.Height()
	if inputH < 1 {
		inputH = 1
	}

	
	fixedH := 2 + 1 + 1 + 1
	vpHeight := m.height - inputH - fixedH
	if vpHeight < 1 {
		vpHeight = 1
	}

	if !m.ready {
		m.viewport = viewport.New(m.width, vpHeight)
		m.viewport.SetContent(m.renderEntries())
		m.ready = true
	} else {
		m.viewport.Width = m.width
		m.viewport.Height = vpHeight
	}

	m.input.SetWidth(m.width - 2)
	return m
}


func (m *Model) syncViewport() {
	m.viewport.SetContent(m.renderEntries())
	m.viewport.GotoBottom()
}


func formatArgs(args map[string]any) string {
	if len(args) == 0 {
		return ""
	}
	var parts []string
	for k, v := range args {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(parts, ", ")
}




func callAI(ai *llm.GroqClient, history []llm.Message, toolsDef []llm.ToolDef) tea.Cmd {
	return func() tea.Msg {
		resp, err := ai.GenerateResponse(history, toolsDef)
		if err != nil {
			return AIErrorMsg{Err: err}
		}

		msg := AIResponseMsg{Text: resp.Text}
		for _, tc := range resp.ToolCalls {
			msg.ToolCalls = append(msg.ToolCalls, ToolCallMsg{
				ID:       tc.ID,
				Function: tc.Function,
				Args:     tc.Arguments,
			})
		}
		return msg
	}
}


func executeTool(dockerClient *client.Client, tc ToolCallMsg) tea.Cmd {
	return func() tea.Msg {
		result, err := docker.Execute(dockerClient, tc.Function, tc.Args)
		return ToolResultMsg{
			ToolCallID: tc.ID,
			Function:   tc.Function,
			Result:     result,
			Err:        err,
		}
	}
}
