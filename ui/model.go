package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/harkirath1511/docker-cli/internal/docker"
	"github.com/harkirath1511/docker-cli/internal/llm"
	"github.com/harkirath1511/docker-cli/internal/tools"
	"github.com/moby/moby/client"
)

// Model is the root Bubble Tea model
type Model struct {
	// terminal dimensions
	width  int
	height int

	// sub-components
	viewport viewport.Model
	input    textinput.Model
	spinner  spinner.Model

	// chat display history
	entries []ChatEntry

	// AI conversation state
	llmHistory []llm.Message
	toolsDef   []llm.ToolDef

	// clients
	ai           *llm.GroqClient
	dockerClient *client.Client

	// flags
	thinking  bool
	ready     bool   // viewport initialised?
	statusMsg string // shown in the status bar
}

// NewModel creates and wires up the initial model
func NewModel(dockerClient *client.Client, ai *llm.GroqClient) Model {
	// text input
	ti := textinput.New()
	ti.Placeholder = "Ask me anything about your Docker setup..."
	ti.Focus()
	ti.CharLimit = 500
	ti.PromptStyle = InputPromptStyle
	ti.Prompt = " ❯ "

	// spinner
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#06B6D4"))

	// system message seeds the conversation
	history := []llm.Message{
		{
			Role:    "system",
			Content: "You are a Docker assistant that helps users manage Docker containers and images. Use the provided tools to interact with Docker. Always use tool calls (function calls) to perform Docker operations - never write out function calls as text or XML. When a task is complete, provide a brief summary to the user.",
		},
	}

	return Model{
		input:        ti,
		spinner:      sp,
		ai:           ai,
		dockerClient: dockerClient,
		llmHistory:   history,
		toolsDef:     tools.GetToolDefs(),
		statusMsg:    "Ready",
	}
}

// Init is called once on startup — start the spinner tick
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.spinner.Tick,
		docker.PingCmd(m.dockerClient), // confirm docker is alive
	)
}
