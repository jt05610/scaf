package couch_test

import (
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/couch"
	"testing"
)

func TestCRUD(t *testing.T) {
	c := couch.NewCouch[*core.System]("http://scaf:fold@localhost:5984/test-systems")
	s := core.NewSystem("test", "JT", "05 Jul 2023")
	if err := c.Create(s); err != nil {
		t.Fatal(err)
	}
	id := s.ID
	s.PortMap.GQL = 1234
	err := c.Update(s)
	if err != nil {
		t.Fatal(err)
	}
	d, err := c.Details(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if d.PortMap.GQL != 1234 {
		t.Fatal("expected port 1234")
	}
	ss, err := c.List()

	if err != nil {
		t.Fatal(err)
	}
	for _, s := range ss {
		if s.ID == id {
			if s.PortMap.GQL != 1234 {
				t.Fatal("expected port 1234")
			}
			if s.Rev[0] != '2' {
				t.Fatal("expected revision 2")
			}
		}
		if err := c.Delete(s); err != nil {
			t.Fatal(err)
		}
	}
}
