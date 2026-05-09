package zpl

import (
	"encoding/hex"
	"fmt"
	"image"
	"strings"
)

func (l *Label) AddGraphicField(x, y int, bytesPerRow int, data []byte) *Label {
	totalBytes := len(data)
	dataBytes := len(data)
	dataStr := strings.ToUpper(hex.EncodeToString(data))

	cmd := fmt.Sprintf("^FO%d,%d^GFA,%d,%d,%d,%s^FS", x, y, totalBytes, dataBytes, bytesPerRow, dataStr)
	l.builder.WriteString(cmd + "\n")
	return l
}

func (l *Label) AddImage(x, y int, img image.Image) *Label {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	bytesPerRow := (width + 7) / 8
	data := make([]byte, bytesPerRow*height)

	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			c := img.At(bounds.Min.X+px, bounds.Min.Y+py)
			r, g, b, _ := c.RGBA()
			lum := (299*uint64(r) + 587*uint64(g) + 114*uint64(b)) / 1000
			if lum < 32768 {
				byteIdx := py*bytesPerRow + px/8
				bitIdx := 7 - (px % 8)
				data[byteIdx] |= (1 << uint(bitIdx))
			}
		}
	}

	return l.AddGraphicField(x, y, bytesPerRow, data)
}
