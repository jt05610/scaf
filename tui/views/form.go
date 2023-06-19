package views

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jt05610/scaf/core"
	"reflect"
	"strconv"
	"strings"
)

const (
	hotPink = lipgloss.Color("#FF06B7")
)

var (
	inputStyle = lipgloss.NewStyle().Foreground(hotPink)
)

type Filler func(dest interface{}, inputs []*FormInput)

func makeValidator(kind reflect.Kind) textinput.ValidateFunc {
	switch kind {
	case reflect.String:
		// For strings, you might want to check if it's not empty
		return func(input string) error {
			if input == "" {
				return fmt.Errorf("input cannot be empty")
			}
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(input string) error {
			_, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("input must be a valid integer")
			}
			return nil
		}
	case reflect.Float32, reflect.Float64:
		return func(input string) error {
			_, err := strconv.ParseFloat(input, 64)
			if err != nil {
				return fmt.Errorf("input must be a valid float")
			}
			return nil
		}
	case reflect.Bool:
		return func(input string) error {
			_, err := strconv.ParseBool(input)
			if err != nil {
				return fmt.Errorf("input must be a valid boolean")
			}
			return nil
		}
	default:
		return func(input string) error {
			return nil
		}
	}
}

func GenerateFormInputs(i interface{}) *InputSet {
	var t reflect.Type
	var v reflect.Value

	t = reflect.TypeOf(i)
	v = reflect.ValueOf(i)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	result := make([]*FormInput, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check if the field is exported.
		if field.PkgPath != "" {
			continue
		}

		// Check if the json tag starts with a _, if yes ignore it
		jsonTag := field.Tag.Get("json")
		if strings.HasPrefix(jsonTag, "_") {
			continue
		}

		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			nestedStruct := v.Field(i).Elem().Interface()
			result = append(result, &FormInput{
				Name:   field.Name,
				Nested: GenerateFormInputs(nestedStruct),
			})
		} else {
			model := textinput.New()
			model.Prompt = ""
			model.Validate = makeValidator(field.Type.Kind())
			result = append(result, &FormInput{
				Model:    model,
				Name:     field.Name,
				Prompt:   "",
				Validate: makeValidator(field.Type.Kind()),
			})
		}
	}

	filler := &InputSet{
		inputs: result,
		fill: func(dest interface{}, inputs []*FormInput) {
			destVal := reflect.ValueOf(dest).Elem()
			for _, input := range inputs {
				if input.Nested != nil {
					input.Nested.fill(destVal.FieldByName(input.Name).Interface(), input.Nested.inputs)
					continue
				}
				fieldVal := destVal.FieldByName(input.Name)
				switch fieldVal.Kind() {
				case reflect.String:
					input.Model.SetValue(fieldVal.String())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v := strconv.Itoa(int(fieldVal.Int()))
					input.Model.SetValue(v)
				case reflect.Float32, reflect.Float64:
					v := strconv.FormatFloat(fieldVal.Float(), 'f', -1, 64)
					input.Model.SetValue(v)
				case reflect.Bool:
					v := strconv.FormatBool(fieldVal.Bool())
					input.Model.SetValue(v)
				}
			}
		},
	}
	return filler
}

type FormInput struct {
	textinput.Model
	Name     string
	Prompt   string
	Validate textinput.ValidateFunc
	Nested   *InputSet
}

type errMsg error

type InputSet struct {
	inputs  []*FormInput
	current int
	fill    Filler
}

func (s *InputSet) Focus() {
	if s.inputs[s.current].Nested != nil {
		s.inputs[s.current].Nested.Focus()
	} else {
		s.inputs[s.current].Focus()
	}
}

func (s *InputSet) Blur() {
	if s.inputs[s.current].Nested != nil {
		s.inputs[s.current].Nested.Blur()
	} else {
		s.inputs[s.current].Blur()
	}
}

func (s *InputSet) nextInput() {
	if s.inputs[s.current].Nested != nil {
		if s.inputs[s.current].Nested.current < len(s.inputs[s.current].Nested.inputs)-1 {
			s.inputs[s.current].Nested.nextInput()
		} else {
			s.inputs[s.current].Nested.Blur()
			s.current = (s.current + 1) % len(s.inputs)
		}
	} else {
		s.inputs[s.current].Blur()
		s.current = (s.current + 1) % len(s.inputs)
	}

	s.Focus()
}

func (s *InputSet) prevInput() {
	if s.inputs[s.current].Nested != nil {
		if s.inputs[s.current].Nested.current > 0 {
			s.inputs[s.current].Nested.prevInput()
		} else {
			s.inputs[s.current].Nested.Blur()
			s.current = (s.current - 1 + len(s.inputs)) % len(s.inputs)
		}
	} else {
		s.inputs[s.current].Blur()
		s.current = (s.current - 1 + len(s.inputs)) % len(s.inputs)
	}

	s.Focus()
}

func (s *InputSet) viewWithIndent(level int) string {
	indent := strings.Repeat("  ", level)
	res := ""
	for _, input := range s.inputs {
		if input.Nested != nil {
			res += indent + inputStyle.Render(input.Name) + "\n"
			res += input.Nested.viewWithIndent(level + 1)
		} else {
			res += indent + inputStyle.Render(input.Name) + "\n" + indent + input.View() + "\n"
		}
	}
	return res
}

func (s *InputSet) View() string {
	return s.viewWithIndent(0)
}

func (s *InputSet) Fill(sys *core.System) error {
	s.fill(sys, s.inputs)
	return nil
}

func (s *InputSet) Update(msg tea.Msg) (*InputSet, tea.Cmd) {
	var cmds []tea.Cmd
	for _, input := range s.inputs {
		if input.Nested != nil {
			var cmd tea.Cmd
			input.Nested, cmd = input.Nested.Update(msg)
			cmds = append(cmds, cmd)
		}
		input.Update(msg)
	}
	return s, tea.Batch(cmds...)
}

type Form struct {
	title   string
	inputs  *InputSet
	focused *textinput.Model
	err     error
}

func (f *Form) Init() tea.Cmd {
	return textinput.Blink
}

func (f *Form) Update(msg tea.Msg) (*Form, tea.Cmd) {
	cmds := make([]tea.Cmd, len(f.inputs.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			f.inputs.nextInput()
		case tea.KeyUp:
			f.inputs.prevInput()
		case tea.KeyEnter:
			f.inputs.Focus()
		case tea.KeyEsc:
			f.inputs.Blur()
		}
	case errMsg:
		f.err = msg
		return f, nil
	case ShowDetailMsg:
		if err := f.inputs.Fill(msg.System); err != nil {
			f.err = err
		}
		f.title = msg.System.Name
	}
	f.inputs.Update(msg)
	return f, tea.Batch(cmds...)
}

func (f *Form) View() string {
	result := f.title + "\n\n"
	result += f.inputs.View()

	return result + "\n"
}

func findFirst(input *FormInput) *FormInput {
	if input.Nested == nil {
		return input
	}
	return findFirst(input.Nested.inputs[0])
}

func NewForm[T core.Storable](title string, model T) *Form {
	filler := GenerateFormInputs(model)

	return &Form{
		title:   title,
		inputs:  filler,
		focused: &findFirst(filler.inputs[0]).Model,
		err:     nil,
	}
}
