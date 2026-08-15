package ui

import "github.com/charmbracelet/lipgloss"

var (
	// ── 3-color palette ───────────────────────────────────────────────────────
	colorAccent  = lipgloss.Color("#60A5FA") // subtle blue — prompt & app name only
	colorText    = lipgloss.Color("#D1D5DB") // primary text — user & AI messages
	colorDim     = lipgloss.Color("#4B5563") // muted — hints, tool output, status
	colorVeryDim = lipgloss.Color("#374151") // separator lines

	// ── Viewport (no border, just padding) ───────────────────────────────────
	ViewportStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	// ── User message ─────────────────────────────────────────────────────────
	UserPromptStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	UserMsgStyle = lipgloss.NewStyle().
			Foreground(colorText)

	// ── AI message ───────────────────────────────────────────────────────────
	AIMsgStyle = lipgloss.NewStyle().
			Foreground(colorText)

	// ── Tool call ────────────────────────────────────────────────────────────
	ToolPrefixStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	ToolResultStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	// ── Error ────────────────────────────────────────────────────────────────
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF"))

	// ── Thinking indicator ───────────────────────────────────────────────────
	ThinkingStyle = lipgloss.NewStyle().
			Foreground(colorDim).
			PaddingLeft(2)

	// ── Input bar — textarea sits between two separator lines ────────────────
	InputBarStyle = lipgloss.NewStyle().
			PaddingLeft(1)

	InputPromptStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true)

	// ── Status bar ───────────────────────────────────────────────────────────
	StatusStyle = lipgloss.NewStyle().
			Foreground(colorVeryDim).
			PaddingLeft(2).
			PaddingRight(2)

	StatusTextStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	// ── Welcome / empty state ────────────────────────────────────────────────
	WelcomeTitleStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true).
				PaddingLeft(2)

	WelcomeHintStyle = lipgloss.NewStyle().
				Foreground(colorDim).
				PaddingLeft(2)

	// ── Separator ────────────────────────────────────────────────────────────
	SeparatorStyle = lipgloss.NewStyle().
			Foreground(colorVeryDim)
)
