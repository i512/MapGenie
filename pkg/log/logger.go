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
	Debug = iota
	Info
	Warn
	Error
)

type Row struct {
	time  time.Time
	level LogLevel
	str   string
}

type Logger struct {
	level        LogLevel
	io           io.Writer
	suppressing  bool
	suppressed   []Row
	unsuppressOn LogLevel
	mux          sync.Mutex
	parent       *Logger
	prefix       string
}

func (l *Logger) Logf(lvl LogLevel, format string, args ...any) {
	l.do(Row{time: time.Now(), level: lvl, str: l.prefix + fmt.Sprintf(format, args...)})
}

func (l *Logger) do(r Row) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.suppressing && r.level == l.unsuppressOn {
		l.unsuppress(nil)
	}

	if l.suppressing {
		l.suppressed = append(l.suppressed, r)
		return
	}

	if r.level >= l.level {
		l.print(r)
	}
}

func (l *Logger) print(r Row) {
	fmt.Fprintf(l.io, "%s%s\n", l.levelShort(r.level), r.str)
}

func (l *Logger) levelShort(lvl LogLevel) string {
	switch lvl {
	case Error:
		return "E"
	case Warn:
		return "W"
	case Info:
		return " "
	case Debug:
		return "d"
	}

	return "?"
}

func (l *Logger) unsuppress(rows []Row) {
	l.suppressing = false
	l.level = Debug

	if l.parent != nil {
		l.parent.unsuppress(append(l.suppressed, rows...))
		l.suppressed = nil
		return
	}

	l.suppressed = append(l.suppressed, rows...)
	sort.SliceStable(l.suppressed, func(i, j int) bool { return l.suppressed[i].time.Before(l.suppressed[j].time) })

	for _, r := range l.suppressed {
		l.print(r)
	}

	l.suppressed = nil
}
