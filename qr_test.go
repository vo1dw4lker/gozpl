package zpl

import "testing"

func TestLabel_AddQRCode(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l := NewLabel()
		l.AddQRCode(50, 50, "https://google.com")
		expected := "^XA\n^FO50,50^BQN,2,2^FDQA,https://google.com^FS\n^XZ\n"
		if result := l.String(); result != expected {
			t.Errorf("AddQRCode mismatch\nGot: %q\nExp: %q", result, expected)
		}
	})

	t.Run("magnification", func(t *testing.T) {
		l := NewLabel()
		l.AddQRCode(50, 50, "data", WithQRMagnification(5))
		expected := "^XA\n^FO50,50^BQN,2,5^FDQA,data^FS\n^XZ\n"
		if result := l.String(); result != expected {
			t.Errorf("AddQRCode magnification mismatch\nGot: %q\nExp: %q", result, expected)
		}
	})
}
