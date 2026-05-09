package zpl

import "fmt"

type qrConfig struct {
	magnification int // 1-10
}

type QROption func(*qrConfig)

func WithQRMagnification(m int) QROption {
	return func(c *qrConfig) {
		c.magnification = m
	}
}

func (l *Label) AddQRCode(x, y int, data string, opts ...QROption) *Label {
	cfg := qrConfig{
		magnification: 2,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	cmd := fmt.Sprintf("^FO%d,%d^BQ%s,2,%d^FDQA,%s^FS", x, y, string(OrientationNormal), cfg.magnification, data)
	l.builder.WriteString(cmd + "\n")
	return l
}
