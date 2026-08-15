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

// Update is the Bubble Tea message handler — pure state machine
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	// ── Window resize ─────────────────────────────────────────────────────────
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.recalcLayout()

	// ── Keyboard ──────────────────────────────────────────────────────────────
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC:
			return m, tea.Quit

		// Enter submits; shift+enter inserts a newline (textarea handles it naturally)
		case tea.KeyEnter:
			text := strings.TrimSpace(m.input.Value())
			if text == "" || m.thinking {
				break
			}
			m.input.Reset()
			// After reset, recalculate layout so viewport grows back
			m = m.recalcLayout()

			// add user entry to display
			m.entries = append(m.entries, ChatEntry{Role: RoleUser, Content: text})
			m.syncViewport()

			// add to LLM history and start AI call
			m.llmHistory = append(m.llmHistory, llm.Message{Role: "user", Content: text})
			m.thinking = true
			cmds = append(cmds, callAI(m.ai, m.llmHistory, m.toolsDef))

		case tea.KeyCtrlL:
			// clear chat display (keeps llm history)
			m.entries = nil
			m.syncViewport()
		}

	// ── Spinner tick ──────────────────────────────────────────────────────────
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	// ── Docker ping result ────────────────────────────────────────────────────
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

	// ── AI response ───────────────────────────────────────────────────────────
	case AIResponseMsg:
		if msg.Text != "" {
			// Final text reply — show to user and stop loop
			m.entries = append(m.entries, ChatEntry{Role: RoleAI, Content: msg.Text})
			m.llmHistory = append(m.llmHistory, llm.Message{Role: "assistant", Content: msg.Text})
			m.thinking = false
			m.syncViewport()
			break
		}

		// Tool calls — execute each one
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

	// ── Tool result ───────────────────────────────────────────────────────────
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

		// After every tool result, call AI again to continue the loop
		cmds = append(cmds, callAI(m.ai, m.llmHistory, m.toolsDef))

	// ── AI error ──────────────────────────────────────────────────────────────
	case AIErrorMsg:
		m.entries = append(m.entries, ChatEntry{
			Role:    RoleError,
			Content: fmt.Sprintf("AI error: %v", msg.Err),
		})
		m.thinking = false
		m.syncViewport()
	}

	// Forward all events to sub-components
	var vpCmd, inputCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, vpCmd, inputCmd)

	// Recalculate layout whenever textarea may have grown
	m = m.recalcLayout()

	return m, tea.Batch(cmds...)
}

// recalcLayout recomputes the viewport height based on current textarea height.
// This makes the viewport shrink/grow as the textarea expands/contracts.
func (m Model) recalcLayout() Model {
	if m.width == 0 || m.height == 0 {
		return m
	}

	// textarea actual rendered height (number of lines it's showing)
	inputH := m.input.Height()
	if inputH < 1 {
		inputH = 1
	}

	// fixed chrome: 2 separators + 1 status + 1 thinking row + 1 newline padding
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

// syncViewport re-renders all entries into the viewport and jumps to bottom
func (m *Model) syncViewport() {
	m.viewport.SetContent(m.renderEntries())
	m.viewport.GotoBottom()
}

// formatArgs formats tool call args as key=value pairs
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

// ── Background Cmds ───────────────────────────────────────────────────────────

// callAI runs the LLM in a goroutine and returns the result as a tea.Msg
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

// executeTool runs a Docker tool in a background goroutine
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
