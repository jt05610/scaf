package solidjs_test

import (
	"github.com/jt05610/scaf/solidjs"
	"testing"
)

func TestHost(t *testing.T) {
	out := "testing/actual/host"
	host := &solidjs.Host{
		Port: 3000,
		Remotes: []*solidjs.Remote{
			{
				Name: "test",
				Addr: "localhost",
				Port: 3001,
				Exposes: []*solidjs.PublicAsset{
					{
						Name: "test",
						Path: "/test",
					},
				},
			},
		},
	}
	renderer := solidjs.NewHostRenderer(out)

	err := renderer.Render(nil, host)
	if err != nil {
		t.Error(err)
	}

}
