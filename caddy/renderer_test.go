package caddy_test

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"github.com/jt05610/scaf/caddy"
	"os"
	"testing"
	"time"
)

func TestRenderer_Render(t *testing.T) {
	opt := &caddy.Options{
		UIPortStart:  3000,
		APIPortStart: 8000,
		PortTimeout:  time.Duration(10) * time.Millisecond,
	}

	c := caddy.NewCaddyfile(opt, "test.bot")

	c.AddServer(&caddy.Server{
		Kind: caddy.UI,
		Addr: "localhost",
	})

	c.AddServer(&caddy.Server{
		Kind: caddy.API,
		Addr: "localhost",
		Path: "/blinky",
	})

	buf := new(bytes.Buffer)
	r := caddy.Renderer("caddy")
	err := r.Flush(buf, c)
	if err != nil {
		t.Fatal(err)
	}
	expectF, err := os.ReadFile("testing/Caddyfile")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectF, buf.Bytes()); diff != "" {
		t.Fatalf("Caddyfile mismatch (-want +got):\n%s", diff)
	}
}
