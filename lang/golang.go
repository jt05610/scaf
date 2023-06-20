package lang

import (
	"embed"
	"github.com/jt05610/scaf/core"
)

//go:embed template/go
var goTpl embed.FS

var goService = `
		srv := grpc.NewServer()
		rpc := server.Service()
		{{.Name}}.Register{{pascal .Name}}Server(srv, rpc)

		certData, err := secrets.ReadFile(".secrets/{{.Name}}.local+3.pem")
		if err != nil {
			log.Fatal(err)
		}

		keyData, err := secrets.ReadFile(".secrets/{{.Name}}.local+3-key.pem")
		if err != nil {
			log.Fatal(err)
		}

		cert, err := tls.X509KeyPair(certData, keyData)
		if err != nil {
			log.Fatal(err)
		}

		listener, err := net.Listen("tcp", srvAddr)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Listening on %s\n", srvAddr)

		log.Fatal(
			srv.Serve(
				tls.NewListener(
					listener,
					&tls.Config{
						Certificates:     []tls.Certificate{cert},
						CurvePreferences: []tls.CurveID{tls.CurveP256},
						MinVersion:       tls.VersionTLS12,
					},
				),
			),
		)
`
var goTypes = &core.TypeMap{
	Int:    "int64",
	Float:  "float64",
	String: "string",
	Bool:   "bool",
}

var goScripts = &core.Scripts{
	Init: `
mkcert {{.Name}}.local localhost 127.0.0.1 ::1
mkdir ./cmd/.secrets
mv {{.Name}}.local+3.pem {{.Name}}.local+3-key.pem ./cmd/.secrets
go mod tidy
`,
	Gen: `
go generate ./...
go fmt ./...
`,
	Build: `
go build 
go install
`,
}

var goSysScripts = &core.Scripts{
	Init: `
mkcert {{.Name}}.bot localhost 127.0.0.1 ::1
mkdir cmd/.secrets
mv {{.Name}}.bot+3.pem {{.Name}}.bot+3-key.pem cmd/.secrets
go mod tidy
`,
	Gen: `
go generate ./...
go fmt ./...
`,
	Build: `
go build 
go install
`,
}

func Go(parent string) *core.Language {
	l := core.CreateLanguage(
		"go",
		parent,
		goSysScripts,
		goScripts,
		&goTpl,
		goTypes,
		"[]%s",
	)
	l.Service = goService
	return l
}
