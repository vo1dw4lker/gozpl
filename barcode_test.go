package zpl

import "testing"

func TestLabel_AddCode128(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l := NewLabel()
		l.AddCode128(5, 5, "hello world")

		expected := "^XA\n^FO5,5^BCN,50,Y,N,N,N^FDhello world^FS\n^XZ\n"
		result := l.String()

		if result != expected {
			t.Errorf("AddCode128 not matching expected output.\nExpected:\n%s\nGot:\n%s", expected, result)
		}
	})

	t.Run("with options", func(t *testing.T) {
		l := NewLabel()
		l.AddCode128(10, 20, "12345", WithBarcodeHeight(80), WithBarcodeOrientation(OrientationRotated), WithBarcodeText(false, true))

		expected := "^XA\n^FO10,20^BCR,80,N,Y,N,N^FD12345^FS\n^XZ\n"
		result := l.String()

		if result != expected {
			t.Errorf("AddCode128 with options not matching expected output.\nExpected:\n%s\nGot:\n%s", expected, result)
		}
	})
}
