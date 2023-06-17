package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/tui/views"
)

type CurrentView int

const (
	ListView CurrentView = iota
	DetailView
	ConfirmView
	CreateView
	UpdateView
)

type TUI struct {
	currentView CurrentView
	list        *views.ListView
	detail      *views.Form
}

func (t *TUI) VisitSystem(ctx context.Context, s *core.System) error {
	//TODO implement me
	panic("implement me")
}

func (t *TUI) VisitModule(ctx context.Context, m *core.Module) error {
	//TODO implement me
	panic("implement me")
}

func NewTUI(systems []*core.System) tea.Model {
	items := make([]views.Item, len(systems))
	for i, system := range systems {
		items[i] = system
	}

	return &TUI{
		currentView: ListView,
		list:        views.NewListView("Systems", items),
		detail: views.NewForm("System", []*views.FormInput{
			{
				Name:     "Name",
				Prompt:   "",
				Validate: nil,
			},
			{
				Name:     "Description",
				Prompt:   "",
				Validate: nil,
			},
		}),
	}
}

func (t *TUI) Init() tea.Cmd {
	return t.list.Init()
}

func (t *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.(type) {
	case views.ShowListMsg:
		t.currentView = ListView
	case views.ShowDetailMsg:
		t.currentView = DetailView
	case views.ShowConfirmMsg:
		t.currentView = ConfirmView
	case views.ShowCreateMsg:
		t.currentView = CreateView
	case views.ShowEditMsg:
		t.currentView = UpdateView
	}

	switch t.currentView {
	case ListView:
		t.list, cmd = t.list.Update(msg)
	case DetailView:
		t.detail, cmd = t.detail.Update(msg)
	case ConfirmView:
	case CreateView:
	case UpdateView:

	}
	return t, cmd
}

func (t *TUI) View() string {
	switch t.currentView {
	case ListView:
		return t.list.View()
	case DetailView:
		return t.detail.View()
	case ConfirmView:
	case CreateView:
	case UpdateView:
	}
	return ""
}
