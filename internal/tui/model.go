package tui

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/fatalistix/terminal-file-explorer/internal/model"
)

type DirectoryLoader interface {
	LoadDirectory(string) (model.Directory, error)
}

type Model struct {
	DirectoryLoader DirectoryLoader
	PreviousDir     model.Directory
	CurrentDir      model.Directory
	SelectedDir     model.Directory
}

func NewModel(DirectoryLoader DirectoryLoader) Model {
	return Model{
		DirectoryLoader: DirectoryLoader,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		const op = "tui.Model.Init"

		var err error

		m.CurrentDir, err = m.DirectoryLoader.LoadDirectory(".")
		if err != nil {
			return fmt.Errorf("%s: failed to load directory: %w", op, err)
		}

		return nil
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			return m, m.handleKeyMsg(msg)
		}
	}
	return m, nil
}

func (m Model) View() string {
    s :=
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit
	}

	return nil
}
