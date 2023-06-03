package lang

import (
	"bufio"
	"github.com/google/uuid"
	"github.com/jt05610/scaf/context"
	"log"
	"os"
	"os/exec"
	"sync"
	"text/template"
)

type Command struct {
	Shell  string
	Script string
}

func NewCommand(shell string, script string) *Command {
	return &Command{
		Shell:  shell,
		Script: script,
	}
}

func (c *Command) Execute() error {
	cmd := exec.Command(c.Shell, "-c", c.Script)
	err := cmd.Run()
	if err != nil {
		log.Println("Error executing command: ", err)
	}
	return err
}

func installName() string {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return id.String()
}

type Installable struct {
	ParDir string
}

func Install(ctx context.Context, t *template.Template, i *Installable) error {
	path := "install.sh"
	df, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = os.Remove(path)
		if err != nil {
			panic(err)
		}
	}()
	err = t.Execute(df, i)
	if err != nil {
		panic(err)
	}
	err = df.Close()
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("bash", path)
	stdout, err := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(stdout)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for scanner.Scan() {
			ctx.Logger.Info(scanner.Text())
		}
	}()
	go func() {
		defer wg.Done()
		err = cmd.Wait()
	}()
	wg.Wait()
	return err
}
