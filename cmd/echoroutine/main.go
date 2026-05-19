package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kurojs/schedule-announcer/internal/tui"
)

func main() {
	m, err := tui.NewModel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing EchoRoutine: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running EchoRoutine: %v\n", err)
		os.Exit(1)
	}
}
