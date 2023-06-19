package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/couch"
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
	systems     []*core.System
	client      *couch.Couch[*core.System]
	currentView CurrentView
	list        *views.ListView
	form        *views.Form
}

func (t *TUI) VisitSystem(ctx context.Context, s *core.System) error {
	//TODO implement me
	panic("implement me")
}

func (t *TUI) VisitModule(ctx context.Context, m *core.Module) error {
	//TODO implement me
	panic("implement me")
}

func NewTUI(url string) tea.Model {
	return &TUI{
		currentView: ListView,
		client:      couch.NewCouch[*core.System](url),
	}
}

var SysInputs = []*views.FormInput{
	{
		Name: "Name",
	},
	{
		Name: "Description",
	},
	{
		Name: "Author",
	},
}

func (t *TUI) Init() tea.Cmd {
	var err error
	t.systems, err = t.client.List()
	if err != nil {
		panic(err)
	}
	items := make([]views.Item, len(t.systems))
	for i, s := range t.systems {
		items[i] = s
	}
	t.list = views.NewListView("Systems", items)
	t.form = views.NewForm[*core.System]("System", t.systems[0])
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
		t.form, cmd = t.form.Update(msg)
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
		return t.form.View()
	case ConfirmView:
	case CreateView:
	case UpdateView:
	}
	return ""
}
