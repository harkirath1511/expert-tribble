package ui

import "github.com/charmbracelet/lipgloss"

var (
	
	colorAccent  = lipgloss.Color("#60A5FA") 
	colorText    = lipgloss.Color("#D1D5DB") 
	colorDim     = lipgloss.Color("#4B5563") 
	colorVeryDim = lipgloss.Color("#374151") 

	
	ViewportStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	
	UserPromptStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	UserMsgStyle = lipgloss.NewStyle().
			Foreground(colorText)

	
	AIMsgStyle = lipgloss.NewStyle().
			Foreground(colorText)

	
	ToolPrefixStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	ToolResultStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF"))

	
	ThinkingStyle = lipgloss.NewStyle().
			Foreground(colorDim).
			PaddingLeft(2)

	
	InputBarStyle = lipgloss.NewStyle().
			PaddingLeft(1)

	InputPromptStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true)

	
	StatusStyle = lipgloss.NewStyle().
			Foreground(colorVeryDim).
			PaddingLeft(2).
			PaddingRight(2)

	StatusTextStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	
	WelcomeTitleStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true).
				PaddingLeft(2)

	WelcomeHintStyle = lipgloss.NewStyle().
				Foreground(colorDim).
				PaddingLeft(2)

	
	SeparatorStyle = lipgloss.NewStyle().
			Foreground(colorVeryDim)
)
