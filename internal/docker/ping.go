package docker

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/moby/moby/client"
)

// PingResultMsg is sent back to the TUI after a Docker ping attempt
type PingResultMsg struct {
	Err error
}

// PingCmd pings the Docker daemon and returns the result as a tea.Msg
func PingCmd(apiClient *client.Client) tea.Cmd {
	return func() tea.Msg {
		_, err := apiClient.Ping(context.Background(), client.PingOptions{
			NegotiateAPIVersion: true,
		})
		return PingResultMsg{Err: err}
	}
}
