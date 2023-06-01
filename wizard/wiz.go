package wizard

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/jt0610/scaf/context"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

type FieldType string

const (
	String FieldType = "string"
	Int              = "int"
	Bool             = "bool"
	Enum             = "enum"
)

type Field struct {
	Name    string
	Type    FieldType
	Default interface{}
	Options []string
}

type Fielder interface {
	Fields() []*Field
	Set(field string, value interface{}) error
}

func (f *Field) stringPrompt() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter a string for %s (%v): ", f.Name, f.Default)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text), nil
}

func (f *Field) intPrompt() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter a integer for %s (%v): ", f.Name, f.Default)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text), nil
}

func (f *Field) boolPrompt() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter a boolean for %s (%v): ", f.Name, f.Default)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text), nil
}
func (f *Field) enumPrompt() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Choose a value for %s:\n", f.Name)
	for i, option := range f.Options {
		fmt.Printf("%d: %s\n", i+1, option)
	}
	fmt.Printf("Enter the number for your choice (%v): ", f.Default)
	text, _ := reader.ReadString('\n')
	if text == "" {
		return f.Default.(string), nil
	}
	choice, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		return "", err
	}
	if choice < 1 || choice > len(f.Options) {
		return "", errors.New("invalid choice")
	}
	return f.Options[choice-1], nil
}

func (f *Field) Prompt(ctx context.Context) (string, error) {
	switch f.Type {
	case String:
		ctx.Logger.Info("prompting for string", zap.String("field", f.Name))
		return f.stringPrompt()
	case Int:
		ctx.Logger.Info("prompting for int", zap.String("field", f.Name))
		return f.intPrompt()
	case Bool:
		ctx.Logger.Info("prompting for bool", zap.String("field", f.Name))
		return f.boolPrompt()
	case Enum:
		ctx.Logger.Info("prompting for enum", zap.String("field", f.Name))
		return f.enumPrompt()
	default:
		return "", errors.New("invalid field type")
	}
}

type Wizard struct{}

func (w *Wizard) Run(ctx context.Context, t Fielder) error {
	ctx.Logger.Info("starting wizard")
	for _, f := range t.Fields() {
		v, err := f.Prompt(ctx)
		ctx.Logger.Info("got value", zap.String("field", f.Name), zap.String("value", v))
		if err != nil {
			return err
		}
		if v == "" {
			err = t.Set(f.Name, f.Default)
		} else {
			err = t.Set(f.Name, v)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
