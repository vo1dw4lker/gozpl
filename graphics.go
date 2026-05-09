package zpl

import "fmt"

type boxConfig struct {
	thickness int
	color     string // "B" or "W"
	rounding  int    // 0-8
}

type BoxOption func(*boxConfig)

func WithBoxThickness(t int) BoxOption {
	return func(c *boxConfig) {
		c.thickness = t
	}
}

func WithBoxRounding(r int) BoxOption {
	return func(c *boxConfig) {
		c.rounding = r
	}
}

func WithBoxColor(color string) BoxOption {
	return func(c *boxConfig) {
		c.color = color
	}
}

func (l *Label) AddGraphicBox(x, y, width, height int, opts ...BoxOption) *Label {
	cfg := boxConfig{
		thickness: 1,
		color:     "B",
		rounding:  0,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	cmd := fmt.Sprintf("^FO%d,%d^GB%d,%d,%d,%s,%d^FS", x, y, width, height, cfg.thickness, cfg.color, cfg.rounding)
	l.builder.WriteString(cmd + "\n")
	return l
}
