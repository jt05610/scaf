package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DetailView struct {
	Title    string
	Content  string
	Style    lipgloss.Style
	Commands []string
}

func (v *DetailView) View() string {
	//TODO implement me
	panic("implement me")
}

func NewDetailView(title string, content string) *DetailView {
	return &DetailView{
		Title:    title,
		Content:  content,
		Style:    lipgloss.NewStyle(),
		Commands: []string{"esc: go back"},
	}
}

func (v *DetailView) Init() tea.Cmd {
	return nil
}

func (v *DetailView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return v, tea.Quit
		}
	}
	return v, nil
}
