package zpl

import "fmt"

type barcodeConfig struct {
	height      int
	orientation Orientation
	printText   bool
	textAbove   bool
}

type BarcodeOption func(*barcodeConfig)

func WithBarcodeHeight(h int) BarcodeOption {
	return func(c *barcodeConfig) {
		c.height = h
	}
}

func WithBarcodeOrientation(o Orientation) BarcodeOption {
	return func(c *barcodeConfig) {
		c.orientation = o
	}
}

func WithBarcodeText(print bool, above bool) BarcodeOption {
	return func(c *barcodeConfig) {
		c.printText = print
		c.textAbove = above
	}
}

func (l *Label) AddCode128(x, y int, data string, opts ...BarcodeOption) *Label {
	cfg := barcodeConfig{
		height:      50,
		orientation: OrientationNormal,
		printText:   true,
		textAbove:   false,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	pText := "N"
	if cfg.printText {
		pText = "Y"
	}
	tAbove := "N"
	if cfg.textAbove {
		tAbove = "Y"
	}
	cmd := fmt.Sprintf("^FO%d,%d^BC%s,%d,%s,%s,N,N^FD%s^FS", x, y, string(cfg.orientation), cfg.height, pText, tAbove, data)
	l.builder.WriteString(cmd + "\n")
	return l
}
