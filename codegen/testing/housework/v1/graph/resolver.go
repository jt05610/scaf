// Code generated by github.com/jt05610/scaf, DO NOT EDIT.
// Author: Jonathan Taylor
// Date: 07 Jun 2023

package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"housework/v1"
	"log"
)

type Resolver struct {
	housework housework.HouseworkClient
}

func NewResolver(cert []byte) *Resolver {
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatal("failed to append certs")
	}

	conn, err := grpc.Dial(
		"https://localhost:5001",
		grpc.WithTransportCredentials(
			credentials.NewTLS(
				&tls.Config{
					CurvePreferences: []tls.CurveID{tls.CurveP256},
					MinVersion:       tls.VersionTLS12,
					RootCAs:          certPool,
				},
			),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &Resolver{
		housework: housework.NewHouseworkClient(conn),
	}
}
