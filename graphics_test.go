package zpl

import "testing"

func TestLabel_AddGraphicBox(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l := NewLabel()
		l.AddGraphicBox(10, 10, 100, 50)
		expected := "^XA\n^FO10,10^GB100,50,1,B,0^FS\n^XZ\n"
		if result := l.String(); result != expected {
			t.Errorf("AddGraphicBox mismatch\nGot: %q\nExp: %q", result, expected)
		}
	})

	t.Run("options", func(t *testing.T) {
		l := NewLabel()
		l.AddGraphicBox(10, 10, 100, 50, WithBoxThickness(3), WithBoxRounding(4), WithBoxColor("W"))
		expected := "^XA\n^FO10,10^GB100,50,3,W,4^FS\n^XZ\n"
		if result := l.String(); result != expected {
			t.Errorf("AddGraphicBox options mismatch\nGot: %q\nExp: %q", result, expected)
		}
	})
}
