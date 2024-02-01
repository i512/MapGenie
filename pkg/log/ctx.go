package log

import (
	"context"
	"fmt"
	"io"
	"time"
)

type ctxKey string

const key = ctxKey("log")

func Ctx(ctx context.Context, lvl LogLevel, writer io.Writer) context.Context {
	l := &Logger{
		level: lvl,
		io:    writer,
	}

	return context.WithValue(ctx, key, l)
}

func Suppress(ctx context.Context, unsuppressOn LogLevel, f func(context.Context)) context.Context {
	if f == nil {
		f = func(ctx context.Context) {}
	}

	l, ok := ctx.Value(key).(*Logger)
	if !ok {
		fmt.Println("No log in ctx")
		f(ctx)
		return ctx
	}

	newLogger := &Logger{
		level:        l.level,
		io:           l.io,
		suppressing:  true,
		unsuppressOn: unsuppressOn,
		parent:       l,
		prefix:       l.prefix + "  ",
	}

	newCtx := context.WithValue(ctx, key, newLogger)

	f(newCtx)

	return newCtx
}

func Fold(ctx context.Context, format string, args ...any) context.Context {
	l, ok := ctx.Value(key).(*Logger)
	if !ok {
		fmt.Println("No log in ctx!")
		return ctx
	}

	newLogger := &Logger{
		level:        l.level,
		io:           l.io,
		suppressing:  true,
		unsuppressOn: Error,
		parent:       l,
		prefix:       l.prefix + "  ",
		suppressed: []Row{
			{
				time:  time.Now(),
				level: Info,
				str:   l.prefix + fmt.Sprintf(format, args...),
			},
		},
	}

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

func logf(ctx context.Context, lvl LogLevel, format string, args ...any) {
	l, ok := ctx.Value(key).(*Logger)
	if !ok {
		fmt.Printf("No log in ctx! -- %s\n", fmt.Sprintf(format, args...))
		return
	}

	l.Logf(lvl, format, args...)
}
