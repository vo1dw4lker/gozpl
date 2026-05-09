package zpl

import (
	"bytes"
	"strings"
	"testing"
)

func TestLabel_Settings(t *testing.T) {
	t.Run("quantity", func(t *testing.T) {
		l := NewLabel()
		l.SetQuantity(5)
		expected := "^XA\n^PQ5\n^XZ\n"
		if result := l.String(); result != expected {
			t.Errorf("SetQuantity mismatch\nGot: %q\nExp: %q", result, expected)
		}
	})

	t.Run("print width", func(t *testing.T) {
		l := NewLabel()
		l.SetPrintWidth(800)
		expected := "^XA\n^PW800\n^XZ\n"
		if result := l.String(); result != expected {
			t.Errorf("SetPrintWidth mismatch\nGot: %q\nExp: %q", result, expected)
		}
	})

	t.Run("label length", func(t *testing.T) {
		l := NewLabel()
		l.SetLabelLength(1200)
		expected := "^XA\n^LL1200\n^XZ\n"
		if result := l.String(); result != expected {
			t.Errorf("SetLabelLength mismatch\nGot: %q\nExp: %q", result, expected)
		}
	})

	t.Run("print rate", func(t *testing.T) {
		l := NewLabel()
		l.SetPrintRate(4)
		expected := "^XA\n^PR4\n^XZ\n"
		if result := l.String(); result != expected {
			t.Errorf("SetPrintRate mismatch\nGot: %q\nExp: %q", result, expected)
		}
	})

	t.Run("orientations", func(t *testing.T) {
		orientations := []Orientation{OrientationNormal, OrientationRotated, OrientationInverted, OrientationBottomUp}
		for _, o := range orientations {
			l := NewLabel()
			l.AddText(0, 0, "test", WithTextOrientation(o))
			expected := "^XA\n^FO0,0^AA," + string(o) + ",15,15^FDtest^FS\n^XZ\n"
			if result := l.String(); result != expected {
				t.Errorf("Orientation %s mismatch\nGot: %q\nExp: %q", o, result, expected)
			}
		}
	})
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

func TestLabel_Utilities(t *testing.T) {
	t.Run("Reset", func(t *testing.T) {
		l := NewLabel()
		l.AddText(0, 0, "test")
		l.Reset()
		expected := "^XA\n^XZ\n"
		if result := l.String(); result != expected {
			t.Errorf("Reset mismatch\nGot: %q\nExp: %q", result, expected)
		}
	})

	t.Run("AddRaw", func(t *testing.T) {
		l := NewLabel()
		l.AddRaw("^FO10,10^FDRAW^FS")
		expected := "^XA\n^FO10,10^FDRAW^FS\n^XZ\n"
		if result := l.String(); result != expected {
			t.Errorf("AddRaw mismatch\nGot: %q\nExp: %q", result, expected)
		}
	})

	t.Run("WriteTo", func(t *testing.T) {
		l := NewLabel()
		l.AddText(0, 0, "test")
		var buf bytes.Buffer
		n, err := l.WriteTo(&buf)
		if err != nil {
			t.Fatalf("WriteTo failed: %v", err)
		}
		expected := l.String()
		if int64(buf.Len()) != n {
			t.Errorf("WriteTo length mismatch: got %d, want %d", buf.Len(), n)
		}
		if buf.String() != expected {
			t.Errorf("WriteTo content mismatch\nGot: %q\nExp: %q", buf.String(), expected)
		}
	})
}

func TestLabel_Complex(t *testing.T) {
	l := NewLabel()
	l.SetPrintWidth(400).
		SetLabelLength(400).
		AddText(50, 50, "HEADER", WithFont("0", 40, 40)).
		AddGraphicBox(50, 100, 300, 5).
		AddCode128(50, 150, "12345678", WithBarcodeHeight(60)).
		AddQRCode(50, 250, "https://example.com", WithQRMagnification(4))

	result := l.String()

	// Check if all parts are present
	parts := []string{
		"^XA",
		"^PW400",
		"^LL400",
		"^FO50,50^A0,N,40,40^FDHEADER^FS",
		"^FO50,100^GB300,5,1,B,0^FS",
		"^FO50,150^BCN,60,Y,N,N,N^FD12345678^FS",
		"^FO50,250^BQN,2,4^FDQA,https://example.com^FS",
		"^XZ",
	}

	for _, p := range parts {
		if !strings.Contains(result, p) {
			t.Errorf("Complex label missing part: %s", p)
		}
	}
}
