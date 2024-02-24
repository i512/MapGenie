package log

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"time"
)

type LogLevel int

const (
	Debug = LogLevel(iota)
	Info
	Warn
	Error
	Fatal
)

type Row struct {
	time  time.Time
	level LogLevel
	str   string
}

type Logger struct {
	mux sync.Mutex
	io  io.Writer

	level  LogLevel
	parent *Logger
	prefix string

	folding  bool
	unfoldOn LogLevel
	folded   []Row
}

func (l *Logger) Logf(lvl LogLevel, format string, args ...any) {
	l.do(Row{time: time.Now(), level: lvl, str: l.prefix + fmt.Sprintf(format, args...)})
}

func (l *Logger) do(r Row) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.folding && r.level >= l.unfoldOn {
		l.unfold(nil)
	}

	if l.folding {
		l.folded = append(l.folded, r)
		return
	}

	if r.level >= l.level || !l.folding {
		l.print(r)
	}

	if r.level == Fatal {
		panic("fatalf called: " + r.str)
	}
}

func (l *Logger) print(r Row) {
	_, err := fmt.Fprintf(l.io, "%s\n", l.levelColorize(r.level, r.str))
	if err != nil {
		panic("failed to print log: " + err.Error())
	}
}

func (l *Logger) levelColorize(lvl LogLevel, s string) string {
	switch lvl {
	case Fatal:
		return Color(FColor, s)
	case Error:
		return Color(EColor, s)
	case Warn:
		return Color(WColor, s)
	case Info:
		return Color(IColor, s)
	case Debug:
		return Color(DColor, s)
	}

	return s
}

func (l *Logger) Fold(format string, args ...any) *Logger {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.level != Debug {
		l.folding = true // new fold started, stop showing logs
	}

	newLogger := &Logger{
		level:    l.level,
		io:       l.io,
		folding:  l.folding,
		unfoldOn: Error,
		parent:   l,
		prefix:   l.prefix,
	}

	newLogger.Logf(Info, format, args...)

	newLogger.prefix = FoldPrefix + newLogger.prefix

	return newLogger
}

func (l *Logger) unfold(rows []Row) {
	l.folding = false

	if l.parent != nil {
		l.parent.unfold(append(l.folded, rows...))
		l.folded = nil
		return
	}

	l.folded = append(l.folded, rows...)
	sort.SliceStable(l.folded, func(i, j int) bool { return l.folded[i].time.Before(l.folded[j].time) })

	for _, r := range l.folded {
		l.print(r)
	}

	l.folded = nil
}
