package fragments

import "strings"

type LineWriter struct {
	lines []string
}

func writer() *LineWriter {
	return &LineWriter{}
}

func (w *LineWriter) s(str ...string) *LineWriter {
	w.lines = append(w.lines, strings.Join(str, ""))
	return w
}

func (w *LineWriter) a(strs []string) *LineWriter {
	w.lines = append(w.lines, strs...)
	return w
}

func (w *LineWriter) w(w2 *LineWriter) *LineWriter {
	w.lines = append(w.lines, w2.lines...)
	return w
}

func (w *LineWriter) ident(f func(lineWriter *LineWriter)) *LineWriter {
	f(w)
	return w
}

func (w *LineWriter) String() string {
	return strings.Join(w.lines, "\n")
}

func (w *LineWriter) Lines() []string {
	return w.lines
}
