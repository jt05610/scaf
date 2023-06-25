package cmd

import (
	"fmt"
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/core"
	"github.com/jt05610/scaf/testData"
	"github.com/jt05610/scaf/yaml"
	"github.com/jt05610/scaf/zap"
	"os"
	"testing"
)

func TestGen(t *testing.T) {
	parent := "core"
	err := os.RemoveAll(parent)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	s := testData.SCAFSystem(parent)
	logger := zap.NewDev(context.Background(), "testing", "gen_test")
	ctx := context.NewContext(logger)
	gen(ctx, parent, s)
	srv := yaml.NewYAMLService[*core.System]()
	path := "core/system.yaml"
	fp, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = fp.Close()
	}()
	err = srv.Flush(fp, s)
	if err != nil {
		t.Fatal(err)
	}
	modSrv := yaml.NewYAMLService[*core.Module]()
	apiSrv := yaml.NewYAMLService[*core.API]()
	for _, m := range s.Modules {
		path := "core/" + m.Name + "/module.yaml"
		fp, err := os.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		err = modSrv.Flush(fp, m)
		if err != nil {
			t.Fatal(err)
		}
		for _, a := range m.APIs() {
			path := fmt.Sprintf("core/%s/v%d/api.yaml", m.Name, a.Version)
			fp, err := os.Create(path)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = fp.Close()
			}()
			err = apiSrv.Flush(fp, a)
			if err != nil {
				t.Fatal(err)
			}
		}
		_ = fp.Close()
	}
}
