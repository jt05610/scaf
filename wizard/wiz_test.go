package wizard_test

import (
	"errors"
	"github.com/jt0610/scaf/context"
	"github.com/jt0610/scaf/wizard"
	"github.com/jt0610/scaf/zap"
	"testing"
)

type Person struct {
	Name   string
	Age    int
	IsCool bool
	Hobby  string
}

func (p *Person) Fields() []*wizard.Field {
	return []*wizard.Field{
		{Name: "Name", Type: wizard.String, Default: "Default Name"},
		{Name: "Age", Type: wizard.Int, Default: 25},
		{Name: "IsCool", Type: wizard.Bool, Default: false},
		{Name: "Hobby", Type: wizard.Enum, Default: "Swimming", Options: []string{"Swimming", "Running", "Gaming"}},
	}
}

func (p *Person) Set(field string, value interface{}) error {
	switch field {
	case "Name":
		p.Name, _ = value.(string)
	case "Age":
		p.Age, _ = value.(int)
	case "IsCool":
		p.IsCool, _ = value.(bool)
	case "Hobby":
		p.Hobby, _ = value.(string)
	default:
		return errors.New("invalid field")
	}
	return nil
}

func TestWizard(t *testing.T) {
	p := &Person{}

	w := &wizard.Wizard{}

	// Create a context with a logger
	logger := zap.NewDev(context.Background(), "wizard_test")
	ctx := context.NewContext(logger)
	// Run the wizard with the Fielder
	if err := w.Run(ctx, p); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the results
	if p.Name == "" {
		t.Errorf("Expected Name to be set, got %v", p.Name)
		t.Fail()
	}

	if p.Age == 0 {
		t.Errorf("Expected Age to be set, got %v", p.Age)
		t.Fail()
	}

	if p.Hobby == "" {
		t.Errorf("Expected Hobby to be set, got %v", p.Hobby)
		t.Fail()
	}
}
