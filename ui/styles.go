package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	colorPurple  = lipgloss.Color("#7C3AED")
	colorCyan    = lipgloss.Color("#06B6D4")
	colorYellow  = lipgloss.Color("#F59E0B")
	colorGreen   = lipgloss.Color("#10B981")
	colorRed     = lipgloss.Color("#EF4444")
	colorGray    = lipgloss.Color("#6B7280")
	colorWhite   = lipgloss.Color("#F9FAFB")
	colorBg      = lipgloss.Color("#0F0F1A")
	colorBorder  = lipgloss.Color("#2D2D4E")
	colorHeader  = lipgloss.Color("#1A1A2E")

	// Header bar
	HeaderStyle = lipgloss.NewStyle().
			Background(colorHeader).
			Foreground(colorCyan).
			Bold(true).
			Padding(0, 2)

	HeaderTitleStyle = lipgloss.NewStyle().
				Foreground(colorCyan).
				Bold(true)

	HeaderSubStyle = lipgloss.NewStyle().
			Foreground(colorGray)

	// Chat viewport
	ViewportStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	// Message styles
	UserLabelStyle = lipgloss.NewStyle().
			Foreground(colorPurple).
			Bold(true)

	UserMsgStyle = lipgloss.NewStyle().
			Foreground(colorWhite)

	AILabelStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true)

	AIMsgStyle = lipgloss.NewStyle().
			Foreground(colorWhite)

	// Tool call styles
	ToolLabelStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	ToolNameStyle = lipgloss.NewStyle().
			Foreground(colorYellow)

	ToolArgStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			Italic(true)

	ToolSuccessStyle = lipgloss.NewStyle().
				Foreground(colorGreen)

	ToolErrorStyle = lipgloss.NewStyle().
			Foreground(colorRed)

	// Spinner / thinking
	ThinkingStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Italic(true)

	// Input bar
	InputBarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPurple).
			Padding(0, 1)

	InputPromptStyle = lipgloss.NewStyle().
				Foreground(colorPurple).
				Bold(true)

	// Status bar
	StatusStyle = lipgloss.NewStyle().
			Background(colorHeader).
			Foreground(colorGray).
			Padding(0, 2)

	StatusReadyStyle = lipgloss.NewStyle().
				Foreground(colorGreen).
				Bold(true)

	StatusBusyStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	// Divider
	DimStyle = lipgloss.NewStyle().
			Foreground(colorGray)
)
