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
var goTypes = TypeMap{
	core.Int:    "int64",
	core.Float:  "float64",
	core.String: "string",
	core.Bool:   "bool",
	core.ID:     "string",
}

var goScripts = &Scripts{
	Map: map[string]string{
		"init": `
mkcert {{.Name}}.local localhost 127.0.0.1 ::1
mkdir ./cmd/.secrets
mv {{.Name}}.local+3.pem {{.Name}}.local+3-key.pem ./cmd/.secrets
go mod tidy`,
		"gen": `
go generate ./...
go fmt ./...
`,
		"build": `
go build 
go install
`,
	},
}

var goSysScripts = &Scripts{
	Map: map[string]string{
		"init":    "mkcert {{.Name}}.module.local localhost 127.0.0.1 ::1 && mkdir cmd/.secrets && mv {{.Name}}.module.local+3.pem {{.Name}}.module.local+3-key.pem cmd/.secrets && go mod tidy",
		"gen":     "go generate ./...",
		"format":  "go fmt",
		"build":   "go build",
		"relay":   "go run main.go relay",
		"service": "go run main.go service",
	},
}

func Go(parent string) *Language {
	l := CreateLanguage(
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
