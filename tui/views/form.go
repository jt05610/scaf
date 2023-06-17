package views

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type Validator func(string) error

type FormInput struct {
	Name     string
	Prompt   string
	Validate Validator
}

type errMsg error

type Form struct {
	title   string
	inputs  []textinput.Model
	names   []string
	focused int
	err     error
}

func (f *Form) Init() tea.Cmd {
	return textinput.Blink
}

func (f *Form) nextInput() {
	f.focused = (f.focused + 1) % len(f.inputs)
}

func (f *Form) prevInput() {
	f.focused--
	if f.focused < 0 {
		f.focused = len(f.inputs) - 1
	}
}

func (f *Form) Update(msg tea.Msg) (*Form, tea.Cmd) {
	cmds := make([]tea.Cmd, len(f.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if f.focused == len(f.inputs)-1 {
				return f, func() tea.Msg {
					return ShowConfirmMsg{}
				}
			}
			f.nextInput()
		case tea.KeyCtrlC:
			return f, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			f.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			f.nextInput()
		}
		for i := range f.inputs {
			f.inputs[i].Blur()
		}
		f.inputs[f.focused].Focus()
	case errMsg:
		f.err = msg
		return f, nil
	}
	for i := range f.inputs {
		f.inputs[i], cmds[i] = f.inputs[i].Update(msg)
	}
	return f, tea.Batch(cmds...)
}

func (f *Form) View() string {
	result := f.title

	for i, input := range f.inputs {
		result += fmt.Sprintf(`

%s
%s`,
			inputStyle.Render(f.names[i]),
			input.View(),
		)
	}
	result += fmt.Sprintf(`

%s
`, continueStyle.Render("Continue ->"))
	return result + "\n"
}

func NewForm(title string, fields []*FormInput) *Form {
	inputs := make([]textinput.Model, len(fields))
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.Name
		inputs[i] = textinput.New()
		inputs[i].Placeholder = ""
		inputs[i].Prompt = f.Prompt
	}
	inputs[0].Focus()
	return &Form{
		title:   title,
		inputs:  inputs,
		names:   names,
		focused: 0,
		err:     nil,
	}
}
