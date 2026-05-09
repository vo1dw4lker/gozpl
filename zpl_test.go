package zpl

import "testing"

func TestLabel_AddCode128(t *testing.T) {
	l := NewLabel()
	l.AddCode128(5, 5, 30, "hello world")

	expected := "^XA\n^FO5,5^BCN,30,N,N,Y,A^FDhello world^FS\n^XZ\n"
	result := l.String()

	if result != expected {
		t.Errorf("AddCode128 not matching expected output.\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestLabel_AddText(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l := NewLabel()
		l.AddText(10, 20, "Test Label")

		expected := "^XA\n^FO10,20^AA,N,15,15^FDTest Label^FS\n^XZ\n"
		result := l.String()

		if result != expected {
			t.Errorf("AddText not matching expected output.\nExpected:\n%s\nGot:\n%s", expected, result)
		}
	})

	t.Run("with options", func(t *testing.T) {
		l := NewLabel()
		l.AddText(10, 20, "Test Label", WithFont("0", 30, 20), WithTextOrientation(OrientationInverted), WithFieldBlock(200, 2, 0, "C"))

		expected := "^XA\n^FO10,20^A0,I,30,20^FB200,2,0,C^FDTest Label^FS\n^XZ\n"
		result := l.String()

		if result != expected {
			t.Errorf("AddText with options not matching expected output.\nExpected:\n%s\nGot:\n%s", expected, result)
		}
	})
}
