package zap

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

func newLogger(ctx context.Context, parDir, name string, l zapcore.Level) *zap.Logger {
	par := filepath.Join(parDir, "logs")
	err := os.MkdirAll(par, 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	z := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.AddSync(
				&lumberjack.Logger{
					Filename:   filepath.Join(par, name+".log"),
					MaxSize:    100,
					MaxAge:     7,
					MaxBackups: 5,
					LocalTime:  true,
					Compress:   true,
				}),
			l,
		),
	)

	z = z.WithOptions(
		zap.WrapCore(
			func(core zapcore.Core) zapcore.Core {
				ucEncoder := encoderCfg
				ucEncoder.EncodeLevel = zapcore.CapitalLevelEncoder
				return zapcore.NewTee(
					core,
					zapcore.NewCore(
						zapcore.NewConsoleEncoder(ucEncoder),
						zapcore.Lock(os.Stdout),
						zap.InfoLevel,
					),
				)
			}))
	go func() {
		<-ctx.Done()
		err := z.Sync()
		if err != nil {
			panic(err)
		}
	}()
	return z
}

func newSampler(ctx context.Context, name string, l zapcore.Level, freq time.Duration, initial, skip int) *zap.Logger {
	z := zap.New(
		zapcore.NewSamplerWithOptions(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderCfg),
				zapcore.AddSync(
					&lumberjack.Logger{
						Filename:   name,
						MaxSize:    100,
						MaxAge:     7,
						MaxBackups: 5,
						LocalTime:  true,
						Compress:   true,
					}),
				l,
			), freq, initial, skip,
		),
	)

	z = z.WithOptions(
		zap.WrapCore(
			func(core zapcore.Core) zapcore.Core {
				ucEncoder := encoderCfg
				ucEncoder.EncodeLevel = zapcore.CapitalLevelEncoder
				return zapcore.NewTee(
					core,
					zapcore.NewCore(
						zapcore.NewConsoleEncoder(ucEncoder),
						zapcore.Lock(os.Stdout),
						zapcore.DebugLevel,
					),
				)
			}))
	go func() {
		<-ctx.Done()
		err := z.Sync()
		if err != nil {
			panic(err)
		}
	}()
	return z
}

func NewProd(ctx context.Context, parDir, name string) *zap.Logger {
	l := newLogger(ctx, parDir, name, zapcore.InfoLevel)
	return l
}

func NewDev(ctx context.Context, parDir, name string) *zap.Logger {
	l := newLogger(ctx, parDir, name, zapcore.DebugLevel)
	return l
}

func NewProdSampling(ctx context.Context, name string, freq time.Duration, initial, skip int) *zap.Logger {
	l := newSampler(ctx, name, zapcore.InfoLevel, freq, initial, skip)
	return l
}

func NewDevSampling(ctx context.Context, name string, freq time.Duration, initial, skip int) *zap.Logger {
	l := newSampler(ctx, name, zapcore.DebugLevel, freq, initial, skip)
	return l
}
