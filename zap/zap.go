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

type Logger struct {
	z *zap.Logger
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.z.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.z.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.z.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.z.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.z.Fatal(msg, fields...)
}

func newLogger(ctx context.Context, parDir, name string, l zapcore.Level) *Logger {
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
	return &Logger{z: z}
}

func newSampler(ctx context.Context, name string, l zapcore.Level, freq time.Duration, initial, skip int) *Logger {
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
	return &Logger{z: z}
}

func NewProd(ctx context.Context, parDir, name string) *Logger {
	l := newLogger(ctx, parDir, name, zapcore.InfoLevel)
	return l
}

func NewDev(ctx context.Context, parDir, name string) *Logger {
	l := newLogger(ctx, parDir, name, zapcore.DebugLevel)
	return l
}

func NewProdSampling(ctx context.Context, name string, freq time.Duration, initial, skip int) *Logger {
	l := newSampler(ctx, name, zapcore.InfoLevel, freq, initial, skip)
	return l
}

func NewDevSampling(ctx context.Context, name string, freq time.Duration, initial, skip int) *Logger {
	l := newSampler(ctx, name, zapcore.DebugLevel, freq, initial, skip)
	return l
}
