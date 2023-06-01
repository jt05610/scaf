package context

import (
	"context"
	"go.uber.org/zap"
)

type Context struct {
	context.Context
	Logger *zap.Logger
}

func NewContext(l *zap.Logger) Context {
	return Context{
		Context: context.Background(),
		Logger:  l,
	}
}

func Background() context.Context {
	return context.Background()
}
