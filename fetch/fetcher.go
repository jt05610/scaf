package fetch

import "os/exec"

type Fetcher interface {
	Fetch(args ...interface{}) error
}

func NPMInstall(workDir string) error {
	cmd := exec.Command("npm", "install")
	cmd.Dir = workDir
	return cmd.Run()
}
