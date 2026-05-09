package zpl

import (
	"fmt"
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
	return "^XA\n" + l.builder.String() + "^XZ\n"
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

func (l *Label) Reset() {
	l.builder.Reset()
}

func (l *Label) AddText(x, y int, font string, size int, text string) *Label {
	cmd := "^FO%d,%d^A%s,%d^FD%s^FS"
	l.builder.WriteString(fmt.Sprintf(
		cmd,
		x, y,
		font,
		size,
		text,
	))
	return l
}
