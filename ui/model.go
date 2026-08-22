package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/harkirath1511/docker-cli/internal/docker"
	"github.com/harkirath1511/docker-cli/internal/llm"
	"github.com/harkirath1511/docker-cli/internal/tools"
	"github.com/moby/moby/client"
)


type Model struct {
	
	width  int
	height int

	
	viewport viewport.Model
	input    textarea.Model
	spinner  spinner.Model

	
	entries []ChatEntry

	
	llmHistory []llm.Message
	toolsDef   []llm.ToolDef

	
	ai           *llm.GroqClient
	dockerClient *client.Client

	
	thinking  bool
	ready     bool   
	statusMsg string 
}


func NewModel(dockerClient *client.Client, ai *llm.GroqClient) Model {
	
	ta := textarea.New()
	ta.Placeholder = "Ask anything..."
	ta.Focus()
	ta.CharLimit = 2000
	ta.ShowLineNumbers = false
	ta.EndOfBufferCharacter = 0 
	ta.SetHeight(1)             
	ta.MaxHeight = 6            

	
	noStyle := lipgloss.NewStyle()
	ta.FocusedStyle.Base = noStyle
	ta.FocusedStyle.CursorLine = noStyle
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(colorDim)
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(colorText)
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(colorAccent)
	ta.BlurredStyle.Base = noStyle
	ta.BlurredStyle.CursorLine = noStyle
	ta.BlurredStyle.Placeholder = lipgloss.NewStyle().Foreground(colorDim)
	ta.BlurredStyle.Text = lipgloss.NewStyle().Foreground(colorDim)
	ta.BlurredStyle.Prompt = lipgloss.NewStyle().Foreground(colorDim)

	ta.Prompt = "> "

	
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(colorDim)

	
	history := []llm.Message{
		{
			Role:    "system",
			Content: "You are a Docker assistant that helps users manage Docker containers and images. Use the provided tools to interact with Docker. Always use tool calls (function calls) to perform Docker operations - never write out function calls as text or XML. When a task is complete, provide a brief summary to the user.",
		},
	}

	return Model{
		input:        ta,
		spinner:      sp,
		ai:           ai,
		dockerClient: dockerClient,
		llmHistory:   history,
		toolsDef:     tools.GetToolDefs(),
		statusMsg:    "Groq · llama-3.3-70b",
	}
}


func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		m.spinner.Tick,
		docker.PingCmd(m.dockerClient),
	)
}
