package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.ready {
		return "\n  Initialising...\n"
	}

	sep := SeparatorStyle.Render(strings.Repeat("─", m.width))

	var sections []string
	sections = append(sections, m.renderViewport())

	if m.thinking {
		sections = append(sections, m.renderThinking())
	}

	sections = append(sections, sep)
	sections = append(sections, m.renderInput())
	sections = append(sections, sep)
	sections = append(sections, m.renderStatus())

	return strings.Join(sections, "\n")
}

func (m Model) renderViewport() string {
	return ViewportStyle.Render(m.viewport.View())
}

func (m Model) renderThinking() string {
	return ThinkingStyle.Render(m.spinner.View() + " Thinking...")
}

func (m Model) renderInput() string {
	if m.awaitingApproval {
		prompt := ApprovalPromptStyle.Render("  Y  approve   ·   N  deny  ")
		return InputBarStyle.Render(prompt)
	}
	return InputBarStyle.Render(m.input.View())
}

func (m Model) renderStatus() string {
	hint := StatusTextStyle.Render("? for shortcuts · ctrl+c to quit · enter to send · shift+enter for newline")
	provider := StatusTextStyle.Render(m.statusMsg)

	hintW := lipgloss.Width(hint)
	provW := lipgloss.Width(provider)
	gap := m.width - hintW - provW - 4
	if gap < 0 {
		gap = 0
	}

	return StatusStyle.Width(m.width).Render(
		hint + strings.Repeat(" ", gap) + provider,
	)
}

func (m Model) renderEntries() string {
	if len(m.entries) == 0 {
		return strings.Join([]string{
			"",
			WelcomeTitleStyle.Render("Docker AI"),
			"",
			WelcomeHintStyle.Render(`Ask about your containers — e.g. "list my containers" or "stop nginx"`),
			"",
		}, "\n")
	}

	var sb strings.Builder
	for _, e := range m.entries {
		sb.WriteString(renderEntry(e))
	}
	return sb.String()
}

func renderEntry(e ChatEntry) string {
	switch e.Role {
	case RoleUser:
		prompt := UserPromptStyle.Render(">")
		msg := UserMsgStyle.Render(e.Content)
		return fmt.Sprintf("%s %s\n\n", prompt, msg)

	case RoleAI:
		msg := AIMsgStyle.Render(e.Content)
		lines := strings.Split(msg, "\n")
		var indented []string
		for _, l := range lines {
			indented = append(indented, "  "+l)
		}
		return strings.Join(indented, "\n") + "\n\n"

	case RoleTool:
		prefix := ToolPrefixStyle.Render("  ⎿ ")
		content := ToolResultStyle.Render(e.Content)
		return prefix + content + "\n"

	case RoleError:
		return ErrorStyle.Render("  ✗ "+e.Content) + "\n"

	case RoleApprovalRequest:
		// Yellow bordered warning card
		card := ApprovalCardStyle.Render(e.Content)
		return "\n" + card + "\n\n"

	case RoleApprovalGiven:
		return ApprovalGivenStyle.Render("  ✓ "+e.Content) + "\n\n"

	case RoleApprovalDenied:
		return ApprovalDeniedStyle.Render("  ✗ "+e.Content) + "\n\n"

	default:
		return e.Content + "\n"
	}
}

