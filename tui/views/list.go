package views

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

type ListView struct {
	list           list.Model
	keys           *listKeyMap
	delegateKeys   *delegateKeyMap
	showTitle      bool
	showFilter     bool
	showStatus     bool
	showPagination bool
	showHelp       bool
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

func NewListView(title string, items []Item) *ListView {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	l := list.New(listItems, newItemDelegate(delegateKeys), 0, 0)
	l.Title = title
	l.Styles.Title = TitleStyle
	l.SetShowStatusBar(true)
	l.SetShowPagination(true)
	l.SetShowTitle(true)
	l.SetShowHelp(true)
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}
	return &ListView{
		list:           l,
		keys:           listKeys,
		delegateKeys:   delegateKeys,
		showTitle:      true,
		showFilter:     true,
		showStatus:     true,
		showPagination: true,
		showHelp:       true,
	}
}

func (v *ListView) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (v *ListView) Update(msg tea.Msg) (*ListView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, vert := AppStyle.GetFrameSize()
		v.list.SetSize(msg.Width-h, msg.Height-vert)
	case tea.KeyMsg:
		if v.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, v.keys.toggleSpinner):
			cmd := v.list.ToggleSpinner()
			return v, cmd
		case key.Matches(msg, v.keys.toggleTitleBar):
			val := !v.list.ShowTitle()
			v.list.SetShowTitle(val)
			v.list.SetShowFilter(val)
			v.list.SetFilteringEnabled(val)
			return v, nil
		case key.Matches(msg, v.keys.toggleStatusBar):
			v.list.SetShowStatusBar(!v.list.ShowStatusBar())
			return v, nil
		case key.Matches(msg, v.keys.togglePagination):
			v.list.SetShowPagination(!v.list.ShowPagination())
			return v, nil
		case key.Matches(msg, v.keys.toggleHelpMenu):
			v.list.SetShowHelp(!v.list.ShowHelp())
			return v, nil
		case key.Matches(msg, v.keys.insertItem):
			return v, tea.Batch(v.list.NewStatusMessage(StatusMessageStyle("Loading create")), func() tea.Msg {
				return ShowCreateMsg{}
			})
		}
	}

	newListModel, cmd := v.list.Update(msg)
	v.list = newListModel
	cmds = append(cmds, cmd)
	return v, tea.Batch(cmds...)
}

func (v *ListView) View() string {
	return v.list.View()
}
