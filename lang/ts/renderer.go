package ts

import (
	"github.com/jt05610/scaf/lang"
	"path/filepath"
)

func Files(parDir string) []*lang.File {
	return []*lang.File{
		{
			Path:    filepath.Join(parDir, "package.json"),
			TmpFunc: PackageJson,
		},
		{
			Path:    filepath.Join(parDir, "tsconfig.json"),
			TmpFunc: tsconfigJson,
		},
		{
			Path:    filepath.Join(parDir, "index.html"),
			TmpFunc: indexHtml,
		},
		{
			Path:    filepath.Join(parDir, "vite.config.ts"),
			TmpFunc: viteConfigTs,
		},
		{
			Path:    filepath.Join(parDir, "src", "index.tsx"),
			TmpFunc: indexTsx,
		},
		{
			Path:    filepath.Join(parDir, "src", "app.tsx"),
			TmpFunc: appTsx,
		},
	}
}
