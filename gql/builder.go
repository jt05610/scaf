package gql

import (
	"github.com/jt05610/scaf/builder"
	"github.com/jt05610/scaf/core"
)

func NewBuilder() core.Visitor {
	return builder.NewBuilder()
}
