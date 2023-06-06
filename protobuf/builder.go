package protobuf

import (
	"fmt"
	"github.com/jt05610/scaf/builder"
	"github.com/jt05610/scaf/core"
	"os"
	"os/exec"
	"path/filepath"
)

var apiDir = builder.NewDir(func(m *core.Module) string { return "v1" }, []*builder.File{
	{
		Name:     func(m *core.Module) string { return fmt.Sprintf("%s.proto", m.Name) },
		Template: protoTemplate,
	},
})

func getRPC(m *core.Module) *exec.Cmd {
	cmd := exec.Command("go", "get", "-u", "google.golang.org/grpc")
	cmd.Dir = m.Name
	return cmd
}
func getProtocGenRPC(m *core.Module) *exec.Cmd {
	cmd := exec.Command("go", "install", "google.golang.org/protobuf/cmd/protoc-gen-go@v1.28")
	cmd.Dir = m.Name
	return cmd
}

func getProtoc(m *core.Module) *exec.Cmd {
	cmd := exec.Command("go", "install", "google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2")
	cmd.Dir = m.Name
	return cmd
}

func runProtoc(m *core.Module) *exec.Cmd {
	vars := os.Environ()
	pth := filepath.Join(os.Getenv("GOPATH"), "bin")
	vars = append(vars, "PATH="+os.Getenv("PATH")+":"+pth)
	cmd := exec.Command(
		"protoc",
		"--go_out=paths=source_relative:.",
		"--go-grpc_out=paths=source_relative:.",
		fmt.Sprintf("v1/%s.proto", m.Name),
	)
	cmd.Dir = m.Name
	cmd.Env = append(cmd.Env, vars...)
	return cmd
}

func NewBuilder() core.Visitor {
	return builder.NewBuilder(
		builder.NewDirBuilder(apiDir),
		builder.NewRunner(builder.GoModInit),
		builder.NewRunner(builder.GoModTidy),
		builder.NewRunner(getProtoc),
		builder.NewRunner(getRPC),
		builder.NewRunner(getProtocGenRPC),
		builder.NewRunner(runProtoc),
		builder.NewRunner(builder.GoFmt),
	)
}
