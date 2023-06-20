package schema

import (
	"encoding/json"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/lang"
	"io"
	"strings"
)

type Properties map[string]*Entry

type Entry struct {
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Enum        []string   `json:"enum,omitempty"`
	Default     string     `json:"default,omitempty"`
	Items       *Entry     `json:"items,omitempty"`
	Properties  Properties `json:"properties,omitempty"`
}

type schemer struct {
	Schema Properties `json:"schema"`
	lang   *core.Language
	wr     io.Writer
	seen   map[string]bool
}

func (s *schemer) typeString(t *core.Type) string {
	v, found := s.lang.TypeMap.Map(core.BaseType(t.Name))
	if !found {
		v = "object"
	}
	return v
}

func (s *schemer) VisitType(ctx context.Context, t *core.Type) *Entry {
	if !t.IsPrimitive() {
		if _, seen := s.seen[t.Name]; seen {
			return nil
		}
		s.seen[t.Name] = true
	}

	item := &Entry{
		Type:  s.typeString(t),
		Title: t.Name,
	}
	if !t.IsPrimitive() {
		item.Properties = make(Properties)
		for _, f := range t.Fields {
			entry := s.VisitType(ctx, f.Type)
			if entry != nil {
				entry.Title = f.Name
				entry.Description = f.Description
				if f.Type.IsArray {
					arr := &Entry{
						Type:  "array",
						Items: entry,
					}
					item.Properties[strings.ToLower(f.Name)] = arr
				} else {
					item.Properties[strings.ToLower(f.Name)] = entry
				}
			}
		}
	}
	return item
}

func (s *schemer) VisitModule(ctx context.Context, m *core.Module) error {
	api := m.API["v1"]
	for _, t := range api.Types {
		if t.Query {
			if entry := s.VisitType(ctx, t); entry != nil {
				entry.Title = t.Name
				entry.Description = t.Description
				s.Schema[strings.ToLower(t.Name)] = entry
			}
		}
	}
	return json.NewEncoder(s.wr).Encode(s)
}

func (s *schemer) VisitSystem(ctx context.Context, sys *core.System) error {
	for _, m := range sys.Modules {
		if err := s.VisitModule(ctx, m); err != nil {
			return err
		}
	}
	return nil
}

func NewSchemer(w io.Writer) core.Visitor {
	return &schemer{
		Schema: make(map[string]*Entry),
		wr:     w,
		seen:   make(map[string]bool),
		lang:   lang.TypeScript("ts"),
	}
}
