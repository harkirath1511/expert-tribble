package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/harkirath1511/docker-cli/internal/agent"
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

	agent        *agent.Agent
	dockerClient *client.Client

	// flags
	thinking    bool
	ready       bool
	quitPending bool
	statusMsg   string

	// HITL state
	awaitingApproval bool
	pendingToolCalls []llm.ToolCall // entire batch held until user decides
}

func NewModel(dockerClient *client.Client, ai llm.LLMClient) Model {

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

	ag := agent.NewAgent(dockerClient, ai, tools.GetToolDefs())

	return Model{
		input:        ta,
		spinner:      sp,
		agent:        ag,
		dockerClient: dockerClient,
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
