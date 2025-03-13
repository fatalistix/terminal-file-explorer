package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/fatalistix/terminal-file-explorer/internal/tui"
)

func main() {
	app := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	if _, err := app.Run(); err != nil {
		log.Fatal("error running app: ", err)
	}
}
