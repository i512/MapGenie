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
	Fatal
)

type Row struct {
	time  time.Time
	level LogLevel
	str   string
}

type Logger struct {
	level    LogLevel
	io       io.Writer
	folding  bool
	folded   []Row
	unfoldOn LogLevel
	mux      sync.Mutex
	parent   *Logger
	prefix   string
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

	if r.level >= l.level {
		l.print(r)
	}

	if r.level == Fatal {
		panic("fatalf called: " + r.str)
	}
}

func (l *Logger) print(r Row) {
	_, err := fmt.Fprintf(l.io, "%s%s\n", l.levelShort(r.level), r.str)
	if err != nil {
		panic("failed to print log: " + err.Error())
	}
}

func (l *Logger) levelShort(lvl LogLevel) string {
	switch lvl {
	case Fatal:
		return "F"
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

func (l *Logger) unfold(rows []Row) {
	l.folding = false
	l.level = Debug

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
