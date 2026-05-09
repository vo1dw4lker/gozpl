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
	l := NewLabel()
	l.AddText(5, 5, "A", 30, "hello world")

	expected := "^XA\n^FO5,5^AA,30^FDhello world^FS\n^XZ\n"
	result := l.String()

	if result != expected {
		t.Errorf("AddText not matching expected output.\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}
