package wizard_test

import (
	"github.com/jt0610/scaf/context"
	"github.com/jt0610/scaf/wizard"
	"github.com/jt0610/scaf/zap"
	"testing"
)

type Hobby string

type Person struct {
	Name    string `prompt:"What is your name?" default:"John Doe"`
	Age     int    `prompt:"What is your age?" default:"27"`
	IsCool  bool   `prompt:"Are you cool?" default:"yes"`
	Hobby   Hobby  `prompt:"What is your favorite hobby?" options:"coding,labwork"  default:"coding"`
	Friends []*Person
}

func TestWizard(t *testing.T) {

	w := &wizard.Wizard[Person]{}

	logger := zap.NewDev(context.Background(), "wizard_test")
	ctx := context.NewContext(logger)
	p, err := w.Run(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if p.Name == "" {
		t.Errorf("Expected Name to be John Doe, got %v", p.Name)
		t.Fail()
	}

	if p.Age == 0 {
		t.Errorf("Expected Age to be 27, got %v", p.Age)
		t.Fail()
	}
	if p.Hobby != "coding" {
		t.Errorf("Expected Hobby to be coding, got %v", p.Hobby)
		t.Fail()
	}
	if !p.IsCool {
		t.Errorf("Expected IsCool to be true, got %v", p.IsCool)
		t.Fail()
	}
	if len(p.Friends) != 0 {
		t.Errorf("Expected Friends to be empty, got %v", p.Friends)
		t.Fail()
	}
}
