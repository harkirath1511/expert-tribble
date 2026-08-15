package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the full terminal UI every frame
func (m Model) View() string {
	if !m.ready {
		return "\n  Initialising...\n"
	}

	return strings.Join([]string{
		m.renderHeader(),
		m.renderViewport(),
		m.renderThinking(),
		m.renderInput(),
		m.renderStatus(),
	}, "\n")
}

// ── Sections ──────────────────────────────────────────────────────────────────

func (m Model) renderHeader() string {
	title := HeaderTitleStyle.Render("🐳 Docker AI")
	hint := HeaderSubStyle.Render("  ctrl+c to quit · ctrl+l to clear")
	content := lipgloss.JoinHorizontal(lipgloss.Top, title, hint)
	return HeaderStyle.Width(m.width).Render(content)
}

func (m Model) renderViewport() string {
	return ViewportStyle.
		Width(m.width - 2).
		Height(m.viewport.Height).
		Render(m.viewport.View())
}

func (m Model) renderThinking() string {
	if !m.thinking {
		return ""
	}
	return "  " + m.spinner.View() + ThinkingStyle.Render(" Thinking...")
}

func (m Model) renderInput() string {
	return InputBarStyle.
		Width(m.width - 2).
		Render(m.input.View())
}

func (m Model) renderStatus() string {
	var status string
	if m.thinking {
		status = StatusBusyStyle.Render("● busy")
	} else {
		status = StatusReadyStyle.Render("● " + m.statusMsg)
	}

	scrollPct := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
	scroll := DimStyle.Render(scrollPct)

	gap := m.width - lipgloss.Width(status) - lipgloss.Width(scroll) - 4
	if gap < 0 {
		gap = 0
	}

	return StatusStyle.Width(m.width).Render(
		status + strings.Repeat(" ", gap) + scroll,
	)
}

// ── Entry rendering ───────────────────────────────────────────────────────────

// renderEntries converts all ChatEntry items into a single string for the viewport
func (m Model) renderEntries() string {
	if len(m.entries) == 0 {
		return DimStyle.Render("\n  Start by asking something — e.g. \"list my containers\" or \"start nginx\"\n")
	}

	var sb strings.Builder
	for _, e := range m.entries {
		sb.WriteString(renderEntry(e))
		sb.WriteString("\n")
	}
	return sb.String()
}

func renderEntry(e ChatEntry) string {
	switch e.Role {
	case RoleUser:
		label := UserLabelStyle.Render("  You")
		msg := UserMsgStyle.Render(e.Content)
		return fmt.Sprintf("%s\n  %s\n", label, msg)

	case RoleAI:
		label := AILabelStyle.Render("  🤖 Docker AI")
		msg := AIMsgStyle.Render(e.Content)
		return fmt.Sprintf("%s\n  %s\n", label, msg)

	case RoleTool:
		return ToolNameStyle.Render(e.Content) + "\n"

	case RoleError:
		return ToolErrorStyle.Render("  " + e.Content) + "\n"

	default:
		return e.Content + "\n"
	}
}
