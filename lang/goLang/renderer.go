package goLang

import (
	"github.com/jt05610/scaf/lang"
	"os"
	"path/filepath"
)

func Files(parDir string) []*lang.File {
	d := filepath.Join(parDir, "cmd")
	err := os.MkdirAll(d, 0755)
	if err != nil {
		panic(err)
	}
	return []*lang.File{
		{
			Path:    filepath.Join(parDir, "main.go"),
			TmpFunc: main,
		},
		{
			Path:    filepath.Join(parDir, "go.mod"),
			TmpFunc: goMod,
		},
		{
			Path:    filepath.Join(parDir, "cmd/root.go"),
			TmpFunc: root,
		},
	}
}
