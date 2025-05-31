package main

import (
	"github.com/fatalistix/terminal-file-explorer/internal/service"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/fatalistix/terminal-file-explorer/internal/tui"
)

func main() {
	s, err := service.NewStateHolder(".")
	m, err := tui.NewRootModel(s)
	if err != nil {
		log.Fatal("error running app: ", err)
	}

	app := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := app.Run(); err != nil {
		log.Fatal("error running app: ", err)
	}
}
