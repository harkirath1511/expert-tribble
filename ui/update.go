package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/harkirath1511/docker-cli/internal/agent"
	"github.com/harkirath1511/docker-cli/internal/docker"
	"github.com/harkirath1511/docker-cli/internal/llm"
	"github.com/harkirath1511/docker-cli/internal/tools"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.recalcLayout()

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			if m.quitPending {
				return m, tea.Quit
			}
			m.quitPending = true
			m.syncViewport()
			return m, nil
		}

		if m.quitPending {
			m.quitPending = false
		}

		// ── HITL approval mode — swallow all keys except Y and N ────────────
		if m.awaitingApproval {
			switch msg.String() {

			case "y", "Y":
				// User approved the entire batch — execute all pending tool calls
				m.awaitingApproval = false
				m.thinking = true
				m.entries = append(m.entries, ChatEntry{
					Role:    RoleApprovalGiven,
					Content: "Approved — executing...",
				})
				for _, tc := range m.pendingToolCalls {
					m.entries = append(m.entries, ChatEntry{
						Role:    RoleTool,
						Content: fmt.Sprintf("%s(%v)", tc.Function, formatArgs(tc.Arguments)),
					})
				}
				m.syncViewport()
				cmds = append(cmds, executeToolsCmd(m.agent, m.pendingToolCalls))
				m.pendingToolCalls = nil

			case "n", "N", "esc":
				// User denied — record denial into agent memory and continue loop
				m.awaitingApproval = false
				m.thinking = true
				m.entries = append(m.entries, ChatEntry{
					Role:    RoleApprovalDenied,
					Content: "Denied — telling AI to stop.",
				})
				m.syncViewport()
				m.agent.ApplyDenial(m.pendingToolCalls)
				m.pendingToolCalls = nil
				cmds = append(cmds, agentStepCmd(m.agent))
			}
			// swallow all other keys while awaiting approval
			return m, tea.Batch(cmds...)
		}

		// ── Normal input mode ────────────────────────────────────────────────
		switch msg.Type {

		case tea.KeyEnter:
			text := strings.TrimSpace(m.input.Value())
			if text == "" || m.thinking {
				break
			}
			m.input.Reset()
			m = m.recalcLayout()

			m.entries = append(m.entries, ChatEntry{Role: RoleUser, Content: text})
			m.syncViewport()

			m.thinking = true
			m.agent.StartNewGoal(text)
			cmds = append(cmds, agentStepCmd(m.agent))

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

	case agent.StepResult:
		if msg.Error != nil {
			m.entries = append(m.entries, ChatEntry{
				Role:    RoleError,
				Content: fmt.Sprintf("Error: %v", msg.Error),
			})
			m.thinking = false
			m.syncViewport()
			break
		}

		// Display agent reasoning thought if present
		if msg.Thought != "" {
			m.entries = append(m.entries, ChatEntry{
				Role:    RoleThought,
				Content: fmt.Sprintf("[Step %d] %s", msg.StepNumber, msg.Thought),
			})
			m.syncViewport()
		}

		// Goal accomplished
		if msg.IsDone {
			if msg.FinalText != "" {
				m.entries = append(m.entries, ChatEntry{
					Role:    RoleAI,
					Content: msg.FinalText,
				})
			}
			m.thinking = false
			m.syncViewport()
			break
		}

		// Process tool calls requested by the agent
		if len(msg.ToolCalls) > 0 {
			// Scan for sensitive tools requiring HITL approval
			var sensitiveLines []string
			for _, tc := range msg.ToolCalls {
				if reason, yes := tools.IsSensitive(tc.Function); yes {
					sensitiveLines = append(sensitiveLines,
						fmt.Sprintf("  ⚠  %-22s  %s", tc.Function, reason),
					)
				}
			}

			if len(sensitiveLines) > 0 {
				m.awaitingApproval = true
				m.thinking = false
				m.pendingToolCalls = msg.ToolCalls

				var cardLines []string
				cardLines = append(cardLines, "AI wants to run the following operations:")
				cardLines = append(cardLines, "")
				for _, tc := range msg.ToolCalls {
					if reason, yes := tools.IsSensitive(tc.Function); yes {
						cardLines = append(cardLines,
							fmt.Sprintf("  ⚠  %-22s  %s", tc.Function, reason))
					} else {
						cardLines = append(cardLines,
							fmt.Sprintf("  ·  %-22s  (safe)", tc.Function))
					}
				}
				cardLines = append(cardLines, "")
				cardLines = append(cardLines, "Press  Y  to approve all   ·   N  to deny all")

				m.entries = append(m.entries, ChatEntry{
					Role:    RoleApprovalRequest,
					Content: strings.Join(cardLines, "\n"),
				})
				m.syncViewport()
				break
			}

			// Safe tools — execute immediately
			for _, tc := range msg.ToolCalls {
				m.entries = append(m.entries, ChatEntry{
					Role:    RoleTool,
					Content: fmt.Sprintf("%s(%v)", tc.Function, formatArgs(tc.Arguments)),
				})
			}
			m.syncViewport()
			cmds = append(cmds, executeToolsCmd(m.agent, msg.ToolCalls))
		}

	case agent.BatchToolExecutionResult:
		for _, res := range msg.Results {
			if res.Error != nil {
				m.entries = append(m.entries, ChatEntry{
					Role:    RoleError,
					Content: fmt.Sprintf("%s failed: %v", res.Function, res.Error),
				})
			} else {
				m.entries = append(m.entries, ChatEntry{
					Role:    RoleTool,
					Content: res.Result,
				})
			}
			// Feed observation back into the agent
			m.agent.ApplyObservation(res)
		}
		m.syncViewport()

		// Continue the ReAct loop
		cmds = append(cmds, agentStepCmd(m.agent))
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
		m.viewport = viewport.New(maxContentWidth(m.width), vpHeight)
		m.ready = true
	} else {
		m.viewport.Width = maxContentWidth(m.width)
		m.viewport.Height = vpHeight
	}

	m.input.SetWidth(m.width - 2)
	m.viewport.SetContent(lipgloss.NewStyle().
		Width(m.viewport.Width).
		Render(m.renderEntries()))
	return m
}

func maxContentWidth(width int) int {
	contentWidth := width - 2
	if contentWidth < 1 {
		return 1
	}
	return contentWidth
}

func (m *Model) syncViewport() {
	contentWidth := m.viewport.Width
	if contentWidth < 1 {
		contentWidth = 1
	}
	content := lipgloss.NewStyle().Width(contentWidth).Render(m.renderEntries())
	m.viewport.SetContent(content)
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

func agentStepCmd(ag *agent.Agent) tea.Cmd {
	return func() tea.Msg {
		return ag.NextStep()
	}
}

func executeToolsCmd(ag *agent.Agent, tcs []llm.ToolCall) tea.Cmd {
	return func() tea.Msg {
		return ag.ExecuteTools(tcs)
	}
}
