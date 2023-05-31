package solidjs_test

import (
	solidjs2 "github.com/jt0610/scaf/solidjs"
	"testing"
)

func TestHost(t *testing.T) {
	out := "testing/actual/host"
	host := &solidjs2.Host{
		Port: 3000,
		Remotes: []*solidjs2.Remote{
			{
				Name: "test",
				Addr: "localhost",
				Port: 3001,
				Exposes: []*solidjs2.PublicAsset{
					{
						Name: "test",
						Path: "/test",
					},
				},
			},
		},
	}
	renderer := solidjs2.NewHostRenderer(out)

	err := renderer.Render(nil, host)
	if err != nil {
		t.Error(err)
	}

}
