package zpl

import (
	"fmt"
	"io"
	"strings"
)

type Orientation string

const (
	OrientationNormal   Orientation = "N"
	OrientationRotated  Orientation = "R"
	OrientationInverted Orientation = "I"
	OrientationBottomUp Orientation = "B"
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

func (l *Label) SetQuantity(q int) *Label {
	l.builder.WriteString(fmt.Sprintf("^PQ%d\n", q))
	return l
}

func (l *Label) SetPrintWidth(w int) *Label {
	l.builder.WriteString(fmt.Sprintf("^PW%d\n", w))
	return l
}

func (l *Label) SetLabelLength(len int) *Label {
	l.builder.WriteString(fmt.Sprintf("^LL%d\n", len))
	return l
}

func (l *Label) SetPrintRate(rate int) *Label {
	l.builder.WriteString(fmt.Sprintf("^PR%d\n", rate))
	return l
}
