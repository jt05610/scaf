package tree

import (
	"os"
	"path/filepath"
)

type Tree map[string][]string

func (t *Tree) Render(parent string) error {
	for dir, files := range *t {
		p := filepath.Join(parent, dir)
		err := os.MkdirAll(p, 0755)
		if err != nil {
			return err
		}
		for _, file := range files {
			fp := filepath.Join(p, file)
			_, err := os.OpenFile(fp, os.O_RDONLY|os.O_CREATE, 0666)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func New(name string, blinky bool) *Tree {
	if blinky {
		return &Tree{
			".":                          {name + ".yaml", "Caddyfile"},
			"api":                        {name + ".graphql"},
			"ui":                         {},
			"modules":                    {},
			"modules/blinky":             {"blinky.yaml"},
			"modules/blinky/firmware":    {},
			"modules/blinky/service":     {},
			"modules/blinky/service/api": {},
			"modules/blinky/ui":          {},
			"external":                   {},
		}

	} else {
		return &Tree{
			".":       {name + ".yaml", "Caddyfile"},
			"api":     {name + ".graphql"},
			"modules": {},
		}
	}
}
