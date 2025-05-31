package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatalistix/terminal-file-explorer/internal/model"
	"github.com/fatalistix/terminal-file-explorer/internal/service"
	"github.com/fatalistix/terminal-file-explorer/internal/tui/directory"
)

type DirectoryLoader interface {
	LoadDirectoryContent(string) (model.Directory, error)
	LoadDirectoryName() (string, error)
}

type RootModel struct {
	width, height int
	stateHolder   *service.StateHolder
}

func NewRootModel(stateHolder *service.StateHolder) (*RootModel, error) {
	const op = "tui.NewModel"

	return &RootModel{
		width:       0,
		height:      0,
		stateHolder: stateHolder,
	}, nil
}

func (m *RootModel) Init() tea.Cmd {
	return nil
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		{
			return m.handleWindowSizeMsg(msg)
		}
	case tea.KeyMsg:
		{
			return m.handleKeyMsg(msg)
		}
	}
	return m, nil
}

func (m *RootModel) handleWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	return m, nil
}

func (m *RootModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		m.stateHolder.MoveUp()
		return m, nil
	case "down", "j":
		m.stateHolder.MoveDown()
		return m, nil
	case "enter", " ", "l":
		m.stateHolder.Select()
		return m, nil
	case "backspace", "h":
		m.stateHolder.GoToPreviousDirectory()
		return m, nil
	default:
		return m, nil
	}
}

func (m *RootModel) View() string {
	state := m.stateHolder.GetState()

	averageWidth := m.width / 3
	remainingWidth := m.width - averageWidth*3

	//m.prevDirView.SetHeight(m.height)
	//m.prevDirView.SetWidth(averageWidth)
	//
	//m.currDirView.SetHeight(m.height)
	//m.currDirView.SetWidth(averageWidth + remainingWidth)

	prevDirBlock := lipgloss.NewStyle().
		Height(m.height).
		Width(averageWidth)

	currDirBlock := lipgloss.NewStyle().
		Height(m.height).
		Width(averageWidth + remainingWidth)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		prevDirBlock.Render(renderView(state.PrevViewContent, averageWidth, m.height)),
		currDirBlock.Render(renderView(state.CurrViewContent, averageWidth+remainingWidth, m.height)),
	)
}

func renderView(content model.Content, width, height int) string {
	switch content := content.(type) {
	case *model.EmptyContent:
		return ""
	case *model.DirectoryContent:
		if content.SelectedId < 0 {
			view := directory.NewView(content.Directory)
			view.SetWidth(width)
			view.SetHeight(height)
			return view.View()
		} else {
			selectable := directory.NewSelectable(content.Directory, content.SelectedId)
			selectable.SetWidth(width)
			selectable.SetHeight(height)
			return selectable.View()
		}
	case *model.TextContent:
		return content.Text
	}

	panic("unreachable")
}
