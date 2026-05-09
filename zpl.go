package zpl

import (
	"io"
	"strings"
)

type Label struct {
	builder strings.Builder
}

func NewLabel() *Label {
	return &Label{}
}

func (l *Label) String() string {
	return l.builder.String()
}

func (l *Label) WriteTo(w io.Writer) (int64, error) {
	s := l.String()
	nn, err := w.Write([]byte(s))
	return int64(nn), err
}

func (l *Label) AddRaw(cmd string) *Label {
	l.builder.WriteString(cmd + "\n")
	return l
}
