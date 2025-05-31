package directory

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatalistix/terminal-file-explorer/internal/model"
	"io"
)

var (
	itemStyle         = lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center)
	selectedItemStyle = lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Foreground(lipgloss.Color("170"))
)

type item struct {
	name string
}

func (i item) FilterValue() string {
	return i.name
}

func entryToListItem(entry model.DirectoryEntry) list.Item {
	return item{name: entry.Name}
}

func entriesToListItems(entries []model.DirectoryEntry) []list.Item {
	items := make([]list.Item, len(entries))
	for i, e := range entries {
		items[i] = entryToListItem(e)
	}

	return items
}

type itemDelegate struct{}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := i.name

	if index == m.Index() {
		_, _ = fmt.Fprint(w, selectedItemStyle.Render(str))
	} else {
		_, _ = fmt.Fprint(w, itemStyle.Render(str))
	}
}

func (_ itemDelegate) Height() int {
	return 1
}

func (_ itemDelegate) Spacing() int {
	return 0
}

func (_ itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

type directoryBase struct {
	list list.Model
}

func newBase(d model.Directory) *directoryBase {
	items := entriesToListItems(d.Entries)

	l := list.New(items, itemDelegate{}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowFilter(false)
	l.ResetSelected()

	return &directoryBase{
		list: l,
	}
}

func (b *directoryBase) View() string {
	return b.list.View()
}

func (b *directoryBase) SetHeight(height int) {
	b.list.SetHeight(height)
}

func (b *directoryBase) SetWidth(width int) {
	b.list.SetWidth(width)
}

type SelectableDirectory struct {
	*directoryBase
}

func NewSelectable(d model.Directory, selectedId int) *SelectableDirectory {
	base := newBase(d)
	base.list.Select(selectedId)

	return &SelectableDirectory{base}
}

type ViewDirectory struct {
	*directoryBase
}

func NewView(d model.Directory) *ViewDirectory {
	base := newBase(d)

	return &ViewDirectory{base}
}
