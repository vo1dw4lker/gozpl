package zpl

import "fmt"

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
