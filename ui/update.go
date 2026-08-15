package ui

import (
	"fmt"

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

		headerH := 3
		inputH := 3
		statusH := 1
		vpHeight := m.height - headerH - inputH - statusH - 4 // paddings

		if !m.ready {
			m.viewport = viewport.New(m.width-4, vpHeight)
			m.viewport.SetContent(m.renderEntries())
			m.ready = true
		} else {
			m.viewport.Width = m.width - 4
			m.viewport.Height = vpHeight
		}
		m.input.Width = m.width - 8

	// ── Keyboard ──────────────────────────────────────────────────────────────
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			text := m.input.Value()
			if text == "" || m.thinking {
				break
			}
			m.input.Reset()

			// add user entry to display
			m.entries = append(m.entries, ChatEntry{Role: RoleUser, Content: text})
			m.syncViewport()

			// add to LLM history and start AI call
			m.llmHistory = append(m.llmHistory, llm.Message{Role: "user", Content: text})
			m.thinking = true
			m.statusMsg = "Thinking..."
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
			m.statusMsg = "Docker connected ✓"
		}
		m.syncViewport()

	// ── AI response ───────────────────────────────────────────────────────────
	case AIResponseMsg:
		if msg.Text != "" {
			// Final text reply — show to user and stop loop
			m.entries = append(m.entries, ChatEntry{Role: RoleAI, Content: msg.Text})
			m.llmHistory = append(m.llmHistory, llm.Message{Role: "assistant", Content: msg.Text})
			m.thinking = false
			m.statusMsg = "Ready"
			m.syncViewport()
			break
		}

		// Tool calls — execute each one sequentially by chaining cmds
		for _, tc := range msg.ToolCalls {
			// Show tool call entry in chat
			m.entries = append(m.entries, ChatEntry{
				Role:    RoleTool,
				Content: fmt.Sprintf("⚙  %s  %v", tc.Function, tc.Args),
			})

			// Save assistant message with tool call to LLM history
			m.llmHistory = append(m.llmHistory, llm.Message{
				Role:    "assistant",
				Content: "",
				ToolCalls: []llm.ToolCall{
					{ID: tc.ID, Function: tc.Function, Arguments: tc.Args},
				},
			})

			// Fire tool execution as a Cmd
			cmds = append(cmds, executeTool(m.dockerClient, tc))
		}
		m.syncViewport()

	// ── Tool result ───────────────────────────────────────────────────────────
	case ToolResultMsg:
		if msg.Err != nil {
			m.entries = append(m.entries, ChatEntry{
				Role:    RoleError,
				Content: fmt.Sprintf("  ❌ %s failed: %v", msg.Function, msg.Err),
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
				Content: fmt.Sprintf("  ✅ %s", msg.Result),
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
		m.statusMsg = "Error — try again"
		m.syncViewport()
	}

	// Forward keyboard events to sub-components
	var vpCmd, inputCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, vpCmd, inputCmd)

	return m, tea.Batch(cmds...)
}

// syncViewport re-renders all entries into the viewport and jumps to bottom
func (m *Model) syncViewport() {
	m.viewport.SetContent(m.renderEntries())
	m.viewport.GotoBottom()
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
