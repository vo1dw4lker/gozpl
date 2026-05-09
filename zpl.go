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

type textConfig struct {
	fontName      string
	fontHeight    int
	fontWidth     int
	orientation   Orientation
	useFieldBlock bool
	fbWidth       int
	fbMaxLines    int
	fbLineSpacing int
	fbAlignment   string
}

type TextOption func(*textConfig)

func WithFont(name string, height, width int) TextOption {
	return func(c *textConfig) {
		c.fontName = name
		c.fontHeight = height
		c.fontWidth = width
	}
}

func WithTextOrientation(o Orientation) TextOption {
	return func(c *textConfig) {
		c.orientation = o
	}
}

func WithFieldBlock(width, maxLines, lineSpacing int, alignment string) TextOption {
	return func(c *textConfig) {
		c.useFieldBlock = true
		c.fbWidth = width
		c.fbMaxLines = maxLines
		c.fbLineSpacing = lineSpacing
		c.fbAlignment = alignment
	}
}

func (l *Label) AddText(x, y int, text string, opts ...TextOption) *Label {
	cfg := textConfig{
		fontName:    "A",
		fontHeight:  15,
		fontWidth:   15,
		orientation: OrientationNormal,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	cmd := fmt.Sprintf("^FO%d,%d^A%s,%s,%d,%d", x, y, cfg.fontName, string(cfg.orientation), cfg.fontHeight, cfg.fontWidth)
	if cfg.useFieldBlock {
		cmd += fmt.Sprintf("^FB%d,%d,%d,%s", cfg.fbWidth, cfg.fbMaxLines, cfg.fbLineSpacing, cfg.fbAlignment)
	}
	cmd += fmt.Sprintf("^FD%s^FS", text)
	l.builder.WriteString(cmd + "\n")
	return l
}

func (l *Label) AddCode128(x, y int, height int, data string) *Label {
	cmd := "^FO%d,%d^BCN,%d,N,N,Y,A^FD%s^FS\n"
	l.builder.WriteString(fmt.Sprintf(
		cmd,
		x, y,
		height,
		data,
	))
	return l
}
