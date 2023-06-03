package py

import (
	"github.com/jt05610/scaf/lang"
	"path/filepath"
)

func Files(parDir string) []*lang.File {
	return []*lang.File{
		{
			Path:    filepath.Join(parDir, "main.py"),
			TmpFunc: main,
		},
		{
			Path:    filepath.Join(parDir, "test_main.http"),
			TmpFunc: testMain,
		},
	}
}
