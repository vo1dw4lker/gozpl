package zpl

import (
	"image"
	"image/color"
	"testing"
)

func TestAddGraphicField(t *testing.T) {
	l := NewLabel()

	// 2x2 pixels image, all black
	// 11000000 11000000 -> 0xC0 0xC0
	// Bytes per row = 1
	data := []byte{0xC0, 0xC0}
	l.AddGraphicField(10, 20, 1, data)

	got := l.String()
	want := "^XA\n^FO10,20^GFA,2,2,1,C0C0^FS\n^XZ\n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAddImage(t *testing.T) {
	l := NewLabel()

	// Create a 8x1 image, half black, half white
	// Black pixels at 0,1,2,3 -> bits 7,6,5,4 are 1 -> 11110000 -> 0xF0
	img := image.NewRGBA(image.Rect(0, 0, 8, 1))
	for x := 0; x < 4; x++ {
		img.Set(x, 0, color.Black)
	}
	for x := 4; x < 8; x++ {
		img.Set(x, 0, color.White)
	}

	l.AddImage(0, 0, img)

	got := l.String()
	// Total bytes = 1 (8 bits / 8) * 1 row = 1
	// Data = F0
	want := "^XA\n^FO0,0^GFA,1,1,1,F0^FS\n^XZ\n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
