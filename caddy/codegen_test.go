package caddy_test

import (
	"github.com/jt0610/scaf/caddy"
	"os"
	"testing"
)

func TestRenderer_Render(t *testing.T) {
	err := os.Remove("testResult")
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	c := caddy.Caddyfile{
		APIs: []*Server{},
	}
}
