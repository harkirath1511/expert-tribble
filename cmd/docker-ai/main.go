package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/harkirath1511/docker-cli/internal/docker"
	"github.com/harkirath1511/docker-cli/internal/llm"
	"github.com/harkirath1511/docker-cli/ui"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Init Docker
	dockerClient, err := docker.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n  🐳 \033[1;31mDocker Error\033[0m\n\n  %v\n\n  Start Docker Desktop and try again.\n\n", err)
		os.Exit(1)
	}
	defer dockerClient.Close()

	// Init AI
	ai, err := llm.NewGroqClient()
	if err != nil {
		log.Fatal("AI init error: ", err)
	}

	// Build TUI model
	m := ui.NewModel(dockerClient, ai)

	// Run with full-screen alt-screen (like Claude Code / agy)
	p := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
