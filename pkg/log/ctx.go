package log

import (
	"context"
	"fmt"
	"io"
)

type ctxKey string

const key = ctxKey("log")

const FoldPrefix = "  "

func Ctx(ctx context.Context, lvl LogLevel, writer io.Writer) context.Context {
	l := &Logger{
		level: lvl,
		io:    writer,
	}

	return context.WithValue(ctx, key, l)
}

func Fold(ctx context.Context, format string, args ...any) context.Context {
	l, ok := ctx.Value(key).(*Logger)
	if !ok {
		fmt.Println("No logger in ctx!")
		return ctx
	}

	newLogger := l.Fold(format, args...)
	return context.WithValue(ctx, key, newLogger)
}

func Debugf(ctx context.Context, format string, args ...any) {
	logf(ctx, Debug, format, args...)
}

func Infof(ctx context.Context, format string, args ...any) {
	logf(ctx, Info, format, args...)
}

func Warnf(ctx context.Context, format string, args ...any) {
	logf(ctx, Warn, format, args...)
}

func Errorf(ctx context.Context, format string, args ...any) {
	logf(ctx, Error, format, args...)
}

func Fatalf(ctx context.Context, format string, args ...any) {
	logf(ctx, Fatal, format, args...)
}

func logf(ctx context.Context, lvl LogLevel, format string, args ...any) {
	l, ok := ctx.Value(key).(*Logger)
	if !ok {
		fmt.Printf("No log in ctx! -- %s\n", fmt.Sprintf(format, args...))
		return
	}

	l.Logf(lvl, format, args...)
}
