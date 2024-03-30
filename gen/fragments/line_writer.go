package fragments

import (
	"fmt"
	"github.com/i512/mapgenie/entities"
	"strings"
)

type LineWriter struct {
	lines []string
}

func writer() *LineWriter {
	return &LineWriter{}
}

func (w *LineWriter) Ln(str ...string) entities.Writer {
	w.lines = append(w.lines, strings.Join(str, ""))
	return w
}

func (w *LineWriter) Lnf(pattern string, args ...any) entities.Writer {
	w.lines = append(w.lines, fmt.Sprintf(pattern, args...))
	return w
}

func (w *LineWriter) Merge(w2 entities.Writer) entities.Writer {
	w.lines = append(w.lines, w2.Lines()...)
	return w
}

func (w *LineWriter) Indent(f func(lineWriter entities.Writer)) entities.Writer {
	f(w)
	return w
}

func (w *LineWriter) String() string {
	return strings.Join(w.lines, "\n")
}

func (w *LineWriter) Lines() []string {
	return w.lines
}

func (w *LineWriter) PreLastLine() entities.Writer {
	return &LineWriter{lines: w.lines[:len(w.lines)-1]}
}

func (w *LineWriter) LastLine() string {
	return w.lines[len(w.lines)-1]
}
